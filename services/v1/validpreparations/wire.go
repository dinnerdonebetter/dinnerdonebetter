package validpreparations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideValidPreparationsService,
		ProvideValidPreparationDataManager,
		ProvideValidPreparationDataServer,
	)
)

// ProvideValidPreparationDataManager turns a database into an ValidPreparationDataManager.
func ProvideValidPreparationDataManager(db database.DataManager) models.ValidPreparationDataManager {
	return db
}

// ProvideValidPreparationDataServer is an arbitrary function for dependency injection's sake.
func ProvideValidPreparationDataServer(s *Service) models.ValidPreparationDataServer {
	return s
}
