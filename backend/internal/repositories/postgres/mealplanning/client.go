package mealplanning

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

const (
	o11yName = "meal_planning_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	recipeEnumsRepository recipeenums.Repository
	tracer                tracing.Tracer
	logger                logging.Logger
	generatedQuerier      generated.Querier
	identityRepo          identity.Repository
	auditLogEntryRepo     audit.Repository
	db                    *sql.DB
}

// ProvideRepository provides a new client.
func ProvideRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	identityRepo identity.Repository,
	recipeEnumsRepository recipeenums.Repository,
	client database.Client,
) mealplanning.Repository {
	c := &Querier{
		Client:                client,
		db:                    client.DB(),
		tracer:                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:      generated.New(),
		auditLogEntryRepo:     auditLogEntryRepo,
		identityRepo:          identityRepo,
		recipeEnumsRepository: recipeEnumsRepository,
		logger:                logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
