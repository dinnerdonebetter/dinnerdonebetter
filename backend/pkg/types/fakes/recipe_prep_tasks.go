package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

func BuildFakeRecipePrepTask() *types.RecipePrepTask {
	recipePrepTaskSteps := []*types.RecipePrepTaskStep{}
	for i := 0; i < exampleQuantity; i++ {
		recipePrepTaskSteps = append(recipePrepTaskSteps, BuildFakeRecipePrepTaskStep())
	}

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
		TaskSteps:                       recipePrepTaskSteps,
		StorageTemperatureInCelsius:     BuildFakeOptionalFloat32Range(),
		TimeBufferBeforeRecipeInSeconds: BuildFakeUint32RangeWithOptionalMax(),
		BelongsToRecipe:                 BuildFakeID(),
		CreatedAt:                       BuildFakeTime(),
		LastUpdatedAt:                   nil,
		ArchivedAt:                      nil,
	}
}

func BuildFakeRecipePrepTasksList() *types.QueryFilteredResult[types.RecipePrepTask] {
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
		BelongsToRecipeStep:     pointer.To(BuildFakeID()),
		BelongsToRecipePrepTask: pointer.To(BuildFakeID()),
		SatisfiesRecipeStep:     pointer.To(fake.Bool()),
	}
}

func BuildFakeRecipePrepTaskCreationRequestInput() *types.RecipePrepTaskCreationRequestInput {
	taskSteps := []*types.RecipePrepTaskStepCreationRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepCreationRequestInput())
	}

	return &types.RecipePrepTaskCreationRequestInput{
		Notes:                           buildUniqueString(),
		ExplicitStorageInstructions:     buildUniqueString(),
		Name:                            buildUniqueString(),
		Optional:                        fake.Bool(),
		Description:                     buildUniqueString(),
		StorageType:                     types.RecipePrepTaskStorageTypeUncovered,
		BelongsToRecipe:                 BuildFakeID(),
		RecipeSteps:                     taskSteps,
		TimeBufferBeforeRecipeInSeconds: BuildFakeUint32RangeWithOptionalMax(),
		StorageTemperatureInCelsius:     BuildFakeOptionalFloat32Range(),
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInput() *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for i := 0; i < exampleQuantity; i++ {
		taskSteps = append(taskSteps, BuildFakeRecipePrepTaskStepUpdateRequestInput())
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                           pointer.To(buildUniqueString()),
		ExplicitStorageInstructions:     pointer.To(buildUniqueString()),
		Name:                            pointer.To(buildUniqueString()),
		Description:                     pointer.To(buildUniqueString()),
		Optional:                        pointer.To(fake.Bool()),
		StorageType:                     pointer.To(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                 pointer.To(BuildFakeID()),
		TimeBufferBeforeRecipeInSeconds: BuildFakeUint32RangeWithOptionalMaxUpdateRequestInput(),
		StorageTemperatureInCelsius:     BuildFakeOptionalFloat32Range(),
		TaskSteps:                       taskSteps,
	}
}

func BuildFakeRecipePrepTaskUpdateRequestInputFromRecipePrepTask(input *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	taskSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, x := range input.TaskSteps {
		taskSteps = append(taskSteps, converters.ConvertRecipePrepTaskStepToRecipePrepTaskStepUpdateRequestInput(x))
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                           pointer.To(buildUniqueString()),
		ExplicitStorageInstructions:     pointer.To(buildUniqueString()),
		Name:                            pointer.To(buildUniqueString()),
		Description:                     pointer.To(buildUniqueString()),
		Optional:                        pointer.To(fake.Bool()),
		StorageType:                     pointer.To(types.RecipePrepTaskStorageTypeUncovered),
		BelongsToRecipe:                 pointer.To(BuildFakeID()),
		TaskSteps:                       taskSteps,
		TimeBufferBeforeRecipeInSeconds: BuildFakeUint32RangeWithOptionalMaxUpdateRequestInput(),
		StorageTemperatureInCelsius:     BuildFakeOptionalFloat32Range(),
	}
}
