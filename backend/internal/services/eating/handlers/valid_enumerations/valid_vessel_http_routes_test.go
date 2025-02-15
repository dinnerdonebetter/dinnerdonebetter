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
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	mocksearch "github.com/dinnerdonebetter/backend/internal/lib/search/text/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidVesselsService_CreateValidVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidVesselDataManagerMock.On(
			"CreateValidVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidVesselDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidVessel, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidVessel)
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

		helper.service.CreateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidVesselCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidVesselDataManagerMock.On(
			"CreateValidVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidVesselDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidVessel)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidVesselsService_ReadValidVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(helper.exampleValidVessel, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ReadValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidVessel)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid vessel in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return((*types.ValidVessel)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ReadValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return((*types.ValidVessel)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ReadValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})
}

func TestValidVesselsService_ListValidVesselsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidVesselList := fakes.BuildFakeValidVesselsList()

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessels",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleValidVesselList, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ListValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidVesselList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidVesselList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessels",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidVessel])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ListValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving valid vessels from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessels",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidVessel])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ListValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})
}

func TestValidVesselsService_SearchValidVesselsHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidVesselList := fakes.BuildFakeValidVesselsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"SearchForValidVessels",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(exampleValidVesselList.Data, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.SearchValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidVesselList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("using external service", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.useSearchService = true

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		expectedIDs := []string{}
		validVesselSearchSubsets := make([]*indexing.ValidVesselSearchSubset, len(exampleValidVesselList.Data))
		for i := range exampleValidVesselList.Data {
			expectedIDs = append(expectedIDs, exampleValidVesselList.Data[i].ID)
			validVesselSearchSubsets[i] = indexing.ConvertValidVesselToValidVesselSearchSubset(exampleValidVesselList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[indexing.ValidVesselSearchSubset]{}
		searchIndex.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(validVesselSearchSubsets, nil)
		helper.service.validVesselsSearchIndex = searchIndex

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVesselsWithIDs",
			testutils.ContextMatcher,
			expectedIDs,
		).Return(exampleValidVesselList.Data, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.SearchValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidVesselList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validVesselDataManager, searchIndex)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.SearchValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"SearchForValidVessels",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidVessel{}, sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.SearchValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"SearchForValidVessels",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidVessel{}, errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.SearchValidVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})
}

func TestValidVesselsService_UpdateValidVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(helper.exampleValidVessel, nil)

		dbManager.ValidVesselDataManagerMock.On(
			"UpdateValidVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidVessel) bool { return true }),
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidVessel)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidVesselUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
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

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid vessel", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return((*types.ValidVessel)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving valid vessel from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return((*types.ValidVessel)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidVesselDataManagerMock.On(
			"GetValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(helper.exampleValidVessel, nil)

		dbManager.ValidVesselDataManagerMock.On(
			"UpdateValidVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidVessel) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidVesselsService_ArchiveValidVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ValidVesselExists",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(true, nil)

		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ArchiveValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, validVesselDataManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid vessel in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ValidVesselExists",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ArchiveValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ValidVesselExists",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ArchiveValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ValidVesselExists",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(true, nil)

		validVesselDataManager.ValidVesselDataManagerMock.On(
			"ArchiveValidVessel",
			testutils.ContextMatcher,
			helper.exampleValidVessel.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.ArchiveValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})
}

func TestValidVesselsService_RandomValidVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetRandomValidVessel",
			testutils.ContextMatcher,
		).Return(helper.exampleValidVessel, nil)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.RandomValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidVessel)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.RandomValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid vessel in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetRandomValidVessel",
			testutils.ContextMatcher,
		).Return((*types.ValidVessel)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.RandomValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validVesselDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validVesselDataManager.ValidVesselDataManagerMock.On(
			"GetRandomValidVessel",
			testutils.ContextMatcher,
		).Return((*types.ValidVessel)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validVesselDataManager

		helper.service.RandomValidVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Nil(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validVesselDataManager)
	})
}
