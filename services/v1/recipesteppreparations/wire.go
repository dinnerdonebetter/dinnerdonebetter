package recipesteppreparations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeStepPreparationsService,
		ProvideRecipeStepPreparationDataManager,
		ProvideRecipeStepPreparationDataServer,
	)
)

// ProvideRecipeStepPreparationDataManager turns a database into an RecipeStepPreparationDataManager.
func ProvideRecipeStepPreparationDataManager(db database.Database) models.RecipeStepPreparationDataManager {
	return db
}

// ProvideRecipeStepPreparationDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepPreparationDataServer(s *Service) models.RecipeStepPreparationDataServer {
	return s
}
