package settings

import (
	"testing"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/implementations/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	exampleQuantity = 3
)

func buildDatabaseClientForTest(t *testing.T) (*Querier, *postgres.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db := pgtesting.BuildDatabaseClientForTest(t)

	config := &databasecfg.Config{
		Provider:          databasecfg.ProviderPostgres,
		ConnectionDetails: databasecfg.ConnectionDetails{},
		Debug:             false,
		LogQueries:        false,
		RunMigrations:     true,
		MaxPingAttempts:   10,
		PingWaitPeriod:    time.Second,
	}

	require.NoError(t, config.LoadConnectionDetailsFromURL(container.MustConnectionString(ctx)))
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), db, config).Migrate(ctx))

	auditLogEntryRepo, err := auditlogentries.ProvideAuditLogRepository(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), db)
	require.NoError(t, err)

	c, err := ProvideSettingsRepository(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, db)
	require.NoError(t, err)

	return c.(*Querier), container
}

func buildInertClientForTest(t *testing.T) *Querier {
	t.Helper()

	c, err := ProvideSettingsRepository(t.Context(), logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil, nil)
	require.NoError(t, err)

	return c.(*Querier)
}
