package auditlogentries

import (
	"testing"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/implementations/auditlogentries/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	exampleQuantity = 3
)

type sqlmockExpecterWrapper struct {
	sqlmock.Sqlmock
}

func (e *sqlmockExpecterWrapper) AssertExpectations(t mock.TestingT) bool {
	return assert.NoError(t, e.Sqlmock.ExpectationsWereMet(), "not all database expectations were met")
}

func buildMockSQLTestClient(t *testing.T) (*Querier, *sqlmockExpecterWrapper) {
	t.Helper()

	fakeDB, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	c := &Querier{
		db:               fakeDB,
		logger:           logging.NewNoopLogger(),
		generatedQuerier: generated.New(),
		timeFunc:         defaultTimeFunc,
		tracer:           tracing.NewTracerForTest("test"),
	}

	return c, &sqlmockExpecterWrapper{Sqlmock: sqlMock}
}

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

	c, err := ProvideAuthRepository(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), db)
	require.NoError(t, err)

	return c.(*Querier), container
}

func buildInertClientForTest(t *testing.T) *Querier {
	t.Helper()

	c, err := ProvideAuthRepository(t.Context(), logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil)
	require.NoError(t, err)

	return c.(*Querier)
}
