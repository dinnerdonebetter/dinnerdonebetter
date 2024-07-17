package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUserIngredientPreferenceForTest(t *testing.T, ctx context.Context, exampleUserIngredientPreference *types.UserIngredientPreference, dbc *Querier) *types.UserIngredientPreference {
	t.Helper()

	// create
	if exampleUserIngredientPreference == nil {
		exampleUserIngredientPreference = fakes.BuildFakeUserIngredientPreference()
	}
	dbInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceDatabaseCreationInput(exampleUserIngredientPreference)

	createdPreferences, err := dbc.CreateUserIngredientPreference(ctx, dbInput)
	require.Len(t, createdPreferences, 1)
	created := createdPreferences[0]
	exampleUserIngredientPreference.ID = created.ID
	exampleUserIngredientPreference.CreatedAt = created.CreatedAt
	exampleUserIngredientPreference.Ingredient = created.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, exampleUserIngredientPreference, created)

	userIngredientPreference, err := dbc.GetUserIngredientPreference(ctx, created.ID, dbInput.BelongsToUser)
	exampleUserIngredientPreference.CreatedAt = userIngredientPreference.CreatedAt
	exampleUserIngredientPreference.Ingredient = userIngredientPreference.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, userIngredientPreference, exampleUserIngredientPreference)

	return created
}

func TestQuerier_Integration_UserIngredientPreferences(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	ingredient := createValidIngredientForTest(t, ctx, nil, dbc)

	exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
	exampleUserIngredientPreference.BelongsToUser = user.ID
	exampleUserIngredientPreference.Ingredient = *ingredient
	createdUserIngredientPreferences := []*types.UserIngredientPreference{}

	// create
	createdUserIngredientPreferences = append(createdUserIngredientPreferences, createUserIngredientPreferenceForTest(t, ctx, exampleUserIngredientPreference, dbc))

	// update
	ingredient2 := createValidIngredientForTest(t, ctx, nil, dbc)
	updatedUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
	updatedUserIngredientPreference.ID = createdUserIngredientPreferences[0].ID
	updatedUserIngredientPreference.BelongsToUser = user.ID
	updatedUserIngredientPreference.Ingredient = *ingredient2
	assert.NoError(t, dbc.UpdateUserIngredientPreference(ctx, updatedUserIngredientPreference))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeUserIngredientPreference()
		ingredient3 := createValidIngredientForTest(t, ctx, nil, dbc)
		input.BelongsToUser = user.ID
		input.Ingredient = *ingredient3
		createdUserIngredientPreferences = append(createdUserIngredientPreferences, createUserIngredientPreferenceForTest(t, ctx, input, dbc))
	}

	// fetch as list
	userIngredientPreferences, err := dbc.GetUserIngredientPreferences(ctx, user.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, userIngredientPreferences.Data)
	assert.Equal(t, len(createdUserIngredientPreferences), len(userIngredientPreferences.Data))

	// delete
	for _, userIngredientPreference := range createdUserIngredientPreferences {
		assert.NoError(t, dbc.ArchiveUserIngredientPreference(ctx, userIngredientPreference.ID, user.ID))

		var exists bool
		exists, err = dbc.UserIngredientPreferenceExists(ctx, userIngredientPreference.ID, user.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.UserIngredientPreference
		y, err = dbc.GetUserIngredientPreference(ctx, userIngredientPreference.ID, user.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_UserIngredientPreferenceExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIngredientPreferenceExists(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserIngredientPreferenceID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIngredientPreferenceExists(ctx, exampleUserIngredientPreferenceID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserIngredientPreference(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUserIngredientPreference(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserIngredientPreference(ctx, nil))
	})
}

func TestQuerier_ArchiveUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveUserIngredientPreference(ctx, "", exampleUserID))
	})
}
