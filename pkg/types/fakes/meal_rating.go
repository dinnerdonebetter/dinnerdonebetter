package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeMealRating builds a faked valid ingredient.
func BuildFakeMealRating() *types.MealRating {
	var x *types.MealRating
	fake.Struct(&x)
	return x
}

// BuildFakeMealRatingList builds a faked MealRatingList.
func BuildFakeMealRatingList() *types.QueryFilteredResult[types.MealRating] {
	var examples []*types.MealRating
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMealRating())
	}

	return &types.QueryFilteredResult[types.MealRating]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealRatingUpdateRequestInput builds a faked MealRatingUpdateRequestInput from a valid ingredient.
func BuildFakeMealRatingUpdateRequestInput() *types.MealRatingUpdateRequestInput {
	mealRating := BuildFakeMealRating()
	return converters.ConvertMealRatingToMealRatingUpdateRequestInput(mealRating)
}

// BuildFakeMealRatingCreationRequestInput builds a faked MealRatingCreationRequestInput.
func BuildFakeMealRatingCreationRequestInput() *types.MealRatingCreationRequestInput {
	mealRating := BuildFakeMealRating()
	return converters.ConvertMealRatingToMealRatingCreationRequestInput(mealRating)
}
