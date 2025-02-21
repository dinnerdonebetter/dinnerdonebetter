package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/services/eating/database/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	tracingName = "db_client"

	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	postgresDuplicateEntryErrorCode = "23505"
)

var _ types.MealPlanningDataManager = (*Querier)(nil)
var _ types.RecipeManagementDataManager = (*Querier)(nil)
var _ types.ValidEnumerationDataManager = (*Querier)(nil)

// Querier is the primary database querying client.
type Querier struct {
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	secretGenerator  random.Generator
	timeFunc         func() time.Time
	config           *databasecfg.Config
	db               *sql.DB
	migrateOnce      sync.Once
}

// ProvideDatabaseClient provides a new DataManager client.
func ProvideDatabaseClient(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *databasecfg.Config) (*Querier, error) {
	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(tracingName))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	db, err := otelsql.Open("pgx", cfg.ConnectionDetails.String(), otelsql.WithAttributes(
		attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("database"),
		},
	))
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)
	db.SetConnMaxLifetime(30 * time.Minute)

	c := &Querier{
		db:               db,
		config:           cfg,
		secretGenerator:  random.NewGenerator(logger, tracerProvider),
		tracer:           tracer,
		timeFunc:         defaultTimeFunc,
		generatedQuerier: generated.New(),
		logger:           logging.EnsureLogger(logger).WithName("querier"),
	}

	if cfg.RunMigrations {
		c.logger.Info("migrating querier")

		start := time.Now()
		if err = c.Migrate(ctx); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "migrating database")
		}

		c.logger.WithValue("elapsed", time.Since(start).Milliseconds()).Info("querier migrated!")
	}

	return c, nil
}

// DB provides the database object.
func (q *Querier) DB() *sql.DB {
	return q.db
}

// Close closes the database connection.
func (q *Querier) Close() {
	if err := q.db.Close(); err != nil {
		q.logger.Error("closing database connection", err)
	}
}

// IsReady returns whether the database is ready for the querier.
func (q *Querier) IsReady(ctx context.Context) bool {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("connection_url", q.config.ConnectionDetails.String())

	attemptCount := 0
	for {
		if err := q.db.PingContext(ctx); err != nil {
			logger.WithValue("attempt_count", attemptCount).Info("ping failed, waiting for db")
			time.Sleep(q.config.PingWaitPeriod)

			attemptCount++
			if attemptCount >= int(q.config.MaxPingAttempts) {
				break
			}
		} else {
			return true
		}
	}

	return false
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

func (q *Querier) checkRowsForErrorAndClose(ctx context.Context, rows database.ResultIterator) error {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if err := rows.Err(); err != nil {
		q.logger.Error("row error", err)
		return observability.PrepareAndLogError(err, q.logger, span, "row error")
	}

	if err := rows.Close(); err != nil {
		q.logger.Error("closing database rows", err)
		return observability.PrepareAndLogError(err, q.logger, span, "closing database rows")
	}

	return nil
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

// Destroy deletes all data in the database.
func (q *Querier) Destroy(ctx context.Context) error {
	return q.generatedQuerier.DestroyAllData(ctx, q.db)
}
