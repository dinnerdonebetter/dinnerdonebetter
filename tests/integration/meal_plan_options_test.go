package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanOptionEquality(t *testing.T, expected, actual *types.MealPlanOption) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Day, actual.Day, "expected Day for meal plan option %s to be %v, but it was %v", expected.ID, expected.Day, actual.Day)
	assert.Equal(t, expected.Meal.ID, actual.Meal.ID, "expected MealID for meal plan option %s to be %v, but it was %v", expected.ID, expected.Meal.ID, actual.Meal.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanOptionToMealPlanOptionUpdateInput creates an MealPlanOptionUpdateRequestInput struct from a meal plan option.
func convertMealPlanOptionToMealPlanOptionUpdateInput(x *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		Day:    &x.Day,
		MealID: &x.Meal.ID,
		Notes:  &x.Notes,
	}
}

func (s *TestSuite) TestMealPlanOptions_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, testClients.admin, testClients.user)

			var createdMealPlanOption *types.MealPlanOption
			for _, opt := range createdMealPlan.Options {
				createdMealPlanOption = opt
				break
			}
			require.NotNil(t, createdMealPlanOption)

			t.Log("changing meal plan option")
			newMealPlanOption := fakes.BuildFakeMealPlanOption()
			newMealPlanOption.Meal.ID = createdMealPlanOption.Meal.ID
			newMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
			createdMealPlanOption.Update(convertMealPlanOptionToMealPlanOptionUpdateInput(newMealPlanOption))
			assert.NoError(t, testClients.user.UpdateMealPlanOption(ctx, createdMealPlanOption))

			t.Log("fetching changed meal plan option")
			actual, err := testClients.user.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, newMealPlanOption, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.user.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}

func (s *TestSuite) TestMealPlanOptions_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdMealPlan := createMealPlanForTest(ctx, t, testClients.admin, testClients.user)

			t.Log("creating meal plan options")
			var expected []*types.MealPlanOption
			for i := 0; i < 5; i++ {
				exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
				exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID

				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)
				exampleMealPlanOption.Meal.ID = createdMeal.ID

				exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
				createdMealPlanOption, err := testClients.user.CreateMealPlanOption(ctx, exampleMealPlanOptionInput)
				require.NoError(t, err)
				t.Logf("meal plan option %q created", createdMealPlanOption.ID)

				checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

				createdMealPlanOption, err = testClients.user.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID)
				requireNotNilAndNoProblems(t, createdMealPlanOption, err)
				require.Equal(t, createdMealPlan.ID, createdMealPlanOption.BelongsToMealPlan)

				expected = append(expected, createdMealPlanOption)
			}

			// assert meal plan option list equality
			actual, err := testClients.user.GetMealPlanOptions(ctx, createdMealPlan.ID, nil)
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
				assert.NoError(t, testClients.user.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))
			}

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}
