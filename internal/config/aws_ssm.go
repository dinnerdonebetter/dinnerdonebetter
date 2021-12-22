package config

import (
	"context"
	"encoding/json"
	"strings"

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
	writesQueueNameSSMKey          = "PRIXFIXE_WRITES_QUEUE_URL"
	updatesQueueNameSSMKey         = "PRIXFIXE_UPDATES_QUEUE_URL"
	archivesQueueNameSSMKey        = "PRIXFIXE_ARCHIVES_QUEUE_URL"
	dataChangesQueueNameSSMKey     = "PRIXFIXE_DATA_CHANGES_QUEUE_URL"
	cookieBlockKeySSMKey           = "PRIXFIXE_COOKIE_BLOCK_KEY"
	cookieHashKeySSMKey            = "PRIXFIXE_COOKIE_HASH_KEY"
	cookiePASETOLocalModeKeySSMKey = "PRIXFIXE_PASETO_LOCAL_MODE_KEY"
	pubsubServerURLSSMKey          = "PRIXFIXE_PUBSUB_SERVER_URLS"
	pubsubServerUsernameSSMKey     = "PRIXFIXE_PUBSUB_SERVER_USERNAME"
	/* #nosec G101 */
	pubsubServerPasswordSSMKey = "PRIXFIXE_PUBSUB_SERVER_PASSWORD"
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

	cfg.Events.Consumers.RedisConfig.QueueAddresses = strings.Split(mustGetParameter(parameterStore, pubsubServerURLSSMKey), ",")
	cfg.Events.Publishers.RedisConfig.QueueAddresses = strings.Split(mustGetParameter(parameterStore, pubsubServerURLSSMKey), ",")

	writesTopicName := mustGetParameter(parameterStore, writesQueueNameSSMKey)
	updatesTopicName := mustGetParameter(parameterStore, updatesQueueNameSSMKey)
	archivesTopicName := mustGetParameter(parameterStore, archivesQueueNameSSMKey)
	dataChangesTopicName := mustGetParameter(parameterStore, dataChangesQueueNameSSMKey)

	cfg.Services.Auth.Cookies.BlockKey = mustGetParameter(parameterStore, cookieBlockKeySSMKey)
	cfg.Services.Auth.Cookies.HashKey = mustGetParameter(parameterStore, cookieHashKeySSMKey)
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(mustGetParameter(parameterStore, cookiePASETOLocalModeKeySSMKey))

	cfg.Services.ValidInstruments.PreWritesTopicName = writesTopicName
	cfg.Services.ValidInstruments.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidInstruments.PreArchivesTopicName = archivesTopicName

	cfg.Services.ValidIngredients.PreWritesTopicName = writesTopicName
	cfg.Services.ValidIngredients.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidIngredients.PreArchivesTopicName = archivesTopicName

	cfg.Services.ValidPreparations.PreWritesTopicName = writesTopicName
	cfg.Services.ValidPreparations.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidPreparations.PreArchivesTopicName = archivesTopicName

	cfg.Services.MealPlanOptionVotes.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlanOptionVotes.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlanOptionVotes.PreArchivesTopicName = archivesTopicName

	cfg.Services.ValidIngredientPreparations.PreWritesTopicName = writesTopicName
	cfg.Services.ValidIngredientPreparations.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidIngredientPreparations.PreArchivesTopicName = archivesTopicName

	cfg.Services.Meals.PreWritesTopicName = writesTopicName
	cfg.Services.Meals.PreUpdatesTopicName = updatesTopicName
	cfg.Services.Meals.PreArchivesTopicName = archivesTopicName

	cfg.Services.Recipes.PreWritesTopicName = writesTopicName
	cfg.Services.Recipes.PreUpdatesTopicName = updatesTopicName
	cfg.Services.Recipes.PreArchivesTopicName = archivesTopicName

	cfg.Services.RecipeSteps.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeSteps.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeSteps.PreArchivesTopicName = archivesTopicName

	cfg.Services.RecipeStepInstruments.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepInstruments.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepInstruments.PreArchivesTopicName = archivesTopicName

	cfg.Services.RecipeStepIngredients.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepIngredients.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepIngredients.PreArchivesTopicName = archivesTopicName

	cfg.Services.MealPlans.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlans.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlans.PreArchivesTopicName = archivesTopicName

	cfg.Services.MealPlanOptions.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlanOptions.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlanOptions.PreArchivesTopicName = archivesTopicName

	cfg.Services.RecipeStepProducts.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepProducts.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepProducts.PreArchivesTopicName = archivesTopicName

	cfg.Services.Households.PreWritesTopicName = writesTopicName

	cfg.Services.HouseholdInvitations.PreWritesTopicName = writesTopicName

	cfg.Services.Webhooks.PreWritesTopicName = writesTopicName
	cfg.Services.Webhooks.PreArchivesTopicName = archivesTopicName

	cfg.Services.Websockets.DataChangesTopicName = dataChangesTopicName

	ctx := context.Background()
	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}
