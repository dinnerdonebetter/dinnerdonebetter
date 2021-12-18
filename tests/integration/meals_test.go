package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkMealEquality(t *testing.T, expected, actual *types.Meal) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for meal %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for meal %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.NotZero(t, actual.CreatedOn)
}

/*

func createMealWithNotificationChannel(ctx context.Context, t *testing.T, notificationsChan chan *types.DataChangeMessage, client *httpclient.Client) *types.Meal {
	t.Helper()

	var n *types.DataChangeMessage

	createdRecipes := []*types.Recipe{}
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, client)
		createdRecipes = append(createdRecipes, recipe)
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
	}

	t.Log("creating meal")
	exampleMeal := fakes.BuildFakeMeal()
	exampleMealInput := fakes.BuildFakeMealCreationRequestInputFromMeal(exampleMeal)
	exampleMealInput.Recipes = createdRecipeIDs

	createdMealID, err := client.CreateMeal(ctx, exampleMealInput)
	require.NoError(t, err)

	n = <-notificationsChan
	assert.Equal(t, types.MealDataType, n.DataType)
	require.NotNil(t, n.Meal)
	checkMealEquality(t, exampleMeal, n.Meal)
	t.Logf("meal %q created", createdMealID)

	createdMeal, err := client.GetMeal(ctx, createdMealID)
	requireNotNilAndNoProblems(t, createdMeal, err)
	checkMealEquality(t, exampleMeal, createdMeal)

	return createdMeal
}

func (s *TestSuite) TestMeals_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			createdMeal := createMealWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			t.Log("cleaning up meal")
			assert.NoError(t, testClients.main.ArchiveMeal(ctx, createdMeal.ID))
		}
	})
}

func (s *TestSuite) TestMeals_Listing() {
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

			t.Log("creating meals")
			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				createdMeal := createMealWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

				expected = append(expected, createdMeal)
			}

			// assert meal list equality
			actual, err := testClients.main.GetMeals(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Meals),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Meals),
			)

			t.Log("cleaning up")
			for _, createdMeal := range expected {
				assert.NoError(t, testClients.main.ArchiveMeal(ctx, createdMeal.ID))
			}
		}
	})
}



*/
