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

func TestWritesWorker_createRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput(),
		}

		expectedRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
		).Return(expectedRecipeStepInstrument, nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.NoError(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
		).Return((*types.RecipeStepInstrument)(nil), errors.New("blah"))

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

		assert.Error(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput(),
		}

		expectedRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
		).Return(expectedRecipeStepInstrument, nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.Error(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}

func TestWritesWorker_updateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.NoError(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error updating recipe step instrument", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
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

		assert.Error(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepInstrumentDataType,
			RecipeStepInstrument: fakes.BuildFakeRecipeStepInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			body.RecipeStepInstrument,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.Error(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}

func TestWritesWorker_archiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.NoError(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
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

		assert.Error(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
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

		assert.Error(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
