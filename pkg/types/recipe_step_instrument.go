package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

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
	gob.Register(new(RecipeStepInstrumentCreationRequestInput))
	gob.Register(new(RecipeStepInstrumentUpdateRequestInput))
}

type (
	// RecipeStepInstrument represents a recipe step instrument.
	RecipeStepInstrument struct {
		_ struct{}

		CreatedAt           time.Time        `json:"createdAt"`
		Instrument          *ValidInstrument `json:"instrument"`
		LastUpdatedAt       *time.Time       `json:"lastUpdatedAt"`
		RecipeStepProductID *string          `json:"recipeStepProductID"`
		ArchivedAt          *time.Time       `json:"archivedAt"`
		Notes               string           `json:"notes"`
		Name                string           `json:"name"`
		BelongsToRecipeStep string           `json:"belongsToRecipeStep"`
		ID                  string           `json:"id"`
		MinimumQuantity     uint32           `json:"minimumQuantity"`
		MaximumQuantity     uint32           `json:"maximumQuantity"`
		OptionIndex         uint16           `json:"optionIndex"`
		PreferenceRank      uint8            `json:"preferenceRank"`
		Optional            bool             `json:"optional"`
	}

	// RecipeStepInstrumentCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentCreationRequestInput struct {
		_                               struct{}
		InstrumentID                    *string `json:"instrumentID"`
		RecipeStepProductID             *string `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64 `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64 `json:"productIndexOfRecipeStep"`
		Name                            string  `json:"name"`
		Notes                           string  `json:"notes"`
		MinimumQuantity                 uint32  `json:"minimumQuantity"`
		MaximumQuantity                 uint32  `json:"maximumQuantity"`
		OptionIndex                     uint16  `json:"optionIndex"`
		Optional                        bool    `json:"optional"`
		PreferenceRank                  uint8   `json:"preferenceRank"`
	}

	// RecipeStepInstrumentDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentDatabaseCreationInput struct {
		_                               struct{}
		InstrumentID                    *string
		RecipeStepProductID             *string
		ProductOfRecipeStepIndex        *uint64
		ProductOfRecipeStepProductIndex *uint64
		Notes                           string
		BelongsToRecipeStep             string
		Name                            string
		ID                              string
		MinimumQuantity                 uint32
		MaximumQuantity                 uint32
		OptionIndex                     uint16
		Optional                        bool
		PreferenceRank                  uint8
	}

	// RecipeStepInstrumentUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepInstrumentUpdateRequestInput struct {
		_ struct{}

		InstrumentID        *string `json:"instrumentID"`
		RecipeStepProductID *string `json:"recipeStepProductID"`
		Notes               *string `json:"notes"`
		PreferenceRank      *uint8  `json:"preferenceRank"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		Name                *string `json:"name"`
		Optional            *bool   `json:"optional"`
		OptionIndex         *uint16 `json:"optionIndex"`
		MinimumQuantity     *uint32 `json:"minimumQuantity"`
		MaximumQuantity     *uint32 `json:"maximumQuantity"`
	}

	// RecipeStepInstrumentDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepInstrumentDataManager interface {
		RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepInstrument, error)
		GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepInstrument], error)
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
	if input.InstrumentID != nil && (x.Instrument == nil || (*input.InstrumentID != "" && *input.InstrumentID != x.Instrument.ID)) {
		x.Instrument.ID = *input.InstrumentID
	}

	if input.RecipeStepProductID != nil && *input.RecipeStepProductID != *x.RecipeStepProductID {
		x.RecipeStepProductID = input.RecipeStepProductID
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.PreferenceRank != nil && *input.PreferenceRank != x.PreferenceRank {
		x.PreferenceRank = *input.PreferenceRank
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && *input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = *input.MaximumQuantity
	}

	if input.OptionIndex != nil && *input.OptionIndex != x.OptionIndex {
		x.OptionIndex = *input.OptionIndex
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentCreationRequestInput.
func (x *RecipeStepInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.InstrumentID, validation.Required),
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
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentUpdateRequestInput.
func (x *RecipeStepInstrumentUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
