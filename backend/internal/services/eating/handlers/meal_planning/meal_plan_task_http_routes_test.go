package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanTasksService_ReadMealPlanTaskHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(helper.exampleMealPlanTask, nil)
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ReadMealPlanTaskHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanTask)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadMealPlanTaskHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return((*types.MealPlanTask)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ReadMealPlanTaskHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return((*types.MealPlanTask)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ReadMealPlanTaskHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})
}

func TestMealPlanTasksService_ListMealPlanTasksByMealPlanHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanTaskList := fakes.BuildFakeMealPlanTasksList().Data

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return(exampleMealPlanTaskList, nil)
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ListMealPlanTasksByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealPlanTaskList)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListMealPlanTasksByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanTask(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ListMealPlanTasksByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanTaskDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanTaskDataManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTasksForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanTask(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanTaskDataManager

		helper.service.ListMealPlanTasksByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanTaskDataManager)
	})
}

func TestMealPlanTasksService_MealPlanTaskStatusChangeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleStatusChangeInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		exampleStatusChangeInput.ID = helper.exampleMealPlanTask.ID
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleStatusChangeInput)

		expectedPrepStep := helper.exampleMealPlanTask
		expectedPrepStep.Status = *exampleStatusChangeInput.Status
		expectedPrepStep.StatusExplanation = exampleStatusChangeInput.StatusExplanation

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanTaskDataManagerMock.On(
			"ChangeMealPlanTaskStatus",
			testutils.ContextMatcher,
			exampleStatusChangeInput,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.MealPlanTaskStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanTask)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.MealPlanTaskStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleCreationInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		exampleCreationInput.ID = helper.exampleMealPlanTask.ID
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanTaskDataManagerMock.On(
			"GetMealPlanTask",
			testutils.ContextMatcher,
			helper.exampleMealPlanTask.ID,
		).Return(helper.exampleMealPlanTask, nil)

		dbManager.MealPlanTaskDataManagerMock.On(
			"ChangeMealPlanTaskStatus",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.MealPlanTaskStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanTask]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
