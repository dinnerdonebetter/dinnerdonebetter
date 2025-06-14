package mealplanning

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccountInstrumentOwnershipsService_CreateAccountInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"CreateAccountInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.AccountInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleAccountInstrumentOwnership, nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleAccountInstrumentOwnership)
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

		helper.service.CreateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.AccountInstrumentOwnershipCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"CreateAccountInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.AccountInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return((*types.AccountInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.CreateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestAccountInstrumentOwnershipsService_ReadAccountInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(helper.exampleAccountInstrumentOwnership, nil)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ReadAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleAccountInstrumentOwnership)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such account instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return((*types.AccountInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ReadAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return((*types.AccountInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ReadAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})
}

func TestAccountInstrumentOwnershipsService_ListAccountInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAccountInstrumentOwnershipList := fakes.BuildFakeAccountInstrumentOwnershipsList()

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleAccount.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleAccountInstrumentOwnershipList, nil)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ListAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleAccountInstrumentOwnershipList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleAccount.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.AccountInstrumentOwnership])(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ListAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving account instrument ownerships from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleAccount.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.AccountInstrumentOwnership])(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ListAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})
}

func TestAccountInstrumentOwnershipsService_UpdateAccountInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(helper.exampleAccountInstrumentOwnership, nil)

		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"UpdateAccountInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.AccountInstrumentOwnership) bool { return true }),
		).Return(nil)
		helper.service.mealPlanningDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleAccountInstrumentOwnership)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.AccountInstrumentOwnershipUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
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

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such account instrument ownership", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return((*types.AccountInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error retrieving account instrument ownership from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return((*types.AccountInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeAccountInstrumentOwnershipUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"GetAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(helper.exampleAccountInstrumentOwnership, nil)

		dbManager.AccountInstrumentOwnershipDataManagerMock.On(
			"UpdateAccountInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.AccountInstrumentOwnership) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = dbManager

		helper.service.UpdateAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestAccountInstrumentOwnershipsService_ArchiveAccountInstrumentOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"AccountInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(true, nil)

		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"ArchiveAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(nil)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such account instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"AccountInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(false, nil)
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ArchiveAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"AccountInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ArchiveAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		accountInstrumentOwnershipDataManager := mocktypes.NewMealPlanningDataManagerMock()
		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"AccountInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(true, nil)

		accountInstrumentOwnershipDataManager.AccountInstrumentOwnershipDataManagerMock.On(
			"ArchiveAccountInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleAccountInstrumentOwnership.ID,
			helper.exampleAccount.ID,
		).Return(errors.New("blah"))
		helper.service.mealPlanningDataManager = accountInstrumentOwnershipDataManager

		helper.service.ArchiveAccountInstrumentOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.AccountInstrumentOwnership]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, accountInstrumentOwnershipDataManager)
	})
}
