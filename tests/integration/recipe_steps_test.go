package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkRecipeStepEquality(t *testing.T, expected, actual *types.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step %s to be %v, but it was %v", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for recipe step %s to be %v, but it was %v", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.PrerequisiteStep, actual.PrerequisiteStep, "expected PrerequisiteStep for recipe step %s to be %v, but it was %v", expected.ID, expected.PrerequisiteStep, actual.PrerequisiteStep)
	assert.Equal(t, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds, "expected MinEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds)
	assert.Equal(t, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds, "expected MaxEstimatedTimeInSeconds for recipe step %s to be %v, but it was %v", expected.ID, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds)
	assert.Equal(t, expected.TemperatureInCelsius, actual.TemperatureInCelsius, "expected TemperatureInCelsius for recipe step %s to be %v, but it was %v", expected.ID, expected.TemperatureInCelsius, actual.TemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.RecipeID, actual.RecipeID, "expected RecipeID for recipe step %s to be %v, but it was %v", expected.ID, expected.RecipeID, actual.RecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepToRecipeStepUpdateInput creates an RecipeStepUpdateRequestInput struct from a recipe step.
func convertRecipeStepToRecipeStepUpdateInput(x *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	return &types.RecipeStepUpdateRequestInput{
		Index:                     x.Index,
		PreparationID:             x.PreparationID,
		PrerequisiteStep:          x.PrerequisiteStep,
		MinEstimatedTimeInSeconds: x.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: x.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      x.TemperatureInCelsius,
		Notes:                     x.Notes,
		RecipeID:                  x.RecipeID,
	}
}

func (s *TestSuite) TestRecipeSteps_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}

			t.Log("changing recipe step")
			newRecipeStep := fakes.BuildFakeRecipeStep()
			newRecipeStep.BelongsToRecipe = createdRecipe.ID
			newRecipeStep.PreparationID = createdValidPreparation.ID
			for j := range newRecipeStep.Ingredients {
				newRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
			}

			createdRecipeStep.Update(convertRecipeStepToRecipeStepUpdateInput(newRecipeStep))
			assert.NoError(t, testClients.main.UpdateRecipeStep(ctx, createdRecipeStep))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepDataType)

			t.Log("fetching changed recipe step")
			actual, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, newRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}

			// change recipe step
			newRecipeStep := fakes.BuildFakeRecipeStep()
			newRecipeStep.BelongsToRecipe = createdRecipe.ID
			newRecipeStep.PreparationID = createdValidPreparation.ID
			for j := range newRecipeStep.Ingredients {
				newRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
			}

			createdRecipeStep.Update(convertRecipeStepToRecipeStepUpdateInput(newRecipeStep))
			assert.NoError(t, testClients.main.UpdateRecipeStep(ctx, createdRecipeStep))

			time.Sleep(2 * time.Second)

			// retrieve changed recipe step
			var (
				actual *types.RecipeStep
				err    error
			)
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
				return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, newRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStep.PreparationID = createdValidPreparation.ID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
				}

				exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
				createdRecipeStepID, createdRecipeStepErr := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
				require.NoError(t, createdRecipeStepErr)
				t.Logf("recipe step %q created", createdRecipeStepID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepDataType)
				require.NotNil(t, n.RecipeStep)
				checkRecipeStepEquality(t, exampleRecipeStep, n.RecipeStep)

				createdRecipeStep, createdRecipeStepErr := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
				requireNotNilAndNoProblems(t, createdRecipeStep, createdRecipeStepErr)
				require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

				expected = append(expected, createdRecipeStep)
			}

			// assert recipe step list equality
			actual, err := testClients.main.GetRecipeSteps(ctx, createdRecipe.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredients, createdValidPreparation, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStep.PreparationID = createdValidPreparation.ID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredients[j].ID)
				}

				exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
				createdRecipeStepID, createdRecipeStepErr := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
				require.NoError(t, createdRecipeStepErr)

				var createdRecipeStep *types.RecipeStep
				checkFunc = func() bool {
					createdRecipeStep, createdRecipeStepErr = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
					return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, createdRecipeStepErr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

				expected = append(expected, createdRecipeStep)
			}

			// assert recipe step list equality
			actual, err := testClients.main.GetRecipeSteps(ctx, createdRecipe.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
