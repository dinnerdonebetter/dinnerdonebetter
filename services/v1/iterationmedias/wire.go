package iterationmedias

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideIterationMediasService,
		ProvideIterationMediaDataManager,
		ProvideIterationMediaDataServer,
	)
)

// ProvideIterationMediaDataManager turns a database into an IterationMediaDataManager.
func ProvideIterationMediaDataManager(db database.DataManager) models.IterationMediaDataManager {
	return db
}

// ProvideIterationMediaDataServer is an arbitrary function for dependency injection's sake.
func ProvideIterationMediaDataServer(s *Service) models.IterationMediaDataServer {
	return s
}
