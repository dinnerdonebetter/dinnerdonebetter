package identity

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/implementations/identity/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

// Querier is the audit log entry client.
type Querier struct {
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	secretGenerator  random.Generator
	timeFunc         func() time.Time
	db               *sql.DB
}

// ProvideAuthRepository provides a new client.
func ProvideAuthRepository(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, db *sql.DB) (identity.Repository, error) {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("auth_db_client"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	c := &Querier{
		db:               db,
		tracer:           tracer,
		timeFunc:         defaultTimeFunc,
		generatedQuerier: generated.New(),
		secretGenerator:  random.NewGenerator(logger, tracerProvider),
		logger:           logging.EnsureLogger(logger).WithName("querier"),
	}

	return c, nil
}

func defaultTimeFunc() time.Time {
	return time.Now()
}

func (q *Querier) currentTime() time.Time {
	if q == nil || q.timeFunc == nil {
		return defaultTimeFunc()
	}

	return q.timeFunc()
}

func (q *Querier) rollbackTransaction(ctx context.Context, tx database.SQLQueryExecutorAndTransactionManager) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Debug("rolling back transaction")

	if err := tx.Rollback(); err != nil {
		observability.AcknowledgeError(err, q.logger, span, "rolling back transaction")
	}

	q.logger.Debug("transaction rolled back")
}
