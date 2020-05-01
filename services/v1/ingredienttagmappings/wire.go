package ingredienttagmappings

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideIngredientTagMappingsService,
		ProvideIngredientTagMappingDataManager,
		ProvideIngredientTagMappingDataServer,
	)
)

// ProvideIngredientTagMappingDataManager turns a database into an IngredientTagMappingDataManager.
func ProvideIngredientTagMappingDataManager(db database.Database) models.IngredientTagMappingDataManager {
	return db
}

// ProvideIngredientTagMappingDataServer is an arbitrary function for dependency injection's sake.
func ProvideIngredientTagMappingDataServer(s *Service) models.IngredientTagMappingDataServer {
	return s
}
