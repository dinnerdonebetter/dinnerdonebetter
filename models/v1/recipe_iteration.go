package models

import (
	"context"
	"net/http"
)

type (
	// RecipeIteration represents a recipe iteration
	RecipeIteration struct {
		ID                  uint64  `json:"id"`
		RecipeID            uint64  `json:"recipe_id"`
		EndDifficultyRating float32 `json:"end_difficulty_rating"`
		EndComplexityRating float32 `json:"end_complexity_rating"`
		EndTasteRating      float32 `json:"end_taste_rating"`
		EndOverallRating    float32 `json:"end_overall_rating"`
		CreatedOn           uint64  `json:"created_on"`
		UpdatedOn           *uint64 `json:"updated_on"`
		ArchivedOn          *uint64 `json:"archived_on"`
		BelongsTo           uint64  `json:"belongs_to"`
	}

	// RecipeIterationList represents a list of recipe iterations
	RecipeIterationList struct {
		Pagination
		RecipeIterations []RecipeIteration `json:"recipe_iterations"`
	}

	// RecipeIterationCreationInput represents what a user could set as input for creating recipe iterations
	RecipeIterationCreationInput struct {
		RecipeID            uint64  `json:"recipe_id"`
		EndDifficultyRating float32 `json:"end_difficulty_rating"`
		EndComplexityRating float32 `json:"end_complexity_rating"`
		EndTasteRating      float32 `json:"end_taste_rating"`
		EndOverallRating    float32 `json:"end_overall_rating"`
		BelongsTo           uint64  `json:"-"`
	}

	// RecipeIterationUpdateInput represents what a user could set as input for updating recipe iterations
	RecipeIterationUpdateInput struct {
		RecipeID            uint64  `json:"recipe_id"`
		EndDifficultyRating float32 `json:"end_difficulty_rating"`
		EndComplexityRating float32 `json:"end_complexity_rating"`
		EndTasteRating      float32 `json:"end_taste_rating"`
		EndOverallRating    float32 `json:"end_overall_rating"`
		BelongsTo           uint64  `json:"-"`
	}

	// RecipeIterationDataManager describes a structure capable of storing recipe iterations permanently
	RecipeIterationDataManager interface {
		GetRecipeIteration(ctx context.Context, recipeIterationID, userID uint64) (*RecipeIteration, error)
		GetRecipeIterationCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeIterationsCount(ctx context.Context) (uint64, error)
		GetRecipeIterations(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeIterationList, error)
		GetAllRecipeIterationsForUser(ctx context.Context, userID uint64) ([]RecipeIteration, error)
		CreateRecipeIteration(ctx context.Context, input *RecipeIterationCreationInput) (*RecipeIteration, error)
		UpdateRecipeIteration(ctx context.Context, updated *RecipeIteration) error
		ArchiveRecipeIteration(ctx context.Context, id, userID uint64) error
	}

	// RecipeIterationDataServer describes a structure capable of serving traffic related to recipe iterations
	RecipeIterationDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeIterationInput with a recipe iteration
func (x *RecipeIteration) Update(input *RecipeIterationUpdateInput) {
	if input.RecipeID != x.RecipeID {
		x.RecipeID = input.RecipeID
	}

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

// ToInput creates a RecipeIterationUpdateInput struct for a recipe iteration
func (x *RecipeIteration) ToInput() *RecipeIterationUpdateInput {
	return &RecipeIterationUpdateInput{
		RecipeID:            x.RecipeID,
		EndDifficultyRating: x.EndDifficultyRating,
		EndComplexityRating: x.EndComplexityRating,
		EndTasteRating:      x.EndTasteRating,
		EndOverallRating:    x.EndOverallRating,
	}
}
