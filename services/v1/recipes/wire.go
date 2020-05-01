package recipes

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipesService,
		ProvideRecipeDataManager,
		ProvideRecipeDataServer,
	)
)

// ProvideRecipeDataManager turns a database into an RecipeDataManager.
func ProvideRecipeDataManager(db database.Database) models.RecipeDataManager {
	return db
}

// ProvideRecipeDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeDataServer(s *Service) models.RecipeDataServer {
	return s
}
