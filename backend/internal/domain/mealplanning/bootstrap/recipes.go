package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
)

// AllRecipes returns all bootstrap recipe creation inputs.
// Each recipe is created with the provided userID as the creator.
func AllRecipes(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	return []*mealplanning.RecipeDatabaseCreationInput{
		PanSearedButterBastedSteakRecipe(userID, enums),
		SousVideChickenBreastRecipe(userID, enums),
		PerfectRoastChickenRecipe(userID, enums),
	}
}

