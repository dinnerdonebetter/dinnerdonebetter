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

func TestWritesWorker_createValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrumentDatabaseCreationInput(),
		}

		expectedValidInstrument := fakes.BuildFakeValidInstrument()

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
		).Return(expectedValidInstrument, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrumentDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrumentDatabaseCreationInput(),
		}

		expectedValidInstrument := fakes.BuildFakeValidInstrument()

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
		).Return(expectedValidInstrument, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
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

		assert.NoError(t, worker.updateValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating valid instrument", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:        types.ValidInstrumentDataType,
			ValidInstrument: fakes.BuildFakeValidInstrument(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrument,
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

		assert.Error(t, worker.updateValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			body.ValidInstrumentID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.ValidInstrumentDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveValidInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
