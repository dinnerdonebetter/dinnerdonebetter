package datachangesfunction

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/observability"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/observability/keys"

	"github.com/prixfixeco/api_server/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data types.DataChangeMessage `json:"data"`
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

	switch m.Data.MessageType {
	case "meal_plan_finalized":
		// alert family members or whatever
		if err = customerDataCollector.EventOccurred(ctx, "meal_plan_finalized", "meal-plan-finalizer", map[string]interface{}{
			keys.HouseholdIDKey: m.Data.MealPlan.BelongsToHousehold,
			keys.MealPlanIDKey:  m.Data.MealPlan.ID,
		}); err != nil {
			observability.AcknowledgeError(err, logger, nil, "notifying customer data platform")
		}
	}

	logger.WithValue("payload", m.Data).Info("invoked")

	return nil
}
