package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

func BuildFakeRecipePrepTask() *types.RecipePrepTask {
	recipePrepTaskSteps := []*types.RecipePrepTaskStep{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTaskSteps = append(recipePrepTaskSteps, BuildFakeRecipePrepTaskStep())
	}

	return &types.RecipePrepTask{
		ID:                          BuildFakeID(),
		Notes:                       buildUniqueString(),
		ExplicitStorageInstructions: buildUniqueString(),
		StorageType: fake.RandomString([]string{
			types.RecipePrepTaskStorageTypeUncovered,
			types.RecipePrepTaskStorageTypeCovered,
			types.RecipePrepTaskStorageTypeAirtightContainer,
			types.RecipePrepTaskStorageTypeWireRack,
		}),
		TaskSteps:                              recipePrepTaskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: fake.Uint32(),
		MaximumTimeBufferBeforeRecipeInSeconds: fake.Uint32(),
		MinimumStorageTemperatureInCelsius:     fake.Uint32(),
		MaximumStorageTemperatureInCelsius:     fake.Uint32(),
		BelongsToRecipe:                        BuildFakeID(),
		CreatedAt:                              fake.Date(),
		LastUpdatedAt:                          nil,
		ArchivedAt:                             nil,
	}
}

func BuildFakeRecipePrepTaskList() *types.RecipePrepTaskList {
	recipePrepTasks := &types.RecipePrepTaskList{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTasks.RecipePrepTasks = append(recipePrepTasks.RecipePrepTasks, BuildFakeRecipePrepTask())
	}

	return recipePrepTasks
}

// BuildFakeRecipePrepTaskDatabaseCreationInputFromRecipePrepTask builds a faked RecipePrepTaskDatabaseCreationInput from a recipe prep task.
func BuildFakeRecipePrepTaskDatabaseCreationInputFromRecipePrepTask(recipePrepTask *types.RecipePrepTask) *types.RecipePrepTaskDatabaseCreationInput {
	taskSteps := []*types.RecipePrepTaskStepDatabaseCreationInput{}
	for _, step := range recipePrepTask.TaskSteps {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepDatabaseCreationInputFromRecipePrepTaskStep(step))
	}

	return &types.RecipePrepTaskDatabaseCreationInput{
		ID:                                     recipePrepTask.ID,
		Notes:                                  recipePrepTask.Notes,
		ExplicitStorageInstructions:            recipePrepTask.ExplicitStorageInstructions,
		StorageType:                            recipePrepTask.StorageType,
		TaskSteps:                              taskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: recipePrepTask.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: recipePrepTask.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     recipePrepTask.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     recipePrepTask.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        recipePrepTask.BelongsToRecipe,
	}
}

func BuildFakeRecipePrepTaskStep() *types.RecipePrepTaskStep {
	return &types.RecipePrepTaskStep{
		ID:                      BuildFakeID(),
		BelongsToRecipeStep:     BuildFakeID(),
		BelongsToRecipePrepTask: BuildFakeID(),
		SatisfiesRecipeStep:     fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepDatabaseCreationInputFromRecipePrepTaskStep(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepDatabaseCreationInput {
	return &types.RecipePrepTaskStepDatabaseCreationInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}
