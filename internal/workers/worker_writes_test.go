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
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestProvideWritesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}

func TestWritesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		worker := newTestWritesWorker(t)

		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))
	})

	T.Run("with WebhookDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.WebhookDataType,
			Webhook:  fakes.BuildFakeWebhookDatabaseCreationInput(),
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		expectedWebhook := fakes.BuildFakeWebhook()

		dbManager := database.NewMockDatabase()
		dbManager.WebhookDataManager.On(
			"CreateWebhook",
			testutils.ContextMatcher,
			body.Webhook,
		).Return(expectedWebhook, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
