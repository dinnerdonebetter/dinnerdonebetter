package identity

import (
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/identity/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

const (
	o11yName = "identity_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogEntryRepo auditlogentries.Repository
	secretGenerator   random.Generator
	timeFunc          func() time.Time
	db                *sql.DB
}

// ProvideIdentityRepository provides a new client.
func ProvideIdentityRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	client database.Client,
) identity.Repository {
	c := &Querier{
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
