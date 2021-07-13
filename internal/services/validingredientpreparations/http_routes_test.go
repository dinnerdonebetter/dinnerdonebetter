package validingredientpreparations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestValidIngredientPreparationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidIngredientPreparationCreationInput{}),
			helper.exampleUser.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Increment", testutils.ContextMatcher).Return()
		helper.service.validIngredientPreparationCounter = unitCounter

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, unitCounter)
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

		exampleCreationInput := &types.ValidIngredientPreparationCreationInput{}
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

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error creating valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidIngredientPreparationCreationInput{}),
			helper.exampleUser.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidIngredientPreparation{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
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

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})
}

func TestValidIngredientPreparationsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(true, nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
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

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no result in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(false, sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error checking database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ValidIngredientPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(false, errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})
}

func TestValidIngredientPreparationsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidIngredientPreparationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
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

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientPreparationList)(nil), sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidIngredientPreparationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error retrieving valid ingredient preparations from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientPreparationList)(nil), errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})
}

func TestValidIngredientPreparationsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)

		validIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidIngredientPreparation{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidIngredientPreparationUpdateInput{}
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

	T.Run("with no such valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error retrieving valid ingredient preparation from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error updating valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidIngredientPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(helper.exampleValidIngredientPreparation, nil)

		validIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidIngredientPreparation{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Decrement", testutils.ContextMatcher).Return()
		helper.service.validIngredientPreparationCounter = unitCounter

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, unitCounter)
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

	T.Run("with no such valid ingredient preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
			helper.exampleUser.ID,
		).Return(sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error saving as archived", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})
}

func TestAccountsService_AuditEntryHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAuditLogEntries := fakes.BuildFakeAuditLogEntryList().Entries

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetAuditLogEntriesForValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return(exampleAuditLogEntries, nil)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.AuditLogEntry{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetAuditLogEntriesForValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return([]*types.AuditLogEntry(nil), sql.ErrNoRows)
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validIngredientPreparationDataManager := &mocktypes.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On(
			"GetAuditLogEntriesForValidIngredientPreparation",
			testutils.ContextMatcher,
			helper.exampleValidIngredientPreparation.ID,
		).Return([]*types.AuditLogEntry(nil), errors.New("blah"))
		helper.service.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, encoderDecoder)
	})
}
