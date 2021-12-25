package config

import (
	"github.com/google/wire"
)

var (
	// Providers represents this package's offering to the dependency injector.
	Providers = wire.NewSet(
		wire.FieldsOf(
			new(*InstanceConfig),
			"Database",
			"Observability",
			"Meta",
			"Encoding",
			"Uploads",
			"Email",
			"CustomerData",
			"Events",
			"Server",
			"Routing",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfigurations),
			"Auth",
			"Webhooks",
			"Websockets",
			"Households",
			"HouseholdInvitations",
			"ValidInstruments",
			"ValidIngredients",
			"ValidPreparations",
			"ValidIngredientPreparations",
			"Meals",
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
