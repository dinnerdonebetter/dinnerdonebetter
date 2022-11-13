package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/googleapis/gax-go/v2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
)

const (
	gcpConfigFilePathEnvVarKey           = "CONFIGURATION_FILEPATH"
	gcpPortEnvVarKey                     = "PORT"
	gcpDatabaseSocketDirEnvVarKey        = "DB_SOCKET_DIR"
	gcpDatabaseUserEnvVarKey             = "PRIXFIXE_DATABASE_USER"
	gcpDatabaseNameEnvVarKey             = "PRIXFIXE_DATABASE_NAME"
	gcpDatabaseInstanceConnNameEnvVarKey = "PRIXFIXE_DATABASE_INSTANCE_CONNECTION_NAME"
	gcpCookieHashKeyEnvVarKey            = "PRIXFIXE_COOKIE_HASH_KEY"
	gcpCookieBlockKeyEnvVarKey           = "PRIXFIXE_COOKIE_BLOCK_KEY"
	gcpPASETOLocalKeyEnvVarKey           = "PRIXFIXE_PASETO_LOCAL_KEY"
	/* #nosec G101 */
	gcpDatabaseUserPasswordEnvVarKey = "PRIXFIXE_DATABASE_PASSWORD"
	/* #nosec G101 */
	gcpSendgridTokenEnvVarKey = "PRIXFIXE_SENDGRID_API_TOKEN"
	/* #nosec G101 */
	gcpSegmentTokenEnvVarKey = "PRIXFIXE_SEGMENT_API_TOKEN"

	dataChangesTopicAccessName = "data_changes_topic_name"
	googleCloudCloudSQLSocket  = "/cloudsql"
)

// SecretVersionAccessor is an interface abstraction of the GCP Secret Manager API call we use during config hydration.
// This interface exists for testing purposes. Yes, you're not supposed to write arbitrary interfaces for testing.
// Yes I'm still doing it.
type SecretVersionAccessor interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

// GetAPIServerConfigFromGoogleCloudRunEnvironment fetches an InstanceConfig from GCP Secret Manager.
func GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx context.Context, client SecretVersionAccessor) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Debug("setting up secret manager client")

	var cfg *InstanceConfig
	configFilepath := os.Getenv(gcpConfigFilePathEnvVarKey)

	configBytes, configReadErr := os.ReadFile(configFilepath)
	if configReadErr != nil {
		return nil, configReadErr
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil || cfg == nil {
		return nil, encodeErr
	}

	rawPort := os.Getenv(gcpPortEnvVarKey)
	port, portParseErr := strconv.ParseUint(rawPort, 10, 64)
	if portParseErr != nil {
		return nil, fmt.Errorf("parsing port: %w", portParseErr)
	}
	cfg.Server.HTTPPort = uint16(port)

	socketDir, isSet := os.LookupEnv(gcpDatabaseSocketDirEnvVarKey)
	if !isSet {
		socketDir = googleCloudCloudSQLSocket
	}

	// fetch supplementary data from env vars
	dbURI := fmt.Sprintf(
		"user=%s password=%s database=%s host=%s/%s",
		os.Getenv(gcpDatabaseUserEnvVarKey),
		os.Getenv(gcpDatabaseUserPasswordEnvVarKey),
		os.Getenv(gcpDatabaseNameEnvVarKey),
		socketDir,
		os.Getenv(gcpDatabaseInstanceConnNameEnvVarKey),
	)

	cfg.Database.ConnectionDetails = database.ConnectionDetails(dbURI)
	cfg.Services.Auth.Cookies.HashKey = os.Getenv(gcpCookieHashKeyEnvVarKey)
	cfg.Services.Auth.Cookies.BlockKey = os.Getenv(gcpCookieBlockKeyEnvVarKey)
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(os.Getenv(gcpPASETOLocalKeyEnvVarKey))

	changesTopic, dataChangesNameFetchErr := fetchSecretFromSecretStore(ctx, client, dataChangesTopicAccessName)
	if dataChangesNameFetchErr != nil {
		return nil, fmt.Errorf("getting data changes topic name from secret store: %w", dataChangesNameFetchErr)
	}

	dataChangesTopicName := string(changesTopic)
	cfg.Events.Publishers.PubSubConfig.TopicName = dataChangesTopicName

	cfg.Email.Sendgrid.APIToken = os.Getenv(gcpSendgridTokenEnvVarKey)
	cfg.CustomerData.APIToken = os.Getenv(gcpSegmentTokenEnvVarKey)

	cfg.Services.ValidMeasurementUnits.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparationInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidInstrumentMeasurementUnits.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Recipes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeSteps.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepProducts.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Meals.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlans.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanEvents.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanTasks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Households.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInvitations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Users.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Auth.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipePrepTasks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanGroceryListItems.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidMeasurementConversions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStates.DataChangesTopicName = dataChangesTopicName

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

func fetchSecretFromSecretStore(ctx context.Context, client SecretVersionAccessor, secretName string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: buildSecretPathForSecretStore(secretName),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	return result.Payload.Data, nil
}

// GetMealPlanFinalizerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
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
		socketDir = googleCloudCloudSQLSocket
	}

	// fetch supplementary data from env vars
	dbURI := fmt.Sprintf(
		"user=%s password=%s database=%s host=%s/%s",
		os.Getenv(gcpDatabaseUserEnvVarKey),
		os.Getenv(gcpDatabaseUserPasswordEnvVarKey),
		os.Getenv(gcpDatabaseNameEnvVarKey),
		socketDir,
		os.Getenv(gcpDatabaseInstanceConnNameEnvVarKey),
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

// GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
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
		socketDir = googleCloudCloudSQLSocket
	}

	// fetch supplementary data from env vars
	dbURI := fmt.Sprintf(
		"user=%s password=%s database=%s host=%s/%s",
		os.Getenv(gcpDatabaseUserEnvVarKey),
		os.Getenv(gcpDatabaseUserPasswordEnvVarKey),
		os.Getenv(gcpDatabaseNameEnvVarKey),
		socketDir,
		os.Getenv(gcpDatabaseInstanceConnNameEnvVarKey),
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

// GetDataChangesWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
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
