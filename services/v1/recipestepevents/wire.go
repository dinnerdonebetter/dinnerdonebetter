package recipestepevents

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideRecipeStepEventsService,
		ProvideRecipeStepEventDataManager,
		ProvideRecipeStepEventDataServer,
	)
)

// ProvideRecipeStepEventDataManager turns a database into an RecipeStepEventDataManager
func ProvideRecipeStepEventDataManager(db database.Database) models.RecipeStepEventDataManager {
	return db
}

// ProvideRecipeStepEventDataServer is an arbitrary function for dependency injection's sake
func ProvideRecipeStepEventDataServer(s *Service) models.RecipeStepEventDataServer {
	return s
}
