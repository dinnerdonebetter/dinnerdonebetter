package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertMealCreationRequestInputToMealDatabaseCreationInput creates a MealDatabaseCreationInput from a MealCreationRequestInput.
func ConvertMealCreationRequestInputToMealDatabaseCreationInput(input *types.MealCreationRequestInput) *types.MealDatabaseCreationInput {
	convertedComponents := []*types.MealComponentDatabaseCreationInput{}
	for _, x := range input.Components {
		convertedComponents = append(convertedComponents, ConvertMealComponentCreationRequestInputToMealComponentDatabaseCreationInput(x))
	}

	x := &types.MealDatabaseCreationInput{
		ID:                       identifiers.New(),
		Name:                     input.Name,
		Description:              input.Description,
		MinimumEstimatedPortions: input.MinimumEstimatedPortions,
		MaximumEstimatedPortions: input.MaximumEstimatedPortions,
		Components:               convertedComponents,
	}

	return x
}

// ConvertMealComponentCreationRequestInputToMealComponentDatabaseCreationInput creates a MealComponentDatabaseCreationInput from a MealComponentCreationRequestInput.
func ConvertMealComponentCreationRequestInputToMealComponentDatabaseCreationInput(input *types.MealComponentCreationRequestInput) *types.MealComponentDatabaseCreationInput {
	x := &types.MealComponentDatabaseCreationInput{
		RecipeID:      input.RecipeID,
		ComponentType: input.ComponentType,
		RecipeScale:   input.RecipeScale,
	}

	return x
}

// ConvertMealToMealCreationRequestInput builds a faked MealCreationRequestInput from a Meal.
func ConvertMealToMealCreationRequestInput(meal *types.Meal) *types.MealCreationRequestInput {
	convertedComponents := []*types.MealComponentCreationRequestInput{}
	for _, x := range meal.Components {
		convertedComponents = append(convertedComponents, ConvertMealComponentToMealComponentCreationRequestInput(x))
	}

	return &types.MealCreationRequestInput{
		Name:                     meal.Name,
		Description:              meal.Description,
		MinimumEstimatedPortions: meal.MinimumEstimatedPortions,
		MaximumEstimatedPortions: meal.MaximumEstimatedPortions,
		Components:               convertedComponents,
	}
}

// ConvertMealComponentToMealComponentCreationRequestInput creates a MealComponentCreationRequestInput from a MealComponent.
func ConvertMealComponentToMealComponentCreationRequestInput(input *types.MealComponent) *types.MealComponentCreationRequestInput {
	x := &types.MealComponentCreationRequestInput{
		RecipeID:      input.Recipe.ID,
		RecipeScale:   input.RecipeScale,
		ComponentType: input.ComponentType,
	}

	return x
}

// ConvertMealToMealDatabaseCreationInput builds a faked MealDatabaseCreationInput from a recipe.
func ConvertMealToMealDatabaseCreationInput(meal *types.Meal) *types.MealDatabaseCreationInput {
	convertedComponents := []*types.MealComponentDatabaseCreationInput{}
	for _, x := range meal.Components {
		convertedComponents = append(convertedComponents, ConvertMealComponentToMealComponentDatabaseCreationInput(x))
	}

	return &types.MealDatabaseCreationInput{
		ID:                       meal.ID,
		Name:                     meal.Name,
		Description:              meal.Description,
		MinimumEstimatedPortions: meal.MinimumEstimatedPortions,
		MaximumEstimatedPortions: meal.MaximumEstimatedPortions,
		CreatedByUser:            meal.CreatedByUser,
		Components:               convertedComponents,
	}
}

// ConvertMealComponentToMealComponentDatabaseCreationInput creates a MealComponentDatabaseCreationInput from a MealComponent.
func ConvertMealComponentToMealComponentDatabaseCreationInput(input *types.MealComponent) *types.MealComponentDatabaseCreationInput {
	x := &types.MealComponentDatabaseCreationInput{
		RecipeID:      input.Recipe.ID,
		RecipeScale:   input.RecipeScale,
		ComponentType: input.ComponentType,
	}

	return x
}

// ConvertMealToMealUpdateRequestInput builds a faked MealUpdateRequestInput from a Meal.
func ConvertMealToMealUpdateRequestInput(meal *types.Meal) *types.MealUpdateRequestInput {
	convertedComponents := []*types.MealComponentUpdateRequestInput{}
	for _, x := range meal.Components {
		convertedComponents = append(convertedComponents, ConvertMealComponentToMealComponentUpdateRequestInput(x))
	}

	return &types.MealUpdateRequestInput{
		Name:                     &meal.Name,
		Description:              &meal.Description,
		MinimumEstimatedPortions: &meal.MinimumEstimatedPortions,
		MaximumEstimatedPortions: meal.MaximumEstimatedPortions,
		CreatedByUser:            &meal.CreatedByUser,
		Components:               convertedComponents,
	}
}

// ConvertMealComponentToMealComponentUpdateRequestInput creates a MealComponentUpdateRequestInput from a MealComponent.
func ConvertMealComponentToMealComponentUpdateRequestInput(input *types.MealComponent) *types.MealComponentUpdateRequestInput {
	x := &types.MealComponentUpdateRequestInput{
		RecipeID:      &input.Recipe.ID,
		RecipeScale:   &input.RecipeScale,
		ComponentType: &input.ComponentType,
	}

	return x
}
