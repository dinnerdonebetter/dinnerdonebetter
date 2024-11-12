package logic

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search/text/indexing"

	"go.opentelemetry.io/otel"
)

func HandleIndexDataRequest(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.InstanceConfig,
	searchIndexRequest *indexing.IndexRequest,
) error {
	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search_indexer_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
	if err = indexing.HandleIndexRequest(ctx, logger, tracerProvider, &cfg.Search, dataManager, searchIndexRequest); err != nil {
		observability.AcknowledgeError(err, logger, span, "handling index request")
	}

	return nil
}
