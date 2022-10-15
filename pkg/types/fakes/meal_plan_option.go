package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeMealPlanOption builds a faked meal plan option.
func BuildFakeMealPlanOption() *types.MealPlanOption {
	var examples []*types.MealPlanOptionVote
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanOptionVote())
	}

	meal := BuildFakeMeal()
	meal.Recipes = nil

	return &types.MealPlanOption{
		ID:                     BuildFakeID(),
		AssignedCook:           func(s string) *string { return &s }(BuildFakeID()),
		Meal:                   *meal,
		Votes:                  examples,
		Chosen:                 false,
		Notes:                  buildUniqueString(),
		CreatedAt:              fake.Date(),
		BelongsToMealPlanEvent: fake.UUID(),
		PrepStepsCreated:       false,
	}
}

// BuildFakeMealPlanOptionList builds a faked MealPlanOptionList.
func BuildFakeMealPlanOptionList() *types.MealPlanOptionList {
	var examples []*types.MealPlanOption
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealPlanOption())
	}

	return &types.MealPlanOptionList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		MealPlanOptions: examples,
	}
}

// BuildFakeMealPlanOptionUpdateRequestInput builds a faked MealPlanOptionUpdateRequestInput from a meal plan option.
func BuildFakeMealPlanOptionUpdateRequestInput() *types.MealPlanOptionUpdateRequestInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return &types.MealPlanOptionUpdateRequestInput{
		MealID:                 &mealPlanOption.Meal.ID,
		Notes:                  &mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: &mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       &mealPlanOption.PrepStepsCreated,
	}
}

// BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption builds a faked MealPlanOptionUpdateRequestInput from a meal plan option.
func BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		MealID:                 &mealPlanOption.Meal.ID,
		Notes:                  &mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: &mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       &mealPlanOption.PrepStepsCreated,
	}
}

// BuildFakeMealPlanOptionCreationRequestInput builds a faked MealPlanOptionCreationRequestInput.
func BuildFakeMealPlanOptionCreationRequestInput() *types.MealPlanOptionCreationRequestInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(mealPlanOption)
}

// BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption builds a faked MealPlanOptionCreationRequestInput from a meal plan option.
func BuildFakeMealPlanOptionCreationRequestInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionCreationRequestInput {
	return &types.MealPlanOptionCreationRequestInput{
		ID:                     mealPlanOption.ID,
		MealID:                 mealPlanOption.Meal.ID,
		Notes:                  mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       mealPlanOption.PrepStepsCreated,
	}
}

// BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption builds a faked MealPlanOptionDatabaseCreationInput from a meal plan option.
func BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionDatabaseCreationInput {
	return &types.MealPlanOptionDatabaseCreationInput{
		ID:                     mealPlanOption.ID,
		MealID:                 mealPlanOption.Meal.ID,
		Notes:                  mealPlanOption.Notes,
		AssignedCook:           mealPlanOption.AssignedCook,
		AssignedDishwasher:     mealPlanOption.AssignedDishwasher,
		BelongsToMealPlanEvent: mealPlanOption.BelongsToMealPlanEvent,
		PrepStepsCreated:       mealPlanOption.PrepStepsCreated,
	}
}
