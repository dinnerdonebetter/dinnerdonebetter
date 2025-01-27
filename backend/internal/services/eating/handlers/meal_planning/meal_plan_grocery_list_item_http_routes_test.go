package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	testutils "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanGroceryListItemsService_ReadMealPlanGroceryListItemHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(helper.exampleMealPlanGroceryListItem, nil)
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ReadMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanGroceryListItem)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return((*types.MealPlanGroceryListItem)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ReadMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return((*types.MealPlanGroceryListItem)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ReadMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})
}

func TestMealPlanGroceryListItemsService_ListMealPlanGroceryListItemsByMealPlanHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanGroceryListItemList := fakes.BuildFakeMealPlanGroceryListItemsList().Data

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return(exampleMealPlanGroceryListItemList, nil)
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ListMealPlanGroceryListItemsByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealPlanGroceryListItemList)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListMealPlanGroceryListItemsByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanGroceryListItem(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ListMealPlanGroceryListItemsByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanGroceryListItem(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ListMealPlanGroceryListItemsByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})
}

func TestMealPlanGroceryListItemsService_UpdateMealPlanGroceryListItemHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(mealPlanGroceryListItem)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		expectedPrepStep := helper.exampleMealPlanGroceryListItem
		expectedPrepStep.Status = *exampleInput.Status
		expectedPrepStep.StatusExplanation = *exampleInput.StatusExplanation

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPut, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanGroceryListItemDataManagerMock.On(
			"UpdateMealPlanGroceryListItem",
			testutils.ContextMatcher,
			mock.MatchedBy(func(input *types.MealPlanGroceryListItem) bool { return true }),
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanGroceryListItem)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemUpdateRequestInput(mealPlanGroceryListItem)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPut, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanGroceryListItemDataManagerMock.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(helper.exampleMealPlanGroceryListItem, nil)

		dbManager.MealPlanGroceryListItemDataManagerMock.On(
			"UpdateMealPlanGroceryListItem",
			testutils.ContextMatcher,
			mock.MatchedBy(func(input *types.MealPlanGroceryListItem) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.UpdateMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealPlanGroceryListItemsService_ArchiveMealPlanGroceryListItemHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(true, nil)

		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"ArchiveMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error checking existence", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(false, errors.New("blah"))

		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ArchiveMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error archiving item", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(true, nil)

		mealPlanGroceryListItemDataManager.MealPlanGroceryListItemDataManagerMock.On(
			"ArchiveMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanGroceryListItemDataManager

		helper.service.ArchiveMealPlanGroceryListItemHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanGroceryListItem]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})
}
