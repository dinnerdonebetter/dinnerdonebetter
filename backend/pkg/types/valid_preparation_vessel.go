package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidPreparationVesselCreatedServiceEventType indicates a valid preparation instrument was created.
	ValidPreparationVesselCreatedServiceEventType = "valid_preparation_instrument_created"
	// ValidPreparationVesselUpdatedServiceEventType indicates a valid preparation instrument was updated.
	ValidPreparationVesselUpdatedServiceEventType = "valid_preparation_instrument_updated"
	// ValidPreparationVesselArchivedServiceEventType indicates a valid preparation instrument was archived.
	ValidPreparationVesselArchivedServiceEventType = "valid_preparation_instrument_archived"
)

func init() {
	gob.Register(new(ValidPreparationVessel))
	gob.Register(new(ValidPreparationVesselCreationRequestInput))
	gob.Register(new(ValidPreparationVesselUpdateRequestInput))
}

type (
	// ValidPreparationVessel represents a valid preparation instrument.
	ValidPreparationVessel struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time        `json:"createdAt"`
		LastUpdatedAt *time.Time       `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time       `json:"archivedAt"`
		ID            string           `json:"id"`
		Notes         string           `json:"notes"`
		Preparation   ValidPreparation `json:"preparation"`
		Vessel        ValidVessel      `json:"instrument"`
	}

	// ValidPreparationVesselCreationRequestInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationVesselCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidPreparationID string `json:"validPreparationID"`
		ValidVesselID      string `json:"validVesselID"`
	}

	// ValidPreparationVesselDatabaseCreationInput represents what a user could set as input for creating valid preparation instruments.
	ValidPreparationVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string `json:"-"`
		Notes              string `json:"-"`
		ValidPreparationID string `json:"-"`
		ValidVesselID      string `json:"-"`
	}

	// ValidPreparationVesselUpdateRequestInput represents what a user could set as input for updating valid preparation instruments.
	ValidPreparationVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes              *string `json:"notes,omitempty"`
		ValidPreparationID *string `json:"validPreparationID,omitempty"`
		ValidVesselID      *string `json:"validVesselID,omitempty"`
	}

	// ValidPreparationVesselDataManager describes a structure capable of storing valid preparation instruments permanently.
	ValidPreparationVesselDataManager interface {
		ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (bool, error)
		GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*ValidPreparationVessel, error)
		GetValidPreparationVessels(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		GetValidPreparationVesselsForVessel(ctx context.Context, instrumentID string, filter *QueryFilter) (*QueryFilteredResult[ValidPreparationVessel], error)
		CreateValidPreparationVessel(ctx context.Context, input *ValidPreparationVesselDatabaseCreationInput) (*ValidPreparationVessel, error)
		UpdateValidPreparationVessel(ctx context.Context, updated *ValidPreparationVessel) error
		ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error
	}

	// ValidPreparationVesselDataService describes a structure capable of serving traffic related to valid preparation instruments.
	ValidPreparationVesselDataService interface {
		ListValidPreparationVesselsHandler(http.ResponseWriter, *http.Request)
		CreateValidPreparationVesselHandler(http.ResponseWriter, *http.Request)
		ReadValidPreparationVesselHandler(http.ResponseWriter, *http.Request)
		UpdateValidPreparationVesselHandler(http.ResponseWriter, *http.Request)
		ArchiveValidPreparationVesselHandler(http.ResponseWriter, *http.Request)
		SearchValidPreparationVesselsByPreparationHandler(http.ResponseWriter, *http.Request)
		SearchValidPreparationVesselsByVesselHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ValidPreparationVesselUpdateRequestInput with a valid preparation instrument.
func (x *ValidPreparationVessel) Update(input *ValidPreparationVesselUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidPreparationID != nil && *input.ValidPreparationID != x.Preparation.ID {
		x.Preparation.ID = *input.ValidPreparationID
	}

	if input.ValidVesselID != nil && *input.ValidVesselID != x.Vessel.ID {
		x.Vessel.ID = *input.ValidVesselID
	}
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselCreationRequestInput.
func (x *ValidPreparationVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselDatabaseCreationInput.
func (x *ValidPreparationVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPreparationVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPreparationVesselUpdateRequestInput.
func (x *ValidPreparationVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidVesselID, validation.Required),
	)
}
