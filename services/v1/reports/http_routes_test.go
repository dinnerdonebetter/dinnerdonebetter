package reports

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestReportsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReportList := fakemodels.BuildFakeReportList()

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReports", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleReportList, nil)
		s.reportDataManager = reportDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ReportList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReports", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ReportList)(nil), sql.ErrNoRows)
		s.reportDataManager = reportDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ReportList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, ed)
	})

	T.Run("with error fetching reports from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReports", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ReportList)(nil), errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReportList := fakemodels.BuildFakeReportList()

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReports", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleReportList, nil)
		s.reportDataManager = reportDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ReportList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, ed)
	})
}

func TestReportsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("CreateReport", mock.Anything, mock.AnythingOfType("*models.ReportCreationInput")).Return(exampleReport, nil)
		s.reportDataManager = reportDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.reportCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating report", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("CreateReport", mock.Anything, mock.AnythingOfType("*models.ReportCreationInput")).Return(exampleReport, errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportCreationInputFromReport(exampleReport)

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("CreateReport", mock.Anything, mock.AnythingOfType("*models.ReportCreationInput")).Return(exampleReport, nil)
		s.reportDataManager = reportDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.reportCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, mc, r, ed)
	})
}

func TestReportsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ReportExists", mock.Anything, exampleReport.ID).Return(true, nil)
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with no such report in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ReportExists", mock.Anything, exampleReport.ID).Return(false, sql.ErrNoRows)
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error fetching report from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ReportExists", mock.Anything, exampleReport.ID).Return(false, errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})
}

func TestReportsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)
		s.reportDataManager = reportDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, ed)
	})

	T.Run("with no such report in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return((*models.Report)(nil), sql.ErrNoRows)
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error fetching report from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return((*models.Report)(nil), errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)
		s.reportDataManager = reportDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, ed)
	})
}

func TestReportsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)
		reportDataManager.On("UpdateReport", mock.Anything, mock.AnythingOfType("*models.Report")).Return(nil)
		s.reportDataManager = reportDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, reportDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching report", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return((*models.Report)(nil), sql.ErrNoRows)
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error fetching report", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return((*models.Report)(nil), errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error updating report", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)
		reportDataManager.On("UpdateReport", mock.Anything, mock.AnythingOfType("*models.Report")).Return(errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeReportUpdateInputFromReport(exampleReport)

		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("GetReport", mock.Anything, exampleReport.ID).Return(exampleReport, nil)
		reportDataManager.On("UpdateReport", mock.Anything, mock.AnythingOfType("*models.Report")).Return(nil)
		s.reportDataManager = reportDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Report")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, reportDataManager, ed)
	})
}

func TestReportsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ArchiveReport", mock.Anything, exampleReport.ID, exampleUser.ID).Return(nil)
		s.reportDataManager = reportDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.reportCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager, mc, r)
	})

	T.Run("with no report in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ArchiveReport", mock.Anything, exampleReport.ID, exampleUser.ID).Return(sql.ErrNoRows)
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleReport := fakemodels.BuildFakeReport()
		exampleReport.BelongsToUser = exampleUser.ID
		s.reportIDFetcher = func(req *http.Request) uint64 {
			return exampleReport.ID
		}

		reportDataManager := &mockmodels.ReportDataManager{}
		reportDataManager.On("ArchiveReport", mock.Anything, exampleReport.ID, exampleUser.ID).Return(errors.New("blah"))
		s.reportDataManager = reportDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, reportDataManager)
	})
}
