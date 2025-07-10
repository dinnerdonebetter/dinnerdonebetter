package maintenance

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/admin"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	generated2 "github.com/dinnerdonebetter/backend/internal/repositories/postgres/maintenance/generated"
)

const (
	o11yName = "webhook_db_client"
)

// Querier is the webhook repository client.
type Querier struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated2.Querier
	db               *sql.DB
}

// ProvideMaintenanceRepository provides a new client.
func ProvideMaintenanceRepository(logger logging.Logger, tracerProvider tracing.TracerProvider, client database.Client) admin.MaintenanceDataManager {
	c := &Querier{
		Client:           client,
		db:               client.DB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated2.New(),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
