package datachangesfunction

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/prixfixeco/backend/internal/config"
	customerdataconfig "github.com/prixfixeco/backend/internal/customerdata/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/pkg/types"
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

	eventContext := map[string]any{
		keys.HouseholdIDKey: changeMessage.AttributableToHouseholdID,
		keys.UserIDKey:      changeMessage.AttributableToUserID,
	}

	for k, v := range changeMessage.Context {
		eventContext[k] = v
	}

	switch changeMessage.EventType {
	case types.UserSignedUpCustomerEventType:
		if err = customerDataCollector.AddUser(ctx, changeMessage.AttributableToUserID, eventContext); err != nil {
			return observability.PrepareError(err, nil, "notifying customer data platform")
		}
		return nil
	case types.APIClientCreatedCustomerEventType:
		eventContext[keys.APIClientDatabaseIDKey] = changeMessage.APIClientID
	case types.APIClientArchivedCustomerEventType:
		eventContext[keys.APIClientDatabaseIDKey] = changeMessage.APIClientID
	case types.HouseholdCreatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.Household.ID
	case types.HouseholdUpdatedCustomerEventType:
		eventContext[keys.HouseholdIDKey] = changeMessage.Household.ID
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
	case types.APIClientUpdatedCustomerEventType:
	case types.TwoFactorSecretVerifiedCustomerEventType:
	case types.UserLoggedInCustomerEventType:
	case types.UserLoggedOutCustomerEventType:
	case types.UserChangedActiveHouseholdCustomerEventType:
	case types.HouseholdArchivedCustomerEventType:
	case types.HouseholdMemberRemovedCustomerEventType:
	case types.HouseholdOwnershipTransferredCustomerEventType:
	case types.HouseholdUserMembershipCreatedCustomerEventType:
	case types.HouseholdUserMembershipUpdatedCustomerEventType:
	case types.HouseholdUserMembershipArchivedCustomerEventType:
	case types.MealUpdatedCustomerEventType:
	case types.MealPlanEventCreatedCustomerEventType:
	case types.MealPlanEventUpdatedCustomerEventType:
	case types.MealPlanEventArchivedCustomerEventType:
	case types.MealPlanOptionCreatedCustomerEventType:
	case types.MealPlanOptionUpdatedCustomerEventType:
	case types.MealPlanOptionArchivedCustomerEventType:
	case types.MealPlanOptionFinalizedCreatedCustomerEventType:
	case types.MealPlanTaskCreatedCustomerEventType:
	case types.MealPlanTaskStatusChangedCustomerEventType:
	case types.RecipeMediaCreatedCustomerEventType:
	case types.RecipeMediaUpdatedCustomerEventType:
	case types.RecipeMediaArchivedCustomerEventType:
	case types.RecipeStepCreatedCustomerEventType:
	case types.RecipeStepUpdatedCustomerEventType:
	case types.RecipeStepArchivedCustomerEventType:
	case types.RecipeStepIngredientCreatedCustomerEventType:
	case types.RecipeStepIngredientUpdatedCustomerEventType:
	case types.RecipeStepIngredientArchivedCustomerEventType:
	case types.RecipeStepInstrumentCreatedCustomerEventType:
	case types.RecipeStepInstrumentUpdatedCustomerEventType:
	case types.RecipeStepInstrumentArchivedCustomerEventType:
	case types.RecipeStepProductCreatedCustomerEventType:
	case types.RecipeStepProductUpdatedCustomerEventType:
	case types.RecipeStepProductArchivedCustomerEventType:
	case types.ValidIngredientCreatedCustomerEventType:
	case types.ValidIngredientUpdatedCustomerEventType:
	case types.ValidIngredientArchivedCustomerEventType:
	case types.ValidIngredientMeasurementUnitCreatedCustomerEventType:
	case types.ValidIngredientMeasurementUnitUpdatedCustomerEventType:
	case types.ValidIngredientMeasurementUnitArchivedCustomerEventType:
	case types.ValidIngredientPreparationCreatedCustomerEventType:
	case types.ValidIngredientPreparationUpdatedCustomerEventType:
	case types.ValidIngredientPreparationArchivedCustomerEventType:
	case types.ValidInstrumentCreatedCustomerEventType:
	case types.ValidInstrumentUpdatedCustomerEventType:
	case types.ValidInstrumentArchivedCustomerEventType:
	case types.ValidMeasurementConversionCreatedCustomerEventType:
	case types.ValidMeasurementConversionUpdatedCustomerEventType:
	case types.ValidMeasurementConversionArchivedCustomerEventType:
	case types.ValidMeasurementUnitCreatedCustomerEventType:
	case types.ValidMeasurementUnitUpdatedCustomerEventType:
	case types.ValidMeasurementUnitArchivedCustomerEventType:
	case types.ValidPreparationInstrumentCreatedCustomerEventType:
	case types.ValidPreparationInstrumentUpdatedCustomerEventType:
	case types.ValidPreparationInstrumentArchivedCustomerEventType:
	case types.WebhookCreatedCustomerEventType:
	case types.WebhookArchivedCustomerEventType:
		// TODO: flesh these out
	default:
		logger.Debug("unknown event type")
	}

	if dataCollectionErr := customerDataCollector.EventOccurred(ctx, changeMessage.EventType, changeMessage.AttributableToUserID, eventContext); dataCollectionErr != nil {
		observability.AcknowledgeError(dataCollectionErr, logger, nil, "notifying customer data platform")
	}

	return nil
}
