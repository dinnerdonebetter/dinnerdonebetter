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

	// InstrumentOwnershipCreationRequestInput represents what a user could set as input for creating household instrument ownerships.
	InstrumentOwnershipCreationRequestInput struct {
		_ struct{} `json:"-"`

		Notes              string `json:"notes"`
		ValidInstrumentID  string `json:"validInstrumentID"`
		BelongsToHousehold string `json:"belongsToHousehold"`
		Quantity           uint16 `json:"quantity"`
	}

	// InstrumentOwnershipDatabaseCreationInput represents what a user could set as input for creating household instrument ownerships.
	InstrumentOwnershipDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string `json:"-"`
		Notes              string `json:"-"`
		ValidInstrumentID  string `json:"-"`
		BelongsToHousehold string `json:"-"`
		Quantity           uint16 `json:"-"`
	}

	// InstrumentOwnershipUpdateRequestInput represents what a user could set as input for updating household instrument ownerships.
	InstrumentOwnershipUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes             *string `json:"notes"`
		Quantity          *uint16 `json:"quantity"`
		ValidInstrumentID *string `json:"validInstrumentID"`
	}

	// InstrumentOwnershipDataManager describes a structure capable of storing household instrument ownerships permanently.
	InstrumentOwnershipDataManager interface {
		InstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (bool, error)
		GetInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*InstrumentOwnership, error)
		GetInstrumentOwnerships(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[InstrumentOwnership], error)
		CreateInstrumentOwnership(ctx context.Context, input *InstrumentOwnershipDatabaseCreationInput) (*InstrumentOwnership, error)
		UpdateInstrumentOwnership(ctx context.Context, updated *InstrumentOwnership) error
		ArchiveInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error
	}

	// InstrumentOwnershipDataService describes a structure capable of serving traffic related to household instrument ownerships.
	InstrumentOwnershipDataService interface {
		ListInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		CreateInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		ReadInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		UpdateInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
		ArchiveInstrumentOwnershipHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an InstrumentOwnershipUpdateRequestInput with a household instrument ownership.
func (x *InstrumentOwnership) Update(input *InstrumentOwnershipUpdateRequestInput) {
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

var _ validation.ValidatableWithContext = (*InstrumentOwnershipCreationRequestInput)(nil)

// ValidateWithContext validates a InstrumentOwnershipCreationRequestInput.
func (x *InstrumentOwnershipCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*InstrumentOwnershipDatabaseCreationInput)(nil)

// ValidateWithContext validates a InstrumentOwnershipDatabaseCreationInput.
func (x *InstrumentOwnershipDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*InstrumentOwnershipUpdateRequestInput)(nil)

// ValidateWithContext validates a InstrumentOwnershipUpdateRequestInput.
func (x *InstrumentOwnershipUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
		validation.Field(&x.ValidInstrumentID, validation.Required),
	)
}
