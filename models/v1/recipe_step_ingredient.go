package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStepIngredient represents a recipe step ingredient
	RecipeStepIngredient struct {
		ID              uint64  `json:"id"`
		IngredientID    *uint64 `json:"ingredient_id"`
		QuantityType    string  `json:"quantity_type"`
		QuantityValue   float32 `json:"quantity_value"`
		QuantityNotes   string  `json:"quantity_notes"`
		ProductOfRecipe bool    `json:"product_of_recipe"`
		IngredientNotes string  `json:"ingredient_notes"`
		RecipeStepID    uint64  `json:"recipe_step_id"`
		CreatedOn       uint64  `json:"created_on"`
		UpdatedOn       *uint64 `json:"updated_on"`
		ArchivedOn      *uint64 `json:"archived_on"`
		BelongsTo       uint64  `json:"belongs_to"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients
	RecipeStepIngredientList struct {
		Pagination
		RecipeStepIngredients []RecipeStepIngredient `json:"recipe_step_ingredients"`
	}

	// RecipeStepIngredientCreationInput represents what a user could set as input for creating recipe step ingredients
	RecipeStepIngredientCreationInput struct {
		IngredientID    *uint64 `json:"ingredient_id"`
		QuantityType    string  `json:"quantity_type"`
		QuantityValue   float32 `json:"quantity_value"`
		QuantityNotes   string  `json:"quantity_notes"`
		ProductOfRecipe bool    `json:"product_of_recipe"`
		IngredientNotes string  `json:"ingredient_notes"`
		RecipeStepID    uint64  `json:"recipe_step_id"`
		BelongsTo       uint64  `json:"-"`
	}

	// RecipeStepIngredientUpdateInput represents what a user could set as input for updating recipe step ingredients
	RecipeStepIngredientUpdateInput struct {
		IngredientID    *uint64 `json:"ingredient_id"`
		QuantityType    string  `json:"quantity_type"`
		QuantityValue   float32 `json:"quantity_value"`
		QuantityNotes   string  `json:"quantity_notes"`
		ProductOfRecipe bool    `json:"product_of_recipe"`
		IngredientNotes string  `json:"ingredient_notes"`
		RecipeStepID    uint64  `json:"recipe_step_id"`
		BelongsTo       uint64  `json:"-"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently
	RecipeStepIngredientDataManager interface {
		GetRecipeStepIngredient(ctx context.Context, recipeStepIngredientID, userID uint64) (*RecipeStepIngredient, error)
		GetRecipeStepIngredientCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error)
		GetRecipeStepIngredients(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeStepIngredientList, error)
		GetAllRecipeStepIngredientsForUser(ctx context.Context, userID uint64) ([]RecipeStepIngredient, error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientCreationInput) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient) error
		ArchiveRecipeStepIngredient(ctx context.Context, id, userID uint64) error
	}

	// RecipeStepIngredientDataServer describes a structure capable of serving traffic related to recipe step ingredients
	RecipeStepIngredientDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepIngredientInput with a recipe step ingredient
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

	if input.RecipeStepID != x.RecipeStepID {
		x.RecipeStepID = input.RecipeStepID
	}
}

// ToInput creates a RecipeStepIngredientUpdateInput struct for a recipe step ingredient
func (x *RecipeStepIngredient) ToInput() *RecipeStepIngredientUpdateInput {
	return &RecipeStepIngredientUpdateInput{
		IngredientID:    x.IngredientID,
		QuantityType:    x.QuantityType,
		QuantityValue:   x.QuantityValue,
		QuantityNotes:   x.QuantityNotes,
		ProductOfRecipe: x.ProductOfRecipe,
		IngredientNotes: x.IngredientNotes,
		RecipeStepID:    x.RecipeStepID,
	}
}
