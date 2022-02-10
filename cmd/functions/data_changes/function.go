package datachangesfunction

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProcessDataChange handles a data change.
func ProcessDataChange(ctx context.Context, m PubSubMessage) error {
	logger := zerolog.NewZerologLogger()

	logger.Info("invoked")

	return nil
}
