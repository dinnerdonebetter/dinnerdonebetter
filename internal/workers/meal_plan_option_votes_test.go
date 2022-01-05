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

func TestWritesWorker_createMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanID:         fakes.BuildFakeID(),
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			expectedMealPlanOptionVote.BelongsToMealPlanOption,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeCompleteMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanDataType }),
		).Return(nil)

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.NoError(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error writing vote", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return((*types.MealPlanOptionVote)(nil), errors.New("blah"))

		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error finalizing meal plan option", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanID:         fakes.BuildFakeID(),
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			expectedMealPlanOptionVote.BelongsToMealPlanOption,
			body.AttributableToHouseholdID,
		).Return(false, errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing message about meal plan option finalization", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanID:         fakes.BuildFakeID(),
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			expectedMealPlanOptionVote.BelongsToMealPlanOption,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error finalizing meal plan", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanID:         fakes.BuildFakeID(),
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			expectedMealPlanOptionVote.BelongsToMealPlanOption,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeCompleteMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(false, errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing message about meal plan finalization", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreWriteMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanID:         fakes.BuildFakeID(),
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput(),
		}

		expectedMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"CreateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(expectedMealPlanOptionVote, nil)

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(nil)

		dbManager.MealPlanOptionDataManager.On(
			"FinalizeMealPlanOption",
			testutils.ContextMatcher,
			body.MealPlanID,
			expectedMealPlanOptionVote.BelongsToMealPlanOption,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"AttemptToFinalizeCompleteMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanDataType }),
		).Return(errors.New("blah"))

		worker := newTestWritesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVote(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
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

		assert.NoError(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error updating meal plan option vote", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVote(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
		).Return(errors.New("blah"))

		dataChangesPublisher := &mockpublishers.Publisher{}

		worker := newTestUpdatesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreUpdateMessage{
			DataType:           types.MealPlanOptionVoteDataType,
			MealPlanOptionVote: fakes.BuildFakeMealPlanOptionVote(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"UpdateMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVote,
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

		assert.Error(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_archiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
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

		assert.NoError(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
		).Return(errors.New("blah"))

		dataChangesPublisher := &mockpublishers.Publisher{}
		worker := newTestArchivesWorker(t)
		worker.dataManager = dbManager
		worker.dataChangesPublisher = dataChangesPublisher

		assert.Error(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := &types.PreArchiveMessage{
			DataType: types.MealPlanOptionVoteDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanOptionVoteDataManager.On(
			"ArchiveMealPlanOptionVote",
			testutils.ContextMatcher,
			body.MealPlanOptionVoteID,
			body.AttributableToHouseholdID,
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

		assert.Error(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
