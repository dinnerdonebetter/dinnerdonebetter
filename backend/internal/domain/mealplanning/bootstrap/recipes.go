package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
)

// AllRecipes returns all bootstrap recipe creation inputs.
// Each recipe is created with the provided userID as the creator.
func AllRecipes(userID string, enums *Enumerations) []*mealplanning.RecipeDatabaseCreationInput {
	var recipes []*mealplanning.RecipeDatabaseCreationInput

	recipes = append(recipes, PanSearedButterBastedSteakRecipe(userID, enums)...)
	recipes = append(recipes, SousVideChickenBreastRecipe(userID, enums)...)
	recipes = append(recipes, PerfectRoastChickenRecipe(userID, enums)...)
	recipes = append(recipes, SousVidePorkChopsRecipe(userID, enums)...)
	recipes = append(recipes, ClassicSmashBurgersRecipe(userID, enums)...)
	recipes = append(recipes, SimpleWhiteRiceRecipe(userID, enums)...)
	recipes = append(recipes, UltraFluffyMashedPotatoesRecipe(userID, enums)...)
	recipes = append(recipes, CaesarRoastedBroccoliRecipe(userID, enums)...)
	recipes = append(recipes, HaricotsVertsAmandineRecipe(userID, enums)...)
	recipes = append(recipes, MixedGreenSaladRecipe(userID, enums)...)
	recipes = append(recipes, SoySauceBraisedChickenThighsRecipe(userID, enums)...)

	return recipes
}

