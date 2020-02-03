package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepInstrument represents a recipe step instrument
	RecipeStepInstrument struct {
		ID           uint64  `json:"id"`
		InstrumentID *uint64 `json:"instrument_id"`
		RecipeStepID uint64  `json:"recipe_step_id"`
		Notes        string  `json:"notes"`
		CreatedOn    uint64  `json:"created_on"`
		UpdatedOn    *uint64 `json:"updated_on"`
		ArchivedOn   *uint64 `json:"archived_on"`
		BelongsTo    uint64  `json:"belongs_to"`
	}

	// RecipeStepInstrumentList represents a list of recipe step instruments
	RecipeStepInstrumentList struct {
		Pagination
		RecipeStepInstruments []RecipeStepInstrument `json:"recipe_step_instruments"`
	}

	// RecipeStepInstrumentCreationInput represents what a user could set as input for creating recipe step instruments
	RecipeStepInstrumentCreationInput struct {
		InstrumentID *uint64 `json:"instrument_id"`
		RecipeStepID uint64  `json:"recipe_step_id"`
		Notes        string  `json:"notes"`
		BelongsTo    uint64  `json:"-"`
	}

	// RecipeStepInstrumentUpdateInput represents what a user could set as input for updating recipe step instruments
	RecipeStepInstrumentUpdateInput struct {
		InstrumentID *uint64 `json:"instrument_id"`
		RecipeStepID uint64  `json:"recipe_step_id"`
		Notes        string  `json:"notes"`
		BelongsTo    uint64  `json:"-"`
	}

	// RecipeStepInstrumentDataManager describes a structure capable of storing recipe step instruments permanently
	RecipeStepInstrumentDataManager interface {
		GetRecipeStepInstrument(ctx context.Context, recipeStepInstrumentID, userID uint64) (*RecipeStepInstrument, error)
		GetRecipeStepInstrumentCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeStepInstrumentsCount(ctx context.Context) (uint64, error)
		GetRecipeStepInstruments(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeStepInstrumentList, error)
		GetAllRecipeStepInstrumentsForUser(ctx context.Context, userID uint64) ([]RecipeStepInstrument, error)
		CreateRecipeStepInstrument(ctx context.Context, input *RecipeStepInstrumentCreationInput) (*RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, updated *RecipeStepInstrument) error
		ArchiveRecipeStepInstrument(ctx context.Context, id, userID uint64) error
	}

	// RecipeStepInstrumentDataServer describes a structure capable of serving traffic related to recipe step instruments
	RecipeStepInstrumentDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepInstrumentInput with a recipe step instrument
func (x *RecipeStepInstrument) Update(input *RecipeStepInstrumentUpdateInput) {
	if input.InstrumentID != nil && input.InstrumentID != x.InstrumentID {
		x.InstrumentID = input.InstrumentID
	}

	if input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToInput creates a RecipeStepInstrumentUpdateInput struct for a recipe step instrument
func (x *RecipeStepInstrument) ToInput() *RecipeStepInstrumentUpdateInput {
	return &RecipeStepInstrumentUpdateInput{
		InstrumentID: x.InstrumentID,
		RecipeStepID: x.RecipeStepID,
		Notes:        x.Notes,
	}
}
