package validenumerations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	encoding "github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidIngredientStateIngredientsService_CreateValidIngredientStateIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"CreateValidIngredientStateIngredient",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientStateIngredientDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidIngredientStateIngredient, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientStateIngredient)
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

		helper.service.CreateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientStateIngredientCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"CreateValidIngredientStateIngredient",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientStateIngredientDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidIngredientStateIngredient)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientStateIngredientsService_ReadValidIngredientStateIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(helper.exampleValidIngredientStateIngredient, nil)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ReadValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientStateIngredient)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return((*types.ValidIngredientStateIngredient)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ReadValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return((*types.ValidIngredientStateIngredient)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ReadValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})
}

func TestValidIngredientStateIngredientsService_ListValidIngredientStateIngredientsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientsList()

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientStateIngredientList, nil)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ListValidIngredientStateIngredientsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientStateIngredientList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidIngredientStateIngredientList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidIngredientStateIngredientsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidIngredientStateIngredient])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ListValidIngredientStateIngredientsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error retrieving valid ingredient preparations from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredients",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidIngredientStateIngredient])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ListValidIngredientStateIngredientsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})
}

func TestValidIngredientStateIngredientsService_UpdateValidIngredientStateIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(helper.exampleValidIngredientStateIngredient, nil)

		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"UpdateValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientStateIngredient)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientStateIngredientUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
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

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return((*types.ValidIngredientStateIngredient)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error retrieving valid ingredient preparation from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return((*types.ValidIngredientStateIngredient)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateIngredientUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(helper.exampleValidIngredientStateIngredient, nil)

		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"UpdateValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientStateIngredientsService_ArchiveValidIngredientStateIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ValidIngredientStateIngredientExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(true, nil)

		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ArchiveValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ValidIngredientStateIngredientExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ArchiveValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ValidIngredientStateIngredientExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.ArchiveValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ValidIngredientStateIngredientExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(true, nil)

		dbManager.ValidIngredientStateIngredientDataManagerMock.On(
			"ArchiveValidIngredientStateIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredientStateIngredient.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidIngredientStateIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientStateIngredientsService_SearchValidIngredientStateIngredientsByIngredientHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientsList()

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredientsForIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleValidIngredientStateIngredientList, nil)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.SearchValidIngredientStateIngredientsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientStateIngredientList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidIngredientStateIngredientList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidIngredientStateIngredientsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredientsForIngredient",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidIngredientStateIngredient])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.SearchValidIngredientStateIngredientsByIngredientHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})
}

func TestValidIngredientStateIngredientsService_SearchValidIngredientStateIngredientsByPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientStateIngredientList := fakes.BuildFakeValidIngredientStateIngredientsList()

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredientsForIngredientState",
			testutils.ContextMatcher,
			helper.exampleValidIngredientState.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleValidIngredientStateIngredientList, nil)
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.SearchValidIngredientStateIngredientsByIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientStateIngredientList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidIngredientStateIngredientList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidIngredientStateIngredientsByIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateIngredientDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateIngredientDataManager.ValidIngredientStateIngredientDataManagerMock.On(
			"GetValidIngredientStateIngredientsForIngredientState",
			testutils.ContextMatcher,
			helper.exampleValidIngredientState.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidIngredientStateIngredient])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateIngredientDataManager

		helper.service.SearchValidIngredientStateIngredientsByIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientStateIngredient]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateIngredientDataManager)
	})
}
