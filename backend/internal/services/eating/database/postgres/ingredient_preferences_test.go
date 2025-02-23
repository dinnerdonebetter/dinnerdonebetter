package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUserIngredientPreferenceForTest(t *testing.T, ctx context.Context, exampleUserIngredientPreference *types.IngredientPreference, dbc *Querier) *types.IngredientPreference {
	t.Helper()

	// create
	if exampleUserIngredientPreference == nil {
		exampleUserIngredientPreference = fakes.BuildFakeIngredientPreference()
	}
	dbInput := converters.ConvertIngredientPreferenceToIngredientPreferenceDatabaseCreationInput(exampleUserIngredientPreference)

	createdPreferences, err := dbc.CreateIngredientPreference(ctx, dbInput)
	require.Len(t, createdPreferences, 1)
	created := createdPreferences[0]
	exampleUserIngredientPreference.ID = created.ID
	exampleUserIngredientPreference.CreatedAt = created.CreatedAt
	exampleUserIngredientPreference.Ingredient = created.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, exampleUserIngredientPreference, created)

	userIngredientPreference, err := dbc.GetIngredientPreference(ctx, created.ID, dbInput.BelongsToUser)
	require.NotNil(t, userIngredientPreference)
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

	ctx := t.Context()
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

	exampleUserIngredientPreference := fakes.BuildFakeIngredientPreference()
	exampleUserIngredientPreference.BelongsToUser = user.ID
	exampleUserIngredientPreference.Ingredient = *ingredient
	createdUserIngredientPreferences := []*types.IngredientPreference{}

	// create
	createdUserIngredientPreferences = append(createdUserIngredientPreferences, createUserIngredientPreferenceForTest(t, ctx, exampleUserIngredientPreference, dbc))

	// update
	ingredient2 := createValidIngredientForTest(t, ctx, nil, dbc)
	updatedUserIngredientPreference := fakes.BuildFakeIngredientPreference()
	updatedUserIngredientPreference.ID = createdUserIngredientPreferences[0].ID
	updatedUserIngredientPreference.BelongsToUser = user.ID
	updatedUserIngredientPreference.Ingredient = *ingredient2
	assert.NoError(t, dbc.UpdateIngredientPreference(ctx, updatedUserIngredientPreference))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeIngredientPreference()
		ingredient1 := createValidIngredientForTest(t, ctx, nil, dbc)
		input.BelongsToUser = user.ID
		input.Ingredient = *ingredient1
		createdUserIngredientPreferences = append(createdUserIngredientPreferences, createUserIngredientPreferenceForTest(t, ctx, input, dbc))
	}

	// fetch as list
	userIngredientPreferences, err := dbc.GetIngredientPreferences(ctx, user.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, userIngredientPreferences.Data)
	assert.Equal(t, len(createdUserIngredientPreferences), len(userIngredientPreferences.Data))

	// delete
	for _, userIngredientPreference := range createdUserIngredientPreferences {
		assert.NoError(t, dbc.ArchiveIngredientPreference(ctx, userIngredientPreference.ID, user.ID))

		var exists bool
		exists, err = dbc.IngredientPreferenceExists(ctx, userIngredientPreference.ID, user.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.IngredientPreference
		y, err = dbc.GetIngredientPreference(ctx, userIngredientPreference.ID, user.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_UserIngredientPreferenceExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.IngredientPreferenceExists(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleUserIngredientPreferenceID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.IngredientPreferenceExists(ctx, exampleUserIngredientPreferenceID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetIngredientPreference(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		actual, err := c.CreateIngredientPreference(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateIngredientPreference(ctx, nil))
	})
}

func TestQuerier_ArchiveUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveIngredientPreference(ctx, "", exampleUserID))
	})
}
