package recipestepinstruments

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	expectations := map[string]bool{
		"1":      true,
		t.Name(): false,
		"true":   true,
		"troo":   false,
		"t":      true,
		"false":  false,
	}

	for input, expected := range expectations {
		assert.Equal(t, expected, parseBool(input))
	}
}

func TestRecipeStepInstrumentsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepInstrumentCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepInstrumentDatabaseCreationInput) bool { return true }),
		).Return((*types.RecipeStepInstrument)(nil), errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = dbManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"CreateRecipeStepInstrument",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.RecipeStepInstrumentDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestRecipeStepInstrumentsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStepInstrument{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such recipe step instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return((*types.RecipeStepInstrument)(nil), sql.ErrNoRows)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return((*types.RecipeStepInstrument)(nil), errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})
}

func TestRecipeStepInstrumentsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstruments",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepInstrumentList, nil)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStepInstrumentList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstruments",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepInstrumentList)(nil), sql.ErrNoRows)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStepInstrumentList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving recipe step instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstruments",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepInstrumentList)(nil), errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})
}

func TestRecipeStepInstrumentsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStepInstrument,
		).Return(nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepInstrumentUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without input attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with no such recipe step instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return((*types.RecipeStepInstrument)(nil), sql.ErrNoRows)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error retrieving recipe step instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return((*types.RecipeStepInstrument)(nil), errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStepInstrument,
		).Return(errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = dbManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepInstrumentUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"GetRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(helper.exampleRecipeStepInstrument, nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dbManager.RecipeStepInstrumentDataManager.On(
			"UpdateRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStepInstrument,
		).Return(nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}

func TestRecipeStepInstrumentsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"RecipeStepInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(true, nil)

		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such recipe step instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"RecipeStepInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(false, nil)
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepInstrumentDataManager := &mocktypes.RecipeStepInstrumentDataManager{}
		recipeStepInstrumentDataManager.On(
			"RecipeStepInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = recipeStepInstrumentDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepInstrumentDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"RecipeStepInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(true, nil)

		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(errors.New("blah"))
		helper.service.recipeStepInstrumentDataManager = dbManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.RecipeStepInstrumentDataManager.On(
			"RecipeStepInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(true, nil)

		dbManager.RecipeStepInstrumentDataManager.On(
			"ArchiveRecipeStepInstrument",
			testutils.ContextMatcher,
			helper.exampleRecipeStep.ID,
			helper.exampleRecipeStepInstrument.ID,
		).Return(nil)
		helper.service.recipeStepInstrumentDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
