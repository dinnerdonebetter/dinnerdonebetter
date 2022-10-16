package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanTaskEquality(t *testing.T, expected, actual *types.MealPlanTask) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.CreationExplanation, actual.CreationExplanation, "expected CreationExplanation for meal plan %s to be %v, but it was %v", expected.CreationExplanation, expected.CreationExplanation, actual.CreationExplanation)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan %s to be %v, but it was %v", expected.Status, expected.Status, actual.Status)
	assert.Equal(t, expected.StatusExplanation, actual.StatusExplanation, "expected StatusExplanation for meal plan %s to be %v, but it was %v", expected.StatusExplanation, expected.StatusExplanation, actual.StatusExplanation)
	assert.Equal(t, expected.AssignedToUser, actual.AssignedToUser, "expected AssignedToUser for meal plan %s to be %v, but it was %v", expected.AssignedToUser, expected.AssignedToUser, actual.AssignedToUser)
	assert.Equal(t, expected.CompletedAt, actual.CompletedAt, "expected CompletedAt for meal plan %s to be %v, but it was %v", expected.CompletedAt, expected.CompletedAt, actual.CompletedAt)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestMealPlanTasks_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			t.Log("creating meal plan task")
			exampleMealPlanTask := fakes.BuildFakeMealPlanTask()
			exampleMealPlanTaskInput := fakes.BuildFakeMealPlanTaskCreationRequestInputFromMealPlanTask(exampleMealPlanTask)

			exampleMealPlanTaskInput.MealPlanOptionID = createdMealPlan.Events[0].Options[0].ID
			exampleMealPlanTaskInput.RecipePrepTaskID = createdMealPlan.Events[0].Options[0].Meal.Recipes[0].PrepTasks[0].ID

			createdMealPlanTask, err := testClients.admin.CreateMealPlanTask(ctx, createdMealPlan.ID, exampleMealPlanTaskInput)
			require.NoError(t, err)
			t.Logf("meal plan task %q created", createdMealPlanTask.ID)
			checkMealPlanTaskEquality(t, exampleMealPlanTask, createdMealPlanTask)

			t.Log("fetching changed meal plan task")
			actual, err := testClients.admin.GetMealPlanTask(ctx, createdMealPlan.ID, createdMealPlanTask.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan task equality
			checkMealPlanTaskEquality(t, exampleMealPlanTask, actual)
		}
	})
}
