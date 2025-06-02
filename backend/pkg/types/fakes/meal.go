package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeMeal builds a faked meal.
func BuildFakeMeal() *types.Meal {
	recipes := []*types.MealComponent{}
	for i := 0; i < exampleQuantity; i++ {
		recipes = append(recipes, BuildFakeMealComponent())
	}

	return &types.Meal{
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
func BuildFakeMealComponent() *types.MealComponent {
	return &types.MealComponent{
		Recipe:        *BuildFakeRecipe(),
		RecipeScale:   float32(1.0),
		ComponentType: types.MealComponentTypesMain,
	}
}

// BuildFakeMealsList builds a faked MealList.
func BuildFakeMealsList() *filtering.QueryFilteredResult[types.Meal] {
	var examples []*types.Meal
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeMeal())
	}

	return &filtering.QueryFilteredResult[types.Meal]{
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
func BuildFakeMealCreationRequestInput() *types.MealCreationRequestInput {
	recipe := BuildFakeMeal()
	return converters.ConvertMealToMealCreationRequestInput(recipe)
}
