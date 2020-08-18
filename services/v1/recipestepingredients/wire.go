package recipestepingredients

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeStepIngredientsService,
		ProvideRecipeStepIngredientDataManager,
		ProvideRecipeStepIngredientDataServer,
	)
)

// ProvideRecipeStepIngredientDataManager turns a database into an RecipeStepIngredientDataManager.
func ProvideRecipeStepIngredientDataManager(db database.DataManager) models.RecipeStepIngredientDataManager {
	return db
}

// ProvideRecipeStepIngredientDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepIngredientDataServer(s *Service) models.RecipeStepIngredientDataServer {
	return s
}
