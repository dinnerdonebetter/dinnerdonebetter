package integration

import (
	"bytes"
	"encoding/json"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"
)

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *types.RecipeStepInstrument, checkInstrument bool) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	if checkInstrument {
		checkValidInstrumentEquality(t, expected.Instrument, actual.Instrument)
	} else {
		assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected Instrument.ID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Instrument.ID, actual.Instrument.ID)
	}
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep, "expected ProductOfRecipeStep for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep)
	assert.Equal(t, expected.RecipeStepProductID, actual.RecipeStepProductID, "expected RecipeStepProductID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.RecipeStepProductID, actual.RecipeStepProductID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.PreferenceRank, actual.PreferenceRank, "expected PreferenceRank for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.PreferenceRank, actual.PreferenceRank)
	assert.Equal(t, expected.Optional, actual.Optional, "expected Optional for recipe step instrument %s to be %v, but was %v", expected.ID, expected.Optional, actual.Optional)
	assert.Equal(t, expected.MinimumQuantity, actual.MinimumQuantity, "expected MinimumQuantity for recipe step instrument %s to be %v, but was %v", expected.ID, expected.MinimumQuantity, actual.MinimumQuantity)
	assert.Equal(t, expected.MaximumQuantity, actual.MaximumQuantity, "expected MaximumQuantity for recipe step instrument %s to be %v, but was %v", expected.ID, expected.MaximumQuantity, actual.MaximumQuantity)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrument.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
			exampleRecipeStepInstrumentInput := converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := testClients.user.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

			createdRecipeStepInstrument, err = testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
			exampleRecipeStepInstrument.Instrument = createdValidInstrument
			exampleRecipeStepInstrument.Instrument.CreatedAt = createdRecipeStepInstrument.Instrument.CreatedAt

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

			t.Log("creating valid instrument")
			newExampleValidInstrument := fakes.BuildFakeValidInstrument()
			newExampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(newExampleValidInstrument)
			newValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, newExampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, newExampleValidInstrument, newValidInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			newRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepInstrument.Instrument = newValidInstrument
			createdRecipeStepInstrument.Update(converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentUpdateRequestInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.user.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual, false)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.user.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepInstruments_AsRecipeStepProducts() {
	s.runForEachClient("should be able to use a recipe step instrument that was the product of a prior recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid preparation")
			lineBase := fakes.BuildFakeValidPreparation()
			lineInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(lineBase)
			line, err := testClients.admin.CreateValidPreparation(ctx, lineInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", line.ID)

			t.Log("creating prerequisite valid preparation")
			roastBase := fakes.BuildFakeValidPreparation()
			roastInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(roastBase)
			roast, err := testClients.admin.CreateValidPreparation(ctx, roastInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", roast.ID)

			t.Log("creating valid instrument")
			bakingSheetBase := fakes.BuildFakeValidInstrument()
			bakingSheetBaseInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(bakingSheetBase)
			bakingSheet, err := testClients.admin.CreateValidInstrument(ctx, bakingSheetBaseInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", bakingSheet.ID)
			checkValidInstrumentEquality(t, bakingSheetBase, bakingSheet)

			t.Log("creating valid measurement units")
			sheetsBase := fakes.BuildFakeValidMeasurementUnit()
			sheetsBaseInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(sheetsBase)
			sheets, err := testClients.admin.CreateValidMeasurementUnit(ctx, sheetsBaseInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", sheets.ID)
			checkValidMeasurementUnitEquality(t, sheetsBase, sheets)

			t.Log("creating valid measurement units")
			headsBase := fakes.BuildFakeValidMeasurementUnit()
			headsBaseInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(headsBase)
			head, err := testClients.admin.CreateValidMeasurementUnit(ctx, headsBaseInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", head.ID)
			checkValidMeasurementUnitEquality(t, headsBase, head)

			t.Log("creating valid measurement units")
			exampleUnits := fakes.BuildFakeValidMeasurementUnit()
			exampleUnitsInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleUnits)
			unit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleUnitsInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", unit.ID)
			checkValidMeasurementUnitEquality(t, exampleUnits, unit)

			t.Log("creating prerequisite valid ingredient")
			aluminumFoilBase := fakes.BuildFakeValidIngredient()
			aluminumFoilInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(aluminumFoilBase)
			aluminumFoil, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, aluminumFoilInput)
			require.NoError(t, createdValidIngredientErr)

			t.Log("creating prerequisite valid ingredient")
			garlic := fakes.BuildFakeValidIngredient()
			garlicInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(garlic)
			garlic, garlicErr := testClients.admin.CreateValidIngredient(ctx, garlicInput)
			require.NoError(t, garlicErr)

			linedBakingSheetName := "lined baking sheet"

			t.Log("creating recipe")
			expected := &types.Recipe{
				Name:        t.Name(),
				Description: "",
				Steps: []*types.RecipeStep{
					{
						Products: []*types.RecipeStepProduct{
							{
								Name:            linedBakingSheetName,
								Type:            types.RecipeStepProductInstrumentType,
								MeasurementUnit: *unit,
								QuantityNotes:   "",
								MinimumQuantity: 1,
							},
						},
						Notes:       "first step",
						Preparation: *line,
						Ingredients: []*types.RecipeStepIngredient{
							{
								RecipeStepProductID: nil,
								Ingredient:          aluminumFoil,
								Name:                "aluminum foil",
								MeasurementUnit:     *sheets,
								MinimumQuantity:     3,
								ProductOfRecipeStep: false,
							},
						},
						Instruments: []*types.RecipeStepInstrument{
							{
								Instrument: bakingSheet,
							},
						},
						Index: 0,
					},
					{
						Instruments: []*types.RecipeStepInstrument{
							{
								Name:                linedBakingSheetName,
								Instrument:          nil,
								ProductOfRecipeStep: true,
							},
						},
						Products: []*types.RecipeStepProduct{
							{
								Name:            "roasted garlic",
								Type:            types.RecipeStepProductIngredientType,
								MeasurementUnit: *head,
								QuantityNotes:   "",
								MinimumQuantity: 1,
							},
						},
						Notes:       "second step",
						Preparation: *roast,
						Ingredients: []*types.RecipeStepIngredient{
							{
								Ingredient:          garlic,
								Name:                "garlic",
								MeasurementUnit:     *head,
								MinimumQuantity:     1,
								ProductOfRecipeStep: false,
							},
						},
						Index: 1,
					},
				},
			}

			exampleRecipeInput := &types.RecipeCreationRequestInput{
				Name:        expected.Name,
				Description: expected.Description,
			}
			for _, step := range expected.Steps {
				newStep := &types.RecipeStepCreationRequestInput{
					MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
					Notes:                         step.Notes,
					PreparationID:                 step.Preparation.ID,
					BelongsToRecipe:               step.BelongsToRecipe,
					ID:                            step.ID,
					Index:                         step.Index,
					MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
					MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
					Optional:                      step.Optional,
				}

				for _, ingredient := range step.Ingredients {
					newIngredient := &types.RecipeStepIngredientCreationRequestInput{
						IngredientID:        &ingredient.Ingredient.ID,
						ID:                  ingredient.ID,
						BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
						Name:                ingredient.Name,
						MeasurementUnitID:   ingredient.MeasurementUnit.ID,
						QuantityNotes:       ingredient.QuantityNotes,
						IngredientNotes:     ingredient.IngredientNotes,
						MinimumQuantity:     ingredient.MinimumQuantity,
						ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
					}
					newStep.Ingredients = append(newStep.Ingredients, newIngredient)
				}

				for _, product := range step.Products {
					newProduct := &types.RecipeStepProductCreationRequestInput{
						ID:                  product.ID,
						Name:                product.Name,
						Type:                product.Type,
						MeasurementUnitID:   product.MeasurementUnit.ID,
						QuantityNotes:       product.QuantityNotes,
						BelongsToRecipeStep: product.BelongsToRecipeStep,
						MinimumQuantity:     product.MinimumQuantity,
					}
					newStep.Products = append(newStep.Products, newProduct)
				}

				for _, instrument := range step.Instruments {
					newInstrument := &types.RecipeStepInstrumentCreationRequestInput{
						Name:                instrument.Name,
						Notes:               instrument.Notes,
						RecipeStepProductID: instrument.RecipeStepProductID,
						ProductOfRecipeStep: instrument.ProductOfRecipeStep,
						PreferenceRank:      instrument.PreferenceRank,
					}

					if instrument.Instrument != nil {
						newInstrument.InstrumentID = &instrument.Instrument.ID
					}

					newStep.Instruments = append(newStep.Instruments, newInstrument)
				}

				exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
			}

			var b bytes.Buffer
			require.NoError(t, json.NewEncoder(&b).Encode(exampleRecipeInput))
			t.Logf("creating recipe with input: %s", b.String())

			created, err := testClients.user.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", created.ID)
			checkRecipeEquality(t, expected, created)

			created, err = testClients.user.GetRecipe(ctx, created.ID)
			requireNotNilAndNoProblems(t, created, err)
			checkRecipeEquality(t, expected, created)

			createdJSON, _ := json.Marshal(created)
			t.Log(string(createdJSON))

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
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrument.Instrument = &types.ValidInstrument{ID: createdValidInstrument.ID}
				exampleRecipeStepInstrumentInput := converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(exampleRecipeStepInstrument)
				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.user.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrument.ID)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument, false)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr = testClients.user.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.user.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepInstruments),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepInstrument := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
