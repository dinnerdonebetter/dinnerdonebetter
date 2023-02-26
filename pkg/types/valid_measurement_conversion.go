package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidMeasurementConversionDataType indicates an event is related to a valid measurement conversion.
	ValidMeasurementConversionDataType dataType = "valid_measurement_conversion"

	// ValidMeasurementConversionCreatedCustomerEventType indicates a valid measurement conversion was created.
	ValidMeasurementConversionCreatedCustomerEventType CustomerEventType = "valid_measurement_conversion_created"
	// ValidMeasurementConversionUpdatedCustomerEventType indicates a valid measurement conversion was updated.
	ValidMeasurementConversionUpdatedCustomerEventType CustomerEventType = "valid_measurement_conversion_updated"
	// ValidMeasurementConversionArchivedCustomerEventType indicates a valid measurement conversion was archived.
	ValidMeasurementConversionArchivedCustomerEventType CustomerEventType = "valid_measurement_conversion_archived"
)

func init() {
	gob.Register(new(ValidMeasurementUnitConversion))
	gob.Register(new(ValidMeasurementUnitConversionCreationRequestInput))
	gob.Register(new(ValidMeasurementUnitConversionUpdateRequestInput))
}

type (
	// ValidMeasurementUnitConversion represents a valid measurement conversion.
	ValidMeasurementUnitConversion struct {
		_ struct{}

		CreatedAt         time.Time            `json:"createdAt"`
		LastUpdatedAt     *time.Time           `json:"lastUpdatedAt"`
		ArchivedAt        *time.Time           `json:"archivedAt"`
		OnlyForIngredient *ValidIngredient     `json:"onlyForIngredient"`
		Notes             string               `json:"notes"`
		ID                string               `json:"id"`
		From              ValidMeasurementUnit `json:"from"`
		To                ValidMeasurementUnit `json:"to"`
		Modifier          float32              `json:"modifier"`
	}

	// ValidMeasurementUnitConversionCreationRequestInput represents what a user could set as input for creating valid measurement conversions.
	ValidMeasurementUnitConversionCreationRequestInput struct {
		_ struct{}

		OnlyForIngredient *string `json:"onlyForIngredient"`
		From              string  `json:"from"`
		To                string  `json:"to"`
		Notes             string  `json:"notes"`
		Modifier          float32 `json:"modifier"`
	}

	// ValidMeasurementConversionDatabaseCreationInput represents what a user could set as input for creating valid measurement conversions.
	ValidMeasurementConversionDatabaseCreationInput struct {
		_ struct{}

		OnlyForIngredient *string
		ID                string
		From              string
		To                string
		Notes             string
		Modifier          float32
	}

	// ValidMeasurementUnitConversionUpdateRequestInput represents what a user could set as input for updating valid measurement conversions.
	ValidMeasurementUnitConversionUpdateRequestInput struct {
		_ struct{}

		From              *string  `json:"from,omitempty"`
		To                *string  `json:"to,omitempty"`
		OnlyForIngredient *string  `json:"onlyForIngredient,omitempty"`
		Modifier          *float32 `json:"modifier,omitempty"`
		Notes             *string  `json:"notes,omitempty"`
	}

	// ValidMeasurementConversionDataManager describes a structure capable of storing valid measurement conversions permanently.
	ValidMeasurementConversionDataManager interface {
		ValidMeasurementConversionExists(ctx context.Context, validMeasurementConversionID string) (bool, error)
		GetValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) (*ValidMeasurementUnitConversion, error)
		CreateValidMeasurementConversion(ctx context.Context, input *ValidMeasurementConversionDatabaseCreationInput) (*ValidMeasurementUnitConversion, error)
		UpdateValidMeasurementConversion(ctx context.Context, updated *ValidMeasurementUnitConversion) error
		ArchiveValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) error
		GetValidMeasurementConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementUnitConversion, error)
		GetValidMeasurementConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementUnitConversion, error)
	}

	// ValidMeasurementConversionDataService describes a structure capable of serving traffic related to valid measurement conversions.
	ValidMeasurementConversionDataService interface {
		FromMeasurementUnitHandler(res http.ResponseWriter, req *http.Request)
		ToMeasurementUnitHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidMeasurementUnitConversionUpdateRequestInput with a valid measurement conversion.
func (x *ValidMeasurementUnitConversion) Update(input *ValidMeasurementUnitConversionUpdateRequestInput) {
	if input.From != nil && *input.From != x.From.ID {
		x.From.ID = *input.From
	}

	if input.To != nil && *input.To != x.To.ID {
		x.To.ID = *input.To
	}

	if input.OnlyForIngredient != nil && *input.OnlyForIngredient != x.OnlyForIngredient.ID {
		x.OnlyForIngredient.ID = *input.OnlyForIngredient
	}

	if input.Modifier != nil && *input.Modifier != x.Modifier {
		x.Modifier = *input.Modifier
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitConversionCreationRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitConversionCreationRequestInput.
func (x *ValidMeasurementUnitConversionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.From, validation.Required),
		validation.Field(&x.To, validation.Required),
		validation.Field(&x.Modifier, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidMeasurementConversionDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidMeasurementConversionDatabaseCreationInput.
func (x *ValidMeasurementConversionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.From, validation.Required),
		validation.Field(&x.To, validation.Required),
		validation.Field(&x.Modifier, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidMeasurementUnitConversionUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementUnitConversionUpdateRequestInput.
func (x *ValidMeasurementUnitConversionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.From, validation.Required),
		validation.Field(&x.To, validation.Required),
		validation.Field(&x.Modifier, validation.Required),
	)
}
