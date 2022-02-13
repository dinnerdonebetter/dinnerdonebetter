package datachangesfunction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/prixfixeco/api_server/internal/observability"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/observability/keys"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data string `json:"data"`
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

	var changeMessage types.DataChangeMessage
	if unmarshalErr := json.Unmarshal([]byte(m.Data), &changeMessage); unmarshalErr != nil {
		observability.AcknowledgeError(unmarshalErr, logger, nil, "unmarshalling data change message")
	}

	switch changeMessage.MessageType {
	case "meal_plan_finalized":
		// alert family members or whatever
		if err = customerDataCollector.EventOccurred(ctx, "meal_plan_finalized", "meal-plan-finalizer", map[string]interface{}{
			keys.HouseholdIDKey: changeMessage.MealPlan.BelongsToHousehold,
			keys.MealPlanIDKey:  changeMessage.MealPlan.ID,
		}); err != nil {
			observability.AcknowledgeError(err, logger, nil, "notifying customer data platform")
		}
	}

	logger.WithValue("payload", changeMessage).Info("invoked")

	return nil
}
