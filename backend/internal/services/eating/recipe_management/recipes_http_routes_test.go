package recipemanagement

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
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	mocksearch "github.com/dinnerdonebetter/backend/internal/search/text/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecipesService_CreateRecipeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipe, nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipe)
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

		helper.service.CreateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeDatabaseCreationInput) bool { return true }),
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.CreateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipesService_ReadRecipeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ReadRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipe)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ReadRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ReadRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})
}

func TestRecipesService_ListRecipesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeList := fakes.BuildFakeRecipesList()

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ListRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.Recipe])(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ListRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error retrieving recipes from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipes",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.Recipe])(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ListRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})
}

func TestRecipesService_SearchRecipesHandler(T *testing.T) {
	T.Parallel()

	const exampleQuery = "example"
	exampleRecipeList := fakes.BuildFakeRecipesList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{types.QueryKeySearch: []string{exampleQuery}}.Encode()

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"SearchForRecipes",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeList, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.SearchRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("using external service", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.cfg.UseSearchService = true

		exampleLimit := uint8(123)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		expectedIDs := []string{}
		recipeSearchSubsets := make([]*types.RecipeSearchSubset, len(exampleRecipeList.Data))
		for i := range exampleRecipeList.Data {
			expectedIDs = append(expectedIDs, exampleRecipeList.Data[i].ID)
			recipeSearchSubsets[i] = converters.ConvertRecipeToRecipeSearchSubset(exampleRecipeList.Data[i])
		}

		searchIndex := &mocksearch.IndexManager[types.RecipeSearchSubset]{}
		searchIndex.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
		).Return(recipeSearchSubsets, nil)
		helper.service.searchIndex = searchIndex

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipesWithIDs",
			testutils.ContextMatcher,
			expectedIDs,
		).Return(exampleRecipeList.Data, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.SearchRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager, searchIndex)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.SearchRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{types.QueryKeySearch: []string{exampleQuery}}.Encode()

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"SearchForRecipes",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.Recipe])(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.SearchRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{types.QueryKeySearch: []string{exampleQuery}}.Encode()

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"SearchForRecipes",
			testutils.ContextMatcher,
			exampleQuery,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.Recipe])(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.SearchRecipesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})
}

func TestRecipesService_UpdateRecipeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)

		dbManager.RecipeDataManagerMock.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe,
		).Return(nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipe)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
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

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error retrieving recipe from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)

		dbManager.RecipeDataManagerMock.On(
			"UpdateRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe,
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.UpdateRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipesService_ArchiveRecipeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"RecipeExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(true, nil)

		dbManager.RecipeDataManagerMock.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"RecipeExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(false, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ArchiveRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"RecipeExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.ArchiveRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeDataManagerMock.On(
			"RecipeExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(true, nil)

		dbManager.RecipeDataManagerMock.On(
			"ArchiveRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.ArchiveRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipesService_RecipeEstimatedPrepStepsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		mockGrapher := &recipeanalysis.MockRecipeAnalyzer{}
		mockGrapher.On(
			"GenerateMealPlanTasksForRecipe",
			testutils.ContextMatcher,
			"",
			helper.exampleRecipe,
		).Return([]*types.MealPlanTaskDatabaseCreationInput{}, nil)
		helper.service.recipeAnalyzer = mockGrapher

		helper.service.RecipeEstimatedPrepStepsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.MealPlanTaskDatabaseCreationEstimate]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeDataManager, mockGrapher)
	})
}

func TestRecipesService_RecipeMermaidHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		fakeResult := fakes.BuildFakeID()
		mockGrapher := &recipeanalysis.MockRecipeAnalyzer{}
		mockGrapher.On(
			"RenderMermaidDiagramForRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe,
		).Return(fakeResult)
		helper.service.recipeAnalyzer = mockGrapher

		helper.service.RecipeMermaidHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeDataManager, mockGrapher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.RecipeMermaidHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.RecipeMermaidHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.RecipeMermaidHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})
}

func TestRecipesService_CloneRecipeHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)

		recipeDataManager.RecipeDataManagerMock.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			mock.MatchedBy(func(recipe *types.RecipeDatabaseCreationInput) bool {
				return assert.NotEqual(t, helper.exampleRecipe.ID, recipe.ID)
			}),
		).Return(helper.exampleRecipe, nil)
		helper.service.recipeManagementDataManager = recipeDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CloneRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipe)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, recipeDataManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error getting recipe", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.CloneRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})

	T.Run("with error creating cloned recipe", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.recipeIDFetcher = func(_ *http.Request) string {
			return helper.exampleRecipe.ID
		}

		recipeDataManager := NewRecipeManagementDataManagerMock()
		recipeDataManager.RecipeDataManagerMock.On(
			"GetRecipe",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
		).Return(helper.exampleRecipe, nil)

		recipeDataManager.RecipeDataManagerMock.On(
			"CreateRecipe",
			testutils.ContextMatcher,
			mock.MatchedBy(func(recipe *types.RecipeDatabaseCreationInput) bool {
				return assert.NotEqual(t, helper.exampleRecipe.ID, recipe.ID)
			}),
		).Return((*types.Recipe)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeDataManager

		helper.service.CloneRecipeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Recipe]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeDataManager)
	})
}
