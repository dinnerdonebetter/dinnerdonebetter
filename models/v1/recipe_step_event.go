package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepEvent represents a recipe step event
	RecipeStepEvent struct {
		ID                uint64  `json:"id"`
		EventType         string  `json:"event_type"`
		Done              bool    `json:"done"`
		RecipeIterationID uint64  `json:"recipe_iteration_id"`
		RecipeStepID      uint64  `json:"recipe_step_id"`
		CreatedOn         uint64  `json:"created_on"`
		UpdatedOn         *uint64 `json:"updated_on"`
		ArchivedOn        *uint64 `json:"archived_on"`
		BelongsTo         uint64  `json:"belongs_to"`
	}

	// RecipeStepEventList represents a list of recipe step events
	RecipeStepEventList struct {
		Pagination
		RecipeStepEvents []RecipeStepEvent `json:"recipe_step_events"`
	}

	// RecipeStepEventCreationInput represents what a user could set as input for creating recipe step events
	RecipeStepEventCreationInput struct {
		EventType         string `json:"event_type"`
		Done              bool   `json:"done"`
		RecipeIterationID uint64 `json:"recipe_iteration_id"`
		RecipeStepID      uint64 `json:"recipe_step_id"`
		BelongsTo         uint64 `json:"-"`
	}

	// RecipeStepEventUpdateInput represents what a user could set as input for updating recipe step events
	RecipeStepEventUpdateInput struct {
		EventType         string `json:"event_type"`
		Done              bool   `json:"done"`
		RecipeIterationID uint64 `json:"recipe_iteration_id"`
		RecipeStepID      uint64 `json:"recipe_step_id"`
		BelongsTo         uint64 `json:"-"`
	}

	// RecipeStepEventDataManager describes a structure capable of storing recipe step events permanently
	RecipeStepEventDataManager interface {
		GetRecipeStepEvent(ctx context.Context, recipeStepEventID, userID uint64) (*RecipeStepEvent, error)
		GetRecipeStepEventCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeStepEventsCount(ctx context.Context) (uint64, error)
		GetRecipeStepEvents(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeStepEventList, error)
		GetAllRecipeStepEventsForUser(ctx context.Context, userID uint64) ([]RecipeStepEvent, error)
		CreateRecipeStepEvent(ctx context.Context, input *RecipeStepEventCreationInput) (*RecipeStepEvent, error)
		UpdateRecipeStepEvent(ctx context.Context, updated *RecipeStepEvent) error
		ArchiveRecipeStepEvent(ctx context.Context, id, userID uint64) error
	}

	// RecipeStepEventDataServer describes a structure capable of serving traffic related to recipe step events
	RecipeStepEventDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepEventInput with a recipe step event
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

// ToInput creates a RecipeStepEventUpdateInput struct for a recipe step event
func (x *RecipeStepEvent) ToInput() *RecipeStepEventUpdateInput {
	return &RecipeStepEventUpdateInput{
		EventType:         x.EventType,
		Done:              x.Done,
		RecipeIterationID: x.RecipeIterationID,
		RecipeStepID:      x.RecipeStepID,
	}
}
