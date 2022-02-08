package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// GetConfigFromCloudSecretManager fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Debug("setting up secret manager client")

	var cfg *InstanceConfig
	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")

	configBytes, configReadErr := os.ReadFile(configFilepath)
	if configReadErr != nil {
		return nil, configReadErr
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil || cfg == nil {
		return nil, encodeErr
	}

	rawPort := os.Getenv("PORT")
	port, err := strconv.ParseUint(rawPort, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing port: %w", err)
	}
	cfg.Server.HTTPPort = uint16(port)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	// fetch supplementary data from env vars
	dbURI := fmt.Sprintf(
		"user=%s password=%s database=%s host=%s/%s",
		os.Getenv("PRIXFIXE_DATABASE_USER"),
		os.Getenv("PRIXFIXE_DATABASE_PASSWORD"),
		os.Getenv("PRIXFIXE_DATABASE_NAME"),
		socketDir,
		os.Getenv("PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME"),
	)

	cfg.Database.ConnectionDetails = database.ConnectionDetails(dbURI)

	logger.WithValues(map[string]interface{}{
		"DB_SOCKET_DIR":                              os.Getenv("DB_SOCKET_DIR"),
		"PRIXFIXE_DATABASE_USER":                     os.Getenv("PRIXFIXE_DATABASE_USER"),
		"PRIXFIXE_DATABASE_PASSWORD":                 os.Getenv("PRIXFIXE_DATABASE_PASSWORD"),
		"PRIXFIXE_DATABASE_NAME":                     os.Getenv("PRIXFIXE_DATABASE_NAME"),
		"PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME": os.Getenv("PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME"),
	}).Debug("fetched database values")

	cfg.Services.Auth.Cookies.HashKey = os.Getenv("PRIXFIXE_COOKIE_HASH_KEY")
	cfg.Services.Auth.Cookies.BlockKey = os.Getenv("PRIXFIXE_COOKIE_BLOCK_KEY")
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(os.Getenv("PRIXFIXE_PASETO_LOCAL_KEY"))

	dataChangesTopicName := os.Getenv("PRIXFIXE_DATA_CHANGES_TOPIC")

	cfg.Email.APIToken = os.Getenv("PRIXFIXE_SENDGRID_API_TOKEN")
	cfg.CustomerData.APIToken = os.Getenv("PRIXFIXE_SEGMENT_API_TOKEN")

	cfg.Services.ValidInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Recipes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeSteps.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepProducts.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Meals.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlans.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Households.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInvitations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}
