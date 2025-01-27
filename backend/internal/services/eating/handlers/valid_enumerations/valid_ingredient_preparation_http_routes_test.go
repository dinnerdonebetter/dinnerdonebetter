package validenumerations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessioncontext"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	testutils "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidIngredientPreparationsService_CreateValidIngredientPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientPreparationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidIngredientPreparation, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientPreparation)
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

		helper.service.CreateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientPreparationCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientPreparationDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientPreparationsService_ReadValidIngredientPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ReadValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientPreparation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ReadValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ReadValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_ListValidIngredientPreparationsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationsList()

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ListValidIngredientPreparationsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientPreparationList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidIngredientPreparationList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidIngredientPreparationsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientPreparation])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ListValidIngredientPreparationsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error retrieving valid ingredient preparations from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientPreparation])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ListValidIngredientPreparationsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_UpdateValidIngredientPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)

		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientPreparation)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientPreparationUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
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

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error retrieving valid ingredient preparation from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)

		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientPreparationsService_ArchiveValidIngredientPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(true, nil)

		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ArchiveValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.ArchiveValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(true, nil)

		dbManager.ValidIngredientPreparationDataManagerMock.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidIngredientPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientPreparationsService_SearchValidIngredientPreparationsByIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleResponse := fakes.BuildFakeValidIngredientPreparationsList()

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparationsForIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleResponse, nil)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.SearchValidIngredientPreparationsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleResponse.Data)
		assert.Equal(t, *actual.Pagination, exampleResponse.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*sessioncontext.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidIngredientPreparationsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparationsForIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			testutils.QueryFilterMatcher,
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientPreparation])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.SearchValidIngredientPreparationsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_SearchValidIngredientPreparationsByPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleResponse := fakes.BuildFakeValidIngredientPreparationsList()

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparationsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleResponse, nil)
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.SearchValidIngredientPreparationsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleResponse.Data)
		assert.Equal(t, *actual.Pagination, exampleResponse.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*sessioncontext.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidIngredientPreparationsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientPreparationDataManager.ValidIngredientPreparationDataManagerMock.On(
			"GetValidIngredientPreparationsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientPreparation])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientPreparationDataManager

		helper.service.SearchValidIngredientPreparationsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientPreparation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}
