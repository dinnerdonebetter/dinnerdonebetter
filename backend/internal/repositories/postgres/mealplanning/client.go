package mealplanning

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

const (
	o11yName = "meal_planning_db_client"
)

// repository is the meal planning repository implementation.
type repository struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	identityRepo      identity.Repository
	auditLogEntryRepo audit.Repository
	db                *sql.DB
}

// ProvideMealPlanningRepository provides a new repository.
func ProvideMealPlanningRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	identityRepo identity.Repository,
	client database.Client,
) mealplanning.Repository {
	c := &repository{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		identityRepo:      identityRepo,
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
