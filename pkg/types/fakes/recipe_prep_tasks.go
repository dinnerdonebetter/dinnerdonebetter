package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

func BuildFakeRecipePrepTask() *types.RecipePrepTask {
	recipePrepTaskSteps := []*types.RecipePrepTaskStep{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTaskSteps = append(recipePrepTaskSteps, BuildFakeRecipePrepTaskStep())
	}

	minStorageTemp := BuildFakeNumber()
	minTimeBuffer := BuildFakeNumber()

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
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MaximumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer + 1),
		MinimumStorageTemperatureInCelsius:     float32(minStorageTemp),
		MaximumStorageTemperatureInCelsius:     float32(minStorageTemp + 1),
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
		BelongsToRecipeStep:     BuildFakeID(),
		BelongsToRecipePrepTask: BuildFakeID(),
		SatisfiesRecipeStep:     fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepUpdateRequestInput() *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		ID:                      BuildFakeID(),
		BelongsToRecipeStep:     pointers.String(BuildFakeID()),
		BelongsToRecipePrepTask: pointers.String(BuildFakeID()),
		SatisfiesRecipeStep:     pointers.Bool(fake.Bool()),
	}
}

func BuildFakeRecipePrepTaskCreationRequestInput() *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepCreationRequestInput())
	}

	minStorageTemp := BuildFakeNumber()
	minTimeBuffer := BuildFakeNumber()

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                                  buildUniqueString(),
		ExplicitStorageInstructions:            buildUniqueString(),
		StorageType:                            types.RecipePrepTaskStorageTypeUncovered,
		BelongsToRecipe:                        BuildFakeID(),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer) + 1,
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MinimumStorageTemperatureInCelsius:     float32(minStorageTemp),
		MaximumStorageTemperatureInCelsius:     float32(minStorageTemp + 1),
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInput() *types.RecipePrepTaskUpdateRequestInput {
	minTemp := BuildFakeNumber()

	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInput())
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointers.String(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.String(buildUniqueString()),
		StorageType:                            pointers.String(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.String(BuildFakeID()),
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Uint32(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Float32(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Float32(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Uint32(fake.Uint32()),
		TaskSteps:                              taskSteps,
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInputFromRecipePrepTask(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	minTemp := BuildFakeNumber()

	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, converters.ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(x))
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointers.String(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.String(buildUniqueString()),
		StorageType:                            pointers.String(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.String(BuildFakeID()),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Uint32(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Float32(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Float32(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Uint32(fake.Uint32()),
	}
}
