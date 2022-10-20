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

	// ValidMeasurementConversionQuantityModifier is what we multiply / divide by for float quantity values.
	ValidMeasurementConversionQuantityModifier = 10000

	// ValidMeasurementConversionCreatedCustomerEventType indicates a valid measurement conversion was created.
	ValidMeasurementConversionCreatedCustomerEventType CustomerEventType = "valid_measurement_conversion_created"
	// ValidMeasurementConversionUpdatedCustomerEventType indicates a valid measurement conversion was updated.
	ValidMeasurementConversionUpdatedCustomerEventType CustomerEventType = "valid_measurement_conversion_updated"
	// ValidMeasurementConversionArchivedCustomerEventType indicates a valid measurement conversion was archived.
	ValidMeasurementConversionArchivedCustomerEventType CustomerEventType = "valid_measurement_conversion_archived"
)

func init() {
	gob.Register(new(ValidMeasurementConversion))
	gob.Register(new(ValidMeasurementConversionList))
	gob.Register(new(ValidMeasurementConversionCreationRequestInput))
	gob.Register(new(ValidMeasurementConversionUpdateRequestInput))
}

type (
	// ValidMeasurementConversion represents a valid measurement conversion.
	ValidMeasurementConversion struct {
		_                 struct{}
		CreatedAt         time.Time            `json:"createdAt"`
		LastUpdatedAt     *time.Time           `json:"lastUpdatedAt"`
		ArchivedAt        *time.Time           `json:"archivedAt"`
		OnlyForIngredient *ValidIngredient     `json:"forIngredient"`
		Notes             string               `json:"notes"`
		ID                string               `json:"id"`
		From              ValidMeasurementUnit `json:"from"`
		To                ValidMeasurementUnit `json:"to"`
		Modifier          float32              `json:"modifier"`
	}

	// ValidMeasurementConversionList represents a list of valid measurement conversions.
	ValidMeasurementConversionList struct {
		_                           struct{}
		ValidMeasurementConversions []*ValidMeasurementConversion `json:"data"`
		Pagination
	}

	// ValidMeasurementConversionCreationRequestInput represents what a user could set as input for creating valid measurement conversions.
	ValidMeasurementConversionCreationRequestInput struct {
		_             struct{}
		ForIngredient *string `json:"forIngredient"`
		From          string  `json:"from"`
		To            string  `json:"to"`
		Notes         string  `json:"notes"`
		Modifier      float32 `json:"modifier"`
	}

	// ValidMeasurementConversionDatabaseCreationInput represents what a user could set as input for creating valid measurement conversions.
	ValidMeasurementConversionDatabaseCreationInput struct {
		_             struct{}
		ForIngredient *string
		ID            string
		From          string
		To            string
		Notes         string
		Modifier      float32
	}

	// ValidMeasurementConversionUpdateRequestInput represents what a user could set as input for updating valid measurement conversions.
	ValidMeasurementConversionUpdateRequestInput struct {
		_ struct{}

		From          *string  `json:"from"`
		To            *string  `json:"to"`
		ForIngredient *string  `json:"forIngredient"`
		Modifier      *float32 `json:"modifier"`
		Notes         *string  `json:"notes"`
	}

	// ValidMeasurementConversionDataManager describes a structure capable of storing valid measurement conversions permanently.
	ValidMeasurementConversionDataManager interface {
		ValidMeasurementConversionExists(ctx context.Context, validMeasurementConversionID string) (bool, error)
		GetValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) (*ValidMeasurementConversion, error)
		CreateValidMeasurementConversion(ctx context.Context, input *ValidMeasurementConversionDatabaseCreationInput) (*ValidMeasurementConversion, error)
		UpdateValidMeasurementConversion(ctx context.Context, updated *ValidMeasurementConversion) error
		ArchiveValidMeasurementConversion(ctx context.Context, validMeasurementConversionID string) error
		GetValidMeasurementConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementConversion, error)
		GetValidMeasurementConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*ValidMeasurementConversion, error)
	}

	// ValidMeasurementConversionDataService describes a structure capable of serving traffic related to valid measurement conversions.
	ValidMeasurementConversionDataService interface {
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidMeasurementConversionUpdateRequestInput with a valid measurement conversion.
func (x *ValidMeasurementConversion) Update(input *ValidMeasurementConversionUpdateRequestInput) {
	if input.From != nil && *input.From != x.From.ID {
		x.From.ID = *input.From
	}

	if input.To != nil && *input.To != x.To.ID {
		x.To.ID = *input.To
	}

	if input.ForIngredient != nil && *input.ForIngredient != x.OnlyForIngredient.ID {
		x.OnlyForIngredient.ID = *input.ForIngredient
	}

	if input.Modifier != nil && *input.Modifier != x.Modifier {
		x.Modifier = *input.Modifier
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*ValidMeasurementConversionCreationRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementConversionCreationRequestInput.
func (x *ValidMeasurementConversionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.From, validation.Required),
		validation.Field(&x.To, validation.Required),
		validation.Field(&x.Modifier, validation.Required),
		validation.Field(&x.Notes, validation.Required),
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
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidMeasurementConversionUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidMeasurementConversionUpdateRequestInput.
func (x *ValidMeasurementConversionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.From, validation.Required),
		validation.Field(&x.To, validation.Required),
		validation.Field(&x.Modifier, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
