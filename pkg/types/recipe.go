package types

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
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
		_                  struct{}
		CreatedAt          time.Time         `json:"createdAt"`
		ArchivedAt         *time.Time        `json:"archivedAt"`
		InspiredByRecipeID *string           `json:"inspiredByRecipeID"`
		LastUpdatedAt      *time.Time        `json:"lastUpdatedAt"`
		Source             string            `json:"source"`
		Description        string            `json:"description"`
		Name               string            `json:"name"`
		CreatedByUser      string            `json:"belongsToUser"`
		ID                 string            `json:"id"`
		Steps              []*RecipeStep     `json:"steps"`
		PrepTasks          []*RecipePrepTask `json:"prepTasks"`
		SealOfApproval     bool              `json:"sealOfApproval"`
		YieldsPortions     uint8             `json:"yieldsPortions"`
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
		InspiredByRecipeID *string                                           `json:"inspiredByRecipeID"`
		CreatedByUser      string                                            `json:"-"`
		ID                 string                                            `json:"-"`
		Name               string                                            `json:"name"`
		Source             string                                            `json:"source"`
		Description        string                                            `json:"description"`
		Steps              []*RecipeStepCreationRequestInput                 `json:"steps"`
		PrepTasks          []*RecipePrepTaskWithinRecipeCreationRequestInput `json:"prepTasks"`
		AlsoCreateMeal     bool                                              `json:"alsoCreateMeal"`
		SealOfApproval     bool                                              `json:"sealOfApproval"`
		YieldsPortions     uint8                                             `json:"yieldsPortions"`
	}

	// RecipeDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipeDatabaseCreationInput struct {
		_                  struct{}
		InspiredByRecipeID *string                                `json:"inspiredByRecipeID"`
		CreatedByUser      string                                 `json:"belongsToHousehold"`
		ID                 string                                 `json:"id"`
		Name               string                                 `json:"name"`
		Source             string                                 `json:"source"`
		Description        string                                 `json:"description"`
		Steps              []*RecipeStepDatabaseCreationInput     `json:"steps"`
		PrepTasks          []*RecipePrepTaskDatabaseCreationInput `json:"prepTasks"`
		AlsoCreateMeal     bool                                   `json:"alsoCreateMeal"`
		SealOfApproval     bool                                   `json:"sealOfApproval"`
		YieldsPortions     uint8                                  `json:"yieldsPortions"`
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
		DAGHandler(res http.ResponseWriter, req *http.Request)
		EstimatedPrepStepsHandler(res http.ResponseWriter, req *http.Request)
		ImageUploadHandler(res http.ResponseWriter, req *http.Request)
	}
)

// FindStepForIndex finds a step for a given index.
func (x *Recipe) FindStepForIndex(index uint32) *RecipeStep {
	for _, step := range x.Steps {
		if step.Index == index {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepByID finds a step for a given ID.
func (x *Recipe) FindStepByID(id string) *RecipeStep {
	for _, step := range x.Steps {
		if step.ID == id {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

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
	stepIndicesMentionedInPrepTasks := map[uint32]bool{}

	var errResult *multierror.Error
	for i, task := range x.PrepTasks {
		for j, step := range task.TaskSteps {
			if _, ok := stepIndicesMentionedInPrepTasks[step.BelongsToRecipeStepIndex]; ok {
				errResult = multierror.Append(fmt.Errorf("duplicate step mentioned in step %d for task %d", i+1, j+1), errResult)
			} else {
				stepIndicesMentionedInPrepTasks[step.BelongsToRecipeStepIndex] = true
			}
		}
	}
	if errResult != nil {
		return errResult
	}

	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Steps, validation.Required),
	)
}

// FindStepByIndex finds a step for a given index.
func (x *RecipeDatabaseCreationInput) FindStepByIndex(index uint32) *RecipeStepDatabaseCreationInput {
	for _, step := range x.Steps {
		if step.Index == index {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
}

// FindStepByID finds a step for a given ID.
func (x *RecipeDatabaseCreationInput) FindStepByID(id string) *RecipeStepDatabaseCreationInput {
	for _, step := range x.Steps {
		if step.ID == id {
			return step
		}
	}

	// we could return an error here, but that would make my life a little harder
	// so if you fuck up and submit a wrong value, and it's your fault here.
	return nil
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
