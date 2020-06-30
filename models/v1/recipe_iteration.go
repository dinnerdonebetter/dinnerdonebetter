package models

import (
	"context"
	"net/http"
)

type (
	// RecipeIteration represents a recipe iteration.
	RecipeIteration struct {
		ID                  uint64  `json:"id"`
		EndDifficultyRating float32 `json:"endDifficultyRating"`
		EndComplexityRating float32 `json:"endComplexityRating"`
		EndTasteRating      float32 `json:"endTasteRating"`
		EndOverallRating    float32 `json:"endOverallRating"`
		CreatedOn           uint64  `json:"createdOn"`
		UpdatedOn           *uint64 `json:"updatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipe     uint64  `json:"belongsToRecipe"`
	}

	// RecipeIterationList represents a list of recipe iterations.
	RecipeIterationList struct {
		Pagination
		RecipeIterations []RecipeIteration `json:"recipeIterations"`
	}

	// RecipeIterationCreationInput represents what a user could set as input for creating recipe iterations.
	RecipeIterationCreationInput struct {
		EndDifficultyRating float32 `json:"endDifficultyRating"`
		EndComplexityRating float32 `json:"endComplexityRating"`
		EndTasteRating      float32 `json:"endTasteRating"`
		EndOverallRating    float32 `json:"endOverallRating"`
		BelongsToRecipe     uint64  `json:"-"`
	}

	// RecipeIterationUpdateInput represents what a user could set as input for updating recipe iterations.
	RecipeIterationUpdateInput struct {
		EndDifficultyRating float32 `json:"endDifficultyRating"`
		EndComplexityRating float32 `json:"endComplexityRating"`
		EndTasteRating      float32 `json:"endTasteRating"`
		EndOverallRating    float32 `json:"endOverallRating"`
		BelongsToRecipe     uint64  `json:"belongsToRecipe"`
	}

	// RecipeIterationDataManager describes a structure capable of storing recipe iterations permanently.
	RecipeIterationDataManager interface {
		RecipeIterationExists(ctx context.Context, recipeID, recipeIterationID uint64) (bool, error)
		GetRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) (*RecipeIteration, error)
		GetAllRecipeIterationsCount(ctx context.Context) (uint64, error)
		GetRecipeIterations(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeIterationList, error)
		CreateRecipeIteration(ctx context.Context, input *RecipeIterationCreationInput) (*RecipeIteration, error)
		UpdateRecipeIteration(ctx context.Context, updated *RecipeIteration) error
		ArchiveRecipeIteration(ctx context.Context, recipeID, recipeIterationID uint64) error
	}

	// RecipeIterationDataServer describes a structure capable of serving traffic related to recipe iterations.
	RecipeIterationDataServer interface {
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

// Update merges an RecipeIterationInput with a recipe iteration.
func (x *RecipeIteration) Update(input *RecipeIterationUpdateInput) {
	if input.EndDifficultyRating != x.EndDifficultyRating {
		x.EndDifficultyRating = input.EndDifficultyRating
	}

	if input.EndComplexityRating != x.EndComplexityRating {
		x.EndComplexityRating = input.EndComplexityRating
	}

	if input.EndTasteRating != x.EndTasteRating {
		x.EndTasteRating = input.EndTasteRating
	}

	if input.EndOverallRating != x.EndOverallRating {
		x.EndOverallRating = input.EndOverallRating
	}
}

// ToUpdateInput creates a RecipeIterationUpdateInput struct for a recipe iteration.
func (x *RecipeIteration) ToUpdateInput() *RecipeIterationUpdateInput {
	return &RecipeIterationUpdateInput{
		EndDifficultyRating: x.EndDifficultyRating,
		EndComplexityRating: x.EndComplexityRating,
		EndTasteRating:      x.EndTasteRating,
		EndOverallRating:    x.EndOverallRating,
	}
}
