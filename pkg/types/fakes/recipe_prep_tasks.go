package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v6"
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
		Name:                        buildUniqueString(),
		Description:                 buildUniqueString(),
		Optional:                    fake.Bool(),
		ExplicitStorageInstructions: buildUniqueString(),
		StorageType: fake.RandomString([]string{
			types.RecipePrepTaskStorageTypeUncovered,
			types.RecipePrepTaskStorageTypeCovered,
			types.RecipePrepTaskStorageTypeAirtightContainer,
			types.RecipePrepTaskStorageTypeWireRack,
		}),
		TaskSteps:                              recipePrepTaskSteps,
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(uint32(minTimeBuffer + 1)),
		MinimumStorageTemperatureInCelsius:     pointers.Pointer(float32(minStorageTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Pointer(float32(minStorageTemp + 1)),
		BelongsToRecipe:                        BuildFakeID(),
		CreatedAt:                              BuildFakeTime(),
		LastUpdatedAt:                          nil,
		ArchivedAt:                             nil,
	}
}

func BuildFakeRecipePrepTaskList() *types.QueryFilteredResult[types.RecipePrepTask] {
	recipePrepTasks := &types.QueryFilteredResult[types.RecipePrepTask]{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTasks.Data = append(recipePrepTasks.Data, BuildFakeRecipePrepTask())
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
		BelongsToRecipeStep: BuildFakeID(),
		SatisfiesRecipeStep: fake.Bool(),
	}
}

func BuildFakeRecipePrepTaskStepUpdateRequestInput() *types.RecipePrepTaskStepUpdateRequestInput {
	return &types.RecipePrepTaskStepUpdateRequestInput{
		BelongsToRecipeStep:     pointers.Pointer(BuildFakeID()),
		BelongsToRecipePrepTask: pointers.Pointer(BuildFakeID()),
		SatisfiesRecipeStep:     pointers.Pointer(fake.Bool()),
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
		Name:                                   buildUniqueString(),
		Optional:                               fake.Bool(),
		Description:                            buildUniqueString(),
		StorageType:                            types.RecipePrepTaskStorageTypeUncovered,
		BelongsToRecipe:                        BuildFakeID(),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(uint32(minTimeBuffer) + 1),
		MinimumTimeBufferBeforeRecipeInSeconds: uint32(minTimeBuffer),
		MinimumStorageTemperatureInCelsius:     pointers.Pointer(float32(minStorageTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Pointer(float32(minStorageTemp + 1)),
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInput() *types.RecipePrepTaskUpdateRequestInput {
	minTemp := BuildFakeNumber()

	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInput())
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  pointers.Pointer(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.Pointer(buildUniqueString()),
		Name:                                   pointers.Pointer(buildUniqueString()),
		Description:                            pointers.Pointer(buildUniqueString()),
		Optional:                               pointers.Pointer(fake.Bool()),
		StorageType:                            pointers.Pointer(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.Pointer(BuildFakeID()),
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Pointer(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Pointer(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(fake.Uint32()),
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
		Notes:                                  pointers.Pointer(buildUniqueString()),
		ExplicitStorageInstructions:            pointers.Pointer(buildUniqueString()),
		Name:                                   pointers.Pointer(buildUniqueString()),
		Description:                            pointers.Pointer(buildUniqueString()),
		Optional:                               pointers.Pointer(fake.Bool()),
		StorageType:                            pointers.Pointer(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                        pointers.Pointer(BuildFakeID()),
		TaskSteps:                              taskSteps,
		MaximumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(fake.Uint32()),
		MinimumStorageTemperatureInCelsius:     pointers.Pointer(float32(minTemp)),
		MaximumStorageTemperatureInCelsius:     pointers.Pointer(float32(minTemp + 1)),
		MinimumTimeBufferBeforeRecipeInSeconds: pointers.Pointer(fake.Uint32()),
	}
}
