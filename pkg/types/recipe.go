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

	// RecipeCreatedCustomerEventType indicates a recipe was created.
	RecipeCreatedCustomerEventType CustomerEventType = "recipe_created"
	// RecipeUpdatedCustomerEventType indicates a recipe was updated.
	RecipeUpdatedCustomerEventType CustomerEventType = "recipe_updated"
	// RecipeArchivedCustomerEventType indicates a recipe was archived.
	RecipeArchivedCustomerEventType CustomerEventType = "recipe_archived"
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

		LastUpdatedAt      *uint64       `json:"lastUpdatedAt"`
		ArchivedAt         *uint64       `json:"archivedAt"`
		InspiredByRecipeID *string       `json:"inspiredByRecipeID"`
		Source             string        `json:"source"`
		Description        string        `json:"description"`
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		CreatedByUser      string        `json:"belongsToUser"`
		Steps              []*RecipeStep `json:"steps"`
		SealOfApproval     bool          `json:"sealOfApproval"`
		YieldsPortions     uint8         `json:"yieldsPortions"`
		CreatedAt          uint64        `json:"createdAt"`
	}

	// RecipeList represents a list of recipes.
	RecipeList struct {
		_ struct{}

		Recipes []*Recipe `json:"data"`
		Pagination
	}

	// RecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipeCreationRequestInput struct {
		_                  struct{}
		InspiredByRecipeID *string                           `json:"inspiredByRecipeID"`
		CreatedByUser      string                            `json:"-"`
		ID                 string                            `json:"-"`
		Name               string                            `json:"name"`
		Source             string                            `json:"source"`
		Description        string                            `json:"description"`
		Steps              []*RecipeStepCreationRequestInput `json:"steps"`
		AlsoCreateMeal     bool                              `json:"alsoCreateMeal"`
		SealOfApproval     bool                              `json:"sealOfApproval"`
		YieldsPortions     uint8                             `json:"yieldsPortions"`
	}

	// RecipeDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipeDatabaseCreationInput struct {
		_                  struct{}
		InspiredByRecipeID *string                            `json:"inspiredByRecipeID"`
		CreatedByUser      string                             `json:"belongsToHousehold"`
		ID                 string                             `json:"id"`
		Name               string                             `json:"name"`
		Source             string                             `json:"source"`
		Description        string                             `json:"description"`
		Steps              []*RecipeStepDatabaseCreationInput `json:"steps"`
		AlsoCreateMeal     bool                               `json:"alsoCreateMeal"`
		SealOfApproval     bool                               `json:"sealOfApproval"`
		YieldsPortions     uint8                              `json:"yieldsPortions"`
	}

	// RecipeUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipeUpdateRequestInput struct {
		_ struct{}

		Name        *string `json:"name"`
		Source      *string `json:"source"`
		Description *string `json:"description"`
		// InspiredByRecipeID is already a pointer, I'm not about to make it a double pointer.
		InspiredByRecipeID *string `json:"inspiredByRecipeID"`
		CreatedByUser      *string `json:"-"`
		SealOfApproval     *bool   `json:"sealOfApproval"`
		YieldsPortions     *uint8  `json:"yieldsPortions"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID string) (bool, error)
		GetRecipe(ctx context.Context, recipeID string) (*Recipe, error)
		GetRecipeByIDAndUser(ctx context.Context, recipeID, userID string) (*Recipe, error)
		GetRecipes(ctx context.Context, filter *QueryFilter) (*RecipeList, error)
		SearchForRecipes(ctx context.Context, query string, filter *QueryFilter) (*RecipeList, error)
		GetRecipesWithIDs(ctx context.Context, userID string, limit uint8, ids []string) ([]*Recipe, error)
		CreateRecipe(ctx context.Context, input *RecipeDatabaseCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		ArchiveRecipe(ctx context.Context, recipeID, userID string) error
	}

	// RecipeDataService describes a structure capable of serving traffic related to recipes.
	RecipeDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		SearchHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeUpdateRequestInput with a recipe.
func (x *Recipe) Update(input *RecipeUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Source != nil && *input.Source != x.Source {
		x.Source = *input.Source
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.InspiredByRecipeID != nil && (x.InspiredByRecipeID == nil || (*input.InspiredByRecipeID != "" && *input.InspiredByRecipeID != *x.InspiredByRecipeID)) {
		x.InspiredByRecipeID = input.InspiredByRecipeID
	}

	if input.SealOfApproval != nil && *input.SealOfApproval != x.SealOfApproval {
		x.SealOfApproval = *input.SealOfApproval
	}

	if input.YieldsPortions != nil && *input.YieldsPortions != x.YieldsPortions {
		x.YieldsPortions = *input.YieldsPortions
	}
}

var _ validation.ValidatableWithContext = (*RecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeCreationRequestInput.
func (x *RecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Steps, validation.Required),
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
		validation.Field(&x.CreatedByUser, validation.Required),
	)
}

// RecipeUpdateRequestInputFromRecipe creates a DatabaseCreationInput from a CreationInput.
func RecipeUpdateRequestInputFromRecipe(input *Recipe) *RecipeUpdateRequestInput {
	x := &RecipeUpdateRequestInput{
		Name:               &input.Name,
		Source:             &input.Source,
		Description:        &input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		CreatedByUser:      &input.CreatedByUser,
		SealOfApproval:     &input.SealOfApproval,
		YieldsPortions:     &input.YieldsPortions,
	}

	return x
}

// RecipeDatabaseCreationInputFromRecipeCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeDatabaseCreationInputFromRecipeCreationInput(input *RecipeCreationRequestInput) *RecipeDatabaseCreationInput {
	steps := []*RecipeStepDatabaseCreationInput{}
	for _, step := range input.Steps {
		steps = append(steps, RecipeStepDatabaseCreationInputFromRecipeStepCreationInput(step))
	}

	x := &RecipeDatabaseCreationInput{
		AlsoCreateMeal:     input.AlsoCreateMeal,
		Name:               input.Name,
		Source:             input.Source,
		Description:        input.Description,
		InspiredByRecipeID: input.InspiredByRecipeID,
		Steps:              steps,
		SealOfApproval:     input.SealOfApproval,
		YieldsPortions:     input.YieldsPortions,
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
