package householdinstrumentownerships

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

func TestHouseholdInstrumentOwnershipsService_CreateHandler(T *testing.T) {
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
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"CreateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)
		helper.service.householdInstrumentOwnershipDataManager = dbManager

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

		exampleCreationInput := &types.HouseholdInstrumentOwnershipCreationRequestInput{}
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

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()
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

		exampleCreationInput := fakes.BuildFakeHouseholdInstrumentOwnershipCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"CreateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
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
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"CreateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnershipDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)
		helper.service.householdInstrumentOwnershipDataManager = dbManager

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

func TestHouseholdInstrumentOwnershipsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdInstrumentOwnership{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
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

	T.Run("with no such household instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
	})
}

func TestHouseholdInstrumentOwnershipsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleHouseholdInstrumentOwnershipList := fakes.BuildFakeHouseholdInstrumentOwnershipList()

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInstrumentOwnershipList, nil)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
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

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInstrumentOwnership])(nil), sql.ErrNoRows)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
	})

	T.Run("with error retrieving household instrument ownerships from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnerships",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInstrumentOwnership])(nil), errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
	})
}

func TestHouseholdInstrumentOwnershipsService_UpdateHandler(T *testing.T) {
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
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)

		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"UpdateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnership) bool { return true }),
		).Return(nil)
		helper.service.householdInstrumentOwnershipDataManager = dbManager

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

		exampleCreationInput := &types.HouseholdInstrumentOwnershipUpdateRequestInput{}
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

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), sql.ErrNoRows)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

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

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInstrumentOwnership)(nil), errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

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
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)

		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"UpdateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnership) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
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
		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"GetHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInstrumentOwnership, nil)

		dbManager.HouseholdInstrumentOwnershipDataManager.On(
			"UpdateHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInstrumentOwnership) bool { return true }),
		).Return(nil)
		helper.service.householdInstrumentOwnershipDataManager = dbManager

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

func TestHouseholdInstrumentOwnershipsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		householdInstrumentOwnershipDataManager.On(
			"ArchiveHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(nil)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, dataChangesPublisher)
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

	T.Run("with no such household instrument ownership in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(false, nil)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(false, errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		householdInstrumentOwnershipDataManager.On(
			"ArchiveHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(errors.New("blah"))
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdInstrumentOwnershipDataManager := &mocktypes.HouseholdInstrumentOwnershipDataManager{}
		householdInstrumentOwnershipDataManager.On(
			"HouseholdInstrumentOwnershipExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(true, nil)

		householdInstrumentOwnershipDataManager.On(
			"ArchiveHouseholdInstrumentOwnership",
			testutils.ContextMatcher,
			helper.exampleHouseholdInstrumentOwnership.ID,
			helper.exampleHousehold.ID,
		).Return(nil)
		helper.service.householdInstrumentOwnershipDataManager = householdInstrumentOwnershipDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdInstrumentOwnershipDataManager, dataChangesPublisher)
	})
}