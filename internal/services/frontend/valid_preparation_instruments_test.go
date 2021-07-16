package frontend

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_fetchValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(exampleValidPreparationInstrument, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstrument(s.ctx, req)
		assert.Equal(t, exampleValidPreparationInstrument, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstrument(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstrument(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidPreparationInstrumentCreationInputToRequest(input *types.ValidPreparationInstrumentCreationInput) *http.Request {
	form := url.Values{
		validPreparationInstrumentCreationInputInstrumentIDFormKey:  {anyToString(input.InstrumentID)},
		validPreparationInstrumentCreationInputPreparationIDFormKey: {anyToString(input.PreparationID)},
		validPreparationInstrumentCreationInputNotesFormKey:         {anyToString(input.Notes)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_preparation_instruments", strings.NewReader(form.Encode()))
}

func TestService_buildValidPreparationInstrumentCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedValidPreparationInstrumentCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeValidPreparationInstrumentCreationInput()
		req := attachValidPreparationInstrumentCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedValidPreparationInstrumentCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidPreparationInstrumentCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationInstrumentCreationInput{}
		req := attachValidPreparationInstrumentCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidPreparationInstrumentCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidPreparationInstrumentCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidPreparationInstrument, nil)
		s.service.dataStore = mockDB

		s.service.handleValidPreparationInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentCreationInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentCreationInputToRequest(&types.ValidPreparationInstrumentCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidPreparationInstrument, nil)
		s.service.dataStore = mockDB

		s.service.handleValidPreparationInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating valid preparation instrument in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"CreateValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleValidPreparationInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidPreparationInstrumentEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(exampleValidPreparationInstrument, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(exampleValidPreparationInstrument, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentEditorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstruments(s.ctx, req)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstruments(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		actual, err := s.service.fetchValidPreparationInstruments(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidPreparationInstrumentsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentsTableView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparation_instruments", nil)

		s.service.buildValidPreparationInstrumentsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidPreparationInstrumentUpdateInputToRequest(input *types.ValidPreparationInstrumentUpdateInput) *http.Request {
	form := url.Values{
		validPreparationInstrumentUpdateInputInstrumentIDFormKey:  {anyToString(input.InstrumentID)},
		validPreparationInstrumentUpdateInputPreparationIDFormKey: {anyToString(input.PreparationID)},
		validPreparationInstrumentUpdateInputNotesFormKey:         {anyToString(input.Notes)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_preparation_instruments", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedValidPreparationInstrumentUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		expected := fakes.BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		req := attachValidPreparationInstrumentUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedValidPreparationInstrumentUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidPreparationInstrumentUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationInstrumentUpdateInput{}

		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidPreparationInstrumentUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidPreparationInstrumentUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(exampleValidPreparationInstrument, nil)

		mockDB.ValidPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationInstrumentUpdateInput{}

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return((*types.ValidPreparationInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(exampleValidPreparationInstrument, nil)

		mockDB.ValidPreparationInstrumentDataManager.On(
			"UpdateValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleValidPreparationInstrumentArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparation_instruments", nil)

		s.service.handleValidPreparationInstrumentArchiveRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparation_instruments", nil)

		s.service.handleValidPreparationInstrumentArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving valid preparation instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparation_instruments", nil)

		s.service.handleValidPreparationInstrumentArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of valid preparation instruments", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		s.service.validPreparationInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparationInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationInstrumentDataManager.On(
			"ArchiveValidPreparationInstrument",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidPreparationInstrumentDataManager.On(
			"GetValidPreparationInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparation_instruments", nil)

		s.service.handleValidPreparationInstrumentArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
