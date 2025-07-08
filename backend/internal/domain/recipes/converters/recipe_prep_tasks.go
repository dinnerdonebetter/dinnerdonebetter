package converters

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput creates a RecipePrepTaskUpdateRequestInput from a RecipePrepTask.
func ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(input *recipes.RecipePrepTask) *recipes.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*recipes.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		y := x
		taskSteps = append(taskSteps, &recipes.RecipePrepTaskStepUpdateRequestInput{
			BelongsToRecipeStep:     &y.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &y.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     &y.SatisfiesRecipeStep,
		})
	}
	x := &recipes.RecipePrepTaskUpdateRequestInput{
		Name:                        &input.Name,
		Description:                 &input.Description,
		Notes:                       &input.Notes,
		ExplicitStorageInstructions: &input.ExplicitStorageInstructions,
		Optional:                    &input.Optional,
		StorageTemperatureInCelsius: input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: types.Uint32RangeWithOptionalMaxUpdateRequestInput{
			Min: &input.TimeBufferBeforeRecipeInSeconds.Min,
			Max: input.TimeBufferBeforeRecipeInSeconds.Max,
		},
		StorageType:     &input.StorageType,
		BelongsToRecipe: &input.BelongsToRecipe,
		TaskSteps:       taskSteps,
	}

	return x
}

// ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput(input *recipes.RecipePrepTaskCreationRequestInput) *recipes.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*recipes.RecipePrepTaskStepDatabaseCreationInput{}
	for _, x := range input.RecipeSteps {
		taskSteps = append(taskSteps, &recipes.RecipePrepTaskStepDatabaseCreationInput{
			BelongsToRecipeStep: x.BelongsToRecipeStep,
			SatisfiesRecipeStep: x.SatisfiesRecipeStep,
		})
	}

	x := &recipes.RecipePrepTaskDatabaseCreationInput{
		ID:                              identifiers.New(),
		Name:                            input.Name,
		Description:                     input.Description,
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Optional:                        input.Optional,
		StorageType:                     input.StorageType,
		BelongsToRecipe:                 input.BelongsToRecipe,
		TaskSteps:                       taskSteps,
		StorageTemperatureInCelsius:     input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
	}

	return x
}

// ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(recipe *recipes.RecipeDatabaseCreationInput, input *recipes.RecipePrepTaskWithinRecipeCreationRequestInput) (*recipes.RecipePrepTaskDatabaseCreationInput, error) {
	x := &recipes.RecipePrepTaskDatabaseCreationInput{
		ID:                              identifiers.New(),
		Name:                            input.Name,
		Description:                     input.Description,
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Optional:                        input.Optional,
		StorageType:                     input.StorageType,
		BelongsToRecipe:                 input.BelongsToRecipe,
		StorageTemperatureInCelsius:     input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
	}

	x.TaskSteps = []*recipes.RecipePrepTaskStepDatabaseCreationInput{}
	for i, ts := range input.RecipeSteps {
		if rs := recipe.FindStepByIndex(ts.BelongsToRecipeStepIndex); rs != nil {
			x.TaskSteps = append(x.TaskSteps, &recipes.RecipePrepTaskStepDatabaseCreationInput{
				ID:                      identifiers.New(),
				BelongsToRecipeStep:     rs.ID,
				BelongsToRecipePrepTask: x.ID,
				SatisfiesRecipeStep:     ts.SatisfiesRecipeStep,
			})
		} else {
			return nil, fmt.Errorf("task step #%d has an invalid recipe step index", i+1)
		}
	}

	return x, nil
}

// ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput builds a RecipePrepTaskDatabaseCreationInput from a recipe prep task.
func ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(input *recipes.RecipePrepTask) *recipes.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*recipes.RecipePrepTaskStepDatabaseCreationInput{}
	for _, step := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(step))
	}

	return &recipes.RecipePrepTaskDatabaseCreationInput{
		ID:                              input.ID,
		Name:                            input.Name,
		Description:                     input.Description,
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Optional:                        input.Optional,
		StorageType:                     input.StorageType,
		TaskSteps:                       taskSteps,
		StorageTemperatureInCelsius:     input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
		BelongsToRecipe:                 input.BelongsToRecipe,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(input *recipes.RecipePrepTaskStep) *recipes.RecipePrepTaskStepCreationRequestInput {
	return &recipes.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		SatisfiesRecipeStep: input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(input *recipes.RecipePrepTask) *recipes.RecipePrepTaskCreationRequestInput {
	taskSteps := []*recipes.RecipePrepTaskStepCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(x))
	}

	return &recipes.RecipePrepTaskCreationRequestInput{
		Name:                            input.Name,
		Description:                     input.Description,
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Optional:                        input.Optional,
		StorageType:                     input.StorageType,
		BelongsToRecipe:                 input.BelongsToRecipe,
		RecipeSteps:                     taskSteps,
		StorageTemperatureInCelsius:     input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(recipe *recipes.Recipe, input *recipes.RecipePrepTask) *recipes.RecipePrepTaskWithinRecipeCreationRequestInput {
	taskSteps := []*recipes.RecipePrepTaskStepWithinRecipeCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe, x))
	}

	return &recipes.RecipePrepTaskWithinRecipeCreationRequestInput{
		Name:                            input.Name,
		Description:                     input.Description,
		Notes:                           input.Notes,
		ExplicitStorageInstructions:     input.ExplicitStorageInstructions,
		Optional:                        input.Optional,
		StorageType:                     input.StorageType,
		BelongsToRecipe:                 input.BelongsToRecipe,
		RecipeSteps:                     taskSteps,
		StorageTemperatureInCelsius:     input.StorageTemperatureInCelsius,
		TimeBufferBeforeRecipeInSeconds: input.TimeBufferBeforeRecipeInSeconds,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe *recipes.Recipe, input *recipes.RecipePrepTaskStep) *recipes.RecipePrepTaskStepWithinRecipeCreationRequestInput {
	var belongsToIndex uint32
	if x := recipe.FindStepByID(input.BelongsToRecipeStep); x != nil {
		belongsToIndex = x.Index
	}

	return &recipes.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		BelongsToRecipeStepIndex: belongsToIndex,
		SatisfiesRecipeStep:      input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(input *recipes.RecipePrepTaskStep) *recipes.RecipePrepTaskStepUpdateRequestInput {
	return &recipes.RecipePrepTaskStepUpdateRequestInput{
		BelongsToRecipeStep:     pointer.To(input.BelongsToRecipeStep),
		BelongsToRecipePrepTask: pointer.To(input.BelongsToRecipePrepTask),
		SatisfiesRecipeStep:     pointer.To(input.SatisfiesRecipeStep),
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(input *recipes.RecipePrepTaskStep) *recipes.RecipePrepTaskStepDatabaseCreationInput {
	return &recipes.RecipePrepTaskStepDatabaseCreationInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}
