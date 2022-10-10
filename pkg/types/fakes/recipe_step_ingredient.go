package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
// NOTE: this currently represents a typical recipe step ingredient with a valid ingredient and not a product.
func BuildFakeRecipeStepIngredient() *types.RecipeStepIngredient {
	return &types.RecipeStepIngredient{
		ID:                  BuildFakeID(),
		Name:                buildUniqueString(),
		Ingredient:          BuildFakeValidIngredient(),
		MeasurementUnit:     *BuildFakeValidMeasurementUnit(),
		MinimumQuantity:     float32(fake.Uint32()),
		MaximumQuantity:     float32(fake.Uint32()),
		QuantityNotes:       buildUniqueString(),
		ProductOfRecipeStep: false,
		Optional:            fake.Bool(),
		IngredientNotes:     buildUniqueString(),
		CreatedAt:           fake.Date(),
		BelongsToRecipeStep: BuildFakeID(),
	}
}

// BuildFakeRecipeStepIngredientList builds a faked RecipeStepIngredientList.
func BuildFakeRecipeStepIngredientList() *types.RecipeStepIngredientList {
	var examples []*types.RecipeStepIngredient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepIngredient())
	}

	return &types.RecipeStepIngredientList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeStepIngredients: examples,
	}
}

// BuildFakeRecipeStepIngredientUpdateRequestInput builds a faked RecipeStepIngredientUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateRequestInput() *types.RecipeStepIngredientUpdateRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return &types.RecipeStepIngredientUpdateRequestInput{
		Name:                &recipeStepIngredient.Name,
		Optional:            &recipeStepIngredient.Optional,
		IngredientID:        &recipeStepIngredient.Ingredient.ID,
		MeasurementUnitID:   &recipeStepIngredient.MeasurementUnit.ID,
		MinimumQuantity:     &recipeStepIngredient.MinimumQuantity,
		MaximumQuantity:     &recipeStepIngredient.MaximumQuantity,
		QuantityNotes:       &recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: &recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     &recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: &recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient builds a faked RecipeStepIngredientUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	return &types.RecipeStepIngredientUpdateRequestInput{
		Name:                &recipeStepIngredient.Name,
		Optional:            &recipeStepIngredient.Optional,
		IngredientID:        &recipeStepIngredient.Ingredient.ID,
		MeasurementUnitID:   &recipeStepIngredient.MeasurementUnit.ID,
		MinimumQuantity:     &recipeStepIngredient.MinimumQuantity,
		MaximumQuantity:     &recipeStepIngredient.MaximumQuantity,
		QuantityNotes:       &recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: &recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     &recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: &recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientCreationRequestInput builds a faked RecipeStepIngredientCreationRequestInput.
func BuildFakeRecipeStepIngredientCreationRequestInput() *types.RecipeStepIngredientCreationRequestInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(recipeStepIngredient)
}

// BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient builds a faked RecipeStepIngredientCreationRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientCreationRequestInput {
	return &types.RecipeStepIngredientCreationRequestInput{
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
	}
}

// BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient builds a faked RecipeStepIngredientDatabaseCreationInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientDatabaseCreationInput {
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
	}
}
