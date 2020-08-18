package models

import (
	"context"
	"net/http"
)

type (
	// Recipe represents a recipe.
	Recipe struct {
		ID                 uint64  `json:"id"`
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspiredByRecipeID"`
		CreatedOn          uint64  `json:"createdOn"`
		LastUpdatedOn      *uint64 `json:"lastUpdatedOn"`
		ArchivedOn         *uint64 `json:"archivedOn"`
		BelongsToUser      uint64  `json:"belongsToUser"`
	}

	// RecipeList represents a list of recipes.
	RecipeList struct {
		Pagination
		Recipes []Recipe `json:"recipes"`
	}

	// RecipeCreationInput represents what a user could set as input for creating recipes.
	RecipeCreationInput struct {
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspiredByRecipeID"`
		BelongsToUser      uint64  `json:"-"`
	}

	// RecipeUpdateInput represents what a user could set as input for updating recipes.
	RecipeUpdateInput struct {
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *uint64 `json:"inspiredByRecipeID"`
		BelongsToUser      uint64  `json:"-"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID uint64) (bool, error)
		GetRecipe(ctx context.Context, recipeID uint64) (*Recipe, error)
		GetAllRecipesCount(ctx context.Context) (uint64, error)
		GetAllRecipes(ctx context.Context, resultChannel chan []Recipe) error
		GetRecipes(ctx context.Context, filter *QueryFilter) (*RecipeList, error)
		GetRecipesWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]Recipe, error)
		CreateRecipe(ctx context.Context, input *RecipeCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		ArchiveRecipe(ctx context.Context, recipeID, userID uint64) error
	}

	// RecipeDataServer describes a structure capable of serving traffic related to recipes.
	RecipeDataServer interface {
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

// Update merges an RecipeInput with a recipe.
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

// ToUpdateInput creates a RecipeUpdateInput struct for a recipe.
func (x *Recipe) ToUpdateInput() *RecipeUpdateInput {
	return &RecipeUpdateInput{
		Name:               x.Name,
		Source:             x.Source,
		Description:        x.Description,
		InspiredByRecipeID: x.InspiredByRecipeID,
	}
}
