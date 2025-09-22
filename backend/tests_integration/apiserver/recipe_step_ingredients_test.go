package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepIngredientSliceEquality(t *testing.T, stepIndex int, expected, actual []*mealplanning.RecipeStepIngredient) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "expected recipe step %d ingredients length", stepIndex)
	for i := range expected {
		checkRecipeStepIngredientEquality(t, stepIndex, i, expected[i], actual[i])
	}
}

func checkRecipeStepIngredientEquality(t *testing.T, stepIndex, ingIndex int, expected, actual *mealplanning.RecipeStepIngredient) {
	t.Helper()
	assert.NotEmpty(t, actual.ID, "expected step %d ingredient %d to have ID", stepIndex, ingIndex)
	assert.False(t, actual.CreatedAt.IsZero(), "expected step %d ingredient %d to have CreatedAt", stepIndex, ingIndex)
	assert.NotEmpty(t, actual.BelongsToRecipeStep, "expected step %d ingredient %d to have BelongsToRecipeStep", stepIndex, ingIndex)
	assert.Equal(t, expected.Name, actual.Name, "expected step %d ingredient %d Name", stepIndex, ingIndex)
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected step %d ingredient %d Quantity", stepIndex, ingIndex)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected step %d ingredient %d QuantityNotes", stepIndex, ingIndex)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected step %d ingredient %d IngredientNotes", stepIndex, ingIndex)
	assert.Equal(t, expected.OptionIndex, actual.OptionIndex, "expected step %d ingredient %d OptionIndex", stepIndex, ingIndex)
	assert.Equal(t, expected.Optional, actual.Optional, "expected step %d ingredient %d Optional", stepIndex, ingIndex)
	assert.Equal(t, expected.ToTaste, actual.ToTaste, "expected step %d ingredient %d ToTaste", stepIndex, ingIndex)
	if expected.VesselIndex != nil {
		require.NotNil(t, actual.VesselIndex, "expected step %d ingredient %d VesselIndex non-nil", stepIndex, ingIndex)
		assert.Equal(t, *expected.VesselIndex, *actual.VesselIndex, "expected step %d ingredient %d VesselIndex", stepIndex, ingIndex)
	}
	if expected.ProductPercentageToUse != nil {
		require.NotNil(t, actual.ProductPercentageToUse, "expected step %d ingredient %d ProductPercentageToUse non-nil", stepIndex, ingIndex)
		assert.Equal(t, *expected.ProductPercentageToUse, *actual.ProductPercentageToUse, "expected step %d ingredient %d ProductPercentageToUse", stepIndex, ingIndex)
	}
	// MeasurementUnit comparison by ID (and ranges already compared above)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected step %d ingredient %d MeasurementUnit.ID", stepIndex, ingIndex)
	// Ingredient pointer may be nil if this ingredient refers to a product of a prior step
	if expected.Ingredient != nil {
		require.NotNil(t, actual.Ingredient, "expected step %d ingredient %d Ingredient non-nil", stepIndex, ingIndex)
		assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected step %d ingredient %d Ingredient.ID", stepIndex, ingIndex)
	}
}

func TestRecipeStepIngredients_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		firstStep := createdRecipe.Steps[0]
		createdRecipeStepID := firstStep.ID
		createdRecipeStepIngredientID := firstStep.Ingredients[0].ID

		require.NotEmpty(t, createdRecipeStepID, "created recipe step ID must not be empty")
		require.NotEmpty(t, createdRecipeStepIngredientID, "created recipe step ingredient ID must not be empty")

		createdRecipeStepIngredientRes, err := userClient.GetRecipeStepIngredient(ctx, &mealplanninggrpc.GetRecipeStepIngredientRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepIngredientID: createdRecipeStepIngredientID,
		})
		require.NotNil(t, createdRecipeStepIngredientRes)
		require.NoError(t, err)

		createdRecipeStepIngredient := converters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(createdRecipeStepIngredientRes.Result)

		createdValidIngredient := createValidIngredientForTest(t)

		newRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		newRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
		newRecipeStepIngredient.ID = createdRecipeStepIngredientID
		newRecipeStepIngredient.Ingredient = createdValidIngredient
		newRecipeStepIngredient.MeasurementUnit = createdRecipeStepIngredient.MeasurementUnit

		updateInput := mpconverters.ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(newRecipeStepIngredient)
		createdRecipeStepIngredient.Update(updateInput)

		_, err = adminClient.UpdateRecipeStepIngredient(ctx, &mealplanninggrpc.UpdateRecipeStepIngredientRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepIngredientID: createdRecipeStepIngredientID,
			Input:                  converters.ConvertRecipeStepIngredientUpdateRequestInputToGRPCRecipeStepIngredientUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		retrievedRes, err := userClient.GetRecipeStepIngredient(ctx, &mealplanninggrpc.GetRecipeStepIngredientRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepIngredientID: createdRecipeStepIngredientID,
		})
		require.NotNil(t, retrievedRes)
		require.NoError(t, err)

		actual := converters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(retrievedRes.Result)

		// assert recipe step ingredient equality
		checkRecipeStepIngredientEquality(t, -1, -1, newRecipeStepIngredient, actual)
		assert.NotNil(t, actual.LastUpdatedAt)

		_, err = userClient.ArchiveRecipeStepIngredient(ctx, &mealplanninggrpc.ArchiveRecipeStepIngredientRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepIngredientID: createdRecipeStepIngredientID,
		})
		assert.NoError(t, err)

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeID: createdRecipe.ID, RecipeStepID: createdRecipeStepID})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipeStepIngredients_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var (
			createdRecipeStepID string
		)
		for _, step := range createdRecipe.Steps {
			createdRecipeStepID = step.ID
			break
		}

		createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)

		var expected []*mealplanning.RecipeStepIngredient
		for i := 0; i < 5; i++ {
			x, _, _ := createRecipeForTest(t, nil)

			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepIngredient.Ingredient = &mealplanning.ValidIngredient{ID: x[0].ID}
			exampleRecipeStepIngredient.MeasurementUnit = mealplanning.ValidMeasurementUnit{ID: createdValidMeasurementUnit.ID}

			exampleRecipeStepIngredientInput := mpconverters.ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(exampleRecipeStepIngredient)
			createdRecipeStepIngredientRes, err := adminClient.CreateRecipeStepIngredient(ctx, &mealplanninggrpc.CreateRecipeStepIngredientRequest{
				RecipeID:     createdRecipe.ID,
				RecipeStepID: createdRecipeStepID,
				Input:        converters.ConvertRecipeStepIngredientCreationRequestInputToGRPCRecipeStepIngredientCreationRequestInput(exampleRecipeStepIngredientInput),
			})
			require.NoError(t, err)

			createdRecipeStepIngredient := converters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(createdRecipeStepIngredientRes.Created)
			checkRecipeStepIngredientEquality(t, -1, -1, exampleRecipeStepIngredient, createdRecipeStepIngredient)

			retrievedRecipeStepIngredientRes, err := userClient.GetRecipeStepIngredient(ctx, &mealplanninggrpc.GetRecipeStepIngredientRequest{
				RecipeID:               createdRecipe.ID,
				RecipeStepID:           createdRecipeStepID,
				RecipeStepIngredientID: createdRecipeStepIngredient.ID,
			})
			require.NotNil(t, retrievedRecipeStepIngredientRes.Result)
			require.NoError(t, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepIngredient.BelongsToRecipeStep)

			expected = append(expected, createdRecipeStepIngredient)
		}

		// assert recipe step ingredient list equality
		actual, err := userClient.GetRecipeStepIngredients(ctx, &mealplanninggrpc.GetRecipeStepIngredientsRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStepID,
		})
		require.NotNil(t, actual)
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipeStepIngredient := range expected {
			_, err = userClient.ArchiveRecipeStepIngredient(ctx, &mealplanninggrpc.ArchiveRecipeStepIngredientRequest{
				RecipeID:               createdRecipe.ID,
				RecipeStepID:           createdRecipeStepID,
				RecipeStepIngredientID: createdRecipeStepIngredient.ID,
			})
			assert.NoError(t, err)
		}

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStepID,
		})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}
