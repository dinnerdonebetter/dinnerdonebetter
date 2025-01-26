package recipemanagement

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	testutils2 "github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecipeRatingsService_CreateRecipeRatingHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeRatingDataManagerMock.On(
			"CreateRecipeRating",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeRatingDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeRating, nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeRating)
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

		helper.service.CreateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeRatingCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.CreateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeRatingDataManagerMock.On(
			"CreateRecipeRating",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeRatingDatabaseCreationInput) bool { return true }),
		).Return((*types.RecipeRating)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.CreateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipeRatingsService_ReadRecipeRatingHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(helper.exampleRecipeRating, nil)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ReadRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeRating)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ReadRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return((*types.RecipeRating)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ReadRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return((*types.RecipeRating)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ReadRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})
}

func TestRecipeRatingsService_ListRecipeRatingsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeRatingList := fakes.BuildFakeRecipeRatingsList()

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRatingsForRecipe",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return(exampleRecipeRatingList, nil)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ListRecipeRatingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeRatingList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ListRecipeRatingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRatingsForRecipe",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.RecipeRating])(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ListRecipeRatingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error retrieving valid instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRatingsForRecipe",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&filtering.QueryFilter{}),
		).Return((*filtering.QueryFilteredResult[types.RecipeRating])(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ListRecipeRatingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})
}

func TestRecipeRatingsService_UpdateRecipeRatingHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(helper.exampleRecipeRating, nil)

		dbManager.RecipeRatingDataManagerMock.On(
			"UpdateRecipeRating",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeRating) bool { return true }),
		).Return(nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeRating)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeRatingUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
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

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return((*types.RecipeRating)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error retrieving valid instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return((*types.RecipeRating)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeRatingUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeRatingDataManagerMock.On(
			"GetRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(helper.exampleRecipeRating, nil)

		dbManager.RecipeRatingDataManagerMock.On(
			"UpdateRecipeRating",
			testutils2.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeRating) bool { return true }),
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.UpdateRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipeRatingsService_ArchiveRecipeRatingHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"RecipeRatingExists",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(true, nil)

		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"ArchiveRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(nil)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils2.ContextMatcher,
			testutils2.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool {
			return mock.AssertExpectationsForObjects(t, recipeRatingDataManager, dataChangesPublisher)
		}, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ArchiveRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such valid instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"RecipeRatingExists",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(false, nil)
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ArchiveRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"RecipeRatingExists",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ArchiveRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeRatingDataManager := NewRecipeManagementDataManagerMock()
		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"RecipeRatingExists",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(true, nil)

		recipeRatingDataManager.RecipeRatingDataManagerMock.On(
			"ArchiveRecipeRating",
			testutils2.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeRating.ID,
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeRatingDataManager

		helper.service.ArchiveRecipeRatingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeRating]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeRatingDataManager)
	})
}
