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
			"Auth",
			"Users",
			"Webhooks",
			"Websockets",
			"Households",
			"HouseholdInvitations",
			"ValidInstruments",
			"ValidIngredients",
			"ValidPreparations",
			"ValidMeasurementUnits",
			"ValidIngredientPreparations",
			"Meals",
			"Recipes",
			"RecipeSteps",
			"RecipeStepProducts",
			"RecipeStepInstruments",
			"RecipeStepIngredients",
			"MealPlans",
			"MealPlanOptions",
			"MealPlanOptionVotes",
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
