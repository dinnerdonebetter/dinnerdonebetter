package recipeiterations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeIterationsService,
		ProvideRecipeIterationDataManager,
		ProvideRecipeIterationDataServer,
	)
)

// ProvideRecipeIterationDataManager turns a database into an RecipeIterationDataManager.
func ProvideRecipeIterationDataManager(db database.DataManager) models.RecipeIterationDataManager {
	return db
}

// ProvideRecipeIterationDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeIterationDataServer(s *Service) models.RecipeIterationDataServer {
	return s
}
