package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepInstrumentDataType indicates an event is related to a recipe step instrument.
	RecipeStepInstrumentDataType dataType = "recipe_step_instrument"

	// RecipeStepInstrumentCreatedCustomerEventType indicates a recipe step instrument was created.
	RecipeStepInstrumentCreatedCustomerEventType CustomerEventType = "recipe_step_instrument_created"
	// RecipeStepInstrumentUpdatedCustomerEventType indicates a recipe step instrument was updated.
	RecipeStepInstrumentUpdatedCustomerEventType CustomerEventType = "recipe_step_instrument_updated"
	// RecipeStepInstrumentArchivedCustomerEventType indicates a recipe step instrument was archived.
	RecipeStepInstrumentArchivedCustomerEventType CustomerEventType = "recipe_step_instrument_archived"
)

func init() {
	gob.Register(new(RecipeStepInstrument))
	gob.Register(new(RecipeStepInstrumentList))
	gob.Register(new(RecipeStepInstrumentCreationRequestInput))
	gob.Register(new(RecipeStepInstrumentUpdateRequestInput))
}

type (
	// RecipeStepInstrument represents a recipe step instrument.
	RecipeStepInstrument struct {
		_                   struct{}
		ArchivedOn          *uint64 `json:"archivedOn"`
		InstrumentID        *string `json:"instrumentID"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		RecipeStepID        string  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		ID                  string  `json:"id"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		CreatedOn           uint64  `json:"createdOn"`
	}

	// RecipeStepInstrumentList represents a list of recipe step instruments.
	RecipeStepInstrumentList struct {
		_                     struct{}
		RecipeStepInstruments []*RecipeStepInstrument `json:"data"`
		Pagination
	}

	// RecipeStepInstrumentCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentCreationRequestInput struct {
		_ struct{}

		ID                  string  `json:"-"`
		InstrumentID        *string `json:"instrumentID"`
		RecipeStepID        string  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		BelongsToRecipeStep string  `json:"-"`
	}

	// RecipeStepInstrumentDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentDatabaseCreationInput struct {
		_ struct{}

		ID                  string  `json:"id"`
		InstrumentID        *string `json:"instrumentID"`
		RecipeStepID        string  `json:"recipeStepID"`
		Notes               string  `json:"notes"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
	}

	// RecipeStepInstrumentUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepInstrumentUpdateRequestInput struct {
		_ struct{}

		// InstrumentID is already a pointer, I'm not about to make it a double pointer.
		InstrumentID        *string `json:"instrumentID"`
		RecipeStepID        *string `json:"recipeStepID"`
		Notes               *string `json:"notes"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
	}

	// RecipeStepInstrumentDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepInstrumentDataManager interface {
		RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepInstrument, error)
		GetTotalRecipeStepInstrumentCount(ctx context.Context) (uint64, error)
		GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepInstrumentList, error)
		GetRecipeStepInstrumentsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*RecipeStepInstrument, error)
		CreateRecipeStepInstrument(ctx context.Context, input *RecipeStepInstrumentDatabaseCreationInput) (*RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, updated *RecipeStepInstrument) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}

	// RecipeStepInstrumentDataService describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepInstrumentDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepInstrumentUpdateRequestInput with a recipe step instrument.
func (x *RecipeStepInstrument) Update(input *RecipeStepInstrumentUpdateRequestInput) {
	if input.InstrumentID != nil && (x.InstrumentID == nil || (*input.InstrumentID != "" && *input.InstrumentID != *x.InstrumentID)) {
		x.InstrumentID = input.InstrumentID
	}

	if input.RecipeStepID != nil && *input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = *input.RecipeStepID
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentCreationRequestInput.
func (x *RecipeStepInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.InstrumentID, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentDatabaseCreationInput.
func (x *RecipeStepInstrumentDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.InstrumentID, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

// RecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrumentCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrumentCreationInput(input *RecipeStepInstrumentCreationRequestInput) *RecipeStepInstrumentDatabaseCreationInput {
	x := &RecipeStepInstrumentDatabaseCreationInput{
		InstrumentID:        input.InstrumentID,
		RecipeStepID:        input.RecipeStepID,
		Notes:               input.Notes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentUpdateRequestInput.
func (x *RecipeStepInstrumentUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.RecipeStepID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
