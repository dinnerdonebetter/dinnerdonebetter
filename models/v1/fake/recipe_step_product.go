package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepProduct builds a faked recipe step product.
func BuildFakeRecipeStepProduct() *models.RecipeStepProduct {
	return &models.RecipeStepProduct{
		ID:                  fake.Uint64(),
		Name:                fake.Word(),
		RecipeStepID:        uint64(fake.Uint32()),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepProductList builds a faked RecipeStepProductList.
func BuildFakeRecipeStepProductList() *models.RecipeStepProductList {
	exampleRecipeStepProduct1 := BuildFakeRecipeStepProduct()
	exampleRecipeStepProduct2 := BuildFakeRecipeStepProduct()
	exampleRecipeStepProduct3 := BuildFakeRecipeStepProduct()

	return &models.RecipeStepProductList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		RecipeStepProducts: []models.RecipeStepProduct{
			*exampleRecipeStepProduct1,
			*exampleRecipeStepProduct2,
			*exampleRecipeStepProduct3,
		},
	}
}

// BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct builds a faked RecipeStepProductUpdateInput from a recipe step product.
func BuildFakeRecipeStepProductUpdateInputFromRecipeStepProduct(recipeStepProduct *models.RecipeStepProduct) *models.RecipeStepProductUpdateInput {
	return &models.RecipeStepProductUpdateInput{
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepProductCreationInput builds a faked RecipeStepProductCreationInput.
func BuildFakeRecipeStepProductCreationInput() *models.RecipeStepProductCreationInput {
	recipeStepProduct := BuildFakeRecipeStepProduct()
	return BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(recipeStepProduct)
}

// BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct builds a faked RecipeStepProductCreationInput from a recipe step product.
func BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(recipeStepProduct *models.RecipeStepProduct) *models.RecipeStepProductCreationInput {
	return &models.RecipeStepProductCreationInput{
		Name:                recipeStepProduct.Name,
		RecipeStepID:        recipeStepProduct.RecipeStepID,
		BelongsToRecipeStep: recipeStepProduct.BelongsToRecipeStep,
	}
}
