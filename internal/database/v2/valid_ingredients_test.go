package v2

import (
	"context"
	"testing"
	"time"

	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage = "postgres:15"
)

func buildDatabaseClientForTest(t *testing.T, ctx context.Context) (*DatabaseClient, *postgres.PostgresContainer) {
	t.Helper()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(defaultImage),
		postgres.WithDatabase("database"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	require.NotNil(t, container)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbc, err := NewDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), &dbconfig.Config{ConnectionDetails: connStr, RunMigrations: true})
	require.NoError(t, err)
	require.NotNil(t, dbc)

	return dbc, container
}

func TestDatabaseClient_ValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		dbc, container := buildDatabaseClientForTest(t, ctx)

		defer func(tt *testing.T) {
			assert.NoError(tt, container.Terminate(ctx))
		}(t)

		var f ValidIngredient
		require.NoError(t, copier.Copy(&f, fakes.BuildFakeValidIngredient()))

		created, err := dbc.CreateValidIngredient(ctx, &f)
		assert.NoError(t, err)
		assert.NotNil(t, created)

		x, err := dbc.GetValidIngredient(ctx, created.ID)
		created.CreatedAt = x.CreatedAt

		assert.NoError(t, err)
		assert.Equal(t, x, created)
	})
}
