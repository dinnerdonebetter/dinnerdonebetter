package auditlogentries

import (
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/testing"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/database"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

const (
	exampleQuantity = 3
)

func buildDatabaseClientForTest(t *testing.T) (*repository, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config, nil)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	c := ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)

	return c.(*repository), container
}

func buildInertClientForTest(t *testing.T) *repository {
	t.Helper()

	c := ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), &database.MockClient{})

	return c.(*repository)
}
