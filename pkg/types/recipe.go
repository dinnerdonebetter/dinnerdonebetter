package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeDataType indicates an event is related to a recipe.
	RecipeDataType dataType = "recipe"
)

func init() {
	gob.Register(new(Recipe))
	gob.Register(new(RecipeList))
	gob.Register(new(RecipeCreationRequestInput))
	gob.Register(new(RecipeUpdateRequestInput))
}

type (
	// Recipe represents a recipe.
	Recipe struct {
		_ struct{}

		LastUpdatedOn      *uint64       `json:"lastUpdatedOn"`
		ArchivedOn         *uint64       `json:"archivedOn"`
		InspiredByRecipeID *string       `json:"inspiredByRecipeID"`
		Source             string        `json:"source"`
		Description        string        `json:"description"`
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		BelongsToHousehold string        `json:"belongsToHousehold"`
		Steps              []*RecipeStep `json:"steps"`
		CreatedOn          uint64        `json:"createdOn"`
	}

	// FullRecipe represents a recipe.
	FullRecipe struct {
		_ struct{}

		LastUpdatedOn      *uint64           `json:"lastUpdatedOn"`
		ArchivedOn         *uint64           `json:"archivedOn"`
		InspiredByRecipeID *string           `json:"inspiredByRecipeID"`
		Source             string            `json:"source"`
		Description        string            `json:"description"`
		ID                 string            `json:"id"`
		Name               string            `json:"name"`
		BelongsToHousehold string            `json:"belongsToHousehold"`
		Steps              []*FullRecipeStep `json:"steps"`
		CreatedOn          uint64            `json:"createdOn"`
	}

	// RecipeList represents a list of recipes.
	RecipeList struct {
		_ struct{}

		Recipes []*Recipe `json:"recipes"`
		Pagination
	}

	// RecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipeCreationRequestInput struct {
		_ struct{}

		InspiredByRecipeID *string                           `json:"inspiredByRecipeID"`
		ID                 string                            `json:"-"`
		Name               string                            `json:"name"`
		Source             string                            `json:"source"`
		Description        string                            `json:"description"`
		BelongsToHousehold string                            `json:"-"`
		Steps              []*RecipeStepCreationRequestInput `json:"steps"`
	}

	// RecipeDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipeDatabaseCreationInput struct {
		_ struct{}

		InspiredByRecipeID *string                            `json:"inspiredByRecipeID"`
		ID                 string                             `json:"id"`
		Name               string                             `json:"name"`
		Source             string                             `json:"source"`
		Description        string                             `json:"description"`
		BelongsToHousehold string                             `json:"belongsToHousehold"`
		Steps              []*RecipeStepDatabaseCreationInput `json:"steps"`
	}

	// RecipeUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipeUpdateRequestInput struct {
		_ struct{}

		Name               string  `json:"name"`
		Source             string  `json:"source"`
		Description        string  `json:"description"`
		InspiredByRecipeID *string `json:"inspiredByRecipeID"`
		BelongsToHousehold string  `json:"-"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID string) (bool, error)
		GetRecipe(ctx context.Context, recipeID string) (*Recipe, error)
		GetTotalRecipeCount(ctx context.Context) (uint64, error)
		GetRecipes(ctx context.Context, filter *QueryFilter) (*RecipeList, error)
		GetRecipesWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*Recipe, error)
		CreateRecipe(ctx context.Context, input *RecipeDatabaseCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		ArchiveRecipe(ctx context.Context, recipeID, householdID string) error
	}

	// RecipeDataService describes a structure capable of serving traffic related to recipes.
	RecipeDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeUpdateRequestInput with a recipe.
func (x *Recipe) Update(input *RecipeUpdateRequestInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Source != "" && input.Source != x.Source {
		x.Source = input.Source
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.InspiredByRecipeID != nil && (x.InspiredByRecipeID == nil || (*input.InspiredByRecipeID != "" && *input.InspiredByRecipeID != *x.InspiredByRecipeID)) {
		x.InspiredByRecipeID = input.InspiredByRecipeID
	}
}

var _ validation.ValidatableWithContext = (*RecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeCreationRequestInput.
func (x *RecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Source, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.InspiredByRecipeID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeDatabaseCreationInput.
func (x *RecipeDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Source, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.InspiredByRecipeID, validation.Required),
		validation.Field(&x.BelongsToHousehold, validation.Required),
	)
}

// RecipeDatabaseCreationInputFromRecipeCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeDatabaseCreationInputFromRecipeCreationInput(input *RecipeCreationRequestInput) *RecipeDatabaseCreationInput {
	x := &RecipeDatabaseCreationInput{
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeUpdateRequestInput.
func (x *RecipeUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Source, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.InspiredByRecipeID, validation.Required),
	)
}
