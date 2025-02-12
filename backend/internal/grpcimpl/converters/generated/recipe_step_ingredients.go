package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredient(input *messages.RecipeStepIngredientCreationRequestInput) *messages.RecipeStepIngredient {

output := &messages.RecipeStepIngredient{
    Quantity: input.Quantity,
    QuantityNotes: input.QuantityNotes,
    Name: input.Name,
    OptionIndex: input.OptionIndex,
    RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
    IngredientNotes: input.IngredientNotes,
    VesselIndex: input.VesselIndex,
    ProductPercentageToUse: input.ProductPercentageToUse,
    Optional: input.Optional,
    ToTaste: input.ToTaste,
}

return output
}
func ConvertRecipeStepIngredientUpdateRequestInputToRecipeStepIngredient(input *messages.RecipeStepIngredientUpdateRequestInput) *messages.RecipeStepIngredient {

output := &messages.RecipeStepIngredient{
    Quantity: ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input.Quantity),
    VesselIndex: input.VesselIndex,
    OptionIndex: input.OptionIndex,
    ProductPercentageToUse: input.ProductPercentageToUse,
    Optional: input.Optional,
    RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
    RecipeStepProductID: input.RecipeStepProductID,
    BelongsToRecipeStep: input.BelongsToRecipeStep,
    QuantityNotes: input.QuantityNotes,
    IngredientNotes: input.IngredientNotes,
    Name: input.Name,
    ToTaste: input.ToTaste,
}

return output
}
