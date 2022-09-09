package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMealPlanEvent builds a faked meal plan event.
func BuildFakeMealPlanEvent() *types.MealPlanEvent {
	mealPlanID := ksuid.New().String()

	now := time.Now().Add(0)
	inTenMinutes := now.Add(time.Minute * 10)
	inOneWeek := now.Add((time.Hour * 24) * 7)

	return &types.MealPlanEvent{
		ID:                mealPlanID,
		Notes:             buildUniqueString(),
		StartsAt:          inTenMinutes,
		EndsAt:            inOneWeek,
		CreatedAt:         fake.Date(),
		BelongsToMealPlan: BuildFakeID(),
		Options:           BuildFakeMealPlanOptionList().MealPlanOptions,
	}
}

// BuildFakeMealPlanEventList builds a faked MealPlanEventList.
func BuildFakeMealPlanEventList() *types.MealPlanEventList {
	var examples []*types.MealPlanEvent
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanEvent())
	}

	return &types.MealPlanEventList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlanEvents: examples,
	}
}

// BuildFakeMealPlanEventUpdateRequestInput builds a faked MealPlanEventUpdateRequestInput from a meal plan.
func BuildFakeMealPlanEventUpdateRequestInput() *types.MealPlanEventUpdateRequestInput {
	mealPlan := BuildFakeMealPlanEvent()
	return &types.MealPlanEventUpdateRequestInput{
		Notes:             &mealPlan.Notes,
		StartsAt:          &mealPlan.StartsAt,
		EndsAt:            &mealPlan.EndsAt,
		BelongsToMealPlan: mealPlan.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanEventUpdateRequestInputFromMealPlanEvent builds a faked MealPlanEventUpdateRequestInput from a meal plan.
func BuildFakeMealPlanEventUpdateRequestInputFromMealPlanEvent(mealPlan *types.MealPlanEvent) *types.MealPlanEventUpdateRequestInput {
	return &types.MealPlanEventUpdateRequestInput{
		Notes:             &mealPlan.Notes,
		StartsAt:          &mealPlan.StartsAt,
		EndsAt:            &mealPlan.EndsAt,
		BelongsToMealPlan: mealPlan.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanEventCreationRequestInput builds a faked MealPlanEventCreationRequestInput.
func BuildFakeMealPlanEventCreationRequestInput() *types.MealPlanEventCreationRequestInput {
	mealPlan := BuildFakeMealPlanEvent()
	return BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent(mealPlan)
}

// BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent builds a faked MealPlanEventCreationRequestInput from a meal plan.
func BuildFakeMealPlanEventCreationRequestInputFromMealPlanEvent(mealPlan *types.MealPlanEvent) *types.MealPlanEventCreationRequestInput {
	return &types.MealPlanEventCreationRequestInput{
		ID:                mealPlan.ID,
		Notes:             mealPlan.Notes,
		StartsAt:          mealPlan.StartsAt,
		EndsAt:            mealPlan.EndsAt,
		BelongsToMealPlan: mealPlan.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanEventDatabaseCreationInputFromMealPlanEvent builds a faked MealPlanEventDatabaseCreationInput from a meal plan.
func BuildFakeMealPlanEventDatabaseCreationInputFromMealPlanEvent(mealPlan *types.MealPlanEvent) *types.MealPlanEventDatabaseCreationInput {
	return &types.MealPlanEventDatabaseCreationInput{
		ID:                mealPlan.ID,
		Notes:             mealPlan.Notes,
		StartsAt:          mealPlan.StartsAt,
		EndsAt:            mealPlan.EndsAt,
		BelongsToMealPlan: mealPlan.BelongsToMealPlan,
	}
}
