package recipestepcompletionconditions

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
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

func TestRecipeStepCompletionConditionsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"CreateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepCompletionConditionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeStepCompletionCondition, nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepCompletionCondition)
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

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"CreateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepCompletionConditionDatabaseCreationInput) bool { return true }),
		).Return((*types.RecipeStepCompletionCondition)(nil), errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"CreateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepCompletionConditionDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeStepCompletionCondition, nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepCompletionCondition)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestRecipeStepCompletionConditionsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(helper.exampleRecipeStepCompletionCondition, nil)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepCompletionCondition)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step ingredient in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return((*types.RecipeStepCompletionCondition)(nil), sql.ErrNoRows)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return((*types.RecipeStepCompletionCondition)(nil), errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})
}

func TestRecipeStepCompletionConditionsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeStepCompletionConditionList := fakes.BuildFakeRecipeStepCompletionConditionList()

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionConditions",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepCompletionConditionList, nil)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeStepCompletionConditionList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionConditions",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.RecipeStepCompletionCondition])(nil), sql.ErrNoRows)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error retrieving recipe step ingredients from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionConditions",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.RecipeStepCompletionCondition])(nil), errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})
}

func TestRecipeStepCompletionConditionsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(helper.exampleRecipeStepCompletionCondition, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"UpdateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStepCompletionCondition,
		).Return(nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepCompletionCondition)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepCompletionConditionUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
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

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return((*types.RecipeStepCompletionCondition)(nil), sql.ErrNoRows)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error retrieving recipe step ingredient from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return((*types.RecipeStepCompletionCondition)(nil), errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(helper.exampleRecipeStepCompletionCondition, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"UpdateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStepCompletionCondition,
		).Return(errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepCompletionConditionUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"GetRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(helper.exampleRecipeStepCompletionCondition, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"UpdateRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStepCompletionCondition,
		).Return(nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepCompletionCondition)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestRecipeStepCompletionConditionsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"RecipeStepCompletionConditionExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(true, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"ArchiveRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step ingredient in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"RecipeStepCompletionConditionExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(false, nil)
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepCompletionConditionDataManager := &mocktypes.RecipeStepCompletionConditionDataManagerMock{}
		recipeStepCompletionConditionDataManager.On(
			"RecipeStepCompletionConditionExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = recipeStepCompletionConditionDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepCompletionConditionDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"RecipeStepCompletionConditionExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(true, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"ArchiveRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(errors.New("blah"))
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"RecipeStepCompletionConditionExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(true, nil)

		dbManager.RecipeStepCompletionConditionDataManagerMock.On(
			"ArchiveRecipeStepCompletionCondition",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepCompletionCondition.ID,
		).Return(nil)
		helper.service.recipeStepCompletionConditionDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepCompletionCondition]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
