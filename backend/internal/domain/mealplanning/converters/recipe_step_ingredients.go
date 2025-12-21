package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput creates a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredientCreationRequestInput.
func ConvertRecipeStepIngredientCreationRequestInputToRecipeStepIngredientDatabaseCreationInput(input *mealplanning.RecipeStepIngredientCreationRequestInput) *mealplanning.RecipeStepIngredientDatabaseCreationInput {
	x := &mealplanning.RecipeStepIngredientDatabaseCreationInput{
		ID:                               identifiers.New(),
		IngredientID:                     input.IngredientID,
		ValidIngredientPreparationID:     input.ValidIngredientPreparationID,
		ValidIngredientMeasurementUnitID: input.ValidIngredientMeasurementUnitID,
		Name:                             input.Name,
		MeasurementUnitID:                input.MeasurementUnitID,
		Quantity: types.Float32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		QuantityNotes:                   input.QuantityNotes,
		IngredientNotes:                 input.IngredientNotes,
		Optional:                        input.Optional,
		OptionIndex:                     input.OptionIndex,
		ProductOfRecipeStepIndex:        input.ProductOfRecipeStepIndex,
		ProductOfRecipeStepProductIndex: input.ProductOfRecipeStepProductIndex,
		VesselIndex:                     input.VesselIndex,
		ToTaste:                         input.ToTaste,
		ProductPercentageToUse:          input.ProductPercentageToUse,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput creates a RecipeStepIngredientUpdateRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientUpdateRequestInput(input *mealplanning.RecipeStepIngredient) *mealplanning.RecipeStepIngredientUpdateRequestInput {
	x := &mealplanning.RecipeStepIngredientUpdateRequestInput{
		IngredientID:        &input.Ingredient.ID,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		MeasurementUnitID:   &input.MeasurementUnit.ID,
		QuantityNotes:       &input.QuantityNotes,
		IngredientNotes:     &input.IngredientNotes,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		Quantity: types.Float32RangeWithOptionalMaxUpdateRequestInput{
			Max: input.Quantity.Max,
			Min: &input.Quantity.Min,
		},
		Optional:               &input.Optional,
		OptionIndex:            &input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                &input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}

	return x
}

// ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput builds a RecipeStepIngredientCreationRequestInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientCreationRequestInput(input *mealplanning.RecipeStepIngredient) *mealplanning.RecipeStepIngredientCreationRequestInput {
	var ingredientID *string
	if input.Ingredient != nil {
		ingredientID = &input.Ingredient.ID
	}

	return &mealplanning.RecipeStepIngredientCreationRequestInput{
		Name:              input.Name,
		Optional:          input.Optional,
		IngredientID:      ingredientID,
		MeasurementUnitID: input.MeasurementUnit.ID,
		Quantity: types.Float32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		QuantityNotes:          input.QuantityNotes,
		IngredientNotes:        input.IngredientNotes,
		OptionIndex:            input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}
}

// ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput builds a RecipeStepIngredientDatabaseCreationInput from a RecipeStepIngredient.
func ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(input *mealplanning.RecipeStepIngredient) *mealplanning.RecipeStepIngredientDatabaseCreationInput {
	return &mealplanning.RecipeStepIngredientDatabaseCreationInput{
		ID:                input.ID,
		Name:              input.Name,
		Optional:          input.Optional,
		IngredientID:      &input.Ingredient.ID,
		MeasurementUnitID: input.MeasurementUnit.ID,
		Quantity: types.Float32RangeWithOptionalMax{
			Max: input.Quantity.Max,
			Min: input.Quantity.Min,
		},
		QuantityNotes:          input.QuantityNotes,
		IngredientNotes:        input.IngredientNotes,
		BelongsToRecipeStep:    input.BelongsToRecipeStep,
		OptionIndex:            input.OptionIndex,
		VesselIndex:            input.VesselIndex,
		ToTaste:                input.ToTaste,
		ProductPercentageToUse: input.ProductPercentageToUse,
	}
}
