package validmeasurementunits

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/internal/search/text/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

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
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnit, nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

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

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = dbManager

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
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
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidMeasurementUnitDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidMeasurementUnit, nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidMeasurementUnitsService_ReadValidMeasurementUnitHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnits",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidMeasurementUnitList, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnits",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnits",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"SearchForValidMeasurementUnits",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(exampleValidMeasurementUnitList.Data, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(validMeasurementUnitSearchSubsets, nil)
		helper.service.validMeasurementUnitSearchIndex = searchIndex

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnitsWithIDs",
			testutils.ContextMatcher,
			expectedIDs,
		).Return(exampleValidMeasurementUnitList.Data, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"SearchForValidMeasurementUnits",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidMeasurementUnit{}, sql.ErrNoRows)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"SearchForValidMeasurementUnits",
			testutils.ContextMatcher,
			exampleQuery,
		).Return([]*types.ValidMeasurementUnit(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidMeasurementUnitList, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"ValidMeasurementUnitsForIngredientID",
			testutils.ContextMatcher,
			helper.exampleValidIngredient.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidMeasurementUnit])(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"UpdateValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit,
		).Return(nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), sql.ErrNoRows)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"GetValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return((*types.ValidMeasurementUnit)(nil), errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"UpdateValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit,
		).Return(errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = dbManager

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
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
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(helper.exampleValidMeasurementUnit, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"UpdateValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit,
		).Return(nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidMeasurementUnit)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ArchiveValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"ValidMeasurementUnitExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(false, nil)
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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

		validMeasurementUnitDataManager := &mocktypes.ValidMeasurementUnitDataManagerMock{}
		validMeasurementUnitDataManager.On(
			"ValidMeasurementUnitExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(false, errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = validMeasurementUnitDataManager

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
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ArchiveValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(errors.New("blah"))
		helper.service.validMeasurementUnitDataManager = dbManager

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ValidMeasurementUnitExists",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(true, nil)

		dbManager.ValidMeasurementUnitDataManagerMock.On(
			"ArchiveValidMeasurementUnit",
			testutils.ContextMatcher,
			helper.exampleValidMeasurementUnit.ID,
		).Return(nil)
		helper.service.validMeasurementUnitDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidMeasurementUnitHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidMeasurementUnit]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
