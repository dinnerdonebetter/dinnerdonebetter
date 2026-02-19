package waitlists

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func buildDatabaseClientForTest(t *testing.T) (*Repository, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	c := ProvideWaitlistsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)

	return c, container
}

func buildInertClientForTest(t *testing.T) *Repository {
	t.Helper()

	c := ProvideWaitlistsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), &database.MockClient{})

	return c
}
