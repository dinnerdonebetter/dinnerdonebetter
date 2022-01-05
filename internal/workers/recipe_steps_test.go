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

func TestWritesWorker_createRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStepDatabaseCreationInput(),
		}

		expectedRecipeStep := fakes.BuildFakeRecipeStep()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"CreateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return(expectedRecipeStep, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStepDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"CreateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStepDatabaseCreationInput(),
		}

		expectedRecipeStep := fakes.BuildFakeRecipeStep()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"CreateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return(expectedRecipeStep, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStep(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"UpdateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postUpdatesPublisher

		assert.NoError(t, worker.updateRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating recipe step", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStep(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"UpdateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:   types.RecipeStepDataType,
			RecipeStep: fakes.BuildFakeRecipeStep(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"UpdateRecipeStep",
			testutils.ContextMatcher,
			body.RecipeStep,
		).Return(nil)

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postUpdatesPublisher

		assert.Error(t, worker.updateRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepDataManager.On(
			"ArchiveRecipeStep",
			testutils.ContextMatcher,
			body.RecipeID,
			body.RecipeStepID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveRecipeStep(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
