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

func TestWritesWorker_createRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProductDatabaseCreationInput(),
		}

		expectedRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
		).Return(expectedRecipeStepProduct, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProductDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
		).Return((*types.RecipeStepProduct)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProductDatabaseCreationInput(),
		}

		expectedRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"CreateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
		).Return(expectedRecipeStepProduct, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProduct(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"UpdateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
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

		assert.NoError(t, worker.updateRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating recipe step product", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProduct(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"UpdateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:          types.RecipeStepProductDataType,
			RecipeStepProduct: fakes.BuildFakeRecipeStepProduct(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"UpdateRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepProduct,
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

		assert.Error(t, worker.updateRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
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

		assert.NoError(t, worker.archiveRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.RecipeStepProductDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepProductDataManager.On(
			"ArchiveRecipeStepProduct",
			testutils.ContextMatcher,
			body.RecipeStepID,
			body.RecipeStepProductID,
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

		assert.Error(t, worker.archiveRecipeStepProduct(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
