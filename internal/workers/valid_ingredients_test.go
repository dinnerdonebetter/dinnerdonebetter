package workers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredientDatabaseCreationInput(),
		}

		expectedValidIngredient := fakes.BuildFakeValidIngredient()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(expectedValidIngredient, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredient.ID,
			expectedValidIngredient,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideWritesWorker(
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

		assert.NoError(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredientDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return((*types.ValidIngredient)(nil), errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
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

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredientDatabaseCreationInput(),
		}

		expectedValidIngredient := fakes.BuildFakeValidIngredient()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(expectedValidIngredient, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredient.ID,
			expectedValidIngredient,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
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

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredientDatabaseCreationInput(),
		}

		expectedValidIngredient := fakes.BuildFakeValidIngredient()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"CreateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(expectedValidIngredient, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredient.ID,
			expectedValidIngredient,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
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

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_updateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredient.ID,
			body.ValidIngredient,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideUpdatesWorker(
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

		assert.NoError(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error updating valid ingredient", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideUpdatesWorker(
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

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredient.ID,
			body.ValidIngredient,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideUpdatesWorker(
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

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:        types.ValidIngredientDataType,
			ValidIngredient: fakes.BuildFakeValidIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"UpdateValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredient,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredient.ID,
			body.ValidIngredient,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideUpdatesWorker(
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

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_archiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker, err := ProvideArchivesWorker(
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

		assert.NoError(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker, err := ProvideArchivesWorker(
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

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker, err := ProvideArchivesWorker(
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

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker, err := ProvideArchivesWorker(
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

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}
