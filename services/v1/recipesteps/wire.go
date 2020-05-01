package recipesteps

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeStepsService,
		ProvideRecipeStepDataManager,
		ProvideRecipeStepDataServer,
	)
)

// ProvideRecipeStepDataManager turns a database into an RecipeStepDataManager.
func ProvideRecipeStepDataManager(db database.Database) models.RecipeStepDataManager {
	return db
}

// ProvideRecipeStepDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepDataServer(s *Service) models.RecipeStepDataServer {
	return s
}
