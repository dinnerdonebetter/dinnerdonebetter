package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		ID                    uint64  `json:"id"`
		ValidIngredientID     uint64  `json:"valid_ingredient_id"`
		IngredientNotes       string  `json:"ingredient_notes"`
		QuantityType          string  `json:"quantity_type"`
		QuantityValue         float32 `json:"quantity_value"`
		QuantityNotes         string  `json:"quantity_notes"`
		ProductOfRecipeStepID *uint64 `json:"product_of_recipe_step_id"`
		CreatedOn             uint64  `json:"created_on"`
		UpdatedOn             *uint64 `json:"updated_on"`
		ArchivedOn            *uint64 `json:"archived_on"`
		BelongsToRecipeStep   uint64  `json:"belongs_to_recipe_step"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		Pagination
		RecipeStepIngredients []RecipeStepIngredient `json:"recipe_step_ingredients"`
	}

	// RecipeStepIngredientCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationInput struct {
		ValidIngredientID     uint64  `json:"valid_ingredient_id"`
		IngredientNotes       string  `json:"ingredient_notes"`
		QuantityType          string  `json:"quantity_type"`
		QuantityValue         float32 `json:"quantity_value"`
		QuantityNotes         string  `json:"quantity_notes"`
		ProductOfRecipeStepID *uint64 `json:"product_of_recipe_step_id"`
		BelongsToRecipeStep   uint64  `json:"-"`
	}

	// RecipeStepIngredientUpdateInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateInput struct {
		ValidIngredientID     uint64  `json:"valid_ingredient_id"`
		IngredientNotes       string  `json:"ingredient_notes"`
		QuantityType          string  `json:"quantity_type"`
		QuantityValue         float32 `json:"quantity_value"`
		QuantityNotes         string  `json:"quantity_notes"`
		ProductOfRecipeStepID *uint64 `json:"product_of_recipe_step_id"`
		BelongsToRecipeStep   uint64  `json:"belongs_to_recipe_step"`
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
