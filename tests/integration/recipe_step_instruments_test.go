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

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *types.RecipeStepInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput creates an RecipeStepInstrumentUpdateRequestInput struct from a recipe step instrument.
func convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(x *types.RecipeStepInstrument) *types.RecipeStepInstrumentUpdateRequestInput {
	return &types.RecipeStepInstrumentUpdateRequestInput{
		InstrumentID: x.InstrumentID,
		RecipeStepID: x.RecipeStepID,
		Notes:        x.Notes,
	}
}

func (s *TestSuite) TestRecipeStepInstruments_CompleteLifecycle() {
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

			t.Log("creating recipe")
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

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
			require.NotNil(t, n.RecipeStepInstrument)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

			createdRecipeStepInstrument, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, createdRecipeStepInstrument, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			t.Log("changing recipe step instrument")
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)

			t.Log("fetching changed recipe step instrument")
			actual, err := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

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

			t.Log("creating recipe")
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

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instrument")
			exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrumentID, err := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			require.NoError(t, err)
			t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

			var createdRecipeStepInstrument *types.RecipeStepInstrument
			checkFunc = func() bool {
				createdRecipeStepInstrument, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// change recipe step instrument
			newRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
			createdRecipeStepInstrument.Update(convertRecipeStepInstrumentToRecipeStepInstrumentUpdateInput(newRecipeStepInstrument))
			assert.NoError(t, testClients.main.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument))

			time.Sleep(time.Second)

			// retrieve changed recipe step instrument
			var actual *types.RecipeStepInstrument
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
				return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step instrument equality
			checkRecipeStepInstrumentEquality(t, newRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step instrument")
			assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}

func (s *TestSuite) TestRecipeStepInstruments_Listing() {
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

			t.Log("creating recipe")
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

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrumentID, createdRecipeStepInstrumentErr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, createdRecipeStepInstrumentErr)
				t.Logf("recipe step instrument %q created", createdRecipeStepInstrumentID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepInstrumentDataType)
				require.NotNil(t, n.RecipeStepInstrument)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, n.RecipeStepInstrument)

				createdRecipeStepInstrument, createdRecipeStepInstrumentErr := testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
				requireNotNilAndNoProblems(t, createdRecipeStepInstrument, createdRecipeStepInstrumentErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepInstrument.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.main.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

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

			t.Log("creating recipe")
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

			var createdRecipeStepID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step instruments")
			var expected []*types.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepInstrumentInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrumentID, EcreatedRecipeStepInstrumentrr := testClients.main.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				require.NoError(t, EcreatedRecipeStepInstrumentrr)

				var createdRecipeStepInstrument *types.RecipeStepInstrument
				checkFunc = func() bool {
					createdRecipeStepInstrument, EcreatedRecipeStepInstrumentrr = testClients.main.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrumentID)
					return assert.NotNil(t, createdRecipeStepInstrument) && assert.NoError(t, EcreatedRecipeStepInstrumentrr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// assert recipe step instrument list equality
			actual, err := testClients.main.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStepID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepInstrument.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}
