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
	validIngredientPreparationIDURLParamKey = "valid_ingredient_preparation"
)

func (s *service) fetchValidIngredientPreparation(ctx context.Context, req *http.Request) (validIngredientPreparation *types.ValidIngredientPreparation, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
	} else {
		// determine valid ingredient preparation ID.
		validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
		tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
		logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

		validIngredientPreparation, err = s.dataStore.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient preparation data")
		}
	}

	return validIngredientPreparation, nil
}

//go:embed templates/partials/generated/creators/valid_ingredient_preparation_creator.gotpl
var validIngredientPreparationCreatorTemplate string

func (s *service) buildValidIngredientPreparationCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredientPreparation := &types.ValidIngredientPreparation{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validIngredientPreparationCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Valid Ingredient Preparation",
				ContentData: validIngredientPreparation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validIngredientPreparationCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, validIngredientPreparation, res)
		}
	}
}

const (
	validIngredientPreparationNotesFormKey              = "notes"
	validIngredientPreparationValidIngredientIDFormKey  = "validIngredientID"
	validIngredientPreparationValidPreparationIDFormKey = "validPreparationID"

	validIngredientPreparationCreationInputNotesFormKey              = validIngredientPreparationNotesFormKey
	validIngredientPreparationCreationInputValidIngredientIDFormKey  = validIngredientPreparationValidIngredientIDFormKey
	validIngredientPreparationCreationInputValidPreparationIDFormKey = validIngredientPreparationValidPreparationIDFormKey

	validIngredientPreparationUpdateInputNotesFormKey              = validIngredientPreparationNotesFormKey
	validIngredientPreparationUpdateInputValidIngredientIDFormKey  = validIngredientPreparationValidIngredientIDFormKey
	validIngredientPreparationUpdateInputValidPreparationIDFormKey = validIngredientPreparationValidPreparationIDFormKey
)

// parseFormEncodedValidIngredientPreparationCreationInput checks a request for an ValidIngredientPreparationCreationInput.
func (s *service) parseFormEncodedValidIngredientPreparationCreationInput(ctx context.Context, req *http.Request) (creationInput *types.ValidIngredientPreparationCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient preparation creation input")
		return nil
	}

	creationInput = &types.ValidIngredientPreparationCreationInput{
		Notes:              form.Get(validIngredientPreparationCreationInputNotesFormKey),
		ValidIngredientID:  s.stringToUint64(form, validIngredientPreparationCreationInputValidIngredientIDFormKey),
		ValidPreparationID: s.stringToUint64(form, validIngredientPreparationCreationInputValidPreparationIDFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid ingredient preparation creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleValidIngredientPreparationCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("valid ingredient preparation creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for valid ingredient preparation creation route")

	creationInput := s.parseFormEncodedValidIngredientPreparationCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient preparation creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("valid ingredient preparation creation input parsed successfully")

	if _, err = s.dataStore.CreateValidIngredientPreparation(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing valid ingredient preparation to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("valid ingredient preparation created")

	htmxRedirectTo(res, "/valid_ingredient_preparations")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/valid_ingredient_preparation_editor.gotpl
var validIngredientPreparationEditorTemplate string

func (s *service) buildValidIngredientPreparationEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredientPreparation, err := s.fetchValidIngredientPreparation(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid ingredient preparation from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.ValidIngredientPreparation) string {
				return fmt.Sprintf("ValidIngredientPreparation #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validIngredientPreparationEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Valid Ingredient Preparation #%d", validIngredientPreparation.ID),
				ContentData: validIngredientPreparation,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validIngredientPreparationEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validIngredientPreparation, res)
		}
	}
}

func (s *service) fetchValidIngredientPreparations(ctx context.Context, req *http.Request) (validIngredientPreparations *types.ValidIngredientPreparationList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		validIngredientPreparations = fakes.BuildFakeValidIngredientPreparationList()
	} else {
		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		validIngredientPreparations, err = s.dataStore.GetValidIngredientPreparations(ctx, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient preparation data")
		}
	}

	return validIngredientPreparations, nil
}

//go:embed templates/partials/generated/tables/valid_ingredient_preparations_table.gotpl
var validIngredientPreparationsTableTemplate string

func (s *service) buildValidIngredientPreparationsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredientPreparations, err := s.fetchValidIngredientPreparations(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid ingredient preparations from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.ValidIngredientPreparation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/valid_ingredient_preparations/%d", x.ID))
			},
			"pushURL": func(x *types.ValidIngredientPreparation) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/valid_ingredient_preparations/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(validIngredientPreparationsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Valid Ingredient Preparations",
				ContentData: validIngredientPreparations,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", validIngredientPreparationsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validIngredientPreparations, res)
		}
	}
}

// parseFormEncodedValidIngredientPreparationUpdateInput checks a request for an ValidIngredientPreparationUpdateInput.
func (s *service) parseFormEncodedValidIngredientPreparationUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.ValidIngredientPreparationUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient preparation creation input")
		return nil
	}

	updateInput = &types.ValidIngredientPreparationUpdateInput{
		Notes:              form.Get(validIngredientPreparationUpdateInputNotesFormKey),
		ValidIngredientID:  s.stringToUint64(form, validIngredientPreparationUpdateInputValidIngredientIDFormKey),
		ValidPreparationID: s.stringToUint64(form, validIngredientPreparationUpdateInputValidPreparationIDFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid ingredient preparation creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleValidIngredientPreparationUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedValidIngredientPreparationUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	validIngredientPreparation, err := s.fetchValidIngredientPreparation(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredient preparation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := validIngredientPreparation.Update(updateInput)

	if err = s.dataStore.UpdateValidIngredientPreparation(ctx, validIngredientPreparation, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredient preparation from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.ValidIngredientPreparation) string {
			return fmt.Sprintf("ValidIngredientPreparation #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", validIngredientPreparationEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validIngredientPreparation, res)
}

func (s *service) handleValidIngredientPreparationArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	if err = s.dataStore.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient preparations in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	validIngredientPreparations, err := s.fetchValidIngredientPreparations(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredient preparations from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.ValidIngredientPreparation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/valid_ingredient_preparations/%d", x.ID))
		},
		"pushURL": func(x *types.ValidIngredientPreparation) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/valid_ingredient_preparations/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", validIngredientPreparationsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validIngredientPreparations, res)
}
