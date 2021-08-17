package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Recipe represents a recipe.
	Recipe struct {
		LastUpdatedOn      *uint64       `json:"lastUpdatedOn"`
		ArchivedOn         *uint64       `json:"archivedOn"`
		InspiredByRecipeID *uint64       `json:"inspiredByRecipeID"`
		Source             string        `json:"source"`
		Description        string        `json:"description"`
		ExternalID         string        `json:"externalID"`
		DisplayImageURL    string        `json:"displayImageURL"`
		Name               string        `json:"name"`
		Steps              []*RecipeStep `json:"steps"`
		ID                 uint64        `json:"id"`
		CreatedOn          uint64        `json:"createdOn"`
		BelongsToAccount   uint64        `json:"belongsToAccount"`
	}

	// RecipeList represents a list of recipes.
	RecipeList struct {
		Recipes []*Recipe `json:"recipes"`
		Pagination
	}

	// RecipeCreationInput represents what a user could set as input for creating recipes.
	RecipeCreationInput struct {
		InspiredByRecipeID *uint64                    `json:"inspiredByRecipeID"`
		Name               string                     `json:"name"`
		Source             string                     `json:"source"`
		DisplayImageURL    string                     `json:"displayImageURL"`
		Description        string                     `json:"description"`
		Steps              []*RecipeStepCreationInput `json:"steps"`
		BelongsToAccount   uint64                     `json:"-"`
	}

	// RecipeUpdateInput represents what a user could set as input for updating recipes.
	RecipeUpdateInput struct {
		InspiredByRecipeID *uint64 `json:"inspiredByRecipeID"`
		Name               string  `json:"name"`
		Source             string  `json:"source"`
		DisplayImageURL    string  `json:"displayImageURL"`
		Description        string  `json:"description"`
		BelongsToAccount   uint64  `json:"-"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID uint64) (bool, error)
		GetRecipe(ctx context.Context, recipeID uint64) (*Recipe, error)
		GetAllRecipesCount(ctx context.Context) (uint64, error)
		GetAllRecipes(ctx context.Context, resultChannel chan []*Recipe, bucketSize uint16) error
		GetRecipes(ctx context.Context, filter *QueryFilter) (*RecipeList, error)
		GetRecipesWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*Recipe, error)
		CreateRecipe(ctx context.Context, input *RecipeCreationInput, createdByUser uint64) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveRecipe(ctx context.Context, recipeID, accountID, archivedBy uint64) error
		GetAuditLogEntriesForRecipe(ctx context.Context, recipeID uint64) ([]*AuditLogEntry, error)
	}

	// RecipeDataService describes a structure capable of serving traffic related to recipes.
	RecipeDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeUpdateInput with a recipe.
func (x *Recipe) Update(input *RecipeUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	if input.Source != x.Source {
		out = append(out, &FieldChangeSummary{
			FieldName: "Source",
			OldValue:  x.Source,
			NewValue:  input.Source,
		})

		x.Source = input.Source
	}

	if input.Description != x.Description {
		out = append(out, &FieldChangeSummary{
			FieldName: "Description",
			OldValue:  x.Description,
			NewValue:  input.Description,
		})

		x.Description = input.Description
	}

	if input.InspiredByRecipeID != nil && (x.InspiredByRecipeID == nil || (*input.InspiredByRecipeID != 0 && *input.InspiredByRecipeID != *x.InspiredByRecipeID)) {
		out = append(out, &FieldChangeSummary{
			FieldName: "InspiredByRecipeID",
			OldValue:  x.InspiredByRecipeID,
			NewValue:  input.InspiredByRecipeID,
		})

		x.InspiredByRecipeID = input.InspiredByRecipeID
	}

	return out
}

var _ validation.ValidatableWithContext = (*RecipeCreationInput)(nil)

// ValidateWithContext validates a RecipeCreationInput.
func (x *RecipeCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeUpdateInput)(nil)

// ValidateWithContext validates a RecipeUpdateInput.
func (x *RecipeUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Source, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.InspiredByRecipeID, validation.Required),
	)
}
