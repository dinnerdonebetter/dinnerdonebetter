package integration

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	authgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitygrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealPlanEquality(t *testing.T, expected, actual *mealplanning.MealPlan) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.WithinDuration(t, expected.VotingDeadline, actual.VotingDeadline, time.Nanosecond*1000, "expected VotingDeadline for meal plan %s to be %v, but it was %v", expected.ID, expected.VotingDeadline, actual.VotingDeadline)
	assert.Equal(t, expected.TasksCreated, actual.TasksCreated, "expected TasksCreated for meal plan %s to be %v, but it was %v", expected.ID, expected.TasksCreated, actual.TasksCreated)
	assert.Equal(t, expected.ElectionMethod, actual.ElectionMethod, "expected ElectionMethod for meal plan %s to be %v, but it was %v", expected.ID, expected.ElectionMethod, actual.ElectionMethod)
	assert.Equal(t, expected.GroceryListInitialized, actual.GroceryListInitialized, "expected GroceryListInitialized for meal plan %s to be %v, but it was %v", expected.ID, expected.GroceryListInitialized, actual.GroceryListInitialized)
	assert.NotZero(t, actual.CreatedAt)
}

func createMealPlanForTest(t *testing.T, clientToUse client.Client, mealPlan *mealplanning.MealPlan) *mealplanning.MealPlan {
	t.Helper()
	ctx := t.Context()

	if mealPlan == nil {
		mealPlan = fakes.BuildFakeMealPlan()
		for i, evt := range mealPlan.Events {
			for j := range evt.Options {
				createdMeal := createMealForTest(t, clientToUse, nil)
				mealPlan.Events[i].Options[j].Meal.ID = createdMeal.ID
				mealPlan.Events[i].Options[j].AssignedCook = nil
			}
		}
	}

	exampleMealPlanInput := mpconverters.ConvertMealPlanToMealPlanCreationRequestInput(mealPlan)
	createdMealPlanRes, err := clientToUse.CreateMealPlan(ctx, &mealplanninggrpc.CreateMealPlanRequest{
		Input: converters.ConvertMealPlanCreationRequestInputToGRPCMealPlanCreationRequestInput(exampleMealPlanInput),
	})
	require.NoError(t, err)
	require.NotEmpty(t, createdMealPlanRes.Created.Id)

	mealPlanRes, err := clientToUse.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanId: createdMealPlanRes.Created.Id})
	require.NoError(t, err)

	actual := converters.ConvertGRPCMealPlanToMealPlan(mealPlanRes.Result)
	checkMealPlanEquality(t, mealPlan, actual)

	return actual
}

func TestMealPlans_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should be readable in paginated form", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		var expected []*mealplanning.MealPlan
		for range 5 {
			createdMealPlan := createMealPlanForTest(t, userClient, nil)
			expected = append(expected, createdMealPlan)
		}

		// assert meal plan list equality
		actual, err := userClient.GetMealPlansForAccount(ctx, &mealplanninggrpc.GetMealPlansForAccountRequest{})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdMealPlan := range expected {
			_, err = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanId: createdMealPlan.ID})
			assert.NoError(t, err)
		}
	})
}

func TestMealPlans_CompleteLifecycleForAllVotesReceived(T *testing.T) {
	T.Parallel()

	T.Run("should resolve the meal plan status upon receiving all votes", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// create a userClient for the meal plan account
		_, accountAdminUserClient := createUserAndClientForTest(t)

		// create account members
		currentStatus, statusErr := accountAdminUserClient.GetAuthStatus(ctx, &authgrpc.GetAuthStatusRequest{})
		require.NotNil(t, currentStatus)
		require.NoError(t, statusErr)
		relevantAccountID := currentStatus.ActiveAccount

		createdClients := []client.Client{}
		for range 2 {
			u, c := createUserAndClientForTest(t)

			invitation, err := accountAdminUserClient.CreateAccountInvitation(ctx, &identitygrpc.CreateAccountInvitationRequest{
				Input: &identitygrpc.AccountInvitationCreationRequestInput{
					Note:    t.Name(),
					ToName:  t.Name(),
					ToEmail: u.EmailAddress,
				},
			})
			require.NoError(t, err)

			sentInvitations, err := accountAdminUserClient.GetSentAccountInvitations(ctx, &identitygrpc.GetSentAccountInvitationsRequest{})
			require.NotNil(t, sentInvitations)
			require.NoError(t, err)
			assert.NotEmpty(t, sentInvitations.Results)

			invitations, err := c.GetReceivedAccountInvitations(ctx, &identitygrpc.GetReceivedAccountInvitationsRequest{})
			require.NotNil(t, invitations)
			require.NoError(t, err)
			assert.NotEmpty(t, invitations.Results)

			_, err = c.AcceptAccountInvitation(ctx, &identitygrpc.AcceptAccountInvitationRequest{
				AccountInvitationId: invitation.Created.Id,
				Input: &identitygrpc.AccountInvitationUpdateRequestInput{
					Token: invitation.Created.Token,
					Note:  t.Name(),
				},
			})

			require.NoError(t, err)
			_, err = c.SetDefaultAccount(ctx, &identitygrpc.SetDefaultAccountRequest{AccountId: relevantAccountID})
			require.NoError(t, err)

			tokenResponse, err := c.LoginForToken(ctx, &authgrpc.LoginForTokenRequest{Input: &authgrpc.UserLoginInput{
				Username:  u.Username,
				Password:  u.HashedPassword,
				TotpToken: generateTOTPCodeForUserForTest(t, u),
			}})
			require.NoError(t, err)
			assert.NotNil(t, tokenResponse)

			createdClients = append(createdClients, c)
		}

		// create recipes for meal plan
		createdMeals := []*mealplanning.Meal{}
		for range 3 {
			createdMeal := createMealForTest(t, accountAdminUserClient, nil)
			createdMeals = append(createdMeals, createdMeal)
		}

		const baseDeadline = 10 * time.Second
		now := time.Now()

		exampleMealPlan := &mealplanning.MealPlan{
			Notes:          t.Name(),
			Status:         string(mealplanning.MealPlanStatusAwaitingVotes),
			VotingDeadline: now.Add(baseDeadline),
			ElectionMethod: mealplanning.MealPlanElectionMethodSchulze,
			Events: []*mealplanning.MealPlanEvent{
				{
					StartsAt: now.Add(24 * time.Hour),
					EndsAt:   now.Add(72 * time.Hour),
					MealName: mealplanning.BreakfastMealName,
					Options: []*mealplanning.MealPlanOption{
						{
							Meal:  mealplanning.Meal{ID: createdMeals[0].ID},
							Notes: "option A",
						},
						{
							Meal:  mealplanning.Meal{ID: createdMeals[1].ID},
							Notes: "option B",
						},
						{
							Meal:  mealplanning.Meal{ID: createdMeals[2].ID},
							Notes: "option C",
						},
					},
				},
			},
		}

		createdMealPlan := createMealPlanForTest(t, accountAdminUserClient, exampleMealPlan)
		createdMealPlanEvent := createdMealPlan.Events[0]
		require.NotNil(t, createdMealPlanEvent)

		userAVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
			Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
					Rank:                    2,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
					Rank:                    1,
				},
			},
		}

		userBVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
			Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
					Rank:                    1,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
					Rank:                    2,
				},
			},
		}

		userCVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
			Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
					Rank:                    1,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
					Rank:                    2,
				},
			},
		}

		createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userAVotes),
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanOptionVotesA)

		createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userBVotes),
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanOptionVotesB)

		createdMealPlanOptionVotesC, err := accountAdminUserClient.CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userCVotes),
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanOptionVotesC)

		createdMealPlan.VotingDeadline = time.Now().Add(-time.Minute)

		q := generated.New()
		rowsAffected, err := q.UpdateMealPlan(ctx, databaseClient.DB(), &generated.UpdateMealPlanParams{
			Notes:            createdMealPlan.Notes,
			Status:           generated.MealPlanStatus(createdMealPlan.Status),
			VotingDeadline:   time.Now().Add(-time.Minute),
			BelongsToAccount: relevantAccountID,
			ID:               createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotZero(t, rowsAffected)

		runRes, err := adminClient.RunFinalizeMealPlanWorker(ctx, &mealplanninggrpc.RunFinalizeMealPlanWorkerRequest{})
		require.NoError(t, err)
		require.NotNil(t, runRes)

		createdMealPlanRes, err := accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanId: createdMealPlan.ID})
		require.NotNil(t, createdMealPlanRes)
		require.NoError(t, err)

		actual := converters.ConvertGRPCMealPlanToMealPlan(createdMealPlanRes.Result)
		assert.Equal(t, string(mealplanning.MealPlanStatusFinalized), actual.Status)

		for _, event := range actual.Events {
			selectionMade := false
			for _, opt := range event.Options {
				if opt.Chosen {
					selectionMade = true
					break
				}
			}
			require.True(t, selectionMade)
		}
	})
}

func TestMealPlans_CompleteLifecycleForSomeVotesReceived(T *testing.T) {
	T.Parallel()

	T.Run("should resolve the meal plan status upon receiving some votes", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// create a userClient for the meal plan account
		_, accountAdminUserClient := createUserAndClientForTest(t)

		// create account members
		currentStatus, statusErr := accountAdminUserClient.GetAuthStatus(ctx, &authgrpc.GetAuthStatusRequest{})
		require.NotNil(t, currentStatus)
		require.NoError(t, statusErr)
		relevantAccountID := currentStatus.ActiveAccount

		createdClients := []client.Client{}
		for range 2 {
			u, c := createUserAndClientForTest(t)

			invitation, err := accountAdminUserClient.CreateAccountInvitation(ctx, &identitygrpc.CreateAccountInvitationRequest{
				Input: &identitygrpc.AccountInvitationCreationRequestInput{
					Note:    t.Name(),
					ToName:  t.Name(),
					ToEmail: u.EmailAddress,
				},
			})
			require.NoError(t, err)

			sentInvitations, err := accountAdminUserClient.GetSentAccountInvitations(ctx, &identitygrpc.GetSentAccountInvitationsRequest{})
			require.NotNil(t, sentInvitations)
			require.NoError(t, err)
			assert.NotEmpty(t, sentInvitations.Results)

			invitations, err := c.GetReceivedAccountInvitations(ctx, &identitygrpc.GetReceivedAccountInvitationsRequest{})
			require.NotNil(t, invitations)
			require.NoError(t, err)
			assert.NotEmpty(t, invitations.Results)

			_, err = c.AcceptAccountInvitation(ctx, &identitygrpc.AcceptAccountInvitationRequest{
				AccountInvitationId: invitation.Created.Id,
				Input: &identitygrpc.AccountInvitationUpdateRequestInput{
					Token: invitation.Created.Token,
					Note:  t.Name(),
				},
			})

			require.NoError(t, err)
			_, err = c.SetDefaultAccount(ctx, &identitygrpc.SetDefaultAccountRequest{AccountId: relevantAccountID})
			require.NoError(t, err)

			tokenResponse, err := c.LoginForToken(ctx, &authgrpc.LoginForTokenRequest{Input: &authgrpc.UserLoginInput{
				Username:  u.Username,
				Password:  u.HashedPassword,
				TotpToken: generateTOTPCodeForUserForTest(t, u),
			}})
			require.NoError(t, err)
			assert.NotNil(t, tokenResponse)

			createdClients = append(createdClients, c)
		}

		// create recipes for meal plan
		createdMeals := []*mealplanning.Meal{}
		for range 3 {
			createdMeal := createMealForTest(t, accountAdminUserClient, nil)
			createdMeals = append(createdMeals, createdMeal)
		}

		const baseDeadline = time.Minute
		now := time.Now()

		exampleMealPlan := &mealplanning.MealPlan{
			Notes:          t.Name(),
			Status:         string(mealplanning.MealPlanStatusAwaitingVotes),
			VotingDeadline: now.Add(baseDeadline),
			ElectionMethod: mealplanning.MealPlanElectionMethodSchulze,
			Events: []*mealplanning.MealPlanEvent{
				{
					StartsAt: now.Add(24 * time.Hour),
					EndsAt:   now.Add(72 * time.Hour),
					MealName: mealplanning.BreakfastMealName,
					Options: []*mealplanning.MealPlanOption{
						{
							Meal:  mealplanning.Meal{ID: createdMeals[0].ID},
							Notes: "option A",
						},
						{
							Meal:  mealplanning.Meal{ID: createdMeals[1].ID},
							Notes: "option B",
						},
						{
							Meal:  mealplanning.Meal{ID: createdMeals[2].ID},
							Notes: "option C",
						},
					},
				},
			},
		}

		exampleMealPlanInput := mpconverters.ConvertMealPlanToMealPlanCreationRequestInput(exampleMealPlan)
		createMealPlanRes, err := accountAdminUserClient.CreateMealPlan(ctx, &mealplanninggrpc.CreateMealPlanRequest{Input: converters.ConvertMealPlanCreationRequestInputToGRPCMealPlanCreationRequestInput(exampleMealPlanInput)})
		require.NotEmpty(t, createMealPlanRes.Created.Id)
		require.NoError(t, err)

		createdMealPlanRes, err := accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanId: createMealPlanRes.Created.Id})
		require.NotNil(t, createdMealPlanRes)
		require.NoError(t, err)

		createdMealPlan := converters.ConvertGRPCMealPlanToMealPlan(createdMealPlanRes.Result)

		checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

		createdMealPlanEvent := createdMealPlan.Events[0]

		userAVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
			Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
					Rank:                    2,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
					Rank:                    1,
				},
			},
		}

		userBVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
			Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
					Rank:                    0,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
					Rank:                    1,
				},
				{
					BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
					Rank:                    2,
				},
			},
		}

		createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userAVotes),
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanOptionVotesA)

		createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
			MealPlanId:      createdMealPlan.ID,
			MealPlanEventId: createdMealPlanEvent.ID,
			Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userBVotes),
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanOptionVotesB)

		createdMealPlanRes, err = accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanId: createdMealPlan.ID})
		require.NotNil(t, createdMealPlanRes)
		require.NoError(t, err)

		createdMealPlan = converters.ConvertGRPCMealPlanToMealPlan(createdMealPlanRes.Result)
		assert.Equal(t, string(mealplanning.MealPlanStatusAwaitingVotes), createdMealPlan.Status)

		q := generated.New()
		rowsAffected, err := q.UpdateMealPlan(ctx, databaseClient.DB(), &generated.UpdateMealPlanParams{
			Notes:            createdMealPlan.Notes,
			Status:           generated.MealPlanStatus(mealplanning.MealPlanStatusAwaitingVotes),
			VotingDeadline:   time.Now().Add(-10 * time.Hour),
			BelongsToAccount: relevantAccountID,
			ID:               createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotZero(t, rowsAffected)

		time.Sleep(10 * time.Second)

		runRes, err := adminClient.RunFinalizeMealPlanWorker(ctx, &mealplanninggrpc.RunFinalizeMealPlanWorkerRequest{})
		require.NoError(t, err)
		require.NotNil(t, runRes)

		createdMealPlanRes, err = accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{MealPlanId: createdMealPlan.ID})
		require.NotNil(t, createdMealPlan)
		require.NoError(t, err)

		createdMealPlan = converters.ConvertGRPCMealPlanToMealPlan(createdMealPlanRes.Result)

		assert.Equal(t, string(mealplanning.MealPlanStatusFinalized), createdMealPlan.Status)

		for _, event := range createdMealPlan.Events {
			selectionMade := false
			for _, opt := range event.Options {
				if opt.Chosen {
					selectionMade = true
					break
				}
			}
			require.True(t, selectionMade)
		}
	})
}
