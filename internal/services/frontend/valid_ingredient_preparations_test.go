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

func TestService_fetchValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return(exampleValidIngredientPreparation, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparation(s.ctx, req)
		assert.Equal(t, exampleValidIngredientPreparation, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparation(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparation(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidIngredientPreparationCreationInputToRequest(input *types.ValidIngredientPreparationCreationInput) *http.Request {
	form := url.Values{
		validIngredientPreparationCreationInputNotesFormKey:              {anyToString(input.Notes)},
		validIngredientPreparationCreationInputValidIngredientIDFormKey:  {anyToString(input.ValidIngredientID)},
		validIngredientPreparationCreationInputValidPreparationIDFormKey: {anyToString(input.ValidPreparationID)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_ingredient_preparations", strings.NewReader(form.Encode()))
}

func TestService_buildValidIngredientPreparationCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedValidIngredientPreparationCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeValidIngredientPreparationCreationInput()
		req := attachValidIngredientPreparationCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedValidIngredientPreparationCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidIngredientPreparationCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientPreparationCreationInput{}
		req := attachValidIngredientPreparationCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidIngredientPreparationCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidIngredientPreparationCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidIngredientPreparation, nil)
		s.service.dataStore = mockDB

		s.service.handleValidIngredientPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationCreationInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationCreationInputToRequest(&types.ValidIngredientPreparationCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidIngredientPreparation, nil)
		s.service.dataStore = mockDB

		s.service.handleValidIngredientPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating valid ingredient preparation in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleValidIngredientPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidIngredientPreparationEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return(exampleValidIngredientPreparation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return(exampleValidIngredientPreparation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparations(s.ctx, req)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparations(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		actual, err := s.service.fetchValidIngredientPreparations(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidIngredientPreparationsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_ingredient_preparations", nil)

		s.service.buildValidIngredientPreparationsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidIngredientPreparationUpdateInputToRequest(input *types.ValidIngredientPreparationUpdateInput) *http.Request {
	form := url.Values{
		validIngredientPreparationUpdateInputNotesFormKey:              {anyToString(input.Notes)},
		validIngredientPreparationUpdateInputValidIngredientIDFormKey:  {anyToString(input.ValidIngredientID)},
		validIngredientPreparationUpdateInputValidPreparationIDFormKey: {anyToString(input.ValidPreparationID)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_ingredient_preparations", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedValidIngredientPreparationUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		expected := fakes.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		req := attachValidIngredientPreparationUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedValidIngredientPreparationUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidIngredientPreparationUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientPreparationUpdateInput{}

		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidIngredientPreparationUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidIngredientPreparationUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return(exampleValidIngredientPreparation, nil)

		mockDB.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidIngredientPreparationUpdateInput{}

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
		).Return(exampleValidIngredientPreparation, nil)

		mockDB.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidIngredientPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidIngredientPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleValidIngredientPreparationArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidIngredientPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredient_preparations", nil)

		s.service.handleValidIngredientPreparationArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredient_preparations", nil)

		s.service.handleValidIngredientPreparationArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredient_preparations", nil)

		s.service.handleValidIngredientPreparationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of valid ingredient preparations", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		s.service.validIngredientPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			exampleValidIngredientPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidIngredientPreparationDataManager.On(
			"GetValidIngredientPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidIngredientPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_ingredient_preparations", nil)

		s.service.handleValidIngredientPreparationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
