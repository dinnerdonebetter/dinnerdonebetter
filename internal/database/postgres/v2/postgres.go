package v2

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	idColumn         = "id"
	archivedAtColumn = "archived_at"
)

type (
	DatabaseClient struct {
		_           struct{}
		tracer      tracing.Tracer
		logger      logging.Logger
		xdb         *goqu.Database
		db          *sql.DB
		migrateOnce sync.Once
	}
)

func NewDatabaseClient(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *dbconfig.Config,
) (*DatabaseClient, error) {
	x := &DatabaseClient{
		tracer: tracing.NewTracer(tracerProvider.Tracer("postgres")),
		logger: logging.EnsureLogger(logger).WithName("postgres"),
	}

	db, err := sql.Open("pgx", cfg.ConnectionDetails)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	x.db = db
	x.xdb = goqu.New("postgres", db)

	if cfg.RunMigrations {
		if err = x.Migrate(ctx, cfg.PingWaitPeriod, cfg.MaxPingAttempts); err != nil {
			return nil, fmt.Errorf("migrating database: %w", err)
		}
	}

	return x, nil
}

func queryFilterToGoqu(q *goqu.SelectDataset, filter *types.QueryFilter) goqu.Ex {
	x := goqu.Ex{}

	if filter == nil {
		return x
	}

	var (
		createdAtOp = goqu.Op{}
		updatedAtOp = goqu.Op{}
	)

	if filter.CreatedAfter != nil {
		createdAtOp["gt"] = filter.CreatedAfter
	}

	if filter.CreatedBefore != nil {
		createdAtOp["lt"] = filter.CreatedAfter
	}

	if len(createdAtOp) > 0 {
		x["created_at"] = createdAtOp
	}

	if filter.UpdatedAfter != nil {
		updatedAtOp["gt"] = filter.UpdatedAfter
	}

	if filter.UpdatedBefore != nil {
		updatedAtOp["lt"] = filter.UpdatedBefore
	}

	if len(updatedAtOp) > 0 {
		x["last_updated_at"] = updatedAtOp
	}

	q = q.Where(x)

	if filter.Page != nil {
		q = q.Offset(uint(*filter.Page))
	}

	if filter.Limit != nil {
		q = q.Limit(uint(*filter.Limit))
	}

	return x
}
