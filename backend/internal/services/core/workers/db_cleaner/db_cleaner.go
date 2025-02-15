package dbcleaner

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "db_cleaner"
)

type Job struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	handledRecordsCounter metrics.Int64Counter
	dataManager           types.MaintenanceDataManager
}

func NewDBCleaner(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	dataManager types.MaintenanceDataManager,
) (*Job, error) {
	handledRecordsCounter, err := metricsProvider.NewInt64Counter("db_cleaner.handled_records")
	if err != nil {
		return nil, err
	}

	return &Job{
		logger:                logging.EnsureLogger(logger).WithName(serviceName),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		handledRecordsCounter: handledRecordsCounter,
		dataManager:           dataManager,
	}, nil
}

func (j *Job) Do(ctx context.Context) error {
	ctx, span := j.tracer.StartSpan(ctx)
	defer span.End()

	deleted, err := j.dataManager.DeleteExpiredOAuth2ClientTokens(ctx)
	if err != nil {
		j.logger.Error("deleting expired oauth2 client tokens", err)
		return err
	}

	j.handledRecordsCounter.Add(ctx, deleted, metric.WithAttributes(
		attribute.KeyValue{
			Key:   "db_table",
			Value: attribute.StringValue("oauth2_clients"),
		},
	))

	return nil
}
