package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

type (
	// RecipeStepInstrument represents a recipe step instrument.
	RecipeStepInstrument struct {
		_ struct{} `json:"-"`

		CreatedAt           time.Time                  `json:"createdAt"`
		Instrument          *ValidInstrument           `json:"instrument"`
		LastUpdatedAt       *time.Time                 `json:"lastUpdatedAt"`
		RecipeStepProductID *string                    `json:"recipeStepProductID"`
		ArchivedAt          *time.Time                 `json:"archivedAt"`
		Notes               string                     `json:"notes"`
		Name                string                     `json:"name"`
		BelongsToRecipeStep string                     `json:"belongsToRecipeStep"`
		ID                  string                     `json:"id"`
		Quantity            Uint32RangeWithOptionalMax `json:"quantity"`
		OptionIndex         uint16                     `json:"optionIndex"`
		PreferenceRank      uint8                      `json:"preferenceRank"`
		Optional            bool                       `json:"optional"`
	}

	// RecipeStepInstrumentCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentCreationRequestInput struct {
		_ struct{} `json:"-"`

		InstrumentID                    *string                    `json:"instrumentID"`
		RecipeStepProductID             *string                    `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64                    `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64                    `json:"productOfRecipeStepProductIndex"`
		Quantity                        Uint32RangeWithOptionalMax `json:"quantity"`
		Notes                           string                     `json:"notes"`
		Name                            string                     `json:"name"`
		OptionIndex                     uint16                     `json:"optionIndex"`
		Optional                        bool                       `json:"optional"`
		PreferenceRank                  uint8                      `json:"preferenceRank"`
	}

	// RecipeStepInstrumentDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepInstrumentDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		InstrumentID                    *string                    `json:"-"`
		RecipeStepProductID             *string                    `json:"-"`
		ProductOfRecipeStepIndex        *uint64                    `json:"-"`
		ProductOfRecipeStepProductIndex *uint64                    `json:"-"`
		BelongsToRecipeStep             string                     `json:"-"`
		Name                            string                     `json:"-"`
		ID                              string                     `json:"-"`
		Notes                           string                     `json:"-"`
		Quantity                        Uint32RangeWithOptionalMax `json:"-"`
		OptionIndex                     uint16                     `json:"-"`
		Optional                        bool                       `json:"-"`
		PreferenceRank                  uint8                      `json:"-"`
	}

	// RecipeStepInstrumentUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepInstrumentUpdateRequestInput struct {
		_ struct{} `json:"-"`

		InstrumentID        *string                                      `json:"instrumentID,omitempty"`
		RecipeStepProductID *string                                      `json:"recipeStepProductID,omitempty"`
		Notes               *string                                      `json:"notes,omitempty"`
		PreferenceRank      *uint8                                       `json:"preferenceRank,omitempty"`
		BelongsToRecipeStep *string                                      `json:"belongsToRecipeStep,omitempty"`
		Name                *string                                      `json:"name,omitempty"`
		Optional            *bool                                        `json:"optional,omitempty"`
		OptionIndex         *uint16                                      `json:"optionIndex,omitempty"`
		Quantity            Uint32RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
	}

	// RecipeStepInstrumentDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepInstrumentDataManager interface {
		RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepInstrument, error)
		GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipeStepInstrument], error)
		CreateRecipeStepInstrument(ctx context.Context, input *RecipeStepInstrumentDatabaseCreationInput) (*RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, updated *RecipeStepInstrument) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}

	// RecipeStepInstrumentDataService describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepInstrumentDataService interface {
		ListRecipeStepInstrumentsHandler(http.ResponseWriter, *http.Request)
		CreateRecipeStepInstrumentHandler(http.ResponseWriter, *http.Request)
		ReadRecipeStepInstrumentHandler(http.ResponseWriter, *http.Request)
		UpdateRecipeStepInstrumentHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipeStepInstrumentHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepInstrumentUpdateRequestInput with a recipe step instrument.
func (x *RecipeStepInstrument) Update(input *RecipeStepInstrumentUpdateRequestInput) {
	if input.InstrumentID != nil && (x.Instrument == nil || (*input.InstrumentID != "" && *input.InstrumentID != x.Instrument.ID)) {
		x.Instrument = &ValidInstrument{ID: *input.InstrumentID}
	}

	if input.RecipeStepProductID != nil && x.RecipeStepProductID != nil && *input.RecipeStepProductID != *x.RecipeStepProductID {
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

	if input.Quantity.Min != nil && *input.Quantity.Min != x.Quantity.Min {
		x.Quantity.Min = *input.Quantity.Min
	}

	if input.Quantity.Max != nil && x.Quantity.Max != nil && *input.Quantity.Max != *x.Quantity.Max {
		x.Quantity.Max = input.Quantity.Max
	}

	if input.OptionIndex != nil && *input.OptionIndex != x.OptionIndex {
		x.OptionIndex = *input.OptionIndex
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepInstrumentCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepInstrumentCreationRequestInput.
func (x *RecipeStepInstrumentCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	err := &multierror.Error{}

	if x.InstrumentID == nil && x.ProductOfRecipeStepIndex == nil && x.ProductOfRecipeStepProductIndex == nil {
		err = multierror.Append(err, ErrInstrumentIDOrProductIndicesRequired)
	}

	validationErr := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
	if validationErr != nil {
		err = multierror.Append(err, validationErr)
	}

	return err.ErrorOrNil()
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
