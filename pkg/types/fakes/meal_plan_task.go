package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeMealPlanTask builds a faked meal plan task.
func BuildFakeMealPlanTask() *types.MealPlanTask {
	return &types.MealPlanTask{
		ID:                  BuildFakeID(),
		CreatedAt:           fake.Date(),
		Status:              "unfinished",
		StatusExplanation:   buildUniqueString(),
		CreationExplanation: buildUniqueString(),
		CompletedAt:         nil,
	}
}

// BuildFakeMealPlanTaskCreationRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskCreationRequestInput() *types.MealPlanTaskCreationRequestInput {
	x := BuildFakeMealPlanTask()

	return converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(x)
}

// BuildFakeMealPlanTaskList builds a faked MealPlanTaskList.
func BuildFakeMealPlanTaskList() *types.MealPlanTaskList {
	var examples []*types.MealPlanTask
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanTask())
	}

	return &types.MealPlanTaskList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlanTasks: examples,
	}
}

// BuildFakeMealPlanTaskDatabaseCreationInputs builds a faked MealPlanTaskList.
func BuildFakeMealPlanTaskDatabaseCreationInputs() []*types.MealPlanTaskDatabaseCreationInput {
	var examples []*types.MealPlanTaskDatabaseCreationInput
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, &types.MealPlanTaskDatabaseCreationInput{
			MealPlanOptionID:    "",
			ID:                  BuildFakeID(),
			StatusExplanation:   buildUniqueString(),
			CreationExplanation: buildUniqueString(),
		})
	}

	return examples
}

// BuildFakeMealPlanTaskStatusChangeRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskStatusChangeRequestInput() *types.MealPlanTaskStatusChangeRequestInput {
	return &types.MealPlanTaskStatusChangeRequestInput{
		ID:                BuildFakeID(),
		Status:            pointers.String("unfinished"),
		StatusExplanation: pointers.String(buildUniqueString()),
	}
}
