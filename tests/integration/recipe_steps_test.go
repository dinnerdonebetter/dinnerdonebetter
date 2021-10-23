package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
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

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidIngredientDataType)
			require.NotNil(t, n.ValidIngredient)
			checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

			createdValidIngredient, err := testClients.main.GetValidIngredient(ctx, createdValidIngredientID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidPreparationDataType)
			require.NotNil(t, n.ValidPreparation)
			checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

			createdValidPreparation, err := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()

			for i, recipeStep := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range recipeStep.Ingredients {
					exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)
			t.Logf("recipe %q created", createdRecipeID)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)
			require.NotEmpty(t, createdRecipe.Steps)

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
				newRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
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
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			var createdValidPreparation *types.ValidPreparation
			checkFunc = func() bool {
				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			var createdValidIngredient *types.ValidIngredient
			checkFunc = func() bool {
				createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredientID)
				return assert.NotNil(t, createdValidIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()

			for i, recipeStep := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range recipeStep.Ingredients {
					exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

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
				newRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
			}

			createdRecipeStep.Update(convertRecipeStepToRecipeStepUpdateInput(newRecipeStep))
			assert.NoError(t, testClients.main.UpdateRecipeStep(ctx, createdRecipeStep))

			time.Sleep(time.Second)

			// retrieve changed recipe step
			var actual *types.RecipeStep
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
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
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

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidIngredientDataType)
			require.NotNil(t, n.ValidIngredient)
			checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

			createdValidIngredient, err := testClients.main.GetValidIngredient(ctx, createdValidIngredientID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.ValidPreparationDataType)
			require.NotNil(t, n.ValidPreparation)
			checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

			createdValidPreparation, err := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()

			for i, recipeStep := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range recipeStep.Ingredients {
					exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)
			t.Logf("recipe %q created", createdRecipeID)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)
			require.NotEmpty(t, createdRecipe.Steps)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStep.PreparationID = createdValidPreparationID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
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
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparationID, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			var createdValidPreparation *types.ValidPreparation
			checkFunc = func() bool {
				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			var createdValidIngredient *types.ValidIngredient
			checkFunc = func() bool {
				createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredientID)
				return assert.NotNil(t, createdValidIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()

			for i, recipeStep := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range recipeStep.Ingredients {
					recipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)
			t.Logf("recipe %q created", createdRecipeID)

			t.Log("creating recipe steps")
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStep.PreparationID = createdValidPreparationID
				for j := range exampleRecipeStep.Ingredients {
					exampleRecipeStep.Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
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
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}
