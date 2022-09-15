package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"

	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name        = "db_client"
	tracingName = name

	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	postgresDuplicateEntryErrorCode = "23505"
)

var _ database.DataManager = (*Querier)(nil)

// Querier is the primary database querying client. All tracing/logging/query execution happens here. Query building generally happens elsewhere.
type Querier struct {
	tracer        tracing.Tracer
	sqlBuilder    squirrel.StatementBuilderType
	logger        logging.Logger
	timeFunc      func() time.Time
	config        *dbconfig.Config
	db            *sql.DB
	connectionURL string
	migrateOnce   sync.Once
	logQueries    bool
}

var instrumentedDriverRegistration sync.Once

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

	const driverName = "instrumented-postgres"

	instrumentedDriverRegistration.Do(func() {
		sql.Register(
			driverName,
			instrumentedsql.WrapDriver(
				&pq.Driver{},
				instrumentedsql.WithOmitArgs(),
				instrumentedsql.WithTracer(tracing.NewInstrumentedSQLTracer(tracerProvider, "postgres_connection")),
				instrumentedsql.WithLogger(tracing.NewInstrumentedSQLLogger(logger)),
			),
		)
	})

	db, err := sql.Open(driverName, string(cfg.ConnectionDetails))
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	c := &Querier{
		db:            db,
		config:        cfg,
		tracer:        tracer,
		logQueries:    cfg.LogQueries,
		timeFunc:      defaultTimeFunc,
		connectionURL: string(cfg.ConnectionDetails),
		logger:        logging.EnsureLogger(logger).WithName("querier"),
		sqlBuilder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	if cfg.Debug {
		c.logger.SetLevel(logging.DebugLevel)
	}

	if cfg.RunMigrations {
		c.logger.Debug("migrating querier")

		if err = c.Migrate(ctx, cfg.MaxPingAttempts); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "migrating database")
		}

		c.logger.Debug("querier migrated!")
	}

	c.db.SetMaxOpenConns(5)
	c.db.SetMaxOpenConns(7)

	return c, nil
}

// DB provides the scs Store for MySQL.
func (q *Querier) DB() *sql.DB {
	return q.db
}

// ProvideSessionStore provides the scs Store for MySQL.
func (q *Querier) ProvideSessionStore() scs.Store {
	return postgresstore.New(q.db)
}

// IsReady is a simple wrapper around the core querier IsReady call.
func (q *Querier) IsReady(ctx context.Context, maxAttempts uint8) (ready bool) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	attemptCount := 0

	logger := q.logger.WithValue("connection_url", q.connectionURL)

	for !ready {
		err := q.db.PingContext(ctx)
		if err != nil {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for db")
			time.Sleep(time.Second)

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

func (q *Querier) getOneRow(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) *sql.Row {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query = minimizeSQL(query)

	logger := q.logger.WithValue("query_desc", queryDescription)
	if q.logQueries {
		logger = logger.WithValue("query", query).WithValue("args", args)
	}

	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s single row fetch query", queryDescription), query, args)

	row := querier.QueryRowContext(ctx, query, args...)

	if q.logQueries {
		logger.Debug("single row query performed")
	}

	return row
}

func (q *Querier) performReadQuery(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) (*sql.Rows, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query = minimizeSQL(query)

	logger := q.logger.WithValue("query_desc", queryDescription)
	if q.logQueries {
		logger = logger.WithValue("query", query).WithValue("args", args)
	}

	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s fetch query", queryDescription), query, args)

	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing read query")
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, observability.PrepareError(rowsErr, span, "scanning results")
	}

	if q.logQueries {
		logger.Debug("read query performed")
	}

	return rows, nil
}

func (q *Querier) performBooleanQuery(ctx context.Context, querier database.SQLQueryExecutor, query string, args []interface{}) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	var exists bool
	query = minimizeSQL(query)

	logger := q.logger.WithValue(keys.DatabaseQueryKey, query).WithValue("args", args)
	tracing.AttachDatabaseQueryToSpan(span, "boolean query", query, args)

	err := querier.QueryRowContext(ctx, query, args...).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "executing boolean query")
	}

	if q.logQueries {
		logger.Debug("boolean query performed")
	}

	return exists, nil
}

func (q *Querier) performWriteQuery(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	query = minimizeSQL(query)

	logger := q.logger.WithValue("query_desc", queryDescription)
	if q.logQueries {
		logger = logger.WithValue("query", query).WithValue("args", args)
	}

	tracing.AttachDatabaseQueryToSpan(span, queryDescription, query, args)

	res, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "executing %s query", queryDescription)
	}

	var affectedRowCount int64
	if affectedRowCount, err = res.RowsAffected(); affectedRowCount == 0 || err != nil {
		// the only errors returned by the currently supported drivers are either
		// always nil or simply indicate that no rows were affected by the query.

		logger.Debug("no rows modified by query")
		span.AddEvent("no_rows_modified")

		return sql.ErrNoRows
	}

	if q.logQueries {
		logger.Debug("query executed")
	}

	return nil
}
