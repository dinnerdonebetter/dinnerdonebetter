package preparations

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvidePreparationsService,
		ProvidePreparationDataManager,
		ProvidePreparationDataServer,
	)
)

// ProvidePreparationDataManager turns a database into an PreparationDataManager
func ProvidePreparationDataManager(db database.Database) models.PreparationDataManager {
	return db
}

// ProvidePreparationDataServer is an arbitrary function for dependency injection's sake
func ProvidePreparationDataServer(s *Service) models.PreparationDataServer {
	return s
}
