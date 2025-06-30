package webhooks

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/migrations"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/auditlogentries"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	exampleQuantity = 3
)

func buildDatabaseClientForTest(t *testing.T) (*Querier, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseClientForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), db, config).Migrate(ctx))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	auditLogEntryRepo := auditlogentries.ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)
	c := ProvideWebhooksRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	require.NoError(t, err)

	return c.(*Querier), container
}

func buildInertClientForTest(t *testing.T) *Querier {
	t.Helper()

	c := ProvideWebhooksRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil, &database.MockClient{})

	return c.(*Querier)
}
