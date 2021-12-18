package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparationDatabaseCreationInput(),
		}

		expectedValidPreparation := fakes.BuildFakeValidPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(expectedValidPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidPreparation.ID,
			expectedValidPreparation,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher, searchIndexManager)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparationDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparationDatabaseCreationInput(),
		}

		expectedValidPreparation := fakes.BuildFakeValidPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(expectedValidPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidPreparation.ID,
			expectedValidPreparation,
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager

		assert.Error(t, worker.createValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparationDatabaseCreationInput(),
		}

		expectedValidPreparation := fakes.BuildFakeValidPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(expectedValidPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidPreparation.ID,
			expectedValidPreparation,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_updateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidPreparation.ID,
			body.ValidPreparation,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.NoError(t, worker.updateValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher, searchIndexManager)
	})

	T.Run("with error updating valid preparation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidPreparation.ID,
			body.ValidPreparation,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager

		assert.Error(t, worker.updateValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:         types.ValidPreparationDataType,
			ValidPreparation: fakes.BuildFakeValidPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidPreparation.ID,
			body.ValidPreparation,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.Error(t, worker.updateValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_archiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			body.ValidPreparationID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager

		assert.Error(t, worker.archiveValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, searchIndexManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.validPreparationsIndexManager = searchIndexManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveValidPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}
