package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOptionDatabaseCreationInput(),
		}

		expectedMealPlanOption := fakes.BuildFakeMealPlanOption()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"CreateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
		).Return(expectedMealPlanOption, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOptionDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"CreateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
		).Return((*types.MealPlanOption)(nil), errors.New("blah"))

		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOptionDatabaseCreationInput(),
		}

		expectedMealPlanOption := fakes.BuildFakeMealPlanOption()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"CreateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
		).Return(expectedMealPlanOption, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOption(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"UpdateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
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

		assert.NoError(t, worker.updateMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating meal plan option", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOption(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"UpdateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
		).Return(errors.New("blah"))

		postUpdatesPublisher := &mockpublishers.Publisher{}

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.postUpdatesPublisher = postUpdatesPublisher

		assert.Error(t, worker.updateMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:       types.MealPlanOptionDataType,
			MealPlanOption: fakes.BuildFakeMealPlanOption(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"UpdateMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanOption,
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

		assert.Error(t, worker.updateMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.MealPlanOptionID,
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

		assert.NoError(t, worker.archiveMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.MealPlanOptionID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.postArchivesPublisher = postArchivesPublisher

		assert.Error(t, worker.archiveMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionDataManager.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.MealPlanOptionID,
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

		assert.Error(t, worker.archiveMealPlanOption(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
