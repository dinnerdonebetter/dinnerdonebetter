package converters

import (
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertRecipeStepToRecipeStepUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeStepToRecipeStepUpdateRequestInput(input *types.RecipeStep) *types.RecipeStepUpdateRequestInput {
	x := &types.RecipeStepUpdateRequestInput{
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         &input.Notes,
		BelongsToRecipe:               input.BelongsToRecipe,
		Preparation:                   &input.Preparation,
		Index:                         &input.Index,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		Optional:                      &input.Optional,
		ExplicitInstructions:          &input.ExplicitInstructions,
	}

	return x
}

// ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input *types.RecipeStepCreationRequestInput) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(ingredient))
	}

	instruments := []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		instruments = append(instruments, types.RecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrumentCreationInput(instrument))
	}

	products := []*types.RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		products = append(products, types.RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(product))
	}

	x := &types.RecipeStepDatabaseCreationInput{
		Index:                         input.Index,
		PreparationID:                 input.PreparationID,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         input.Notes,
		Products:                      products,
		Optional:                      input.Optional,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		ExplicitInstructions:          input.ExplicitInstructions,
	}

	// we need to set this here or later converters will fail
	x.ID = ksuid.New().String()

	return x
}
