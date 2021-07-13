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
	validPreparationInstrumentIDURLParamKey = "valid_preparation_instrument"
)

func (s *service) fetchValidPreparationInstrument(ctx context.Context, req *http.Request) (validPreparationInstrument *types.ValidPreparationInstrument, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
	} else {
		// determine valid preparation instrument ID.
		validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
		tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
		logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

		validPreparationInstrument, err = s.dataStore.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid preparation instrument data")
		}
	}

	return validPreparationInstrument, nil
}

//go:embed templates/partials/generated/creators/valid_preparation_instrument_creator.gotpl
var validPreparationInstrumentCreatorTemplate string

func (s *service) buildValidPreparationInstrumentCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparationInstrument := &types.ValidPreparationInstrument{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validPreparationInstrumentCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New ValidPreparationInstrument",
				ContentData: validPreparationInstrument,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validPreparationInstrumentCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, validPreparationInstrument, res)
		}
	}
}

const (
	validPreparationInstrumentInstrumentIDFormKey  = "instrumentID"
	validPreparationInstrumentPreparationIDFormKey = "preparationID"
	validPreparationInstrumentNotesFormKey         = "notes"

	validPreparationInstrumentCreationInputInstrumentIDFormKey  = validPreparationInstrumentInstrumentIDFormKey
	validPreparationInstrumentCreationInputPreparationIDFormKey = validPreparationInstrumentPreparationIDFormKey
	validPreparationInstrumentCreationInputNotesFormKey         = validPreparationInstrumentNotesFormKey

	validPreparationInstrumentUpdateInputInstrumentIDFormKey  = validPreparationInstrumentInstrumentIDFormKey
	validPreparationInstrumentUpdateInputPreparationIDFormKey = validPreparationInstrumentPreparationIDFormKey
	validPreparationInstrumentUpdateInputNotesFormKey         = validPreparationInstrumentNotesFormKey
)

// parseFormEncodedValidPreparationInstrumentCreationInput checks a request for an ValidPreparationInstrumentCreationInput.
func (s *service) parseFormEncodedValidPreparationInstrumentCreationInput(ctx context.Context, req *http.Request) (creationInput *types.ValidPreparationInstrumentCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation instrument creation input")
		return nil
	}

	creationInput = &types.ValidPreparationInstrumentCreationInput{
		InstrumentID:  s.stringToUint64(form, validPreparationInstrumentCreationInputInstrumentIDFormKey),
		PreparationID: s.stringToUint64(form, validPreparationInstrumentCreationInputPreparationIDFormKey),
		Notes:         form.Get(validPreparationInstrumentCreationInputNotesFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid preparation instrument creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleValidPreparationInstrumentCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("valid preparation instrument creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for valid preparation instrument creation route")

	creationInput := s.parseFormEncodedValidPreparationInstrumentCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation instrument creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("valid preparation instrument creation input parsed successfully")

	if _, err = s.dataStore.CreateValidPreparationInstrument(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing valid preparation instrument to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("valid preparation instrument created")

	htmxRedirectTo(res, "/valid_preparation_instruments")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/valid_preparation_instrument_editor.gotpl
var validPreparationInstrumentEditorTemplate string

func (s *service) buildValidPreparationInstrumentEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparationInstrument, err := s.fetchValidPreparationInstrument(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid preparation instrument from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.ValidPreparationInstrument) string {
				return fmt.Sprintf("ValidPreparationInstrument #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validPreparationInstrumentEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("ValidPreparationInstrument #%d", validPreparationInstrument.ID),
				ContentData: validPreparationInstrument,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validPreparationInstrumentEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validPreparationInstrument, res)
		}
	}
}

func (s *service) fetchValidPreparationInstruments(ctx context.Context, req *http.Request) (validPreparationInstruments *types.ValidPreparationInstrumentList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validPreparationInstruments = fakes.BuildFakeValidPreparationInstrumentList()
	} else {
		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		validPreparationInstruments, err = s.dataStore.GetValidPreparationInstruments(ctx, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid preparation instrument data")
		}
	}

	return validPreparationInstruments, nil
}

//go:embed templates/partials/generated/tables/valid_preparation_instruments_table.gotpl
var validPreparationInstrumentsTableTemplate string

func (s *service) buildValidPreparationInstrumentsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparationInstruments, err := s.fetchValidPreparationInstruments(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid preparation instruments from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.ValidPreparationInstrument) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/valid_preparation_instruments/%d", x.ID))
			},
			"pushURL": func(x *types.ValidPreparationInstrument) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/valid_preparation_instruments/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(validPreparationInstrumentsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "ValidPreparationInstruments",
				ContentData: validPreparationInstruments,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", validPreparationInstrumentsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validPreparationInstruments, res)
		}
	}
}

// parseFormEncodedValidPreparationInstrumentUpdateInput checks a request for an ValidPreparationInstrumentUpdateInput.
func (s *service) parseFormEncodedValidPreparationInstrumentUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.ValidPreparationInstrumentUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation instrument creation input")
		return nil
	}

	updateInput = &types.ValidPreparationInstrumentUpdateInput{
		InstrumentID:  s.stringToUint64(form, validPreparationInstrumentUpdateInputInstrumentIDFormKey),
		PreparationID: s.stringToUint64(form, validPreparationInstrumentUpdateInputPreparationIDFormKey),
		Notes:         form.Get(validPreparationInstrumentUpdateInputNotesFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid preparation instrument creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleValidPreparationInstrumentUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedValidPreparationInstrumentUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	validPreparationInstrument, err := s.fetchValidPreparationInstrument(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparation instrument from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := validPreparationInstrument.Update(updateInput)

	if err = s.dataStore.UpdateValidPreparationInstrument(ctx, validPreparationInstrument, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparation instrument from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.ValidPreparationInstrument) string {
			return fmt.Sprintf("ValidPreparationInstrument #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", validPreparationInstrumentEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validPreparationInstrument, res)
}

func (s *service) handleValidPreparationInstrumentArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	if err = s.dataStore.ArchiveValidPreparationInstrument(ctx, validPreparationInstrumentID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparation instruments in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	validPreparationInstruments, err := s.fetchValidPreparationInstruments(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparation instruments from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.ValidPreparationInstrument) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/valid_preparation_instruments/%d", x.ID))
		},
		"pushURL": func(x *types.ValidPreparationInstrument) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/valid_preparation_instruments/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", validPreparationInstrumentsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validPreparationInstruments, res)
}
