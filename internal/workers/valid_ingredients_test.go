package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.dataChangesPublisher = dataChangesPublisher

		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher, searchIndexManager)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_updateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.NoError(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher, searchIndexManager)
	})

	T.Run("with error updating valid ingredient", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.Error(t, worker.updateValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_archiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientDataManager.On(
			"ArchiveValidIngredient",
			testutils.ContextMatcher,
			body.ValidIngredientID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validIngredientsIndexManager = searchIndexManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveValidIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}
