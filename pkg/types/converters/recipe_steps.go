package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
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
		ConditionExpression:           &input.ConditionExpression,
	}

	return x
}

// ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput creates a RecipeStepDatabaseCreationInput from a RecipeStepCreationRequestInput.
func ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input *types.RecipeStepCreationRequestInput) *types.RecipeStepDatabaseCreationInput {
	x := &types.RecipeStepDatabaseCreationInput{
		ID:                            identifiers.New(),
		Index:                         input.Index,
		PreparationID:                 input.PreparationID,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         input.Notes,
		Optional:                      input.Optional,
		ExplicitInstructions:          input.ExplicitInstructions,
		ConditionExpression:           input.ConditionExpression,
	}

	x.Ingredients = []*types.RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		convertedIngredient := ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(ingredient)
		convertedIngredient.ID = identifiers.New()
		convertedIngredient.BelongsToRecipeStep = x.ID
		x.Ingredients = append(x.Ingredients, convertedIngredient)
	}

	x.Instruments = []*types.RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		convertedInstrument := ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(instrument)
		convertedInstrument.ID = identifiers.New()
		convertedInstrument.BelongsToRecipeStep = x.ID
		x.Instruments = append(x.Instruments, convertedInstrument)
	}

	x.Products = []*types.RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		convertedProduct := ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(product)
		convertedProduct.ID = identifiers.New()
		convertedProduct.BelongsToRecipeStep = x.ID
		x.Products = append(x.Products, convertedProduct)
	}

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

	completionConditions := []*types.RecipeStepCompletionConditionCreationRequestInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(completionCondition))
	}

	return &types.RecipeStepCreationRequestInput{
		Optional:                      recipeStep.Optional,
		Index:                         recipeStep.Index,
		PreparationID:                 recipeStep.Preparation.ID,
		MinimumEstimatedTimeInSeconds: recipeStep.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: recipeStep.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   recipeStep.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   recipeStep.MaximumTemperatureInCelsius,
		Notes:                         recipeStep.Notes,
		ExplicitInstructions:          recipeStep.ExplicitInstructions,
		ConditionExpression:           recipeStep.ConditionExpression,
		Products:                      products,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		CompletionConditions:          completionConditions,
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

	completionConditions := []*types.RecipeStepCompletionConditionDatabaseCreationInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(completionCondition))
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
		ConditionExpression:           recipeStep.ConditionExpression,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		Products:                      products,
		BelongsToRecipe:               recipeStep.BelongsToRecipe,
		CompletionConditions:          completionConditions,
	}
}
