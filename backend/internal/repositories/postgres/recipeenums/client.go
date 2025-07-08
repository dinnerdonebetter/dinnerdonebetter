package recipeenums

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	generated2 "github.com/dinnerdonebetter/backend/internal/repositories/postgres/recipeenums/generated"
)

const (
	o11yName = "recipe_enums_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated2.Querier
	auditLogEntryRepo audit.Repository
	db                *sql.DB
	database.Client
}

// ProvideRepository provides a new client.
func ProvideRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	client database.Client,
) recipeenums.Repository {
	c := &Querier{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated2.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
