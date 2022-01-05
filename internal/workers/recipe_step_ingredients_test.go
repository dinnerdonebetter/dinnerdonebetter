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

func TestWritesWorker_createRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredientDatabaseCreationInput(),
		}

		expectedRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return(expectedRecipeStepIngredient, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredientDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return((*types.RecipeStepIngredient)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredientDatabaseCreationInput(),
		}

		expectedRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"CreateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return(expectedRecipeStepIngredient, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"UpdateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.updateRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error updating recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"UpdateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:             types.RecipeStepIngredientDataType,
			RecipeStepIngredient: fakes.BuildFakeRecipeStepIngredient(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"UpdateRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepIngredient,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.updateRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_archiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.archiveRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepIngredientDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepIngredientDataManager.On(
			"ArchiveRecipeStepIngredient",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepIngredientID,
		).Return(nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.archiveRecipeStepIngredient(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
