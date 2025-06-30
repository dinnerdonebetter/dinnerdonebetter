package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
)

// BuildFakeMeal builds a faked meal.
func BuildFakeMeal() *mealplanning.Meal {
	recipes := []*mealplanning.MealComponent{}
	for i := 0; i < exampleQuantity; i++ {
		recipes = append(recipes, BuildFakeMealComponent())
	}

	return &mealplanning.Meal{
		ID:          BuildFakeID(),
		Name:        buildUniqueString(),
		Description: buildUniqueString(),
		EstimatedPortions: types.Float32RangeWithOptionalMax{
			Min: float32(buildFakeNumber()),
			Max: nil,
		},
		CreatedAt:            BuildFakeTime(),
		CreatedByUser:        BuildFakeID(),
		Components:           recipes,
		EligibleForMealPlans: true,
	}
}

// BuildFakeMealComponent builds a faked meal component.
func BuildFakeMealComponent() *mealplanning.MealComponent {
	return &mealplanning.MealComponent{
		Recipe:        *BuildFakeRecipe(),
		RecipeScale:   float32(1.0),
		ComponentType: mealplanning.MealComponentTypesMain,
	}
}

// BuildFakeMealsList builds a faked MealList.
func BuildFakeMealsList() *filtering.QueryFilteredResult[mealplanning.Meal] {
	var examples []*mealplanning.Meal
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMeal())
	}

	return &filtering.QueryFilteredResult[mealplanning.Meal]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeMealCreationRequestInput builds a faked MealCreationRequestInput.
func BuildFakeMealCreationRequestInput() *mealplanning.MealCreationRequestInput {
	recipe := BuildFakeMeal()
	return converters.ConvertMealToMealCreationRequestInput(recipe)
}
