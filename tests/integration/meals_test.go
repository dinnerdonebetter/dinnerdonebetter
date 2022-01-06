package integration

import (
	"context"
	"testing"

	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

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

func createMealForTest(ctx context.Context, t *testing.T, client *httpclient.Client) *types.Meal {
	t.Helper()

	createdRecipes := []*types.Recipe{}
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeForTest(ctx, t, client)
		createdRecipes = append(createdRecipes, recipe)
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
	}

	t.Log("creating meal")
	exampleMeal := fakes.BuildFakeMeal()
	exampleMealInput := fakes.BuildFakeMealCreationRequestInputFromMeal(exampleMeal)
	exampleMealInput.Recipes = createdRecipeIDs

	createdMeal, err := client.CreateMeal(ctx, exampleMealInput)
	require.NoError(t, err)

	t.Logf("meal %q created", createdMeal.ID)

	createdMeal, err = client.GetMeal(ctx, createdMeal.ID)
	requireNotNilAndNoProblems(t, createdMeal, err)
	checkMealEquality(t, exampleMeal, createdMeal)

	return createdMeal
}

func (s *TestSuite) TestMeals_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMeal := createMealForTest(ctx, t, testClients.main)

			t.Log("cleaning up meal")
			assert.NoError(t, testClients.main.ArchiveMeal(ctx, createdMeal.ID))
		}
	})
}

func (s *TestSuite) TestMeals_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationRequestInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.main.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating meals")
			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.main)

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
