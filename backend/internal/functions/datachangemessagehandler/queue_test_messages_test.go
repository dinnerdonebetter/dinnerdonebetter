package datachangemessagehandler

import (
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	internalopsmock "github.com/dinnerdonebetter/backend/internal/domain/internalops/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAsyncDataChangeMessageHandler_handleQueueTestMessage(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		repo := &internalopsmock.InternalOpsDataManager{}
		handler.internalOpsRepo = repo

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id":    "test-123",
				"queue_name": "data-changes",
			},
		}

		repo.On(reflection.GetMethodName(repo.AcknowledgeQueueTestMessage), mock.Anything, "test-123").Return(nil).Once()
		repo.On(reflection.GetMethodName(repo.PruneQueueTestMessages), mock.Anything, "data-changes").Return(nil).Once()

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repo)
	})

	t.Run("missing test_id", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"queue_name": "data-changes",
			},
		}

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test_id")
	})

	t.Run("empty test_id", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id":    "",
				"queue_name": "data-changes",
			},
		}

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test_id")
	})

	t.Run("missing queue_name", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id": "test-123",
			},
		}

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queue name")
	})

	t.Run("empty queue_name", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id":    "test-123",
				"queue_name": "",
			},
		}

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queue name")
	})

	t.Run("nil context map", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		changeMessage := &audit.DataChangeMessage{
			Context: nil,
		}

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test_id")
	})

	t.Run("acknowledge error", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		repo := &internalopsmock.InternalOpsDataManager{}
		handler.internalOpsRepo = repo

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id":    "test-123",
				"queue_name": "data-changes",
			},
		}

		repo.On(reflection.GetMethodName(repo.AcknowledgeQueueTestMessage), mock.Anything, "test-123").Return(errors.New("db error")).Once()

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.Error(t, err)
		mock.AssertExpectationsForObjects(t, repo)
	})

	t.Run("prune error is not fatal", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		_, span := tracing.NewTracerForTest(t.Name()).StartSpan(ctx)
		logger := logging.NewNoopLogger()

		repo := &internalopsmock.InternalOpsDataManager{}
		handler.internalOpsRepo = repo

		changeMessage := &audit.DataChangeMessage{
			Context: map[string]any{
				"test_id":    "test-123",
				"queue_name": "data-changes",
			},
		}

		repo.On(reflection.GetMethodName(repo.AcknowledgeQueueTestMessage), mock.Anything, "test-123").Return(nil).Once()
		repo.On(reflection.GetMethodName(repo.PruneQueueTestMessages), mock.Anything, "data-changes").Return(errors.New("prune error")).Once()

		err := handler.handleQueueTestMessage(ctx, logger, span, changeMessage)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, repo)
	})
}
