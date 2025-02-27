package fakes

import (
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeMealPlan builds a faked meal plan.
func BuildFakeMealPlan() *types.MealPlan {
	mealPlanID := BuildFakeID()

	var events []*types.MealPlanEvent
	for range exampleQuantity {
		event := BuildFakeMealPlanEvent()
		event.BelongsToMealPlan = mealPlanID
		events = append(events, event)
	}

	now := time.Now().Add(30 * time.Minute).Truncate(time.Second).UTC()
	return &types.MealPlan{
		ID:                     mealPlanID,
		Notes:                  buildUniqueString(),
		Status:                 string(types.MealPlanStatusAwaitingVotes),
		VotingDeadline:         now,
		CreatedAt:              BuildFakeTime(),
		BelongsToHousehold:     fake.UUID(),
		TasksCreated:           false,
		GroceryListInitialized: false,
		ElectionMethod:         types.MealPlanElectionMethodSchulze,
		Events:                 events,
		CreatedByUser:          BuildFakeID(),
	}
}

// BuildFakeMealPlansList builds a faked MealPlanList.
func BuildFakeMealPlansList() *filtering.QueryFilteredResult[types.MealPlan] {
	var examples []*types.MealPlan
	for range exampleQuantity {
		examples = append(examples, BuildFakeMealPlan())
	}

	return &filtering.QueryFilteredResult[types.MealPlan]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
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
