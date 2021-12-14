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
	"go.opentelemetry.io/otel/trace"

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
)

var _ database.DataManager = (*SQLQuerier)(nil)

// SQLQuerier is the primary database querying client. All tracing/logging/query execution happens here. Query building generally happens elsewhere.
type SQLQuerier struct {
	tracer        tracing.Tracer
	sqlBuilder    squirrel.StatementBuilderType
	logger        logging.Logger
	db            *sql.DB
	timeFunc      func() uint64
	config        *dbconfig.Config
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
	tracerProvider trace.TracerProvider,
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

	c := &SQLQuerier{
		db:            db,
		config:        cfg,
		tracer:        tracer,
		logQueries:    true,
		timeFunc:      defaultTimeFunc,
		connectionURL: string(cfg.ConnectionDetails),
		logger:        logging.EnsureLogger(logger),
		sqlBuilder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	if cfg.Debug {
		c.logger.SetLevel(logging.DebugLevel)
	}

	if cfg.RunMigrations {
		c.logger.Debug("migrating querier")

		if err = c.Migrate(ctx, cfg.MaxPingAttempts); err != nil {
			return nil, observability.PrepareError(err, logger, span, "migrating database")
		}

		c.logger.Debug("querier migrated!")
	}

	return c, nil
}

// ProvideSessionStore provides the scs Store for MySQL.
func (q *SQLQuerier) ProvideSessionStore() scs.Store {
	return postgresstore.New(q.db)
}

// IsReady is a simple wrapper around the core querier IsReady call.
func (q *SQLQuerier) IsReady(ctx context.Context, maxAttempts uint8) (ready bool) {
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

func defaultTimeFunc() uint64 {
	return uint64(time.Now().Unix())
}

func (q *SQLQuerier) currentTime() uint64 {
	if q == nil || q.timeFunc == nil {
		return defaultTimeFunc()
	}

	return q.timeFunc()
}

func (q *SQLQuerier) checkRowsForErrorAndClose(ctx context.Context, rows database.ResultIterator) error {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if err := rows.Err(); err != nil {
		q.logger.Error(err, "row error")
		return observability.PrepareError(err, logger, span, "row error")
	}

	if err := rows.Close(); err != nil {
		q.logger.Error(err, "closing database rows")
		return observability.PrepareError(err, logger, span, "closing database rows")
	}

	return nil
}

func (q *SQLQuerier) rollbackTransaction(ctx context.Context, tx *sql.Tx) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if err := tx.Rollback(); err != nil {
		observability.AcknowledgeError(err, q.logger, span, "rolling back transaction")
	}

	q.logger.Debug("transaction rolled back")
}

func (q *SQLQuerier) getOneRow(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) *sql.Row {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("query", query).WithValue("args", args)
	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s single row fetch query", queryDescription), query, args)

	row := querier.QueryRowContext(ctx, query, args...)

	if q.logQueries {
		logger.Debug("single row query performed")
	}

	return row
}

func (q *SQLQuerier) performReadQuery(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) (*sql.Rows, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("query", query).WithValue("args", args)

	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s fetch query", queryDescription), query, args)

	rows, err := querier.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing query")
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, observability.PrepareError(rowsErr, logger, span, "scanning results")
	}

	if q.logQueries {
		logger.Debug("read query performed")
	}

	return rows, nil
}

func (q *SQLQuerier) performCountQuery(ctx context.Context, querier database.SQLQueryExecutor, query, queryDesc string) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("query", query)

	tracing.AttachDatabaseQueryToSpan(span, fmt.Sprintf("%s count query", queryDesc), query, nil)

	var count uint64
	if err := q.getOneRow(ctx, querier, queryDesc, query, []interface{}{}).Scan(&count); err != nil {
		return 0, observability.PrepareError(err, logger, span, "executing count query")
	}

	if q.logQueries {
		logger.Debug("count query performed")
	}

	return count, nil
}

func (q *SQLQuerier) performBooleanQuery(ctx context.Context, querier database.SQLQueryExecutor, query string, args []interface{}) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	var exists bool

	logger := q.logger.WithValue(keys.DatabaseQueryKey, query).WithValue("args", args)
	tracing.AttachDatabaseQueryToSpan(span, "boolean query", query, args)

	err := querier.QueryRowContext(ctx, query, args...).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "executing boolean query")
	}

	if q.logQueries {
		logger.Debug("boolean query performed")
	}

	return exists, nil
}

func (q *SQLQuerier) performWriteQuery(ctx context.Context, querier database.SQLQueryExecutor, queryDescription, query string, args []interface{}) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("query", query).WithValue("description", queryDescription).WithValue("args", args)
	tracing.AttachDatabaseQueryToSpan(span, queryDescription, query, args)

	res, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return observability.PrepareError(err, logger, span, "executing %s query", queryDescription)
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
