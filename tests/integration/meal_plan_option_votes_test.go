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
	assert.Equal(t, expected.MealPlanOptionID, actual.MealPlanOptionID, "expected MealPlanOptionID for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.MealPlanOptionID, actual.MealPlanOptionID)
	assert.Equal(t, expected.DayOfWeek, actual.DayOfWeek, "expected DayOfWeek for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.DayOfWeek, actual.DayOfWeek)
	assert.Equal(t, expected.Points, actual.Points, "expected Points for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Points, actual.Points)
	assert.Equal(t, expected.Abstain, actual.Abstain, "expected Abstain for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Abstain, actual.Abstain)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput creates an MealPlanOptionVoteUpdateRequestInput struct from a meal plan option vote.
func convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(x *types.MealPlanOptionVote) *types.MealPlanOptionVoteUpdateRequestInput {
	return &types.MealPlanOptionVoteUpdateRequestInput{
		MealPlanOptionID: x.MealPlanOptionID,
		DayOfWeek:        x.DayOfWeek,
		Points:           x.Points,
		Abstain:          x.Abstain,
		Notes:            x.Notes,
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

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)
			require.NotNil(t, n.MealPlanOptionVote)
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, n.MealPlanOptionVote)

			createdMealPlanOptionVote, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, createdMealPlanOptionVote, err)

			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			t.Log("changing meal plan option vote")
			newMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			createdMealPlanOptionVote.Update(convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(newMealPlanOptionVote))
			assert.NoError(t, testClients.main.UpdateMealPlanOptionVote(ctx, createdMealPlanOptionVote))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)

			t.Log("fetching changed meal plan option vote")
			actual, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, newMealPlanOptionVote, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option vote")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlanOptionVoteID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			var createdMealPlanOptionVote *types.MealPlanOptionVote
			checkFunc = func() bool {
				createdMealPlanOptionVote, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
				return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			// change meal plan option vote
			newMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			createdMealPlanOptionVote.Update(convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(newMealPlanOptionVote))
			assert.NoError(t, testClients.main.UpdateMealPlanOptionVote(ctx, createdMealPlanOptionVote))

			time.Sleep(time.Second)

			// retrieve changed meal plan option vote
			var actual *types.MealPlanOptionVote
			checkFunc = func() bool {
				actual, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
				return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, newMealPlanOptionVote, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option vote")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlanOptionVoteID))
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

			t.Log("creating meal plan option votes")
			var expected []*types.MealPlanOptionVote
			for i := 0; i < 5; i++ {
				exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
				exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
				createdMealPlanOptionVoteID, createdMealPlanOptionVoteErr := testClients.main.CreateMealPlanOptionVote(ctx, exampleMealPlanOptionVoteInput)
				require.NoError(t, createdMealPlanOptionVoteErr)
				t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.MealPlanOptionVoteDataType)
				require.NotNil(t, n.MealPlanOptionVote)
				checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, n.MealPlanOptionVote)

				createdMealPlanOptionVote, createdMealPlanOptionVoteErr := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
				requireNotNilAndNoProblems(t, createdMealPlanOptionVote, createdMealPlanOptionVoteErr)

				expected = append(expected, createdMealPlanOptionVote)
			}

			// assert meal plan option vote list equality
			actual, err := testClients.main.GetMealPlanOptionVotes(ctx, nil)
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
				assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlanOptionVote.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating meal plan option votes")
			var expected []*types.MealPlanOptionVote
			for i := 0; i < 5; i++ {
				exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
				exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
				createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, exampleMealPlanOptionVoteInput)
				require.NoError(t, err)

				var createdMealPlanOptionVote *types.MealPlanOptionVote
				checkFunc = func() bool {
					createdMealPlanOptionVote, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlanOptionVoteID)
					return assert.NotNil(t, createdMealPlanOptionVote) && assert.NoError(t, err)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

				expected = append(expected, createdMealPlanOptionVote)
			}

			// assert meal plan option vote list equality
			actual, err := testClients.main.GetMealPlanOptionVotes(ctx, nil)
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
				assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlanOptionVote.ID))
			}
		}
	})
}
