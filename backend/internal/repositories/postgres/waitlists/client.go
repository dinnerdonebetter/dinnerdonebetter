package waitlists

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/waitlists"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/waitlists/generated"
)

const (
	o11yName = "waitlists_db_client"
)

// Repository is the waitlists repository implementation.
type Repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	readDB           *sql.DB
	writeDB          *sql.DB
}

// ProvideWaitlistsRepository provides a new repository.
// Returns concrete *repository so the manager can wrap it; the manager is the sole provider of waitlists.Repository for services.
func ProvideWaitlistsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) *Repository {
	c := &Repository{
		Client:           client,
		readDB:           client.ReadDB(),
		writeDB:          client.WriteDB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated.New(),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}

// Ensure *Repository implements the interface.
var _ waitlists.Repository = (*Repository)(nil)
