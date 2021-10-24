package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func checkMealPlanOptionEquality(t *testing.T, expected, actual *types.MealPlanOption) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.DayOfWeek, actual.DayOfWeek, "expected DayOfWeek for meal plan option %s to be %v, but it was %v", expected.ID, expected.DayOfWeek, actual.DayOfWeek)
	assert.Equal(t, expected.RecipeID, actual.RecipeID, "expected RecipeID for meal plan option %s to be %v, but it was %v", expected.ID, expected.RecipeID, actual.RecipeID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanOptionToMealPlanOptionUpdateInput creates an MealPlanOptionUpdateRequestInput struct from a meal plan option.
func convertMealPlanOptionToMealPlanOptionUpdateInput(x *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		DayOfWeek: x.DayOfWeek,
		RecipeID:  x.RecipeID,
		Notes:     x.Notes,
	}
}

func (s *TestSuite) TestMealPlanOptions_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite meal plan")
			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlanID, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlanID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanDataType)
			require.NotNil(t, n.MealPlan)
			checkMealPlanEquality(t, exampleMealPlan, n.MealPlan)

			createdMealPlan, err := testClients.main.GetMealPlan(ctx, createdMealPlanID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)

			t.Log("creating meal plan option")
			exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
			exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
			exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
			createdMealPlanOptionID, err := testClients.main.CreateMealPlanOption(ctx, exampleMealPlanOptionInput)
			require.NoError(t, err)
			t.Logf("meal plan option %q created", createdMealPlanOptionID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionDataType)
			require.NotNil(t, n.MealPlanOption)
			checkMealPlanOptionEquality(t, exampleMealPlanOption, n.MealPlanOption)

			createdMealPlanOption, err := testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
			requireNotNilAndNoProblems(t, createdMealPlanOption, err)
			require.Equal(t, createdMealPlan.ID, createdMealPlanOption.BelongsToMealPlan)

			checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

			t.Log("changing meal plan option")
			newMealPlanOption := fakes.BuildFakeMealPlanOption()
			createdMealPlanOption.Update(convertMealPlanOptionToMealPlanOptionUpdateInput(newMealPlanOption))
			assert.NoError(t, testClients.main.UpdateMealPlanOption(ctx, createdMealPlanOption))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionDataType)

			t.Log("fetching changed meal plan option")
			actual, err := testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, newMealPlanOption, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite meal plan")
			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlanID, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlanID)

			var createdMealPlan *types.MealPlan
			checkFunc = func() bool {
				createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlanID)
				return assert.NotNil(t, createdMealPlan) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			t.Log("creating meal plan option")
			exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
			exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
			exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
			createdMealPlanOptionID, err := testClients.main.CreateMealPlanOption(ctx, exampleMealPlanOptionInput)
			require.NoError(t, err)
			t.Logf("meal plan option %q created", createdMealPlanOptionID)

			var createdMealPlanOption *types.MealPlanOption
			checkFunc = func() bool {
				createdMealPlanOption, err = testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
				return assert.NotNil(t, createdMealPlanOption) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdMealPlan.ID, createdMealPlanOption.BelongsToMealPlan)
			checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

			// change meal plan option
			newMealPlanOption := fakes.BuildFakeMealPlanOption()
			createdMealPlanOption.Update(convertMealPlanOptionToMealPlanOptionUpdateInput(newMealPlanOption))
			assert.NoError(t, testClients.main.UpdateMealPlanOption(ctx, createdMealPlanOption))

			time.Sleep(time.Second)

			// retrieve changed meal plan option
			var actual *types.MealPlanOption
			checkFunc = func() bool {
				actual, err = testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
				return assert.NotNil(t, createdMealPlanOption) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option equality
			checkMealPlanOptionEquality(t, newMealPlanOption, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})
}

func (s *TestSuite) TestMealPlanOptions_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite meal plan")
			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlanID, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlanID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanDataType)
			require.NotNil(t, n.MealPlan)
			checkMealPlanEquality(t, exampleMealPlan, n.MealPlan)

			createdMealPlan, err := testClients.main.GetMealPlan(ctx, createdMealPlanID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)

			t.Log("creating meal plan options")
			var expected []*types.MealPlanOption
			for i := 0; i < 5; i++ {
				exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
				exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
				exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
				createdMealPlanOptionID, err := testClients.main.CreateMealPlanOption(ctx, exampleMealPlanOptionInput)
				require.NoError(t, err)
				t.Logf("meal plan option %q created", createdMealPlanOptionID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.MealPlanOptionDataType)
				require.NotNil(t, n.MealPlanOption)
				checkMealPlanOptionEquality(t, exampleMealPlanOption, n.MealPlanOption)

				createdMealPlanOption, err := testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
				requireNotNilAndNoProblems(t, createdMealPlanOption, err)
				require.Equal(t, createdMealPlan.ID, createdMealPlanOption.BelongsToMealPlan)

				expected = append(expected, createdMealPlanOption)
			}

			// assert meal plan option list equality
			actual, err := testClients.main.GetMealPlanOptions(ctx, createdMealPlan.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))
			}

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite meal plan")
			exampleMealPlan := fakes.BuildFakeMealPlan()
			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlanID, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlanID)

			var createdMealPlan *types.MealPlan
			checkFunc = func() bool {
				createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlanID)
				return assert.NotNil(t, createdMealPlan) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			t.Log("creating meal plan options")
			var expected []*types.MealPlanOption
			for i := 0; i < 5; i++ {
				exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
				exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
				exampleMealPlanOptionInput := fakes.BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(exampleMealPlanOption)
				createdMealPlanOptionID, mealPlanOptionErr := testClients.main.CreateMealPlanOption(ctx, exampleMealPlanOptionInput)
				require.NoError(t, mealPlanOptionErr)

				var createdMealPlanOption *types.MealPlanOption
				checkFunc = func() bool {
					createdMealPlanOption, mealPlanOptionErr = testClients.main.GetMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID)
					return assert.NotNil(t, createdMealPlanOption) && assert.NoError(t, mealPlanOptionErr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkMealPlanOptionEquality(t, exampleMealPlanOption, createdMealPlanOption)

				expected = append(expected, createdMealPlanOption)
			}

			// assert meal plan option list equality
			actual, err := testClients.main.GetMealPlanOptions(ctx, createdMealPlan.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))
			}

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})
}
