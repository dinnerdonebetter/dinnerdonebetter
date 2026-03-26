package payments

import (
	"database/sql"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/payments/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/tracing"
)

const (
	o11yName = "payments_db_client"
)

type repository struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogEntryRepo audit.Repository
	readDB            *sql.DB
	writeDB           *sql.DB
}

// ProvidePaymentsRepository provides a new payments repository.
func ProvidePaymentsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	client database.Client,
) payments.Repository {
	r := &repository{
		Client:            client,
		readDB:            client.ReadDB(),
		writeDB:           client.WriteDB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}
	var _ payments.Repository = r
	return r
}
