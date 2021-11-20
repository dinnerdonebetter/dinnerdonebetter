package workers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		postWritesPublisher := &mockpublishers.Publisher{}
		postWritesPublisher.On(
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
			true,
		).Return(true, nil)

		postWritesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"FinalizeMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
			true,
		).Return(true, nil)

		postWritesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanDataType }),
		).Return(nil)

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postWritesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postWritesPublisher)
	})

	T.Run("with error writing vote", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		dataChangesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool {
				return message.DataType == types.MealPlanOptionVoteDataType
			}),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error finalizing meal plan option", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

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
			true,
		).Return(false, errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing message about meal plan option finalization", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

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
			true,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error finalizing meal plan", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

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
			true,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher, mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"FinalizeMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
			true,
		).Return(false, errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing message about meal plan finalization", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

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
			true,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanOptionDataType }),
		).Return(nil)

		dbManager.MealPlanDataManager.On(
			"FinalizeMealPlan",
			testutils.ContextMatcher,
			body.MealPlanID,
			body.AttributableToHouseholdID,
			true,
		).Return(true, nil)

		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return message.DataType == types.MealPlanDataType }),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			dataChangesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestWritesWorker_updateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postUpdatesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error updating meal plan option vote", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postUpdatesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postUpdatesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		postUpdatesPublisher := &mockpublishers.Publisher{}
		postUpdatesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postUpdatesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.updateMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postUpdatesPublisher)
	})
}

func TestWritesWorker_archiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

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

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return &mocksearch.IndexManager{}, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.archiveMealPlanOptionVote(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
