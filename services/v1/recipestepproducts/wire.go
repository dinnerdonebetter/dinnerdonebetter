package recipestepproducts

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideRecipeStepProductsService,
		ProvideRecipeStepProductDataManager,
		ProvideRecipeStepProductDataServer,
	)
)

// ProvideRecipeStepProductDataManager turns a database into an RecipeStepProductDataManager
func ProvideRecipeStepProductDataManager(db database.Database) models.RecipeStepProductDataManager {
	return db
}

// ProvideRecipeStepProductDataServer is an arbitrary function for dependency injection's sake
func ProvideRecipeStepProductDataServer(s *Service) models.RecipeStepProductDataServer {
	return s
}
