package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientMeasurementUnitDataType indicates an event is related to a valid ingredient measurement unit.
	ValidIngredientMeasurementUnitDataType dataType = "valid_ingredient_preparation"

	// ValidIngredientMeasurementUnitCreatedCustomerEventType indicates a valid ingredient measurement unit was created.
	ValidIngredientMeasurementUnitCreatedCustomerEventType CustomerEventType = "valid_ingredient_preparation_created"
	// ValidIngredientMeasurementUnitUpdatedCustomerEventType indicates a valid ingredient measurement unit was updated.
	ValidIngredientMeasurementUnitUpdatedCustomerEventType CustomerEventType = "valid_ingredient_preparation_updated"
	// ValidIngredientMeasurementUnitArchivedCustomerEventType indicates a valid ingredient measurement unit was archived.
	ValidIngredientMeasurementUnitArchivedCustomerEventType CustomerEventType = "valid_ingredient_preparation_archived"
)

func init() {
	gob.Register(new(ValidIngredientMeasurementUnit))
	gob.Register(new(ValidIngredientMeasurementUnitList))
	gob.Register(new(ValidIngredientMeasurementUnitCreationRequestInput))
	gob.Register(new(ValidIngredientMeasurementUnitUpdateRequestInput))
}

type (
	// ValidIngredientMeasurementUnit represents a valid ingredient measurement unit.
	ValidIngredientMeasurementUnit struct {
		_                      struct{}
		ArchivedOn             *uint64 `json:"archivedOn"`
		LastUpdatedOn          *uint64 `json:"lastUpdatedOn"`
		Notes                  string  `json:"notes"`
		ValidMeasurementUnitID string  `json:"validPreparationID"`
		ValidIngredientID      string  `json:"validIngredientID"`
		ID                     string  `json:"id"`
		CreatedOn              uint64  `json:"createdOn"`
	}

	// ValidIngredientMeasurementUnitList represents a list of valid ingredient measurement units.
	ValidIngredientMeasurementUnitList struct {
		_                               struct{}
		ValidIngredientMeasurementUnits []*ValidIngredientMeasurementUnit `json:"data"`
		Pagination
	}

	// ValidIngredientMeasurementUnitCreationRequestInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitCreationRequestInput struct {
		_                      struct{}
		ID                     string `json:"-"`
		Notes                  string `json:"notes"`
		ValidMeasurementUnitID string `json:"validPreparationID"`
		ValidIngredientID      string `json:"validIngredientID"`
	}

	// ValidIngredientMeasurementUnitDatabaseCreationInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitDatabaseCreationInput struct {
		_ struct{}

		ID                     string `json:"id"`
		Notes                  string `json:"notes"`
		ValidMeasurementUnitID string `json:"validPreparationID"`
		ValidIngredientID      string `json:"validIngredientID"`
	}

	// ValidIngredientMeasurementUnitUpdateRequestInput represents what a user could set as input for updating valid ingredient measurement units.
	ValidIngredientMeasurementUnitUpdateRequestInput struct {
		_ struct{}

		Notes                  *string `json:"notes"`
		ValidMeasurementUnitID *string `json:"validPreparationID"`
		ValidIngredientID      *string `json:"validIngredientID"`
	}

	// ValidIngredientMeasurementUnitDataManager describes a structure capable of storing valid ingredient measurement units permanently.
	ValidIngredientMeasurementUnitDataManager interface {
		ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error)
		GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*ValidIngredientMeasurementUnit, error)
		GetValidIngredientMeasurementUnits(ctx context.Context, filter *QueryFilter) (*ValidIngredientMeasurementUnitList, error)
		GetValidIngredientMeasurementUnitsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*ValidIngredientMeasurementUnit, error)
		CreateValidIngredientMeasurementUnit(ctx context.Context, input *ValidIngredientMeasurementUnitDatabaseCreationInput) (*ValidIngredientMeasurementUnit, error)
		UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *ValidIngredientMeasurementUnit) error
		ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error
	}

	// ValidIngredientMeasurementUnitDataService describes a structure capable of serving traffic related to valid ingredient measurement units.
	ValidIngredientMeasurementUnitDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientMeasurementUnitUpdateRequestInput with a valid ingredient measurement unit.
func (x *ValidIngredientMeasurementUnit) Update(input *ValidIngredientMeasurementUnitUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidMeasurementUnitID != nil && *input.ValidMeasurementUnitID != x.ValidMeasurementUnitID {
		x.ValidMeasurementUnitID = *input.ValidMeasurementUnitID
	}

	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.ValidIngredientID {
		x.ValidIngredientID = *input.ValidIngredientID
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitCreationRequestInput.
func (x *ValidIngredientMeasurementUnitCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitDatabaseCreationInput.
func (x *ValidIngredientMeasurementUnitDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}

// ValidIngredientMeasurementUnitFromValidIngredientMeasurementUnit creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientMeasurementUnitFromValidIngredientMeasurementUnit(input *ValidIngredientMeasurementUnit) *ValidIngredientMeasurementUnitUpdateRequestInput {
	x := &ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                  &input.Notes,
		ValidMeasurementUnitID: &input.ValidMeasurementUnitID,
		ValidIngredientID:      &input.ValidIngredientID,
	}

	return x
}

// ValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnitCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnitCreationInput(input *ValidIngredientMeasurementUnitCreationRequestInput) *ValidIngredientMeasurementUnitDatabaseCreationInput {
	x := &ValidIngredientMeasurementUnitDatabaseCreationInput{
		Notes:                  input.Notes,
		ValidMeasurementUnitID: input.ValidMeasurementUnitID,
		ValidIngredientID:      input.ValidIngredientID,
	}

	return x
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitUpdateRequestInput.
func (x *ValidIngredientMeasurementUnitUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}