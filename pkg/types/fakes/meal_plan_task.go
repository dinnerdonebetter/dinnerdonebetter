package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMealPlanTask builds a faked meal plan task.
func BuildFakeMealPlanTask() *types.MealPlanTask {
	now := time.Now().Add(0).Truncate(time.Second).UTC()
	inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()

	return &types.MealPlanTask{
		ID:                   BuildFakeID(),
		CannotCompleteBefore: now,
		CannotCompleteAfter:  inOneWeek,
		CreatedAt:            fake.Date(),
		Status:               "unfinished",
		StatusExplanation:    buildUniqueString(),
		CreationExplanation:  buildUniqueString(),
		CompletedAt:          nil,
	}
}

// BuildFakeMealPlanTaskCreationRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskCreationRequestInput() *types.MealPlanTaskCreationRequestInput {
	x := BuildFakeMealPlanTask()

	return BuildFakeMealPlanTaskCreationRequestInputFromMealPlanTask(x)
}

// BuildFakeMealPlanTaskCreationRequestInputFromMealPlanTask builds a faked meal plan task.
func BuildFakeMealPlanTaskCreationRequestInputFromMealPlanTask(x *types.MealPlanTask) *types.MealPlanTaskCreationRequestInput {
	return &types.MealPlanTaskCreationRequestInput{
		CannotCompleteBefore: x.CannotCompleteBefore,
		CannotCompleteAfter:  x.CannotCompleteAfter,
		Status:               x.Status,
		StatusExplanation:    x.StatusExplanation,
		CreationExplanation:  x.CreationExplanation,
	}
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
		now := time.Now().Add(0).Truncate(time.Second).UTC()
		inOneWeek := now.Add((time.Hour * 24) * 7).Add(0).Truncate(time.Second).UTC()

		examples = append(examples, &types.MealPlanTaskDatabaseCreationInput{
			MealPlanOptionID: "",
			// AppliesToRecipeStep:         "",
			ID:                   BuildFakeID(),
			CannotCompleteBefore: now,
			CannotCompleteAfter:  inOneWeek,
			StatusExplanation:    buildUniqueString(),
			CreationExplanation:  buildUniqueString(),
		})
	}

	return examples
}

// BuildFakeMealPlanTaskStatusChangeRequestInput builds a faked meal plan task.
func BuildFakeMealPlanTaskStatusChangeRequestInput() *types.MealPlanTaskStatusChangeRequestInput {
	return &types.MealPlanTaskStatusChangeRequestInput{
		ID:                BuildFakeID(),
		Status:            pointers.StringPointer("unfinished"),
		StatusExplanation: pointers.StringPointer(buildUniqueString()),
	}
}
