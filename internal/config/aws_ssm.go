package config

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/prixfixeco/api_server/internal/database"
)

const (
	baseConfigSSMKey               = "PRIXFIXE_BASE_CONFIG"
	databaseConnectionURLSSMKey    = "PRIXFIXE_DATABASE_URL"
	writesQueueNameSSMKey          = "PRIXFIXE_WRITES_QUEUE_NAME"
	updatesQueueNameSSMKey         = "PRIXFIXE_UPDATES_QUEUE_NAME"
	archivesQueueNameSSMKey        = "PRIXFIXE_ARCHIVES_QUEUE_NAME"
	elasticsearchInstanceURLSSMKEy = "PRIXFIXE_ELASTICSEARCH_INSTANCE_URL"
	/* #nosec G101 */
	sendgridAPITokenSSMKey = "PRIXFIXE_SENDGRID_API_TOKEN"
	/* #nosec G101 */
	segmentAPITokenSSMKey = "PRIXFIXE_SEGMENT_API_TOKEN"
)

func mustGetParameter(ps *ssm.SSM, paramName string) string {
	input := &ssm.GetParameterInput{Name: aws.String(paramName)}
	rawParam, err := ps.GetParameter(input)
	if err != nil {
		panic(err)
	}

	return *rawParam.Parameter.Value
}

// GetConfigFromParameterStore fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromParameterStore() (*InstanceConfig, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ssm.New(sess)

	rawPartialConfig := mustGetParameter(svc, baseConfigSSMKey)

	var cfg *InstanceConfig
	if err := json.Unmarshal([]byte(rawPartialConfig), &cfg); err != nil {
		return nil, err
	}

	// fetch supplementary data from SSM
	cfg.Database.ConnectionDetails = database.ConnectionDetails(mustGetParameter(svc, databaseConnectionURLSSMKey))
	cfg.Email.APIToken = mustGetParameter(svc, sendgridAPITokenSSMKey)
	cfg.CustomerData.APIToken = mustGetParameter(svc, segmentAPITokenSSMKey)

	writesTopicName := mustGetParameter(svc, writesQueueNameSSMKey)
	updatesTopicName := mustGetParameter(svc, updatesQueueNameSSMKey)
	archivesTopicName := mustGetParameter(svc, archivesQueueNameSSMKey)
	elasticsearchInstanceURL := mustGetParameter(svc, elasticsearchInstanceURLSSMKEy)

	cfg.Services.ValidInstruments.PreWritesTopicName = writesTopicName
	cfg.Services.ValidIngredients.PreWritesTopicName = writesTopicName
	cfg.Services.ValidPreparations.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlanOptionVotes.PreWritesTopicName = writesTopicName
	cfg.Services.ValidIngredientPreparations.PreWritesTopicName = writesTopicName
	cfg.Services.Meals.PreWritesTopicName = writesTopicName
	cfg.Services.Recipes.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeSteps.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepInstruments.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepIngredients.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlans.PreWritesTopicName = writesTopicName
	cfg.Services.MealPlanOptions.PreWritesTopicName = writesTopicName
	cfg.Services.RecipeStepProducts.PreWritesTopicName = writesTopicName
	cfg.Services.Households.PreWritesTopicName = writesTopicName
	cfg.Services.HouseholdInvitations.PreWritesTopicName = writesTopicName
	cfg.Services.Webhooks.PreWritesTopicName = writesTopicName

	cfg.Services.ValidInstruments.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidIngredients.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidPreparations.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlanOptionVotes.PreUpdatesTopicName = updatesTopicName
	cfg.Services.ValidIngredientPreparations.PreUpdatesTopicName = updatesTopicName
	cfg.Services.Meals.PreUpdatesTopicName = updatesTopicName
	cfg.Services.Recipes.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeSteps.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepInstruments.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepIngredients.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlans.PreUpdatesTopicName = updatesTopicName
	cfg.Services.MealPlanOptions.PreUpdatesTopicName = updatesTopicName
	cfg.Services.RecipeStepProducts.PreUpdatesTopicName = updatesTopicName

	cfg.Services.ValidInstruments.PreArchivesTopicName = archivesTopicName
	cfg.Services.ValidIngredients.PreArchivesTopicName = archivesTopicName
	cfg.Services.ValidPreparations.PreArchivesTopicName = archivesTopicName
	cfg.Services.MealPlanOptionVotes.PreArchivesTopicName = archivesTopicName
	cfg.Services.ValidIngredientPreparations.PreArchivesTopicName = archivesTopicName
	cfg.Services.Meals.PreArchivesTopicName = archivesTopicName
	cfg.Services.Recipes.PreArchivesTopicName = archivesTopicName
	cfg.Services.RecipeSteps.PreArchivesTopicName = archivesTopicName
	cfg.Services.RecipeStepInstruments.PreArchivesTopicName = archivesTopicName
	cfg.Services.RecipeStepIngredients.PreArchivesTopicName = archivesTopicName
	cfg.Services.MealPlans.PreArchivesTopicName = archivesTopicName
	cfg.Services.MealPlanOptions.PreArchivesTopicName = archivesTopicName
	cfg.Services.RecipeStepProducts.PreArchivesTopicName = archivesTopicName
	cfg.Services.Webhooks.PreArchivesTopicName = archivesTopicName

	cfg.Services.ValidInstruments.SearchIndexPath = elasticsearchInstanceURL
	cfg.Services.ValidPreparations.SearchIndexPath = elasticsearchInstanceURL
	cfg.Services.ValidIngredients.SearchIndexPath = elasticsearchInstanceURL

	ctx := context.Background()
	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}
