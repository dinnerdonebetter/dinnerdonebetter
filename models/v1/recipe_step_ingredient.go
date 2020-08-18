package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		ID                  uint64  `json:"id"`
		IngredientID        *uint64 `json:"ingredientID"`
		CreatedOn           uint64  `json:"createdOn"`
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		Pagination
		RecipeStepIngredients []RecipeStepIngredient `json:"recipe_step_ingredients"`
	}

	// RecipeStepIngredientCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationInput struct {
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		IngredientID        *uint64 `json:"ingredientID"`
		BelongsToRecipeStep uint64  `json:"-"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// RecipeStepIngredientUpdateInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateInput struct {
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		IngredientID        *uint64 `json:"ingredientID"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipe     bool    `json:"productOfRecipe"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*RecipeStepIngredient, error)
		GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepIngredients(ctx context.Context, resultChannel chan []RecipeStepIngredient) error
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepIngredientList, error)
		GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeID, recipeStepID uint64, limit uint8, ids []uint64) ([]RecipeStepIngredient, error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientCreationInput) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) error
	}

	// RecipeStepIngredientDataServer describes a structure capable of serving traffic related to recipe step ingredients.
	RecipeStepIngredientDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepIngredientInput with a recipe step ingredient.
func (x *RecipeStepIngredient) Update(input *RecipeStepIngredientUpdateInput) {
	if input.IngredientID != nil && input.IngredientID != x.IngredientID {
		x.IngredientID = input.IngredientID
	}

	if input.QuantityType != "" && input.QuantityType != x.QuantityType {
		x.QuantityType = input.QuantityType
	}

	if input.QuantityValue != x.QuantityValue {
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

// ToUpdateInput creates a RecipeStepIngredientUpdateInput struct for a recipe step ingredient.
func (x *RecipeStepIngredient) ToUpdateInput() *RecipeStepIngredientUpdateInput {
	return &RecipeStepIngredientUpdateInput{
		IngredientID:    x.IngredientID,
		QuantityType:    x.QuantityType,
		QuantityValue:   x.QuantityValue,
		QuantityNotes:   x.QuantityNotes,
		ProductOfRecipe: x.ProductOfRecipe,
		IngredientNotes: x.IngredientNotes,
	}
}
