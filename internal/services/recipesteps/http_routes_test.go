package recipesteps

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	mockpublishers "gitlab.com/prixfixe/prixfixe/internal/messagequeue/publishers/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
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

func TestRecipeStepsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreWriteMessageMatcher),
		).Return(nil)
		helper.service.preWritesPublisher = mockEventProducer

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreWriteMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preWritesPublisher = mockEventProducer

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})
}

func TestRecipeStepsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(helper.exampleRecipeStep, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStep{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
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

	T.Run("with no such recipe step in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), sql.ErrNoRows)
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
	})
}

func TestRecipeStepsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleRecipeStepList, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStepList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
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

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepList)(nil), sql.ErrNoRows)
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.RecipeStepList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
	})

	T.Run("with error retrieving recipe steps from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeSteps",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.RecipeStepList)(nil), errors.New("blah"))
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
	})
}

func TestRecipeStepsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(helper.exampleRecipeStep, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreUpdateMessageMatcher),
		).Return(nil)
		helper.service.preUpdatesPublisher = mockEventProducer

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, mockEventProducer)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.RecipeStepUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with no such recipe step", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), sql.ErrNoRows)
		helper.service.recipeStepDataManager = recipeStepDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager)
	})

	T.Run("with error retrieving recipe step from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return((*types.RecipeStep)(nil), errors.New("blah"))
		helper.service.recipeStepDataManager = recipeStepDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeRecipeStepUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"GetRecipeStep",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(helper.exampleRecipeStep, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreUpdateMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preUpdatesPublisher = mockEventProducer

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, mockEventProducer)
	})
}

func TestRecipeStepsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"RecipeStepExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(true, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(nil)
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, mockEventProducer)
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

	T.Run("with no such recipe step in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"RecipeStepExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(false, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"RecipeStepExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(false, errors.New("blah"))
		helper.service.recipeStepDataManager = recipeStepDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		recipeStepDataManager := &mocktypes.RecipeStepDataManager{}
		recipeStepDataManager.On(
			"RecipeStepExists",
			testutils.ContextMatcher,
			helper.exampleRecipe.ID,
			helper.exampleRecipeStep.ID,
		).Return(true, nil)
		helper.service.recipeStepDataManager = recipeStepDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, recipeStepDataManager, mockEventProducer)
	})
}
