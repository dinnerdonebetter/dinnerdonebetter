package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/prixfixeco/api_server/internal/database"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

func mustGetSecretParameter(ctx context.Context, client *secretmanager.Client, name string) string {
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: name})
	if err != nil {
		panic(err)
	}

	return string(result.Payload.Data)
}

// GetConfigFromCloudSecretManager fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Debug("setting up secret manager client")

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	var rawPartialConfig string
	logger.Debug("fetching partial server config")
	rawPartialConfig = mustGetSecretParameter(ctx, client, "api_service_config")

	var cfg *InstanceConfig
	if unmarshalErr := json.Unmarshal([]byte(rawPartialConfig), &cfg); unmarshalErr != nil {
		return nil, fmt.Errorf("error unmarshalling configuration: %w", unmarshalErr)
	}

	port, portParseErr := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if portParseErr != nil {
		panic(portParseErr)
	}
	cfg.Server.HTTPPort = uint16(port)

	cfg.Database.ConnectionDetails = database.ConnectionDetails(mustGetSecretParameter(ctx, client, "database_connection_string"))
	cfg.Email.APIToken = mustGetSecretParameter(ctx, client, "sendgrid_api_token")
	cfg.CustomerData.APIToken = mustGetSecretParameter(ctx, client, "segment_api_token")

	dataChangesTopicName := "data_changes_topic_name"

	cfg.Services.Auth.Cookies.BlockKey = mustGetSecretParameter(ctx, client, "cookie_block_key")
	cfg.Services.Auth.Cookies.HashKey = mustGetSecretParameter(ctx, client, "cookie_hash_key")
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(mustGetSecretParameter(ctx, client, "paseto_local_key"))

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
	cfg.Services.Websockets.DataChangesTopicName = dataChangesTopicName

	if validationErr := cfg.ValidateWithContext(ctx); validationErr != nil {
		return nil, fmt.Errorf("error validating configuration: %w", validationErr)
	}

	return cfg, nil
}
