package mealplanning

import (
	"context"
	"database/sql"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildRecipeStepForTestCreation(t *testing.T, ctx context.Context, recipeID string, dbc *repository) *types.RecipeStep {
	t.Helper()

	recipeStepID := identifiers.New()

	validIngredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)

	preparation := createValidPreparationForTest(t, ctx, nil, dbc)
	ingredient := createValidIngredientForTest(t, ctx, nil, dbc)
	measurementUnit1 := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	recipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
	recipeStepIngredient.Ingredient = ingredient
	recipeStepIngredient.MeasurementUnit = *measurementUnit1
	recipeStepIngredient.BelongsToRecipeStep = recipeStepID

	instrument := createValidInstrumentForTest(t, ctx, nil, dbc)
	recipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
	recipeStepInstrument.Instrument = instrument
	recipeStepInstrument.BelongsToRecipeStep = recipeStepID

	measurementUnit2 := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	exampleVessel := fakes.BuildFakeValidVessel()
	exampleVessel.CapacityUnit = measurementUnit2
	vessel := createValidVesselForTest(t, ctx, exampleVessel, dbc)
	recipeStepVessel := fakes.BuildFakeRecipeStepVessel()
	recipeStepVessel.Vessel = vessel
	recipeStepVessel.BelongsToRecipeStep = recipeStepID

	exampleRecipeMedia := fakes.BuildFakeRecipeMedia()
	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = recipeID
	exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			ID:                      identifiers.New(),
			BelongsToRecipeStep:     recipeStepID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}

	exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
	exampleRecipeStepProduct.BelongsToRecipeStep = recipeStepID
	exampleRecipeStepProduct.MeasurementUnit = measurementUnit1

	exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
	exampleRecipeStepCompletionCondition.BelongsToRecipeStep = recipeStepID
	exampleRecipeStepCompletionCondition.IngredientState = *validIngredientState
	exampleRecipeStepCompletionCondition.Ingredients = []*types.RecipeStepCompletionConditionIngredient{
		{
			ID:                                     identifiers.New(),
			BelongsToRecipeStepCompletionCondition: exampleRecipeStepCompletionCondition.ID,
			RecipeStepIngredient:                   recipeStepIngredient.ID,
		},
	}

	exampleRecipeStep := fakes.BuildFakeRecipeStep()
	exampleRecipeStep.ID = recipeStepID
	exampleRecipeStep.Index = 0
	exampleRecipeStep.BelongsToRecipe = recipeID
	exampleRecipeStep.Preparation = *preparation
	exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{
		recipeStepIngredient,
	}
	exampleRecipeStep.Instruments = []*types.RecipeStepInstrument{
		recipeStepInstrument,
	}
	exampleRecipeStep.Vessels = []*types.RecipeStepVessel{
		recipeStepVessel,
	}
	exampleRecipeStep.Media = []*types.RecipeMedia{
		exampleRecipeMedia,
	}
	exampleRecipeStep.Products = []*types.RecipeStepProduct{
		exampleRecipeStepProduct,
	}
	exampleRecipeStep.CompletionConditions = []*types.RecipeStepCompletionCondition{
		exampleRecipeStepCompletionCondition,
	}

	return exampleRecipeStep
}

func createRecipeStepForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStep *types.RecipeStep, dbc *repository) *types.RecipeStep {
	t.Helper()

	// create
	if exampleRecipeStep == nil {
		exampleRecipeStep = fakes.BuildFakeRecipeStep()
	}
	dbInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

	created, err := dbc.CreateRecipeStep(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStep.Media = nil
	exampleRecipeStep.CreatedAt = created.CreatedAt
	exampleRecipeStep.Preparation = types.ValidPreparation{ID: created.Preparation.ID}

	for j := range created.Ingredients {
		exampleRecipeStep.Ingredients[j].CreatedAt = created.Ingredients[j].CreatedAt

		assert.Equal(t, exampleRecipeStep.Ingredients[j].Ingredient.ID, created.Ingredients[j].Ingredient.ID)
		exampleRecipeStep.Ingredients[j].Ingredient = created.Ingredients[j].Ingredient

		assert.Equal(t, exampleRecipeStep.Ingredients[j].MeasurementUnit.ID, created.Ingredients[j].MeasurementUnit.ID)
		exampleRecipeStep.Ingredients[j].MeasurementUnit = created.Ingredients[j].MeasurementUnit
	}

	for j := range created.Instruments {
		exampleRecipeStep.Instruments[j].CreatedAt = created.Instruments[j].CreatedAt

		assert.Equal(t, exampleRecipeStep.Instruments[j].Instrument.ID, created.Instruments[j].Instrument.ID)
		exampleRecipeStep.Instruments[j].Instrument = created.Instruments[j].Instrument
	}

	for j := range created.Vessels {
		exampleRecipeStep.Vessels[j].CreatedAt = created.Vessels[j].CreatedAt

		assert.Equal(t, exampleRecipeStep.Vessels[j].Vessel.ID, created.Vessels[j].Vessel.ID)
		exampleRecipeStep.Vessels[j].Vessel = created.Vessels[j].Vessel

		if created.Vessels[j].Vessel.CapacityUnit != nil {
			assert.Equal(t, exampleRecipeStep.Vessels[j].Vessel.CapacityUnit.ID, created.Vessels[j].Vessel.CapacityUnit.ID)
			exampleRecipeStep.Vessels[j].Vessel.CapacityUnit = created.Vessels[j].Vessel.CapacityUnit
		}
	}

	for j := range created.CompletionConditions {
		exampleRecipeStep.CompletionConditions[j].CreatedAt = created.CompletionConditions[j].CreatedAt

		assert.Equal(t, exampleRecipeStep.CompletionConditions[j].IngredientState.ID, created.CompletionConditions[j].IngredientState.ID)
		exampleRecipeStep.CompletionConditions[j].IngredientState = created.CompletionConditions[j].IngredientState

		// Update the CreatedAt for each ingredient in the completion condition
		for k := range created.CompletionConditions[j].Ingredients {
			if k < len(exampleRecipeStep.CompletionConditions[j].Ingredients) {
				exampleRecipeStep.CompletionConditions[j].Ingredients[k].CreatedAt = created.CompletionConditions[j].Ingredients[k].CreatedAt
			}
		}
	}

	for j := range created.Products {
		exampleRecipeStep.Products[j].CreatedAt = created.Products[j].CreatedAt
		assert.Equal(t, exampleRecipeStep.Products[j].MeasurementUnit.ID, created.Products[j].MeasurementUnit.ID)
		exampleRecipeStep.Products[j].MeasurementUnit = created.Products[j].MeasurementUnit
	}

	for j := range created.Media {
		assert.Equal(t, exampleRecipeStep.Media[j].ID, created.Media[j].ID)
		exampleRecipeStep.Media[j] = created.Media[j]
	}

	assert.Equal(t, exampleRecipeStep.Products, created.Products)
	assert.Equal(t, exampleRecipeStep.Instruments, created.Instruments)
	assert.Equal(t, exampleRecipeStep.Vessels, created.Vessels)
	assert.Equal(t, exampleRecipeStep.Ingredients, created.Ingredients)
	assert.Equal(t, exampleRecipeStep.Media, created.Media)
	assert.Equal(t, exampleRecipeStep.CompletionConditions, created.CompletionConditions)

	exampleRecipeStep.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleRecipeStep, created)

	recipeStep, err := dbc.GetRecipeStep(ctx, recipeID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipeStep)

	assert.Equal(t, exampleRecipeStep.Preparation.ID, recipeStep.Preparation.ID)
	exampleRecipeStep.Preparation = recipeStep.Preparation
	exampleRecipeStep.CreatedAt = recipeStep.CreatedAt

	require.Equal(t, exampleRecipeStep, recipeStep)

	assert.NoError(t, err)
	assert.Equal(t, recipeStep, exampleRecipeStep)

	return created
}

func TestQuerier_Integration_RecipeSteps(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := t.Context()
	dbc, container := buildDatabaseClientForTest(t)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := pgtesting.CreateUserForTest(t, nil, dbc.db)

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)

	exampleRecipeStep := buildRecipeStepForTestCreation(t, ctx, createdRecipe.ID, dbc)
	exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
	createdRecipeSteps := []*types.RecipeStep{
		exampleRecipe.Steps[0],
	}

	// create
	createdRecipeSteps = append(createdRecipeSteps, createRecipeStepForTest(t, ctx, exampleRecipe.ID, exampleRecipeStep, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := buildRecipeStepForTestCreation(t, ctx, createdRecipe.ID, dbc)
		input.BelongsToRecipe = createdRecipe.ID
		createdRecipeSteps = append(createdRecipeSteps, createRecipeStepForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeSteps, err := dbc.GetRecipeSteps(ctx, exampleRecipe.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeSteps.Data)
	assert.Equal(t, len(createdRecipeSteps), len(recipeSteps.Data))

	// delete
	for _, recipeStep := range createdRecipeSteps {
		assert.NoError(t, dbc.ArchiveRecipeStep(ctx, exampleRecipe.ID, recipeStep.ID))

		var exists bool
		exists, err = dbc.RecipeStepExists(ctx, exampleRecipe.ID, recipeStep.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStep
		y, err = dbc.GetRecipeStep(ctx, exampleRecipe.ID, recipeStep.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_RecipeStepExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		c := buildInertClientForTest(t)

		actual, err := c.RecipeStepExists(ctx, "", exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		exampleRecipeID := fakes.BuildFakeID()

		c := buildInertClientForTest(t)

		actual, err := c.RecipeStepExists(ctx, exampleRecipeID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.GetRecipeStep(ctx, "", exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.GetRecipeStep(ctx, exampleRecipeID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getRecipeStepByID(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.getRecipeStepByID(ctx, c.db, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := filtering.DefaultQueryFilter()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.GetRecipeSteps(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.CreateRecipeStep(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		actual, err := c.createRecipeStep(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateRecipeStep(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveRecipeStep(ctx, "", exampleRecipeStep.ID))
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := t.Context()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveRecipeStep(ctx, exampleRecipeID, ""))
	})
}
