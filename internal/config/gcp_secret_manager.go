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

// GetAPIServerConfigFromGoogleCloudRunEnvironment fetches and InstanceConfig from GCP Secret Manager.
func GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Debug("setting up secret manager client")

	client, secretManagerCreationErr := secretmanager.NewClient(ctx)
	if secretManagerCreationErr != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", secretManagerCreationErr)
	}

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
	port, portParseErr := strconv.ParseUint(rawPort, 10, 64)
	if portParseErr != nil {
		return nil, fmt.Errorf("parsing port: %w", portParseErr)
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
	cfg.Services.Auth.Cookies.HashKey = os.Getenv("PRIXFIXE_COOKIE_HASH_KEY")
	cfg.Services.Auth.Cookies.BlockKey = os.Getenv("PRIXFIXE_COOKIE_BLOCK_KEY")
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(os.Getenv("PRIXFIXE_PASETO_LOCAL_KEY"))

	changesTopic, dataChangesNameFetchErr := fetchSecretFromSecretStore(ctx, client, "data_changes_topic_name")
	if dataChangesNameFetchErr != nil {
		return nil, fmt.Errorf("getting data changes topic name from secret store: %w", dataChangesNameFetchErr)
	}

	dataChangesTopicName := string(changesTopic)
	cfg.Events.Publishers.PubSubConfig.TopicName = dataChangesTopicName

	cfg.Email.Sendgrid.APIToken = os.Getenv("PRIXFIXE_SENDGRID_API_TOKEN")
	cfg.CustomerData.APIToken = os.Getenv("PRIXFIXE_SEGMENT_API_TOKEN")

	cfg.Services.ValidInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Recipes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeSteps.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Meals.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlans.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Households.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInvitations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Users.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Auth.DataChangesTopicName = dataChangesTopicName

	if validationErr := cfg.ValidateWithContext(ctx, true); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

var secretStorePrefix = os.Getenv("GOOGLE_CLOUD_SECRET_STORE_PREFIX")

func buildSecretPathForSecretStore(secretName string) string {
	return fmt.Sprintf(
		"%s/%s/versions/latest",
		secretStorePrefix,
		secretName,
	)
}

func fetchSecretFromSecretStore(ctx context.Context, client *secretmanager.Client, secretName string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: buildSecretPathForSecretStore(secretName),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	return result.Payload.Data, nil
}

// GetMealPlanFinalizerConfigFromGoogleCloudSecretManager fetches and InstanceConfig from GCP Secret Manager.
func GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()

	client, secretManagerCreationErr := secretmanager.NewClient(ctx)
	if secretManagerCreationErr != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", secretManagerCreationErr)
	}

	var cfg *InstanceConfig
	configBytes, configFetchErr := fetchSecretFromSecretStore(ctx, client, "api_service_config")
	if configFetchErr != nil {
		return nil, fmt.Errorf("fetching config from secret store: %w", configFetchErr)
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil || cfg == nil {
		return nil, encodeErr
	}

	rawPort := os.Getenv("PORT")
	port, portParseErr := strconv.ParseUint(rawPort, 10, 64)
	if portParseErr != nil {
		return nil, fmt.Errorf("parsing port: %w", portParseErr)
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
	cfg.Database.RunMigrations = false

	logger.Debug("fetched database values")

	cfg.Events.Publishers.PubSubConfig.TopicName = "data_changes"

	// we don't actually need these, except for validation
	cfg.CustomerData.APIToken = " "
	cfg.Email.Sendgrid.APIToken = " "

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetDataChangesWorkerConfigFromGoogleCloudSecretManager fetches and InstanceConfig from GCP Secret Manager.
func GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	var cfg *InstanceConfig
	configBytes, configFetchErr := fetchSecretFromSecretStore(ctx, client, "api_service_config")
	if configFetchErr != nil {
		return nil, fmt.Errorf("fetching config from secret store: %w", configFetchErr)
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil || cfg == nil {
		return nil, encodeErr
	}

	rawPort := os.Getenv("PORT")
	port, portParseErr := strconv.ParseUint(rawPort, 10, 64)
	if portParseErr != nil {
		return nil, fmt.Errorf("parsing port: %w", portParseErr)
	}
	cfg.Server.HTTPPort = uint16(port)

	// don't worry about it
	cfg.Database.ConnectionDetails = " "

	cfg.Email.Sendgrid.APIToken = os.Getenv("PRIXFIXE_SENDGRID_API_TOKEN")
	// we don't need the HouseholdInviteTemplateID here since we invoke that elsewhere
	cfg.CustomerData.APIToken = os.Getenv("PRIXFIXE_SEGMENT_API_TOKEN")

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}
