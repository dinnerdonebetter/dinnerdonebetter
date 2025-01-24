package validenumerations

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
	testutils "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidMeasurementUnitConversionsService_CreateValidMeasurementUnitConversionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"CreateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnitConversion, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
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

		helper.service.CreateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitConversionCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"CreateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversionDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidMeasurementUnitConversionsService_ReadValidMeasurementUnitConversionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ReadValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}

func TestValidMeasurementUnitConversionsService_UpdateValidMeasurementUnitConversionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"UpdateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnitConversion)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitConversionUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
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

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving valid measurement conversion from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with problem writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitConversionUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(helper.exampleValidMeasurementUnitConversion, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"UpdateValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidMeasurementUnitConversionsService_ArchiveValidMeasurementUnitConversionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ArchiveValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ArchiveValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ArchiveValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ValidMeasurementUnitConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"ArchiveValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnitConversion.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidMeasurementUnitConversionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidMeasurementUnitConversionsService_ValidMeasurementUnitConversionsFromMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		expected := []*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementUnitConversion}

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(expected, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsFromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, expected)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ValidMeasurementUnitConversionsFromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsFromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsFromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}

func TestValidMeasurementUnitConversionsService_ValidMeasurementUnitConversionsToMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		expected := []*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementUnitConversion}

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(expected, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, expected)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ValidMeasurementUnitConversionsToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitConversionDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitConversionDataManager.ValidMeasurementUnitConversionDataManagerMock.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitConversionDataManager

		helper.service.ValidMeasurementUnitConversionsToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnitConversion]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitConversionDataManager)
	})
}
