package integration

import (
	"fmt"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealEquality(t *testing.T, expected, actual *types.Meal) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Name, actual.Name, "expected Name for meal %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for meal %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.EstimatedPortions, actual.EstimatedPortions, "expected EstimatedPortions for meal %s to be %v, but it was %v", expected.ID, expected.EstimatedPortions, actual.EstimatedPortions)
	assert.Equal(t, expected.EligibleForMealPlans, actual.EligibleForMealPlans, "expected EligibleForMealPlans for meal %s to be %v, but it was %v", expected.ID, expected.EligibleForMealPlans, actual.EligibleForMealPlans)

	assert.NotZero(t, actual.CreatedAt)
}

func createMealForTest(t *testing.T, clientToUse client.Client, mealInput *types.Meal) *types.Meal {
	t.Helper()

	ctx := t.Context()

	createdComponents := []*types.MealComponentCreationRequestInput{}
	createdRecipeIDs := []string{}
	for i := 0; i < 3; i++ {
		_, _, recipe := createRecipeForTest(t, nil)
		createdRecipeIDs = append(createdRecipeIDs, recipe.ID)
		createdComponents = append(createdComponents, &types.MealComponentCreationRequestInput{
			RecipeID:      recipe.ID,
			ComponentType: types.MealComponentTypesMain,
			RecipeScale:   1.0,
		})
	}

	exampleMeal := mealInput
	if exampleMeal == nil {
		exampleMeal = fakes.BuildFakeMeal()
	}

	exampleMealInput := mpconverters.ConvertMealToMealCreationRequestInput(exampleMeal)
	exampleMealInput.Components = createdComponents

	createdMealRes, err := clientToUse.CreateMeal(ctx, &mealplanninggrpc.CreateMealRequest{Input: converters.ConvertMealCreationRequestInputToGRPCMealCreationRequestInput(exampleMealInput)})
	require.NoError(t, err)

	fetchedMealRes, err := clientToUse.GetMeal(ctx, &mealplanninggrpc.GetMealRequest{MealID: createdMealRes.Created.ID})
	require.NoError(t, err)

	createdMeal := converters.ConvertGRPCMealToMeal(fetchedMealRes.Result)
	checkMealEquality(t, exampleMeal, createdMeal)

	componentRecipeIDs := []string{}
	for _, component := range createdMeal.Components {
		componentRecipeIDs = append(componentRecipeIDs, component.Recipe.ID)
	}

	assert.ElementsMatch(t, createdRecipeIDs, componentRecipeIDs)

	return createdMeal
}

func TestMeals_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		createdMeal := createMealForTest(t, userClient, nil)

		_, err := userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: createdMeal.ID})
		assert.NoError(t, err)
	})

	T.Run("requires auth for creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMealInput := mpconverters.ConvertMealToMealCreationRequestInput(exampleMeal)
		convertedInput := converters.ConvertMealCreationRequestInputToGRPCMealCreationRequestInput(exampleMealInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateMeal(ctx, &mealplanninggrpc.CreateMealRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestMeals_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		var expected []*types.Meal
		for i := 0; i < 5; i++ {
			createdMeal := createMealForTest(t, userClient, nil)

			expected = append(expected, createdMeal)
		}

		// assert meal list equality
		actual, err := userClient.GetMeals(ctx, &mealplanninggrpc.GetMealsRequest{})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdMeal := range expected {
			_, err = userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: createdMeal.ID})
			assert.NoError(t, err)
		}
	})

	T.Run("requires auth for listing", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		meals, err := c.GetMeals(ctx, &mealplanninggrpc.GetMealsRequest{})
		assert.Error(t, err)
		assert.Nil(t, meals)
	})
}

func TestMeals_Searching(T *testing.T) {
	T.Parallel()

	T.Run("should be searchable", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		exampleMeal := fakes.BuildFakeMeal()
		var expected []*types.Meal
		for i := 0; i < 5; i++ {
			exampleMeal.Name = fmt.Sprintf("example%d", i)
			createdMeal := createMealForTest(t, userClient, exampleMeal)

			expected = append(expected, createdMeal)
		}

		// assert meal list equality
		actual, err := userClient.SearchForMeals(ctx, &mealplanninggrpc.SearchForMealsRequest{
			Query:            "example",
			UseSearchService: false,
		})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdMeal := range expected {
			_, err = userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: createdMeal.ID})
			assert.NoError(t, err)
		}
	})

	T.Run("requires auth for searching", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)
		results, err := c.SearchForMeals(ctx, &mealplanninggrpc.SearchForMealsRequest{
			Query:            "test",
			UseSearchService: false,
		})
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestMeals_Reading(T *testing.T) {
	T.Parallel()

	T.Run("requires auth for reading", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMeal := createMealForTest(t, userClient, nil)

		c := buildUnauthenticatedGRPCClientForTest(t)
		meal, err := c.GetMeal(ctx, &mealplanninggrpc.GetMealRequest{MealID: createdMeal.ID})
		assert.Error(t, err)
		assert.Nil(t, meal)

		// Clean up
		_, err = userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: createdMeal.ID})
		assert.NoError(t, err)
	})

	T.Run("nonexistent meal", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		meal, err := userClient.GetMeal(ctx, &mealplanninggrpc.GetMealRequest{MealID: nonexistentID})
		assert.Error(t, err)
		assert.Nil(t, meal)
	})
}

func TestMeals_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("requires auth for archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMeal := createMealForTest(t, userClient, nil)

		c := buildUnauthenticatedGRPCClientForTest(t)
		_, err := c.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: createdMeal.ID})
		assert.Error(t, err)
	})

	T.Run("nonexistent meal for archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		_, err := userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealID: nonexistentID})
		assert.Error(t, err)
	})
}
