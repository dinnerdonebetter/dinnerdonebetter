package converters

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput creates a RecipePrepTaskUpdateRequestInput from a RecipePrepTask.
func ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		y := x
		taskSteps = append(taskSteps, &types.RecipePrepTaskStepUpdateRequestInput{
			BelongsToRecipeStep:     &y.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &y.BelongsToRecipePrepTask,
			SatisfiesRecipeStep:     &y.SatisfiesRecipeStep,
		})
	}
	x := &types.RecipePrepTaskUpdateRequestInput{
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
func ConvertRecipePrepTaskCreationRequestInputToRecipePrepTaskDatabaseCreationInput(input *types.RecipePrepTaskCreationRequestInput) *types.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for _, x := range input.RecipeSteps {
		taskSteps = append(taskSteps, &types.RecipePrepTaskStepDatabaseCreationInput{
			BelongsToRecipeStep: x.BelongsToRecipeStep,
			SatisfiesRecipeStep: x.SatisfiesRecipeStep,
		})
	}

	x := &types.RecipePrepTaskDatabaseCreationInput{
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
func ConvertRecipePrepTaskWithinRecipeCreationRequestInputToRecipePrepTaskDatabaseCreationInput(recipe *types.RecipeDatabaseCreationInput, input *types.RecipePrepTaskWithinRecipeCreationRequestInput) (*types.RecipePrepTaskDatabaseCreationInput, error) {
	x := &types.RecipePrepTaskDatabaseCreationInput{
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

	x.TaskSteps = []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for i, ts := range input.RecipeSteps {
		if rs := recipe.FindStepByIndex(ts.BelongsToRecipeStepIndex); rs != nil {
			x.TaskSteps = append(x.TaskSteps, &types.RecipePrepTaskStepDatabaseCreationInput{
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
func ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(input *types.RecipePrepTask) *types.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for _, step := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(step))
	}

	return &types.RecipePrepTaskDatabaseCreationInput{
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

func ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		SatisfiesRecipeStep: input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(input *types.RecipePrepTask) *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepCreationRequestInput(x))
	}

	return &types.RecipePrepTaskCreationRequestInput{
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

func ConvertRecipePrepTaskToRecipePrepTaskWithinRecipeCreationRequestInput(recipe *types.Recipe, input *types.RecipePrepTask) *types.RecipePrepTaskWithinRecipeCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepWithinRecipeCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe, x))
	}

	return &types.RecipePrepTaskWithinRecipeCreationRequestInput{
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

func ConvertRecipePrepTaskStepToRecipePrepTaskStepWithinRecipeCreationRequestInput(recipe *types.Recipe, input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepWithinRecipeCreationRequestInput {
	var belongsToIndex uint32
	if x := recipe.FindStepByID(input.BelongsToRecipeStep); x != nil {
		belongsToIndex = x.Index
	}

	return &types.RecipePrepTaskStepWithinRecipeCreationRequestInput{
		BelongsToRecipeStepIndex: belongsToIndex,
		SatisfiesRecipeStep:      input.SatisfiesRecipeStep,
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		BelongsToRecipeStep:     pointer.To(input.BelongsToRecipeStep),
		BelongsToRecipePrepTask: pointer.To(input.BelongsToRecipePrepTask),
		SatisfiesRecipeStep:     pointer.To(input.SatisfiesRecipeStep),
	}
}

func ConvertRecipePrepTaskStepToRecipePrepTaskStepDatabaseCreationInput(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepDatabaseCreationInput {
	return &types.RecipePrepTaskStepDatabaseCreationInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}
