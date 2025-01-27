package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/lib/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/lib/email/config"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing/chi"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "go.uber.org/automaxprocs"
)

var (
	errRequiredDataIsNil = errors.New("required data is nil")

	nonWebhookEventTypes = []string{
		types.UserSignedUpServiceEventType,
		types.UserArchivedServiceEventType,
		types.TwoFactorSecretVerifiedServiceEventType,
		types.TwoFactorDeactivatedServiceEventType,
		types.TwoFactorSecretChangedServiceEventType,
		types.PasswordResetTokenCreatedEventType,
		types.PasswordResetTokenRedeemedEventType,
		types.PasswordChangedEventType,
		types.EmailAddressChangedEventType,
		types.UsernameChangedEventType,
		types.UserDetailsChangedEventType,
		types.UsernameReminderRequestedEventType,
		types.UserLoggedInServiceEventType,
		types.UserLoggedOutServiceEventType,
		types.UserChangedActiveHouseholdServiceEventType,
		types.UserEmailAddressVerifiedEventType,
		types.UserEmailAddressVerificationEmailRequestedEventType,
		types.HouseholdMemberRemovedServiceEventType,
		types.HouseholdMembershipPermissionsUpdatedServiceEventType,
		types.HouseholdOwnershipTransferredServiceEventType,
		types.OAuth2ClientCreatedServiceEventType,
		types.OAuth2ClientArchivedServiceEventType,
	}
)

func main() {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.AsyncMessageHandlerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	ctx := context.Background()
	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatalf("could not establish observability pillars: %v", err)
	}

	dataManager,
		consumerProvider,
		analyticsEventReporter,
		outboundEmailsPublisher,
		searchDataIndexPublisher,
		webhookExecutionRequestPublisher,
		emailer,
		uploadManager,
		dataChangesExecutionTimeHistogram,
		outboundEmailsExecutionTimeHistogram,
		searchIndexRequestsExecutionTimeHistogram,
		userDataAggregationExecutionTimeHistogram,
		webhookExecutionTimestampHistogram,
		closeFunc,
		err := setupDependencies(ctx, logger, tracerProvider, metricsProvider, cfg)
	if err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	stopChan := make(chan bool)
	errorsChan := make(chan error)

	if err = doTheThing(
		ctx,
		logger,
		tracerProvider,
		metricsProvider,
		cfg,
		dataManager,
		consumerProvider,
		analyticsEventReporter,
		outboundEmailsPublisher,
		searchDataIndexPublisher,
		webhookExecutionRequestPublisher,
		emailer,
		uploadManager,
		dataChangesExecutionTimeHistogram,
		outboundEmailsExecutionTimeHistogram,
		searchIndexRequestsExecutionTimeHistogram,
		userDataAggregationExecutionTimeHistogram,
		webhookExecutionTimestampHistogram,
		stopChan,
		errorsChan,
	); err != nil {
		closeFunc()
		log.Fatal(err)
	}

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
		stopChan <- true
		closeFunc()
	}()
}

//nolint:gocritic // I know there are too many results, I don't see the value in breaking this out.
func setupDependencies(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *config.AsyncMessageHandlerConfig,
) (
	dataManager database.DataManager,
	consumerProvider messagequeue.ConsumerProvider,
	analyticsEventReporter analytics.EventReporter,
	outboundEmailsPublisher,
	searchDataIndexPublisher,
	webhookExecutionRequestPublisher messagequeue.Publisher,
	emailer email.Emailer,
	uploadManager uploads.UploadManager,
	dataChangesExecutionTimeHistogram,
	outboundEmailsExecutionTimeHistogram,
	searchIndexRequestsExecutionTimeHistogram,
	userDataAggregationExecutionTimeHistogram,
	webhookExecutionTimestampHistogram metrics.Float64Histogram,
	closeFunc func(),
	err error,
) {
	// connect to database
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err = postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error connecting to database: %w", err)
	}
	cancel()

	consumerProvider, err = msgconfig.ProvideConsumerProvider(ctx, logger, &cfg.Events)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error initializing consumer provider: %w", err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error initializing publisher provider: %w", err)
	}

	//nolint:contextcheck // I actually want to use a whatever context here.
	analyticsEventReporter, err = analyticscfg.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up customer data collector: %w", err)
	}

	outboundEmailsPublisher, err = publisherProvider.ProvidePublisher(cfg.Queues.OutboundEmailsTopicName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("configuring outbound emails publisher: %w", err)
	}

	searchDataIndexPublisher, err = publisherProvider.ProvidePublisher(cfg.Queues.SearchIndexRequestsTopicName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("configuring search indexing publisher: %w", err)
	}

	webhookExecutionRequestPublisher, err = publisherProvider.ProvidePublisher(cfg.Queues.WebhookExecutionRequestsTopicName)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("configuring webhook execution requests publisher: %w", err)
	}

	// setup emailer

	//nolint:contextcheck // I actually want to use a whatever context here.
	emailer, err = emailcfg.ProvideEmailer(&cfg.Email, logger, tracerProvider, metricsProvider, tracing.BuildTracedHTTPClient())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("configuring outbound emailer: %w", err)
	}

	// setup uploader
	uploadManager, err = objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Storage, chi.NewRouteParamManager())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("creating upload manager: %w", err)
	}

	// setup execution timers
	dataChangesExecutionTimeHistogram, err = metricsProvider.NewFloat64Histogram("data_changes_execution_time")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up dataChanges execution time histogram: %w", err)
	}

	outboundEmailsExecutionTimeHistogram, err = metricsProvider.NewFloat64Histogram("outbound_emails_execution_time")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up outboundEmails execution time histogram: %w", err)
	}

	searchIndexRequestsExecutionTimeHistogram, err = metricsProvider.NewFloat64Histogram("search_index_requests_execution_time")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up searchIndexRequests execution time histogram: %w", err)
	}

	userDataAggregationExecutionTimeHistogram, err = metricsProvider.NewFloat64Histogram("user_data_aggregation_execution_time")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up userDataAggregation execution time histogram: %w", err)
	}

	webhookExecutionTimestampHistogram, err = metricsProvider.NewFloat64Histogram("webhook_requests_execution_time")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("setting up webhookExecutionRequests execution time histogram: %w", err)
	}

	return dataManager,
		consumerProvider,
		analyticsEventReporter,
		outboundEmailsPublisher,
		searchDataIndexPublisher,
		webhookExecutionRequestPublisher,
		emailer,
		uploadManager,
		dataChangesExecutionTimeHistogram,
		outboundEmailsExecutionTimeHistogram,
		searchIndexRequestsExecutionTimeHistogram,
		userDataAggregationExecutionTimeHistogram,
		webhookExecutionTimestampHistogram,
		func() {
			dataManager.Close()
			publisherProvider.Close()
			analyticsEventReporter.Close()
			outboundEmailsPublisher.Stop()
			searchDataIndexPublisher.Stop()
			webhookExecutionRequestPublisher.Stop()
		},
		nil
}

func doTheThing(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *config.AsyncMessageHandlerConfig,
	dataManager database.DataManager,
	consumerProvider messagequeue.ConsumerProvider,
	analyticsEventReporter analytics.EventReporter,
	outboundEmailsPublisher,
	searchDataIndexPublisher,
	webhookExecutionRequestPublisher messagequeue.Publisher,
	emailer email.Emailer,
	uploadManager uploads.UploadManager,
	dataChangesExecutionTimeHistogram,
	outboundEmailsExecutionTimeHistogram,
	searchIndexRequestsExecutionTimeHistogram,
	userDataAggregationExecutionTimeHistogram,
	webhookExecutionTimestampHistogram metrics.Float64Histogram,
	stopChan chan bool,
	errorsChan chan error,
) error {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("async_message_handler"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	// set up myriad publishers

	dataChangesConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.DataChangesTopicName,
		buildDataChangesEventHandler(
			logger,
			tracer,
			dataManager,
			analyticsEventReporter,
			webhookExecutionRequestPublisher,
			outboundEmailsPublisher,
			searchDataIndexPublisher,
			dataChangesExecutionTimeHistogram,
		),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring data changes consumer")
	}

	outboundEmailsConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.OutboundEmailsTopicName,
		buildOutboundEmailsEventHandler(
			logger,
			tracer,
			emailer,
			analyticsEventReporter,
			outboundEmailsExecutionTimeHistogram,
		),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emails consumer")
	}

	searchIndexRequestsConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.SearchIndexRequestsTopicName,
		buildSearchIndexRequestsEventHandler(
			logger,
			tracer,
			tracerProvider,
			dataManager,
			metricsProvider,
			&cfg.Search,
			searchIndexRequestsExecutionTimeHistogram,
		),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring search index requests consumer")
	}

	userDataAggregationConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.UserDataAggregationTopicName,
		buildUserDataAggregationEventHandler(
			logger,
			tracer,
			dataManager,
			uploadManager,
			userDataAggregationExecutionTimeHistogram,
		),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring user data aggregation consumer")
	}

	webhookExecutionRequestsConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.WebhookExecutionRequestsTopicName,
		buildWebhookExecutionRequestsEventHandler(logger, tracer, dataManager, webhookExecutionTimestampHistogram),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring webhook execution requests consumer")
	}

	go dataChangesConsumer.Consume(stopChan, errorsChan)
	go outboundEmailsConsumer.Consume(stopChan, errorsChan)
	go searchIndexRequestsConsumer.Consume(stopChan, errorsChan)
	go userDataAggregationConsumer.Consume(stopChan, errorsChan)
	go webhookExecutionRequestsConsumer.Consume(stopChan, errorsChan)
	go func() {
		for e := range errorsChan {
			logger.Error("consuming message", e)
		}
	}()

	return nil
}

func buildDataChangesEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	analyticsEventReporter analytics.EventReporter,
	webhookExecutionRequestPublisher,
	outboundEmailsPublisher,
	searchDataIndexPublisher messagequeue.Publisher,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var dataChangeMessage types.DataChangeMessage
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&dataChangeMessage); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		handleDataChangeMessage(ctx, logger, tracer, dataManager, analyticsEventReporter, webhookExecutionRequestPublisher, outboundEmailsPublisher, searchDataIndexPublisher, &dataChangeMessage)

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}

func buildOutboundEmailsEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	emailer email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var emailMessage email.OutboundEmailMessage
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&emailMessage); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		if err := handleEmailRequest(ctx, logger, tracer, emailer, analyticsEventReporter, &emailMessage); err != nil {
			return fmt.Errorf("handling outbound email request: %w", err)
		}

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}

func buildSearchIndexRequestsEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	tracerProvider tracing.TracerProvider,
	dataManager database.DataManager,
	metricsProvider metrics.Provider,
	searchCfg *textsearchcfg.Config,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var searchIndexRequest textsearch.IndexRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&searchIndexRequest); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		switch searchIndexRequest.IndexType {
		case textsearch.IndexTypeRecipes,
			textsearch.IndexTypeMeals,
			textsearch.IndexTypeValidIngredients,
			textsearch.IndexTypeValidInstruments,
			textsearch.IndexTypeValidMeasurementUnits,
			textsearch.IndexTypeValidPreparations,
			textsearch.IndexTypeValidIngredientStates,
			textsearch.IndexTypeValidVessels:
			// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
			if err := eatingindexing.HandleIndexRequest(ctx, logger, tracerProvider, metricsProvider, searchCfg, dataManager, &searchIndexRequest); err != nil {
				return fmt.Errorf("handling search indexing request: %w", err)
			}

		case textsearch.IndexTypeUsers:
			// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
			if err := coreindexing.HandleIndexRequest(ctx, logger, tracerProvider, metricsProvider, searchCfg, dataManager, &searchIndexRequest); err != nil {
				return fmt.Errorf("handling search indexing request: %w", err)
			}
		}

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}

func buildUserDataAggregationEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	uploadManager uploads.UploadManager,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var userDataAggregationRequest types.UserDataAggregationRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&userDataAggregationRequest); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		if err := handleUserDataRequest(ctx, logger, tracer, uploadManager, dataManager, &userDataAggregationRequest); err != nil {
			return fmt.Errorf("handling user data aggregation request: %w", err)
		}

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}

func buildWebhookExecutionRequestsEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var webhookExecutionRequest types.WebhookExecutionRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&webhookExecutionRequest); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		if err := handleWebhookExecutionRequest(ctx, logger, tracer, dataManager, &webhookExecutionRequest); err != nil {
			return fmt.Errorf("handling webhook execution request: %w", err)
		}

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}
