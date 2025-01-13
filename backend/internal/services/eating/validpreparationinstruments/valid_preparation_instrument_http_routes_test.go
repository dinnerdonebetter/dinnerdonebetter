package validpreparationinstruments

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidPreparationInstrumentsService_CreateValidPreparationInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidPreparationInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidPreparationInstrument)
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

		helper.service.CreateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationInstrumentCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidPreparationInstrumentDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.CreateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidPreparationInstrumentsService_ReadValidPreparationInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ReadValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidPreparationInstrument)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ReadValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ReadValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_ListValidPreparationInstrumentsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentsList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ListValidPreparationInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidPreparationInstrumentList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidPreparationInstrumentList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidPreparationInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ListValidPreparationInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving valid preparation instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ListValidPreparationInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_UpdateValidPreparationInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidPreparationInstrument)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationInstrumentUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
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

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving valid preparation instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument,
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.UpdateValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidPreparationInstrumentsService_ArchiveValidPreparationInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)

		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ArchiveValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ArchiveValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)

		dbManager.ValidPreparationInstrumentDataManagerMock.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = dbManager

		helper.service.ArchiveValidPreparationInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidPreparationInstrumentsService_SearchValidPreparationInstrumentsByPreparationHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentsList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleValidPreparationInstrumentList, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchValidPreparationInstrumentsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidPreparationInstrumentList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidPreparationInstrumentList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidPreparationInstrumentsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchValidPreparationInstrumentsByPreparationHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_SearchValidPreparationInstrumentsByInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentsList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
			testutils.QueryFilterMatcher,
		).Return(exampleValidPreparationInstrumentList, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchValidPreparationInstrumentsByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidPreparationInstrumentList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidPreparationInstrumentList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = func(request *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.SearchValidPreparationInstrumentsByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error fetching data from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManagerMock{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrumentsForInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
			testutils.QueryFilterMatcher,
		).Return((*types.QueryFilteredResult[types.ValidPreparationInstrument])(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.SearchValidPreparationInstrumentsByInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidPreparationInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}
