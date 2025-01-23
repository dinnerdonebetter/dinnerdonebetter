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
	testutils2 "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidMeasurementUnitsService_CreateValidMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"CreateValidMeasurementUnit",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnit, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
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

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"CreateValidMeasurementUnit",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidMeasurementUnitsService_ReadValidMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ReadValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ReadValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement unit in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ReadValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ReadValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})
}

func TestValidMeasurementUnitsService_ListValidMeasurementUnitsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitsList()

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnits",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidMeasurementUnitList, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ListValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidMeasurementUnitList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidMeasurementUnitList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ListValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnits",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ListValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving valid measurement units from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnits",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ListValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})
}

func TestValidMeasurementUnitsService_SearchValidMeasurementUnitsHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"SearchForValidMeasurementUnits",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return(exampleValidMeasurementUnitList.Data, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidMeasurementUnitList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
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
		validMeasurementUnitSearchSubsets := make([]*types.ValidMeasurementUnitSearchSubset, len(exampleValidMeasurementUnitList.Data))
		for i := range exampleValidMeasurementUnitList.Data {
			expectedIDs = append(expectedIDs, exampleValidMeasurementUnitList.Data[i].ID)
			validMeasurementUnitSearchSubsets[i] = converters.ConvertValidMeasurementUnitToValidMeasurementUnitSearchSubset(exampleValidMeasurementUnitList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[types.ValidMeasurementUnitSearchSubset]{}
		searchIndex.On(
			"TextSearch",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return(validMeasurementUnitSearchSubsets, nil)
		helper.service.validMeasurementUnitSearchIndex = searchIndex

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnitsWithIDs",
			testutils2.ContextMatcher,
			expectedIDs,
		).Return(exampleValidMeasurementUnitList.Data, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidMeasurementUnitList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager, searchIndex)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.SearchValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
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

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"SearchForValidMeasurementUnits",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidMeasurementUnit{}, sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"SearchForValidMeasurementUnits",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidMeasurementUnit(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})
}

func TestVaValidMeasurementUnitslidMeasurementUnitsService_SearchByIngredientIDHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils2.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidMeasurementUnitList, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsByIngredientIDHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleValidMeasurementUnitList.Data)
		assert.Equal(t, *actual.Pagination, exampleValidMeasurementUnitList.Pagination)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.SearchValidMeasurementUnitsByIngredientIDHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
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

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils2.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsByIngredientIDHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils2.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.SearchValidMeasurementUnitsByIngredientIDHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[[]*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})
}

func TestValidMeasurementUnitsService_UpdateValidMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"UpdateValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidMeasurementUnitUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
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

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement unit", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error retrieving valid measurement unit from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"GetValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"UpdateValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidMeasurementUnitsService_ArchiveValidMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitExists",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ArchiveValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid measurement unit in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitExists",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validMeasurementUnitDataManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitExists",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validMeasurementUnitDataManager

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validMeasurementUnitDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitExists",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ArchiveValidMeasurementUnit",
			testutils2.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
