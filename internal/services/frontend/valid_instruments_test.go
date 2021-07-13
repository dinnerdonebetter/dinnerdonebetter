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

func TestService_fetchValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(exampleValidInstrument, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstrument(s.ctx, req)
		assert.Equal(t, exampleValidInstrument, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstrument(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching valid instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstrument(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidInstrumentCreationInputToRequest(input *types.ValidInstrumentCreationInput) *http.Request {
	form := url.Values{
		validInstrumentCreationInputNameFormKey:        {anyToString(input.Name)},
		validInstrumentCreationInputVariantFormKey:     {anyToString(input.Variant)},
		validInstrumentCreationInputDescriptionFormKey: {anyToString(input.Description)},
		validInstrumentCreationInputIconPathFormKey:    {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_instruments", strings.NewReader(form.Encode()))
}

func TestService_buildValidInstrumentCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedValidInstrumentCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeValidInstrumentCreationInput()
		req := attachValidInstrumentCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedValidInstrumentCreationInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidInstrumentCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidInstrumentCreationInput{}
		req := attachValidInstrumentCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidInstrumentCreationInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidInstrumentCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		res := httptest.NewRecorder()
		req := attachValidInstrumentCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidInstrument, nil)
		s.service.dataStore = mockDB

		s.service.handleValidInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidInstrumentCreationInputToRequest(exampleInput)

		s.service.handleValidInstrumentCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		res := httptest.NewRecorder()
		req := attachValidInstrumentCreationInputToRequest(&types.ValidInstrumentCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleValidInstrument, nil)
		s.service.dataStore = mockDB

		s.service.handleValidInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating valid instrument in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		res := httptest.NewRecorder()
		req := attachValidInstrumentCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleValidInstrumentCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidInstrumentEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(exampleValidInstrument, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(exampleValidInstrument, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching valid instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidInstrumentList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstruments(s.ctx, req)
		assert.Equal(t, exampleValidInstrumentList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with fake mode", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.useFakeData = true

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstruments(s.ctx, req)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		actual, err := s.service.fetchValidInstruments(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildValidInstrumentsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/valid_instruments", nil)

		s.service.buildValidInstrumentsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachValidInstrumentUpdateInputToRequest(input *types.ValidInstrumentUpdateInput) *http.Request {
	form := url.Values{
		validInstrumentUpdateInputNameFormKey:        {anyToString(input.Name)},
		validInstrumentUpdateInputVariantFormKey:     {anyToString(input.Variant)},
		validInstrumentUpdateInputDescriptionFormKey: {anyToString(input.Description)},
		validInstrumentUpdateInputIconPathFormKey:    {anyToString(input.IconPath)},
	}

	return httptest.NewRequest(http.MethodPost, "/valid_instruments", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedValidInstrumentUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		expected := fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(exampleValidInstrument)

		req := attachValidInstrumentUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedValidInstrumentUpdateInput(s.ctx, req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedValidInstrumentUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidInstrumentUpdateInput{}

		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedValidInstrumentUpdateInput(s.ctx, req)
		assert.Nil(t, actual)
	})
}

func TestService_handleValidInstrumentUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(exampleValidInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(exampleValidInstrument, nil)

		mockDB.ValidInstrumentDataManager.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(exampleValidInstrument)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidInstrumentUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ValidInstrumentUpdateInput{}

		res := httptest.NewRecorder()
		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(exampleValidInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return((*types.ValidInstrument)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleInput := fakes.BuildFakeValidInstrumentUpdateInputFromValidInstrument(exampleValidInstrument)

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(exampleValidInstrument, nil)

		mockDB.ValidInstrumentDataManager.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachValidInstrumentUpdateInputToRequest(exampleInput)

		s.service.handleValidInstrumentUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleValidInstrumentArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleValidInstrumentList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_instruments", nil)

		s.service.handleValidInstrumentArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/valid_instruments", nil)

		s.service.handleValidInstrumentArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving valid instrument", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_instruments", nil)

		s.service.handleValidInstrumentArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of valid instruments", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		s.service.validInstrumentIDFetcher = func(*http.Request) uint64 {
			return exampleValidInstrument.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ValidInstrumentDataManager.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ValidInstrumentDataManager.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ValidInstrumentList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/valid_instruments", nil)

		s.service.handleValidInstrumentArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
