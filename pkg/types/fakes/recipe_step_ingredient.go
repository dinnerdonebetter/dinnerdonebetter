package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
func BuildFakeRecipeStepIngredient() *types.RecipeStepIngredient {
	return &types.RecipeStepIngredient{
		ID:                  uint64(fake.Uint32()),
		ExternalID:          fake.UUID(),
		IngredientID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		Name:                fake.Word(),
		QuantityType:        fake.Word(),
		QuantityValue:       fake.Float32(),
		QuantityNotes:       fake.Word(),
		ProductOfRecipeStep: fake.Bool(),
		IngredientNotes:     fake.Word(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
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

// BuildFakeRecipeStepIngredientUpdateInput builds a faked RecipeStepIngredientUpdateInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateInput() *types.RecipeStepIngredientUpdateInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return &types.RecipeStepIngredientUpdateInput{
		IngredientID:        recipeStepIngredient.IngredientID,
		Name:                recipeStepIngredient.Name,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient builds a faked RecipeStepIngredientUpdateInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateInput {
	return &types.RecipeStepIngredientUpdateInput{
		IngredientID:        recipeStepIngredient.IngredientID,
		Name:                recipeStepIngredient.Name,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientCreationInput builds a faked RecipeStepIngredientCreationInput.
func BuildFakeRecipeStepIngredientCreationInput() *types.RecipeStepIngredientCreationInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(recipeStepIngredient)
}

// BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient builds a faked RecipeStepIngredientCreationInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientCreationInput {
	return &types.RecipeStepIngredientCreationInput{
		IngredientID:        recipeStepIngredient.IngredientID,
		Name:                recipeStepIngredient.Name,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipeStep: recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}
