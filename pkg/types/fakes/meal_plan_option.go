package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

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
		ID:           ksuid.New().String(),
		Day:          time.Monday,
		AssignedCook: func(s string) *string { return &s }(BuildFakeID()),
		Meal:         *meal,
		MealName: types.MealName(fake.RandomString([]string{
			string(types.BreakfastMealName),
			string(types.SecondBreakfastMealName),
			string(types.BrunchMealName),
			string(types.LunchMealName),
			string(types.SupperMealName),
			string(types.DinnerMealName),
		})),
		Votes:             examples,
		Chosen:            false,
		Notes:             buildUniqueString(),
		CreatedOn:         uint64(uint32(fake.Date().Unix())),
		BelongsToMealPlan: fake.UUID(),
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
		Day:               &mealPlanOption.Day,
		MealID:            &mealPlanOption.Meal.ID,
		Notes:             &mealPlanOption.Notes,
		MealName:          &mealPlanOption.MealName,
		AssignedCook:      mealPlanOption.AssignedCook,
		BelongsToMealPlan: &mealPlanOption.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption builds a faked MealPlanOptionUpdateRequestInput from a meal plan option.
func BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		Day:               &mealPlanOption.Day,
		MealID:            &mealPlanOption.Meal.ID,
		Notes:             &mealPlanOption.Notes,
		MealName:          &mealPlanOption.MealName,
		AssignedCook:      mealPlanOption.AssignedCook,
		BelongsToMealPlan: &mealPlanOption.BelongsToMealPlan,
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
		ID:                mealPlanOption.ID,
		Day:               mealPlanOption.Day,
		MealID:            mealPlanOption.Meal.ID,
		Notes:             mealPlanOption.Notes,
		MealName:          mealPlanOption.MealName,
		AssignedCook:      mealPlanOption.AssignedCook,
		BelongsToMealPlan: mealPlanOption.BelongsToMealPlan,
	}
}

// BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption builds a faked MealPlanOptionDatabaseCreationInput from a meal plan option.
func BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionDatabaseCreationInput {
	return &types.MealPlanOptionDatabaseCreationInput{
		ID:                mealPlanOption.ID,
		Day:               mealPlanOption.Day,
		MealID:            mealPlanOption.Meal.ID,
		MealName:          mealPlanOption.MealName,
		Notes:             mealPlanOption.Notes,
		AssignedCook:      mealPlanOption.AssignedCook,
		BelongsToMealPlan: mealPlanOption.BelongsToMealPlan,
	}
}
