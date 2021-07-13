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
	validInstrumentIDURLParamKey = "valid_instrument"
)

func (s *service) fetchValidInstrument(ctx context.Context, req *http.Request) (validInstrument *types.ValidInstrument, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validInstrument = fakes.BuildFakeValidInstrument()
	} else {
		// determine valid instrument ID.
		validInstrumentID := s.validInstrumentIDFetcher(req)
		tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
		logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

		validInstrument, err = s.dataStore.GetValidInstrument(ctx, validInstrumentID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid instrument data")
		}
	}

	return validInstrument, nil
}

//go:embed templates/partials/generated/creators/valid_instrument_creator.gotpl
var validInstrumentCreatorTemplate string

func (s *service) buildValidInstrumentCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validInstrument := &types.ValidInstrument{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validInstrumentCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New ValidInstrument",
				ContentData: validInstrument,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validInstrumentCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, validInstrument, res)
		}
	}
}

const (
	validInstrumentNameFormKey        = "name"
	validInstrumentVariantFormKey     = "variant"
	validInstrumentDescriptionFormKey = "description"
	validInstrumentIconPathFormKey    = "iconPath"

	validInstrumentCreationInputNameFormKey        = validInstrumentNameFormKey
	validInstrumentCreationInputVariantFormKey     = validInstrumentVariantFormKey
	validInstrumentCreationInputDescriptionFormKey = validInstrumentDescriptionFormKey
	validInstrumentCreationInputIconPathFormKey    = validInstrumentIconPathFormKey

	validInstrumentUpdateInputNameFormKey        = validInstrumentNameFormKey
	validInstrumentUpdateInputVariantFormKey     = validInstrumentVariantFormKey
	validInstrumentUpdateInputDescriptionFormKey = validInstrumentDescriptionFormKey
	validInstrumentUpdateInputIconPathFormKey    = validInstrumentIconPathFormKey
)

// parseFormEncodedValidInstrumentCreationInput checks a request for an ValidInstrumentCreationInput.
func (s *service) parseFormEncodedValidInstrumentCreationInput(ctx context.Context, req *http.Request) (creationInput *types.ValidInstrumentCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid instrument creation input")
		return nil
	}

	creationInput = &types.ValidInstrumentCreationInput{
		Name:        form.Get(validInstrumentCreationInputNameFormKey),
		Variant:     form.Get(validInstrumentCreationInputVariantFormKey),
		Description: form.Get(validInstrumentCreationInputDescriptionFormKey),
		IconPath:    form.Get(validInstrumentCreationInputIconPathFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid instrument creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleValidInstrumentCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("valid instrument creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for valid instrument creation route")

	creationInput := s.parseFormEncodedValidInstrumentCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid instrument creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("valid instrument creation input parsed successfully")

	if _, err = s.dataStore.CreateValidInstrument(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing valid instrument to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("valid instrument created")

	htmxRedirectTo(res, "/valid_instruments")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/valid_instrument_editor.gotpl
var validInstrumentEditorTemplate string

func (s *service) buildValidInstrumentEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validInstrument, err := s.fetchValidInstrument(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid instrument from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.ValidInstrument) string {
				return fmt.Sprintf("ValidInstrument #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validInstrumentEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("ValidInstrument #%d", validInstrument.ID),
				ContentData: validInstrument,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validInstrumentEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validInstrument, res)
		}
	}
}

func (s *service) fetchValidInstruments(ctx context.Context, req *http.Request) (validInstruments *types.ValidInstrumentList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validInstruments = fakes.BuildFakeValidInstrumentList()
	} else {
		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		validInstruments, err = s.dataStore.GetValidInstruments(ctx, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid instrument data")
		}
	}

	return validInstruments, nil
}

//go:embed templates/partials/generated/tables/valid_instruments_table.gotpl
var validInstrumentsTableTemplate string

func (s *service) buildValidInstrumentsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validInstruments, err := s.fetchValidInstruments(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid instruments from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.ValidInstrument) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/valid_instruments/%d", x.ID))
			},
			"pushURL": func(x *types.ValidInstrument) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/valid_instruments/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(validInstrumentsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "ValidInstruments",
				ContentData: validInstruments,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", validInstrumentsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validInstruments, res)
		}
	}
}

// parseFormEncodedValidInstrumentUpdateInput checks a request for an ValidInstrumentUpdateInput.
func (s *service) parseFormEncodedValidInstrumentUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.ValidInstrumentUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid instrument creation input")
		return nil
	}

	updateInput = &types.ValidInstrumentUpdateInput{
		Name:        form.Get(validInstrumentUpdateInputNameFormKey),
		Variant:     form.Get(validInstrumentUpdateInputVariantFormKey),
		Description: form.Get(validInstrumentUpdateInputDescriptionFormKey),
		IconPath:    form.Get(validInstrumentUpdateInputIconPathFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid instrument creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleValidInstrumentUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedValidInstrumentUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	validInstrument, err := s.fetchValidInstrument(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid instrument from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := validInstrument.Update(updateInput)

	if err = s.dataStore.UpdateValidInstrument(ctx, validInstrument, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid instrument from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.ValidInstrument) string {
			return fmt.Sprintf("ValidInstrument #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", validInstrumentEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validInstrument, res)
}

func (s *service) handleValidInstrumentArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	if err = s.dataStore.ArchiveValidInstrument(ctx, validInstrumentID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid instruments in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	validInstruments, err := s.fetchValidInstruments(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid instruments from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.ValidInstrument) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/valid_instruments/%d", x.ID))
		},
		"pushURL": func(x *types.ValidInstrument) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/valid_instruments/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", validInstrumentsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validInstruments, res)
}
