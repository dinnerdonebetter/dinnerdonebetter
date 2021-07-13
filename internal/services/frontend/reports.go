package frontend

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

const (
	reportIDURLParamKey = "report"
)

func (s *service) fetchReport(ctx context.Context, req *http.Request) (report *types.Report, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		report = fakes.BuildFakeReport()
	} else {
		// determine report ID.
		reportID := s.reportIDFetcher(req)
		tracing.AttachReportIDToSpan(span, reportID)
		logger = logger.WithValue(keys.ReportIDKey, reportID)

		report, err = s.dataStore.GetReport(ctx, reportID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching report data")
		}
	}

	return report, nil
}

//go:embed templates/partials/generated/creators/report_creator.gotpl
var reportCreatorTemplate string

func (s *service) buildReportCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		report := &types.Report{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(reportCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Report",
				ContentData: report,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", reportCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, report, res)
		}
	}
}

const (
	reportReportTypeFormKey = "reportType"
	reportConcernFormKey    = "concern"

	reportCreationInputReportTypeFormKey = reportReportTypeFormKey
	reportCreationInputConcernFormKey    = reportConcernFormKey

	reportUpdateInputReportTypeFormKey = reportReportTypeFormKey
	reportUpdateInputConcernFormKey    = reportConcernFormKey
)

// parseFormEncodedReportCreationInput checks a request for an ReportCreationInput.
func (s *service) parseFormEncodedReportCreationInput(ctx context.Context, req *http.Request, sessionCtxData *types.SessionContextData) (creationInput *types.ReportCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing report creation input")
		return nil
	}

	creationInput = &types.ReportCreationInput{
		ReportType:       form.Get(reportCreationInputReportTypeFormKey),
		Concern:          form.Get(reportCreationInputConcernFormKey),
		BelongsToAccount: sessionCtxData.ActiveAccountID,
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid report creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleReportCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("report creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for report creation route")

	creationInput := s.parseFormEncodedReportCreationInput(ctx, req, sessionCtxData)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing report creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("report creation input parsed successfully")

	if _, err = s.dataStore.CreateReport(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing report to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("report created")

	htmxRedirectTo(res, "/reports")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/report_editor.gotpl
var reportEditorTemplate string

func (s *service) buildReportEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		report, err := s.fetchReport(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching report from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.Report) string {
				return fmt.Sprintf("Report #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(reportEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Report #%d", report.ID),
				ContentData: report,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", reportEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, report, res)
		}
	}
}

func (s *service) fetchReports(ctx context.Context, req *http.Request) (reports *types.ReportList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		reports = fakes.BuildFakeReportList()
	} else {
		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		reports, err = s.dataStore.GetReports(ctx, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching report data")
		}
	}

	return reports, nil
}

//go:embed templates/partials/generated/tables/reports_table.gotpl
var reportsTableTemplate string

func (s *service) buildReportsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		logger := s.logger.WithRequest(req)
		tracing.AttachRequestToSpan(span, req)

		sessionCtxData, err := s.sessionContextDataFetcher(req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
			http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
			return
		}

		reports, err := s.fetchReports(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching reports from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.Report) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/reports/%d", x.ID))
			},
			"pushURL": func(x *types.Report) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/reports/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(reportsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Reports",
				ContentData: reports,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", reportsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, reports, res)
		}
	}
}

// parseFormEncodedReportUpdateInput checks a request for an ReportUpdateInput.
func (s *service) parseFormEncodedReportUpdateInput(ctx context.Context, req *http.Request, sessionCtxData *types.SessionContextData) (updateInput *types.ReportUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing report creation input")
		return nil
	}

	updateInput = &types.ReportUpdateInput{
		ReportType:       form.Get(reportUpdateInputReportTypeFormKey),
		Concern:          form.Get(reportUpdateInputConcernFormKey),
		BelongsToAccount: sessionCtxData.ActiveAccountID,
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid report creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleReportUpdateRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	updateInput := s.parseFormEncodedReportUpdateInput(ctx, req, sessionCtxData)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	report, err := s.fetchReport(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching report from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := report.Update(updateInput)

	if err = s.dataStore.UpdateReport(ctx, report, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching report from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.Report) string {
			return fmt.Sprintf("Report #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", reportEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, report, res)
}

func (s *service) handleReportArchiveRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	reportID := s.reportIDFetcher(req)
	tracing.AttachReportIDToSpan(span, reportID)
	logger = logger.WithValue(keys.ReportIDKey, reportID)

	if err = s.dataStore.ArchiveReport(ctx, reportID, sessionCtxData.ActiveAccountID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving reports in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	reports, err := s.fetchReports(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching reports from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.Report) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/reports/%d", x.ID))
		},
		"pushURL": func(x *types.Report) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/reports/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", reportsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, reports, res)
}
