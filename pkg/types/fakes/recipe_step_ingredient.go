package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeRecipeStepIngredient builds a faked recipe step ingredient.
func BuildFakeRecipeStepIngredient() *types.RecipeStepIngredient {
	return &types.RecipeStepIngredient{
		ID:                  ksuid.New().String(),
		IngredientID:        func(x string) *string { return &x }(fake.LoremIpsumSentence(exampleQuantity)),
		QuantityType:        fake.LoremIpsumSentence(exampleQuantity),
		QuantityValue:       fake.Float32(),
		QuantityNotes:       fake.LoremIpsumSentence(exampleQuantity),
		ProductOfRecipeStep: fake.Bool(),
		IngredientNotes:     fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.UUID(),
	}
}

// BuildFakeFullRecipeStepIngredient builds a faked recipe step ingredient.
func BuildFakeFullRecipeStepIngredient() *types.FullRecipeStepIngredient {
	return &types.FullRecipeStepIngredient{
		ID:                  ksuid.New().String(),
		Ingredient:          *BuildFakeValidIngredient(),
		QuantityType:        fake.LoremIpsumSentence(exampleQuantity),
		QuantityValue:       float32(fake.Uint32()),
		QuantityNotes:       fake.LoremIpsumSentence(exampleQuantity),
		ProductOfRecipeStep: fake.Bool(),
		IngredientNotes:     fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: ksuid.New().String(),
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
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient builds a faked RecipeStepIngredientUpdateRequestInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	return &types.RecipeStepIngredientUpdateRequestInput{
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
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
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepIngredientDatabaseCreationInput builds a faked RecipeStepIngredientDatabaseCreationInput.
func BuildFakeRecipeStepIngredientDatabaseCreationInput() *types.RecipeStepIngredientDatabaseCreationInput {
	recipeStepIngredient := BuildFakeRecipeStepIngredient()
	return BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(recipeStepIngredient)
}

// BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient builds a faked RecipeStepIngredientDatabaseCreationInput from a recipe step ingredient.
func BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(recipeStepIngredient *types.RecipeStepIngredient) *types.RecipeStepIngredientDatabaseCreationInput {
	return &types.RecipeStepIngredientDatabaseCreationInput{
		ID:                  recipeStepIngredient.ID,
		IngredientID:        recipeStepIngredient.IngredientID,
		QuantityType:        recipeStepIngredient.QuantityType,
		QuantityValue:       recipeStepIngredient.QuantityValue,
		QuantityNotes:       recipeStepIngredient.QuantityNotes,
		ProductOfRecipe:     recipeStepIngredient.ProductOfRecipeStep,
		IngredientNotes:     recipeStepIngredient.IngredientNotes,
		BelongsToRecipeStep: recipeStepIngredient.BelongsToRecipeStep,
	}
}
