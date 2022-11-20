package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	flag "github.com/spf13/pflag"

	"github.com/prixfixeco/backend/internal/database"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	logcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	"github.com/prixfixeco/backend/internal/observability/tracing"
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
DELETE FROM "meal_components" WHERE id IS NOT NULL;
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

	dumpBytes, err := os.ReadFile("cmd/tools/db_initializer/dump.sql")
	if err != nil {
		log.Fatal(fmt.Errorf("error reading dump file: %w", err))
	}

	if _, err = dataManager.DB().ExecContext(ctx, clearAllQuery); err != nil {
		log.Fatal(fmt.Errorf("error clearing database: %w", err))
	}

	if err = scaffoldUsers(ctx, dataManager); err != nil {
		log.Fatal(fmt.Errorf("error creating users: %w", err))
	}

	momJones, err := dataManager.GetUserByUsername(ctx, userCollection.MomJones.Username)
	if err != nil {
		log.Fatal(fmt.Errorf("error fetching momJones user: %w", err))
	}

	replacedDump := strings.ReplaceAll(string(dumpBytes), "2751SjGkKN5AzdVbcNP0eblooTC", momJones.ID)

	if _, err = dataManager.DB().ExecContext(ctx, replacedDump); err != nil {
		log.Fatal(fmt.Errorf("initializing running query: %w", err))
	}
}
