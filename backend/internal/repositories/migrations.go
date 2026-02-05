package repositories

import (
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	postgresmigrations "github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
)

// ProvideMigrator returns a Migrator appropriate for the configured database provider.
func ProvideMigrator(
	cfg *databasecfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) database.Migrator {
	switch cfg.Provider {
	case databasecfg.ProviderPostgres:
		return postgresmigrations.NewMigrator(logger, tracerProvider)
	default:
		return nil
	}
}
