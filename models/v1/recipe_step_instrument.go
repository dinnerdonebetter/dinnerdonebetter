package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepInstrument represents a recipe step instrument.
	RecipeStepInstrument struct {
		ID                  uint64  `json:"id"`
		InstrumentID        *uint64 `json:"instrumentID"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		CreatedOn           uint64  `json:"createdOn"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepInstrumentList represents a list of recipe step instruments.
	RecipeStepInstrumentList struct {
		Pagination
		RecipeStepInstruments []RecipeStepInstrument `json:"recipeStepInstruments"`
	}

	// RecipeStepInstrumentCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentCreationInput struct {
		InstrumentID        *uint64 `json:"instrumentID"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		BelongsToRecipeStep uint64  `json:"-"`
	}

	// RecipeStepInstrumentUpdateInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepInstrumentUpdateInput struct {
		InstrumentID        *uint64 `json:"instrumentID"`
		RecipeStepID        uint64  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepInstrumentDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepInstrumentDataManager interface {
		RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (bool, error)
		GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID uint64) (*RecipeStepInstrument, error)
		GetAllRecipeStepInstrumentsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepInstruments(ctx context.Context, resultChannel chan []RecipeStepInstrument) error
		GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepInstrumentList, error)
		GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]RecipeStepInstrument, error)
		CreateRecipeStepInstrument(ctx context.Context, input *RecipeStepInstrumentCreationInput) (*RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, updated *RecipeStepInstrument) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID uint64) error
	}

	// RecipeStepInstrumentDataServer describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepInstrumentDataServer interface {
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

// Update merges an RecipeStepInstrumentInput with a recipe step instrument.
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

// ToUpdateInput creates a RecipeStepInstrumentUpdateInput struct for a recipe step instrument.
func (x *RecipeStepInstrument) ToUpdateInput() *RecipeStepInstrumentUpdateInput {
	return &RecipeStepInstrumentUpdateInput{
		InstrumentID: x.InstrumentID,
		RecipeStepID: x.RecipeStepID,
		Notes:        x.Notes,
	}
}
