package instruments

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideInstrumentsService,
		ProvideInstrumentDataManager,
		ProvideInstrumentDataServer,
	)
)

// ProvideInstrumentDataManager turns a database into an InstrumentDataManager
func ProvideInstrumentDataManager(db database.Database) models.InstrumentDataManager {
	return db
}

// ProvideInstrumentDataServer is an arbitrary function for dependency injection's sake
func ProvideInstrumentDataServer(s *Service) models.InstrumentDataServer {
	return s
}
