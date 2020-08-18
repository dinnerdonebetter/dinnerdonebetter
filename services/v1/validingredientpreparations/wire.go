package validingredientpreparations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideValidIngredientPreparationsService,
		ProvideValidIngredientPreparationDataManager,
		ProvideValidIngredientPreparationDataServer,
	)
)

// ProvideValidIngredientPreparationDataManager turns a database into an ValidIngredientPreparationDataManager.
func ProvideValidIngredientPreparationDataManager(db database.DataManager) models.ValidIngredientPreparationDataManager {
	return db
}

// ProvideValidIngredientPreparationDataServer is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientPreparationDataServer(s *Service) models.ValidIngredientPreparationDataServer {
	return s
}
