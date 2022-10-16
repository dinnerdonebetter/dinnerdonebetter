package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput creates a MealPlanOptionUpdateRequestInput from a MealPlanOption.
func ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(input *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	x := &types.MealPlanOptionUpdateRequestInput{
		MealID:                 &input.Meal.ID,
		Notes:                  &input.Notes,
		BelongsToMealPlanEvent: &input.BelongsToMealPlanEvent,
		PrepStepsCreated:       &input.PrepStepsCreated,
	}

	return x
}

// ConvertMealPlanOptionCreationRequestInputToMealPlanOptionDatabaseCreationInput creates a MealPlanOptionDatabaseCreationInput from a MealPlanOptionCreationRequestInput.
func ConvertMealPlanOptionCreationRequestInputToMealPlanOptionDatabaseCreationInput(input *types.MealPlanOptionCreationRequestInput) *types.MealPlanOptionDatabaseCreationInput {
	x := &types.MealPlanOptionDatabaseCreationInput{
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
		MealID:                 input.MealID,
		Notes:                  input.Notes,
		PrepStepsCreated:       input.PrepStepsCreated,
	}

	return x
}

// ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput builds a MealPlanOptionVoteCreationRequestInput from a meal plan option vote.
func ConvertMealPlanOptionVoteToMealPlanOptionVoteCreationRequestInput(mealPlanOptionVote *types.MealPlanOptionVote) *types.MealPlanOptionVoteCreationRequestInput {
	return &types.MealPlanOptionVoteCreationRequestInput{
		Votes: []*types.MealPlanOptionVoteCreationInput{
			{
				ID:                      mealPlanOptionVote.ID,
				Rank:                    mealPlanOptionVote.Rank,
				Abstain:                 mealPlanOptionVote.Abstain,
				Notes:                   mealPlanOptionVote.Notes,
				BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
				ByUser:                  mealPlanOptionVote.ByUser,
			},
		},
		ByUser: mealPlanOptionVote.ByUser,
	}
}

// ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput builds a MealPlanOptionVotesDatabaseCreationInput from a meal plan option vote.
func ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput(mealPlanOptionVote *types.MealPlanOptionVote) *types.MealPlanOptionVotesDatabaseCreationInput {
	return &types.MealPlanOptionVotesDatabaseCreationInput{
		Votes: []*types.MealPlanOptionVoteCreationInput{
			{
				ID:                      mealPlanOptionVote.ID,
				Rank:                    mealPlanOptionVote.Rank,
				Abstain:                 mealPlanOptionVote.Abstain,
				Notes:                   mealPlanOptionVote.Notes,
				BelongsToMealPlanOption: mealPlanOptionVote.BelongsToMealPlanOption,
				ByUser:                  mealPlanOptionVote.ByUser,
			},
		},
		ByUser: mealPlanOptionVote.ByUser,
	}
}

// ConvertMealPlanOptionToMealPlanOptionCreationRequestInput builds a MealPlanOptionCreationRequestInput from a meal plan option.
func ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionCreationRequestInput {
	return &types.MealPlanOptionCreationRequestInput{
		ID:                     mealPlanOption.ID,
		MealID:                 mealPlanOption.Meal.ID,
		Notes:                  mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       mealPlanOption.PrepStepsCreated,
	}
}

// ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput builds a MealPlanOptionDatabaseCreationInput from a meal plan option.
func ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionDatabaseCreationInput {
	return &types.MealPlanOptionDatabaseCreationInput{
		ID:                     mealPlanOption.ID,
		MealID:                 mealPlanOption.Meal.ID,
		Notes:                  mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       mealPlanOption.PrepStepsCreated,
	}
}
