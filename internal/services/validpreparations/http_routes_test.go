package validpreparations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	mocksearch "gitlab.com/prixfixe/prixfixe/internal/search/mock"
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

func TestValidPreparationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationCreationInput{}),
			helper.exampleUser.ID,
		).Return(helper.exampleValidPreparation, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Increment", testutils.ContextMatcher).Return()
		helper.service.validPreparationCounter = unitCounter

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Index",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleValidPreparation,
		).Return(nil)
		helper.service.search = indexManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, unitCounter, indexManager)
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

		exampleCreationInput := &types.ValidPreparationCreationInput{}
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

		exampleCreationInput := fakes.BuildFakeValidPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error creating valid preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationCreationInput{}),
			helper.exampleUser.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error indexing valid preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparationCreationInput{}),
			helper.exampleUser.ID,
		).Return(helper.exampleValidPreparation, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Increment", testutils.ContextMatcher).Return()
		helper.service.validPreparationCounter = unitCounter

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Index",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleValidPreparation,
		).Return(errors.New("blah"))
		helper.service.search = indexManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, unitCounter, indexManager)
	})
}

func TestValidPreparationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(helper.exampleValidPreparation, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparation{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
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

	T.Run("with no such valid preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})
}

func TestValidPreparationsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ValidPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(true, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
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

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ValidPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(false, sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error checking database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ValidPreparationExists",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(false, errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ExistenceHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})
}

func TestValidPreparationsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationList, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
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

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationList)(nil), sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.ValidPreparationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error retrieving valid preparations from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationList)(nil), errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})
}

func TestValidPreparationsService_SearchHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleValidPreparationList := fakes.BuildFakeValidPreparationList()
	exampleValidPreparationIDs := []uint64{}
	for _, x := range exampleValidPreparationList.ValidPreparations {
		exampleValidPreparationIDs = append(exampleValidPreparationIDs, x.ID)
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.SearchQueryKey: []string{exampleQuery},
			types.LimitQueryKey:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
			helper.exampleHousehold.ID,
		).Return(exampleValidPreparationIDs, nil)
		helper.service.search = indexManager

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparationsWithIDs",
			testutils.ContextMatcher,
			exampleLimit,
			exampleValidPreparationIDs,
		).Return(exampleValidPreparationList.ValidPreparations, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.ValidPreparation{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, indexManager, validPreparationDataManager, encoderDecoder)
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

		helper.service.SearchHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with error conducting search", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{types.SearchQueryKey: []string{exampleQuery}}.Encode()

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
			helper.exampleHousehold.ID,
		).Return([]uint64{}, errors.New("blah"))
		helper.service.search = indexManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, indexManager, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.SearchQueryKey: []string{exampleQuery},
			types.LimitQueryKey:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
			helper.exampleHousehold.ID,
		).Return(exampleValidPreparationIDs, nil)
		helper.service.search = indexManager

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparationsWithIDs",
			testutils.ContextMatcher,
			exampleLimit,
			exampleValidPreparationIDs,
		).Return([]*types.ValidPreparation{}, sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.ValidPreparation{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, indexManager, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.SearchQueryKey: []string{exampleQuery},
			types.LimitQueryKey:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Search",
			testutils.ContextMatcher,
			exampleQuery,
			helper.exampleHousehold.ID,
		).Return(exampleValidPreparationIDs, nil)
		helper.service.search = indexManager

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparationsWithIDs",
			testutils.ContextMatcher,
			exampleLimit,
			exampleValidPreparationIDs,
		).Return([]*types.ValidPreparation{}, errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.SearchHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, indexManager, validPreparationDataManager, encoderDecoder)
	})
}

func TestValidPreparationsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(helper.exampleValidPreparation, nil)

		validPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparation{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Index",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleValidPreparation,
		).Return(nil)
		helper.service.search = indexManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, indexManager)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.ValidPreparationUpdateInput{}
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

	T.Run("with no such valid preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error retrieving valid preparation from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error updating valid preparation", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(helper.exampleValidPreparation, nil)

		validPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparation{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeValidPreparationUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(helper.exampleValidPreparation, nil)

		validPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			mock.IsType(&types.ValidPreparation{}),
			helper.exampleUser.ID,
			mock.IsType([]*types.FieldChangeSummary{}),
		).Return(nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Index",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleValidPreparation,
		).Return(errors.New("blah"))
		helper.service.search = indexManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, indexManager)
	})
}

func TestValidPreparationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Delete",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(nil)
		helper.service.search = indexManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Decrement", testutils.ContextMatcher).Return()
		helper.service.validPreparationCounter = unitCounter

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, unitCounter, indexManager)
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

	T.Run("with no such valid preparation in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleUser.ID,
		).Return(sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error saving as archived", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error removing from search index", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

		indexManager := &mocksearch.IndexManager{}
		indexManager.On(
			"Delete",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(errors.New("blah"))
		helper.service.search = indexManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Decrement", testutils.ContextMatcher).Return()
		helper.service.validPreparationCounter = unitCounter

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, unitCounter, indexManager)
	})
}

func TestHouseholdsService_AuditEntryHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAuditLogEntries := fakes.BuildFakeAuditLogEntryList().Entries

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetAuditLogEntriesForValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return(exampleAuditLogEntries, nil)
		helper.service.validPreparationDataManager = validPreparationDataManager

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

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
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

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetAuditLogEntriesForValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return([]*types.AuditLogEntry(nil), sql.ErrNoRows)
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		validPreparationDataManager := &mocktypes.ValidPreparationDataManager{}
		validPreparationDataManager.On(
			"GetAuditLogEntriesForValidPreparation",
			testutils.ContextMatcher,
			helper.exampleValidPreparation.ID,
		).Return([]*types.AuditLogEntry(nil), errors.New("blah"))
		helper.service.validPreparationDataManager = validPreparationDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.AuditEntryHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, validPreparationDataManager, encoderDecoder)
	})
}
