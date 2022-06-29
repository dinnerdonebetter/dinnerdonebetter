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
		_                   struct{}
		RecipeStepProductID *string `json:"recipeStepProductID"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		IngredientID        *string `json:"ingredientID"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		ID                  string  `json:"id"`
		CreatedOn           uint64  `json:"createdOn"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
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
		QuantityType        string                                 `json:"quantityType"`
		QuantityNotes       string                                 `json:"quantityNotes"`
		IngredientNotes     string                                 `json:"ingredientNotes"`
		QuantityValue       float32                                `json:"quantityValue"`
		ProductOfRecipeStep bool                                   `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_                   struct{}
		IngredientID        *string                                 `json:"ingredientID"`
		RecipeStepProduct   *RecipeStepProductDatabaseCreationInput `json:"recipeStepProduct"`
		ID                  string                                  `json:"id"`
		Name                string                                  `json:"name"`
		QuantityType        string                                  `json:"quantityType"`
		QuantityNotes       string                                  `json:"quantityNotes"`
		IngredientNotes     string                                  `json:"ingredientNotes"`
		BelongsToRecipeStep string                                  `json:"belongsToRecipeStep"`
		QuantityValue       float32                                 `json:"quantityValue"`
		ProductOfRecipeStep bool                                    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_                   struct{}
		IngredientID        *string `json:"ingredientID"`
		RecipeStepProductID *string `json:"recipeStepProductID"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
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

	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.QuantityType != "" && input.QuantityType != x.QuantityType {
		x.QuantityType = input.QuantityType
	}

	if input.QuantityValue != 0 && input.QuantityValue != x.QuantityValue {
		x.QuantityValue = input.QuantityValue
	}

	if input.QuantityNotes != "" && input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = input.QuantityNotes
	}

	if input.ProductOfRecipeStep != x.ProductOfRecipeStep {
		x.ProductOfRecipeStep = input.ProductOfRecipeStep
	}

	if input.IngredientNotes != "" && input.IngredientNotes != x.IngredientNotes {
		x.IngredientNotes = input.IngredientNotes
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientCreationRequestInput.
func (x *RecipeStepIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientDatabaseCreationInput.
func (x *RecipeStepIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
	)
}

// RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(input *RecipeStepIngredientCreationRequestInput) *RecipeStepIngredientDatabaseCreationInput {
	x := &RecipeStepIngredientDatabaseCreationInput{
		IngredientID:        input.IngredientID,
		Name:                input.Name,
		RecipeStepProduct:   RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(input.RecipeStepProduct),
		QuantityType:        input.QuantityType,
		QuantityValue:       input.QuantityValue,
		QuantityNotes:       input.QuantityNotes,
		ProductOfRecipeStep: input.ProductOfRecipeStep,
		IngredientNotes:     input.IngredientNotes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientUpdateRequestInput.
func (x *RecipeStepIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
	)
}
