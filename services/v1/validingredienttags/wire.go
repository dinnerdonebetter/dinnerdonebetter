package validingredienttags

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideValidIngredientTagsService,
		ProvideValidIngredientTagDataManager,
		ProvideValidIngredientTagDataServer,
	)
)

// ProvideValidIngredientTagDataManager turns a database into an ValidIngredientTagDataManager.
func ProvideValidIngredientTagDataManager(db database.Database) models.ValidIngredientTagDataManager {
	return db
}

// ProvideValidIngredientTagDataServer is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientTagDataServer(s *Service) models.ValidIngredientTagDataServer {
	return s
}
