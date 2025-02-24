package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RecipeStepVessel represents a recipe step instrument.
	RecipeStepVessel struct {
		_ struct{} `json:"-"`

		CreatedAt            time.Time                  `json:"createdAt"`
		Quantity             Uint16RangeWithOptionalMax `json:"quantity"`
		LastUpdatedAt        *time.Time                 `json:"lastUpdatedAt"`
		ArchivedAt           *time.Time                 `json:"archivedAt"`
		RecipeStepProductID  *string                    `json:"recipeStepProductID"`
		Vessel               *ValidVessel               `json:"vessel"`
		ID                   string                     `json:"id"`
		Notes                string                     `json:"notes"`
		BelongsToRecipeStep  string                     `json:"belongsToRecipeStep"`
		VesselPreposition    string                     `json:"vesselPreposition"`
		Name                 string                     `json:"name"`
		UnavailableAfterStep bool                       `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselCreationRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepProductID             *string                    `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64                    `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64                    `json:"productOfRecipeStepProductIndex"`
		VesselID                        *string                    `json:"vesselID"`
		Quantity                        Uint16RangeWithOptionalMax `json:"quantity"`
		Name                            string                     `json:"name"`
		Notes                           string                     `json:"notes"`
		VesselPreposition               string                     `json:"vesselPreposition"`
		UnavailableAfterStep            bool                       `json:"unavailableAfterStep"`
	}

	// RecipeStepVesselDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		VesselID                        *string                    `json:"-"`
		RecipeStepProductID             *string                    `json:"-"`
		ProductOfRecipeStepIndex        *uint64                    `json:"-"`
		ProductOfRecipeStepProductIndex *uint64                    `json:"-"`
		Quantity                        Uint16RangeWithOptionalMax `json:"-"`
		ID                              string                     `json:"-"`
		Notes                           string                     `json:"-"`
		BelongsToRecipeStep             string                     `json:"-"`
		VesselPreposition               string                     `json:"-"`
		Name                            string                     `json:"-"`
		UnavailableAfterStep            bool                       `json:"-"`
	}

	// RecipeStepVesselUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepProductID  *string                                      `json:"recipeStepProductID,omitempty"`
		Name                 *string                                      `json:"name,omitempty"`
		Notes                *string                                      `json:"notes,omitempty"`
		BelongsToRecipeStep  *string                                      `json:"belongsToRecipeStep,omitempty"`
		VesselID             *string                                      `json:"vesselID,omitempty"`
		Quantity             Uint16RangeWithOptionalMaxUpdateRequestInput `json:"quantity"`
		VesselPreposition    *string                                      `json:"vesselPreposition,omitempty"`
		UnavailableAfterStep *bool                                        `json:"unavailableAfterStep,omitempty"`
	}

	// RecipeStepVesselDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepVesselDataManager interface {
		RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepVessel, error)
		GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipeStepVessel], error)
		CreateRecipeStepVessel(ctx context.Context, input *RecipeStepVesselDatabaseCreationInput) (*RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, updated *RecipeStepVessel) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}

	// RecipeStepVesselDataService describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepVesselDataService interface {
		ListRecipeStepVesselsHandler(http.ResponseWriter, *http.Request)
		CreateRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		ReadRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		UpdateRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
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

	if input.Quantity.Min != nil && *input.Quantity.Min != x.Quantity.Min {
		x.Quantity.Min = *input.Quantity.Min
	}

	if input.Quantity.Max != nil && x.Quantity.Max != nil && *input.Quantity.Max != *x.Quantity.Max {
		x.Quantity.Max = input.Quantity.Max
	}

	if input.VesselPreposition != nil && *input.VesselPreposition != x.VesselPreposition {
		x.VesselPreposition = *input.VesselPreposition
	}

	if input.UnavailableAfterStep != nil && *input.UnavailableAfterStep != x.UnavailableAfterStep {
		x.UnavailableAfterStep = *input.UnavailableAfterStep
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.VesselID != nil && x.Vessel != nil && *input.VesselID != x.Vessel.ID {
		x.Vessel = &ValidVessel{ID: *input.VesselID}
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselCreationRequestInput.
func (x *RecipeStepVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
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
