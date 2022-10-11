package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/internal/pointers"
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

func BuildFakeRecipePrepTaskStepCreationRequestInput() *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		ID:                      BuildFakeID(),
		BelongsToRecipeStep:     BuildFakeID(),
		BelongsToRecipePrepTask: BuildFakeID(),
		SatisfiesRecipeStep:     fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepCreationRequestInputFromRecipePrepTaskStep(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepCreationRequestInput {
	return &types.RecipePrepTaskStepCreationRequestInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     input.BelongsToRecipeStep,
		BelongsToRecipePrepTask: input.BelongsToRecipePrepTask,
		SatisfiesRecipeStep:     input.SatisfiesRecipeStep,
	}
}

func BuildFakeRecipePrepTaskStepUpdateRequestInput() *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		ID:                      BuildFakeID(),
		BelongsToRecipeStep:     pointers.StringPointer(BuildFakeID()),
		BelongsToRecipePrepTask: pointers.StringPointer(BuildFakeID()),
		SatisfiesRecipeStep:     pointers.BoolPointer(fake.Bool()),
	}
}

func BuildFakeRecipePrepTaskStepUpdateRequestInputFromRecipePrepTaskStep(input *types.RecipePrepTaskStep) *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		ID:                      input.ID,
		BelongsToRecipeStep:     pointers.StringPointer(input.BelongsToRecipeStep),
		BelongsToRecipePrepTask: pointers.StringPointer(input.BelongsToRecipePrepTask),
		SatisfiesRecipeStep:     pointers.BoolPointer(input.SatisfiesRecipeStep),
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

func BuildFakeRecipePrepTaskCreationRequestInput() *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepCreationRequestInput())
	}

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                                  buildUniqueString(),
		ExplicitStorageInstructions:            buildUniqueString(),
		StorageType:                            types.RecipePrepTaskStorageTypeUncovered,
		BelongsToRecipe:                        BuildFakeID(),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: fake.Uint32(),
		MinimumStorageTemperatureInCelsius:     fake.Uint32(),
		MaximumStorageTemperatureInCelsius:     fake.Uint32(),
		MinimumTimeBufferBeforeRecipeInSeconds: fake.Uint32(),
	}
}

func BuildFakeRecipePrepTaskCreationRequestInputFromRecipePrepTask(input *types.RecipePrepTask) *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepCreationRequestInputFromRecipePrepTaskStep(x))
	}

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                                  input.Notes,
		ExplicitStorageInstructions:            input.ExplicitStorageInstructions,
		StorageType:                            input.StorageType,
		BelongsToRecipe:                        input.BelongsToRecipe,
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: input.MaximumTimeBufferBeforeRecipeInSeconds,
		MinimumStorageTemperatureInCelsius:     input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     input.MaximumStorageTemperatureInCelsius,
		MinimumTimeBufferBeforeRecipeInSeconds: input.MinimumTimeBufferBeforeRecipeInSeconds,
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInput() *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInput())
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointers.StringPointer(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.StringPointer(buildUniqueString()),
		StorageType:                            pointers.StringPointer(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.StringPointer(BuildFakeID()),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Uint32Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Uint32Pointer(fake.Uint32()),
		MaximumStorageTemperatureInCelsius:     pointers.Uint32Pointer(fake.Uint32()),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Uint32Pointer(fake.Uint32()),
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInputFromRecipePrepTask(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInputFromRecipePrepTaskStep(x))
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointers.StringPointer(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.StringPointer(buildUniqueString()),
		StorageType:                            pointers.StringPointer(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.StringPointer(BuildFakeID()),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Uint32Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Uint32Pointer(fake.Uint32()),
		MaximumStorageTemperatureInCelsius:     pointers.Uint32Pointer(fake.Uint32()),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Uint32Pointer(fake.Uint32()),
	}
}
