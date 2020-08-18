package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		ID                  uint64  `json:"id"`
		Name                string  `json:"name"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		CreatedOn           uint64  `json:"createdOn"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepProductList represents a list of recipe step products.
	RecipeStepProductList struct {
		Pagination
		RecipeStepProducts []RecipeStepProduct `json:"recipe_step_products"`
	}

	// RecipeStepProductCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationInput struct {
		Name                string `json:"name"`
		RecipeStepID        uint64 `json:"recipeStepID"`
		BelongsToRecipeStep uint64 `json:"-"`
	}

	// RecipeStepProductUpdateInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateInput struct {
		Name                string `json:"name"`
		RecipeStepID        uint64 `json:"recipeStepID"`
		BelongsToRecipeStep uint64 `json:"belongsToRecipeStep"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (*RecipeStepProduct, error)
		GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepProducts(ctx context.Context, resultChannel chan []RecipeStepProduct) error
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepProductList, error)
		GetRecipeStepProductsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]RecipeStepProduct, error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID uint64) error
	}

	// RecipeStepProductDataServer describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepProductInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

// ToUpdateInput creates a RecipeStepProductUpdateInput struct for a recipe step product.
func (x *RecipeStepProduct) ToUpdateInput() *RecipeStepProductUpdateInput {
	return &RecipeStepProductUpdateInput{
		Name:         x.Name,
		RecipeStepID: x.RecipeStepID,
	}
}
