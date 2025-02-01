package mealplangrocerylistinitializer

import (
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.MealPlanGroceryListInitializerConfig),
			"Queues",
			"Analytics",
			"Events",
			"Observability",
			"Database",
		),
	)
)
