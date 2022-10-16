package fakes

import (
	"math"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeMealPlanOptionVote builds a faked meal plan option vote.
func BuildFakeMealPlanOptionVote() *types.MealPlanOptionVote {
	return &types.MealPlanOptionVote{
		ID:                      BuildFakeID(),
		Rank:                    uint8(fake.Number(1, math.MaxUint8)),
		Abstain:                 fake.Bool(),
		Notes:                   buildUniqueString(),
		CreatedAt:               fake.Date(),
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
		Rank:                    &mealPlanOptionVote.Rank,
		Abstain:                 &mealPlanOptionVote.Abstain,
		Notes:                   &mealPlanOptionVote.Notes,
		BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
	}
}

// BuildFakeMealPlanOptionVoteCreationRequestInput builds a faked MealPlanOptionVoteCreationRequestInput.
func BuildFakeMealPlanOptionVoteCreationRequestInput() *types.MealPlanOptionVoteCreationRequestInput {
	mealPlanOptionVote := BuildFakeMealPlanOptionVote()
	return converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput(mealPlanOptionVote)
}

// BuildFakeMealPlanOptionVoteDatabaseCreationInput builds a faked MealPlanOptionVotesDatabaseCreationInput.
func BuildFakeMealPlanOptionVoteDatabaseCreationInput() *types.MealPlanOptionVotesDatabaseCreationInput {
	mealPlanOptionVote := BuildFakeMealPlanOptionVote()
	return converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput(mealPlanOptionVote)
}
