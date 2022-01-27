package config

import (
	"context"
	"encoding/json"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/prixfixeco/api_server/internal/database"
)

const (
	baseAPIServerConfigSSMKey      = "PRIXFIXE_BASE_API_SERVER_CONFIG"
	baseWorkerConfigSSMKey         = "PRIXFIXE_BASE_WORKER_CONFIG"
	databaseConnectionURLSSMKey    = "PRIXFIXE_DATABASE_CONNECTION_STRING"
	dataChangesQueueNameSSMKey     = "PRIXFIXE_DATA_CHANGES_QUEUE_URL"
	cookieBlockKeySSMKey           = "PRIXFIXE_COOKIE_BLOCK_KEY"
	cookieHashKeySSMKey            = "PRIXFIXE_COOKIE_HASH_KEY"
	cookiePASETOLocalModeKeySSMKey = "PRIXFIXE_PASETO_LOCAL_MODE_KEY"
	/* #nosec G101 */
	sendgridAPITokenSSMKey = "PRIXFIXE_SENDGRID_API_TOKEN"
	/* #nosec G101 */
	segmentAPITokenSSMKey = "PRIXFIXE_SEGMENT_API_TOKEN"
)

func mustGetParameter(parameterStore *ssm.SSM, paramName string) string {
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}

	rawParam, err := parameterStore.GetParameter(input)
	if err != nil {
		panic(err)
	}

	return *rawParam.Parameter.Value
}

// GetConfigFromParameterStore fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromParameterStore(worker bool) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger().WithValue("worker", worker)
	logger.Debug("setting up ssm session client")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	parameterStore := ssm.New(sess)

	var rawPartialConfig string
	if worker {
		logger.Debug("fetching partial worker config")
		rawPartialConfig = mustGetParameter(parameterStore, baseWorkerConfigSSMKey)
	} else {
		logger.Debug("fetching partial server config")
		rawPartialConfig = mustGetParameter(parameterStore, baseAPIServerConfigSSMKey)
	}

	logger.Debug("unloading JSON")

	var cfg *InstanceConfig
	if err := json.Unmarshal([]byte(rawPartialConfig), &cfg); err != nil {
		return nil, err
	}

	// fetch supplementary data from SSM
	cfg.Database.ConnectionDetails = database.ConnectionDetails(mustGetParameter(parameterStore, databaseConnectionURLSSMKey))
	cfg.Email.APIToken = mustGetParameter(parameterStore, sendgridAPITokenSSMKey)
	cfg.CustomerData.APIToken = mustGetParameter(parameterStore, segmentAPITokenSSMKey)

	dataChangesTopicName := mustGetParameter(parameterStore, dataChangesQueueNameSSMKey)

	cfg.Services.Auth.Cookies.BlockKey = mustGetParameter(parameterStore, cookieBlockKeySSMKey)
	cfg.Services.Auth.Cookies.HashKey = mustGetParameter(parameterStore, cookieHashKeySSMKey)
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(mustGetParameter(parameterStore, cookiePASETOLocalModeKeySSMKey))

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

	ctx := context.Background()
	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}
