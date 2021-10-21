package config

import (
	"github.com/google/wire"
)

var (
	// Providers represents this package's offering to the dependency injector.
	Providers = wire.NewSet(
		ProvideDatabaseClient,
		wire.FieldsOf(
			new(*InstanceConfig),
			"Database",
			"Observability",
			"Meta",
			"Encoding",
			"Uploads",
			"Search",
			"Events",
			"Server",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfigurations),
			"Auth",
			"Frontend",
			"Webhooks",
			"Websockets",
			"Accounts",
			"ValidInstruments",
			"ValidIngredients",
			"ValidPreparations",
			"ValidIngredientPreparations",
			"Recipes",
			"RecipeSteps",
			"RecipeStepInstruments",
			"RecipeStepIngredients",
			"RecipeStepProducts",
			"MealPlans",
			"MealPlanOptions",
			"MealPlanOptionVotes",
		),
	)
)
