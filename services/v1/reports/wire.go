package reports

import (
	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideReportsService,
		ProvideReportDataManager,
		ProvideReportDataServer,
	)
)

// ProvideReportDataManager turns a database into an ReportDataManager
func ProvideReportDataManager(db database.Database) models.ReportDataManager {
	return db
}

// ProvideReportDataServer is an arbitrary function for dependency injection's sake
func ProvideReportDataServer(s *Service) models.ReportDataServer {
	return s
}
