package internalops

import (
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/testing"

	mockdatabase "github.com/primandproper/platform/database/mock"
	"github.com/primandproper/platform/database/postgres"
	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func buildDatabaseClientForTest(t *testing.T) (*repository, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config, nil)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	c := ProvideInternalOpsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)
	require.NoError(t, err)

	return c.(*repository), container
}

func buildInertClientForTest(t *testing.T) *repository {
	t.Helper()

	c := ProvideInternalOpsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), &mockdatabase.ClientMock{ReadDBFunc: func() *sql.DB { return nil }, WriteDBFunc: func() *sql.DB { return nil }})

	return c.(*repository)
}
