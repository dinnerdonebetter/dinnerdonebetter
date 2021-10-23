package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *types.RecipeStepProduct {
	return &types.RecipeStepProduct{
		ID:                  ksuid.New().String(),
		Name:                fake.LoremIpsumSentence(exampleQuantity),
		RecipeStepID:        fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.UUID(),
	}
}

// BuildFakeRecipeStepProductList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductList() *types.RecipeStepProductList {
	var examples []*types.RecipeStepProduct
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepProduct())
	}

	return &types.RecipeStepProductList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		RecipeStepProducts: examples,
	}
}

// BuildFakeRecipeStepProductUpdateRequestInput builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInput() *types.RecipeStepProductUpdateRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return &types.RecipeStepProductUpdateRequestInput{
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct builds a faked RecipeStepProductUpdateRequestInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateRequestInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	return &types.RecipeStepProductUpdateRequestInput{
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductCreationRequestInput builds a faked RecipeStepProductCreationRequestInput.
func BuildFakeRecipeStepProductCreationRequestInput() *types.RecipeStepProductCreationRequestInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(recipeStepProduct)
}

// BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct builds a faked RecipeStepProductCreationRequestInput from a recipe step product.
func BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductCreationRequestInput {
	return &types.RecipeStepProductCreationRequestInput{
		ID:                  recipeStepProduct.ID,
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductDatabaseCreationInput builds a faked RecipeStepProductDatabaseCreationInput.
func BuildFakeRecipeStepProductDatabaseCreationInput() *types.RecipeStepProductDatabaseCreationInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct(recipeStepProduct)
}

// BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct builds a faked RecipeStepProductDatabaseCreationInput from a recipe step product.
func BuildFakeRecipeStepProductDatabaseCreationInputFromRecipeStepProduct(recipeStepProduct *types.RecipeStepProduct) *types.RecipeStepProductDatabaseCreationInput {
	return &types.RecipeStepProductDatabaseCreationInput{
		ID:                  recipeStepProduct.ID,
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}
