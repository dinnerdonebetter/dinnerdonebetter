package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipeStepEquality(t *testing.T, expected, actual *types.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step %s to be %v, but it was %v", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.MinimumEstimatedTimeInSeconds, actual.MinimumEstimatedTimeInSeconds, "expected MinimumEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MinimumEstimatedTimeInSeconds, actual.MinimumEstimatedTimeInSeconds)
	assert.Equal(t, expected.MaximumEstimatedTimeInSeconds, actual.MaximumEstimatedTimeInSeconds, "expected MaximumEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MaximumEstimatedTimeInSeconds, actual.MaximumEstimatedTimeInSeconds)
	assert.Equal(t, expected.MinimumTemperatureInCelsius, actual.MinimumTemperatureInCelsius, "expected MinimumTemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.MinimumTemperatureInCelsius, actual.MinimumTemperatureInCelsius)
	assert.Equal(t, expected.MaximumTemperatureInCelsius, actual.MaximumTemperatureInCelsius, "expected MaximumTemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.MaximumTemperatureInCelsius, actual.MaximumTemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ExplicitInstructions, actual.ExplicitInstructions, "expected ExplicitInstructions for recipe step %s to be %v, but it was %v", expected.ID, expected.ExplicitInstructions, actual.ExplicitInstructions)
	assert.NotZero(t, actual.CreatedAt)
}

// convertRecipeStepToRecipeStepUpdateInput creates an RecipeStepUpdateRequestInput struct from a recipe step.
func convertRecipeStepToRecipeStepUpdateInput(x *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Index:                         &x.Index,
		MinimumEstimatedTimeInSeconds: &x.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: &x.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   x.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   x.MaximumTemperatureInCelsius,
		Preparation:                   &x.Preparation,
		Optional:                      &x.Optional,
		Notes:                         &x.Notes,
		ExplicitInstructions:          &x.ExplicitInstructions,
		BelongsToRecipe:               x.BelongsToRecipe,
	}
}

func (s *TestSuite) TestRecipeSteps_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}

			t.Log("changing recipe step")
			newRecipeStep := fakes.BuildFakeRecipeStep()
			newRecipeStep.BelongsToRecipe = createdRecipe.ID
			for j := range newRecipeStep.Ingredients {
				newRecipeStep.Ingredients[j].Ingredient = createdValidIngredients[j]
			}

			updateInput := convertRecipeStepToRecipeStepUpdateInput(newRecipeStep)
			updateInput.Preparation = createdValidPreparation
			createdRecipeStep.Update(updateInput)
			assert.NoError(t, testClients.user.UpdateRecipeStep(ctx, createdRecipeStep))

			t.Log("fetching changed recipe step")
			actual, err := testClients.user.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, newRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("creating valid instrument")
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakes.BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			t.Logf("valid instrument %q created", createdValidInstrument.ID)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].Ingredient = createdValidIngredients[j]
					exampleRecipeStep.Ingredients[j].MeasurementUnit = *createdValidMeasurementUnit
				}

				for j := range exampleRecipeStep.Products {
					exampleRecipeStep.Products[j].MeasurementUnit = *createdValidMeasurementUnit
				}

				for j := range exampleRecipeStep.Instruments {
					exampleRecipeStep.Instruments[j].Instrument = createdValidInstrument
				}

				exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
				exampleRecipeStepInput.PreparationID = createdValidPreparation.ID

				createdRecipeStep, createdRecipeStepErr := testClients.user.CreateRecipeStep(ctx, exampleRecipeStepInput)
				require.NoError(t, createdRecipeStepErr)
				t.Logf("recipe step %q created", createdRecipeStep.ID)

				checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

				createdRecipeStep, createdRecipeStepErr = testClients.user.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
				requireNotNilAndNoProblems(t, createdRecipeStep, createdRecipeStepErr)
				require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

				expected = append(expected, createdRecipeStep)
			}

			// assert recipe step list equality
			actual, err := testClients.user.GetRecipeSteps(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeSteps),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeSteps),
			)

			t.Log("cleaning up")
			for _, createdRecipeStep := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
