package v2

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage            = "postgres:15"
	defaultDatabaseName     = "dinnerdonebetter"
	defaultDatabaseUsername = "dbuser"
	defaultDatabasePassword = "hunter2"
	exampleQuantity         = 3
)

func buildDatabaseClientForTest(t *testing.T, ctx context.Context) (*DatabaseClient, *postgres.PostgresContainer) {
	t.Helper()

	testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(defaultImage),
		postgres.WithDatabase(defaultDatabaseName),
		postgres.WithUsername(defaultDatabaseUsername),
		postgres.WithPassword(defaultDatabasePassword),
		testcontainers.WithWaitStrategyAndDeadline(
			time.Minute,
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)

	require.NoError(t, err)
	require.NotNil(t, container)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbc, err := NewDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), &dbconfig.Config{ConnectionDetails: connStr, RunMigrations: true})
	require.NoError(t, err)
	require.NotNil(t, dbc)

	return dbc, container
}
