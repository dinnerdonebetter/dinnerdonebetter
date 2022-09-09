package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMealPlan builds a faked meal plan.
func BuildFakeMealPlan() *types.MealPlan {
	mealPlanID := ksuid.New().String()

	var events []*types.MealPlanEvent
	for i := 0; i < exampleQuantity; i++ {
		option := BuildFakeMealPlanEvent()
		option.BelongsToMealPlan = mealPlanID
		events = append(events, option)
	}

	now := time.Now().Add(0)
	return &types.MealPlan{
		ID:                 mealPlanID,
		Notes:              buildUniqueString(),
		Status:             types.AwaitingVotesMealPlanStatus,
		VotingDeadline:     now,
		CreatedAt:          fake.Date(),
		BelongsToHousehold: fake.UUID(),
		Events:             events,
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
	return &types.MealPlanUpdateRequestInput{
		Notes:              &mealPlan.Notes,
		VotingDeadline:     &mealPlan.VotingDeadline,
		BelongsToHousehold: &mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanUpdateRequestInputFromMealPlan builds a faked MealPlanUpdateRequestInput from a meal plan.
func BuildFakeMealPlanUpdateRequestInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanUpdateRequestInput {
	return &types.MealPlanUpdateRequestInput{
		Notes:              &mealPlan.Notes,
		VotingDeadline:     &mealPlan.VotingDeadline,
		BelongsToHousehold: &mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanCreationRequestInput builds a faked MealPlanCreationRequestInput.
func BuildFakeMealPlanCreationRequestInput() *types.MealPlanCreationRequestInput {
	mealPlan := BuildFakeMealPlan()
	return BuildFakeMealPlanCreationRequestInputFromMealPlan(mealPlan)
}

// BuildFakeMealPlanCreationRequestInputFromMealPlan builds a faked MealPlanCreationRequestInput from a meal plan.
func BuildFakeMealPlanCreationRequestInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanCreationRequestInput {
	events := []*types.MealPlanEventCreationRequestInput{}
	for _, opt := range mealPlan.Events {
		events = append(events, BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent(opt))
	}

	return &types.MealPlanCreationRequestInput{
		ID:                 mealPlan.ID,
		Notes:              mealPlan.Notes,
		VotingDeadline:     mealPlan.VotingDeadline,
		Events:             events,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}

// BuildFakeMealPlanDatabaseCreationInputFromMealPlan builds a faked MealPlanDatabaseCreationInput from a meal plan.
func BuildFakeMealPlanDatabaseCreationInputFromMealPlan(mealPlan *types.MealPlan) *types.MealPlanDatabaseCreationInput {
	events := []*types.MealPlanEventDatabaseCreationInput{}
	for _, opt := range mealPlan.Events {
		events = append(events, BuildFakeMealPlanEventDatabaseCreationInputFromMealPlanEvent(opt))
	}

	return &types.MealPlanDatabaseCreationInput{
		ID:                 mealPlan.ID,
		Notes:              mealPlan.Notes,
		VotingDeadline:     mealPlan.VotingDeadline,
		Events:             events,
		BelongsToHousehold: mealPlan.BelongsToHousehold,
	}
}
