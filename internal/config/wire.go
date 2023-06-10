package config

import (
	"github.com/google/wire"
)

var (
	// ServiceConfigProviders represents this package's offering to the dependency injector.
	ServiceConfigProviders = wire.NewSet(
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
			"Search",
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
			"ValidIngredientGroups",
			"RecipeStepCompletionConditions",
			"ValidIngredientStateIngredients",
			"RecipeStepVessels",
			"Auth",
			"VendorProxy",
			"ServiceSettings",
			"ServiceSettingConfigurations",
			"UserIngredientPreferences",
			"RecipeRatings",
			"HouseholdInstrumentOwnerships",
		),
	)
)
