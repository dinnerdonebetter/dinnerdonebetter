package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildRecipeForTestCreation(t *testing.T, ctx context.Context, userID string, dbc *Querier) *types.Recipe {
	t.Helper()

	if userID == "" {
		user := createUserForTest(t, ctx, nil, dbc)
		userID = user.ID
	}

	exampleRecipe := fakes.BuildFakeRecipe()
	preparation := createValidPreparationForTest(t, ctx, nil, dbc)

	exampleRecipe.Steps = []*types.RecipeStep{
		{
			BelongsToRecipe: exampleRecipe.ID,
			ID:              identifiers.New(),
			Preparation:     *preparation,
			Index:           0,
		},
	}
	exampleRecipe.Media = []*types.RecipeMedia{}
	exampleRecipe.PrepTasks = []*types.RecipePrepTask{}
	exampleRecipe.CreatedByUser = userID

	return exampleRecipe
}

func createRecipeForTest(t *testing.T, ctx context.Context, exampleRecipe *types.Recipe, dbc *Querier) *types.Recipe {
	t.Helper()

	// create
	if exampleRecipe == nil {
		exampleRecipe = fakes.BuildFakeRecipe()
	}
	dbInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

	created, err := dbc.CreateRecipe(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	exampleRecipe.ID = created.ID
	exampleRecipe.CreatedAt = created.CreatedAt
	for i := range created.Steps {
		assert.Equal(t, exampleRecipe.Steps[i].Preparation.ID, created.Steps[i].Preparation.ID)
		exampleRecipe.Steps[i].Preparation = created.Steps[i].Preparation
		exampleRecipe.Steps[i].CreatedAt = created.Steps[i].CreatedAt
	}

	require.Equal(t, exampleRecipe, created)

	recipe, err := dbc.GetRecipe(ctx, created.ID)
	require.NoError(t, err)
	for i := range recipe.Steps {
		assert.Equal(t, exampleRecipe.Steps[i].Preparation.ID, recipe.Steps[i].Preparation.ID)
		exampleRecipe.Steps[i].Preparation = recipe.Steps[i].Preparation
		exampleRecipe.Steps[i].CreatedAt = recipe.Steps[i].CreatedAt
	}
	exampleRecipe.CreatedAt = recipe.CreatedAt

	// TODO: actually assert these
	exampleRecipe.PrepTasks = recipe.PrepTasks
	exampleRecipe.Media = recipe.Media
	exampleRecipe.SupportingRecipes = recipe.SupportingRecipes

	assert.Equal(t, recipe, exampleRecipe)

	return recipe
}

func TestQuerier_Integration_Recipes(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)

	createdRecipes := []*types.Recipe{}

	// create
	createdRecipes = append(createdRecipes, createRecipeForTest(t, ctx, exampleRecipe, dbc))

	// update
	updatedRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	updatedRecipe.ID = createdRecipes[0].ID
	assert.NoError(t, dbc.UpdateRecipe(ctx, updatedRecipe))

	// fetch as list
	recipes, err := dbc.GetRecipes(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipes.Data)
	assert.Equal(t, len(createdRecipes), len(recipes.Data))

	// delete
	for _, recipe := range createdRecipes {
		assert.NoError(t, dbc.ArchiveRecipe(ctx, recipe.ID, user.ID))

		var exists bool
		exists, err = dbc.RecipeExists(ctx, recipe.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.Recipe
		y, err = dbc.GetRecipe(ctx, recipe.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
