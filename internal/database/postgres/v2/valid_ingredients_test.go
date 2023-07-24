package v2

import (
	"context"
	"database/sql"
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
	defaultImage            = "postgres:15"
	defaultDatabaseName     = "dinnerdonebetter"
	defaultDatabaseUsername = "dbuser"
	defaultDatabasePassword = "hunter2"
)

func buildDatabaseClientForTest(t *testing.T, ctx context.Context) (*DatabaseClient, *postgres.PostgresContainer) {
	t.Helper()

	// testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)

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

func TestDatabaseClient_ValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		dbc, container := buildDatabaseClientForTest(t, ctx)

		defer func(t *testing.T) {
			t.Helper()
			assert.NoError(t, container.Terminate(ctx))
		}(t)

		// create
		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		var x ValidIngredient
		require.NoError(t, copier.Copy(&x, exampleValidIngredient))

		created, err := dbc.CreateValidIngredient(ctx, &x)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredient, created)

		// read
		validIngredient, err := dbc.GetValidIngredient(ctx, created.ID)
		exampleValidIngredient.CreatedAt = validIngredient.CreatedAt

		assert.NoError(t, err)
		assert.Equal(t, validIngredient, exampleValidIngredient)

		// update
		updatedValidIngredient := fakes.BuildFakeValidIngredient()
		updatedValidIngredient.ID = created.ID
		var y ValidIngredient
		require.NoError(t, copier.Copy(&y, updatedValidIngredient))
		assert.NoError(t, dbc.UpdateValidIngredient(ctx, updatedValidIngredient))

		// delete
		assert.NoError(t, dbc.ArchiveValidIngredient(ctx, created.ID))

		validIngredient, err = dbc.GetValidIngredient(ctx, created.ID)
		assert.Nil(t, validIngredient)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
