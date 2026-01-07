package bootstrap

import (
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
)

// AllRecipes returns all bootstrap recipe creation inputs.
// Each recipe is created with the provided userID as the creator.
// Note: PanSearedButterBastedSteakRecipe is excluded as it's created separately using RecipeManager.
func AllRecipes(enums *Enumerations) []*mealplanning.RecipeCreationRequestInput {
	var recipes []*mealplanning.RecipeCreationRequestInput

	// working recipes here

	// recipes = append(recipes, PanSearedButterBastedSteakRecipe(enums)...)
	// recipes = append(recipes, SousVideChickenBreastRecipe(enums)...)
	// recipes = append(recipes, PerfectRoastChickenRecipe(enums)...)
	// recipes = append(recipes, SousVidePorkChopsRecipe(enums)...)
	// recipes = append(recipes, ClassicSmashBurgersRecipe(enums)...)
	// recipes = append(recipes, SimpleWhiteRiceRecipe(enums)...)
	// recipes = append(recipes, UltraFluffyMashedPotatoesRecipe(enums)...)
	// recipes = append(recipes, CaesarRoastedBroccoliRecipe(enums)...)
	// recipes = append(recipes, HaricotsVertsAmandineRecipe(enums)...)
	// recipes = append(recipes, MixedGreenSaladRecipe(enums)...)
	// recipes = append(recipes, RoastedBrusselsSproutsRecipe(enums)...)
	// recipes = append(recipes, StovetopMacAndCheeseRecipe(enums)...)
	// recipes = append(recipes, CaesarSaladRecipe(enums)...)
	// recipes = append(recipes, GlazedCarrotsWithBrownButterAndSageRecipe(enums)...)
	// recipes = append(recipes, StirFriedGreenBeansRecipe(enums)...)
	// recipes = append(recipes, TortillasRecipe(enums)...)

	// non-working recipes here

	recipes = append(recipes, SoySauceBraisedChickenThighsRecipe(enums)...)
	//recipes = append(recipes, GrilledPorkTenderloinRecipe(enums)...)
	//recipes = append(recipes, PanSearedSalmonFilletsRecipe(enums)...)
	//recipes = append(recipes, RefriedBeansRecipe(enums)...)
	//recipes = append(recipes, CarneAsadaRecipe(enums)...)
	//recipes = append(recipes, ButterChickenRecipe(enums)...)
	//recipes = append(recipes, CornbreadRecipe(enums)...)
	//recipes = append(recipes, GrilledWholeCauliflowerRecipe(enums)...)

	return recipes
}
