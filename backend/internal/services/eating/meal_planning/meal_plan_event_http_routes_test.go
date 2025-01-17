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
	testutils2 "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMealPlanEventsService_CreateMealPlanEventHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"CreateMealPlanEvent",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanEventDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleMealPlanEvent, nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanEvent)
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

		helper.service.CreateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealPlanEventCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.CreateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"CreateMealPlanEvent",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.MealPlanEventDatabaseCreationInput) bool { return true }),
		).Return((*types.MealPlanEvent)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.CreateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealPlanEventsService_ReadMealPlanEventHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(helper.exampleMealPlanEvent, nil)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ReadMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanEvent)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ReadMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return((*types.MealPlanEvent)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ReadMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return((*types.MealPlanEvent)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ReadMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})
}

func TestMealPlanEventsService_ListMealPlanEventHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanEventList := fakes.BuildFakeMealPlanEventsList()

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvents",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleMealPlanEventList, nil)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ListMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleMealPlanEventList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ListMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvents",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanEvent])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ListMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvents",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.MealPlanEvent])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ListMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})
}

func TestMealPlanEventsService_UpdateMealPlanEventHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(helper.exampleMealPlanEvent, nil)

		dbManager.MealPlanEventDataManagerMock.On(
			"UpdateMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlanEvent,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleMealPlanEvent)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealPlanEventUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
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

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return((*types.MealPlanEvent)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error retrieving meal plan from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return((*types.MealPlanEvent)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanEventUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"GetMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(helper.exampleMealPlanEvent, nil)

		dbManager.MealPlanEventDataManagerMock.On(
			"UpdateMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlanEvent,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.UpdateMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestMealPlanEventsService_ArchiveMealPlanEventHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"MealPlanEventExists",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanEventDataManagerMock.On(
			"ArchiveMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ArchiveMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"MealPlanEventExists",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(false, nil)
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ArchiveMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanEventDataManager := mocktypes.NewMealPlanningDataManagerMock()
		mealPlanEventDataManager.MealPlanEventDataManagerMock.On(
			"MealPlanEventExists",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanningDataManager = mealPlanEventDataManager

		helper.service.ArchiveMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, mealPlanEventDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := mocktypes.NewMealPlanningDataManagerMock()
		dbManager.MealPlanEventDataManagerMock.On(
			"MealPlanEventExists",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(true, nil)

		dbManager.MealPlanEventDataManagerMock.On(
			"ArchiveMealPlanEvent",
			testutils2.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanEvent.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.ArchiveMealPlanEventHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.MealPlanEvent]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
