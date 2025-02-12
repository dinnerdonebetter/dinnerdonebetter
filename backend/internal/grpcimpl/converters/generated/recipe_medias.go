package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeMediaCreationRequestInputToRecipeMedia(input *messages.RecipeMediaCreationRequestInput) *messages.RecipeMedia {

output := &messages.RecipeMedia{
    BelongsToRecipe: input.BelongsToRecipe,
    MimeType: input.MimeType,
    InternalPath: input.InternalPath,
    ExternalPath: input.ExternalPath,
    Index: input.Index,
    BelongsToRecipeStep: input.BelongsToRecipeStep,
}

return output
}
func ConvertRecipeMediaUpdateRequestInputToRecipeMedia(input *messages.RecipeMediaUpdateRequestInput) *messages.RecipeMedia {

output := &messages.RecipeMedia{
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    BelongsToRecipe: input.BelongsToRecipe,
    MimeType: input.MimeType,
    InternalPath: input.InternalPath,
    ExternalPath: input.ExternalPath,
    Index: input.Index,
}

return output
}
