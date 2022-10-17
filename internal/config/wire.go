package config

import (
	"github.com/google/wire"

	"github.com/prixfixeco/api_server/internal/database"
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
			"Auth",
		),
	)
)

// ProvideCloserFunc provides a closer function.
func ProvideCloserFunc(dbm database.DataManager) func() error {
	return func() error {
		db := dbm.DB()
		if err := db.Close(); err != nil {
			return err
		}

		return nil
	}
}
