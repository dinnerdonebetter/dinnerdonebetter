package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMealPlanOptionVote builds a faked meal plan option vote.
func BuildFakeMealPlanOptionVote() *types.MealPlanOptionVote {
	return &types.MealPlanOptionVote{
		ID:                      ksuid.New().String(),
		Rank:                    fake.Uint8(),
		Abstain:                 fake.Bool(),
		Notes:                   fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:               uint64(uint32(fake.Date().Unix())),
		BelongsToMealPlanOption: fake.UUID(),
	}
}

// BuildFakeMealPlanOptionVoteList builds a faked MealPlanOptionVoteList.
func BuildFakeMealPlanOptionVoteList() *types.MealPlanOptionVoteList {
	var examples []*types.MealPlanOptionVote
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanOptionVote())
	}

	return &types.MealPlanOptionVoteList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlanOptionVotes: examples,
	}
}

// BuildFakeMealPlanOptionVoteUpdateRequestInput builds a faked MealPlanOptionVoteUpdateRequestInput from a meal plan option vote.
func BuildFakeMealPlanOptionVoteUpdateRequestInput() *types.MealPlanOptionVoteUpdateRequestInput {
	mealPlanOptionVote := BuildFakeMealPlanOptionVote()
	return &types.MealPlanOptionVoteUpdateRequestInput{
		Rank:                    mealPlanOptionVote.Rank,
		Abstain:                 mealPlanOptionVote.Abstain,
		Notes:                   mealPlanOptionVote.Notes,
		BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
	}
}

// BuildFakeMealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote builds a faked MealPlanOptionVoteUpdateRequestInput from a meal plan option vote.
func BuildFakeMealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote(mealPlanOptionVote *types.MealPlanOptionVote) *types.MealPlanOptionVoteUpdateRequestInput {
	return &types.MealPlanOptionVoteUpdateRequestInput{
		Rank:                    mealPlanOptionVote.Rank,
		Abstain:                 mealPlanOptionVote.Abstain,
		Notes:                   mealPlanOptionVote.Notes,
		BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
	}
}

// BuildFakeMealPlanOptionVoteCreationRequestInput builds a faked MealPlanOptionVoteCreationRequestInput.
func BuildFakeMealPlanOptionVoteCreationRequestInput() *types.MealPlanOptionVoteCreationRequestInput {
	mealPlanOptionVote := BuildFakeMealPlanOptionVote()
	return BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(mealPlanOptionVote)
}

// BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote builds a faked MealPlanOptionVoteCreationRequestInput from a meal plan option vote.
func BuildFakeMealPlanOptionVoteCreationRequestInputFromMealPlanOptionVote(mealPlanOptionVote *types.MealPlanOptionVote) *types.MealPlanOptionVoteCreationRequestInput {
	return &types.MealPlanOptionVoteCreationRequestInput{
		ID:                      mealPlanOptionVote.ID,
		Rank:                    mealPlanOptionVote.Rank,
		Abstain:                 mealPlanOptionVote.Abstain,
		Notes:                   mealPlanOptionVote.Notes,
		BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
	}
}

// BuildFakeMealPlanOptionVoteDatabaseCreationInput builds a faked MealPlanOptionVoteDatabaseCreationInput.
func BuildFakeMealPlanOptionVoteDatabaseCreationInput() *types.MealPlanOptionVoteDatabaseCreationInput {
	mealPlanOptionVote := BuildFakeMealPlanOptionVote()
	return BuildFakeMealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVote(mealPlanOptionVote)
}

// BuildFakeMealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVote builds a faked MealPlanOptionVoteDatabaseCreationInput from a meal plan option vote.
func BuildFakeMealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVote(mealPlanOptionVote *types.MealPlanOptionVote) *types.MealPlanOptionVoteDatabaseCreationInput {
	return &types.MealPlanOptionVoteDatabaseCreationInput{
		ID:                      mealPlanOptionVote.ID,
		Rank:                    mealPlanOptionVote.Rank,
		Abstain:                 mealPlanOptionVote.Abstain,
		Notes:                   mealPlanOptionVote.Notes,
		BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
	}
}
