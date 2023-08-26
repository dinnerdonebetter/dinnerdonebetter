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

func createRecipeStepProductForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepProduct *types.RecipeStepProduct, dbc *Querier) *types.RecipeStepProduct {
	t.Helper()

	// create
	if exampleRecipeStepProduct == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
		exampleRecipeStep := createdRecipe.Steps[0]

		exampleRecipeStepProduct = fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
	}
	dbInput := converters.ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(exampleRecipeStepProduct)

	created, err := dbc.CreateRecipeStepProduct(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepProduct.CreatedAt = created.CreatedAt
	exampleRecipeStepProduct.MeasurementUnit = created.MeasurementUnit
	assert.Equal(t, exampleRecipeStepProduct, created)

	recipeStepProduct, err := dbc.GetRecipeStepProduct(ctx, recipeID, exampleRecipeStepProduct.BelongsToRecipeStep, exampleRecipeStepProduct.ID)

	exampleRecipeStepProduct.CreatedAt = recipeStepProduct.CreatedAt
	exampleRecipeStepProduct.MeasurementUnit = recipeStepProduct.MeasurementUnit

	require.Equal(t, exampleRecipeStepProduct, recipeStepProduct)

	assert.NoError(t, err)
	assert.Equal(t, recipeStepProduct, exampleRecipeStepProduct)

	return created
}

func TestQuerier_Integration_RecipeStepProducts(t *testing.T) {
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
	exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
	exampleRecipeStepProduct.MeasurementUnit = validMeasurementUnit
	exampleRecipeStepProduct.BelongsToRecipeStep = exampleRecipeStep.ID
	createdRecipeStepProducts := []*types.RecipeStepProduct{
		exampleRecipeStep.Products[0],
	}

	// create
	createdRecipeStepProducts = append(createdRecipeStepProducts, createRecipeStepProductForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepProduct, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		validMeasurementUnit = createValidMeasurementUnitForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeRecipeStepProduct()
		input.MeasurementUnit = validMeasurementUnit
		input.BelongsToRecipeStep = exampleRecipeStep.ID
		createdRecipeStepProducts = append(createdRecipeStepProducts, createRecipeStepProductForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeStepProducts, err := dbc.GetRecipeStepProducts(ctx, exampleRecipe.ID, createdRecipeStepProducts[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepProducts.Data)
	assert.Equal(t, len(createdRecipeStepProducts), len(recipeStepProducts.Data))

	// delete
	for _, recipeStepProduct := range createdRecipeStepProducts {
		assert.NoError(t, dbc.ArchiveRecipeStepProduct(ctx, exampleRecipeStep.ID, recipeStepProduct.ID))

		var exists bool
		exists, err = dbc.RecipeStepProductExists(ctx, exampleRecipe.ID, recipeStepProduct.BelongsToRecipeStep, recipeStepProduct.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepProduct
		y, err = dbc.GetRecipeStepProduct(ctx, exampleRecipe.ID, recipeStepProduct.BelongsToRecipeStep, recipeStepProduct.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
