package oauth

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/testing"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

const (
	exampleQuantity = 3
)

func buildDatabaseClientForTest(t *testing.T) (*repository, audit.Repository, *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config, nil)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	auditLogEntryRepo := auditlogentries.ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)

	c := ProvideOAuthRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, config, pgc)
	require.NoError(t, err)

	return c.(*repository), auditLogEntryRepo, container
}

func buildInertClientForTest(t *testing.T) *repository {
	t.Helper()

	config := &databasecfg.Config{
		Provider:                 databasecfg.ProviderPostgres,
		ReadConnection:           databasecfg.ConnectionDetails{},
		Debug:                    false,
		LogQueries:               false,
		RunMigrations:            true,
		MaxPingAttempts:          10,
		PingWaitPeriod:           time.Second,
		OAuth2TokenEncryptionKey: "blahblahblahblahblahblahblahblah",
	}

	c := ProvideOAuthRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil, config, &database.MockClient{})

	return c.(*repository)
}
