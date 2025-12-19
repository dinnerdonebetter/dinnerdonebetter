package datachangemessagehandler

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/config"
	identitymock "github.com/dinnerdonebetter/backend/internal/domain/identity/mock"
	webhooksmock "github.com/dinnerdonebetter/backend/internal/domain/webhooks/mock"
	analyticsmock "github.com/dinnerdonebetter/backend/internal/platform/analytics/mock"
	emailmock "github.com/dinnerdonebetter/backend/internal/platform/email/mock"
	encodingmock "github.com/dinnerdonebetter/backend/internal/platform/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	msgqueuemock "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	mockmetrics "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	uploadsmock "github.com/dinnerdonebetter/backend/internal/platform/uploads/mock"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	mealplanningindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:gocritic // I know this returns too many things
func buildTestAsyncDataChangeMessageHandler(t *testing.T) (*AsyncDataChangeMessageHandler, *identitymock.RepositoryMock, *webhooksmock.Repository, *msgqueuemock.ConsumerProvider, *msgqueuemock.PublisherProvider, *analyticsmock.EventReporter, *emailmock.Emailer, *uploadsmock.MockUploadManager, *mockmetrics.MetricsProvider, *encodingmock.EncoderDecoder) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())

	identityRepo := &identitymock.RepositoryMock{}
	webhookRepo := &webhooksmock.Repository{}
	consumerProvider := &msgqueuemock.ConsumerProvider{}
	publisherProvider := &msgqueuemock.PublisherProvider{}
	analyticsEventReporter := &analyticsmock.EventReporter{}
	emailer := &emailmock.Emailer{}
	uploadManager := &uploadsmock.MockUploadManager{}
	metricsProvider := &mockmetrics.MetricsProvider{}
	decoder := &encodingmock.EncoderDecoder{}

	// Create mock indexers with noop implementations for testing
	userDataIndexer := &identityindexing.UserDataIndexer{}
	mealPlanningDataIndexer := &mealplanningindexing.MealPlanningDataIndexer{}

	// Set up mock publishers for the indexers to prevent nil pointer dereferences
	mockPublisher := &msgqueuemock.Publisher{}
	publisherProvider.On("ProvidePublisher", mock.AnythingOfType("string")).Return(mockPublisher, nil).Maybe()

	// Set up mock histograms
	mockHistogram := metrics.NewNoopMetricsProvider()
	noop, _ := mockHistogram.NewFloat64Histogram("test")
	metricsProvider.On("NewFloat64Histogram", mock.AnythingOfType("string"), mock.Anything).Return(noop, nil).Maybe()

	handler := &AsyncDataChangeMessageHandler{
		identityRepo:                         identityRepo,
		webhookRepo:                          webhookRepo,
		consumerProvider:                     consumerProvider,
		analyticsEventReporter:               analyticsEventReporter,
		emailer:                              emailer,
		uploadManager:                        uploadManager,
		decoder:                              decoder,
		userDataIndexer:                      userDataIndexer,
		mealPlanningDataIndexer:              mealPlanningDataIndexer,
		logger:                               logger,
		tracer:                               tracer,
		nonWebhookEventTypes:                 []string{},
		dataChangesExecutionTimeHistogram:    noop,
		outboundEmailsExecutionTimeHistogram: noop,
		webhookExecutionTimestampHistogram:   noop,
		userDataAggregationExecutionTimeHistogram: noop,
		queuesConfig: msgconfig.QueuesConfig{
			SearchIndexRequestsTopicName: "search-index-requests",
		},
		searchDataIndexPublisher:         mockPublisher,
		outboundEmailsPublisher:          mockPublisher,
		webhookExecutionRequestPublisher: mockPublisher,
	}

	return handler, identityRepo, webhookRepo, consumerProvider, publisherProvider, analyticsEventReporter, emailer, uploadManager, metricsProvider, decoder
}

func TestNewAsyncDataChangeMessageHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		cfg := &config.AsyncMessageHandlerConfig{
			Queues: msgconfig.QueuesConfig{
				OutboundEmailsTopicName:           "outbound-emails",
				SearchIndexRequestsTopicName:      "search-index-requests",
				WebhookExecutionRequestsTopicName: "webhook-execution-requests",
			},
		}
		identityRepo := &identitymock.RepositoryMock{}
		webhookRepo := &webhooksmock.Repository{}
		consumerProvider := &msgqueuemock.ConsumerProvider{}
		publisherProvider := &msgqueuemock.PublisherProvider{}
		analyticsEventReporter := &analyticsmock.EventReporter{}
		emailer := &emailmock.Emailer{}
		uploadManager := &uploadsmock.MockUploadManager{}
		metricsProvider := &mockmetrics.MetricsProvider{}
		decoder := &encodingmock.EncoderDecoder{}
		coreDataIndexer := &identityindexing.UserDataIndexer{}
		eatingDataIndexer := &mealplanningindexing.MealPlanningDataIndexer{}

		// Set up metrics expectations
		mockHistogram := metrics.NewNoopMetricsProvider()
		noop, _ := mockHistogram.NewFloat64Histogram("test")
		metricsProvider.On("NewFloat64Histogram", "data_changes_execution_time", mock.Anything).Return(noop, nil)
		metricsProvider.On("NewFloat64Histogram", "outbound_emails_execution_time", mock.Anything).Return(noop, nil)
		metricsProvider.On("NewFloat64Histogram", "search_index_requests_execution_time", mock.Anything).Return(noop, nil)
		metricsProvider.On("NewFloat64Histogram", "user_data_aggregation_execution_time", mock.Anything).Return(noop, nil)
		metricsProvider.On("NewFloat64Histogram", "webhook_requests_execution_time", mock.Anything).Return(noop, nil)

		// Set up publisher expectations
		mockPublisher := &msgqueuemock.Publisher{}
		publisherProvider.On("ProvidePublisher", "outbound-emails").Return(mockPublisher, nil)
		publisherProvider.On("ProvidePublisher", "search-index-requests").Return(mockPublisher, nil)
		publisherProvider.On("ProvidePublisher", "webhook-execution-requests").Return(mockPublisher, nil)

		handler, err := NewAsyncDataChangeMessageHandler(
			ctx,
			logger,
			tracerProvider,
			cfg,
			identityRepo,
			webhookRepo,
			consumerProvider,
			publisherProvider,
			analyticsEventReporter,
			emailer,
			uploadManager,
			metricsProvider,
			decoder,
			coreDataIndexer,
			eatingDataIndexer,
		)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
		assert.Equal(t, identityRepo, handler.identityRepo)
		assert.Equal(t, webhookRepo, handler.webhookRepo)
		assert.Equal(t, consumerProvider, handler.consumerProvider)
		assert.Equal(t, analyticsEventReporter, handler.analyticsEventReporter)
		assert.Equal(t, emailer, handler.emailer)
		assert.Equal(t, uploadManager, handler.uploadManager)
		assert.Equal(t, decoder, handler.decoder)
		assert.Equal(t, coreDataIndexer, handler.userDataIndexer)
		assert.Equal(t, eatingDataIndexer, handler.mealPlanningDataIndexer)

		mock.AssertExpectationsForObjects(t, metricsProvider, publisherProvider)
	})
}

func TestAsyncDataChangeMessageHandler_SetNonWebhookEventTypes(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		eventTypes := []string{"event1", "event2", "event3"}
		handler.SetNonWebhookEventTypes(eventTypes)

		handler.nonWebhookEventTypesHat.RLock()
		assert.Equal(t, eventTypes, handler.nonWebhookEventTypes)
		handler.nonWebhookEventTypesHat.RUnlock()
	})
}
