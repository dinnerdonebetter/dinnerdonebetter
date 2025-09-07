package integration

import (
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeEquality(t *testing.T, expected, actual *mealplanning.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)

	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe %s to be %v, but it was %v", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.Zero(t, actual.LastUpdatedAt)
	assert.Zero(t, actual.ArchivedAt)
	assert.Equal(t, expected.EstimatedPortions, actual.EstimatedPortions, "expected EstimatedPortions for recipe %s to be %v, but it was %v", expected.ID, expected.EstimatedPortions, actual.EstimatedPortions)
	assert.Equal(t, expected.PluralPortionName, actual.PluralPortionName, "expected PluralPortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PluralPortionName, actual.PluralPortionName)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.PortionName, actual.PortionName, "expected PortionName for recipe %s to be %v, but it was %v", expected.ID, expected.PortionName, actual.PortionName)
	assert.NotZero(t, actual.CreatedByUser)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe %s to be %v, but it was %v", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for recipe %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.YieldsComponentType, actual.YieldsComponentType, "expected YieldsComponentType for recipe %s to be %v, but it was %v", expected.ID, expected.YieldsComponentType, actual.YieldsComponentType)
	assert.Equal(t, expected.PrepTasks, actual.PrepTasks, "expected PrepTasks for recipe %s to be %v, but it was %v", expected.ID, expected.PrepTasks, actual.PrepTasks)
	//assert.Equal(t, expected.Steps, actual.Steps, "expected Steps for recipe %s to be %v, but it was %v", expected.ID, expected.Steps, actual.Steps)
	assert.Equal(t, expected.Media, actual.Media, "expected Media for recipe %s to be %v, but it was %v", expected.ID, expected.Media, actual.Media)
	assert.Equal(t, expected.SealOfApproval, actual.SealOfApproval, "expected SealOfApproval for recipe %s to be %v, but it was %v", expected.ID, expected.SealOfApproval, actual.SealOfApproval)
	assert.Equal(t, expected.EligibleForMeals, actual.EligibleForMeals, "expected EligibleForMeals for recipe %s to be %v, but it was %v", expected.ID, expected.EligibleForMeals, actual.EligibleForMeals)

	assert.NotZero(t, actual.CreatedAt)
}

func TestRecipes_Realistic(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		soak := createValidPreparationForTest(t)
		mix := createValidPreparationForTest(t)
		grams := createValidMeasurementUnitForTest(t)
		cups := createValidMeasurementUnitForTest(t)
		pintoBeans := createValidIngredientForTest(t)
		water := createValidIngredientForTest(t)
		garlicPaste := createValidIngredientForTest(t)
		createdValidInstrument := createValidInstrumentForTest(t)

		expected := &mealplanning.Recipe{
			Name:                "sopa de frijol",
			Slug:                "sopa-de-frijol",
			YieldsComponentType: mealplanning.MealComponentTypesMain,
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Media:               []*mealplanning.RecipeMedia{},
			PrepTasks:           []*mealplanning.RecipePrepTask{},
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: nil,
				Min: 1,
			},
			Steps: []*mealplanning.RecipeStep{
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "soaked pinto beans",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1000)),
							},
						},
					},
					Notes:       "first step",
					Preparation: *soak,
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Ingredient:      pintoBeans,
							Name:            "pinto beans",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 500,
							},
						},
						{
							Ingredient:      water,
							Name:            "water",
							MeasurementUnit: *cups,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 5,
							},
						},
					},
					Index: 0,
				},
				{
					Products: []*mealplanning.RecipeStepProduct{
						{
							Name:            "final output",
							Type:            mealplanning.RecipeStepProductIngredientType,
							MeasurementUnit: grams,
							QuantityNotes:   "",
							Quantity: types.OptionalFloat32Range{
								Max: nil,
								Min: pointer.To(float32(1010)),
							},
						},
					},
					Notes:       "second step",
					Preparation: *mix,
					Instruments: []*mealplanning.RecipeStepInstrument{
						{
							Name:       "whatever",
							Instrument: createdValidInstrument,
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredient{
						{
							Name:            "soaked pinto beans",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 1000,
							},
						},
						{
							Ingredient:      garlicPaste,
							Name:            "garlic paste",
							MeasurementUnit: *grams,
							Quantity: types.Float32RangeWithOptionalMax{
								Min: 10,
							},
						},
					},
					Index: 1,
				},
			},
		}

		expectedInput := &mealplanning.RecipeCreationRequestInput{
			Name:                expected.Name,
			Description:         expected.Description,
			Slug:                expected.Slug,
			YieldsComponentType: expected.YieldsComponentType,
			PortionName:         expected.PortionName,
			PluralPortionName:   expected.PluralPortionName,
			EstimatedPortions: types.Float32RangeWithOptionalMax{
				Max: expected.EstimatedPortions.Max,
				Min: expected.EstimatedPortions.Min,
			},
			Steps: []*mealplanning.RecipeStepCreationRequestInput{
				{
					TemperatureInCelsius: expected.Steps[0].TemperatureInCelsius,
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              expected.Steps[0].Products[0].Name,
							Type:              expected.Steps[0].Products[0].Type,
							MeasurementUnitID: &expected.Steps[0].Products[0].MeasurementUnit.ID,
							QuantityNotes:     expected.Steps[0].Products[0].QuantityNotes,
							Quantity:          expected.Steps[0].Products[0].Quantity,
						},
					},
					Notes:         expected.Steps[0].Notes,
					PreparationID: expected.Steps[0].Preparation.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							IngredientID:      &expected.Steps[0].Ingredients[0].Ingredient.ID,
							Name:              expected.Steps[0].Ingredients[0].Name,
							MeasurementUnitID: expected.Steps[0].Ingredients[0].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[0].Ingredients[0].Quantity.Min,
							},
						},
						{
							IngredientID:      &expected.Steps[0].Ingredients[1].Ingredient.ID,
							Name:              expected.Steps[0].Ingredients[1].Name,
							MeasurementUnitID: expected.Steps[0].Ingredients[1].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[0].Ingredients[1].Quantity.Min,
							},
						},
					},
					Index: expected.Steps[0].Index,
				},
				{
					TemperatureInCelsius: expected.Steps[1].TemperatureInCelsius,
					Products: []*mealplanning.RecipeStepProductCreationRequestInput{
						{
							Name:              expected.Steps[1].Products[0].Name,
							Type:              expected.Steps[1].Products[0].Type,
							MeasurementUnitID: &expected.Steps[1].Products[0].MeasurementUnit.ID,
							QuantityNotes:     expected.Steps[1].Products[0].QuantityNotes,
							Quantity:          expected.Steps[1].Products[0].Quantity,
						},
					},
					Notes:         expected.Steps[1].Notes,
					PreparationID: expected.Steps[1].Preparation.ID,
					Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
						{
							Name:         "whatever",
							InstrumentID: pointer.To(createdValidInstrument.ID),
						},
					},
					Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
						{
							Name:                            expected.Steps[1].Ingredients[0].Name,
							MeasurementUnitID:               expected.Steps[1].Ingredients[0].MeasurementUnit.ID,
							ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
							ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[1].Ingredients[0].Quantity.Min,
							},
						},
						{
							IngredientID:      &expected.Steps[1].Ingredients[1].Ingredient.ID,
							Name:              expected.Steps[1].Ingredients[1].Name,
							MeasurementUnitID: expected.Steps[1].Ingredients[1].MeasurementUnit.ID,
							Quantity: types.Float32RangeWithOptionalMax{
								Max: nil,
								Min: expected.Steps[1].Ingredients[1].Quantity.Min,
							},
						},
					},
					Index: expected.Steps[1].Index,
				},
			},
		}

		createdRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(expectedInput)})
		require.NoError(t, err)
		created := converters.ConvertGRPCRecipeToRecipe(createdRes.Created)
		//checkRecipeEquality(t, expected, created)

		e, c := dumpToJSON(expected), dumpToJSON(created)
		_, _ = e, c

		assertRoughEquality(t, expected, created, defaultIgnoredFields("CreatedByUser", "ID", "BelongsToRecipe", "BelongsToRecipeStep")...)

		//_, testClient := createUserAndClientForTest(t)
		//
		//recipeRes, err := testClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{RecipeID: created.ID})
		//require.NoError(t, err)
		//checkRecipeEquality(t, expected, converters.ConvertGRPCRecipeToRecipe(recipeRes.Result))
		//
		//recipeStepProductIndex := -1
		//for i, ingredient := range created.Steps[1].Ingredients {
		//	if ingredient.RecipeStepProductID != nil {
		//		recipeStepProductIndex = i
		//	}
		//}
		//
		//require.NotEqual(t, -1, recipeStepProductIndex)
		//require.NotNil(t, created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
		//assert.Equal(t, created.Steps[0].Products[0].ID, *created.Steps[1].Ingredients[recipeStepProductIndex].RecipeStepProductID)
	})
}

/*

func (s *TestSuite) TestRecipes_Updating() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			newRecipe := fakes.BuildFakeRecipe()
			updateInput := converters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe)
			createdRecipe.Update(updateInput)
			assert.NoError(t, testClients.adminClient.UpdateRecipe(ctx, createdRecipe.ID, updateInput))

			actual, err := testClients.userClient.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// TODO: uncomment me
//func (s *TestSuite) TestRecipes_UploadRecipeMedia() {
//	s.runTest("should be able to upload content for a recipe", func(testClients *testClientWrapper) func() {
//		return func() {
//			t := s.T()
//
//			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
//			defer span.End()
//
//			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)
//
//			newRecipe := fakes.BuildFakeRecipe()
//			updateInput := converters.ConvertRecipeToRecipeUpdateRequestInput(newRecipe)
//			createdRecipe.Update(updateInput)
//			assert.NoError(t, testClients.adminClient.UpdateRecipe(ctx, createdRecipe.ID, updateInput))
//
//			actual, err := testClients.userClient.GetRecipe(ctx, createdRecipe.ID)
//			requireNotNilAndNoProblems(t, actual, err)
//
//			// assert recipe equality
//			checkRecipeEquality(t, newRecipe, actual)
//			assert.NotNil(t, actual.LastUpdatedAt)
//
//			_, img1Bytes := testutils.BuildArbitraryImagePNGBytes(200)
//			_, img2Bytes := testutils.BuildArbitraryImagePNGBytes(250)
//			_, img3Bytes := testutils.BuildArbitraryImagePNGBytes(300)
//
//			files := map[string][]byte{
//				"image_1.png": img1Bytes,
//				"image_2.png": img2Bytes,
//				"image_3.png": img3Bytes,
//			}
//
//			require.NoError(t, testClients.userClient.UploadRecipeMedia(ctx, files, actual.ID))
//
//			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
//		}
//	})
//}

func (s *TestSuite) TestRecipes_AlsoCreateMeal() {
	s.runTest("should be able to create a meal and a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil, func(input *mealplanning.RecipeCreationRequestInput) {
				input.AlsoCreateMeal = true
			})

			mealResults, err := testClients.userClient.SearchForMeals(ctx, createdRecipe.Name, nil)
			requireNotNilAndNoProblems(t, mealResults, err)

			foundMealID := ""
			for _, m := range mealResults.Data {
				meal, mealFetchErr := testClients.userClient.GetMeal(ctx, m.ID)
				requireNotNilAndNoProblems(t, meal, mealFetchErr)

				for _, component := range meal.Components {
					if component.Recipe.ID == createdRecipe.ID {
						foundMealID = meal.ID
					}
				}
			}

			require.NotEmpty(t, foundMealID)

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*mealplanning.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.userClient.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_Searching() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.userClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.adminClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.userClient.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleRecipe := fakes.BuildFakeRecipe()

			var expected []*mealplanning.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe.Name = fmt.Sprintf("example%d", i)
				_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, exampleRecipe)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.userClient.SearchForRecipes(ctx, "example", nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_GetMealPlanTasksForRecipe() {
	s.runTest("meal plan tasks with frozen chicken breast", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			diceBase := fakes.BuildFakeValidPreparation()
			diceInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(diceBase)
			dice, err := testClients.adminClient.CreateValidPreparation(ctx, diceInput)
			require.NoError(t, err)

			exampleGrams := fakes.BuildFakeValidMeasurementUnit()
			exampleGramsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleGrams)
			grams, err := testClients.adminClient.CreateValidMeasurementUnit(ctx, exampleGramsInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			grams, err = testClients.adminClient.GetValidMeasurementUnit(ctx, grams.ID)
			requireNotNilAndNoProblems(t, grams, err)
			checkValidMeasurementUnitEquality(t, exampleGrams, grams)

			chickenBreastBase := fakes.BuildFakeValidIngredient()
			chickenBreastInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(chickenBreastBase)
			chickenBreastInput.StorageTemperatureInCelsius.Min = pointer.To(float32(2.5))
			chickenBreast, createdValidIngredientErr := testClients.adminClient.CreateValidIngredient(ctx, chickenBreastInput)
			require.NoError(t, createdValidIngredientErr)

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.adminClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			sauteeBase := fakes.BuildFakeValidPreparation()
			sauteeInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(sauteeBase)
			sautee, err := testClients.adminClient.CreateValidPreparation(ctx, sauteeInput)
			require.NoError(t, err)

			expected := &mealplanning.Recipe{
				Name:                "sopa de frijol",
				Slug:                "whatever-who-cares-sopa-de-frijol",
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
								Name:            "diced chicken breast",
								Type:            mealplanning.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								Quantity: mealplanning.OptionalFloat32Range{
									Max: nil,
									Min: pointer.To(float32(1000)),
								},
							},
						},
						Notes:       "first step",
						Preparation: *dice,
						Instruments: []*mealplanning.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*mealplanning.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          chickenBreast,
								Name:                "pinto beans",
								MeasurementUnit:     *grams,
								Quantity: mealplanning.Float32RangeWithOptionalMax{
									Min: 500,
								},
							},
						},
						Index: 0,
					},
					{
						Products: []*mealplanning.RecipeStepProduct{
							{
								Name:            "final output",
								Type:            mealplanning.RecipeStepProductIngredientType,
								MeasurementUnit: grams,
								QuantityNotes:   "",
								Quantity: mealplanning.OptionalFloat32Range{
									Max: nil,
									Min: pointer.To(float32(1010)),
								},
							},
						},
						Notes:       "second step",
						Preparation: *sautee,
						Instruments: []*mealplanning.RecipeStepInstrument{
							{
								Name:       "whatever",
								Instrument: createdValidInstrument,
							},
						},
						Ingredients: []*mealplanning.RecipeStepIngredient{
							{
								Name:            "diced chicken breast",
								MeasurementUnit: *grams,
								Quantity: mealplanning.Float32RangeWithOptionalMax{
									Min: 1000,
								},
							},
						},
						Index: 1,
					},
				},
			}

			expectedInput := &mealplanning.RecipeCreationRequestInput{
				Name:                expected.Name,
				Slug:                expected.Slug,
				YieldsComponentType: expected.YieldsComponentType,
				PortionName:         expected.PortionName,
				PluralPortionName:   expected.PluralPortionName,
				EstimatedPortions: mealplanning.Float32RangeWithOptionalMax{
					Min: expected.EstimatedPortions.Min,
					Max: expected.EstimatedPortions.Max,
				},
				Steps: []*mealplanning.RecipeStepCreationRequestInput{
					{
						Products: []*mealplanning.RecipeStepProductCreationRequestInput{
							{
								Name:              "diced chicken breast",
								Type:              mealplanning.RecipeStepProductIngredientType,
								MeasurementUnitID: &grams.ID,
								QuantityNotes:     "",
								Quantity:          mealplanning.OptionalFloat32Range{Min: pointer.To(float32(1000))},
							},
						},
						Notes:         "first step",
						PreparationID: dice.ID,
						Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
							{
								IngredientID:      &chickenBreast.ID,
								Name:              "pinto beans",
								MeasurementUnitID: grams.ID,
								Quantity:          mealplanning.Float32RangeWithOptionalMax{Min: 500},
							},
						},
						Index: 0,
					},
					{
						Products: []*mealplanning.RecipeStepProductCreationRequestInput{
							{
								Name:              "final output",
								Type:              mealplanning.RecipeStepProductIngredientType,
								MeasurementUnitID: &grams.ID,
								QuantityNotes:     "",
								Quantity:          mealplanning.OptionalFloat32Range{Min: pointer.To(float32(1010))},
							},
						},
						Notes:         "second step",
						PreparationID: sautee.ID,
						Instruments: []*mealplanning.RecipeStepInstrumentCreationRequestInput{
							{
								Name:         "whatever",
								InstrumentID: pointer.To(createdValidInstrument.ID),
							},
						},
						Ingredients: []*mealplanning.RecipeStepIngredientCreationRequestInput{
							{
								Name:                            "diced chicken breast",
								MeasurementUnitID:               grams.ID,
								Quantity:                        mealplanning.Float32RangeWithOptionalMax{Min: 1000},
								ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
								ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
							},
						},
						Index: 1,
					},
				},
			}

			created, err := testClients.adminClient.CreateRecipe(ctx, expectedInput)
			require.NoError(t, err)
			checkRecipeEquality(t, expected, created)

			steps, err := testClients.userClient.GetMealPlanTasks(ctx, created.ID, nil)
			requireNotNilAndNoProblems(t, created, err)

			require.NotEmpty(t, steps)
		}
	})
}

func (s *TestSuite) TestRecipes_Cloning() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			actual, err := testClients.userClient.CloneRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			require.Equal(t, createdRecipe.Name, actual.Name)
			require.Equal(t, len(createdRecipe.Steps), len(actual.Steps))

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// TODO: uncomment me
//func (s *TestSuite) TestRecipes_DAGGeneration() {
//	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
//		return func() {
//			t := s.T()
//
//			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
//			defer span.End()
//
//			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)
//
//			actual, err := testClients.userClient.GetRecipeDAG(ctx, createdRecipe.ID)
//			requireNotNilAndNoProblems(t, actual, err)
//
//			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
//		}
//	})
//}

*/
