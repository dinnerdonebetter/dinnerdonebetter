package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput creates a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredientCreationRequestInput.
func ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(input *types.RecipeStepIngredientCreationRequestInput) *types.RecipeStepIngredientDatabaseCreationInput {
	x := &types.RecipeStepIngredientDatabaseCreationInput{
		ID:                  identifiers.New(),
		IngredientID:        input.IngredientID,
		Name:                input.Name,
		MeasurementUnitID:   input.MeasurementUnitID,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		RequiresDefrost:     input.RequiresDefrost,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput creates a RecipeStepIngredientUpdateRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(input *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	x := &types.RecipeStepIngredientUpdateRequestInput{
		IngredientID:        &input.Ingredient.ID,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		MeasurementUnitID:   &input.MeasurementUnit.ID,
		QuantityNotes:       &input.QuantityNotes,
		IngredientNotes:     &input.IngredientNotes,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		MinimumQuantity:     &input.MinimumQuantity,
		MaximumQuantity:     &input.MaximumQuantity,
		ProductOfRecipeStep: &input.ProductOfRecipeStep,
		Optional:            &input.Optional,
		OptionIndex:         &input.OptionIndex,
		RequiresDefrost:     &input.RequiresDefrost,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput builds a RecipeStepIngredientCreationRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientCreationRequestInput {
	return &types.RecipeStepIngredientCreationRequestInput{
		Name:                recipeStepIngredient.Name,
		Optional:            recipeStepIngredient.Optional,
		IngredientID:        &recipeStepIngredient.Ingredient.ID,
		MeasurementUnitID:   recipeStepIngredient.MeasurementUnit.ID,
		MinimumQuantity:     recipeStepIngredient.MinimumQuantity,
		MaximumQuantity:     recipeStepIngredient.MaximumQuantity,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		OptionIndex:         recipeStepIngredient.OptionIndex,
		RequiresDefrost:     recipeStepIngredient.RequiresDefrost,
	}
}

// ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput builds a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientDatabaseCreationInput {
	return &types.RecipeStepIngredientDatabaseCreationInput{
		ID:                  recipeStepIngredient.ID,
		Name:                recipeStepIngredient.Name,
		Optional:            recipeStepIngredient.Optional,
		IngredientID:        &recipeStepIngredient.Ingredient.ID,
		MeasurementUnitID:   recipeStepIngredient.MeasurementUnit.ID,
		MinimumQuantity:     recipeStepIngredient.MinimumQuantity,
		MaximumQuantity:     recipeStepIngredient.MaximumQuantity,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
		OptionIndex:         recipeStepIngredient.OptionIndex,
		RequiresDefrost:     recipeStepIngredient.RequiresDefrost,
	}
}
