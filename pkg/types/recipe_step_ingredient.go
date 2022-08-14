package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepIngredientDataType indicates an event is related to a recipe step ingredient.
	RecipeStepIngredientDataType dataType = "recipe_step_ingredient"

	// RecipeStepIngredientCreatedCustomerEventType indicates a recipe step ingredient was created.
	RecipeStepIngredientCreatedCustomerEventType CustomerEventType = "recipe_step_ingredient_created"
	// RecipeStepIngredientUpdatedCustomerEventType indicates a recipe step ingredient was updated.
	RecipeStepIngredientUpdatedCustomerEventType CustomerEventType = "recipe_step_ingredient_updated"
	// RecipeStepIngredientArchivedCustomerEventType indicates a recipe step ingredient was archived.
	RecipeStepIngredientArchivedCustomerEventType CustomerEventType = "recipe_step_ingredient_archived"
)

func init() {
	gob.Register(new(RecipeStepIngredient))
	gob.Register(new(RecipeStepIngredientList))
	gob.Register(new(RecipeStepIngredientCreationRequestInput))
	gob.Register(new(RecipeStepIngredientUpdateRequestInput))
}

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		_                    struct{}
		IngredientID         *string              `json:"ingredientID"`
		ArchivedOn           *uint64              `json:"archivedOn"`
		RecipeStepProductID  *string              `json:"recipeStepProductID"`
		LastUpdatedOn        *uint64              `json:"lastUpdatedOn"`
		QuantityNotes        string               `json:"quantityNotes"`
		BelongsToRecipeStep  string               `json:"belongsToRecipeStep"`
		ID                   string               `json:"id"`
		IngredientNotes      string               `json:"ingredientNotes"`
		Name                 string               `json:"name"`
		MeasurementUnit      ValidMeasurementUnit `json:"measurementUnit"`
		CreatedOn            uint64               `json:"createdOn"`
		MinimumQuantityValue float32              `json:"minimumQuantityValue"`
		MaximumQuantityValue float32              `json:"maximumQuantityValue"`
		ProductOfRecipeStep  bool                 `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		_                     struct{}
		RecipeStepIngredients []*RecipeStepIngredient `json:"data"`
		Pagination
	}

	// RecipeStepIngredientCreationRequestInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationRequestInput struct {
		_                    struct{}
		IngredientID         *string                                `json:"ingredientID"`
		RecipeStepProduct    *RecipeStepProductCreationRequestInput `json:"recipeStepProduct"`
		ID                   string                                 `json:"-"`
		BelongsToRecipeStep  string                                 `json:"-"`
		Name                 string                                 `json:"name"`
		MeasurementUnitID    string                                 `json:"measurementUnitID"`
		QuantityNotes        string                                 `json:"quantityNotes"`
		IngredientNotes      string                                 `json:"ingredientNotes"`
		MinimumQuantityValue float32                                `json:"minimumQuantityValue"`
		MaximumQuantityValue float32                                `json:"maximumQuantityValue"`
		ProductOfRecipeStep  bool                                   `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_                    struct{}
		IngredientID         *string `json:"ingredientID"`
		RecipeStepProductID  *string `json:"recipeStepProductID"`
		ID                   string  `json:"id"`
		Name                 string  `json:"name"`
		MeasurementUnitID    string  `json:"measurementUnitID"`
		QuantityNotes        string  `json:"quantityNotes"`
		IngredientNotes      string  `json:"ingredientNotes"`
		BelongsToRecipeStep  string  `json:"belongsToRecipeStep"`
		MinimumQuantityValue float32 `json:"minimumQuantityValue"`
		MaximumQuantityValue float32 `json:"maximumQuantityValue"`
		ProductOfRecipeStep  bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_ struct{}

		// IngredientID is already a pointer, and I don't feel like making it a double pointer.
		IngredientID *string `json:"ingredientID"`
		// RecipeStepProductID is already a pointer, and I don't feel like making it a double pointer.
		RecipeStepProductID  *string  `json:"recipeStepProductID"`
		Name                 *string  `json:"name"`
		MeasurementUnitID    *string  `json:"measurementUnitID"`
		QuantityNotes        *string  `json:"quantityNotes"`
		IngredientNotes      *string  `json:"ingredientNotes"`
		BelongsToRecipeStep  *string  `json:"belongsToRecipeStep"`
		MinimumQuantityValue *float32 `json:"minimumQuantityValue"`
		MaximumQuantityValue *float32 `json:"maximumQuantityValue"`
		ProductOfRecipeStep  *bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepIngredient, error)
		GetTotalRecipeStepIngredientCount(ctx context.Context) (uint64, error)
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepIngredientList, error)
		GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*RecipeStepIngredient, error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientDatabaseCreationInput) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error
	}

	// RecipeStepIngredientDataService describes a structure capable of serving traffic related to recipe step ingredients.
	RecipeStepIngredientDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepIngredientUpdateRequestInput with a recipe step ingredient.
func (x *RecipeStepIngredient) Update(input *RecipeStepIngredientUpdateRequestInput) {
	if input.IngredientID != nil && (x.IngredientID == nil || (*input.IngredientID != "" && *input.IngredientID != *x.IngredientID)) {
		x.IngredientID = input.IngredientID
	}

	if input.RecipeStepProductID != nil && (x.RecipeStepProductID == nil || (*input.RecipeStepProductID != "" && *input.RecipeStepProductID != *x.RecipeStepProductID)) {
		x.RecipeStepProductID = input.RecipeStepProductID
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
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

	if input.ProductOfRecipeStep != nil && *input.ProductOfRecipeStep != x.ProductOfRecipeStep {
		x.ProductOfRecipeStep = *input.ProductOfRecipeStep
	}

	if input.IngredientNotes != nil && *input.IngredientNotes != x.IngredientNotes {
		x.IngredientNotes = *input.IngredientNotes
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientCreationRequestInput.
func (x *RecipeStepIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientDatabaseCreationInput.
func (x *RecipeStepIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
	)
}

// RecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient(input *RecipeStepIngredient) *RecipeStepIngredientUpdateRequestInput {
	x := &RecipeStepIngredientUpdateRequestInput{
		IngredientID:         input.IngredientID,
		RecipeStepProductID:  input.RecipeStepProductID,
		Name:                 &input.Name,
		MeasurementUnitID:    &input.MeasurementUnit.ID,
		QuantityNotes:        &input.QuantityNotes,
		IngredientNotes:      &input.IngredientNotes,
		BelongsToRecipeStep:  &input.BelongsToRecipeStep,
		MinimumQuantityValue: &input.MinimumQuantityValue,
		MaximumQuantityValue: &input.MaximumQuantityValue,
		ProductOfRecipeStep:  &input.ProductOfRecipeStep,
	}

	return x
}

// RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(input *RecipeStepIngredientCreationRequestInput) *RecipeStepIngredientDatabaseCreationInput {
	x := &RecipeStepIngredientDatabaseCreationInput{
		IngredientID:         input.IngredientID,
		Name:                 input.Name,
		MeasurementUnitID:    input.MeasurementUnitID,
		MinimumQuantityValue: input.MinimumQuantityValue,
		MaximumQuantityValue: input.MaximumQuantityValue,
		QuantityNotes:        input.QuantityNotes,
		ProductOfRecipeStep:  input.ProductOfRecipeStep,
		IngredientNotes:      input.IngredientNotes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientUpdateRequestInput.
func (x *RecipeStepIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantityValue, validation.Required),
	)
}
