package internalops

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/internalops"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops/generated"
)

const (
	o11yName = "internalops_db_client"
)

type repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	readDB           *sql.DB
	writeDB          *sql.DB
}

// ProvideInternalOpsRepository provides a new repository.
func ProvideInternalOpsRepository(logger logging.Logger, tracerProvider tracing.TracerProvider, client database.Client) internalops.InternalOpsDataManager {
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
