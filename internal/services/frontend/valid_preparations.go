package frontend

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	validPreparationIDURLParamKey = "valid_preparation"
)

func (s *service) fetchValidPreparation(ctx context.Context, req *http.Request) (validPreparation *types.ValidPreparation, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	validPreparation, err = s.dataStore.GetValidPreparation(ctx, validPreparationID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparation data")
	}

	return validPreparation, nil
}

//go:embed templates/partials/generated/creators/valid_preparation_creator.gotpl
var validPreparationCreatorTemplate string

func (s *service) buildValidPreparationCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparation := &types.ValidPreparation{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validPreparationCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Valid Preparation",
				ContentData: validPreparation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validPreparationCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, validPreparation, res)
		}
	}
}

const (
	validPreparationNameFormKey        = "name"
	validPreparationDescriptionFormKey = "description"
	validPreparationIconPathFormKey    = "iconPath"

	validPreparationCreationInputNameFormKey        = validPreparationNameFormKey
	validPreparationCreationInputDescriptionFormKey = validPreparationDescriptionFormKey
	validPreparationCreationInputIconPathFormKey    = validPreparationIconPathFormKey

	validPreparationUpdateInputNameFormKey        = validPreparationNameFormKey
	validPreparationUpdateInputDescriptionFormKey = validPreparationDescriptionFormKey
	validPreparationUpdateInputIconPathFormKey    = validPreparationIconPathFormKey
)

// parseFormEncodedValidPreparationCreationInput checks a request for an ValidPreparationCreationInput.
func (s *service) parseFormEncodedValidPreparationCreationInput(ctx context.Context, req *http.Request) (creationInput *types.ValidPreparationCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation creation input")
		return nil
	}

	creationInput = &types.ValidPreparationCreationInput{
		Name:        form.Get(validPreparationCreationInputNameFormKey),
		Description: form.Get(validPreparationCreationInputDescriptionFormKey),
		IconPath:    form.Get(validPreparationCreationInputIconPathFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid preparation creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleValidPreparationCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("valid preparation creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for valid preparation creation route")

	creationInput := s.parseFormEncodedValidPreparationCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("valid preparation creation input parsed successfully")

	if _, err = s.dataStore.CreateValidPreparation(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing valid preparation to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("valid preparation created")

	htmxRedirectTo(res, "/valid_preparations")
	res.WriteHeader(http.StatusCreated)
}

func (s *service) validPreparationsSearchResults(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	query := req.URL.Query().Get(types.SearchQueryKey)
	tracing.AttachSearchQueryToSpan(span, query)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachQueryFilterToSpan(span, filter)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	renderID := uuid.New().String()

	_, err = s.validPreparationsService.SearchForValidPreparations(ctx, sessionCtxData, query, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "error searching for valid preparations")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("").Parse(`<input class="form-control" type="text" id="preparationSearch{{ .RenderID }}" list="preparationSuggestions{{ .RenderID }}" placeholder="" name="q" hx-target="#preparationSearch{{ .RenderID }}" hx-trigger="keyup changed" value={{ .Query }} hx-get="/elements/valid_preparations/search"/>
<datalist id="preparationSuggestions{{ .RenderID }}">
	{{ for $i, $preparation range .Preparations }}
	
	{{ end }}
</datalist>
	`))
	x := &struct {
		RenderID string
		Query    string
	}{
		RenderID: renderID,
		Query:    query,
	}

	s.renderTemplateToResponse(ctx, tmpl, x, res)
}

//go:embed templates/partials/generated/editors/valid_preparation_editor.gotpl
var validPreparationEditorTemplate string

func (s *service) buildValidPreparationEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparation, err := s.fetchValidPreparation(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid preparation from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.ValidPreparation) string {
				return fmt.Sprintf("ValidPreparation #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validPreparationEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Valid Preparation #%d", validPreparation.ID),
				ContentData: validPreparation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validPreparationEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validPreparation, res)
		}
	}
}

func (s *service) fetchValidPreparations(ctx context.Context, req *http.Request) (validPreparations *types.ValidPreparationList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachQueryFilterToSpan(span, filter)

	validPreparations, err = s.dataStore.GetValidPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid preparation data")
	}

	return validPreparations, nil
}

//go:embed templates/partials/generated/tables/valid_preparations_table.gotpl
var validPreparationsTableTemplate string

func (s *service) buildValidPreparationsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validPreparations, err := s.fetchValidPreparations(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid preparations from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.ValidPreparation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/valid_preparations/%d", x.ID))
			},
			"pushURL": func(x *types.ValidPreparation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/valid_preparations/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(validPreparationsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Valid Preparations",
				ContentData: validPreparations,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", validPreparationsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validPreparations, res)
		}
	}
}

// parseFormEncodedValidPreparationUpdateInput checks a request for an ValidPreparationUpdateInput.
func (s *service) parseFormEncodedValidPreparationUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.ValidPreparationUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid preparation creation input")
		return nil
	}

	updateInput = &types.ValidPreparationUpdateInput{
		Name:        form.Get(validPreparationUpdateInputNameFormKey),
		Description: form.Get(validPreparationUpdateInputDescriptionFormKey),
		IconPath:    form.Get(validPreparationUpdateInputIconPathFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid preparation creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleValidPreparationUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedValidPreparationUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	validPreparation, err := s.fetchValidPreparation(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := validPreparation.Update(updateInput)

	if err = s.dataStore.UpdateValidPreparation(ctx, validPreparation, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.ValidPreparation) string {
			return fmt.Sprintf("ValidPreparation #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", validPreparationEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validPreparation, res)
}

func (s *service) handleValidPreparationArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	if err = s.dataStore.ArchiveValidPreparation(ctx, validPreparationID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparations in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	validPreparations, err := s.fetchValidPreparations(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid preparations from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.ValidPreparation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/valid_preparations/%d", x.ID))
		},
		"pushURL": func(x *types.ValidPreparation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/valid_preparations/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", validPreparationsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validPreparations, res)
}
