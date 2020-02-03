package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepProduct represents a recipe step product
	RecipeStepProduct struct {
		ID           uint64  `json:"id"`
		Name         string  `json:"name"`
		RecipeStepID uint64  `json:"recipe_step_id"`
		CreatedOn    uint64  `json:"created_on"`
		UpdatedOn    *uint64 `json:"updated_on"`
		ArchivedOn   *uint64 `json:"archived_on"`
		BelongsTo    uint64  `json:"belongs_to"`
	}

	// RecipeStepProductList represents a list of recipe step products
	RecipeStepProductList struct {
		Pagination
		RecipeStepProducts []RecipeStepProduct `json:"recipe_step_products"`
	}

	// RecipeStepProductCreationInput represents what a user could set as input for creating recipe step products
	RecipeStepProductCreationInput struct {
		Name         string `json:"name"`
		RecipeStepID uint64 `json:"recipe_step_id"`
		BelongsTo    uint64 `json:"-"`
	}

	// RecipeStepProductUpdateInput represents what a user could set as input for updating recipe step products
	RecipeStepProductUpdateInput struct {
		Name         string `json:"name"`
		RecipeStepID uint64 `json:"recipe_step_id"`
		BelongsTo    uint64 `json:"-"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently
	RecipeStepProductDataManager interface {
		GetRecipeStepProduct(ctx context.Context, recipeStepProductID, userID uint64) (*RecipeStepProduct, error)
		GetRecipeStepProductCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeStepProductsCount(ctx context.Context) (uint64, error)
		GetRecipeStepProducts(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeStepProductList, error)
		GetAllRecipeStepProductsForUser(ctx context.Context, userID uint64) ([]RecipeStepProduct, error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, id, userID uint64) error
	}

	// RecipeStepProductDataServer describes a structure capable of serving traffic related to recipe step products
	RecipeStepProductDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepProductInput with a recipe step product
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

// ToInput creates a RecipeStepProductUpdateInput struct for a recipe step product
func (x *RecipeStepProduct) ToInput() *RecipeStepProductUpdateInput {
	return &RecipeStepProductUpdateInput{
		Name:         x.Name,
		RecipeStepID: x.RecipeStepID,
	}
}
