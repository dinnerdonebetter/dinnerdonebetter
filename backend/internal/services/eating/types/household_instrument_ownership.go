package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// InstrumentOwnership represents a household instrument ownership.
	InstrumentOwnership struct {
		_ struct{} `json:"-"`

		CreatedAt          time.Time       `json:"createdAt"`
		ArchivedAt         *time.Time      `json:"archivedAt"`
		LastUpdatedAt      *time.Time      `json:"lastUpdatedAt"`
		ID                 string          `json:"id"`
		Notes              string          `json:"notes"`
		BelongsToHousehold string          `json:"belongsToHousehold"`
		Instrument         ValidInstrument `json:"instrument"`
		Quantity           uint16          `json:"quantity"`
	}

	// HouseholdInstrumentOwnershipCreationRequestInput represents what a user could set as input for creating household instrument ownerships.
	HouseholdInstrumentOwnershipCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidInstrumentID  string `json:"validInstrumentID"`
		BelongsToHousehold string `json:"belongsToHousehold"`
		Quantity           uint16 `json:"quantity"`
	}

	// HouseholdInstrumentOwnershipDatabaseCreationInput represents what a user could set as input for creating household instrument ownerships.
	HouseholdInstrumentOwnershipDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string `json:"-"`
		Notes              string `json:"-"`
		ValidInstrumentID  string `json:"-"`
		BelongsToHousehold string `json:"-"`
		Quantity           uint16 `json:"-"`
	}

	// HouseholdInstrumentOwnershipUpdateRequestInput represents what a user could set as input for updating household instrument ownerships.
	HouseholdInstrumentOwnershipUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes             *string `json:"notes"`
		Quantity          *uint16 `json:"quantity"`
		ValidInstrumentID *string `json:"validInstrumentID"`
	}

	// HouseholdInstrumentOwnershipDataManager describes a structure capable of storing household instrument ownerships permanently.
	HouseholdInstrumentOwnershipDataManager interface {
		HouseholdInstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (bool, error)
		GetHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*InstrumentOwnership, error)
		GetHouseholdInstrumentOwnerships(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[InstrumentOwnership], error)
		CreateHouseholdInstrumentOwnership(ctx context.Context, input *HouseholdInstrumentOwnershipDatabaseCreationInput) (*InstrumentOwnership, error)
		UpdateHouseholdInstrumentOwnership(ctx context.Context, updated *InstrumentOwnership) error
		ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error
	}

	// HouseholdInstrumentOwnershipDataService describes a structure capable of serving traffic related to household instrument ownerships.
	HouseholdInstrumentOwnershipDataService interface {
		ListHouseholdInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		CreateHouseholdInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		ReadHouseholdInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		UpdateHouseholdInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		ArchiveHouseholdInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an HouseholdInstrumentOwnershipUpdateRequestInput with a household instrument ownership.
func (x *InstrumentOwnership) Update(input *HouseholdInstrumentOwnershipUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.Quantity != nil && *input.Quantity != x.Quantity {
		x.Quantity = *input.Quantity
	}

	if input.ValidInstrumentID != nil && *input.ValidInstrumentID != x.Instrument.ID {
		x.Instrument = ValidInstrument{ID: *input.ValidInstrumentID}
	}
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipCreationRequestInput.
func (x *HouseholdInstrumentOwnershipCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipDatabaseCreationInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipDatabaseCreationInput.
func (x *HouseholdInstrumentOwnershipDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdInstrumentOwnershipUpdateRequestInput)(nil)

// ValidateWithContext validates a HouseholdInstrumentOwnershipUpdateRequestInput.
func (x *HouseholdInstrumentOwnershipUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}
