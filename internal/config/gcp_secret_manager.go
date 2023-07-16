package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	analyticscfg "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
)

const (
	dataChangesTopicAccessName           = "data_changes_topic_name"
	googleCloudCloudSQLSocket            = "/cloudsql"
	gcpConfigFilePathEnvVarKey           = "CONFIGURATION_FILEPATH"
	gcpPortEnvVarKey                     = "PORT"
	gcpDatabaseSocketDirEnvVarKey        = "DB_SOCKET_DIR"
	gcpDatabaseUserEnvVarKey             = "DINNER_DONE_BETTER_DATABASE_USER"
	gcpDatabaseNameEnvVarKey             = "DINNER_DONE_BETTER_DATABASE_NAME"
	gcpDatabaseInstanceConnNameEnvVarKey = "DINNER_DONE_BETTER_DATABASE_INSTANCE_CONNECTION_NAME"
	gcpCookieHashKeyEnvVarKey            = "DINNER_DONE_BETTER_COOKIE_HASH_KEY"
	gcpCookieBlockKeyEnvVarKey           = "DINNER_DONE_BETTER_COOKIE_BLOCK_KEY"
	gcpPASETOLocalKeyEnvVarKey           = "DINNER_DONE_BETTER_PASETO_LOCAL_KEY"
	gcpAlgoliaAPIKeyEnvVarKey            = "DINNER_DONE_BETTER_ALGOLIA_API_KEY"
	gcpAlgoliaAppIDEnvVarKey             = "DINNER_DONE_BETTER_ALGOLIA_APPLICATION_ID"
	/* #nosec G101 */
	gcpOauth2TokenEncryptionKeyEnvVarKey = "DINNER_DONE_BETTER_OAUTH2_TOKEN_ENCRYPTION_KEY"
	/* #nosec G101 */
	gcpDatabaseUserPasswordEnvVarKey = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
	/* #nosec G101 */
	gcpSendgridTokenEnvVarKey = "DINNER_DONE_BETTER_SENDGRID_API_TOKEN"
	/* #nosec G101 */
	gcpSegmentTokenEnvVarKey = "DINNER_DONE_BETTER_SEGMENT_API_TOKEN"
)

// SecretVersionAccessor is an interface abstraction of the GCP Secret Manager API call we use during config hydration.
// This interface exists for testing purposes. Yes, you're not supposed to write arbitrary interfaces for testing.
// Yes I'm still doing it.
type SecretVersionAccessor interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

// GetAPIServerConfigFromGoogleCloudRunEnvironment fetches an InstanceConfig from GCP Secret Manager.
func GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx context.Context, client SecretVersionAccessor) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger(logging.DebugLevel)
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
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)

	cfg.Services.Auth.Cookies.HashKey = os.Getenv(gcpCookieHashKeyEnvVarKey)
	cfg.Services.Auth.Cookies.BlockKey = os.Getenv(gcpCookieBlockKeyEnvVarKey)
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(os.Getenv(gcpPASETOLocalKeyEnvVarKey))

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	// TODO: get this from the env var DATA_CHANGES_TOPIC_NAME, dump GOOGLE_CLOUD_SECRET_STORE_PREFIX
	changesTopic, err := fetchSecretFromSecretStore(ctx, client, dataChangesTopicAccessName)
	if err != nil {
		return nil, fmt.Errorf("getting data changes topic name from secret store: %w", err)
	}
	dataChangesTopicName := string(changesTopic)

	cfg.Email.Sendgrid.APIToken = os.Getenv(gcpSendgridTokenEnvVarKey)
	cfg.Analytics.Segment = &segment.Config{APIToken: os.Getenv(gcpSegmentTokenEnvVarKey)}

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
	cfg.Services.ValidIngredientGroups.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Auth.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipePrepTasks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanGroceryListItems.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidMeasurementConversions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStates.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepCompletionConditions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStateIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.VendorProxy.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettingConfigurations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.UserIngredientPreferences.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeRatings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInstrumentOwnerships.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparationVessels.DataChangesTopicName = dataChangesTopicName

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

// getWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func getWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	client, secretManagerCreationErr := secretmanager.NewClient(ctx)
	if secretManagerCreationErr != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", secretManagerCreationErr)
	}

	var cfg *InstanceConfig
	configBytes, configFetchErr := fetchSecretFromSecretStore(ctx, client, "api_service_config")
	if configFetchErr != nil {
		return nil, fmt.Errorf("fetching config from secret store: %w", configFetchErr)
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil {
		return nil, encodeErr
	}

	if cfg == nil {
		return nil, errors.New("config is nil")
	}

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

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	cfg.Database.ConnectionDetails = database.ConnectionDetails(dbURI)
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = " "
	cfg.Email.Sendgrid.APIToken = os.Getenv(gcpSendgridTokenEnvVarKey)
	cfg.Analytics = analyticscfg.Config{
		Segment:  &segment.Config{APIToken: os.Getenv(gcpSegmentTokenEnvVarKey)},
		Provider: analyticscfg.ProviderSegment,
	}

	return cfg, nil
}

// GetDataChangesWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetMealPlanFinalizerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetMealPlanGroceryListInitializerWorkerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetMealPlanGroceryListInitializerWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetOutboundEmailerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetOutboundEmailerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{
		Provider: emailcfg.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{
			APIToken: os.Getenv(gcpSendgridTokenEnvVarKey),
		},
	}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetSearchDataIndexerConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetSearchDataIndexerConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

// GetEmailProberConfigFromGoogleCloudSecretManager fetches an InstanceConfig from GCP Secret Manager.
func GetEmailProberConfigFromGoogleCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database = dbconfig.Config{ConnectionDetails: " ", OAuth2TokenEncryptionKey: " "}
	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{
		Provider: emailcfg.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{
			APIToken: os.Getenv(gcpSendgridTokenEnvVarKey),
		},
	}

	if validationErr := cfg.ValidateWithContext(ctx, false); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}
