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

func createRecipeStepIngredientForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepIngredient *types.RecipeStepIngredient, dbc *Querier) *types.RecipeStepIngredient {
	t.Helper()

	// create
	if exampleRecipeStepIngredient == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
		exampleRecipeStep := createdRecipe.Steps[0]

		exampleRecipeStepIngredient = fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
	}
	dbInput := converters.ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(exampleRecipeStepIngredient)

	created, err := dbc.CreateRecipeStepIngredient(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepIngredient.CreatedAt = created.CreatedAt
	exampleRecipeStepIngredient.Ingredient.CreatedAt = created.Ingredient.CreatedAt
	exampleRecipeStepIngredient.Ingredient = created.Ingredient
	exampleRecipeStepIngredient.MeasurementUnit = created.MeasurementUnit
	assert.Equal(t, exampleRecipeStepIngredient, created)

	recipeStepIngredient, err := dbc.GetRecipeStepIngredient(ctx, recipeID, exampleRecipeStepIngredient.BelongsToRecipeStep, exampleRecipeStepIngredient.ID)

	exampleRecipeStepIngredient.CreatedAt = recipeStepIngredient.CreatedAt
	exampleRecipeStepIngredient.Ingredient.CreatedAt = recipeStepIngredient.Ingredient.CreatedAt
	exampleRecipeStepIngredient.Ingredient = recipeStepIngredient.Ingredient
	exampleRecipeStepIngredient.MeasurementUnit = recipeStepIngredient.MeasurementUnit

	require.Equal(t, exampleRecipeStepIngredient, recipeStepIngredient)

	assert.NoError(t, err)
	assert.Equal(t, recipeStepIngredient, exampleRecipeStepIngredient)

	return created
}

func TestQuerier_Integration_RecipeStepIngredients(t *testing.T) {
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
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
	exampleRecipeStep := createdRecipe.Steps[0]

	validMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	validIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
	exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
	exampleRecipeStepIngredient.Ingredient = validIngredient
	exampleRecipeStepIngredient.MeasurementUnit = *validMeasurementUnit
	exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
	createdRecipeStepIngredients := []*types.RecipeStepIngredient{
		exampleRecipeStep.Ingredients[0],
	}

	// create
	createdRecipeStepIngredients = append(createdRecipeStepIngredients, createRecipeStepIngredientForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepIngredient, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		validMeasurementUnit = createValidMeasurementUnitForTest(t, ctx, nil, dbc)
		validIngredient = createValidIngredientForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeRecipeStepIngredient()
		input.Ingredient = validIngredient
		input.MeasurementUnit = *validMeasurementUnit
		input.BelongsToRecipeStep = exampleRecipeStep.ID
		createdRecipeStepIngredients = append(createdRecipeStepIngredients, createRecipeStepIngredientForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeStepIngredients, err := dbc.GetRecipeStepIngredients(ctx, exampleRecipe.ID, createdRecipeStepIngredients[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepIngredients.Data)
	assert.Equal(t, len(createdRecipeStepIngredients), len(recipeStepIngredients.Data))

	// delete
	for _, recipeStepIngredient := range createdRecipeStepIngredients {
		assert.NoError(t, dbc.ArchiveRecipeStepIngredient(ctx, exampleRecipeStep.ID, recipeStepIngredient.ID))

		var exists bool
		exists, err = dbc.RecipeStepIngredientExists(ctx, exampleRecipe.ID, recipeStepIngredient.BelongsToRecipeStep, recipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepIngredient
		y, err = dbc.GetRecipeStepIngredient(ctx, exampleRecipe.ID, recipeStepIngredient.BelongsToRecipeStep, recipeStepIngredient.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
