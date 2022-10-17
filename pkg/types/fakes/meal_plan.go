package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
)

// BuildFakeMealPlan builds a faked meal plan.
func BuildFakeMealPlan() *types.MealPlan {
	mealPlanID := BuildFakeID()

	var events []*types.MealPlanEvent
	for i := 0; i < exampleQuantity; i++ {
		event := BuildFakeMealPlanEvent()
		event.BelongsToMealPlan = mealPlanID
		events = append(events, event)
	}

	now := time.Now().Add(30 * time.Minute)
	return &types.MealPlan{
		ID:                     mealPlanID,
		Notes:                  buildUniqueString(),
		Status:                 types.AwaitingVotesMealPlanStatus,
		VotingDeadline:         now,
		CreatedAt:              fake.Date(),
		BelongsToHousehold:     fake.UUID(),
		TasksCreated:           false,
		GroceryListInitialized: false,
		Events:                 events,
	}
}

// BuildFakeMealPlanList builds a faked MealPlanList.
func BuildFakeMealPlanList() *types.MealPlanList {
	var examples []*types.MealPlan
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlan())
	}

	return &types.MealPlanList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlans: examples,
	}
}

// BuildFakeMealPlanUpdateRequestInput builds a faked MealPlanUpdateRequestInput from a meal plan.
func BuildFakeMealPlanUpdateRequestInput() *types.MealPlanUpdateRequestInput {
	mealPlan := BuildFakeMealPlan()
	return converters.ConvertMealPlanToMealPlanUpdateRequestInput(mealPlan)
}

// BuildFakeMealPlanCreationRequestInput builds a faked MealPlanCreationRequestInput.
func BuildFakeMealPlanCreationRequestInput() *types.MealPlanCreationRequestInput {
	mealPlan := BuildFakeMealPlan()
	return converters.ConvertMealPlanToMealPlanCreationRequestInput(mealPlan)
}
