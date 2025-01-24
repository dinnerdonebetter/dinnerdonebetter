package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	testutils "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanOptionsService_CreateMealPlanOptionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"CreateMealPlanOption",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanOption, nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanOption)
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

		helper.service.CreateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealPlanOptionCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"CreateMealPlanOption",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanOptionDatabaseCreationInput) bool { return true }),
		).Return((*types.MealPlanOption)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.CreateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealPlanOptionsService_ReadMealPlanOptionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(helper.exampleMealPlanOption, nil)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ReadMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanOption)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan option in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return((*types.MealPlanOption)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ReadMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return((*types.MealPlanOption)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ReadMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})
}

func TestMealPlanOptionsService_ListMealPlanOptionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionsList()

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOptions",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleMealPlanOptionList, nil)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ListMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealPlanOptionList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOptions",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanOption])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ListMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error retrieving meal plan options from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOptions",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanOption])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ListMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})
}

func TestMealPlanOptionsService_UpdateMealPlanOptionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(helper.exampleMealPlanOption, nil)

		dbManager.MealPlanOptionDataManagerMock.On(
			"UpdateMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlanOption,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanOption)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealPlanOptionUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("without input attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan option", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return((*types.MealPlanOption)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error retrieving meal plan option from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return((*types.MealPlanOption)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"GetMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(helper.exampleMealPlanOption, nil)

		dbManager.MealPlanOptionDataManagerMock.On(
			"UpdateMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlanOption,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.UpdateMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealPlanOptionsService_ArchiveMealPlanOptionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"MealPlanOptionExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionDataManagerMock.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan option in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"MealPlanOptionExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(false, nil)
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ArchiveMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanOptionDataManager.MealPlanOptionDataManagerMock.On(
			"MealPlanOptionExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanOptionDataManager

		helper.service.ArchiveMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanOptionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanOptionDataManagerMock.On(
			"MealPlanOptionExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(true, nil)

		dbManager.MealPlanOptionDataManagerMock.On(
			"ArchiveMealPlanOption",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
			helper.exampleMealPlanOption.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.ArchiveMealPlanOptionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanOption]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
