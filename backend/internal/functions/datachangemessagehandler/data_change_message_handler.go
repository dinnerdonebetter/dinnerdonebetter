package datachangemessagehandler

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
)

const (
	o11yName = "async_data_change_message_handler"
)

var (
	errRequiredDataIsNil = errors.New("required data is nil")
)

type AsyncDataChangeMessageHandler struct {
	searchDataIndexPublisher                  messagequeue.Publisher
	identityRepo                              identity.Repository
	logger                                    logging.Logger
	decoder                                   encoding.ServerEncoderDecoder
	webhookExecutionTimestampHistogram        metrics.Float64Histogram
	userDataAggregationExecutionTimeHistogram metrics.Float64Histogram
	outboundEmailsPublisher                   messagequeue.Publisher
	webhookRepo                               webhooks.Repository
	mealPlanningRepo                          mealplanning.Repository
	consumerProvider                          messagequeue.ConsumerProvider
	tracerProvider                            tracing.TracerProvider
	analyticsEventReporter                    analytics.EventReporter
	dataChangesExecutionTimeHistogram         metrics.Float64Histogram
	webhookExecutionRequestPublisher          messagequeue.Publisher
	emailer                                   email.Emailer
	uploadManager                             uploads.UploadManager
	tracer                                    tracing.Tracer
	outboundEmailsExecutionTimeHistogram      metrics.Float64Histogram
	searchIndexRequestsExecutionTimeHistogram metrics.Float64Histogram
	coreDataIndexer                           *coreindexing.CoreDataIndexer
	eatingDataIndexer                         *eatingindexing.EatingDataIndexer
	queuesConfig                              msgconfig.QueuesConfig
	nonWebhookEventTypes                      []string
	nonWebhookEventTypesHat                   sync.RWMutex
}

func (a *AsyncDataChangeMessageHandler) SetNonWebhookEventTypes(nonWebhookEventTypes []string) {
	a.nonWebhookEventTypesHat.Lock()
	defer a.nonWebhookEventTypesHat.Unlock()
	a.nonWebhookEventTypes = nonWebhookEventTypes
}

func NewAsyncDataChangeMessageHandler(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.AsyncMessageHandlerConfig,
	identityRepo identity.Repository,
	webhookRepo webhooks.Repository,
	mealplanningRepo mealplanning.Repository,
	consumerProvider messagequeue.ConsumerProvider,
	publisherProvider messagequeue.PublisherProvider,
	analyticsEventReporter analytics.EventReporter,
	emailer email.Emailer,
	uploadManager uploads.UploadManager,
	metricsProvider metrics.Provider,
	decoder encoding.ServerEncoderDecoder,
	coreDataIndexer *coreindexing.CoreDataIndexer,
	eatingDataIndexer *eatingindexing.EatingDataIndexer,
) (*AsyncDataChangeMessageHandler, error) {
	dataChangesExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("data_changes_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up dataChanges execution time histogram: %w", err)
	}

	outboundEmailsExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("outbound_emails_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up outboundEmails execution time histogram: %w", err)
	}

	searchIndexRequestsExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("search_index_requests_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up searchIndexRequests execution time histogram: %w", err)
	}

	userDataAggregationExecutionTimeHistogram, err := metricsProvider.NewFloat64Histogram("user_data_aggregation_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up userDataAggregation execution time histogram: %w", err)
	}

	webhookExecutionTimestampHistogram, err := metricsProvider.NewFloat64Histogram("webhook_requests_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up webhookExecutionRequests execution time histogram: %w", err)
	}

	outboundEmailsPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.OutboundEmailsTopicName)
	if err != nil {
		return nil, fmt.Errorf("configuring outbound emails publisher: %w", err)
	}

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.SearchIndexRequestsTopicName)
	if err != nil {
		return nil, fmt.Errorf("configuring search indexing publisher: %w", err)
	}

	webhookExecutionRequestPublisher, err := publisherProvider.ProvidePublisher(cfg.Queues.WebhookExecutionRequestsTopicName)
	if err != nil {
		return nil, fmt.Errorf("configuring webhook execution requests publisher: %w", err)
	}

	return &AsyncDataChangeMessageHandler{
		tracer:                               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                               logging.EnsureLogger(logger).WithName(o11yName),
		nonWebhookEventTypes:                 []string{},
		identityRepo:                         identityRepo,
		webhookRepo:                          webhookRepo,
		mealPlanningRepo:                     mealplanningRepo,
		consumerProvider:                     consumerProvider,
		analyticsEventReporter:               analyticsEventReporter,
		outboundEmailsPublisher:              outboundEmailsPublisher,
		searchDataIndexPublisher:             searchDataIndexPublisher,
		queuesConfig:                         cfg.Queues,
		webhookExecutionRequestPublisher:     webhookExecutionRequestPublisher,
		emailer:                              emailer,
		uploadManager:                        uploadManager,
		dataChangesExecutionTimeHistogram:    dataChangesExecutionTimeHistogram,
		outboundEmailsExecutionTimeHistogram: outboundEmailsExecutionTimeHistogram,
		searchIndexRequestsExecutionTimeHistogram: searchIndexRequestsExecutionTimeHistogram,
		userDataAggregationExecutionTimeHistogram: userDataAggregationExecutionTimeHistogram,
		webhookExecutionTimestampHistogram:        webhookExecutionTimestampHistogram,
		decoder:                                   decoder,
		coreDataIndexer:                           coreDataIndexer,
		eatingDataIndexer:                         eatingDataIndexer,
	}, nil
}

func (a *AsyncDataChangeMessageHandler) ConsumeMessages(
	ctx context.Context,
	stopChan chan bool,
	errorsChan chan error,
) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	// set up myriad publishers

	dataChangesConsumer, err := a.consumerProvider.ProvideConsumer(
		ctx,
		a.queuesConfig.DataChangesTopicName,
		a.DataChangesEventHandler,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, a.logger, span, "configuring data changes consumer")
	}

	outboundEmailsConsumer, err := a.consumerProvider.ProvideConsumer(
		ctx,
		a.queuesConfig.OutboundEmailsTopicName,
		a.OutboundEmailsEventHandler,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, a.logger, span, "configuring outbound emails consumer")
	}

	searchIndexRequestsConsumer, err := a.consumerProvider.ProvideConsumer(
		ctx,
		a.queuesConfig.SearchIndexRequestsTopicName,
		a.SearchIndexRequestsEventHandler,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, a.logger, span, "configuring search index requests consumer")
	}

	webhookExecutionRequestsConsumer, err := a.consumerProvider.ProvideConsumer(
		ctx,
		a.queuesConfig.WebhookExecutionRequestsTopicName,
		a.WebhookExecutionRequestsEventHandler,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, a.logger, span, "configuring webhook execution requests consumer")
	}

	//userDataAggregationConsumer, err := a.consumerProvider.ProvideConsumer(
	//	ctx,
	//	a.queuesConfig.UserDataAggregationTopicName,
	//	a.buildUserDataAggregationEventHandler(
	//		logger,
	//		tracer,
	//		dataManager,
	//		uploadManager,
	//		userDataAggregationExecutionTimeHistogram,
	//	),
	//)
	//if err != nil {
	//	return observability.PrepareAndLogError(err, a.logger, span, "configuring user data aggregation consumer")
	//}

	go dataChangesConsumer.Consume(stopChan, errorsChan)
	go outboundEmailsConsumer.Consume(stopChan, errorsChan)
	go searchIndexRequestsConsumer.Consume(stopChan, errorsChan)
	go webhookExecutionRequestsConsumer.Consume(stopChan, errorsChan)
	//go userDataAggregationConsumer.Consume(stopChan, errorsChan)

	go func() {
		for e := range errorsChan {
			a.logger.Error("consuming message", e)
		}
	}()

	return nil
}
