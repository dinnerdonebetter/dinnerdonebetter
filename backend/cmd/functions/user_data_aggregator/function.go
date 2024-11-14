package webhookexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dinnerdonebetter/backend/cmd/functions/user_data_aggregator/logic"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("AggregateUserData", AggregateUserData)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// AggregateUserData handles a user data aggregation request.
func AggregateUserData(ctx context.Context, e event.Event) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.GetUserDataAggregatorConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger := cfg.Observability.Logging.ProvideLogger()

	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	var userDataCollectionRequest *types.UserDataAggregationRequest
	if err = json.Unmarshal(msg.Message.Data, &userDataCollectionRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, nil, "unmarshalling data change message")
	}

	if err = logic.CollectAndSaveUserData(ctx, logger, cfg, userDataCollectionRequest); err != nil {
		return observability.PrepareAndLogError(err, logger, nil, "collecting and saving user data")
	}

	return nil
}
