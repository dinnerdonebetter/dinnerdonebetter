package datachangemessagehandler

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

func (a *AsyncDataChangeMessageHandler) handleQueueTestMessage(
	ctx context.Context,
	logger logging.Logger,
	span tracing.Span,
	testID, topicName string,
) error {
	if testID == "" {
		return fmt.Errorf("missing or invalid test_id in queue test message")
	}

	logger.WithValue("test_id", testID).Info("acknowledging queue test message")

	if err := a.internalOpsRepo.AcknowledgeQueueTestMessage(ctx, testID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "acknowledging queue test message")
	}

	if topicName != "" {
		if err := a.internalOpsRepo.PruneQueueTestMessages(ctx, topicName); err != nil {
			observability.AcknowledgeError(err, logger, span, "pruning queue test messages")
		}
	}

	return nil
}
