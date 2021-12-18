package integration

import (
	"context"
	"testing"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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

func createRecipeWithNotificationChannel(ctx context.Context, t *testing.T, notificationsChan chan *types.DataChangeMessage, client *httpclient.Client) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	t.Helper()

	var n *types.DataChangeMessage

	t.Log("creating prerequisite valid preparation")
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
	createdValidPreparationID, err := client.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, err)
	t.Logf("valid preparation %q created", createdValidPreparationID)

	n = <-notificationsChan
	assert.Equal(t, types.ValidPreparationDataType, n.DataType)
	require.NotNil(t, n.ValidPreparation)
	checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

	createdValidPreparation, err := client.GetValidPreparation(ctx, createdValidPreparationID)
	requireNotNilAndNoProblems(t, createdValidPreparation, err)
	checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

	t.Log("creating recipe")
	exampleRecipe := fakes.BuildFakeRecipe()

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := client.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientDataType, n.DataType)
			require.NotNil(t, n.ValidIngredient)
			checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

			createdValidIngredient, err := client.GetValidIngredient(ctx, createdValidIngredientID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

			exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
		}
	}

	exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	createdRecipeID, err := client.CreateRecipe(ctx, exampleRecipeInput)
	require.NoError(t, err)
	t.Logf("recipe %q created", createdRecipeID)

	n = <-notificationsChan
	assert.Equal(t, types.RecipeDataType, n.DataType)
	require.NotNil(t, n.Recipe)
	checkRecipeEquality(t, exampleRecipe, n.Recipe)

	createdRecipe, err := client.GetRecipe(ctx, createdRecipeID)
	requireNotNilAndNoProblems(t, createdRecipe, err)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

func createRecipeWhilePolling(ctx context.Context, t *testing.T, client *httpclient.Client) ([]*types.ValidIngredient, *types.ValidPreparation, *types.Recipe) {
	t.Helper()

	var checkFunc func() bool

	t.Log("creating valid preparation")
	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
	createdValidPreparationID, validPreparationCreationErr := client.CreateValidPreparation(ctx, exampleValidPreparationInput)
	require.NoError(t, validPreparationCreationErr)

	var (
		createdValidPreparation *types.ValidPreparation
	)
	checkFunc = func() bool {
		createdValidPreparation, validPreparationCreationErr = client.GetValidPreparation(ctx, createdValidPreparationID)
		return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, validPreparationCreationErr)
	}
	assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
	checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)
	t.Logf("valid preparation %q created", createdValidPreparationID)

	t.Log("creating recipe")
	exampleRecipe := fakes.BuildFakeRecipe()

	createdValidIngredients := []*types.ValidIngredient{}
	for i, recipeStep := range exampleRecipe.Steps {
		for j := range recipeStep.Ingredients {
			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, validIngredientCreationErr := client.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, validIngredientCreationErr)

			var (
				createdValidIngredient *types.ValidIngredient
			)
			checkFunc = func() bool {
				createdValidIngredient, validIngredientCreationErr = client.GetValidIngredient(ctx, createdValidIngredientID)
				return assert.NotNil(t, createdValidIngredient) && assert.NoError(t, validIngredientCreationErr)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			createdValidIngredients = append(createdValidIngredients, createdValidIngredient)
			exampleRecipe.Steps[i].Ingredients[j].IngredientID = stringPointer(createdValidIngredient.ID)
		}
	}

	exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID
	}

	createdRecipeID, err := client.CreateRecipe(ctx, exampleRecipeInput)
	require.NoError(t, err)

	var createdRecipe *types.Recipe
	checkFunc = func() bool {
		createdRecipe, err = client.GetRecipe(ctx, createdRecipeID)
		return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
	}
	assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
	checkRecipeEquality(t, exampleRecipe, createdRecipe)
	t.Logf("recipe %q created", createdRecipeID)

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

func (s *TestSuite) TestRecipes_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			_, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			t.Log("changing recipe")
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(newRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			n = <-notificationsChan
			assert.Equal(t, types.RecipeDataType, n.DataType)

			t.Log("fetching changed recipe")
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

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

			_, _, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

			// change recipe
			newRecipe := fakes.BuildFakeRecipe()
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(newRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			time.Sleep(2 * time.Second)

			// retrieve changed recipe
			var actual *types.Recipe
			checkFunc = func() bool {
				var err error
				actual, err = testClients.main.GetRecipe(ctx, createdRecipe.ID)
				return assert.NotNil(t, actual) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			// assert recipe equality
			checkRecipeEquality(t, newRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
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
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
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
			assert.Equal(t, types.ValidIngredientDataType, n.DataType)
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
			assert.Equal(t, types.ValidPreparationDataType, n.DataType)
			require.NotNil(t, n.ValidPreparation)
			checkValidPreparationEquality(t, exampleValidPreparation, n.ValidPreparation)

			createdValidPreparation, err := testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

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

			var createdValidPreparation *types.ValidPreparation
			checkFunc = func() bool {
				createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparationID)
				return assert.NotNil(t, createdValidPreparation) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)
			t.Logf("valid preparation %q created", createdValidPreparationID)

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredientID, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)

			var createdValidIngredient *types.ValidIngredient
			checkFunc = func() bool {
				createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredientID)
				return assert.NotNil(t, createdValidIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredientID)

			t.Log("creating recipes")
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				_, _, createdRecipe := createRecipeWhilePolling(ctx, t, testClients.main)

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
