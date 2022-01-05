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

func TestWritesWorker_createMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlanDatabaseCreationInput(),
		}

		expectedMealPlan := fakes.BuildFakeMealPlan()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"CreateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
		).Return(expectedMealPlan, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlanDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"CreateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
		).Return((*types.MealPlan)(nil), errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.createMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlanDatabaseCreationInput(),
		}

		expectedMealPlan := fakes.BuildFakeMealPlan()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"CreateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
		).Return(expectedMealPlan, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlan(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"UpdateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
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

		assert.NoError(t, worker.updateMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating meal plan", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlan(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"UpdateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
		).Return(errors.New("blah"))

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.updateMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType: types.MealPlanDataType,
			MealPlan: fakes.BuildFakeMealPlan(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"UpdateMealPlan",
			testutils.ContextMatcher,
			body.MealPlan,
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

		assert.Error(t, worker.updateMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
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

		assert.NoError(t, worker.archiveMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager

		assert.Error(t, worker.archiveMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanDataManager.On(
			"ArchiveMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
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

		assert.Error(t, worker.archiveMealPlan(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
