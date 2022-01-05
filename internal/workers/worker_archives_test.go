package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
)

func TestProvidePreArchivesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
		dataChangesPublisher := &mockpublishers.Publisher{}

		actual, err := ProvideArchivesWorker(
			ctx,
			logger,
			dbManager,
			dataChangesPublisher,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestArchivesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		dbManager := database.NewMockDatabase()
		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		dbManager := database.NewMockDatabase()
		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		body := &types.PreArchiveMessage{
			DataType: types.UserMembershipDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
