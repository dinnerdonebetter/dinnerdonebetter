package datachangemessagehandler

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

func (a *AsyncDataChangeMessageHandler) handleQueueTestMessage(
	ctx context.Context,
	logger logging.Logger,
	span tracing.Span,
	changeMessage *audit.DataChangeMessage,
) error {
	testID, ok := changeMessage.Context["test_id"].(string)
	if !ok || testID == "" {
		return fmt.Errorf("missing or invalid test_id in queue test message context")
	}

	queueName, ok := changeMessage.Context["queue_name"].(string)
	if !ok || queueName == "" {
		return fmt.Errorf("missing or invalid queue name in queue test message context")
	}

	logger.WithValue("test_id", testID).Info("acknowledging queue test message")

	if err := a.internalOpsRepo.AcknowledgeQueueTestMessage(ctx, testID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "acknowledging queue test message")
	}

	if queueName != "" {
		if err := a.internalOpsRepo.PruneQueueTestMessages(ctx, queueName); err != nil {
			observability.AcknowledgeError(err, logger, span, "pruning queue test messages")
		}
	}

	return nil
}
