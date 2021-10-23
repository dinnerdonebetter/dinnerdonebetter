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

func checkRecipeEquality(t *testing.T, expected, actual *types.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe %s to be %v, but it was %v", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe %s to be %v, but it was %v", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeToRecipeUpdateInput creates an RecipeUpdateRequestInput struct from a recipe.
func convertRecipeToRecipeUpdateInput(x *types.Recipe) *types.RecipeUpdateRequestInput {
	return &types.RecipeUpdateRequestInput{
		Name:               x.Name,
		Source:             x.Source,
		Description:        x.Description,
		InspiredByRecipeID: x.InspiredByRecipeID,
	}
}

func (s *TestSuite) TestRecipes_CompleteLifecycle() {
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

			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			t.Log("changing recipe")
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(newRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)

			t.Log("fetching changed recipe")
			actual, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

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

			// assert recipe equality
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			// change recipe
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(newRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			time.Sleep(time.Second)

			// retrieve changed recipe
			var actual *types.Recipe
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}

func (s *TestSuite) TestRecipes_Listing() {
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

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe := fakes.BuildFakeRecipe()

				for i, recipeStep := range exampleRecipe.Steps {
					exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
					for j := range recipeStep.Ingredients {
						exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
					}
				}

				exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
				createdRecipeID, createdRecipeErr := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
				require.NoError(t, createdRecipeErr)
				t.Logf("recipe %q created", createdRecipeID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeDataType)
				require.NotNil(t, n.Recipe)
				checkRecipeEquality(t, exampleRecipe, n.Recipe)

				createdRecipe, createdRecipeErr := testClients.main.GetRecipe(ctx, createdRecipeID)
				requireNotNilAndNoProblems(t, createdRecipe, createdRecipeErr)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.main.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			t.Log("cleaning up")
			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
			}
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

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe := fakes.BuildFakeRecipe()

				for j, recipeStep := range exampleRecipe.Steps {
					exampleRecipe.Steps[j].PreparationID = createdValidPreparation.ID
					for k := range recipeStep.Ingredients {
						exampleRecipe.Steps[j].Ingredients[k].IngredientID = stringPointer(createdValidIngredient.ID)
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

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.main.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			t.Log("cleaning up")
			for _, createdRecipe := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}
