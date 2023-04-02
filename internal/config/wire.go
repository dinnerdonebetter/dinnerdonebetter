package config

import (
	"github.com/prixfixeco/backend/internal/database"

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
			"RecipeStepCompletionConditions",
			"ValidIngredientStateIngredients",
			"RecipeStepVessels",
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
