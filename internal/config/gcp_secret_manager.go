package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// GetConfigFromGoogleCloudRunEnvironment fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromGoogleCloudRunEnvironment(ctx context.Context) (*InstanceConfig, error) {
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

	if err = cfg.ValidateWithContext(ctx, true); err != nil {
		return nil, err
	}

	return cfg, nil
}

func fetchSecretFromSecretStore(ctx context.Context, client *secretmanager.Client, secretName string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: secretName}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, nil
}

// GetMealPlanFinalizerConfigFromGoogleCloudSecretManager fetches and InstanceConfig from GCP Secret Manager.
func GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Info("setting up secret manager client")

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	secretPrefix := os.Getenv("GOOGLE_CLOUD_SECRET_STORE_PREFIX")
	configSecretPath := fmt.Sprintf("%s/%s/versions/latest", secretPrefix, "api_service_config")

	var cfg *InstanceConfig
	configBytes, err := fetchSecretFromSecretStore(ctx, client, configSecretPath)
	if err != nil {
		return nil, fmt.Errorf("fetching config from secret store: %w", err)
	}

	logger.WithValue("raw_config", string(configBytes)).Info("config retrieved")

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

	logger.Debug("fetched database values")

	//dataChangesTopicName, err := fetchSecretFromSecretStore(ctx, client, "data_changes_topic_name")
	//if err != nil {
	//	return nil, fmt.Errorf("error getting data changes topic name from secret store: %w", err)
	//}

	// we don't actually need these, except for validation
	cfg.CustomerData.APIToken = " "
	cfg.Email.APIToken = " "

	if err = cfg.ValidateWithContext(ctx, false); err != nil {
		return nil, err
	}

	return cfg, nil
}
