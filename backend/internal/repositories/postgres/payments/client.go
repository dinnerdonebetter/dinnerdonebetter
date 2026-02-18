package payments

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments/generated"
)

const (
	o11yName = "payments_db_client"
)

type repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	readDB           *sql.DB
	writeDB          *sql.DB
}

// ProvidePaymentsRepository provides a new payments repository.
func ProvidePaymentsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) payments.Repository {
	r := &repository{
		Client:           client,
		readDB:           client.ReadDB(),
		writeDB:          client.WriteDB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated.New(),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}
	var _ payments.Repository = r
	return r
}
