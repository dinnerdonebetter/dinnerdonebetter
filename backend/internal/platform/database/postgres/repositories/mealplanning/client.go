package mealplanning

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	settings "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/mealplanning/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "meal_planning_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	identityRepo      identity.Repository
	auditLogEntryRepo auditlogentries.Repository
	db                *sql.DB
}

// ProvideSettingsRepository provides a new client.
func ProvideSettingsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	identityRepo identity.Repository,
	client database.Client,
) settings.Repository {
	c := &Querier{
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
