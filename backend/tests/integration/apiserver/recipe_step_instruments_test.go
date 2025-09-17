package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepInstrumentSliceEquality(t *testing.T, stepIndex int, expected, actual []*mealplanning.RecipeStepInstrument) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "expected recipe step %d instruments length", stepIndex)
	for i := range expected {
		checkRecipeStepInstrumentEquality(t, stepIndex, i, expected[i], actual[i])
	}
}

func checkRecipeStepInstrumentEquality(t *testing.T, stepIndex, instrIndex int, expected, actual *mealplanning.RecipeStepInstrument) {
	t.Helper()
	assert.NotEmpty(t, actual.ID, "expected step %d instrument %d to have ID", stepIndex, instrIndex)
	assert.False(t, actual.CreatedAt.IsZero(), "expected step %d instrument %d to have CreatedAt", stepIndex, instrIndex)
	assert.NotEmpty(t, actual.BelongsToRecipeStep, "expected step %d instrument %d to have BelongsToRecipeStep", stepIndex, instrIndex)
	assert.Equal(t, expected.Name, actual.Name, "expected step %d instrument %d Name", stepIndex, instrIndex)
	assert.Equal(t, expected.Notes, actual.Notes, "expected step %d instrument %d Notes", stepIndex, instrIndex)
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected step %d instrument %d Quantity", stepIndex, instrIndex)
	assert.Equal(t, expected.OptionIndex, actual.OptionIndex, "expected step %d instrument %d OptionIndex", stepIndex, instrIndex)
	assert.Equal(t, expected.PreferenceRank, actual.PreferenceRank, "expected step %d instrument %d PreferenceRank", stepIndex, instrIndex)
	assert.Equal(t, expected.Optional, actual.Optional, "expected step %d instrument %d Optional", stepIndex, instrIndex)
	if expected.Instrument != nil {
		require.NotNil(t, actual.Instrument, "expected step %d instrument %d Instrument non-nil", stepIndex, instrIndex)
		assert.NotEmpty(t, actual.Instrument.ID, "expected step %d instrument %d Instrument.ID", stepIndex, instrIndex)
		assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected step %d instrument %d Instrument.ID", stepIndex, instrIndex)
	}
}

func TestRecipeStepInstruments_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var createdRecipeStepID string
		for _, step := range createdRecipe.Steps {
			createdRecipeStepID = step.ID
			break
		}

		createdValidInstrument := createValidInstrumentForTest(t)

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
		exampleRecipeStepInstrument.Instrument = &mealplanning.ValidInstrument{ID: createdValidInstrument.ID}
		exampleRecipeStepInstrumentInput := mpconverters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrument)
		createdRecipeStepInstrumentRes, err := adminClient.CreateRecipeStepInstrument(ctx, &mealplanninggrpc.CreateRecipeStepInstrumentRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStepID,
			Input:        converters.ConvertRecipeStepInstrumentCreationRequestInputToGRPCRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrumentInput),
		})
		require.NoError(t, err)

		createdRecipeStepInstrument := converters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(createdRecipeStepInstrumentRes.Created)
		checkRecipeStepInstrumentEquality(t, -1, -1, exampleRecipeStepInstrument, createdRecipeStepInstrument)

		retrievedRecipeStepInstrumentRes, err := userClient.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
		})
		require.NoError(t, err)

		createdRecipeStepInstrument = converters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(retrievedRecipeStepInstrumentRes.Result)
		checkRecipeStepInstrumentEquality(t, -1, -1, exampleRecipeStepInstrument, createdRecipeStepInstrument)

		require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
		exampleRecipeStepInstrument.Instrument = createdValidInstrument
		exampleRecipeStepInstrument.Instrument.CreatedAt = createdRecipeStepInstrument.Instrument.CreatedAt

		checkRecipeStepInstrumentEquality(t, -1, -1, exampleRecipeStepInstrument, createdRecipeStepInstrument)

		newValidInstrument := createValidInstrumentForTest(t)

		newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
		newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
		newRecipeStepInstrument.Instrument = newValidInstrument
		updateInput := mpconverters.ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(newRecipeStepInstrument)
		createdRecipeStepInstrument.Update(updateInput)

		_, err = adminClient.UpdateRecipeStepInstrument(ctx, &mealplanninggrpc.UpdateRecipeStepInstrumentRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
			Input:                  converters.ConvertRecipeStepInstrumentUpdateRequestInputToGRPCRecipeStepInstrumentUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		retrievedRecipeStepInstrumentRes, err = userClient.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
		})
		require.NoError(t, err)

		actual := converters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(retrievedRecipeStepInstrumentRes.Result)

		// assert recipe step instrument equality
		checkRecipeStepInstrumentEquality(t, -1, -1, newRecipeStepInstrument, actual)
		assert.NotNil(t, actual.LastUpdatedAt)

		_, err = userClient.ArchiveRecipeStepInstrument(ctx, &mealplanninggrpc.ArchiveRecipeStepInstrumentRequest{
			RecipeID:               createdRecipe.ID,
			RecipeStepID:           createdRecipeStepID,
			RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
		})
		assert.NoError(t, err)

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeID: createdRecipe.ID, RecipeStepID: createdRecipeStepID})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipeStepInstruments_AsRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("should be able to use a recipe step instrument that was the product of a prior recipe step", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		knife := createValidInstrumentForTest(t)
		preheat := createValidPreparationForTest(t)
		cut := createValidPreparationForTest(t)
		stick := createValidMeasurementUnitForTest(t)
		unit := createValidMeasurementUnitForTest(t)
		butter := createValidIngredientForTest(t)

		preheatedKnife := "preheated knife"

		expected := &mealplanning.Recipe{
			Name:                t.Name(),
			Slug:                "whatever-who-cares",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: nil,
				Min: 1,
			},
			Steps: []*mealplanning.RecipeStep{
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            preheatedKnife,
							Type:            mealplanning.RecipeStepProductInstrumentType,
							MeasurementUnit: unit,
							QuantityNotes:   "",
							Quantity:        types.OptionalFloat32Range{Min: pointer.To(float32(1))},
						},
					},
					Notes:       "first step",
					Preparation: *preheat,
					Ingredients: []*mealplanning.RecipeStepIngredient{},
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Instrument: knife,
						},
					},
					Index: 0,
				},
				{
					Preparation: *cut,
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       preheatedKnife,
							Instrument: nil,
						},
					},
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "cut butter",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: stick,
							QuantityNotes:   "",
							Quantity:        types.OptionalFloat32Range{Min: pointer.To(float32(1))},
						},
					},
					Notes: "second step",
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient:      butter,
							Name:            "butter",
							MeasurementUnit: *stick,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: 1,
							},
						},
					},
					Index: 1,
				},
			},
		}

		exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(expected)
		exampleRecipeInput.Steps[1].Instruments[0].ProductOfRecipeStepIndex = pointer.To(uint64(0))
		exampleRecipeInput.Steps[1].Instruments[0].ProductOfRecipeStepProductIndex = pointer.To(uint64(0))

		createdRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput)})
		require.NoError(t, err)

		created := converters.ConvertGRPCRecipeToRecipe(createdRes.Created)
		checkRecipeEquality(t, expected, created)

		fetchedRes, err := userClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: created.ID})
		require.NoError(t, err)

		created = converters.ConvertGRPCRecipeToRecipe(fetchedRes.Result)
		checkRecipeEquality(t, expected, created)

		recipeStepProductIndex := -1
		for i, ingredient := range created.Steps[1].Instruments {
			if ingredient.RecipeStepProductID != nil {
				recipeStepProductIndex = i
			}
		}

		require.NotEqual(t, -1, recipeStepProductIndex)
		require.NotNil(t, created.Steps[1].Instruments[recipeStepProductIndex].RecipeStepProductID)
		assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Instruments[recipeStepProductIndex].RecipeStepProductID)
	})
}

func TestRecipeStepInstruments_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var createdRecipeStepID string
		for _, step := range createdRecipe.Steps {
			createdRecipeStepID = step.ID
			break
		}

		createdValidInstrument := createValidInstrumentForTest(t)

		var expected []*mealplanning.RecipeStepInstrument
		for i := 0; i < 5; i++ {
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrument.Instrument = &mealplanning.ValidInstrument{ID: createdValidInstrument.ID}
			exampleRecipeStepInstrumentInput := mpconverters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentRes, err := adminClient.CreateRecipeStepInstrument(ctx, &mealplanninggrpc.CreateRecipeStepInstrumentRequest{
				RecipeID:     createdRecipe.ID,
				RecipeStepID: createdRecipeStepID,
				Input:        converters.ConvertRecipeStepInstrumentCreationRequestInputToGRPCRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrumentInput),
			})
			require.NoError(t, err)

			createdRecipeStepInstrument := converters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(createdRecipeStepInstrumentRes.Created)
			checkRecipeStepInstrumentEquality(t, i, i, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			fetchedRecipeStepInstrumentRes, err := userClient.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
				RecipeID:               createdRecipe.ID,
				RecipeStepID:           createdRecipeStepID,
				RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
			})
			require.NoError(t, err)

			createdRecipeStepInstrument = converters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(fetchedRecipeStepInstrumentRes.Result)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

			expected = append(expected, createdRecipeStepInstrument)
		}

		// assert recipe step instrument list equality
		actual, err := userClient.GetRecipeStepInstruments(ctx, &mealplanninggrpc.GetRecipeStepInstrumentsRequest{
			RecipeID:     createdRecipe.ID,
			RecipeStepID: createdRecipeStepID,
		})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipeStepInstrument := range expected {
			_, err = userClient.ArchiveRecipeStepInstrument(ctx, &mealplanninggrpc.ArchiveRecipeStepInstrumentRequest{
				RecipeID:               createdRecipe.ID,
				RecipeStepID:           createdRecipeStepID,
				RecipeStepInstrumentID: createdRecipeStepInstrument.ID,
			})
			assert.NoError(t, err)
		}

		_, err = userClient.ArchiveRecipeStep(ctx, &mealplanninggrpc.ArchiveRecipeStepRequest{RecipeID: createdRecipe.ID, RecipeStepID: createdRecipeStepID})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}
