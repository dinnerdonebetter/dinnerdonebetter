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

	"github.com/dinnerdonebetter/backend/internal/analytics"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/search/text/indexing"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

var (
	errRequiredDataIsNil = errors.New("required data is nil")

	nonWebhookEventTypes = []types.ServiceEventType{
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
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.AsyncMessageHandlerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	logger := cfg.Observability.Logging.ProvideLogger()

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error("initializing tracer", err)
	}
	otel.SetTracerProvider(tracerProvider)

	metricsProvider, err := cfg.Observability.Metrics.ProvideMetricsProvider(ctx, logger)
	if err != nil {
		logger.Error("initializing metrics provider", err)
	}

	dataManager, consumerProvider, publisherProvider := requisiteSetup(ctx, logger, tracerProvider, cfg)

	defer dataManager.Close()
	defer publisherProvider.Close()

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
		publisherProvider,
		stopChan,
		errorsChan,
	); err != nil {
		log.Fatal(err)
	}

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
		stopChan <- true
	}()
}

func requisiteSetup(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.AsyncMessageHandlerConfig,
) (
	dataManager database.DataManager,
	consumerProvider messagequeue.ConsumerProvider,
	publisherProvider messagequeue.PublisherProvider,
// analyticsEventReporter analytics.EventReporter,
// outboundEmailsPublisher,
// searchDataIndexPublisher,
// webhookExecutionRequestPublisher messagequeue.Publisher,
// emailer email.Emailer,
// uploadManager uploads.UploadManager,
// dataChangesExecutionTimeHistogram,
// outboundEmailsExecutionTimeHistogram,
// searchIndexRequestsExecutionTimeHistogram,
// userDataAggregationExecutionTimeHistogram,
// webhookExecutionTimestampHistogram metrics.Float64Histogram,
	closeFunc func(),
) {
	// connect to database
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		log.Fatalf("error connecting to database: %v", err)
	}

	cancel()

	consumerProvider, err = msgconfig.ProvideConsumerProvider(ctx, logger, &cfg.Events)
	if err != nil {
		log.Fatalf("error initializing consumer provider: %v", err)
	}

	publisherProvider, err = msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatalf("error initializing publisher provider: %v", err)
	}

	return dataManager, consumerProvider, publisherProvider, func() {
		dataManager.Close()
		publisherProvider.Close()
	}
}

func doTheThing(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *config.AsyncMessageHandlerConfig,
	dataManager database.DataManager,
	consumerProvider messagequeue.ConsumerProvider,
	publisherProvider messagequeue.PublisherProvider,
// analyticsEventReporter analytics.EventReporter,
// outboundEmailsPublisher,
// searchDataIndexPublisher,
// webhookExecutionRequestPublisher messagequeue.Publisher,
// emailer email.Emailer,
// uploadManager uploads.UploadManager,
// dataChangesExecutionTimeHistogram,
// outboundEmailsExecutionTimeHistogram,
// searchIndexRequestsExecutionTimeHistogram,
// userDataAggregationExecutionTimeHistogram,
// webhookExecutionTimestampHistogram metrics.Float64Histogram,
	stopChan chan bool,
	errorsChan chan error,
) error {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("async_message_handler"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	// set up myriad publishers

	//nolint:contextcheck // I actually want to use a whatever context here.
	analyticsEventReporter, err := analyticscfg.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider, metricsProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	outboundEmailsPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.OutboundEmailsTopicName)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emails publisher")
	}

	defer outboundEmailsPublisher.Stop()

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.SearchIndexRequestsTopicName)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring search indexing publisher")
	}

	defer searchDataIndexPublisher.Stop()

	webhookExecutionRequestPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.WebhookExecutionRequestsTopicName)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring webhook execution requests publisher")
	}

	defer webhookExecutionRequestPublisher.Stop()

	// setup emailer

	//nolint:contextcheck // I actually want to use a whatever context here.
	emailer, err := emailcfg.ProvideEmailer(&cfg.Email, logger, tracerProvider, metricsProvider, tracing.BuildTracedHTTPClient())
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emailer")
	}

	// setup uploader

	uploadManager, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Storage, chi.NewRouteParamManager())
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating upload manager")
	}

	// setup message listeners

	dataChangesExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("data_changes_execution_time")
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up dataChanges execution time histogram")
	}

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

	outboundEmailsExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("outbound_emails_execution_time")
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up outboundEmails execution time histogram")
	}

	outboundEmailsConsumer, err := consumerProvider.ProvideConsumer(
		ctx,
		cfg.Queues.OutboundEmailsTopicName,
		buildOutboundEmailsEventHandler(
			logger,
			tracer,
			dataManager,
			emailer,
			analyticsEventReporter,
			outboundEmailsExecutionTimeHistogram,
		),
	)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emails consumer")
	}

	searchIndexRequestsExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("search_index_requests_execution_time")
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up searchIndexRequests execution time histogram")
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

	userDataAggregationExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("user_data_aggregation_execution_time")
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up userDataAggregation execution time histogram")
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

	webhookExecutionTimestampHistogram, err := metricsProvider.NewFloat64Histogram("webhook_requests_execution_time")
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up webhookExecutionRequests execution time histogram")
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
	dataManager database.DataManager,
	emailer email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	executionTimestampHistogram metrics.Float64Histogram,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var emailMessage email.DeliveryRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&emailMessage); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		if err := handleEmailRequest(ctx, logger, tracer, dataManager, emailer, analyticsEventReporter, &emailMessage); err != nil {
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

		var searchIndexRequest indexing.IndexRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&searchIndexRequest); err != nil {
			return fmt.Errorf("decoding JSON body: %w", err)
		}

		// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
		if err := indexing.HandleIndexRequest(ctx, logger, tracerProvider, metricsProvider, searchCfg, dataManager, &searchIndexRequest); err != nil {
			return fmt.Errorf("handling search indexing request: %w", err)
		}

		executionTimestampHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}

func buildUserDataAggregationEventHandler(
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	uploadManager *objectstorage.Uploader,
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
