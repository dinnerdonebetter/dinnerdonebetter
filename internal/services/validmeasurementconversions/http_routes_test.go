package validmeasurementconversions

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidMeasurementConversionsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"CreateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementConversionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementConversion, nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
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

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"CreateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementConversionDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"CreateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementConversionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementConversion, nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementConversionsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(helper.exampleValidMeasurementConversion, nil)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidMeasurementUnitConversion{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
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

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})
}

func TestValidMeasurementConversionsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(helper.exampleValidMeasurementConversion, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"UpdateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without input attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with no such valid measurement conversion", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), sql.ErrNoRows)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager)
	})

	T.Run("with error retrieving valid measurement conversion from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return((*types.ValidMeasurementUnitConversion)(nil), errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager)
	})

	T.Run("with problem writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(helper.exampleValidMeasurementConversion, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"UpdateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementConversionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(helper.exampleValidMeasurementConversion, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"UpdateValidMeasurementConversion",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitConversion) bool { return true }),
		).Return(nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementConversionsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"ValidMeasurementConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"ArchiveValidMeasurementConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

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

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"ValidMeasurementConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(false, nil)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"ValidMeasurementConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(false, errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"ValidMeasurementConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"ArchiveValidMeasurementConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementConversionDataManager.On(
			"ValidMeasurementConversionExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementConversionDataManager.On(
			"ArchiveValidMeasurementConversion",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementConversion.ID,
		).Return(nil)
		helper.service.validMeasurementConversionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementConversionsService_FromMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementConversion}, nil)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.ValidMeasurementUnitConversion{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
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

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsFromUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.FromMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})
}

func TestValidMeasurementConversionsService_ToMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion{helper.exampleValidMeasurementConversion}, nil)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.ValidMeasurementUnitConversion{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
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

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such valid measurement conversion in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), sql.ErrNoRows)
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementConversionDataManager := &mocktypes.ValidMeasurementConversionDataManager{}
		validMeasurementConversionDataManager.On(
			"GetValidMeasurementUnitConversionsToUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return([]*types.ValidMeasurementUnitConversion(nil), errors.New("blah"))
		helper.service.validMeasurementConversionDataManager = validMeasurementConversionDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ToMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validMeasurementConversionDataManager, encoderDecoder)
	})
}
