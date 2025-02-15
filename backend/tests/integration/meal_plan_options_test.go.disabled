package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

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

func (s *TestSuite) TestMealPlanOptions_CompleteLifecycle() {
	s.runTest("should CRUD", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.adminClient, testClients.userClient)

			require.NotEmpty(t, createdMealPlan.Events)
			require.NotEmpty(t, createdMealPlan.Events[0].Options)

			createdMealPlanEvent := createdMealPlan.Events[0]
			createdMealPlanOption := createdMealPlanEvent.Options[0]
			require.NotNil(t, createdMealPlanOption)

			newMealPlanOption := fakes.BuildFakeMealPlanOption()
			newMealPlanOption.Meal.ID = createdMealPlanOption.Meal.ID
			newMealPlanOption.BelongsToMealPlanEvent = createdMealPlanEvent.ID
			newMealPlanOption.AssignedCook = nil

			updateInput := converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(newMealPlanOption)
			createdMealPlanOption.Update(updateInput)
			require.NoError(t, testClients.userClient.UpdateMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID, updateInput))

			actual, err := testClients.userClient.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, newMealPlanOption, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			require.NoError(t, testClients.userClient.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID))

			require.NoError(t, testClients.userClient.ArchiveMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID))

			require.NoError(t, testClients.userClient.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}

func (s *TestSuite) TestMealPlanOptions_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlan.Events = []*types.MealPlanEvent{exampleMealPlan.Events[0]}
			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.adminClient, testClients.userClient)

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

				createdMeal := createMealForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)
				exampleMealPlanOption.Meal.ID = createdMeal.ID

				exampleMealPlanOptionInput := converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(exampleMealPlanOption)
				newlyCreatedMealPlanOption, err := testClients.userClient.CreateMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, exampleMealPlanOptionInput)
				require.NoError(t, err)

				checkMealPlanOptionEquality(t, exampleMealPlanOption, newlyCreatedMealPlanOption)

				newlyCreatedMealPlanOption, err = testClients.userClient.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID)
				requireNotNilAndNoProblems(t, newlyCreatedMealPlanOption, err)
				require.Equal(t, createdMealPlanEvent.ID, newlyCreatedMealPlanOption.BelongsToMealPlanEvent)

				expected = append(expected, newlyCreatedMealPlanOption)
			}

			// assert meal plan option list equality
			actual, err := testClients.userClient.GetMealPlanOptions(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			assert.NoError(t, testClients.userClient.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}
