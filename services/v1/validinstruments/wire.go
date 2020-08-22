package validinstruments

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideValidInstrumentsService,
		ProvideValidInstrumentDataManager,
		ProvideValidInstrumentDataServer,
		ProvideValidInstrumentsServiceSearchIndex,
	)
)

// ProvideValidInstrumentDataManager turns a database into an ValidInstrumentDataManager.
func ProvideValidInstrumentDataManager(db database.DataManager) models.ValidInstrumentDataManager {
	return db
}

// ProvideValidInstrumentDataServer is an arbitrary function for dependency injection's sake.
func ProvideValidInstrumentDataServer(s *Service) models.ValidInstrumentDataServer {
	return s
}
