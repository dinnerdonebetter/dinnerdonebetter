package validingredientgroups

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
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

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

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientGroupDataManagerMock.On(
			"CreateValidIngredientGroup",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.ValidIngredientGroupDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleValidIngredientGroup, nil)
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientGroup)
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
		helper.service.validIngredientGroupDataManager = dbManager

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
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
		).Return(helper.exampleValidIngredientGroup, nil)
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleValidIngredientGroup)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestValidIngredientGroupsService_ReadValidIngredientGroupHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(helper.exampleValidIngredientGroup, nil)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), sql.ErrNoRows)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), errors.New("blah"))
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientGroupList, nil)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidIngredientGroup])(nil), sql.ErrNoRows)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

		helper.service.ListValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving valid ingredients from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroups",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ValidIngredientGroup])(nil), errors.New("blah"))
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientGroupList.Data, nil)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return([]*types.ValidIngredientGroup{}, sql.ErrNoRows)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

		helper.service.SearchValidIngredientGroupsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientGroupDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"SearchForValidIngredientGroups",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return([]*types.ValidIngredientGroup(nil), errors.New("blah"))
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), sql.ErrNoRows)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"GetValidIngredientGroup",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return((*types.ValidIngredientGroup)(nil), errors.New("blah"))
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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
		helper.service.validIngredientGroupDataManager = dbManager

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
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
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(false, nil)
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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

		validIngredientGroupDataManager := &mocktypes.ValidIngredientGroupDataManagerMock{}
		validIngredientGroupDataManager.On(
			"ValidIngredientGroupExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientGroup.ID,
		).Return(false, errors.New("blah"))
		helper.service.validIngredientGroupDataManager = validIngredientGroupDataManager

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
		helper.service.validIngredientGroupDataManager = dbManager

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ValidIngredientGroup]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
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
		helper.service.validIngredientGroupDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveValidIngredientGroupHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
