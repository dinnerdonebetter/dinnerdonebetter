package searchindexer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dinnerdonebetter/backend/internal/asyncfunc"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/search/text/indexing"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("IndexDataForSearch", IndexDataForSearch)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// IndexDataForSearch handles a data change.
func IndexDataForSearch(ctx context.Context, e event.Event) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.GetSearchDataIndexerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger := cfg.Observability.Logging.ProvideLogger()

	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	var searchIndexRequest *indexing.IndexRequest
	if err = json.Unmarshal(msg.Message.Data, &searchIndexRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, nil, "unmarshaling data change message")
	}

	if err = asyncfunc.HandleIndexDataRequest(ctx, logger, cfg, searchIndexRequest); err != nil {
		observability.AcknowledgeError(err, logger, nil, "handling index request")
	}

	return nil
}
