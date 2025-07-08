package recipes

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/recipes/generated"
)

const (
	o11yName = "meal_planning_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	mealsRepo         mealplanning.Repository
	recipeEnumsRepo   recipeenums.Repository
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	identityRepo      identity.Repository
	auditLogEntryRepo audit.Repository
	db                *sql.DB
}

// ProvideSettingsRepository provides a new client.
func ProvideSettingsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	identityRepo identity.Repository,
	mealsRepo mealplanning.Repository,
	recipeEnumsRepo recipeenums.Repository,
	client database.Client,
) recipes.Repository {
	c := &Querier{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		identityRepo:      identityRepo,
		mealsRepo:         mealsRepo,
		recipeEnumsRepo:   recipeEnumsRepo,
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
