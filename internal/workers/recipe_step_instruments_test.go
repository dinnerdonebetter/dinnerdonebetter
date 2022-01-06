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

func TestWritesWorker_createRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error updating recipe step instrument", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.updateRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_archiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

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

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.archiveRecipeStepInstrument(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
