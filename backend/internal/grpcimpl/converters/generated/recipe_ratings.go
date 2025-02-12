package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeRatingCreationRequestInputToRecipeRating(input *messages.RecipeRatingCreationRequestInput) *messages.RecipeRating {

output := &messages.RecipeRating{
    Difficulty: input.Difficulty,
    RecipeID: input.RecipeID,
    Notes: input.Notes,
    ByUser: input.ByUser,
    Taste: input.Taste,
    Instructions: input.Instructions,
    Overall: input.Overall,
    Cleanup: input.Cleanup,
}

return output
}
func ConvertRecipeRatingUpdateRequestInputToRecipeRating(input *messages.RecipeRatingUpdateRequestInput) *messages.RecipeRating {

output := &messages.RecipeRating{
    Taste: input.Taste,
    Instructions: input.Instructions,
    Overall: input.Overall,
    Cleanup: input.Cleanup,
    Difficulty: input.Difficulty,
    RecipeID: input.RecipeID,
    Notes: input.Notes,
    ByUser: input.ByUser,
}

return output
}
