package recipeiterationsteps

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeIterationStepsService,
		ProvideRecipeIterationStepDataManager,
		ProvideRecipeIterationStepDataServer,
	)
)

// ProvideRecipeIterationStepDataManager turns a database into an RecipeIterationStepDataManager.
func ProvideRecipeIterationStepDataManager(db database.Database) models.RecipeIterationStepDataManager {
	return db
}

// ProvideRecipeIterationStepDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeIterationStepDataServer(s *Service) models.RecipeIterationStepDataServer {
	return s
}
