package recipemanagement

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	testutils "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecipeStepVesselsService_CreateRecipeStepVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"CreateRecipeStepVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepVesselDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeStepVessel, nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepVessel)
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

		helper.service.CreateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepVesselCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"CreateRecipeStepVessel",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepVesselDatabaseCreationInput) bool { return true }),
		).Return((*types.RecipeStepVessel)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.CreateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipeStepVesselsService_ReadRecipeStepVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(helper.exampleRecipeStepVessel, nil)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ReadRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepVessel)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step vessel in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return((*types.RecipeStepVessel)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ReadRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return((*types.RecipeStepVessel)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ReadRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})
}

func TestRecipeStepVesselsService_ListRecipeStepVesselsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeStepVesselList := fakes.BuildFakeRecipeStepVesselsList()

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessels",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepVesselList, nil)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ListRecipeStepVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleRecipeStepVesselList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListRecipeStepVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessels",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.RecipeStepVessel])(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ListRecipeStepVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error retrieving recipe step vessels from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessels",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.RecipeStepVessel])(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ListRecipeStepVesselsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})
}

func TestRecipeStepVesselsService_UpdateRecipeStepVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(helper.exampleRecipeStepVessel, nil)
		helper.service.recipeManagementDataManager = dbManager

		dbManager.RecipeStepVesselDataManagerMock.On(
			"UpdateRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipeStepVessel,
		).Return(nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleRecipeStepVessel)
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepVesselUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
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

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step vessel", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return((*types.RecipeStepVessel)(nil), sql.ErrNoRows)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error retrieving recipe step vessel from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return((*types.RecipeStepVessel)(nil), errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepVesselUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"GetRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(helper.exampleRecipeStepVessel, nil)
		helper.service.recipeManagementDataManager = dbManager

		dbManager.RecipeStepVesselDataManagerMock.On(
			"UpdateRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipeStepVessel,
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.UpdateRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}

func TestRecipeStepVesselsService_ArchiveRecipeStepVesselHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"RecipeStepVesselExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(true, nil)

		dbManager.RecipeStepVesselDataManagerMock.On(
			"ArchiveRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(nil)
		helper.service.recipeManagementDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"PublishAsync",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		assert.Eventually(t, func() bool { return mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher) }, time.Second, time.Millisecond*100)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such recipe step vessel in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"RecipeStepVesselExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(false, nil)
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ArchiveRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepVesselDataManager := NewRecipeManagementDataManagerMock()
		recipeStepVesselDataManager.RecipeStepVesselDataManagerMock.On(
			"RecipeStepVesselExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeManagementDataManager = recipeStepVesselDataManager

		helper.service.ArchiveRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, recipeStepVesselDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepVesselDataManagerMock.On(
			"RecipeStepVesselExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(true, nil)

		dbManager.RecipeStepVesselDataManagerMock.On(
			"ArchiveRecipeStepVessel",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepVessel.ID,
		).Return(errors.New("blah"))
		helper.service.recipeManagementDataManager = dbManager

		helper.service.ArchiveRecipeStepVesselHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.RecipeStepVessel]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dbManager)
	})
}
