package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVote(input *messages.MealPlanOptionVoteCreationRequestInput) *messages.MealPlanOptionVote {

output := &messages.MealPlanOptionVote{
}

return output
}
func ConvertMealPlanOptionVoteUpdateRequestInputToMealPlanOptionVote(input *messages.MealPlanOptionVoteUpdateRequestInput) *messages.MealPlanOptionVote {

output := &messages.MealPlanOptionVote{
    BelongsToMealPlanOption: input.BelongsToMealPlanOption,
    Rank: input.Rank,
    Abstain: input.Abstain,
    Notes: input.Notes,
}

return output
}
