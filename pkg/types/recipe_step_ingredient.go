package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

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
		_                   struct{}
		CreatedAt           time.Time            `json:"createdAt"`
		RecipeStepProductID *string              `json:"recipeStepProductID"`
		LastUpdatedAt       *time.Time           `json:"lastUpdatedAt"`
		ArchivedAt          *time.Time           `json:"archivedAt"`
		Ingredient          *ValidIngredient     `json:"ingredient"`
		QuantityNotes       string               `json:"quantityNotes"`
		Name                string               `json:"name"`
		ID                  string               `json:"id"`
		IngredientNotes     string               `json:"ingredientNotes"`
		BelongsToRecipeStep string               `json:"belongsToRecipeStep"`
		MeasurementUnit     ValidMeasurementUnit `json:"measurementUnit"`
		MinimumQuantity     float32              `json:"minimumQuantity"`
		MaximumQuantity     float32              `json:"maximumQuantity"`
		Optional            bool                 `json:"optional"`
		ProductOfRecipeStep bool                 `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		_                     struct{}
		RecipeStepIngredients []*RecipeStepIngredient `json:"data"`
		Pagination
	}

	// RecipeStepIngredientCreationRequestInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationRequestInput struct {
		_                   struct{}
		IngredientID        *string                                `json:"ingredientID"`
		RecipeStepProduct   *RecipeStepProductCreationRequestInput `json:"recipeStepProduct"`
		ID                  string                                 `json:"-"`
		BelongsToRecipeStep string                                 `json:"-"`
		Name                string                                 `json:"name"`
		MeasurementUnitID   string                                 `json:"measurementUnitID"`
		QuantityNotes       string                                 `json:"quantityNotes"`
		IngredientNotes     string                                 `json:"ingredientNotes"`
		MaximumQuantity     float32                                `json:"maximumQuantity"`
		MinimumQuantity     float32                                `json:"minimumQuantity"`
		Optional            bool                                   `json:"optional"`
		ProductOfRecipeStep bool                                   `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_                   struct{}
		IngredientID        *string `json:"ingredientID"`
		RecipeStepProductID *string `json:"recipeStepProductID"`
		MeasurementUnitID   string  `json:"measurementUnitID"`
		Name                string  `json:"name"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		ID                  string  `json:"id"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		MinimumQuantity     float32 `json:"minimumQuantity"`
		MaximumQuantity     float32 `json:"maximumQuantity"`
		Optional            bool    `json:"optional"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_ struct{}

		// IngredientID is already a pointer, and I don't feel like making it a double pointer.
		IngredientID *string `json:"ingredientID"`
		// RecipeStepProductID is already a pointer, and I don't feel like making it a double pointer.
		RecipeStepProductID *string  `json:"recipeStepProductID"`
		Name                *string  `json:"name"`
		Optional            *bool    `json:"optional"`
		MeasurementUnitID   *string  `json:"measurementUnitID"`
		QuantityNotes       *string  `json:"quantityNotes"`
		IngredientNotes     *string  `json:"ingredientNotes"`
		BelongsToRecipeStep *string  `json:"belongsToRecipeStep"`
		MinimumQuantity     *float32 `json:"minimumQuantity"`
		MaximumQuantity     *float32 `json:"maximumQuantity"`
		ProductOfRecipeStep *bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepIngredient, error)
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*RecipeStepIngredientList, error)
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
	if input.IngredientID != nil && (x.Ingredient == nil || (*input.IngredientID != "" && *input.IngredientID != x.Ingredient.ID)) {
		x.Ingredient.ID = *input.IngredientID
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

	if input.MinimumQuantity != nil && *input.MinimumQuantity != x.MinimumQuantity {
		x.MinimumQuantity = *input.MinimumQuantity
	}

	if input.MaximumQuantity != nil && *input.MaximumQuantity != x.MaximumQuantity {
		x.MaximumQuantity = *input.MaximumQuantity
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

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientCreationRequestInput.
func (x *RecipeStepIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MinimumQuantity, validation.Required),
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
		validation.Field(&x.MinimumQuantity, validation.Required),
	)
}

// RecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientUpdateRequestInputFromRecipeStepIngredient(input *RecipeStepIngredient) *RecipeStepIngredientUpdateRequestInput {
	x := &RecipeStepIngredientUpdateRequestInput{
		IngredientID:        &input.Ingredient.ID,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                &input.Name,
		MeasurementUnitID:   &input.MeasurementUnit.ID,
		QuantityNotes:       &input.QuantityNotes,
		IngredientNotes:     &input.IngredientNotes,
		BelongsToRecipeStep: &input.BelongsToRecipeStep,
		MinimumQuantity:     &input.MinimumQuantity,
		MaximumQuantity:     &input.MaximumQuantity,
		ProductOfRecipeStep: &input.ProductOfRecipeStep,
		Optional:            &input.Optional,
	}

	return x
}

// RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(input *RecipeStepIngredientCreationRequestInput) *RecipeStepIngredientDatabaseCreationInput {
	x := &RecipeStepIngredientDatabaseCreationInput{
		IngredientID:        input.IngredientID,
		Name:                input.Name,
		MeasurementUnitID:   input.MeasurementUnitID,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Optional:            input.Optional,
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
		validation.Field(&x.MinimumQuantity, validation.Required),
	)
}
