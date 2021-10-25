package workers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	mockpublishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	mocksearch "gitlab.com/prixfixe/prixfixe/internal/search/mock"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
)

func TestProvidePreArchivesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		actual, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, errors.New("blah")
		}

		actual, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}

func TestPreArchivesWorker_HandleMessage(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := database.BuildMockDatabase()
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, []byte("} bad JSON lol")))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with ValidInstrumentDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidInstrumentDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with ValidInstrumentDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidInstrumentDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with ValidIngredientDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidPreparationDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidPreparationDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with ValidPreparationDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidPreparationDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientPreparationDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientPreparationDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with ValidIngredientPreparationDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with ValidIngredientPreparationDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			body.RecipeID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			body.RecipeID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with RecipeDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			body.RecipeID,
			body.AttributableToHouseholdID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			body.RecipeID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			body.RecipeID,
			body.RecipeStepID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			body.RecipeID,
			body.RecipeStepID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with RecipeStepDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			body.RecipeID,
			body.RecipeStepID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			body.RecipeID,
			body.RecipeStepID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepInstrumentDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepInstrumentID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepInstrumentID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepInstrumentDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepInstrumentID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with RecipeStepInstrumentDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepInstrumentID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepInstrumentID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepInstrumentDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepInstrumentID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepInstrumentID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepIngredientDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepIngredientID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepIngredientDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with RecipeStepIngredientDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepIngredientID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepIngredientDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepIngredientID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepProductDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepProductID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepProductDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with RecipeStepProductDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepProductID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with RecipeStepProductDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.RecipeStepProductID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with MealPlanDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with MealPlanOptionDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
			body.AttributableToHouseholdID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionVoteDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionVoteDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with MealPlanOptionVoteDataType and error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with MealPlanOptionVoteDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with WebhookDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.WebhookDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.WebhookDataManager.On(
			"ArchiveWebhook",
			testutils.ContextMatcher,
			body.WebhookID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with WebhookDataType and error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.WebhookDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.WebhookDataManager.On(
			"ArchiveWebhook",
			testutils.ContextMatcher,
			body.WebhookID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexManager := &mocksearch.IndexManager{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with WebhookDataType and error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.WebhookDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		dbManager := database.BuildMockDatabase()
		dbManager.WebhookDataManager.On(
			"ArchiveWebhook",
			testutils.ContextMatcher,
			body.WebhookID,
			body.AttributableToHouseholdID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with UserMembershipDataType", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}
		dbManager := database.BuildMockDatabase()
		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvidePreArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		body := &types.PreArchiveMessage{
			DataType: types.UserMembershipDataType,
		}
		examplePayload, err := json.Marshal(body)
		require.NoError(t, err)

		assert.NoError(t, worker.HandleMessage(ctx, examplePayload))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
