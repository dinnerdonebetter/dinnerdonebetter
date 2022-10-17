package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
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
	return converters.ConvertMealPlanOptionToMealPlanOptionUpdateRequestInput(mealPlanOption)
}

// BuildFakeMealPlanOptionCreationRequestInput builds a faked MealPlanOptionCreationRequestInput.
func BuildFakeMealPlanOptionCreationRequestInput() *types.MealPlanOptionCreationRequestInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return converters.ConvertMealPlanOptionToMealPlanOptionCreationRequestInput(mealPlanOption)
}
