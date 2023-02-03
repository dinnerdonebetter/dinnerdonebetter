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
			"Email",
			"Analytics",
			"Events",
			"Server",
			"Routing",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfigurations),
			"ValidMeasurementUnits",
			"ValidInstruments",
			"ValidIngredients",
			"ValidPreparations",
			"MealPlanEvents",
			"MealPlanOptionVotes",
			"ValidIngredientPreparations",
			"ValidPreparationInstruments",
			"ValidInstrumentMeasurementUnits",
			"Meals",
			"Recipes",
			"RecipeSteps",
			"RecipeStepProducts",
			"RecipeStepInstruments",
			"RecipeStepIngredients",
			"MealPlans",
			"MealPlanOptions",
			"Households",
			"HouseholdInvitations",
			"Websockets",
			"Webhooks",
			"Users",
			"MealPlanTasks",
			"RecipePrepTasks",
			"MealPlanGroceryListItems",
			"ValidMeasurementConversions",
			"ValidIngredientStates",
			"ValidIngredientStateIngredients",
			"RecipeStepVessels",
			"Auth",
		),
	)
)
