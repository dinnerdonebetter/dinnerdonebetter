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

/*

func (s *TestSuite) TestRecipeStepInstruments_AsRecipeStepProducts() {
	s.runTest("should be able to use a recipe step instrument that was the product of a prior recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			lineBase := fakes.BuildFakeValidPreparation()
			lineInput := mpconverters.ConvertValidPreparationToValidPreparationCreationRequestInput(lineBase)
			line, err := adminClient.CreateValidPreparation(ctx, lineInput)
			require.NoError(t, err)

			roastBase := fakes.BuildFakeValidPreparation()
			roastInput := mpconverters.ConvertValidPreparationToValidPreparationCreationRequestInput(roastBase)
			roast, err := adminClient.CreateValidPreparation(ctx, roastInput)
			require.NoError(t, err)

			bakingSheetBase := fakes.BuildFakeValidInstrument()
			bakingSheetBaseInput := mpconverters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(bakingSheetBase)
			bakingSheet, err := adminClient.CreateValidInstrument(ctx, bakingSheetBaseInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, bakingSheetBase, bakingSheet)

			sheetsBase := fakes.BuildFakeValidMeasurementUnit()
			sheetsBaseInput := mpconverters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(sheetsBase)
			sheets, err := adminClient.CreateValidMeasurementUnit(ctx, sheetsBaseInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, sheetsBase, sheets)

			headsBase := fakes.BuildFakeValidMeasurementUnit()
			headsBaseInput := mpconverters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(headsBase)
			head, err := adminClient.CreateValidMeasurementUnit(ctx, headsBaseInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, headsBase, head)

			exampleUnits := fakes.BuildFakeValidMeasurementUnit()
			exampleUnitsInput := mpconverters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleUnits)
			unit, err := adminClient.CreateValidMeasurementUnit(ctx, exampleUnitsInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleUnits, unit)

			aluminumFoilBase := fakes.BuildFakeValidIngredient()
			aluminumFoilInput := mpconverters.ConvertValidIngredientToValidIngredientCreationRequestInput(aluminumFoilBase)
			aluminumFoil, createdValidIngredientErr := adminClient.CreateValidIngredient(ctx, aluminumFoilInput)
			require.NoError(t, createdValidIngredientErr)

			garlic := fakes.BuildFakeValidIngredient()
			garlicInput := mpconverters.ConvertValidIngredientToValidIngredientCreationRequestInput(garlic)
			garlic, garlicErr := adminClient.CreateValidIngredient(ctx, garlicInput)
			require.NoError(t, garlicErr)

			linedBakingSheetName := "lined baking sheet"

			expected := &mealplanning.Recipe{
				Name:                t.Name(),
				Slug:                "whatever-who-cares",
				YieldsComponentType: mealplanning.MealComponentTypesMain,
				PortionName:         t.Name(),
				PluralPortionName:   t.Name(),
				EstimatedPortions: mealplanning.Float32RangeWithOptionalMax{
					Max: nil,
					Min: 1,
				},
				Steps: []*mealplanning.RecipeStep{
					{
						Products: []*mealplanning.RecipeStepProduct{
							{
								Name:            linedBakingSheetName,
								Type:            mealplanning.RecipeStepProductInstrumentType,
								MeasurementUnit: unit,
								QuantityNotes:   "",
								Quantity:        mealplanning.OptionalFloat32Range{Min: pointer.To(float32(1))},
							},
						},
						Notes:       "first step",
						Preparation: *line,
						Ingredients: []*mealplanning.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          aluminumFoil,
								Name:                "aluminum foil",
								MeasurementUnit:     *sheets,
								Quantity: mealplanning.Float32RangeWithOptionalMax{
									Max: nil,
									Min: 3,
								},
							},
						},
						Instruments: []*mealplanning.RecipeStepInstrument{
							{
								Instrument: bakingSheet,
							},
						},
						Index: 0,
					},
					{
						Preparation: *roast,
						Instruments: []*mealplanning.RecipeStepInstrument{
							{
								Name:       linedBakingSheetName,
								Instrument: nil,
							},
						},
						Products: []*mealplanning.RecipeStepProduct{
							{
								Name:            "roasted garlic",
								Type:            mealplanning.RecipeStepProductIngredientType,
								MeasurementUnit: head,
								QuantityNotes:   "",
								Quantity:        mealplanning.OptionalFloat32Range{Min: pointer.To(float32(1))},
							},
						},
						Notes: "second step",
						Ingredients: []*mealplanning.RecipeStepIngredient{
							{
								Ingredient:      garlic,
								Name:            "garlic",
								MeasurementUnit: *head,
								Quantity: mealplanning.Float32RangeWithOptionalMax{
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

			created, err := adminClient.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			checkRecipeEquality(t, expected, created)

			created, err = userClient.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
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
		}
	})
}

func (s *TestSuite) TestRecipeStepInstruments_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, adminClient, userClient, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := mpconverters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := adminClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			var expected []*mealplanning.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrument.Instrument = &mealplanning.ValidInstrument{ID: createdValidInstrument.ID}
				exampleRecipeStepInstrumentInput := mpconverters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrument)
				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := adminClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr = userClient.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := userClient.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipeStepInstrument := range expected {
				assert.NoError(t, userClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			assert.NoError(t, userClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			assert.NoError(t, adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}


*/
