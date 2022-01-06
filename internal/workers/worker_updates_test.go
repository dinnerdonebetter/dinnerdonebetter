package workers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func TestProvidePreUpdatesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
		dataChangesPublisher := &mockpublishers.Publisher{}

		actual, err := ProvideUpdatesWorker(
			ctx,
			logger,
			dbManager,
			dataChangesPublisher,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestUpdatesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		dbManager := database.NewMockDatabase()
		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
