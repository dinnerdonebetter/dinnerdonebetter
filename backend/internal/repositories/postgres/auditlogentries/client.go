package auditlogentries

import (
	"database/sql"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries/generated"

	"github.com/verygoodsoftwarenotvirus/platform/database"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

const (
	o11yName = "audit_log_entries_db_client"
)

// repository is the audit log entry repository implementation.
type repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	readDB           *sql.DB
	writeDB          *sql.DB
}

// ProvideAuditLogRepository provides a new repository.
func ProvideAuditLogRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) audit.Repository {
	c := &repository{
		Client:           client,
		readDB:           client.ReadDB(),
		writeDB:          client.WriteDB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated.New(),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
