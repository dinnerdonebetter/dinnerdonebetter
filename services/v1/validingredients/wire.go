package validingredients

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideValidIngredientsService,
		ProvideValidIngredientDataManager,
		ProvideValidIngredientDataServer,
	)
)

// ProvideValidIngredientDataManager turns a database into an ValidIngredientDataManager.
func ProvideValidIngredientDataManager(db database.DataManager) models.ValidIngredientDataManager {
	return db
}

// ProvideValidIngredientDataServer is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientDataServer(s *Service) models.ValidIngredientDataServer {
	return s
}
