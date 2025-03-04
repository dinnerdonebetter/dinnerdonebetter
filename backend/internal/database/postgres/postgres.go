package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/lib/cryptography/encryption"
	"github.com/dinnerdonebetter/backend/internal/lib/cryptography/encryption/salsa20"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/internal/lib/random"

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

var _ database.DataManager = (*Querier)(nil)

// Querier is the primary database querying client.
type Querier struct {
	tracer                  tracing.Tracer
	logger                  logging.Logger
	secretGenerator         random.Generator
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	generatedQuerier        generated.Querier
	timeFunc                func() time.Time
	config                  *databasecfg.Config
	db                      *sql.DB
	migrateOnce             sync.Once
}

// ProvideDatabaseClient provides a new DataManager client.
func ProvideDatabaseClient(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *databasecfg.Config) (database.DataManager, error) {
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

	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating encryptor/decryptor with secret length %d", len(cfg.OAuth2TokenEncryptionKey))
	}

	c := &Querier{
		db:                      db,
		config:                  cfg,
		tracer:                  tracer,
		timeFunc:                defaultTimeFunc,
		generatedQuerier:        generated.New(),
		secretGenerator:         random.NewGenerator(logger, tracerProvider),
		logger:                  logging.EnsureLogger(logger).WithName("querier"),
		oauth2ClientTokenEncDec: encDec,
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

func fetchAllRows[T any](fetchFunc func(*filtering.QueryFilter) (*filtering.QueryFilteredResult[T], error)) ([]T, error) {
	var done bool
	allData := []T{}

	filter := &filtering.QueryFilter{
		Page:            pointer.To(uint16(1)),
		Limit:           pointer.To(uint8(math.MaxUint8)),
		IncludeArchived: pointer.To(true),
	}

	for !done {
		data, err := fetchFunc(filter)
		if err != nil {
			return nil, fmt.Errorf("getting data: %w", err)
		}

		for _, x := range data.Data {
			if x != nil {
				allData = append(allData, *x)
			}
		}

		if data.TotalCount <= uint64(len(allData)) {
			done = true
		}
		filter.Page = pointer.To(*filter.Page + 1)
	}

	return allData, nil
}

// Destroy deletes all data in the database.
func (q *Querier) Destroy(ctx context.Context) error {
	return q.generatedQuerier.DestroyAllData(ctx, q.db)
}
