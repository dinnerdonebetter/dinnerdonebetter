package recipetags

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideRecipeTagsService,
		ProvideRecipeTagDataManager,
		ProvideRecipeTagDataServer,
	)
)

// ProvideRecipeTagDataManager turns a database into an RecipeTagDataManager.
func ProvideRecipeTagDataManager(db database.Database) models.RecipeTagDataManager {
	return db
}

// ProvideRecipeTagDataServer is an arbitrary function for dependency injection's sake.
func ProvideRecipeTagDataServer(s *Service) models.RecipeTagDataServer {
	return s
}
