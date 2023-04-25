package config

import (
	"github.com/google/wire"
)

var (
	// Providers represents this package's offering to the dependency injector.
	Providers = wire.NewSet(
		wire.FieldsOf(
			new(*InstanceConfig),
			"Observability",
			"Email",
			"Analytics",
			"FeatureFlags",
			"Encoding",
			"Routing",
			"Database",
			"Meta",
			"Events",
			"Server",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfig),
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
			"RecipeStepCompletionConditions",
			"ValidIngredientStateIngredients",
			"RecipeStepVessels",
			"Auth",
			"VendorProxy",
			"ServiceSettings",
			"ServiceSettingConfigurations",
		),
	)
)
