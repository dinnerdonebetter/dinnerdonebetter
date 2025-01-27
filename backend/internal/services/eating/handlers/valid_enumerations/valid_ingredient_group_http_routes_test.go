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
	testutils "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	"github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidIngredientGroupsService_CreateValidIngredientGroupHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := mocktypes.NewValidEnumerationDataManagerMock()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"CreateValidIngredientGroup",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientGroupDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidIngredientGroup, nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientGroup)
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

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientGroupCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"CreateValidIngredientGroup",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientGroupDatabaseCreationInput) bool { return true }),
		).Return((*types.ValidIngredientGroup)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientGroupsService_ReadValidIngredientGroupHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(helper.exampleValidIngredientGroup, nil)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ReadValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientGroup)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ReadValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ReadValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})
}

func TestValidIngredientGroupsService_ListValidIngredientGroupsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientGroupList := fakes.BuildFakeValidIngredientGroupsList()

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleValidIngredientGroupList, nil)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ListValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientGroup])(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ListValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving valid ingredients from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.ValidIngredientGroup])(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ListValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})
}

func TestValidIngredientGroupsService_SearchValidIngredientGroupsHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidIngredientGroupList := fakes.BuildFakeValidIngredientGroupsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleValidIngredientGroupList.Data, nil)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.SearchValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.SearchValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
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

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return([]*types.ValidIngredientGroup{}, sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.SearchValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			textsearch.QueryKeySearch: []string{exampleQuery},
			filtering.QueryKeyLimit:   []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&filtering.QueryFilter{}),
		).Return([]*types.ValidIngredientGroup(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.SearchValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})
}

func TestValidIngredientGroupsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(helper.exampleValidIngredientGroup, nil)

		dbManager.ValidIngredientGroupDataManagerMock.On(
			"UpdateValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientGroupUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
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

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), sql.ErrNoRows)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving valid ingredient from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(helper.exampleValidIngredientGroup, nil)

		dbManager.ValidIngredientGroupDataManagerMock.On(
			"UpdateValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestValidIngredientGroupsService_ArchiveValidIngredientGroupHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(true, nil)

		dbManager.ValidIngredientGroupDataManagerMock.On(
			"ArchiveValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(nil)
		helper.service.validEnumerationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.MatchType[*types.DataChangeMessage](),
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid ingredient in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(false, nil)
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := mocktypes.NewValidEnumerationDataManagerMock()
		validIngredientGroupDataManager.ValidIngredientGroupDataManagerMock.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(false, errors.New("blah"))
		helper.service.validEnumerationDataManager = validIngredientGroupDataManager

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(true, nil)

		dbManager.ValidIngredientGroupDataManagerMock.On(
			"ArchiveValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(errors.New("blah"))
		helper.service.validEnumerationDataManager = dbManager

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
