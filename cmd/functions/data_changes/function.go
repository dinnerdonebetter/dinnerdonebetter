package datachangesfunction

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/pkg/types"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Base64EncodedDataChangeMessage string `json:"data"`
}

// ProcessDataChange handles a data change.
func ProcessDataChange(ctx context.Context, m PubSubMessage) error {
	logger := zerolog.NewZerologLogger()

	cfg, err := config.GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	customerDataCollector, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		return fmt.Errorf("error setting up customer data collector: %w", err)
	}

	rawMessage, err := base64.StdEncoding.DecodeString(m.Base64EncodedDataChangeMessage)
	if err != nil {
		return fmt.Errorf("decoding raw pubsub message: %w", err)
	}

	var changeMessage types.DataChangeMessage
	if unmarshalErr := json.Unmarshal(rawMessage, &changeMessage); unmarshalErr != nil {
		logger = logger.WithValue("raw_data", rawMessage)
		observability.AcknowledgeError(unmarshalErr, logger, nil, "unmarshalling data change message")
	}

	logger = logger.WithValue("event_type", changeMessage.EventType)

	eventContext := map[string]interface{}{
		keys.HouseholdIDKey: changeMessage.AttributableToHouseholdID,
		keys.UserIDKey:      changeMessage.AttributableToUserID,
	}

	for k, v := range changeMessage.Context {
		eventContext[k] = v
	}

	switch changeMessage.EventType {
	case types.UserSignedUpCustomerEventType:
		if err = customerDataCollector.AddUser(ctx, changeMessage.AttributableToUserID, eventContext); err != nil {
			return observability.PrepareError(err, logger, nil, "notifying customer data platform")
		}
		return nil
	case types.UserLoggedInCustomerEventType:
		break
	case types.UserChangedActiveHouseholdCustomerEventType:
		break
	case types.UserLoggedOutCustomerEventType:
		break
	case types.TwoFactorSecretVerifiedCustomerEventType:
		break
	case types.APIClientCreatedCustomerEventType:
		eventContext[keys.APIClientDatabaseIDKey] = changeMessage.APIClientID
	case types.APIClientArchivedCustomerEventType:
		eventContext[keys.APIClientDatabaseIDKey] = changeMessage.APIClientID
	case types.HouseholdCreatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.Household.ID
	case types.HouseholdUpdatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.Household.ID
	case types.HouseholdArchivedCustomerEventType:
		break
	case types.HouseholdMembershipPermissionsUpdatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.HouseholdID
	case types.HouseholdInvitationCreatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.HouseholdID
	case types.HouseholdInvitationCanceledCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.HouseholdID
	case types.HouseholdInvitationAcceptedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.HouseholdID
	case types.HouseholdInvitationRejectedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.HouseholdID
	case types.MealPlanCreatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.MealPlan.BelongsToHousehold
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlan.ID
	case types.MealPlanUpdatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.MealPlan.BelongsToHousehold
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlan.ID
	case types.MealPlanArchivedCustomerEventType:
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlanID
	case types.MealPlanFinalizedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.MealPlan.BelongsToHousehold
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlanID
	case types.MealPlanOptionVoteCreatedCustomerEventType:
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlanID
		eventContext[keys.MealPlanOptionVoteIDKey] = changeMessage.MealPlanOptionVote.ID
		eventContext[keys.MealPlanOptionIDKey] = changeMessage.MealPlanOptionVote.BelongsToMealPlanOption
	case types.MealPlanOptionVoteUpdatedCustomerEventType:
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlanID
		eventContext[keys.MealPlanOptionVoteIDKey] = changeMessage.MealPlanOptionVote.ID
		eventContext[keys.MealPlanOptionIDKey] = changeMessage.MealPlanOptionVote.BelongsToMealPlanOption
	case types.MealPlanOptionVoteArchivedCustomerEventType:
		eventContext[keys.MealPlanIDKey] = changeMessage.MealPlanID
		eventContext[keys.MealPlanOptionVoteIDKey] = changeMessage.MealPlanOptionVoteID
		eventContext[keys.MealPlanOptionIDKey] = changeMessage.MealPlanOptionID
	case types.MealCreatedCustomerEventType:
		eventContext[keys.MealIDKey] = changeMessage.Meal.ID
	case types.MealArchivedCustomerEventType:
		eventContext[keys.MealIDKey] = changeMessage.MealID
	case types.RecipeCreatedCustomerEventType:
		eventContext[keys.RecipeIDKey] = changeMessage.Recipe.ID
	case types.RecipeUpdatedCustomerEventType:
		eventContext[keys.RecipeIDKey] = changeMessage.Recipe.ID
	case types.RecipeArchivedCustomerEventType:
		eventContext[keys.RecipeIDKey] = changeMessage.RecipeID
	default:
		logger.Debug("unknown event type")
	}

	if dataCollectionErr := customerDataCollector.EventOccurred(ctx, changeMessage.EventType, changeMessage.AttributableToUserID, eventContext); dataCollectionErr != nil {
		observability.AcknowledgeError(dataCollectionErr, logger, nil, "notifying customer data platform")
	}

	return nil
}
