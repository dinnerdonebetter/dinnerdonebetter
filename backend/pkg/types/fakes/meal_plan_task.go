package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeMealPlanTask builds a faked meal plan task.
func BuildFakeMealPlanTask() *types.MealPlanTask {
	return &types.MealPlanTask{
		ID:                  BuildFakeID(),
		CreatedAt:           BuildFakeTime(),
		Status:              "unfinished",
		StatusExplanation:   buildUniqueString(),
		CreationExplanation: buildUniqueString(),
		MealPlanOption:      *BuildFakeMealPlanOption(),
		CompletedAt:         nil,
		RecipePrepTask:      *BuildFakeRecipePrepTask(),
	}
}

// BuildFakeMealPlanTaskCreationRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskCreationRequestInput() *types.MealPlanTaskCreationRequestInput {
	x := BuildFakeMealPlanTask()

	return converters.ConvertMealPlanTaskToMealPlanTaskCreationRequestInput(x)
}

// BuildFakeMealPlanTasksList builds a faked MealPlanTaskList.
func BuildFakeMealPlanTasksList() *filtering.QueryFilteredResult[types.MealPlanTask] {
	var examples []*types.MealPlanTask
	for range exampleQuantity {
		examples = append(examples, BuildFakeMealPlanTask())
	}

	return &filtering.QueryFilteredResult[types.MealPlanTask]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealPlanTaskDatabaseCreationInputs builds a faked MealPlanTaskList.
func BuildFakeMealPlanTaskDatabaseCreationInputs() []*types.MealPlanTaskDatabaseCreationInput {
	var examples []*types.MealPlanTaskDatabaseCreationInput
	for range exampleQuantity {
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
		Status:            pointer.To("unfinished"),
		StatusExplanation: buildUniqueString(),
	}
}
