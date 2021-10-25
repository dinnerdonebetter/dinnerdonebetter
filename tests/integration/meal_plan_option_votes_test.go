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

func checkMealPlanOptionVoteEquality(t *testing.T, expected, actual *types.MealPlanOptionVote) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Points, actual.Points, "expected Points for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Points, actual.Points)
	assert.Equal(t, expected.Abstain, actual.Abstain, "expected Abstain for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Abstain, actual.Abstain)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput creates an MealPlanOptionVoteUpdateRequestInput struct from a meal plan option vote.
func convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(x *types.MealPlanOptionVote) *types.MealPlanOptionVoteUpdateRequestInput {
	return &types.MealPlanOptionVoteUpdateRequestInput{
		Points:  x.Points,
		Abstain: x.Abstain,
		Notes:   x.Notes,
	}
}

func (s *TestSuite) TestMealPlanOptionVotes_CompleteLifecycle() {
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

			t.Log("creating prerequisite meal plan option")
			exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
			exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
			//exampleMealPlanOption.RecipeID = createdRecipe.ID
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

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)
			require.NotNil(t, n.MealPlanOptionVote)
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, n.MealPlanOptionVote)

			createdMealPlanOptionVote, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, createdMealPlanOptionVote, err)
			require.Equal(t, createdMealPlanOption.ID, createdMealPlanOptionVote.BelongsToMealPlanOption)

			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			t.Log("changing meal plan option vote")
			newMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			createdMealPlanOptionVote.Update(convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(newMealPlanOptionVote))
			assert.NoError(t, testClients.main.UpdateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOptionVote))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)

			t.Log("fetching changed meal plan option vote")
			actual, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, newMealPlanOptionVote, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option vote")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID))

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

			_, _, createdRecipe := createRecipeWithPolling(ctx, t, testClients.main)

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

			t.Log("creating prerequisite meal plan option")
			exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
			exampleMealPlanOption.BelongsToMealPlan = createdMealPlan.ID
			exampleMealPlanOption.RecipeID = createdRecipe.ID
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

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			var createdMealPlanOptionVote *types.MealPlanOptionVote
			checkFunc = func() bool {
				createdMealPlanOptionVote, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
				return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdMealPlanOption.ID, createdMealPlanOptionVote.BelongsToMealPlanOption)
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			// change meal plan option vote
			newMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			createdMealPlanOptionVote.Update(convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(newMealPlanOptionVote))
			assert.NoError(t, testClients.main.UpdateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOptionVote))

			time.Sleep(time.Second)

			// retrieve changed meal plan option vote
			var actual *types.MealPlanOptionVote
			checkFunc = func() bool {
				actual, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
				return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, newMealPlanOptionVote, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option vote")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID))

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})
}

func (s *TestSuite) TestMealPlanOptionVotes_Listing() {
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

			t.Log("creating prerequisite meal plan option")
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

			t.Log("creating meal plan option votes")
			var expected []*types.MealPlanOptionVote
			for i := 0; i < 5; i++ {
				exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
				exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
				exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
				createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
				require.NoError(t, err)
				t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)
				require.NotNil(t, n.MealPlanOptionVote)
				checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, n.MealPlanOptionVote)

				createdMealPlanOptionVote, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
				requireNotNilAndNoProblems(t, createdMealPlanOptionVote, err)
				require.Equal(t, createdMealPlanOption.ID, createdMealPlanOptionVote.BelongsToMealPlanOption)

				expected = append(expected, createdMealPlanOptionVote)
			}

			// assert meal plan option vote list equality
			actual, err := testClients.main.GetMealPlanOptionVotes(ctx, createdMealPlan.ID, createdMealPlanOption.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlanOptionVotes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlanOptionVotes),
			)

			t.Log("cleaning up")
			for _, createdMealPlanOptionVote := range expected {
				assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVote.ID))
			}

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID))

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

			t.Log("creating prerequisite meal plan option")
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

			t.Log("creating meal plan option votes")
			var expected []*types.MealPlanOptionVote
			for i := 0; i < 5; i++ {
				exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
				exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
				exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
				createdMealPlanOptionVoteID, mealPlanOptionVoteErr := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
				require.NoError(t, mealPlanOptionVoteErr)

				var createdMealPlanOptionVote *types.MealPlanOptionVote
				checkFunc = func() bool {
					createdMealPlanOptionVote, mealPlanOptionVoteErr = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
					return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, mealPlanOptionVoteErr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

				expected = append(expected, createdMealPlanOptionVote)
			}

			// assert meal plan option vote list equality
			actual, err := testClients.main.GetMealPlanOptionVotes(ctx, createdMealPlan.ID, createdMealPlanOption.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlanOptionVotes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlanOptionVotes),
			)

			t.Log("cleaning up")
			for _, createdMealPlanOptionVote := range expected {
				assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVote.ID))
			}

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOptionID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlanID))
		}
	})
}
