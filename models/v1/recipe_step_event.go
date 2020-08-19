package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepEvent represents a recipe step event.
	RecipeStepEvent struct {
		ID                  uint64  `json:"id"`
		EventType           string  `json:"eventType"`
		Done                bool    `json:"done"`
		RecipeIterationID   uint64  `json:"recipeIterationID"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		CreatedOn           uint64  `json:"createdOn"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepEventList represents a list of recipe step events.
	RecipeStepEventList struct {
		Pagination
		RecipeStepEvents []RecipeStepEvent `json:"recipeStepEvents"`
	}

	// RecipeStepEventCreationInput represents what a user could set as input for creating recipe step events.
	RecipeStepEventCreationInput struct {
		EventType           string `json:"eventType"`
		Done                bool   `json:"done"`
		RecipeIterationID   uint64 `json:"recipeIterationID"`
		RecipeStepID        uint64 `json:"recipeStepID"`
		BelongsToRecipeStep uint64 `json:"-"`
	}

	// RecipeStepEventUpdateInput represents what a user could set as input for updating recipe step events.
	RecipeStepEventUpdateInput struct {
		EventType           string `json:"eventType"`
		Done                bool   `json:"done"`
		RecipeIterationID   uint64 `json:"recipeIterationID"`
		RecipeStepID        uint64 `json:"recipeStepID"`
		BelongsToRecipeStep uint64 `json:"belongsToRecipeStep"`
	}

	// RecipeStepEventDataManager describes a structure capable of storing recipe step events permanently.
	RecipeStepEventDataManager interface {
		RecipeStepEventExists(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (bool, error)
		GetRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*RecipeStepEvent, error)
		GetAllRecipeStepEventsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepEvents(ctx context.Context, resultChannel chan []RecipeStepEvent) error
		GetRecipeStepEvents(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepEventList, error)
		GetRecipeStepEventsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]RecipeStepEvent, error)
		CreateRecipeStepEvent(ctx context.Context, input *RecipeStepEventCreationInput) (*RecipeStepEvent, error)
		UpdateRecipeStepEvent(ctx context.Context, updated *RecipeStepEvent) error
		ArchiveRecipeStepEvent(ctx context.Context, recipeStepID, recipeStepEventID uint64) error
	}

	// RecipeStepEventDataServer describes a structure capable of serving traffic related to recipe step events.
	RecipeStepEventDataServer interface {
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

// Update merges an RecipeStepEventInput with a recipe step event.
func (x *RecipeStepEvent) Update(input *RecipeStepEventUpdateInput) {
	if input.EventType != "" && input.EventType != x.EventType {
		x.EventType = input.EventType
	}

	if input.Done != x.Done {
		x.Done = input.Done
	}

	if input.RecipeIterationID != x.RecipeIterationID {
		x.RecipeIterationID = input.RecipeIterationID
	}

	if input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

// ToUpdateInput creates a RecipeStepEventUpdateInput struct for a recipe step event.
func (x *RecipeStepEvent) ToUpdateInput() *RecipeStepEventUpdateInput {
	return &RecipeStepEventUpdateInput{
		EventType:         x.EventType,
		Done:              x.Done,
		RecipeIterationID: x.RecipeIterationID,
		RecipeStepID:      x.RecipeStepID,
	}
}
