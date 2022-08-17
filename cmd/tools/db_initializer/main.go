package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing"

	_ "github.com/lib/pq"
	flag "github.com/spf13/pflag"
)

const (
	defaultDBURL  = "postgres://dbuser:hunter2@localhost:5432/prixfixe?sslmode=disable"
	clearAllQuery = `
DELETE FROM "users" WHERE id IS NOT NULL;
DELETE FROM "households" WHERE id IS NOT NULL;
DELETE FROM "household_invitations" WHERE id IS NOT NULL;
DELETE FROM "household_user_memberships" WHERE id IS NOT NULL;
DELETE FROM "valid_ingredients" WHERE id IS NOT NULL;
DELETE FROM "valid_instruments" WHERE id IS NOT NULL;
DELETE FROM "valid_preparations" WHERE id IS NOT NULL;
DELETE FROM "valid_measurement_units" WHERE id IS NOT NULL;
DELETE FROM "valid_ingredient_preparations" WHERE id IS NOT NULL;
DELETE FROM "valid_preparation_instruments" WHERE id IS NOT NULL;
DELETE FROM "valid_ingredient_measurement_units" WHERE id IS NOT NULL;
DELETE FROM "recipes" WHERE id IS NOT NULL;
DELETE FROM "recipe_steps" WHERE id IS NOT NULL;
DELETE FROM "recipe_step_products" WHERE id IS NOT NULL;
DELETE FROM "recipe_step_instruments" WHERE id IS NOT NULL;
DELETE FROM "recipe_step_ingredients" WHERE id IS NOT NULL;
DELETE FROM "meals" WHERE id IS NOT NULL;
DELETE FROM "meal_recipes" WHERE id IS NOT NULL;
DELETE FROM "meal_plans" WHERE id IS NOT NULL;
DELETE FROM "meal_plan_options" WHERE id IS NOT NULL;
DELETE FROM "meal_plan_option_votes" WHERE id IS NOT NULL;
DELETE FROM "sessions" WHERE data IS NOT NULL;
`
)

var (
	dbString string
	debug    bool
)

func padID(s string) string {
	var x []byte

	if len(s) > 27 {
		log.Panicf("invalid identifier: %q", s)
	}

	stopIndex := 0
	for i, b := range s {
		x = append(x, byte(b))
		stopIndex = i + 1
	}

	for i := stopIndex; i < 27; i++ {
		x = append(x, []byte("_")[0])
	}

	return string(x)
}

func init() {
	flag.StringVarP(&dbString, "dburl", "u", defaultDBURL, "where the database is hosted")
	flag.BoolVarP(&debug, "debug", "d", false, "whether debug mode is enabled")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	logger, err := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbConfig := &dbconfig.Config{
		ConnectionDetails: database.ConnectionDetails(dbString),
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracing.NewNoopTracerProvider())
	if err != nil {
		log.Fatal(fmt.Errorf("initializing database client: %w", err))
	}

	_, err = dataManager.DB().ExecContext(ctx, clearAllQuery)
	if err != nil {
		log.Fatal(fmt.Errorf("error clearing database: %w", err))
	}

	if err = scaffoldUsers(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating users: %w", err))
	}

	if err = scaffoldValidIngredients(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid ingredients: %w", err))
	}

	if err = scaffoldValidPreparations(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid preparations: %w", err))
	}

	if err = scaffoldValidInstruments(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid instruments: %w", err))
	}

	if err = scaffoldValidMeasurementUnits(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid measurement units: %w", err))
	}

	if err = scaffoldValidPreparationInstruments(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid preparation instruments: %w", err))
	}

	if err = scaffoldValidIngredientPreparations(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid ingredient preparations: %w", err))
	}

	if err = scaffoldValidIngredientMeasurementUnits(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating valid ingredient preparations: %w", err))
	}

	if err = scaffoldRecipes(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating recipes: %w", err))
	}

	if err = scaffoldMeals(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating meals: %w", err))
	}

	if err = scaffoldMealPlans(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating meals: %w", err))
	}

	if err = scaffoldMealPlanVotes(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating meals: %w", err))
	}
}

func sp(s string) *string {
	return &s
}

func uint8p(n uint8) *uint8 {
	return &n
}
