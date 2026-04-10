package mealplangrocerylistinitializer

import (
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"

	analyticscfg "github.com/verygoodsoftwarenotvirus/platform/v5/analytics/config"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v5/database/config"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v5/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability"

	"github.com/samber/do/v2"
)

// RegisterConfigs registers all config sub-fields with the injector.
func RegisterConfigs(i do.Injector) {
	do.Provide[*msgconfig.QueuesConfig](i, func(i do.Injector) (*msgconfig.QueuesConfig, error) {
		cfg := do.MustInvoke[*config.MealPlanGroceryListInitializerConfig](i)
		return &cfg.Queues, nil
	})
	do.Provide[*analyticscfg.Config](i, func(i do.Injector) (*analyticscfg.Config, error) {
		cfg := do.MustInvoke[*config.MealPlanGroceryListInitializerConfig](i)
		return &cfg.Analytics, nil
	})
	do.Provide[*msgconfig.Config](i, func(i do.Injector) (*msgconfig.Config, error) {
		cfg := do.MustInvoke[*config.MealPlanGroceryListInitializerConfig](i)
		return &cfg.Events, nil
	})
	do.Provide[*observability.Config](i, func(i do.Injector) (*observability.Config, error) {
		cfg := do.MustInvoke[*config.MealPlanGroceryListInitializerConfig](i)
		return &cfg.Observability, nil
	})
	do.Provide[*databasecfg.Config](i, func(i do.Injector) (*databasecfg.Config, error) {
		cfg := do.MustInvoke[*config.MealPlanGroceryListInitializerConfig](i)
		return &cfg.Database, nil
	})
}
