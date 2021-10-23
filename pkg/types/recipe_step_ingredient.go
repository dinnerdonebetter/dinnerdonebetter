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
		IngredientID        *string `json:"ingredientID"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		ID                  string  `json:"id"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		QuantityType        string  `json:"quantityType"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		CreatedOn           uint64  `json:"createdOn"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// FullRecipeStepIngredient represents a recipe step ingredient.
	FullRecipeStepIngredient struct {
		_                   struct{}
		IngredientID        *string         `json:"ingredientID"`
		LastUpdatedOn       *uint64         `json:"lastUpdatedOn"`
		ArchivedOn          *uint64         `json:"archivedOn"`
		IngredientNotes     string          `json:"ingredientNotes"`
		QuantityNotes       string          `json:"quantityNotes"`
		ID                  string          `json:"id"`
		QuantityType        string          `json:"quantityType"`
		BelongsToRecipeStep string          `json:"belongsToRecipeStep"`
		Ingredient          ValidIngredient `json:"ingredient"`
		CreatedOn           uint64          `json:"createdOn"`
		QuantityValue       float32         `json:"quantityValue"`
		ProductOfRecipeStep bool            `json:"productOfRecipe"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		_                     struct{}
		RecipeStepIngredients []*RecipeStepIngredient `json:"recipeStepIngredients"`
		Pagination
	}

	// RecipeStepIngredientCreationRequestInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationRequestInput struct {
		_                   struct{}
		IngredientID        *string `json:"ingredientID"`
		ID                  string  `json:"-"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep string  `json:"-"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_                   struct{}
		IngredientID        *string `json:"ingredientID"`
		ID                  string  `json:"id"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_                   struct{}
		IngredientID        *string `json:"ingredientID"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep string  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
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

	if input.QuantityType != "" && input.QuantityType != x.QuantityType {
		x.QuantityType = input.QuantityType
	}

	if input.QuantityValue != 0 && input.QuantityValue != x.QuantityValue {
		x.QuantityValue = input.QuantityValue
	}

	if input.QuantityNotes != "" && input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = input.QuantityNotes
	}

	if input.ProductOfRecipe != x.ProductOfRecipe {
		x.ProductOfRecipe = input.ProductOfRecipe
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
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.IngredientNotes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientDatabaseCreationInput.
func (x *RecipeStepIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.IngredientNotes, validation.Required),
	)
}

// RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(input *RecipeStepIngredientCreationRequestInput) *RecipeStepIngredientDatabaseCreationInput {
	x := &RecipeStepIngredientDatabaseCreationInput{
		IngredientID:    input.IngredientID,
		QuantityType:    input.QuantityType,
		QuantityValue:   input.QuantityValue,
		QuantityNotes:   input.QuantityNotes,
		ProductOfRecipe: input.ProductOfRecipe,
		IngredientNotes: input.IngredientNotes,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientUpdateRequestInput.
func (x *RecipeStepIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.IngredientNotes, validation.Required),
	)
}
