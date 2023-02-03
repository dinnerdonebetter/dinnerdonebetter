package types

import (
	"context"
	"encoding/gob"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepVesselDataType indicates an event is related to a recipe step instrument.
	RecipeStepVesselDataType dataType = "recipe_step_vessel"

	// RecipeStepVesselCreatedCustomerEventType indicates a recipe step instrument was created.
	RecipeStepVesselCreatedCustomerEventType CustomerEventType = "recipe_step_vessel_created"
	// RecipeStepVesselUpdatedCustomerEventType indicates a recipe step instrument was updated.
	RecipeStepVesselUpdatedCustomerEventType CustomerEventType = "recipe_step_vessel_updated"
	// RecipeStepVesselArchivedCustomerEventType indicates a recipe step instrument was archived.
	RecipeStepVesselArchivedCustomerEventType CustomerEventType = "recipe_step_vessel_archived"
)

func init() {
	gob.Register(new(RecipeStepVessel))
	gob.Register(new(RecipeStepVesselCreationRequestInput))
	gob.Register(new(RecipeStepVesselUpdateRequestInput))
}

type (
	// RecipeStepVessel represents a recipe step instrument.
	RecipeStepVessel struct {
		_                    struct{}
		CreatedAt            time.Time        `json:"createdAt"`
		Instrument           *ValidInstrument `json:"instrument"`
		LastUpdatedAt        *time.Time       `json:"lastUpdatedAt"`
		ArchivedAt           *time.Time       `json:"archivedAt"`
		RecipeStepProductID  *string          `json:"recipeStepProductID"`
		Name                 string           `json:"name"`
		ID                   string           `json:"id"`
		Notes                string           `json:"notes"`
		BelongsToRecipeStep  string           `json:"belongsToRecipeStep"`
		VesselPredicate      string           `json:"vesselPredicate"`
		MaximumQuantity      uint32           `json:"maximumQuantity"`
		MinimumQuantity      uint32           `json:"minimumQuantity"`
		UnavailableAfterStep bool             `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselCreationRequestInput struct {
		_                               struct{}
		RecipeStepProductID             *string `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64 `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64 `json:"productOfRecipeStepProductIndex"`
		InstrumentID                    *string `json:"instrumentID"`
		Name                            string  `json:"name"`
		Notes                           string  `json:"notes"`
		VesselPredicate                 string  `json:"vesselPredicate"`
		MinimumQuantity                 uint32  `json:"minimumQuantity"`
		MaximumQuantity                 uint32  `json:"maximumQuantity"`
		UnavailableAfterStep            bool    `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselDatabaseCreationInput struct {
		_                               struct{}
		InstrumentID                    *string
		RecipeStepProductID             *string
		ProductOfRecipeStepIndex        *uint64
		ProductOfRecipeStepProductIndex *uint64
		Name                            string
		ID                              string
		Notes                           string
		BelongsToRecipeStep             string
		VesselPredicate                 string
		MaximumQuantity                 uint32
		MinimumQuantity                 uint32
		UnavailableAfterStep            bool
	}

	// RecipeStepVesselUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepVesselUpdateRequestInput struct {
		_ struct{}

		RecipeStepProductID  *string `json:"recipeStepProductID"`
		Name                 *string `json:"name"`
		Notes                *string `json:"notes"`
		BelongsToRecipeStep  *string `json:"belongsToRecipeStep"`
		InstrumentID         *string `json:"instrumentID"`
		MinimumQuantity      *uint32 `json:"minimumQuantity"`
		MaximumQuantity      *uint32 `json:"maximumQuantity"`
		VesselPredicate      *string `json:"vesselPredicate"`
		UnavailableAfterStep *bool   `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepVesselDataManager interface {
		RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepVessel, error)
		GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepVessel], error)
		CreateRecipeStepVessel(ctx context.Context, input *RecipeStepVesselDatabaseCreationInput) (*RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, updated *RecipeStepVessel) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}
)

// Update merges an RecipeStepVesselUpdateRequestInput with a recipe step instrument.
func (x *RecipeStepVessel) Update(input *RecipeStepVesselUpdateRequestInput) {
	if input.RecipeStepProductID != nil && input.RecipeStepProductID != x.RecipeStepProductID {
		x.RecipeStepProductID = input.RecipeStepProductID
	}
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}
	if input.BelongsToRecipeStep != nil && *input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = *input.BelongsToRecipeStep
	}
	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}
	if input.MaximumQuantity != nil && *input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = *input.MaximumQuantity
	}
	if input.VesselPredicate != nil && *input.VesselPredicate != x.VesselPredicate {
		x.VesselPredicate = *input.VesselPredicate
	}

	if input.UnavailableAfterStep != nil && *input.UnavailableAfterStep != x.UnavailableAfterStep {
		x.UnavailableAfterStep = *input.UnavailableAfterStep
	}

	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && *input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = *input.MaximumQuantity
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.InstrumentID != nil && x.Instrument != nil && *input.InstrumentID != x.Instrument.ID {
		x.Instrument = &ValidInstrument{ID: *input.InstrumentID}
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselCreationRequestInput.
func (x *RecipeStepVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MinimumQuantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepVesselDatabaseCreationInput.
func (x *RecipeStepVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselUpdateRequestInput.
func (x *RecipeStepVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
