package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanOptionVoteEquality(t *testing.T, expected, actual *types.MealPlanOptionVote) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Rank, actual.Rank, "expected Rank for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Rank, actual.Rank)
	assert.Equal(t, expected.Abstain, actual.Abstain, "expected Abstain for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Abstain, actual.Abstain)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan option vote %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput creates an MealPlanOptionVoteUpdateRequestInput struct from a meal plan option vote.
func convertMealPlanOptionVoteToMealPlanOptionVoteUpdateInput(x *types.MealPlanOptionVote) *types.MealPlanOptionVoteUpdateRequestInput {
	return &types.MealPlanOptionVoteUpdateRequestInput{
		Rank:                    x.Rank,
		Abstain:                 x.Abstain,
		Notes:                   x.Notes,
		BelongsToMealPlanOption: x.BelongsToMealPlanOption,
	}
}

func (s *TestSuite) TestMealPlanOptionVotes_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdMealPlan := createMealPlanWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var createdMealPlanOption *types.MealPlanOption
			for _, opt := range createdMealPlan.Options {
				createdMealPlanOption = opt
				break
			}
			require.NotNil(t, createdMealPlanOption)

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			n = <-notificationsChan
			assert.Equal(t, types.MealPlanOptionVoteDataType, n.DataType)
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

			// one for the option vote
			n = <-notificationsChan
			// one for the option
			n = <-notificationsChan
			// can't predict which order :'(

			t.Log("fetching changed meal plan option vote")
			actual, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert meal plan option vote equality
			checkMealPlanOptionVoteEquality(t, newMealPlanOptionVote, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up meal plan option vote")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID))

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
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
			notificationsChan, err := testClients.main.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdMealPlan := createMealPlanWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var createdMealPlanOption *types.MealPlanOption
			for _, opt := range createdMealPlan.Options {
				createdMealPlanOption = opt
				break
			}
			require.NotNil(t, createdMealPlanOption)

			t.Log("creating meal plan option vote")
			exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
			exampleMealPlanOptionVote.BelongsToMealPlanOption = createdMealPlanOption.ID
			exampleMealPlanOptionVoteInput := fakes.BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(exampleMealPlanOptionVote)
			createdMealPlanOptionVoteID, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, exampleMealPlanOptionVoteInput)
			require.NoError(t, err)
			t.Logf("meal plan option vote %q created", createdMealPlanOptionVoteID)

			n = <-notificationsChan
			assert.Equal(t, types.MealPlanOptionVoteDataType, n.DataType)
			require.NotNil(t, n.MealPlanOptionVote)
			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, n.MealPlanOptionVote)

			createdMealPlanOptionVote, err := testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVoteID)
			requireNotNilAndNoProblems(t, createdMealPlanOptionVote, err)
			require.Equal(t, createdMealPlanOption.ID, createdMealPlanOptionVote.BelongsToMealPlanOption)

			checkMealPlanOptionVoteEquality(t, exampleMealPlanOptionVote, createdMealPlanOptionVote)

			// assert meal plan option vote list equality
			actual, err := testClients.main.GetMealPlanOptionVotes(ctx, createdMealPlan.ID, createdMealPlanOption.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.NotEmpty(t, actual.MealPlanOptionVotes)

			t.Log("cleaning up")
			assert.NoError(t, testClients.main.ArchiveMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanOption.ID, createdMealPlanOptionVote.ID))

			t.Log("cleaning up meal plan option")
			assert.NoError(t, testClients.main.ArchiveMealPlanOption(ctx, createdMealPlan.ID, createdMealPlanOption.ID))

			t.Log("cleaning up meal plan")
			assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
		}
	})
}
