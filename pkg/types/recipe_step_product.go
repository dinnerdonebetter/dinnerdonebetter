package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

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
		CreatedAt                          time.Time            `json:"createdAt"`
		MinimumStorageTemperatureInCelsius *float32             `json:"minimumStorageTemperatureInCelsius"`
		LastUpdatedAt                      *time.Time           `json:"lastUpdatedAt"`
		MaximumStorageTemperatureInCelsius *float32             `json:"maximumStorageTemperatureInCelsius"`
		ArchivedAt                         *time.Time           `json:"archivedAt"`
		StorageInstructions                string               `json:"storageInstructions"`
		Name                               string               `json:"name"`
		Type                               string               `json:"type"`
		BelongsToRecipeStep                string               `json:"belongsToRecipeStep"`
		ID                                 string               `json:"id"`
		QuantityNotes                      string               `json:"quantityNotes"`
		MeasurementUnit                    ValidMeasurementUnit `json:"measurementUnit"`
		MaximumStorageDurationInSeconds    uint32               `json:"maximumStorageDurationInSeconds"`
		MaximumQuantity                    float32              `json:"maximumQuantity"`
		MinimumQuantity                    float32              `json:"minimumQuantity"`
		Compostable                        bool                 `json:"compostable"`
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
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius"`
		StorageInstructions                string   `json:"storageInstructions"`
		QuantityNotes                      string   `json:"quantityNotes"`
		MeasurementUnitID                  string   `json:"measurementUnitID"`
		BelongsToRecipeStep                string   `json:"-"`
		Type                               string   `json:"type"`
		Name                               string   `json:"name"`
		ID                                 string   `json:"-"`
		MaximumStorageDurationInSeconds    uint32   `json:"maximumStorageDurationInSeconds"`
		MaximumQuantity                    float32  `json:"maximumQuantity"`
		MinimumQuantity                    float32  `json:"minimumQuantity"`
		Compostable                        bool     `json:"compostable"`
	}

	// RecipeStepProductDatabaseCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductDatabaseCreationInput struct {
		_                                  struct{}
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius"`
		StorageInstructions                string   `json:"storageInstructions"`
		MeasurementUnitID                  string   `json:"measurementUnitID"`
		QuantityNotes                      string   `json:"quantityNotes"`
		BelongsToRecipeStep                string   `json:"belongsToRecipeStep"`
		Type                               string   `json:"type"`
		Name                               string   `json:"name"`
		ID                                 string   `json:"id"`
		MaximumStorageDurationInSeconds    uint32   `json:"maximumStorageDurationInSeconds"`
		MaximumQuantity                    float32  `json:"maximumQuantity"`
		MinimumQuantity                    float32  `json:"minimumQuantity"`
		Compostable                        bool     `json:"compostable"`
	}

	// RecipeStepProductUpdateRequestInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateRequestInput struct {
		_                                  struct{}
		Name                               *string  `json:"name"`
		Type                               *string  `json:"type"`
		MeasurementUnitID                  *string  `json:"measurementUnitID"`
		QuantityNotes                      *string  `json:"quantityNotes"`
		BelongsToRecipeStep                *string  `json:"belongsToRecipeStep"`
		MinimumQuantity                    *float32 `json:"minimumQuantity"`
		MaximumQuantity                    *float32 `json:"maximumQuantity"`
		Compostable                        *bool    `json:"compostable"`
		MaximumStorageDurationInSeconds    *uint32  `json:"maximumStorageDurationInSeconds"`
		MinimumStorageTemperatureInCelsius *float32 `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius *float32 `json:"maximumStorageTemperatureInCelsius"`
		StorageInstructions                *string  `json:"storageInstructions"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*RecipeStepProduct, error)
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepProductList, error)
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

	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && *input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = *input.MaximumQuantity
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

	if input.MinimumStorageTemperatureInCelsius != nil && input.MinimumStorageTemperatureInCelsius != x.MinimumStorageTemperatureInCelsius {
		x.MinimumStorageTemperatureInCelsius = input.MinimumStorageTemperatureInCelsius
	}

	if input.MaximumStorageTemperatureInCelsius != nil && input.MaximumStorageTemperatureInCelsius != x.MaximumStorageTemperatureInCelsius {
		x.MaximumStorageTemperatureInCelsius = input.MaximumStorageTemperatureInCelsius
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
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
		validation.Field(&x.MinimumQuantity, validation.Required),
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
		validation.Field(&x.MinimumQuantity, validation.Required),
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
		MinimumQuantity:                    &input.MinimumQuantity,
		MaximumQuantity:                    &input.MaximumQuantity,
		Compostable:                        &input.Compostable,
		MaximumStorageDurationInSeconds:    &input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                &input.StorageInstructions,
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
		MinimumQuantity:                    input.MinimumQuantity,
		MaximumQuantity:                    input.MaximumQuantity,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                input.StorageInstructions,
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
		validation.Field(&x.MinimumQuantity, validation.Required),
		validation.Field(&x.MaximumQuantity, validation.Required),
	)
}
