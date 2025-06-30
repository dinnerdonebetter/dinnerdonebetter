package converters

import types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"

// ConvertMealPlanOptionVoteToMealPlanOptionVoteUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertMealPlanOptionVoteToMealPlanOptionVoteUpdateRequestInput(input *types.MealPlanOptionVote) *types.MealPlanOptionVoteUpdateRequestInput {
	x := &types.MealPlanOptionVoteUpdateRequestInput{
		Notes:                   &input.Notes,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    &input.Rank,
		Abstain:                 &input.Abstain,
	}

	return x
}

// ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteDatabaseCreationInput(input *types.MealPlanOptionVoteCreationRequestInput) *types.MealPlanOptionVotesDatabaseCreationInput {
	x := &types.MealPlanOptionVotesDatabaseCreationInput{
		Votes: input.Votes,
	}

	return x
}
