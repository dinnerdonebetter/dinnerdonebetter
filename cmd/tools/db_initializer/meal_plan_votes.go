package main

import (
	"context"
	"fmt"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/segmentio/ksuid"
)

func scaffoldMealPlanVotes(ctx context.Context, db database.DataManager) error {
	// we have to do this here so that it runs after we create teh options things below
	var mealPlanVotesCollection = struct {
		MomJonesPastHouseholdMealPlanVotes,
		DadJonesPastHouseholdMealPlanVotes,
		KidJones1PastHouseholdMealPlanVotes,
		KidJones2PastHouseholdMealPlanVotes,
		MomJonesCurrentHouseholdMealPlanVotes,
		DadJonesCurrentHouseholdMealPlanVotes,
		KidJones1CurrentHouseholdMealPlanVotes *types.MealPlanOptionVoteDatabaseCreationInput
	}{
		MomJonesPastHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.MomJones.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// MomJones's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[0].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[2].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[3].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[4].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[5].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[6].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[7].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[8].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[9].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[10].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[11].ID,
					Rank:                    3,
				},
			},
		},
		DadJonesPastHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.DadJones.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// DadJones's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[0].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[2].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[3].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[4].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[5].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[6].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[7].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[8].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[9].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[10].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[11].ID,
					Rank:                    1,
				},
			},
		},
		KidJones1PastHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.KidJones1.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// KidJones1's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[0].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[1].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[2].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[3].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[4].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[5].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[6].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[7].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[8].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[9].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[10].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[11].ID,
					Rank:                    1,
				},
			},
		},
		KidJones2PastHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.KidJones2.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// KidJones2's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[0].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[2].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[3].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[4].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[5].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[6].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[7].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[8].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[9].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[10].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones2.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesPastHouseholdMealPlan.Options[11].ID,
					Rank:                    1,
				},
			},
		},
		MomJonesCurrentHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.MomJones.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// MomJones's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[0].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[2].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[3].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[4].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[5].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[6].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[7].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[8].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[9].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[10].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.MomJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[11].ID,
					Rank:                    3,
				},
			},
		},
		DadJonesCurrentHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.DadJones.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// DadJones's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[0].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[1].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[2].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[3].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[4].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[5].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[6].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[7].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[8].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[9].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[10].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.DadJones.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[11].ID,
					Rank:                    1,
				},
			},
		},
		KidJones1CurrentHouseholdMealPlanVotes: &types.MealPlanOptionVoteDatabaseCreationInput{
			ByUser: userCollection.KidJones1.ID,
			Votes: []*types.MealPlanOptionVoteCreationInput{
				// KidJones1's votes
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[0].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[1].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[2].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[3].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[4].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[5].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[6].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[7].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[8].ID,
					Rank:                    1,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[9].ID,
					Rank:                    2,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[10].ID,
					Rank:                    3,
				},
				{
					ID:                      ksuid.New().String(),
					ByUser:                  userCollection.KidJones1.ID,
					BelongsToMealPlanOption: mealPlanCollection.JonesCurrentHouseholdMealPlan.Options[11].ID,
					Rank:                    1,
				},
			},
		},
	}

	mealPlanVotes := []*types.MealPlanOptionVoteDatabaseCreationInput{
		mealPlanVotesCollection.MomJonesPastHouseholdMealPlanVotes,
		mealPlanVotesCollection.DadJonesPastHouseholdMealPlanVotes,
		mealPlanVotesCollection.KidJones1PastHouseholdMealPlanVotes,
		mealPlanVotesCollection.KidJones2PastHouseholdMealPlanVotes,
		mealPlanVotesCollection.MomJonesCurrentHouseholdMealPlanVotes,
		mealPlanVotesCollection.DadJonesCurrentHouseholdMealPlanVotes,
		mealPlanVotesCollection.KidJones1CurrentHouseholdMealPlanVotes,
	}

	for _, input := range mealPlanVotes {
		if _, err := db.CreateMealPlanOptionVote(ctx, input); err != nil {
			return fmt.Errorf("voting for meal plan: %w", err)
		}
	}

	if _, finalizationErr := db.AttemptToFinalizeCompleteMealPlan(ctx, mealPlanCollection.JonesPastHouseholdMealPlan.ID, jonesHouseholdID); finalizationErr != nil {
		return fmt.Errorf("finalizing meal plan: %w", finalizationErr)
	}

	if _, finalizationErr := db.AttemptToFinalizeCompleteMealPlan(ctx, mealPlanCollection.JonesCurrentHouseholdMealPlan.ID, jonesHouseholdID); finalizationErr != nil {
		return fmt.Errorf("finalizing meal plan: %w", finalizationErr)
	}

	return nil
}
