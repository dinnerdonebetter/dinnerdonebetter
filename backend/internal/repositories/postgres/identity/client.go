package identity

import (
	"database/sql"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/customroles"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/random"
)

const (
	o11yName = "identity_db_client"
)

var _ identity.Repository = (*repository)(nil)

// repository is the identity repository implementation.
type repository struct {
	database.Client
	tracer              tracing.Tracer
	logger              logging.Logger
	generatedQuerier    generated.Querier
	auditLogEntryRepo   audit.Repository
	secretGenerator     random.Generator
	customRolesRepo     customroles.Repository
	rolePermissionCache *authorization.RolePermissionCache
	readDB              *sql.DB
	writeDB             *sql.DB
}

// ProvideIdentityRepository provides a new repository.
func ProvideIdentityRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	client database.Client,
	customRolesRepo customroles.Repository,
	rolePermissionCache *authorization.RolePermissionCache,
) identity.Repository {
	c := &repository{
		Client:              client,
		readDB:              client.ReadDB(),
		writeDB:             client.WriteDB(),
		tracer:              tracing.NewNamedTracer(tracerProvider, o11yName),
		generatedQuerier:    generated.New(),
		auditLogEntryRepo:   auditLogEntryRepo,
		secretGenerator:     random.NewGenerator(logger, tracerProvider),
		logger:              logging.NewNamedLogger(logger, o11yName),
		customRolesRepo:     customRolesRepo,
		rolePermissionCache: rolePermissionCache,
	}

	return c
}
