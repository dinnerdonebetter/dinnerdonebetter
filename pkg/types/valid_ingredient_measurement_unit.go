package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidIngredientMeasurementUnitDataType indicates an event is related to a valid ingredient measurement unit.
	ValidIngredientMeasurementUnitDataType dataType = "valid_ingredient_measurement_unit"

	// ValidIngredientMeasurementUnitCreatedCustomerEventType indicates a valid ingredient measurement unit was created.
	ValidIngredientMeasurementUnitCreatedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_created"
	// ValidIngredientMeasurementUnitUpdatedCustomerEventType indicates a valid ingredient measurement unit was updated.
	ValidIngredientMeasurementUnitUpdatedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_updated"
	// ValidIngredientMeasurementUnitArchivedCustomerEventType indicates a valid ingredient measurement unit was archived.
	ValidIngredientMeasurementUnitArchivedCustomerEventType CustomerEventType = "valid_ingredient_measurement_unit_archived"
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
		_                        struct{}
		CreatedAt                time.Time            `json:"createdAt"`
		LastUpdatedAt            *time.Time           `json:"lastUpdatedAt"`
		ArchivedAt               *time.Time           `json:"archivedAt"`
		MeasurementUnit          ValidMeasurementUnit `json:"measurementUnit"`
		Notes                    string               `json:"notes"`
		ID                       string               `json:"id"`
		Ingredient               ValidIngredient      `json:"ingredient"`
		MinimumAllowableQuantity float32              `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity float32              `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitList represents a list of valid ingredient measurement units.
	ValidIngredientMeasurementUnitList struct {
		_                               struct{}
		ValidIngredientMeasurementUnits []*ValidIngredientMeasurementUnit `json:"data"`
		Pagination
	}

	// ValidIngredientMeasurementUnitCreationRequestInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitCreationRequestInput struct {
		_                        struct{}
		ID                       string  `json:"-"`
		Notes                    string  `json:"notes"`
		ValidMeasurementUnitID   string  `json:"validMeasurementUnitID"`
		ValidIngredientID        string  `json:"validIngredientID"`
		MinimumAllowableQuantity float32 `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity float32 `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitDatabaseCreationInput represents what a user could set as input for creating valid ingredient measurement units.
	ValidIngredientMeasurementUnitDatabaseCreationInput struct {
		_ struct{}

		ID                       string  `json:"id"`
		Notes                    string  `json:"notes"`
		ValidMeasurementUnitID   string  `json:"validMeasurementUnitID"`
		ValidIngredientID        string  `json:"validIngredientID"`
		MinimumAllowableQuantity float32 `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity float32 `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitUpdateRequestInput represents what a user could set as input for updating valid ingredient measurement units.
	ValidIngredientMeasurementUnitUpdateRequestInput struct {
		_ struct{}

		Notes                    *string  `json:"notes"`
		ValidMeasurementUnitID   *string  `json:"validMeasurementUnitID"`
		ValidIngredientID        *string  `json:"validIngredientID"`
		MinimumAllowableQuantity *float32 `json:"minimumAllowableQuantity"`
		MaximumAllowableQuantity *float32 `json:"maximumAllowableQuantity"`
	}

	// ValidIngredientMeasurementUnitDataManager describes a structure capable of storing valid ingredient measurement units permanently.
	ValidIngredientMeasurementUnitDataManager interface {
		ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error)
		GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*ValidIngredientMeasurementUnit, error)
		GetValidIngredientMeasurementUnits(ctx context.Context, filter *QueryFilter) (*ValidIngredientMeasurementUnitList, error)
		GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *QueryFilter) (*ValidIngredientMeasurementUnitList, error)
		GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, ingredientID string, filter *QueryFilter) (*ValidIngredientMeasurementUnitList, error)
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
		SearchByIngredientHandler(res http.ResponseWriter, req *http.Request)
		SearchByMeasurementUnitHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientMeasurementUnitUpdateRequestInput with a valid ingredient measurement unit.
func (x *ValidIngredientMeasurementUnit) Update(input *ValidIngredientMeasurementUnitUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ValidMeasurementUnitID != nil && *input.ValidMeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit.ID = *input.ValidMeasurementUnitID
	}

	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.Ingredient.ID {
		x.Ingredient.ID = *input.ValidIngredientID
	}

	if input.MinimumAllowableQuantity != nil && *input.MinimumAllowableQuantity != x.MinimumAllowableQuantity {
		x.MinimumAllowableQuantity = *input.MinimumAllowableQuantity
	}

	if input.MaximumAllowableQuantity != nil && *input.MaximumAllowableQuantity != x.MaximumAllowableQuantity {
		x.MaximumAllowableQuantity = *input.MaximumAllowableQuantity
	}
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitCreationRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitCreationRequestInput.
func (x *ValidIngredientMeasurementUnitCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.MinimumAllowableQuantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitDatabaseCreationInput.
func (x *ValidIngredientMeasurementUnitDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.MinimumAllowableQuantity, validation.Required),
	)
}

// ValidIngredientMeasurementUnitFromValidIngredientMeasurementUnit creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientMeasurementUnitFromValidIngredientMeasurementUnit(input *ValidIngredientMeasurementUnit) *ValidIngredientMeasurementUnitUpdateRequestInput {
	x := &ValidIngredientMeasurementUnitUpdateRequestInput{
		Notes:                    &input.Notes,
		ValidMeasurementUnitID:   &input.MeasurementUnit.ID,
		ValidIngredientID:        &input.Ingredient.ID,
		MinimumAllowableQuantity: &input.MinimumAllowableQuantity,
		MaximumAllowableQuantity: &input.MaximumAllowableQuantity,
	}

	return x
}

// ValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnitCreationInput creates a DatabaseCreationInput from a CreationInput.
func ValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnitCreationInput(input *ValidIngredientMeasurementUnitCreationRequestInput) *ValidIngredientMeasurementUnitDatabaseCreationInput {
	x := &ValidIngredientMeasurementUnitDatabaseCreationInput{
		Notes:                    input.Notes,
		ValidMeasurementUnitID:   input.ValidMeasurementUnitID,
		ValidIngredientID:        input.ValidIngredientID,
		MinimumAllowableQuantity: input.MinimumAllowableQuantity,
		MaximumAllowableQuantity: input.MaximumAllowableQuantity,
	}

	return x
}

var _ validation.ValidatableWithContext = (*ValidIngredientMeasurementUnitUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidIngredientMeasurementUnitUpdateRequestInput.
func (x *ValidIngredientMeasurementUnitUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidMeasurementUnitID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}
