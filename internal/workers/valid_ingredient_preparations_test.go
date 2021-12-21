package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		expectedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(expectedValidIngredientPreparation, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		expectedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(expectedValidIngredientPreparation, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.NoError(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.Error(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
