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
	"github.com/dinnerdonebetter/backend/internal/analytics/posthog"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/search/text/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
)

const (
	googleCloudCloudSQLSocket             = "/cloudsql"
	gcpPortEnvVarKey                      = "PORT"
	gcpDatabaseSocketDirEnvVarKey         = "DB_SOCKET_DIR"
	gcpDataChangesTopicNameEnvVarKey      = "DINNER_DONE_BETTER_DATA_CHANGES_TOPIC_NAME"
	gcpOutboundEmailsTopicNameEnvVarKey   = "DINNER_DONE_BETTER_OUTBOUND_EMAILS_TOPIC_NAME"
	gcpSearchIndexingTopicNameEnvVarKey   = "DINNER_DONE_BETTER_SEARCH_INDEXING_TOPIC_NAME"
	gcpWebhookExecutionTopicNameEnvVarKey = "DINNER_DONE_BETTER_WEBHOOK_EXECUTION_REQUESTS_TOPIC_NAME"
	gcpUserAggregatorTopicName            = "DINNER_DONE_BETTER_USER_AGGREGATOR_TOPIC_NAME"
	gcpDatabaseUserEnvVarKey              = "DINNER_DONE_BETTER_DATABASE_USER"
	gcpDatabaseNameEnvVarKey              = "DINNER_DONE_BETTER_DATABASE_NAME"
	gcpDatabaseInstanceConnNameEnvVarKey  = "DINNER_DONE_BETTER_DATABASE_INSTANCE_CONNECTION_NAME"
	gcpAlgoliaAPIKeyEnvVarKey             = "DINNER_DONE_BETTER_ALGOLIA_API_KEY"
	gcpAlgoliaAppIDEnvVarKey              = "DINNER_DONE_BETTER_ALGOLIA_APPLICATION_ID"
	gcpGoogleSSOClientIDEnvVarKey         = "DINNER_DONE_BETTER_GOOGLE_SSO_CLIENT_ID"
	gcpGoogleSSOClientSecretEnvVarKey     = "DINNER_DONE_BETTER_GOOGLE_SSO_CLIENT_SECRET"
	/* #nosec G101 */
	gcpJWTSigningKeyEnvVarKey = "DINNER_DONE_BETTER_JWT_SIGNING_KEY"
	/* #nosec G101 */
	gcpOauth2TokenEncryptionKeyEnvVarKey = "DINNER_DONE_BETTER_OAUTH2_TOKEN_ENCRYPTION_KEY"
	/* #nosec G101 */
	gcpDatabaseUserPasswordEnvVarKey = "DINNER_DONE_BETTER_DATABASE_PASSWORD"
	/* #nosec G101 */
	gcpSendgridTokenEnvVarKey = "DINNER_DONE_BETTER_SENDGRID_API_TOKEN"
	/* #nosec G101 */
	gcpSegmentTokenEnvVarKey = "DINNER_DONE_BETTER_SEGMENT_API_TOKEN"
	/* #nosec G101 */
	gcpPostHogKeyEnvVarKey = "DINNER_DONE_BETTER_POSTHOG_API_KEY"
)

// SecretVersionAccessor is an interface abstraction of the GCP Secret Manager API call we use during config hydration.
// This interface exists for testing purposes. Yes, you're not supposed to write arbitrary interfaces for testing.
// Yes I'm still doing it.
type SecretVersionAccessor interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

func buildDatabaseURIFromGCPEnvVars() string {
	socketDir, isSet := os.LookupEnv(gcpDatabaseSocketDirEnvVarKey)
	if !isSet {
		socketDir = googleCloudCloudSQLSocket
	}

	return fmt.Sprintf(
		"user=%s password=%s database=%s host=%s/%s",
		os.Getenv(gcpDatabaseUserEnvVarKey),
		os.Getenv(gcpDatabaseUserPasswordEnvVarKey),
		os.Getenv(gcpDatabaseNameEnvVarKey),
		socketDir,
		os.Getenv(gcpDatabaseInstanceConnNameEnvVarKey),
	)
}

// GetAPIServiceConfigFromGoogleCloudRunEnvironment fetches an APIServiceConfig from GCP Secret Manager.
func GetAPIServiceConfigFromGoogleCloudRunEnvironment(ctx context.Context) (*APIServiceConfig, error) {
	configBytes, err := os.ReadFile(os.Getenv(FilePathEnvVarKey))
	if err != nil {
		return nil, err
	}

	var cfg *APIServiceConfig
	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
		return nil, err
	}

	rawPort := os.Getenv(gcpPortEnvVarKey)
	var port uint64
	port, err = strconv.ParseUint(rawPort, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing port: %w", err)
	}
	cfg.Server.HTTPPort = uint16(port)

	// fetch supplementary data from env vars
	dbURI := buildDatabaseURIFromGCPEnvVars()

	cfg.Database.ConnectionDetails = dbURI
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Services.Auth.SSO.Google.ClientID = os.Getenv(gcpGoogleSSOClientIDEnvVarKey)
	cfg.Services.Auth.SSO.Google.ClientSecret = os.Getenv(gcpGoogleSSOClientSecretEnvVarKey)
	cfg.Services.Auth.JWTSigningKey = os.Getenv(gcpJWTSigningKeyEnvVarKey)

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	dataChangesTopicName := os.Getenv(gcpDataChangesTopicNameEnvVarKey)

	cfg.Queues = QueuesConfig{
		DataChangesTopicName:              os.Getenv(gcpDataChangesTopicNameEnvVarKey),
		OutboundEmailsTopicName:           os.Getenv(gcpOutboundEmailsTopicNameEnvVarKey),
		SearchIndexRequestsTopicName:      os.Getenv(gcpSearchIndexingTopicNameEnvVarKey),
		UserDataAggregationTopicName:      os.Getenv(gcpUserAggregatorTopicName),
		WebhookExecutionRequestsTopicName: os.Getenv(gcpWebhookExecutionTopicNameEnvVarKey),
	}

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
	cfg.Services.ValidMeasurementUnitConversions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStates.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepCompletionConditions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStateIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettingConfigurations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.UserIngredientPreferences.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeRatings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInstrumentOwnerships.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparationVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Workers.DataChangesTopicName = dataChangesTopicName
	cfg.Services.UserNotifications.DataChangesTopicName = dataChangesTopicName
	cfg.Services.DataPrivacy.DataChangesTopicName = dataChangesTopicName
	cfg.Services.DataPrivacy.UserDataAggregationTopicName = os.Getenv(gcpUserAggregatorTopicName)

	cfg.validateServices = true

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

var secretStorePrefix = os.Getenv("GOOGLE_CLOUD_SECRET_STORE_PREFIX")

func buildSecretPathForGCPSecretStore(secretName string) string {
	return fmt.Sprintf(
		"%s/%s/versions/latest",
		secretStorePrefix,
		secretName,
	)
}

func fetchSecretFromSecretStore(ctx context.Context, client SecretVersionAccessor, secretName string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: buildSecretPathForGCPSecretStore(secretName),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	return result.Payload.Data, nil
}

// BEGIN JOBS

// GetDBCleanerConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetDBCleanerConfigFromGoogleCloudSecretManager(ctx context.Context) (*DBCleanerConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[DBCleanerConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetEmailProberConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetEmailProberConfigFromGoogleCloudSecretManager(ctx context.Context) (*EmailProberConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[EmailProberConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database = dbconfig.Config{
		ConnectionDetails:        " ",
		OAuth2TokenEncryptionKey: " ",
	}
	cfg.Email = emailcfg.Config{
		Provider: emailcfg.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{
			APIToken: os.Getenv(gcpSendgridTokenEnvVarKey),
		},
	}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetMealPlanFinalizerConfigFromGoogleCloudSecretManager fetches an MealPlanFinalizerConfig from GCP Secret Manager.
func GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx context.Context) (*MealPlanFinalizerConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[MealPlanFinalizerConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetMealPlanGroceryListInitializerWorkerConfigFromGoogleCloudSecretManager fetches an MealPlanGroceryListInitializerConfig from GCP Secret Manager.
func GetMealPlanGroceryListInitializerWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*MealPlanGroceryListInitializerConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[MealPlanGroceryListInitializerConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager fetches an MealPlanTaskCreatorConfig from GCP Secret Manager.
func GetMealPlanTaskCreatorWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*MealPlanTaskCreatorConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[MealPlanTaskCreatorConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager fetches an SearchDataIndexSchedulerConfig from GCP Secret Manager.
func GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager(ctx context.Context) (*SearchDataIndexSchedulerConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[SearchDataIndexSchedulerConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// END JOBS, BEGIN FUNCTIONS

// getWorkerConfigFromGoogleCloudSecretManager fetches a config from GCP Secret Manager.
func getWorkerConfigFromGoogleCloudSecretManager[T configurations](ctx context.Context) (*T, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	var cfg *T
	configBytes, err := fetchSecretFromSecretStore(ctx, client, "api_service_config")
	if err != nil {
		return nil, fmt.Errorf("fetching config from secret store: %w", err)
	}

	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	return cfg, nil
}

// GetDataChangesWorkerConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx context.Context) (*APIServiceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[APIServiceConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{
		Segment:  &segment.Config{APIToken: os.Getenv(gcpSegmentTokenEnvVarKey)},
		Posthog:  &posthog.Config{APIKey: os.Getenv(gcpPostHogKeyEnvVarKey)},
		Provider: analyticscfg.ProviderSegment,
	}
	cfg.Email = emailcfg.Config{}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetOutboundEmailerConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetOutboundEmailerConfigFromGoogleCloudSecretManager(ctx context.Context) (*APIServiceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[APIServiceConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{
		Provider: emailcfg.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{
			APIToken: os.Getenv(gcpSendgridTokenEnvVarKey),
		},
	}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetSearchDataIndexerConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetSearchDataIndexerConfigFromGoogleCloudSecretManager(ctx context.Context) (*APIServiceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[APIServiceConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetUserDataAggregatorConfigFromGoogleCloudSecretManager fetches an APIServiceConfig from GCP Secret Manager.
func GetUserDataAggregatorConfigFromGoogleCloudSecretManager(ctx context.Context) (*APIServiceConfig, error) {
	cfg, err := getWorkerConfigFromGoogleCloudSecretManager[APIServiceConfig](ctx)
	if err != nil {
		return nil, err
	}

	cfg.Search = searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			APIKey: os.Getenv(gcpAlgoliaAPIKeyEnvVarKey),
			AppID:  os.Getenv(gcpAlgoliaAppIDEnvVarKey),
		},
	}

	cfg.Database.ConnectionDetails = buildDatabaseURIFromGCPEnvVars()
	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = os.Getenv(gcpOauth2TokenEncryptionKeyEnvVarKey)
	cfg.Analytics = analyticscfg.Config{}
	cfg.Email = emailcfg.Config{}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}
