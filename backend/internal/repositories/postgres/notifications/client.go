package notifications

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications/generated"
)

const (
	o11yName = "notifications_db_client"
)

// repository is the notifications repository implementation.
type repository struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogEntryRepo audit.Repository
	secretGenerator   random.Generator
	db                *sql.DB
}

// ProvideNotificationsRepository provides a new repository.
func ProvideNotificationsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	client database.Client,
) notifications.Repository {
	c := &repository{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		secretGenerator:   random.NewGenerator(logger, tracerProvider),
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
