package config

import (
	"github.com/google/wire"
)

var (
	// ServiceConfigProviders represents this package's offering to the dependency injector.
	ServiceConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*APIServiceConfig),
			"Observability",
			"Email",
			"Analytics",
			"FeatureFlags",
			"Encoding",
			"Routing",
			"Database",
			"Meta",
			"Events",
			"Queues",
			"Search",
			"Server",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfig),
			// Core
			"AuditLogEntries",
			"Auth",
			"Households",
			"HouseholdInvitations",
			"ServiceSettings",
			"ServiceSettingConfigurations",
			"Users",
			"UserNotifications",
			"Webhooks",
			"Workers",
			"DataPrivacy",
			"OAuth2Clients",
			// Data
			"ValidEnumerations",
			"MealPlanEvents",
			"MealPlanOptionVotes",
			"Meals",
			"Recipes",
			"MealPlans",
			"MealPlanOptions",
			"MealPlanTasks",
			"MealPlanGroceryListItems",
			"UserIngredientPreferences",
			"HouseholdInstrumentOwnerships",
		),
	)
)
