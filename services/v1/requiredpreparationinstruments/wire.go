package requiredpreparationinstruments

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideRequiredPreparationInstrumentsService,
		ProvideRequiredPreparationInstrumentDataManager,
		ProvideRequiredPreparationInstrumentDataServer,
	)
)

// ProvideRequiredPreparationInstrumentDataManager turns a database into an RequiredPreparationInstrumentDataManager
func ProvideRequiredPreparationInstrumentDataManager(db database.Database) models.RequiredPreparationInstrumentDataManager {
	return db
}

// ProvideRequiredPreparationInstrumentDataServer is an arbitrary function for dependency injection's sake
func ProvideRequiredPreparationInstrumentDataServer(s *Service) models.RequiredPreparationInstrumentDataServer {
	return s
}
