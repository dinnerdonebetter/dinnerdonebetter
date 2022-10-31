package converters

import (
	"github.com/prixfixeco/api_server/internal/identifiers"
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertRecipeStepToRecipeStepUpdateRequestInput creates a RecipeStepUpdateRequestInput from a RecipeStep.
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

// ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput creates a RecipeStepDatabaseCreationInput from a RecipeStepCreationRequestInput.
func ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input *types.RecipeStepCreationRequestInput) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(ingredient))
	}

	instruments := []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(instrument))
	}

	products := []*types.RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		products = append(products, ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(product))
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
	x.ID = identifiers.New()

	return x
}

// ConvertRecipeStepToRecipeStepCreationRequestInput builds a RecipeStepCreationRequestInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep *types.RecipeStep) *types.RecipeStepCreationRequestInput {
	ingredients := []*types.RecipeStepIngredientCreationRequestInput{}
	for _, ingredient := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(ingredient))
	}

	instruments := []*types.RecipeStepInstrumentCreationRequestInput{}
	for _, instrument := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(instrument))
	}

	products := []*types.RecipeStepProductCreationRequestInput{}
	for _, product := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(product))
	}

	return &types.RecipeStepCreationRequestInput{
		ID:                            recipeStep.ID,
		Optional:                      recipeStep.Optional,
		Index:                         recipeStep.Index,
		PreparationID:                 recipeStep.Preparation.ID,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		Notes:                         recipeStep.Notes,
		ExplicitInstructions:          recipeStep.ExplicitInstructions,
		BelongsToRecipe:               recipeStep.BelongsToRecipe,
		Products:                      products,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
	}
}

// ConvertRecipeStepToRecipeStepDatabaseCreationInput builds a RecipeStepDatabaseCreationInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepDatabaseCreationInput(recipeStep *types.RecipeStep) *types.RecipeStepDatabaseCreationInput {
	ingredients := []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(i))
	}

	instruments := []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, i := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(i))
	}

	products := []*types.RecipeStepProductDatabaseCreationInput{}
	for _, p := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(p))
	}

	return &types.RecipeStepDatabaseCreationInput{
		ID:                            recipeStep.ID,
		Index:                         recipeStep.Index,
		PreparationID:                 recipeStep.Preparation.ID,
		Optional:                      recipeStep.Optional,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		Notes:                         recipeStep.Notes,
		ExplicitInstructions:          recipeStep.ExplicitInstructions,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		Products:                      products,
		BelongsToRecipe:               recipeStep.BelongsToRecipe,
	}
}
