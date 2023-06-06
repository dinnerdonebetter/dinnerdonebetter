package mealplangrocerylistitems

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanGroceryListItemsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(helper.exampleMealPlanGroceryListItem, nil)
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealPlanGroceryListItem{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return((*types.MealPlanGroceryListItem)(nil), sql.ErrNoRows)
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return((*types.MealPlanGroceryListItem)(nil), errors.New("blah"))
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})
}

func TestMealPlanGroceryListItemsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanGroceryListItemList := fakes.BuildFakeMealPlanGroceryListItemList().Data

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return(exampleMealPlanGroceryListItemList, nil)
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.MealPlanGroceryListItem{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanGroceryListItem(nil), sql.ErrNoRows)
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.MealPlanGroceryListItem{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItemsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.MealPlanGroceryListItem(nil), errors.New("blah"))
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager, encoderDecoder)
	})
}

func TestMealPlanGroceryListItemsService_UpdateHandler(T *testing.T) {
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

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanGroceryListItemDataManager.On(
			"UpdateMealPlanGroceryListItem",
			testutils.ContextMatcher,
			mock.MatchedBy(func(input *types.MealPlanGroceryListItem) bool { return true }),
		).Return(nil)
		helper.service.mealPlanGroceryListItemDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
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

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(helper.exampleMealPlanGroceryListItem, nil)

		dbManager.MealPlanGroceryListItemDataManager.On(
			"UpdateMealPlanGroceryListItem",
			testutils.ContextMatcher,
			mock.MatchedBy(func(input *types.MealPlanGroceryListItem) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.mealPlanGroceryListItemDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
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

		dbManager := database.NewMockDatabase()
		dbManager.MealPlanGroceryListItemDataManager.On(
			"GetMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(expectedPrepStep, nil)

		dbManager.MealPlanGroceryListItemDataManager.On(
			"UpdateMealPlanGroceryListItem",
			testutils.ContextMatcher,
			mock.MatchedBy(func(input *types.MealPlanGroceryListItem) bool { return true }),
		).Return(nil)
		helper.service.mealPlanGroceryListItemDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestMealPlanGroceryListItemsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(true, nil)

		mealPlanGroceryListItemDataManager.On(
			"ArchiveMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(nil)
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error checking existence", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(false, errors.New("blah"))

		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})

	T.Run("with error archiving item", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanGroceryListItemDataManager := &mocktypes.MealPlanGroceryListItemDataManager{}
		mealPlanGroceryListItemDataManager.On(
			"MealPlanGroceryListItemExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(true, nil)

		mealPlanGroceryListItemDataManager.On(
			"ArchiveMealPlanGroceryListItem",
			testutils.ContextMatcher,
			helper.exampleMealPlanGroceryListItem.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanGroceryListItemDataManager = mealPlanGroceryListItemDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanGroceryListItemDataManager)
	})
}
