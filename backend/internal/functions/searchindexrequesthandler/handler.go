package searchindexrequesthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	mealplanningindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
)

const (
	o11yName = "search_index_request_handler"
)

type SearchIndexRequestHandler struct {
	logger                  logging.Logger
	tracer                  tracing.Tracer
	consumerProvider        messagequeue.ConsumerProvider
	userDataIndexer         *identityindexing.UserDataIndexer
	mealPlanningDataIndexer *mealplanningindexing.MealPlanningDataIndexer
	executionTimeHistogram  metrics.Float64Histogram
	queuesConfig            msgconfig.QueuesConfig
}

func NewSearchIndexRequestHandler(
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.SearchIndexRequestHandlerConfig,
	consumerProvider messagequeue.ConsumerProvider,
	metricsProvider metrics.Provider,
	userDataIndexer *identityindexing.UserDataIndexer,
	mealPlanningDataIndexer *mealplanningindexing.MealPlanningDataIndexer,
) (*SearchIndexRequestHandler, error) {
	executionTimeHistogram, err := metricsProvider.NewFloat64Histogram("search_index_requests_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up search index requests execution time histogram: %w", err)
	}

	return &SearchIndexRequestHandler{
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                  logging.EnsureLogger(logger).WithName(o11yName),
		consumerProvider:        consumerProvider,
		userDataIndexer:         userDataIndexer,
		mealPlanningDataIndexer: mealPlanningDataIndexer,
		executionTimeHistogram:  executionTimeHistogram,
		queuesConfig:            cfg.Queues,
	}, nil
}

func (h *SearchIndexRequestHandler) HandleMessage(ctx context.Context, rawMsg []byte) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var searchIndexRequest textsearch.IndexRequest
	if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&searchIndexRequest); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	switch searchIndexRequest.IndexType {
	case mealplanningindexing.IndexTypeRecipes,
		mealplanningindexing.IndexTypeMeals,
		mealplanningindexing.IndexTypeValidIngredients,
		mealplanningindexing.IndexTypeValidInstruments,
		mealplanningindexing.IndexTypeValidMeasurementUnits,
		mealplanningindexing.IndexTypeValidPreparations,
		mealplanningindexing.IndexTypeValidIngredientStates,
		mealplanningindexing.IndexTypeValidVessels:
		if err := h.mealPlanningDataIndexer.HandleIndexRequest(ctx, &searchIndexRequest); err != nil {
			return fmt.Errorf("handling search indexing request: %w", err)
		}

	case identityindexing.IndexTypeUsers:
		if err := h.userDataIndexer.HandleIndexRequest(ctx, &searchIndexRequest); err != nil {
			return fmt.Errorf("handling search indexing request: %w", err)
		}
	}

	h.executionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	return nil
}

func (h *SearchIndexRequestHandler) ConsumeMessages(
	ctx context.Context,
	stopChan chan bool,
	errorsChan chan error,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	consumer, err := h.consumerProvider.ProvideConsumer(
		ctx,
		h.queuesConfig.SearchIndexRequestsTopicName,
		h.HandleMessage,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, h.logger, span, "configuring search index requests consumer")
	}

	go consumer.Consume(stopChan, errorsChan)

	go func() {
		for e := range errorsChan {
			h.logger.Error("consuming message", e)
		}
	}()

	return nil
}
