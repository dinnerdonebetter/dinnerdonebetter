package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	testutils "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHouseholdInstrumentOwnershipsService_CreateHouseholdInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"CreateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleHouseholdInstrumentOwnership)
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

		helper.service.CreateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.HouseholdInstrumentOwnershipCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"CreateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.CreateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestHouseholdInstrumentOwnershipsService_ReadHouseholdInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ReadHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleHouseholdInstrumentOwnership)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such household instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ReadHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ReadHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})
}

func TestHouseholdInstrumentOwnershipsService_ListHouseholdInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleHouseholdInstrumentOwnershipList := fakes.BuildFakeHouseholdInstrumentOwnershipsList()

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleHouseholdInstrumentOwnershipList, nil)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ListHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleHouseholdInstrumentOwnershipList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.HouseholdInstrumentOwnership])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ListHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving household instrument ownerships from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.HouseholdInstrumentOwnership])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ListHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})
}

func TestHouseholdInstrumentOwnershipsService_UpdateHouseholdInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)

		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"UpdateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnership) bool { return true }),
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleHouseholdInstrumentOwnership)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.HouseholdInstrumentOwnershipUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
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

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such household instrument ownership", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving household instrument ownership from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)

		dbManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"UpdateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnership) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.UpdateHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestHouseholdInstrumentOwnershipsService_ArchiveHouseholdInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"ArchiveHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such household instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(false, nil)
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ArchiveHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ArchiveHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipDataManagerMock.On(
			"ArchiveHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = householdInstrumentOwnershipDataManager

		helper.service.ArchiveHouseholdInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})
}
