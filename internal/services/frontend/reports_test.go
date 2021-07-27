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

func TestService_fetchReport(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return(exampleReport, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		actual, err := s.service.fetchReport(s.ctx, req)
		assert.Equal(t, exampleReport, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching report", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return((*types.Report)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		actual, err := s.service.fetchReport(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachReportCreationInputToRequest(input *types.ReportCreationInput) *http.Request {
	form := url.Values{
		reportCreationInputReportTypeFormKey: {anyToString(input.ReportType)},
		reportCreationInputConcernFormKey:    {anyToString(input.Concern)},
	}

	return httptest.NewRequest(http.MethodPost, "/reports", strings.NewReader(form.Encode()))
}

func TestService_buildReportCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedReportCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeReportCreationInput()
		expected.BelongsToAccount = s.exampleAccount.ID
		req := attachReportCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedReportCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedReportCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ReportCreationInput{}
		req := attachReportCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedReportCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleReportCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachReportCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"CreateReport",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleReport, nil)
		s.service.dataStore = mockDB

		s.service.handleReportCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachReportCreationInputToRequest(exampleInput)

		s.service.handleReportCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachReportCreationInputToRequest(&types.ReportCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"CreateReport",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleReport, nil)
		s.service.dataStore = mockDB

		s.service.handleReportCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating report in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportCreationInputFromReport(exampleReport)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachReportCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"CreateReport",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.Report)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleReportCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildReportEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return(exampleReport, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return(exampleReport, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportEditorView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching report", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return((*types.Report)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchReports(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReportList := fakes.BuildFakeReportList()

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleReportList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		actual, err := s.service.fetchReports(s.ctx, req)
		assert.Equal(t, exampleReportList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ReportList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		actual, err := s.service.fetchReports(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildReportsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReportList := fakes.BuildFakeReportList()
		for _, report := range exampleReportList.Reports {
			report.BelongsToAccount = s.exampleAccount.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleReportList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReportList := fakes.BuildFakeReportList()

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleReportList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportsTableView(false)(res, req)

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
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ReportList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/reports", nil)

		s.service.buildReportsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachReportUpdateInputToRequest(input *types.ReportUpdateInput) *http.Request {
	form := url.Values{
		reportUpdateInputReportTypeFormKey: {anyToString(input.ReportType)},
		reportUpdateInputConcernFormKey:    {anyToString(input.Concern)},
	}

	return httptest.NewRequest(http.MethodPost, "/reports", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedReportUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		expected := fakes.BuildFakeReportUpdateInputFromReport(exampleReport)

		req := attachReportUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedReportUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedReportUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ReportUpdateInput{}

		req := attachReportUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedReportUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleReportUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportUpdateInputFromReport(exampleReport)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return(exampleReport, nil)

		mockDB.ReportDataManager.On(
			"UpdateReport",
			testutils.ContextMatcher,
			exampleReport,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachReportUpdateInputToRequest(exampleInput)

		s.service.handleReportUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachReportUpdateInputToRequest(exampleInput)

		s.service.handleReportUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.ReportUpdateInput{}

		res := httptest.NewRecorder()
		req := attachReportUpdateInputToRequest(exampleInput)

		s.service.handleReportUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportUpdateInputFromReport(exampleReport)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return((*types.Report)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachReportUpdateInputToRequest(exampleInput)

		s.service.handleReportUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleInput := fakes.BuildFakeReportUpdateInputFromReport(exampleReport)

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"GetReport",
			testutils.ContextMatcher,
			exampleReport.ID,
		).Return(exampleReport, nil)

		mockDB.ReportDataManager.On(
			"UpdateReport",
			testutils.ContextMatcher,
			exampleReport,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachReportUpdateInputToRequest(exampleInput)

		s.service.handleReportUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleReportArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		exampleReportList := fakes.BuildFakeReportList()

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"ArchiveReport",
			testutils.ContextMatcher,
			exampleReport.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleReportList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/reports", nil)

		s.service.handleReportArchiveRequest(res, req)

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
		req := httptest.NewRequest(http.MethodDelete, "/reports", nil)

		s.service.handleReportArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving report", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"ArchiveReport",
			testutils.ContextMatcher,
			exampleReport.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/reports", nil)

		s.service.handleReportArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of reports", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleReport := fakes.BuildFakeReport()
		exampleReport.BelongsToAccount = s.exampleAccount.ID
		s.service.reportIDFetcher = func(*http.Request) uint64 {
			return exampleReport.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.ReportDataManager.On(
			"ArchiveReport",
			testutils.ContextMatcher,
			exampleReport.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.ReportDataManager.On(
			"GetReports",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.ReportList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/reports", nil)

		s.service.handleReportArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
