package repositories

import (
	postgresmigrations "github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"

	"github.com/verygoodsoftwarenotvirus/platform/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
)

// ProvideMigrator returns a Migrator appropriate for the configured database provider.
func ProvideMigrator(
	cfg *databasecfg.Config,
	logger logging.Logger,
) database.Migrator {
	switch cfg.Provider {
	case databasecfg.ProviderPostgres:
		return postgresmigrations.NewMigrator(logger)
	default:
		return nil
	}
}
