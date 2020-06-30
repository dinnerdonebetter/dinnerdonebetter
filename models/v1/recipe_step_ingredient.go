package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		ID                    uint64  `json:"id"`
		ValidIngredientID     uint64  `json:"validIngredientID"`
		IngredientNotes       string  `json:"ingredientNotes"`
		QuantityType          string  `json:"quantityType"`
		QuantityValue         float32 `json:"quantityValue"`
		QuantityNotes         string  `json:"quantityNotes"`
		ProductOfRecipeStepID *uint64 `json:"productOfRecipeStepID"`
		CreatedOn             uint64  `json:"createdOn"`
		UpdatedOn             *uint64 `json:"updatedOn"`
		ArchivedOn            *uint64 `json:"archivedOn"`
		BelongsToRecipeStep   uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		Pagination
		RecipeStepIngredients []RecipeStepIngredient `json:"recipeStepIngredients"`
	}

	// RecipeStepIngredientCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationInput struct {
		ValidIngredientID     uint64  `json:"validIngredientID"`
		IngredientNotes       string  `json:"ingredientNotes"`
		QuantityType          string  `json:"quantityType"`
		QuantityValue         float32 `json:"quantityValue"`
		QuantityNotes         string  `json:"quantityNotes"`
		ProductOfRecipeStepID *uint64 `json:"productOfRecipeStepID"`
		BelongsToRecipeStep   uint64  `json:"-"`
	}

	// RecipeStepIngredientUpdateInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateInput struct {
		ValidIngredientID     uint64  `json:"validIngredientID"`
		IngredientNotes       string  `json:"ingredientNotes"`
		QuantityType          string  `json:"quantityType"`
		QuantityValue         float32 `json:"quantityValue"`
		QuantityNotes         string  `json:"quantityNotes"`
		ProductOfRecipeStepID *uint64 `json:"productOfRecipeStepID"`
		BelongsToRecipeStep   uint64  `json:"belongsToRecipeStep"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*RecipeStepIngredient, error)
		GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error)
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepIngredientList, error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientCreationInput) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) error
	}

	// RecipeStepIngredientDataServer describes a structure capable of serving traffic related to recipe step ingredients.
	RecipeStepIngredientDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ExistenceHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepIngredientInput with a recipe step ingredient.
func (x *RecipeStepIngredient) Update(input *RecipeStepIngredientUpdateInput) {
	if input.ValidIngredientID != x.ValidIngredientID {
		x.ValidIngredientID = input.ValidIngredientID
	}

	if input.IngredientNotes != "" && input.IngredientNotes != x.IngredientNotes {
		x.IngredientNotes = input.IngredientNotes
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

	if input.ProductOfRecipeStepID != nil && input.ProductOfRecipeStepID != x.ProductOfRecipeStepID {
		x.ProductOfRecipeStepID = input.ProductOfRecipeStepID
	}
}

// ToUpdateInput creates a RecipeStepIngredientUpdateInput struct for a recipe step ingredient.
func (x *RecipeStepIngredient) ToUpdateInput() *RecipeStepIngredientUpdateInput {
	return &RecipeStepIngredientUpdateInput{
		ValidIngredientID:     x.ValidIngredientID,
		IngredientNotes:       x.IngredientNotes,
		QuantityType:          x.QuantityType,
		QuantityValue:         x.QuantityValue,
		QuantityNotes:         x.QuantityNotes,
		ProductOfRecipeStepID: x.ProductOfRecipeStepID,
	}
}
