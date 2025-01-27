package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealsService_CreateMealHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealDataManagerMock.On(
			"CreateMeal",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMeal, nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMeal)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealDataManagerMock.On(
			"CreateMeal",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealDatabaseCreationInput) bool { return true }),
		).Return((*types.Meal)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.CreateMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealsService_ReadMealHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(helper.exampleMeal, nil)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ReadMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMeal)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return((*types.Meal)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ReadMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return((*types.Meal)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ReadMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})
}

func TestMealsService_ListMealsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealList := fakes.BuildFakeMealsList()

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleMealList, nil)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ListMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.Meal])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ListMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error retrieving meals from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.Meal])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ListMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})
}

func TestMealsService_SearchMealsHandler(T *testing.T) {
	T.Parallel()

	const exampleQuery = "example"
	exampleMealList := fakes.BuildFakeMealsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{textsearch.QueryKeySearch: []string{exampleQuery}}.Encode()

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"SearchForMeals",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleMealList, nil)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.SearchMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("using external service", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.useSearchService = true

		exampleLimit := uint8(123)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		expectedIDs := []string{}
		mealSearchSubsets := make([]*indexing.MealSearchSubset, len(exampleMealList.Data))
		for i := range exampleMealList.Data {
			expectedIDs = append(expectedIDs, exampleMealList.Data[i].ID)
			mealSearchSubsets[i] = indexing.ConvertMealToMealSearchSubset(exampleMealList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[indexing.MealSearchSubset]{}
		searchIndex.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(mealSearchSubsets, nil)
		helper.service.searchIndex = searchIndex

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"GetMealsWithIDs",
			testutils.ContextMatcher,
			expectedIDs,
		).Return(exampleMealList.Data, nil)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.SearchMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager, searchIndex)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.SearchMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{textsearch.QueryKeySearch: []string{exampleQuery}}.Encode()

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"SearchForMeals",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.Meal])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.SearchMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{textsearch.QueryKeySearch: []string{exampleQuery}}.Encode()

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"SearchForMeals",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.Meal])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.SearchMealsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})
}

func TestMealsService_ArchiveMealHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealDataManagerMock.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(true, nil)

		dbManager.MealDataManagerMock.On(
			"ArchiveMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(false, nil)
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ArchiveMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealDataManager.MealDataManagerMock.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanningDataManager = mealDataManager

		helper.service.ArchiveMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealDataManagerMock.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(true, nil)

		dbManager.MealDataManagerMock.On(
			"ArchiveMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.ArchiveMealHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Meal]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
