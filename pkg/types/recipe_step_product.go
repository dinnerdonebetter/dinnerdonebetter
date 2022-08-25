package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepProductDataType indicates an event is related to a recipe step product.
	RecipeStepProductDataType dataType = "recipe_step_product"

	// RecipeStepProductIngredientType represents one of the valid recipe step product type values.
	RecipeStepProductIngredientType = "ingredient"
	// RecipeStepProductInstrumentType represents one of the valid recipe step product type values.
	RecipeStepProductInstrumentType = "instrument"

	// RecipeStepProductCreatedCustomerEventType indicates a recipe step product was created.
	RecipeStepProductCreatedCustomerEventType CustomerEventType = "recipe_step_product_created"
	// RecipeStepProductUpdatedCustomerEventType indicates a recipe step product was updated.
	RecipeStepProductUpdatedCustomerEventType CustomerEventType = "recipe_step_product_updated"
	// RecipeStepProductArchivedCustomerEventType indicates a recipe step product was archived.
	RecipeStepProductArchivedCustomerEventType CustomerEventType = "recipe_step_product_archived"
)

func init() {
	gob.Register(new(RecipeStepProduct))
	gob.Register(new(RecipeStepProductList))
	gob.Register(new(RecipeStepProductCreationRequestInput))
	gob.Register(new(RecipeStepProductUpdateRequestInput))
}

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		_                                  struct{}
		LastUpdatedOn                      *uint64              `json:"lastUpdatedOn"`
		ArchivedOn                         *uint64              `json:"archivedOn"`
		ID                                 string               `json:"id"`
		Name                               string               `json:"name"`
		Type                               string               `json:"type"`
		QuantityNotes                      string               `json:"quantityNotes"`
		BelongsToRecipeStep                string               `json:"belongsToRecipeStep"`
		MeasurementUnit                    ValidMeasurementUnit `json:"measurementUnit"`
		CreatedOn                          uint64               `json:"createdOn"`
		MinimumQuantityValue               float32              `json:"minimumQuantityValue"`
		MaximumQuantityValue               float32              `json:"maximumQuantityValue"`
		Compostable                        bool                 `json:"compostable"`
		MaximumStorageDurationInSeconds    uint32               `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius float32              `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius float32              `json:"maximumStorageTemperatureInCelsius"`
	}

	// RecipeStepProductList represents a list of recipe step products.
	RecipeStepProductList struct {
		_                  struct{}
		RecipeStepProducts []*RecipeStepProduct `json:"data"`
		Pagination
	}

	// RecipeStepProductCreationRequestInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationRequestInput struct {
		_                                  struct{}
		ID                                 string  `json:"-"`
		Name                               string  `json:"name"`
		Type                               string  `json:"type"`
		QuantityNotes                      string  `json:"quantityNotes"`
		MeasurementUnitID                  string  `json:"measurementUnitID"`
		BelongsToRecipeStep                string  `json:"-"`
		MinimumQuantityValue               float32 `json:"minimumQuantityValue"`
		MaximumQuantityValue               float32 `json:"maximumQuantityValue"`
		Compostable                        bool    `json:"compostable"`
		MaximumStorageDurationInSeconds    uint32  `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius float32 `json:"maximumStorageTemperatureInCelsius"`
	}

	// RecipeStepProductDatabaseCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductDatabaseCreationInput struct {
		_                                  struct{}
		ID                                 string  `json:"id"`
		Name                               string  `json:"name"`
		Type                               string  `json:"type"`
		MeasurementUnitID                  string  `json:"measurementUnitID"`
		QuantityNotes                      string  `json:"quantityNotes"`
		BelongsToRecipeStep                string  `json:"belongsToRecipeStep"`
		MinimumQuantityValue               float32 `json:"minimumQuantityValue"`
		MaximumQuantityValue               float32 `json:"maximumQuantityValue"`
		Compostable                        bool    `json:"compostable"`
		MaximumStorageDurationInSeconds    uint32  `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius float32 `json:"maximumStorageTemperatureInCelsius"`
	}

	// RecipeStepProductUpdateRequestInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateRequestInput struct {
		_                                  struct{}
		Name                               *string  `json:"name"`
		Type                               *string  `json:"type"`
		MeasurementUnitID                  *string  `json:"measurementUnitID"`
		QuantityNotes                      *string  `json:"quantityNotes"`
		BelongsToRecipeStep                *string  `json:"belongsToRecipeStep"`
		MinimumQuantityValue               *float32 `json:"minimumQuantityValue"`
		MaximumQuantityValue               *float32 `json:"maximumQuantityValue"`
		Compostable                        *bool    `json:"compostable"`
		MaximumStorageDurationInSeconds    *uint32  `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*RecipeStepProduct, error)
		GetTotalRecipeStepProductCount(ctx context.Context) (uint64, error)
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepProductList, error)
		GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*RecipeStepProduct, error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductDatabaseCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error
	}

	// RecipeStepProductDataService describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepProductUpdateRequestInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Type != nil && *input.Type != x.Type {
		x.Type = *input.Type
	}

	if input.MeasurementUnitID != nil && *input.MeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	if input.MinimumQuantityValue != nil && *input.MinimumQuantityValue != x.MinimumQuantityValue {
		x.MinimumQuantityValue = *input.MinimumQuantityValue
	}

	if input.MaximumQuantityValue != nil && *input.MaximumQuantityValue != x.MaximumQuantityValue {
		x.MaximumQuantityValue = *input.MaximumQuantityValue
	}

	if input.QuantityNotes != nil && *input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = *input.QuantityNotes
	}

	if input.Compostable != nil && *input.Compostable != x.Compostable {
		x.Compostable = *input.Compostable
	}

	if input.MaximumStorageDurationInSeconds != nil && *input.MaximumStorageDurationInSeconds != x.MaximumStorageDurationInSeconds {
		x.MaximumStorageDurationInSeconds = *input.MaximumStorageDurationInSeconds
	}

	if input.MinimumStorageTemperatureInCelsius != nil && *input.MinimumStorageTemperatureInCelsius != x.MinimumStorageTemperatureInCelsius {
		x.MinimumStorageTemperatureInCelsius = *input.MinimumStorageTemperatureInCelsius
	}

	if input.MaximumStorageTemperatureInCelsius != nil && *input.MaximumStorageTemperatureInCelsius != x.MaximumStorageTemperatureInCelsius {
		x.MaximumStorageTemperatureInCelsius = *input.MaximumStorageTemperatureInCelsius
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepProductCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductCreationRequestInput.
func (x *RecipeStepProductCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.Required),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepProductDatabaseCreationInput.
func (x *RecipeStepProductDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType)),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
		validation.Field(&x.MeasurementUnitID, validation.Required),
	)
}

// RecipeStepProductUpdateRequestInputFromRecipeStepProduct creates a DatabaseCreationInput from a CreationInput.
func RecipeStepProductUpdateRequestInputFromRecipeStepProduct(input *RecipeStepProduct) *RecipeStepProductUpdateRequestInput {
	if input == nil {
		return nil
	}

	x := &RecipeStepProductUpdateRequestInput{
		Name:                               &input.Name,
		Type:                               &input.Type,
		MeasurementUnitID:                  &input.MeasurementUnit.ID,
		QuantityNotes:                      &input.QuantityNotes,
		BelongsToRecipeStep:                &input.BelongsToRecipeStep,
		MinimumQuantityValue:               &input.MinimumQuantityValue,
		MaximumQuantityValue:               &input.MaximumQuantityValue,
		Compostable:                        &input.Compostable,
		MaximumStorageDurationInSeconds:    &input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: &input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: &input.MaximumStorageTemperatureInCelsius,
	}

	return x
}

// RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(input *RecipeStepProductCreationRequestInput) *RecipeStepProductDatabaseCreationInput {
	if input == nil {
		return nil
	}

	x := &RecipeStepProductDatabaseCreationInput{
		Name:                               input.Name,
		Type:                               input.Type,
		MeasurementUnitID:                  input.MeasurementUnitID,
		QuantityNotes:                      input.QuantityNotes,
		MinimumQuantityValue:               input.MinimumQuantityValue,
		MaximumQuantityValue:               input.MaximumQuantityValue,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepProductUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductUpdateRequestInput.
func (x *RecipeStepProductUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.Required),
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
		validation.Field(&x.MaximumQuantityValue, validation.Required),
	)
}
