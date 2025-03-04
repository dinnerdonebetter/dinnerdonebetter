package fakes

import (
	"math"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeMealPlanOptionVote builds a faked meal plan option vote.
func BuildFakeMealPlanOptionVote() *types.MealPlanOptionVote {
	return &types.MealPlanOptionVote{
		ID:                      BuildFakeID(),
		Rank:                    uint8(fake.Number(1, math.MaxUint8)),
		Abstain:                 fake.Bool(),
		Notes:                   buildUniqueString(),
		CreatedAt:               BuildFakeTime(),
		BelongsToMealPlanOption: fake.UUID(),
	}
}

// BuildFakeMealPlanOptionVotesList builds a faked MealPlanOptionVoteList.
func BuildFakeMealPlanOptionVotesList() *filtering.QueryFilteredResult[types.MealPlanOptionVote] {
	var examples []*types.MealPlanOptionVote
	for range exampleQuantity {
		examples = append(examples, BuildFakeMealPlanOptionVote())
	}

	return &filtering.QueryFilteredResult[types.MealPlanOptionVote]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
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
