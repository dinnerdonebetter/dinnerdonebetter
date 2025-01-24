package validenumerations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mocksearch "github.com/dinnerdonebetter/backend/internal/search/text/mock"
	testutils "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidInstrumentsService_CreateValidInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManagerMock.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidInstrument, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidInstrument)
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

		helper.service.CreateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidInstrumentCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManagerMock.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidInstrumentDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidInstrumentsService_ReadValidInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(helper.exampleValidInstrument, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ReadValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidInstrument)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ReadValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ReadValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})
}

func TestValidInstrumentsService_ListValidInstrumentsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentsList()

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidInstrumentList, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ListValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidInstrumentList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidInstrumentList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidInstrument])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ListValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving valid instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidInstrument])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ListValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})
}

func TestValidInstrumentsService_SearchValidInstrumentsHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidInstrumentList := fakes.BuildFakeValidInstrumentsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"SearchForValidInstruments",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(exampleValidInstrumentList.Data, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.SearchValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidInstrumentList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("using external service", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.useSearchService = true

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		expectedIDs := []string{}
		validInstrumentSearchSubsets := make([]*types.ValidInstrumentSearchSubset, len(exampleValidInstrumentList.Data))
		for i := range exampleValidInstrumentList.Data {
			expectedIDs = append(expectedIDs, exampleValidInstrumentList.Data[i].ID)
			validInstrumentSearchSubsets[i] = converters.ConvertValidInstrumentToValidInstrumentSearchSubset(exampleValidInstrumentList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[types.ValidInstrumentSearchSubset]{}
		searchIndex.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(validInstrumentSearchSubsets, nil)
		helper.service.validInstrumentSearchIndex = searchIndex

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrumentsWithIDs",
			testutils.ContextMatcher,
			expectedIDs,
		).Return(exampleValidInstrumentList.Data, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.SearchValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager, searchIndex)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.SearchValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"SearchForValidInstruments",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidInstrument{}, sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.SearchValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"SearchForValidInstruments",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidInstrument{}, errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.SearchValidInstrumentsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})
}

func TestValidInstrumentsService_UpdateValidInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(helper.exampleValidInstrument, nil)

		dbManager.ValidInstrumentDataManagerMock.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidInstrument) bool { return true }),
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidInstrument)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidInstrumentUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
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

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving valid instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidInstrumentDataManagerMock.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(helper.exampleValidInstrument, nil)

		dbManager.ValidInstrumentDataManagerMock.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidInstrument) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidInstrumentsService_ArchiveValidInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(true, nil)

		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, validInstrumentDataManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ArchiveValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ArchiveValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(true, nil)

		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			helper.exampleValidInstrument.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.ArchiveValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})
}

func TestValidInstrumentsService_RandomValidInstrumentHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetRandomValidInstrument",
			testutils.ContextMatcher,
		).Return(helper.exampleValidInstrument, nil)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.RandomValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidInstrument)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.RandomValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetRandomValidInstrument",
			testutils.ContextMatcher,
		).Return((*types.ValidInstrument)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.RandomValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validInstrumentDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validInstrumentDataManager.ValidInstrumentDataManagerMock.On(
			"GetRandomValidInstrument",
			testutils.ContextMatcher,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validInstrumentDataManager

		helper.service.RandomValidInstrumentHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidInstrument]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validInstrumentDataManager)
	})
}
