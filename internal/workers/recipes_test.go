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

func TestWritesWorker_createRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipeDatabaseCreationInput(),
		}

		expectedRecipe := fakes.BuildFakeRecipe()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
		).Return(expectedRecipe, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipeDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
		).Return((*types.Recipe)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipeDatabaseCreationInput(),
		}

		expectedRecipe := fakes.BuildFakeRecipe()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
		).Return(expectedRecipe, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipe(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
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

		assert.NoError(t, worker.updateRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating recipe", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipe(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.RecipeDataType,
			Recipe:   fakes.BuildFakeRecipe(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			body.Recipe,
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

		assert.Error(t, worker.updateRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.NoError(t, worker.archiveRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManager.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			body.RecipeID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeDataType,
		}

		dbManager := database.NewMockDatabase()
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

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveRecipe(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
