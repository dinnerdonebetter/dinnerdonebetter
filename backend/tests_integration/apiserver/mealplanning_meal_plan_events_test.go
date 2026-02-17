package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealPlanEventEquality(t *testing.T, expected, actual *mealplanning.MealPlanEvent) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for meal plan event %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.StartsAt, actual.StartsAt, "expected StartsAt for meal plan event %s to be %v, but it was %v", expected.ID, expected.StartsAt, actual.StartsAt)
	assert.Equal(t, expected.EndsAt, actual.EndsAt, "expected EndsAt for meal plan event %s to be %v, but it was %v", expected.ID, expected.EndsAt, actual.EndsAt)
	assert.Equal(t, expected.MealName, actual.MealName, "expected MealName for meal plan event %s to be %v, but it was %v", expected.ID, expected.MealName, actual.MealName)
	assert.Equal(t, expected.BelongsToMealPlan, actual.BelongsToMealPlan, "expected BelongsToMealPlan for meal plan event %s to be %v, but it was %v", expected.ID, expected.BelongsToMealPlan, actual.BelongsToMealPlan)
	assert.NotZero(t, actual.CreatedAt)
}

func TestMealPlanEvents_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		require.NotNil(t, createdMealPlan)
		require.NotEmpty(t, createdMealPlan.Events)
		require.NotNil(t, createdMealPlan.Events[0])
		createdMealPlanEvent := createdMealPlan.Events[0]

		newMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		newMealPlanEvent.BelongsToMealPlan = createdMealPlan.ID

		updateInput := mpconverters.ConvertMealPlanEventToMealPlanEventUpdateRequestInput(newMealPlanEvent)
		createdMealPlanEvent.Update(updateInput)

		_, err := userClient.UpdateMealPlanEvent(ctx, &mealplanninggrpc.UpdateMealPlanEventRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanEventUpdateRequestInputToGRPCMealPlanEventUpdateRequestInput(updateInput),
		})
		assert.NoError(t, err)

		actualRes, err := userClient.GetMealPlanEvent(ctx, &mealplanninggrpc.GetMealPlanEventRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
		})
		require.NoError(t, err)

		// assert meal plan event equality
		actual := converters.ConvertGRPCMealPlanEventToMealPlanEvent(actualRes.Result)
		checkMealPlanEventEquality(t, newMealPlanEvent, actual)
		assert.NotNil(t, actual.LastUpdatedAt)

		_, err = userClient.ArchiveMealPlanEvent(ctx, &mealplanninggrpc.ArchiveMealPlanEventRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
		})
		assert.NoError(t, err)

		_, err = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanId: createdMealPlan.ID})
		assert.NoError(t, err)
	})
}

func TestMealPlanEvents_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		var expected []*mealplanning.MealPlanEvent
		for range 5 {
			exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
			exampleMealPlanEvent.Options = nil
			exampleMealPlanEvent.BelongsToMealPlan = createdMealPlan.ID

			exampleMealPlanEventInput := mpconverters.ConvertMealPlanEventToMealPlanEventCreationRequestInput(exampleMealPlanEvent)
			createdMealPlanEventRes, err := userClient.CreateMealPlanEvent(ctx, &mealplanninggrpc.CreateMealPlanEventRequest{
				MealPlanId: createdMealPlan.ID,
				Input:      converters.ConvertMealPlanEventCreationRequestInputToGRPCMealPlanEventCreationRequestInput(exampleMealPlanEventInput),
			})
			require.NoError(t, err)

			createdMealPlanEvent := converters.ConvertGRPCMealPlanEventToMealPlanEvent(createdMealPlanEventRes.Created)
			checkMealPlanEventEquality(t, exampleMealPlanEvent, createdMealPlanEvent)

			retrievedMealPlanEventRes, err := userClient.GetMealPlanEvent(ctx, &mealplanninggrpc.GetMealPlanEventRequest{
				MealPlanId:      createdMealPlan.ID,
				MealPlanEventId: createdMealPlanEvent.ID,
			})
			require.NoError(t, err)

			createdMealPlanEvent = converters.ConvertGRPCMealPlanEventToMealPlanEvent(retrievedMealPlanEventRes.Result)
			require.Equal(t, createdMealPlan.ID, createdMealPlanEvent.BelongsToMealPlan)

			expected = append(expected, createdMealPlanEvent)
		}

		// assert meal plan event list equality
		actual, err := userClient.GetMealPlanEvents(ctx, &mealplanninggrpc.GetMealPlanEventsRequest{MealPlanId: createdMealPlan.ID})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdMealPlanEvent := range expected {
			_, err = userClient.ArchiveMealPlanEvent(ctx, &mealplanninggrpc.ArchiveMealPlanEventRequest{
				MealPlanId:      createdMealPlan.ID,
				MealPlanEventId: createdMealPlanEvent.ID,
			})
			assert.NoError(t, err)
		}

		_, err = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanId: createdMealPlan.ID})
		assert.NoError(t, err)
	})
}
