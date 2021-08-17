package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		LastUpdatedOn       *uint64 `json:"lastUpdatedOn"`
		IngredientID        *uint64 `json:"ingredientID"`
		ArchivedOn          *uint64 `json:"archivedOn"`
		IngredientNotes     string  `json:"ingredientNotes"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		ExternalID          string  `json:"externalID"`
		ID                  uint64  `json:"id"`
		CreatedOn           uint64  `json:"createdOn"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientList represents a list of recipe step ingredients.
	RecipeStepIngredientList struct {
		RecipeStepIngredients []*RecipeStepIngredient `json:"recipeStepIngredients"`
		Pagination
	}

	// RecipeStepIngredientCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationInput struct {
		IngredientID        *uint64 `json:"ingredientID"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep uint64  `json:"-"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientUpdateInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateInput struct {
		IngredientID        *uint64 `json:"ingredientID"`
		Name                string  `json:"name"`
		QuantityType        string  `json:"quantityType"`
		QuantityNotes       string  `json:"quantityNotes"`
		IngredientNotes     string  `json:"ingredientNotes"`
		BelongsToRecipeStep uint64  `json:"belongsToRecipeStep"`
		QuantityValue       float32 `json:"quantityValue"`
		ProductOfRecipeStep bool    `json:"productOfRecipeStep"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (*RecipeStepIngredient, error)
		GetAllRecipeStepIngredientsCount(ctx context.Context) (uint64, error)
		GetAllRecipeStepIngredients(ctx context.Context, resultChannel chan []*RecipeStepIngredient, bucketSize uint16) error
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID uint64, filter *QueryFilter) (*RecipeStepIngredientList, error)
		GetRecipeStepIngredientsWithIDs(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) ([]*RecipeStepIngredient, error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientCreationInput, createdByUser uint64) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID, archivedBy uint64) error
		GetAuditLogEntriesForRecipeStepIngredient(ctx context.Context, recipeStepIngredientID uint64) ([]*AuditLogEntry, error)
	}

	// RecipeStepIngredientDataService describes a structure capable of serving traffic related to recipe step ingredients.
	RecipeStepIngredientDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepIngredientUpdateInput with a recipe step ingredient.
func (x *RecipeStepIngredient) Update(input *RecipeStepIngredientUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.IngredientID != nil && (x.IngredientID == nil || (*input.IngredientID != 0 && *input.IngredientID != *x.IngredientID)) {
		out = append(out, &FieldChangeSummary{
			FieldName: "IngredientID",
			OldValue:  x.IngredientID,
			NewValue:  input.IngredientID,
		})

		x.IngredientID = input.IngredientID
	}

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	if input.QuantityType != x.QuantityType {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityType",
			OldValue:  x.QuantityType,
			NewValue:  input.QuantityType,
		})

		x.QuantityType = input.QuantityType
	}

	if input.QuantityValue != 0 && input.QuantityValue != x.QuantityValue {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityValue",
			OldValue:  x.QuantityValue,
			NewValue:  input.QuantityValue,
		})

		x.QuantityValue = input.QuantityValue
	}

	if input.QuantityNotes != x.QuantityNotes {
		out = append(out, &FieldChangeSummary{
			FieldName: "QuantityNotes",
			OldValue:  x.QuantityNotes,
			NewValue:  input.QuantityNotes,
		})

		x.QuantityNotes = input.QuantityNotes
	}

	if input.ProductOfRecipeStep != x.ProductOfRecipeStep {
		out = append(out, &FieldChangeSummary{
			FieldName: "ProductOfRecipeStep",
			OldValue:  x.ProductOfRecipeStep,
			NewValue:  input.ProductOfRecipeStep,
		})

		x.ProductOfRecipeStep = input.ProductOfRecipeStep
	}

	if input.IngredientNotes != x.IngredientNotes {
		out = append(out, &FieldChangeSummary{
			FieldName: "IngredientNotes",
			OldValue:  x.IngredientNotes,
			NewValue:  input.IngredientNotes,
		})

		x.IngredientNotes = input.IngredientNotes
	}

	return out
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientCreationInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientCreationInput.
func (x *RecipeStepIngredientCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.IngredientNotes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientUpdateInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientUpdateInput.
func (x *RecipeStepIngredientUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.QuantityType, validation.Required),
		validation.Field(&x.QuantityValue, validation.Required),
		validation.Field(&x.QuantityNotes, validation.Required),
		validation.Field(&x.IngredientNotes, validation.Required),
	)
}
