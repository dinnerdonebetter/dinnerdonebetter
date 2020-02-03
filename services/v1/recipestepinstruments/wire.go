package recipestepinstruments

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideRecipeStepInstrumentsService,
		ProvideRecipeStepInstrumentDataManager,
		ProvideRecipeStepInstrumentDataServer,
	)
)

// ProvideRecipeStepInstrumentDataManager turns a database into an RecipeStepInstrumentDataManager
func ProvideRecipeStepInstrumentDataManager(db database.Database) models.RecipeStepInstrumentDataManager {
	return db
}

// ProvideRecipeStepInstrumentDataServer is an arbitrary function for dependency injection's sake
func ProvideRecipeStepInstrumentDataServer(s *Service) models.RecipeStepInstrumentDataServer {
	return s
}
