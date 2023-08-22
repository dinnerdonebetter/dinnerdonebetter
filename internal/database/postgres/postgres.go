package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/cryptography"
	"github.com/dinnerdonebetter/backend/internal/pkg/cryptography/salsa20"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"

	"github.com/Masterminds/squirrel"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	tracingName = "db_client"

	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	postgresDuplicateEntryErrorCode = "23505"
)

var _ database.DataManager = (*Querier)(nil)

// Querier is the primary database querying client. All tracing/logging/query execution happens here. Query building generally happens elsewhere.
type Querier struct {
	tracer                  tracing.Tracer
	logger                  logging.Logger
	sqlBuilder              squirrel.StatementBuilderType
	secretGenerator         random.Generator
	oauth2ClientTokenEncDec cryptography.EncryptorDecryptor
	generatedQuerier        generated.Querier
	timeFunc                func() time.Time
	config                  *dbconfig.Config
	db                      *sql.DB
	pgxDB                   *pgxpool.Pool
	connectionURL           string
	migrateOnce             sync.Once
	logQueries              bool
}

// ProvideDatabaseClient provides a new DataManager client.
func ProvideDatabaseClient(
	ctx context.Context,
	logger logging.Logger,
	cfg *dbconfig.Config,
	tracerProvider tracing.TracerProvider,
) (database.DataManager, error) {
	tracer := tracing.NewTracer(tracerProvider.Tracer(tracingName))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	db, err := sql.Open("pgx", cfg.ConnectionDetails)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)
	db.SetConnMaxLifetime(1800 * time.Second)

	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating encryptor/decryptor with secret length %d", len(cfg.OAuth2TokenEncryptionKey))
	}

	pgxDB, err := pgxpool.New(ctx, cfg.ConnectionDetails)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "initializing PGX pool")
	}

	c := &Querier{
		db:                      db,
		pgxDB:                   pgxDB,
		config:                  cfg,
		tracer:                  tracer,
		logQueries:              cfg.LogQueries,
		timeFunc:                defaultTimeFunc,
		generatedQuerier:        generated.New(),
		secretGenerator:         random.NewGenerator(logger, tracerProvider),
		connectionURL:           cfg.ConnectionDetails,
		logger:                  logging.EnsureLogger(logger).WithName("querier"),
		sqlBuilder:              squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		oauth2ClientTokenEncDec: encDec,
	}

	if cfg.RunMigrations {
		c.logger.Debug("migrating querier")

		if err = c.Migrate(ctx, cfg.PingWaitPeriod, cfg.MaxPingAttempts); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "migrating database")
		}

		c.logger.Debug("querier migrated!")
	}

	return c, nil
}

// DB provides the database object.
func (q *Querier) DB() *sql.DB {
	q.logger.Info("closing database connection")
	return q.db
}

// Close closes the database connection.
func (q *Querier) Close() {
	if err := q.db.Close(); err != nil {
		q.logger.Error(err, "closing database connection")
	}
}

// ProvideSessionStore provides the scs Store for Postgres.
func (q *Querier) ProvideSessionStore() scs.Store {
	return postgresstore.New(q.db)
}

// IsReady is a simple wrapper around the core querier IsReady call.
func (q *Querier) IsReady(ctx context.Context, waitPeriod time.Duration, maxAttempts uint64) (ready bool) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	attemptCount := 0

	logger := q.logger.WithValue("connection_url", q.connectionURL)

	for !ready {
		err := q.db.PingContext(ctx)
		if err != nil {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for db")
			time.Sleep(waitPeriod)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			ready = true
			return ready
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
		q.logger.Error(err, "row error")
		return observability.PrepareAndLogError(err, q.logger, span, "row error")
	}

	if err := rows.Close(); err != nil {
		q.logger.Error(err, "closing database rows")
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

func (q *Querier) getRows(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []any) (*sql.Rows, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("query_desc", queryDescription)
	if q.logQueries {
		logger = logger.WithValue("args", args)
	}

	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s fetch query", queryDescription), query, args)

	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, span, "performing read query")
	}

	if err = rows.Err(); err != nil {
		return nil, observability.PrepareError(err, span, "scanning results")
	}

	if q.logQueries {
		logger.Debug("read query performed")
	}

	return rows, nil
}
