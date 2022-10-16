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
	// RecipePrepTaskDataType indicates an event is related to a recipe prep task.
	RecipePrepTaskDataType dataType = "recipe_prep_step"

	// RecipePrepTaskStorageTemperatureModifier is what we multiply/divide floats by to store in the database.
	RecipePrepTaskStorageTemperatureModifier = 100

	// RecipePrepTaskStorageTypeUncovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeUncovered = "uncovered"
	// RecipePrepTaskStorageTypeCovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeCovered = "covered"
	// RecipePrepTaskStorageTypeWireRack is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeWireRack = "on a wire rack"
	// RecipePrepTaskStorageTypeAirtightContainer is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeAirtightContainer = "in an airtight container"

	// RecipePrepTaskCreatedCustomerEventType indicates a recipe prep task was created.
	RecipePrepTaskCreatedCustomerEventType CustomerEventType = "recipe_created"
	// RecipePrepTaskUpdatedCustomerEventType indicates a recipe prep task was updated.
	RecipePrepTaskUpdatedCustomerEventType CustomerEventType = "recipe_updated"
	// RecipePrepTaskArchivedCustomerEventType indicates a recipe prep task was archived.
	RecipePrepTaskArchivedCustomerEventType CustomerEventType = "recipe_archived"
)

func init() {
	gob.Register(new(RecipePrepTask))
	gob.Register(new(RecipePrepTaskList))
	gob.Register(new(RecipePrepTaskCreationRequestInput))
	gob.Register(new(RecipePrepTaskUpdateRequestInput))
}

type (
	// RecipePrepTask represents a recipe prep task.
	RecipePrepTask struct {
		_                                      struct{}
		CreatedAt                              time.Time             `json:"createdAt"`
		ArchivedAt                             *time.Time            `json:"archivedAt"`
		LastUpdatedAt                          *time.Time            `json:"lastUpdatedAt"`
		Notes                                  string                `json:"notes"`
		ExplicitStorageInstructions            string                `json:"explicitStorageInstructions"`
		StorageType                            string                `json:"storageType"`
		BelongsToRecipe                        string                `json:"belongsToRecipe"`
		ID                                     string                `json:"id"`
		TaskSteps                              []*RecipePrepTaskStep `json:"recipeSteps"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                `json:"minimumTimeBufferBeforeRecipeInSeconds"`
		MaximumStorageTemperatureInCelsius     float32               `json:"maximumStorageTemperatureInCelsius"`
		MaximumTimeBufferBeforeRecipeInSeconds uint32                `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     float32               `json:"minimumStorageTemperatureInCelsius"`
	}

	// RecipePrepTaskList represents a list of recipe prep tasks.
	RecipePrepTaskList struct {
		_ struct{}

		RecipePrepTasks []*RecipePrepTask `json:"data"`
		Pagination
	}

	// RecipePrepTaskCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskCreationRequestInput struct {
		_                                      struct{}
		Notes                                  string                                    `json:"notes"`
		ExplicitStorageInstructions            string                                    `json:"explicitStorageInstructions"`
		StorageType                            string                                    `json:"storageType"`
		BelongsToRecipe                        string                                    `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepCreationRequestInput `json:"recipeSteps"`
		MaximumTimeBufferBeforeRecipeInSeconds uint32                                    `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     float32                                   `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     float32                                   `json:"maximumStorageTemperatureInCelsius"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                                    `json:"minimumTimeBufferBeforeRecipeInSeconds"`
	}

	// RecipePrepTaskWithinRecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskWithinRecipeCreationRequestInput struct {
		_                                      struct{}
		Notes                                  string                                                `json:"notes"`
		ExplicitStorageInstructions            string                                                `json:"explicitStorageInstructions"`
		StorageType                            string                                                `json:"storageType"`
		BelongsToRecipe                        string                                                `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepWithinRecipeCreationRequestInput `json:"recipeSteps"`
		MaximumTimeBufferBeforeRecipeInSeconds uint32                                                `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     float32                                               `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     float32                                               `json:"maximumStorageTemperatureInCelsius"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                                                `json:"minimumTimeBufferBeforeRecipeInSeconds"`
	}

	// RecipePrepTaskDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipePrepTaskDatabaseCreationInput struct {
		_                                      struct{}
		ID                                     string                                     `json:"id"`
		Notes                                  string                                     `json:"notes"`
		ExplicitStorageInstructions            string                                     `json:"explicitStorageInstructions"`
		StorageType                            string                                     `json:"storageType"`
		BelongsToRecipe                        string                                     `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepDatabaseCreationInput `json:"recipeSteps"`
		MaximumTimeBufferBeforeRecipeInSeconds uint32                                     `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		MinimumStorageTemperatureInCelsius     uint32                                     `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     uint32                                     `json:"maximumStorageTemperatureInCelsius"`
		MinimumTimeBufferBeforeRecipeInSeconds uint32                                     `json:"minimumTimeBufferBeforeRecipeInSeconds"`
	}

	// RecipePrepTaskUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipePrepTaskUpdateRequestInput struct {
		_ struct{}

		Notes                                  *string                                 `json:"notes"`
		ExplicitStorageInstructions            *string                                 `json:"explicitStorageInstructions"`
		MinimumTimeBufferBeforeRecipeInSeconds *uint32                                 `json:"minimumTimeBufferBeforeRecipeInSeconds"`
		MaximumTimeBufferBeforeRecipeInSeconds *uint32                                 `json:"maximumTimeBufferBeforeRecipeInSeconds"`
		StorageType                            *string                                 `json:"storageType"`
		MinimumStorageTemperatureInCelsius     *float32                                `json:"minimumStorageTemperatureInCelsius"`
		MaximumStorageTemperatureInCelsius     *float32                                `json:"maximumStorageTemperatureInCelsius"`
		BelongsToRecipe                        *string                                 `json:"belongsToRecipe"`
		TaskSteps                              []*RecipePrepTaskStepUpdateRequestInput `json:"recipeSteps"`
	}

	// RecipePrepTaskDataManager describes a structure capable of storing recipes permanently.
	RecipePrepTaskDataManager interface {
		RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error)
		GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*RecipePrepTask, error)
		GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*RecipePrepTask, error)
		CreateRecipePrepTask(ctx context.Context, input *RecipePrepTaskDatabaseCreationInput) (*RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, updated *RecipePrepTask) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error
	}

	// RecipePrepTaskDataService describes a structure capable of serving traffic related to recipes.
	RecipePrepTaskDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipePrepTaskUpdateRequestInput with a recipe prep task.
func (x *RecipePrepTask) Update(input *RecipePrepTaskUpdateRequestInput) {
	if input.BelongsToRecipe != nil && *input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ExplicitStorageInstructions != nil && *input.ExplicitStorageInstructions != x.ExplicitStorageInstructions {
		x.ExplicitStorageInstructions = *input.ExplicitStorageInstructions
	}

	if input.MinimumTimeBufferBeforeRecipeInSeconds != nil && *input.MinimumTimeBufferBeforeRecipeInSeconds != x.MinimumTimeBufferBeforeRecipeInSeconds {
		x.MinimumTimeBufferBeforeRecipeInSeconds = *input.MinimumTimeBufferBeforeRecipeInSeconds
	}

	if input.MaximumTimeBufferBeforeRecipeInSeconds != nil && *input.MaximumTimeBufferBeforeRecipeInSeconds != x.MaximumTimeBufferBeforeRecipeInSeconds {
		x.MaximumTimeBufferBeforeRecipeInSeconds = *input.MaximumTimeBufferBeforeRecipeInSeconds
	}

	if input.StorageType != nil && *input.StorageType != x.StorageType {
		x.StorageType = *input.StorageType
	}

	if input.MinimumStorageTemperatureInCelsius != nil && *input.MinimumStorageTemperatureInCelsius != x.MinimumStorageTemperatureInCelsius {
		x.MinimumStorageTemperatureInCelsius = *input.MinimumStorageTemperatureInCelsius
	}

	if input.MaximumStorageTemperatureInCelsius != nil && *input.MaximumStorageTemperatureInCelsius != x.MaximumStorageTemperatureInCelsius {
		x.MaximumStorageTemperatureInCelsius = *input.MaximumStorageTemperatureInCelsius
	}

	if input.BelongsToRecipe != nil && *input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskCreationRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskCreationRequestInput.
func (x *RecipePrepTaskCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
		validation.Field(&x.StorageType, validation.Required),
		validation.Field(&x.MaximumTimeBufferBeforeRecipeInSeconds, validation.Required),
		validation.Field(&x.MinimumStorageTemperatureInCelsius, validation.Required),
		validation.Field(&x.MaximumStorageTemperatureInCelsius, validation.Required),
		validation.Field(&x.MinimumTimeBufferBeforeRecipeInSeconds, validation.Required),
	)

	// TODO: uncomment me
	// if x.MinimumStorageTemperatureInCelsius > x.MaximumStorageTemperatureInCelsius {
	// 	result = multierror.Append(fmt.Errorf("minimum storage temperature (%d) is greater than maximum storage temperature (%d)", x.MinimumStorageTemperatureInCelsius, x.MaximumStorageTemperatureInCelsius))
	// }

	if err != nil {
		result = multierror.Append(err, result)
	}

	if result != nil {
		return result
	}

	return nil
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipePrepTaskDatabaseCreationInput.
func (x *RecipePrepTaskDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}

// RecipePrepTaskUpdateRequestInputFromRecipePrepTask creates a DatabaseCreationInput from a CreationInput.
func RecipePrepTaskUpdateRequestInputFromRecipePrepTask(input *RecipePrepTask) *RecipePrepTaskUpdateRequestInput {
	taskSteps := []*RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, &RecipePrepTaskStepUpdateRequestInput{
			ID:                      x.ID,
			BelongsToRecipeStep:     &x.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &x.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     &x.SatisfiesRecipeStep,
		})
	}
	x := &RecipePrepTaskUpdateRequestInput{
		Notes:                                  &input.Notes,
		ExplicitStorageInstructions:            &input.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: &input.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: &input.MaximumTimeBufferBeforeRecipeInSeconds,
		StorageType:                            &input.StorageType,
		MinimumStorageTemperatureInCelsius:     &input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     &input.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        &input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
	}

	return x
}

// RecipePrepTaskDatabaseCreationInputFromRecipePrepTaskCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipePrepTaskDatabaseCreationInputFromRecipePrepTaskCreationInput(input *RecipePrepTaskCreationRequestInput) *RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*RecipePrepTaskStepDatabaseCreationInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, &RecipePrepTaskStepDatabaseCreationInput{
			BelongsToRecipeStep:     x.BelongsToRecipeStep,
			BelongsToRecipePrepTask: x.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     x.SatisfiesRecipeStep,
		})
	}

	x := &RecipePrepTaskDatabaseCreationInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     uint32(input.MinimumStorageTemperatureInCelsius * RecipePrepTaskStorageTemperatureModifier),
		MaximumStorageTemperatureInCelsius:     uint32(input.MaximumStorageTemperatureInCelsius * RecipePrepTaskStorageTemperatureModifier),
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	return x
}

// RecipePrepTaskDatabaseCreationInputFromRecipePrepTaskWithinRecipeCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipePrepTaskDatabaseCreationInputFromRecipePrepTaskWithinRecipeCreationInput(recipe *RecipeDatabaseCreationInput, input *RecipePrepTaskWithinRecipeCreationRequestInput) (*RecipePrepTaskDatabaseCreationInput, error) {
	taskSteps := []*RecipePrepTaskStepDatabaseCreationInput{}
	for i, x := range input.TaskSteps {
		if y := recipe.FindStepByIndex(x.BelongsToRecipeStepIndex); y != nil {
			taskSteps = append(taskSteps, &RecipePrepTaskStepDatabaseCreationInput{
				BelongsToRecipeStep:     y.ID,
				BelongsToRecipePrepTask: x.BelongsToRecipePrepTask,
				SatisfiesRecipeStep:     x.SatisfiesRecipeStep,
			})
		} else {
			return nil, fmt.Errorf("task step #%d has an invalid recipe step index", i+1)
		}
	}

	x := &RecipePrepTaskDatabaseCreationInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     uint32(input.MinimumStorageTemperatureInCelsius * RecipePrepTaskStorageTemperatureModifier),
		MaximumStorageTemperatureInCelsius:     uint32(input.MaximumStorageTemperatureInCelsius * RecipePrepTaskStorageTemperatureModifier),
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}

	return x, nil
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskUpdateRequestInput.
func (x *RecipePrepTaskUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}
