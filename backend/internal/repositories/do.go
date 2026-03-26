package repositories

import (
	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v3/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v3/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v3/observability/logging"
)

// RegisterMigrator registers the database migrator with the injector.
func RegisterMigrator(i do.Injector) {
	do.Provide[database.Migrator](i, func(i do.Injector) (database.Migrator, error) {
		return ProvideMigrator(
			do.MustInvoke[*databasecfg.Config](i),
			do.MustInvoke[logging.Logger](i),
		), nil
	})
}
