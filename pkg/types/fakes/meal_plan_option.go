package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeMealPlanOption builds a faked meal plan option.
func BuildFakeMealPlanOption() *types.MealPlanOption {
	return &types.MealPlanOption{
		ID:               ksuid.New().String(),
		MealPlanID:       fake.Word(),
		DayOfWeek:        fake.Uint8(),
		RecipeID:         fake.Word(),
		Notes:            fake.Word(),
		CreatedOn:        uint64(uint32(fake.Date().Unix())),
		BelongsToAccount: fake.UUID(),
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
		MealPlanID:       mealPlanOption.MealPlanID,
		DayOfWeek:        mealPlanOption.DayOfWeek,
		RecipeID:         mealPlanOption.RecipeID,
		Notes:            mealPlanOption.Notes,
		BelongsToAccount: mealPlanOption.BelongsToAccount,
	}
}

// BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption builds a faked MealPlanOptionUpdateRequestInput from a meal plan option.
func BuildFakeMealPlanOptionUpdateRequestInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionUpdateRequestInput {
	return &types.MealPlanOptionUpdateRequestInput{
		MealPlanID:       mealPlanOption.MealPlanID,
		DayOfWeek:        mealPlanOption.DayOfWeek,
		RecipeID:         mealPlanOption.RecipeID,
		Notes:            mealPlanOption.Notes,
		BelongsToAccount: mealPlanOption.BelongsToAccount,
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
		ID:               mealPlanOption.ID,
		MealPlanID:       mealPlanOption.MealPlanID,
		DayOfWeek:        mealPlanOption.DayOfWeek,
		RecipeID:         mealPlanOption.RecipeID,
		Notes:            mealPlanOption.Notes,
		BelongsToAccount: mealPlanOption.BelongsToAccount,
	}
}

// BuildFakeMealPlanOptionDatabaseCreationInput builds a faked MealPlanOptionDatabaseCreationInput.
func BuildFakeMealPlanOptionDatabaseCreationInput() *types.MealPlanOptionDatabaseCreationInput {
	mealPlanOption := BuildFakeMealPlanOption()
	return BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(mealPlanOption)
}

// BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption builds a faked MealPlanOptionDatabaseCreationInput from a meal plan option.
func BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(mealPlanOption *types.MealPlanOption) *types.MealPlanOptionDatabaseCreationInput {
	return &types.MealPlanOptionDatabaseCreationInput{
		ID:               mealPlanOption.ID,
		MealPlanID:       mealPlanOption.MealPlanID,
		DayOfWeek:        mealPlanOption.DayOfWeek,
		RecipeID:         mealPlanOption.RecipeID,
		Notes:            mealPlanOption.Notes,
		BelongsToAccount: mealPlanOption.BelongsToAccount,
	}
}
