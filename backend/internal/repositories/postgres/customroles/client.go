package customroles

import (
	"database/sql"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/customroles/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "custom_roles_db_client"
)

var _ customroles.Repository = (*repository)(nil)

type repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	readDB           *sql.DB
	writeDB          *sql.DB
}

// ProvideCustomRolesRepository provides a new repository.
func ProvideCustomRolesRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) customroles.Repository {
	return &repository{
		Client:           client,
		readDB:           client.ReadDB(),
		writeDB:          client.WriteDB(),
		tracer:           tracing.NewNamedTracer(tracerProvider, o11yName),
		generatedQuerier: generated.New(),
		logger:           logging.NewNamedLogger(logger, o11yName),
	}
}
