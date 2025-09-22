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

func checkMealPlanTaskEquality(t *testing.T, expected, actual *mealplanning.MealPlanTask) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.CreationExplanation, actual.CreationExplanation, "expected CreationExplanation for meal plan %s to be %v, but it was %v", expected.CreationExplanation, expected.CreationExplanation, actual.CreationExplanation)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan %s to be %v, but it was %v", expected.Status, expected.Status, actual.Status)
	assert.Equal(t, expected.StatusExplanation, actual.StatusExplanation, "expected StatusExplanation for meal plan %s to be %v, but it was %v", expected.StatusExplanation, expected.StatusExplanation, actual.StatusExplanation)
	assert.Equal(t, expected.AssignedToUser, actual.AssignedToUser, "expected AssignedToUser for meal plan %s to be %v, but it was %v", expected.AssignedToUser, expected.AssignedToUser, actual.AssignedToUser)
	assert.Equal(t, expected.CompletedAt, actual.CompletedAt, "expected CompletedAt for meal plan %s to be %v, but it was %v", expected.CompletedAt, expected.CompletedAt, actual.CompletedAt)

	assert.NotZero(t, actual.CreatedAt)
}

func TestMealPlanTasks_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		createdMealPlan := createMealPlanForTest(t, adminClient, nil)

		exampleMealPlanTask := fakes.BuildFakeMealPlanTask()
		exampleMealPlanTaskInput := mpconverters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(exampleMealPlanTask)

		exampleMealPlanTaskInput.MealPlanOptionID = createdMealPlan.Events[0].Options[0].ID
		exampleMealPlanTaskInput.RecipePrepTaskID = createdMealPlan.Events[0].Options[0].Meal.Components[0].Recipe.PrepTasks[0].ID

		createdMealPlanTaskRes, err := adminClient.CreateMealPlanTask(ctx, &mealplanninggrpc.CreateMealPlanTaskRequest{
			MealPlanID: createdMealPlan.ID,
			Input:      converters.ConvertMealPlanTaskCreationRequestInputToGRPCMealPlanTaskCreationRequestInput(exampleMealPlanTaskInput),
		})
		require.NoError(t, err)

		createdMealPlanTask := converters.ConvertGRPCMealPlanTaskToMealPlanTask(createdMealPlanTaskRes.Created)
		checkMealPlanTaskEquality(t, exampleMealPlanTask, createdMealPlanTask)

		actualRes, err := adminClient.GetMealPlanTask(ctx, &mealplanninggrpc.GetMealPlanTaskRequest{
			MealPlanID:     createdMealPlan.ID,
			MealPlanTaskID: createdMealPlanTask.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, actualRes)

		actual := converters.ConvertGRPCMealPlanTaskToMealPlanTask(actualRes.Result)
		// assert meal plan task equality
		checkMealPlanTaskEquality(t, exampleMealPlanTask, actual)
	})
}
