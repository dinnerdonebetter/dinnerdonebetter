package validpreparationinstruments

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

func TestValidPreparationInstrumentsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationInstrumentCreationInput{}),
			helper.exampleUser.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Increment", testutils.ContextMatcher).Return()
		helper.service.validPreparationInstrumentCounter = unitCounter

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, unitCounter)
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

		exampleCreationInput := &types.ValidPreparationInstrumentCreationInput{}
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

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error creating valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationInstrumentCreationInput{}),
			helper.exampleUser.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationInstrument{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
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

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestValidPreparationInstrumentsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(true, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
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

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error checking database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ValidPreparationInstrumentExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(false, errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestValidPreparationInstrumentsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationInstrumentList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
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

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationInstrumentList)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationInstrumentList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error retrieving valid preparation instruments from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationInstrumentList)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestValidPreparationInstrumentsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		validPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationInstrument{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationInstrumentUpdateInput{}
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

	T.Run("with no such valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error retrieving valid preparation instrument from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})

	T.Run("with error updating valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationInstrumentUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(helper.exampleValidPreparationInstrument, nil)

		validPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationInstrument{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager)
	})
}

func TestValidPreparationInstrumentsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Decrement", testutils.ContextMatcher).Return()
		helper.service.validPreparationInstrumentCounter = unitCounter

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, unitCounter)
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

	T.Run("with no such valid preparation instrument in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
			helper.exampleUser.ID,
		).Return(sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error saving as archived", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}

func TestAccountsService_AuditEntryHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAuditLogEntries := fakes.BuildFakeAuditLogEntryList().Entries

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetAuditLogEntriesForValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return(exampleAuditLogEntries, nil)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
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

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetAuditLogEntriesForValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return([]*types.AuditLogEntry(nil), sql.ErrNoRows)
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationInstrumentDataManager := &mocktypes.ValidPreparationInstrumentDataManager{}
		validPreparationInstrumentDataManager.On(
			"GetAuditLogEntriesForValidPreparationInstrument",
			testutils.ContextMatcher,
			helper.exampleValidPreparationInstrument.ID,
		).Return([]*types.AuditLogEntry(nil), errors.New("blah"))
		helper.service.validPreparationInstrumentDataManager = validPreparationInstrumentDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationInstrumentDataManager, encoderDecoder)
	})
}
