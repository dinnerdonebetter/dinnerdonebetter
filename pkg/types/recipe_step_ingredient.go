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
	gob.Register(new(RecipeStepIngredientCreationRequestInput))
	gob.Register(new(RecipeStepIngredientUpdateRequestInput))
}

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		_                   struct{}
		CreatedAt           time.Time            `json:"createdAt"`
		RecipeStepProductID *string              `json:"recipeStepProductID"`
		ArchivedAt          *time.Time           `json:"archivedAt"`
		Ingredient          *ValidIngredient     `json:"ingredient"`
		LastUpdatedAt       *time.Time           `json:"lastUpdatedAt"`
		MaximumQuantity     *float32             `json:"maximumQuantity"`
		VesselIndex         *uint16              `json:"vesselIndex"`
		Name                string               `json:"name"`
		QuantityNotes       string               `json:"quantityNotes"`
		IngredientNotes     string               `json:"ingredientNotes"`
		ID                  string               `json:"id"`
		BelongsToRecipeStep string               `json:"belongsToRecipeStep"`
		MeasurementUnit     ValidMeasurementUnit `json:"measurementUnit"`
		MinimumQuantity     float32              `json:"minimumQuantity"`
		OptionIndex         uint16               `json:"optionIndex"`
		RequiresDefrost     bool                 `json:"requiresDefrost"`
		Optional            bool                 `json:"optional"`
	}

	// RecipeStepIngredientCreationRequestInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationRequestInput struct {
		_                               struct{}
		IngredientID                    *string  `json:"ingredientID"`
		ProductOfRecipeStepIndex        *uint64  `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64  `json:"productOfRecipeStepProductIndex"`
		MaximumQuantity                 *float32 `json:"maximumQuantity"`
		VesselIndex                     *uint16  `json:"vesselIndex"`
		MeasurementUnitID               string   `json:"measurementUnitID"`
		IngredientNotes                 string   `json:"ingredientNotes"`
		QuantityNotes                   string   `json:"quantityNotes"`
		Name                            string   `json:"name"`
		MinimumQuantity                 float32  `json:"minimumQuantity"`
		OptionIndex                     uint16   `json:"optionIndex"`
		RequiresDefrost                 bool     `json:"requiresDefrost"`
		Optional                        bool     `json:"optional"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_                               struct{}
		IngredientID                    *string
		RecipeStepProductID             *string
		ProductOfRecipeStepIndex        *uint64
		ProductOfRecipeStepProductIndex *uint64
		MaximumQuantity                 *float32
		VesselIndex                     *uint16
		QuantityNotes                   string
		ID                              string
		IngredientNotes                 string
		MeasurementUnitID               string
		BelongsToRecipeStep             string
		Name                            string
		MinimumQuantity                 float32
		OptionIndex                     uint16
		Optional                        bool
		RequiresDefrost                 bool
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_ struct{}
		// IngredientID and RecipeStepProductID are already pointers, and I don't feel like making it a double pointer.
		IngredientID        *string  `json:"ingredientID"`
		RecipeStepProductID *string  `json:"recipeStepProductID"`
		Name                *string  `json:"name"`
		Optional            *bool    `json:"optional"`
		MeasurementUnitID   *string  `json:"measurementUnitID"`
		QuantityNotes       *string  `json:"quantityNotes"`
		IngredientNotes     *string  `json:"ingredientNotes"`
		BelongsToRecipeStep *string  `json:"belongsToRecipeStep"`
		MinimumQuantity     *float32 `json:"minimumQuantity"`
		MaximumQuantity     *float32 `json:"maximumQuantity"`
		OptionIndex         *uint16  `json:"optionIndex"`
		RequiresDefrost     *bool    `json:"requiresDefrost"`
		VesselIndex         *uint16  `json:"vesselIndex"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepIngredient, error)
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepIngredient], error)
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

	if input.MaximumQuantity != nil && x.MaximumQuantity != nil && *input.MaximumQuantity != *x.MaximumQuantity {
		x.MaximumQuantity = input.MaximumQuantity
	}

	if input.QuantityNotes != nil && *input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = *input.QuantityNotes
	}

	if input.IngredientNotes != nil && *input.IngredientNotes != x.IngredientNotes {
		x.IngredientNotes = *input.IngredientNotes
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.OptionIndex != nil && *input.OptionIndex != x.OptionIndex {
		x.OptionIndex = *input.OptionIndex
	}

	if input.RequiresDefrost != nil && *input.RequiresDefrost != x.RequiresDefrost {
		x.RequiresDefrost = *input.RequiresDefrost
	}

	if input.VesselIndex != nil && x.VesselIndex != nil && *input.VesselIndex != *x.VesselIndex {
		x.VesselIndex = input.VesselIndex
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
