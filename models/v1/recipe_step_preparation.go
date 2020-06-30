package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepPreparation represents a recipe step preparation.
	RecipeStepPreparation struct {
		ID                  uint64  `json:"id"`
		ValidPreparationID  uint64  `json:"validPreparationID"`
		Notes               string  `json:"notes"`
		CreatedOn           uint64  `json:"createdOn"`
		UpdatedOn           *uint64 `json:"updatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepPreparationList represents a list of recipe step preparations.
	RecipeStepPreparationList struct {
		Pagination
		RecipeStepPreparations []RecipeStepPreparation `json:"recipeStepPreparations"`
	}

	// RecipeStepPreparationCreationInput represents what a user could set as input for creating recipe step preparations.
	RecipeStepPreparationCreationInput struct {
		ValidPreparationID  uint64 `json:"validPreparationID"`
		Notes               string `json:"notes"`
		BelongsToRecipeStep uint64 `json:"-"`
	}

	// RecipeStepPreparationUpdateInput represents what a user could set as input for updating recipe step preparations.
	RecipeStepPreparationUpdateInput struct {
		ValidPreparationID  uint64 `json:"validPreparationID"`
		Notes               string `json:"notes"`
		BelongsToRecipeStep uint64 `json:"belongsToRecipeStep"`
	}

	// RecipeStepPreparationDataManager describes a structure capable of storing recipe step preparations permanently.
	RecipeStepPreparationDataManager interface {
		RecipeStepPreparationExists(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (bool, error)
		GetRecipeStepPreparation(ctx context.Context, recipeID, recipeStepID, recipeStepPreparationID uint64) (*RecipeStepPreparation, error)
		GetAllRecipeStepPreparationsCount(ctx context.Context) (uint64, error)
		GetRecipeStepPreparations(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepPreparationList, error)
		CreateRecipeStepPreparation(ctx context.Context, input *RecipeStepPreparationCreationInput) (*RecipeStepPreparation, error)
		UpdateRecipeStepPreparation(ctx context.Context, updated *RecipeStepPreparation) error
		ArchiveRecipeStepPreparation(ctx context.Context, recipeStepID, recipeStepPreparationID uint64) error
	}

	// RecipeStepPreparationDataServer describes a structure capable of serving traffic related to recipe step preparations.
	RecipeStepPreparationDataServer interface {
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

// Update merges an RecipeStepPreparationInput with a recipe step preparation.
func (x *RecipeStepPreparation) Update(input *RecipeStepPreparationUpdateInput) {
	if input.ValidPreparationID != x.ValidPreparationID {
		x.ValidPreparationID = input.ValidPreparationID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToUpdateInput creates a RecipeStepPreparationUpdateInput struct for a recipe step preparation.
func (x *RecipeStepPreparation) ToUpdateInput() *RecipeStepPreparationUpdateInput {
	return &RecipeStepPreparationUpdateInput{
		ValidPreparationID: x.ValidPreparationID,
		Notes:              x.Notes,
	}
}
