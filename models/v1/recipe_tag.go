package models

import (
	"context"
	"net/http"
)

type (
	// RecipeTag represents a recipe tag.
	RecipeTag struct {
		ID              uint64  `json:"id"`
		Name            string  `json:"name"`
		CreatedOn       uint64  `json:"created_on"`
		UpdatedOn       *uint64 `json:"updated_on"`
		ArchivedOn      *uint64 `json:"archived_on"`
		BelongsToRecipe uint64  `json:"belongs_to_recipe"`
	}

	// RecipeTagList represents a list of recipe tags.
	RecipeTagList struct {
		Pagination
		RecipeTags []RecipeTag `json:"recipe_tags"`
	}

	// RecipeTagCreationInput represents what a user could set as input for creating recipe tags.
	RecipeTagCreationInput struct {
		Name            string `json:"name"`
		BelongsToRecipe uint64 `json:"-"`
	}

	// RecipeTagUpdateInput represents what a user could set as input for updating recipe tags.
	RecipeTagUpdateInput struct {
		Name            string `json:"name"`
		BelongsToRecipe uint64 `json:"belongs_to_recipe"`
	}

	// RecipeTagDataManager describes a structure capable of storing recipe tags permanently.
	RecipeTagDataManager interface {
		RecipeTagExists(ctx context.Context, recipeID, recipeTagID uint64) (bool, error)
		GetRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) (*RecipeTag, error)
		GetAllRecipeTagsCount(ctx context.Context) (uint64, error)
		GetRecipeTags(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeTagList, error)
		CreateRecipeTag(ctx context.Context, input *RecipeTagCreationInput) (*RecipeTag, error)
		UpdateRecipeTag(ctx context.Context, updated *RecipeTag) error
		ArchiveRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) error
	}

	// RecipeTagDataServer describes a structure capable of serving traffic related to recipe tags.
	RecipeTagDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ExistenceHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeTagInput with a recipe tag.
func (x *RecipeTag) Update(input *RecipeTagUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}
}

// ToUpdateInput creates a RecipeTagUpdateInput struct for a recipe tag.
func (x *RecipeTag) ToUpdateInput() *RecipeTagUpdateInput {
	return &RecipeTagUpdateInput{
		Name: x.Name,
	}
}
