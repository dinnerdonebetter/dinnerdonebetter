package integration

import (
	"context"
	"fmt"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/apiclient"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func checkMealEquality(t *testing.T, expected, actual *types.Meal) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for meal %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for meal %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.NotZero(t, actual.CreatedAt)
}

func createMealForTest(ctx context.Context, t *testing.T, adminClient, client *apiclient.Client, mealInput *types.Meal) *types.Meal {
	t.Helper()

	createdRecipes := []*types.Recipe{}
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeForTest(ctx, t, adminClient, client, nil)
		createdRecipes = append(createdRecipes, recipe)
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
	}

	t.Log("creating meal")
	exampleMeal := mealInput
	if exampleMeal == nil {
		exampleMeal = fakes.BuildFakeMeal()
	}

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

			createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

			t.Log("cleaning up meal")
			assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
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
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)

			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating meals")
			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)

				expected = append(expected, createdMeal)
			}

			// assert meal list equality
			actual, err := testClients.user.GetMeals(ctx, nil)
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
				assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
			}
		}
	})
}

func (s *TestSuite) TestMeals_Searching() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.user.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("creating prerequisite valid preparation")
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)
			t.Logf("valid preparation %q created", createdValidPreparation.ID)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			t.Log("creating meals")
			exampleMeal := fakes.BuildFakeMeal()
			var expected []*types.Meal
			for i := 0; i < 5; i++ {
				exampleMeal.Name = fmt.Sprintf("example%d", i)
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, exampleMeal)

				expected = append(expected, createdMeal)
			}

			// assert meal list equality
			actual, err := testClients.user.SearchForMeals(ctx, "example", nil)
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
				assert.NoError(t, testClients.user.ArchiveMeal(ctx, createdMeal.ID))
			}
		}
	})
}
