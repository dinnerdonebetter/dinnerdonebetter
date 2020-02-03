package ingredients

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideIngredientsService,
		ProvideIngredientDataManager,
		ProvideIngredientDataServer,
	)
)

// ProvideIngredientDataManager turns a database into an IngredientDataManager
func ProvideIngredientDataManager(db database.Database) models.IngredientDataManager {
	return db
}

// ProvideIngredientDataServer is an arbitrary function for dependency injection's sake
func ProvideIngredientDataServer(s *Service) models.IngredientDataServer {
	return s
}
