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

func TestService_fetchValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return(exampleValidPreparation, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		actual, err := s.service.fetchValidPreparation(s.ctx, req)
		assert.Equal(t, exampleValidPreparation, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching valid preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		actual, err := s.service.fetchValidPreparation(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidPreparationCreationInputToRequest(input *types.ValidPreparationCreationInput) *http.Request {
	form := url.Values{
		validPreparationCreationInputNameFormKey:        {anyToString(input.Name)},
		validPreparationCreationInputDescriptionFormKey: {anyToString(input.Description)},
		validPreparationCreationInputIconPathFormKey:    {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_preparations", strings.NewReader(form.Encode()))
}

func TestService_buildValidPreparationCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedValidPreparationCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeValidPreparationCreationInput()
		req := attachValidPreparationCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedValidPreparationCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidPreparationCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationCreationInput{}
		req := attachValidPreparationCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidPreparationCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidPreparationCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		res := httptest.NewRecorder()
		req := attachValidPreparationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidPreparation, nil)
		s.service.dataStore = mockDB

		s.service.handleValidPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidPreparationCreationInputToRequest(exampleInput)

		s.service.handleValidPreparationCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		res := httptest.NewRecorder()
		req := attachValidPreparationCreationInputToRequest(&types.ValidPreparationCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidPreparation, nil)
		s.service.dataStore = mockDB

		s.service.handleValidPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating valid preparation in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		res := httptest.NewRecorder()
		req := attachValidPreparationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"CreateValidPreparation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleValidPreparationCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidPreparationEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return(exampleValidPreparation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return(exampleValidPreparation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching valid preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		actual, err := s.service.fetchValidPreparations(s.ctx, req)
		assert.Equal(t, exampleValidPreparationList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		actual, err := s.service.fetchValidPreparations(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidPreparationsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_preparations", nil)

		s.service.buildValidPreparationsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidPreparationUpdateInputToRequest(input *types.ValidPreparationUpdateInput) *http.Request {
	form := url.Values{
		validPreparationUpdateInputNameFormKey:        {anyToString(input.Name)},
		validPreparationUpdateInputDescriptionFormKey: {anyToString(input.Description)},
		validPreparationUpdateInputIconPathFormKey:    {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_preparations", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedValidPreparationUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		expected := fakes.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		req := attachValidPreparationUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedValidPreparationUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidPreparationUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationUpdateInput{}

		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidPreparationUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidPreparationUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return(exampleValidPreparation, nil)

		mockDB.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidPreparationUpdateInput{}

		res := httptest.NewRecorder()
		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return((*types.ValidPreparation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleInput := fakes.BuildFakeValidPreparationUpdateInputFromValidPreparation(exampleValidPreparation)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
		).Return(exampleValidPreparation, nil)

		mockDB.ValidPreparationDataManager.On(
			"UpdateValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidPreparationUpdateInputToRequest(exampleInput)

		s.service.handleValidPreparationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleValidPreparationArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		exampleValidPreparationList := fakes.BuildFakeValidPreparationList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidPreparationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparations", nil)

		s.service.handleValidPreparationArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparations", nil)

		s.service.handleValidPreparationArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving valid preparation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparations", nil)

		s.service.handleValidPreparationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of valid preparations", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		s.service.validPreparationIDFetcher = func(*http.Request) uint64 {
			return exampleValidPreparation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidPreparationDataManager.On(
			"ArchiveValidPreparation",
			testutils.ContextMatcher,
			exampleValidPreparation.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidPreparationDataManager.On(
			"GetValidPreparations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidPreparationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_preparations", nil)

		s.service.handleValidPreparationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
