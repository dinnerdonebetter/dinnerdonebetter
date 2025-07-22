package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

// ConvertRecipeStepToRecipeStepUpdateRequestInput creates a RecipeStepUpdateRequestInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepUpdateRequestInput(input *mealplanning.RecipeStep) *mealplanning.RecipeStepUpdateRequestInput {
	x := &mealplanning.RecipeStepUpdateRequestInput{
		Notes:                   &input.Notes,
		BelongsToRecipe:         &input.BelongsToRecipe,
		Preparation:             &input.Preparation,
		Index:                   &input.Index,
		EstimatedTimeInSeconds:  input.EstimatedTimeInSeconds,
		TemperatureInCelsius:    input.TemperatureInCelsius,
		Optional:                &input.Optional,
		ExplicitInstructions:    &input.ExplicitInstructions,
		ConditionExpression:     &input.ConditionExpression,
		StartTimerAutomatically: &input.StartTimerAutomatically,
	}

	return x
}

// ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput creates a RecipeStepDatabaseCreationInput from a RecipeStepCreationRequestInput.
func ConvertRecipeStepCreationInputToRecipeStepDatabaseCreationInput(input *mealplanning.RecipeStepCreationRequestInput) *mealplanning.RecipeStepDatabaseCreationInput {
	x := &mealplanning.RecipeStepDatabaseCreationInput{
		ID:                      identifiers.New(),
		Index:                   input.Index,
		PreparationID:           input.PreparationID,
		EstimatedTimeInSeconds:  input.EstimatedTimeInSeconds,
		TemperatureInCelsius:    input.TemperatureInCelsius,
		Notes:                   input.Notes,
		Optional:                input.Optional,
		ExplicitInstructions:    input.ExplicitInstructions,
		ConditionExpression:     input.ConditionExpression,
		StartTimerAutomatically: input.StartTimerAutomatically,
	}

	x.Ingredients = []*mealplanning.RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		convertedIngredient := ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(ingredient)
		convertedIngredient.ID = identifiers.New()
		convertedIngredient.BelongsToRecipeStep = x.ID
		x.Ingredients = append(x.Ingredients, convertedIngredient)
	}

	x.Instruments = []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		convertedInstrument := ConvertRecipeStepInstrumentCreationRequestInputToRecipeStepInstrumentDatabaseCreationInput(instrument)
		convertedInstrument.ID = identifiers.New()
		convertedInstrument.BelongsToRecipeStep = x.ID
		x.Instruments = append(x.Instruments, convertedInstrument)
	}

	x.Vessels = []*mealplanning.RecipeStepVesselDatabaseCreationInput{}
	for _, vessel := range input.Vessels {
		convertedVessel := ConvertRecipeStepVesselCreationRequestInputToRecipeStepVesselDatabaseCreationInput(vessel)
		convertedVessel.ID = identifiers.New()
		convertedVessel.BelongsToRecipeStep = x.ID
		x.Vessels = append(x.Vessels, convertedVessel)
	}

	x.Products = []*mealplanning.RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		convertedProduct := ConvertRecipeStepProductCreationInputToRecipeStepProductDatabaseCreationInput(product)
		convertedProduct.ID = identifiers.New()
		convertedProduct.BelongsToRecipeStep = x.ID
		x.Products = append(x.Products, convertedProduct)
	}

	x.CompletionConditions = []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{}
	for _, product := range input.CompletionConditions {
		convertedCompletionCondition := ConvertRecipeStepCompletionConditionCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput(x, product)
		convertedCompletionCondition.ID = identifiers.New()
		convertedCompletionCondition.BelongsToRecipeStep = x.ID
		x.CompletionConditions = append(x.CompletionConditions, convertedCompletionCondition)
	}

	return x
}

// ConvertRecipeStepToRecipeStepCreationRequestInput builds a RecipeStepCreationRequestInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepCreationRequestInput(recipeStep *mealplanning.RecipeStep) *mealplanning.RecipeStepCreationRequestInput {
	ingredients := []*mealplanning.RecipeStepIngredientCreationRequestInput{}
	for _, ingredient := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(ingredient))
	}

	instruments := []*mealplanning.RecipeStepInstrumentCreationRequestInput{}
	for _, instrument := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentCreationRequestInput(instrument))
	}

	vessels := []*mealplanning.RecipeStepVesselCreationRequestInput{}
	for _, vessel := range recipeStep.Vessels {
		vessels = append(vessels, ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(vessel))
	}

	products := []*mealplanning.RecipeStepProductCreationRequestInput{}
	for _, product := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(product))
	}

	completionConditions := []*mealplanning.RecipeStepCompletionConditionCreationRequestInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionCreationRequestInput(recipeStep, completionCondition))
	}

	return &mealplanning.RecipeStepCreationRequestInput{
		Optional:                recipeStep.Optional,
		Index:                   recipeStep.Index,
		PreparationID:           recipeStep.Preparation.ID,
		EstimatedTimeInSeconds:  recipeStep.EstimatedTimeInSeconds,
		TemperatureInCelsius:    recipeStep.TemperatureInCelsius,
		Notes:                   recipeStep.Notes,
		ExplicitInstructions:    recipeStep.ExplicitInstructions,
		ConditionExpression:     recipeStep.ConditionExpression,
		StartTimerAutomatically: recipeStep.StartTimerAutomatically,
		Products:                products,
		Ingredients:             ingredients,
		Instruments:             instruments,
		Vessels:                 vessels,
		CompletionConditions:    completionConditions,
	}
}

// ConvertRecipeStepToRecipeStepDatabaseCreationInput builds a RecipeStepDatabaseCreationInput from a RecipeStep.
func ConvertRecipeStepToRecipeStepDatabaseCreationInput(recipeStep *mealplanning.RecipeStep) *mealplanning.RecipeStepDatabaseCreationInput {
	ingredients := []*mealplanning.RecipeStepIngredientDatabaseCreationInput{}
	for _, i := range recipeStep.Ingredients {
		ingredients = append(ingredients, ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(i))
	}

	instruments := []*mealplanning.RecipeStepInstrumentDatabaseCreationInput{}
	for _, i := range recipeStep.Instruments {
		instruments = append(instruments, ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(i))
	}

	vessels := []*mealplanning.RecipeStepVesselDatabaseCreationInput{}
	for _, v := range recipeStep.Vessels {
		vessels = append(vessels, ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(v))
	}

	products := []*mealplanning.RecipeStepProductDatabaseCreationInput{}
	for _, p := range recipeStep.Products {
		products = append(products, ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(p))
	}

	completionConditions := []*mealplanning.RecipeStepCompletionConditionDatabaseCreationInput{}
	for _, completionCondition := range recipeStep.CompletionConditions {
		completionConditions = append(completionConditions, ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(completionCondition))
	}

	return &mealplanning.RecipeStepDatabaseCreationInput{
		ID:                      recipeStep.ID,
		Index:                   recipeStep.Index,
		PreparationID:           recipeStep.Preparation.ID,
		Optional:                recipeStep.Optional,
		EstimatedTimeInSeconds:  recipeStep.EstimatedTimeInSeconds,
		TemperatureInCelsius:    recipeStep.TemperatureInCelsius,
		StartTimerAutomatically: recipeStep.StartTimerAutomatically,
		Notes:                   recipeStep.Notes,
		ExplicitInstructions:    recipeStep.ExplicitInstructions,
		ConditionExpression:     recipeStep.ConditionExpression,
		Ingredients:             ingredients,
		Instruments:             instruments,
		Products:                products,
		Vessels:                 vessels,
		BelongsToRecipe:         recipeStep.BelongsToRecipe,
		CompletionConditions:    completionConditions,
	}
}
