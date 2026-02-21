package userdataaggregationhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads"
)

const (
	o11yName = "user_data_aggregation_handler"
)

type UserDataAggregationHandler struct {
	logger                 logging.Logger
	tracer                 tracing.Tracer
	consumerProvider       messagequeue.ConsumerProvider
	dataPrivacyRepo        dataprivacy.Repository
	uploadManager          uploads.UploadManager
	decoder                encoding.ServerEncoderDecoder
	executionTimeHistogram metrics.Float64Histogram
	queuesConfig           msgconfig.QueuesConfig
}

func NewUserDataAggregationHandler(
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.UserDataAggregationHandlerConfig,
	consumerProvider messagequeue.ConsumerProvider,
	dataPrivacyRepo dataprivacy.Repository,
	uploadManager uploads.UploadManager,
	decoder encoding.ServerEncoderDecoder,
	metricsProvider metrics.Provider,
) (*UserDataAggregationHandler, error) {
	executionTimeHistogram, err := metricsProvider.NewFloat64Histogram("user_data_aggregation_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up user data aggregation execution time histogram: %w", err)
	}

	return &UserDataAggregationHandler{
		tracer:                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                 logging.EnsureLogger(logger).WithName(o11yName),
		consumerProvider:       consumerProvider,
		dataPrivacyRepo:        dataPrivacyRepo,
		uploadManager:          uploadManager,
		decoder:                decoder,
		executionTimeHistogram: executionTimeHistogram,
		queuesConfig:           cfg.Queues,
	}, nil
}

func (h *UserDataAggregationHandler) HandleMessage(ctx context.Context, rawMsg []byte) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var userDataCollectionRequest dataprivacy.UserDataAggregationRequest
	if err := h.decoder.DecodeBytes(ctx, rawMsg, &userDataCollectionRequest); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	logger := h.logger.WithValue(keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID).
		WithValue(keys.UserIDKey, userDataCollectionRequest.UserID)
	tracing.AttachToSpan(span, keys.UserDataAggregationReportIDKey, userDataCollectionRequest.ReportID)
	tracing.AttachToSpan(span, keys.UserIDKey, userDataCollectionRequest.UserID)
	logger.Info("loaded payload, aggregating user data")

	collection, err := h.dataPrivacyRepo.FetchUserDataCollection(ctx, userDataCollectionRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching user data collection")
	}

	logger.Info("compiled user data payload, marshaling")

	collectionBytes, err := json.Marshal(collection)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marshaling collection")
	}

	logger.Info("saving file to object storage")

	if err = h.uploadManager.SaveFile(ctx, fmt.Sprintf("%s.json", userDataCollectionRequest.ReportID), collectionBytes); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "saving collection")
	}

	h.executionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	logger.Info("user data aggregation complete")

	return nil
}

func (h *UserDataAggregationHandler) ConsumeMessages(
	ctx context.Context,
	stopChan chan bool,
	errorsChan chan error,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	consumer, err := h.consumerProvider.ProvideConsumer(
		ctx,
		h.queuesConfig.UserDataAggregationTopicName,
		h.HandleMessage,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, h.logger, span, "configuring user data aggregation requests consumer")
	}

	go consumer.Consume(stopChan, errorsChan)

	go func() {
		for e := range errorsChan {
			h.logger.Error("consuming message", e)
		}
	}()

	return nil
}
