package integration

import (
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealPlanOptionEquality(t *testing.T, expected, actual *types.MealPlanOption) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Meal.ID, actual.Meal.ID, "expected RecipeID for meal plan option %s to be %v, but it was %v", expected.ID, expected.Meal.ID, actual.Meal.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for meal plan option %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.AssignedCook, actual.AssignedCook, "expected AssignedCook for meal plan option %s to be %v, but it was %v", expected.ID, expected.AssignedCook, actual.AssignedCook)
	assert.Equal(t, expected.AssignedDishwasher, actual.AssignedDishwasher, "expected AssignedDishwasher for meal plan option %s to be %v, but it was %v", expected.ID, expected.AssignedDishwasher, actual.AssignedDishwasher)
	assert.NotZero(t, actual.CreatedAt)
}

func TestMealPlanOptions_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		require.NotEmpty(t, createdMealPlan.Events)
		require.NotEmpty(t, createdMealPlan.Events[0].Options)

		createdMealPlanEvent := createdMealPlan.Events[0]
		createdMealPlanOption := createdMealPlanEvent.Options[0]
		require.NotNil(t, createdMealPlanOption)

		newMealPlanOption := fakes.BuildFakeMealPlanOption()
		newMealPlanOption.Meal.ID = createdMealPlanOption.Meal.ID
		newMealPlanOption.BelongsToMealPlanEvent = createdMealPlanEvent.ID
		newMealPlanOption.AssignedCook = nil

		updateInput := mpconverters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(newMealPlanOption)
		createdMealPlanOption.Update(updateInput)
		_, err := userClient.UpdateMealPlanOption(ctx, &mealplanninggrpc.UpdateMealPlanOptionRequest{
			MealPlanID:       createdMealPlan.ID,
			MealPlanOptionID: createdMealPlanOption.ID,
			MealPlanEventID:  createdMealPlanEvent.ID,
			Input:            converters.ConvertMealPlanOptionUpdateRequestInputToGRPCMealPlanOptionUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		actualRes, err := userClient.GetMealPlanOption(ctx, &mealplanninggrpc.GetMealPlanOptionRequest{
			MealPlanID:       createdMealPlan.ID,
			MealPlanEventID:  createdMealPlanEvent.ID,
			MealPlanOptionID: createdMealPlanOption.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, actualRes)

		actual := converters.ConvertGRPCMealPlanOptionToMealPlanOption(actualRes.Result)

		// assert meal plan option equality
		checkMealPlanOptionEquality(t, newMealPlanOption, actual)
		assert.NotNil(t, actual.LastUpdatedAt)

		_, err = userClient.ArchiveMealPlanOption(ctx, &mealplanninggrpc.ArchiveMealPlanOptionRequest{
			MealPlanID:       createdMealPlan.ID,
			MealPlanEventID:  createdMealPlanEvent.ID,
			MealPlanOptionID: createdMealPlanOption.ID,
		})
		require.NoError(t, err)

		_, err = userClient.ArchiveMealPlanEvent(ctx, &mealplanninggrpc.ArchiveMealPlanEventRequest{
			MealPlanID:      createdMealPlan.ID,
			MealPlanEventID: createdMealPlanEvent.ID,
		})
		require.NoError(t, err)

		_, err = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanID: createdMealPlan.ID})
		require.NoError(t, err)
	})
}

func TestMealPlanOptions_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		require.NotEmpty(t, createdMealPlan.Events)
		require.NotEmpty(t, createdMealPlan.Events[0].Options)

		createdMealPlanEvent := createdMealPlan.Events[0]
		createdMealPlanOption := createdMealPlanEvent.Options[0]
		require.NotNil(t, createdMealPlanOption)

		var expected []*types.MealPlanOption
		for i := 0; i < 5; i++ {
			exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
			exampleMealPlanOption.Meal.ID = createdMealPlanOption.Meal.ID
			exampleMealPlanOption.BelongsToMealPlanEvent = createdMealPlanEvent.ID
			exampleMealPlanOption.AssignedCook = nil

			createdMeal := createMealForTest(t, userClient, nil)
			exampleMealPlanOption.Meal.ID = createdMeal.ID

			exampleMealPlanOptionInput := mpconverters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(exampleMealPlanOption)
			createdMealPlanOptionRes, err := userClient.CreateMealPlanOption(ctx, &mealplanninggrpc.CreateMealPlanOptionRequest{
				MealPlanID:      createdMealPlan.ID,
				MealPlanEventID: createdMealPlanEvent.ID,
				Input:           converters.ConvertMealPlanOptionCreationRequestInputToGRPCMealPlanOptionCreationRequestInput(exampleMealPlanOptionInput),
			})
			require.NoError(t, err)

			newlyCreatedMealPlanOption := converters.ConvertGRPCMealPlanOptionToMealPlanOption(createdMealPlanOptionRes.Created)
			checkMealPlanOptionEquality(t, exampleMealPlanOption, newlyCreatedMealPlanOption)

			retrievedMealPlanOptionRes, err := userClient.GetMealPlanOption(ctx, &mealplanninggrpc.GetMealPlanOptionRequest{
				MealPlanID:       createdMealPlan.ID,
				MealPlanEventID:  createdMealPlanEvent.ID,
				MealPlanOptionID: newlyCreatedMealPlanOption.ID,
			})
			require.NotNil(t, retrievedMealPlanOptionRes)
			require.NoError(t, err)
			require.Equal(t, createdMealPlanEvent.ID, retrievedMealPlanOptionRes.Result.BelongsToMealPlanEvent)

			expected = append(expected, converters.ConvertGRPCMealPlanOptionToMealPlanOption(retrievedMealPlanOptionRes.Result))
		}

		// assert meal plan option list equality
		actual, err := userClient.GetMealPlanOptions(ctx, &mealplanninggrpc.GetMealPlanOptionsRequest{
			Filter:          nil,
			MealPlanID:      createdMealPlan.ID,
			MealPlanEventID: createdMealPlanEvent.ID,
		})
		require.NotNil(t, actual)
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		_, err = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanID: createdMealPlan.ID})
		assert.NoError(t, err)
	})
}
