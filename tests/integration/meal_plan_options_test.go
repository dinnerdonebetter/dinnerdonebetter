package integration

import (
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/stretchr/testify/assert"
)

func checkMealPlanOptionEquality(t *testing.T, expected, actual *types.MealPlanOption) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Meal.ID, actual.Meal.ID, "expected MealID for meal plan option %s to be %v, but it was %v", expected.ID, expected.Meal.ID, actual.Meal.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.AssignedCook, actual.AssignedCook, "expected AssignedCook for meal plan option %s to be %v, but it was %v", expected.ID, expected.AssignedCook, actual.AssignedCook)
	assert.Equal(t, expected.AssignedDishwasher, actual.AssignedDishwasher, "expected AssignedDishwasher for meal plan option %s to be %v, but it was %v", expected.ID, expected.AssignedDishwasher, actual.AssignedDishwasher)
	assert.Equal(t, expected.PrepStepsCreated, actual.PrepStepsCreated, "expected PrepStepsCreated for meal plan option %s to be %v, but it was %v", expected.ID, expected.PrepStepsCreated, actual.PrepStepsCreated)
	assert.NotZero(t, actual.CreatedAt)
}

// convertMealPlanOptionToMealPlanOptionUpdateInput creates an MealPlanOptionUpdateRequestInput struct from a meal plan option.
func convertMealPlanOptionToMealPlanOptionUpdateInput(x *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		MealID:             &x.Meal.ID,
		Notes:              &x.Notes,
		AssignedCook:       x.AssignedCook,
		AssignedDishwasher: x.AssignedDishwasher,
		PrepStepsCreated:   &x.PrepStepsCreated,
	}
}

func (s *TestSuite) TestMealPlanOptions_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			require.NotEmpty(t, createdMealPlan.Events)
			require.NotEmpty(t, createdMealPlan.Events[0].Options)

			createdMealPlanEvent := createdMealPlan.Events[0]
			createdMealPlanOption := createdMealPlanEvent.Options[0]
			require.NotNil(t, createdMealPlanOption)

			t.Log("changing meal plan option")
			newMealPlanOption := fakes.BuildFakeMealPlanOption()
			newMealPlanOption.Meal.ID = createdMealPlanOption.Meal.ID
			newMealPlanOption.BelongsToMealPlanEvent = createdMealPlanEvent.ID
			newMealPlanOption.AssignedCook = nil

			createdMealPlanOption.Update(convertMealPlanOptionToMealPlanOptionUpdateInput(newMealPlanOption))
			require.NoError(t, testClients.user.UpdateMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption))

			t.Log("fetching changed meal plan option")
			actual, err := testClients.user.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, newMealPlanOption, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up meal plan option")
			require.NoError(t, testClients.user.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID))

			t.Log("cleaning up meal plan event")
			require.NoError(t, testClients.user.ArchiveMealPlanEvent(ctx, createdMealPlan.ID, createdMealPlanEvent.ID))

			t.Log("cleaning up meal plan")
			require.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}

func (s *TestSuite) TestMealPlanOptions_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)

			require.NotEmpty(t, createdMealPlan.Events)
			require.NotEmpty(t, createdMealPlan.Events[0].Options)

			createdMealPlanEvent := createdMealPlan.Events[0]

			t.Log("creating meal plan options")
			var expected []*types.MealPlanOption
			for i := 0; i < 5; i++ {
				exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
				exampleMealPlanOption.BelongsToMealPlanEvent = createdMealPlan.ID
				exampleMealPlanOption.AssignedCook = nil

				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)
				exampleMealPlanOption.Meal.ID = createdMeal.ID

				exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
				createdMealPlanOption, err := testClients.user.CreateMealPlanOption(ctx, createdMealPlan.ID, exampleMealPlanOptionInput)
				require.NoError(t, err)
				t.Logf("meal plan option %q created", createdMealPlanOption.ID)

				checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

				createdMealPlanOption, err = testClients.user.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID)
				requireNotNilAndNoProblems(t, createdMealPlanOption, err)
				require.Equal(t, createdMealPlan.ID, createdMealPlanOption.BelongsToMealPlanEvent)

				expected = append(expected, createdMealPlanOption)
			}

			// assert meal plan option list equality
			actual, err := testClients.user.GetMealPlanOptions(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlanOptions),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlanOptions),
			)

			t.Log("cleaning up")
			for _, createdMealPlanOption := range expected {
				assert.NoError(t, testClients.user.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, createdMealPlanOption.ID))
			}

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}
