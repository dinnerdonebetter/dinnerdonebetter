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
	gob.Register(new(RecipeCreationRequestInput))
	gob.Register(new(RecipeUpdateRequestInput))
}

type (
	// Recipe represents a recipe.
	Recipe struct {
		_ struct{}

		CreatedAt                time.Time         `json:"createdAt"`
		InspiredByRecipeID       *string           `json:"inspiredByRecipeID"`
		LastUpdatedAt            *time.Time        `json:"lastUpdatedAt"`
		ArchivedAt               *time.Time        `json:"archivedAt"`
		MaximumEstimatedPortions *float32          `json:"maximumEstimatedPortions"`
		PluralPortionName        string            `json:"pluralPortionName"`
		Description              string            `json:"description"`
		Name                     string            `json:"name"`
		PortionName              string            `json:"portionName"`
		ID                       string            `json:"id"`
		CreatedByUser            string            `json:"belongsToUser"`
		Source                   string            `json:"source"`
		Slug                     string            `json:"slug"`
		Media                    []*RecipeMedia    `json:"media"`
		PrepTasks                []*RecipePrepTask `json:"prepTasks"`
		Steps                    []*RecipeStep     `json:"steps"`
		MinimumEstimatedPortions float32           `json:"minimumEstimatedPortions"`
		SealOfApproval           bool              `json:"sealOfApproval"`
		EligibleForMeals         bool              `json:"eligibleForMeals"`
	}

	// RecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipeCreationRequestInput struct {
		_ struct{}

		InspiredByRecipeID       *string                                           `json:"inspiredByRecipeID"`
		MaximumEstimatedPortions *float32                                          `json:"maximumEstimatedPortions"`
		Slug                     string                                            `json:"slug"`
		Source                   string                                            `json:"source"`
		Description              string                                            `json:"description"`
		PluralPortionName        string                                            `json:"pluralPortionName"`
		Name                     string                                            `json:"name"`
		PortionName              string                                            `json:"portionName"`
		PrepTasks                []*RecipePrepTaskWithinRecipeCreationRequestInput `json:"prepTasks"`
		Steps                    []*RecipeStepCreationRequestInput                 `json:"steps"`
		MinimumEstimatedPortions float32                                           `json:"minimumEstimatedPortions"`
		AlsoCreateMeal           bool                                              `json:"alsoCreateMeal"`
		SealOfApproval           bool                                              `json:"sealOfApproval"`
		EligibleForMeals         bool                                              `json:"eligibleForMeals"`
	}

	// RecipeDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipeDatabaseCreationInput struct {
		_ struct{}

		InspiredByRecipeID       *string
		MaximumEstimatedPortions *float32
		PluralPortionName        string
		ID                       string
		Name                     string
		Slug                     string
		Source                   string
		PortionName              string
		CreatedByUser            string
		Description              string
		PrepTasks                []*RecipePrepTaskDatabaseCreationInput
		Steps                    []*RecipeStepDatabaseCreationInput
		MinimumEstimatedPortions float32
		AlsoCreateMeal           bool
		SealOfApproval           bool
		EligibleForMeals         bool
	}

	// RecipeUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipeUpdateRequestInput struct {
		_ struct{}

		Name                     *string  `json:"name,omitempty"`
		Slug                     *string  `json:"slug"`
		Source                   *string  `json:"source,omitempty"`
		Description              *string  `json:"description,omitempty"`
		InspiredByRecipeID       *string  `json:"inspiredByRecipeID,omitempty"`
		CreatedByUser            *string  `json:"-"`
		SealOfApproval           *bool    `json:"sealOfApproval,omitempty"`
		MinimumEstimatedPortions *float32 `json:"minimumEstimatedPortions,omitempty"`
		MaximumEstimatedPortions *float32 `json:"maximumEstimatedPortions,omitempty"`
		PortionName              *string  `json:"portionName"`
		PluralPortionName        *string  `json:"pluralPortionName"`
		EligibleForMeals         *bool    `json:"eligibleForMeals"`
	}

	// RecipeDataManager describes a structure capable of storing recipes permanently.
	RecipeDataManager interface {
		RecipeExists(ctx context.Context, recipeID string) (bool, error)
		GetRecipe(ctx context.Context, recipeID string) (*Recipe, error)
		GetRecipeByIDAndUser(ctx context.Context, recipeID, userID string) (*Recipe, error)
		GetRecipes(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[Recipe], error)
		SearchForRecipes(ctx context.Context, query string, filter *QueryFilter) (*QueryFilteredResult[Recipe], error)
		CreateRecipe(ctx context.Context, input *RecipeDatabaseCreationInput) (*Recipe, error)
		UpdateRecipe(ctx context.Context, updated *Recipe) error
		ArchiveRecipe(ctx context.Context, recipeID, userID string) error
	}

	// RecipeDataService describes a structure capable of serving traffic related to recipes.
	RecipeDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		SearchHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		DAGHandler(http.ResponseWriter, *http.Request)
		EstimatedPrepStepsHandler(http.ResponseWriter, *http.Request)
		ImageUploadHandler(http.ResponseWriter, *http.Request)
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

	if input.Slug != nil && *input.Slug != x.Slug {
		x.Slug = *input.Slug
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

	if input.MinimumEstimatedPortions != nil && *input.MinimumEstimatedPortions != x.MinimumEstimatedPortions {
		x.MinimumEstimatedPortions = *input.MinimumEstimatedPortions
	}

	if input.MaximumEstimatedPortions != nil && input.MaximumEstimatedPortions != x.MaximumEstimatedPortions {
		x.MaximumEstimatedPortions = input.MaximumEstimatedPortions
	}

	if input.PortionName != nil && *input.PortionName != x.PortionName {
		x.PortionName = *input.PortionName
	}

	if input.PluralPortionName != nil && *input.PluralPortionName != x.PluralPortionName {
		x.PluralPortionName = *input.PluralPortionName
	}

	if input.EligibleForMeals != nil && *input.EligibleForMeals != x.EligibleForMeals {
		x.EligibleForMeals = *input.EligibleForMeals
	}
}

var _ validation.ValidatableWithContext = (*RecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeCreationRequestInput.
func (x *RecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var errResult *multierror.Error

	if len(x.Steps) < 2 {
		errResult = multierror.Append(fmt.Errorf("recipe must have at least 2 steps"), errResult)
	}

	stepIndicesMentionedInPrepTasks := map[uint32]bool{}
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
		validation.Field(&x.MinimumEstimatedPortions, validation.Required),
		validation.Field(&x.InspiredByRecipeID, validation.Required),
	)
}
