package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ValidMeasurementUnitConversionCreatedCustomerEventType indicates a valid measurement conversion was created.
	ValidMeasurementUnitConversionCreatedCustomerEventType CustomerEventType = "valid_measurement_unit_conversion_created"
	// ValidMeasurementUnitConversionUpdatedCustomerEventType indicates a valid measurement conversion was updated.
	ValidMeasurementUnitConversionUpdatedCustomerEventType CustomerEventType = "valid_measurement_unit_conversion_updated"
	// ValidMeasurementUnitConversionArchivedCustomerEventType indicates a valid measurement conversion was archived.
	ValidMeasurementUnitConversionArchivedCustomerEventType CustomerEventType = "valid_measurement_unit_conversion_archived"
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
		GetValidMeasurementUnitConversion(ctx context.Context, validMeasurementConversionID string) (*ValidMeasurementUnitConversion, error)
		CreateValidMeasurementConversion(ctx context.Context, input *ValidMeasurementConversionDatabaseCreationInput) (*ValidMeasurementUnitConversion, error)
		UpdateValidMeasurementConversion(ctx context.Context, updated *ValidMeasurementUnitConversion) error
		ArchiveValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) error
		GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementUnitConversion, error)
		GetValidMeasurementUnitConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementUnitConversion, error)
	}

	// ValidMeasurementConversionDataService describes a structure capable of serving traffic related to valid measurement conversions.
	ValidMeasurementConversionDataService interface {
		FromMeasurementUnitHandler(http.ResponseWriter, *http.Request)
		ToMeasurementUnitHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
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
