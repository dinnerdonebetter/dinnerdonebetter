package dbcleaner

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/samber/do/v2"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
)

// RegisterConfigs registers all config sub-fields with the injector.
func RegisterConfigs(i do.Injector) {
	do.Provide[*observability.Config](i, func(i do.Injector) (*observability.Config, error) {
		cfg := do.MustInvoke[*config.DBCleanerConfig](i)
		return &cfg.Observability, nil
	})
	do.Provide[*databasecfg.Config](i, func(i do.Injector) (*databasecfg.Config, error) {
		cfg := do.MustInvoke[*config.DBCleanerConfig](i)
		return &cfg.Database, nil
	})
}
