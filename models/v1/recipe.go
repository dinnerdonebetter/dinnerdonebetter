package models

import (
	"context"
	"net/http"
)

type (
	// Recipe represents a recipe
	Recipe struct {
		ID                 uint64  `json:"id"`
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspired_by_recipe_id"`
		CreatedOn          uint64  `json:"created_on"`
		UpdatedOn          *uint64 `json:"updated_on"`
		ArchivedOn         *uint64 `json:"archived_on"`
		BelongsTo          uint64  `json:"belongs_to"`
	}

	// RecipeList represents a list of recipes
	RecipeList struct {
		Pagination
		Recipes []Recipe `json:"recipes"`
	}

	// RecipeCreationInput represents what a user could set as input for creating recipes
	RecipeCreationInput struct {
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspired_by_recipe_id"`
		BelongsTo          uint64  `json:"-"`
	}

	// RecipeUpdateInput represents what a user could set as input for updating recipes
	RecipeUpdateInput struct {
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspired_by_recipe_id"`
		BelongsTo          uint64  `json:"-"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently
	RecipeDataManager interface {
		GetRecipe(ctx context.Context, recipeID, userID uint64) (*Recipe, error)
		GetRecipeCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipesCount(ctx context.Context) (uint64, error)
		GetRecipes(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeList, error)
		GetAllRecipesForUser(ctx context.Context, userID uint64) ([]Recipe, error)
		CreateRecipe(ctx context.Context, input *RecipeCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		ArchiveRecipe(ctx context.Context, id, userID uint64) error
	}

	// RecipeDataServer describes a structure capable of serving traffic related to recipes
	RecipeDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeInput with a recipe
func (x *Recipe) Update(input *RecipeUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Source != "" && input.Source != x.Source {
		x.Source = input.Source
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.InspiredByRecipeID != nil && input.InspiredByRecipeID != x.InspiredByRecipeID {
		x.InspiredByRecipeID = input.InspiredByRecipeID
	}
}

// ToInput creates a RecipeUpdateInput struct for a recipe
func (x *Recipe) ToInput() *RecipeUpdateInput {
	return &RecipeUpdateInput{
		Name:               x.Name,
		Source:             x.Source,
		Description:        x.Description,
		InspiredByRecipeID: x.InspiredByRecipeID,
	}
}
