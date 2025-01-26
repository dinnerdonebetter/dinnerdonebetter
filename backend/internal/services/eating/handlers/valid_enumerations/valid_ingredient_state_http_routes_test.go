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
	"github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/mock"
	testutils2 "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidIngredientStatesService_CreateValidIngredientStateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"CreateValidIngredientState",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientStateDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidIngredientState, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientState)
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

		helper.service.CreateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientStateCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.CreateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"CreateValidIngredientState",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientStateDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidIngredientState)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientStatesService_ReadValidIngredientStateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(helper.exampleValidIngredientState, nil)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ReadValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientState)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ReadValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient state in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return((*types.ValidIngredientState)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ReadValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return((*types.ValidIngredientState)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ReadValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})
}

func TestValidIngredientStatesService_ListValidIngredientStatesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientStateList := fakes.BuildFakeValidIngredientStatesList()

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientStates",
			testutils2.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleValidIngredientStateList, nil)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ListValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientStateList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ListValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientStates",
			testutils2.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientState])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ListValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error retrieving valid ingredient states from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientStates",
			testutils2.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientState])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ListValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})
}

func TestValidIngredientStatesService_SearchValidIngredientStatesHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidIngredientStateList := fakes.BuildFakeValidIngredientStatesList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"SearchForValidIngredientStates",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return(exampleValidIngredientStateList.Data, nil)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.SearchValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidIngredientStateList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
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
		validIngredientStateSearchSubsets := make([]*types.ValidIngredientStateSearchSubset, len(exampleValidIngredientStateList.Data))
		for i := range exampleValidIngredientStateList.Data {
			expectedIDs = append(expectedIDs, exampleValidIngredientStateList.Data[i].ID)
			validIngredientStateSearchSubsets[i] = converters.ConvertValidIngredientStateToValidIngredientStateSearchSubset(exampleValidIngredientStateList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[types.ValidIngredientStateSearchSubset]{}
		searchIndex.On(
			"Search",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return(validIngredientStateSearchSubsets, nil)
		helper.service.validIngredientStatesSearchIndex = searchIndex

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientStatesWithIDs",
			testutils2.ContextMatcher,
			expectedIDs,
		).Return(exampleValidIngredientStateList.Data, nil)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.SearchValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager, searchIndex)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.SearchValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"SearchForValidIngredientStates",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidIngredientState{}, sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.SearchValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"SearchForValidIngredientStates",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidIngredientState{}, errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.SearchValidIngredientStatesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})
}

func TestValidIngredientStatesService_UpdateValidIngredientStateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(helper.exampleValidIngredientState, nil)

		dbManager.ValidIngredientStateDataManagerMock.On(
			"UpdateValidIngredientState",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientState) bool { return true }),
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientState)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientStateUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
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

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient state", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return((*types.ValidIngredientState)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error retrieving valid ingredient state from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return((*types.ValidIngredientState)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with problem writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"GetValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(helper.exampleValidIngredientState, nil)

		dbManager.ValidIngredientStateDataManagerMock.On(
			"UpdateValidIngredientState",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientState) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientStatesService_ArchiveValidIngredientStateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"ValidIngredientStateExists",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(true, nil)

		dbManager.ValidIngredientStateDataManagerMock.On(
			"ArchiveValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient state in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"ValidIngredientStateExists",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ArchiveValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientStateDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientStateDataManager.ValidIngredientStateDataManagerMock.On(
			"ValidIngredientStateExists",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientStateDataManager

		helper.service.ArchiveValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientStateDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientStateDataManagerMock.On(
			"ValidIngredientStateExists",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(true, nil)

		dbManager.ValidIngredientStateDataManagerMock.On(
			"ArchiveValidIngredientState",
			testutils2.ContextMatcher,
			helper.exampleValidIngredientState.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidIngredientStateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientState]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
