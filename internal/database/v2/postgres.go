package v2

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

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

	db, err := sql.Open("pgx", string(cfg.ConnectionDetails))
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
