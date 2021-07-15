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
	recipeStepIDURLParamKey = "recipe_step"
)

func (s *service) fetchRecipeStep(ctx context.Context, req *http.Request) (recipeStep *types.RecipeStep, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeStep = fakes.BuildFakeRecipeStep()
	} else {
		// determine recipe ID.
		recipeID := s.recipeIDFetcher(req)
		tracing.AttachRecipeIDToSpan(span, recipeID)
		logger = logger.WithValue(keys.RecipeIDKey, recipeID)

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

		recipeStep, err = s.dataStore.GetRecipeStep(ctx, recipeID, recipeStepID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step data")
		}
	}

	return recipeStep, nil
}

//go:embed templates/partials/generated/creators/recipe_step_creator.gotpl
var recipeStepCreatorTemplate string

func (s *service) buildRecipeStepCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStep := &types.RecipeStep{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Recipe Step",
				ContentData: recipeStep,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, recipeStep, res)
		}
	}
}

const (
	recipeStepIndexFormKey                     = "index"
	recipeStepPreparationIDFormKey             = "preparationID"
	recipeStepPrerequisiteStepFormKey          = "prerequisiteStep"
	recipeStepMinEstimatedTimeInSecondsFormKey = "minEstimatedTimeInSeconds"
	recipeStepMaxEstimatedTimeInSecondsFormKey = "maxEstimatedTimeInSeconds"
	recipeStepTemperatureInCelsiusFormKey      = "temperatureInCelsius"
	recipeStepNotesFormKey                     = "notes"
	recipeStepWhyFormKey                       = "why"
	recipeStepRecipeIDFormKey                  = "recipeID"

	recipeStepCreationInputIndexFormKey                     = recipeStepIndexFormKey
	recipeStepCreationInputPreparationIDFormKey             = recipeStepPreparationIDFormKey
	recipeStepCreationInputPrerequisiteStepFormKey          = recipeStepPrerequisiteStepFormKey
	recipeStepCreationInputMinEstimatedTimeInSecondsFormKey = recipeStepMinEstimatedTimeInSecondsFormKey
	recipeStepCreationInputMaxEstimatedTimeInSecondsFormKey = recipeStepMaxEstimatedTimeInSecondsFormKey
	recipeStepCreationInputTemperatureInCelsiusFormKey      = recipeStepTemperatureInCelsiusFormKey
	recipeStepCreationInputNotesFormKey                     = recipeStepNotesFormKey
	recipeStepCreationInputWhyFormKey                       = recipeStepWhyFormKey
	recipeStepCreationInputRecipeIDFormKey                  = recipeStepRecipeIDFormKey

	recipeStepUpdateInputIndexFormKey                     = recipeStepIndexFormKey
	recipeStepUpdateInputPreparationIDFormKey             = recipeStepPreparationIDFormKey
	recipeStepUpdateInputPrerequisiteStepFormKey          = recipeStepPrerequisiteStepFormKey
	recipeStepUpdateInputMinEstimatedTimeInSecondsFormKey = recipeStepMinEstimatedTimeInSecondsFormKey
	recipeStepUpdateInputMaxEstimatedTimeInSecondsFormKey = recipeStepMaxEstimatedTimeInSecondsFormKey
	recipeStepUpdateInputTemperatureInCelsiusFormKey      = recipeStepTemperatureInCelsiusFormKey
	recipeStepUpdateInputNotesFormKey                     = recipeStepNotesFormKey
	recipeStepUpdateInputWhyFormKey                       = recipeStepWhyFormKey
	recipeStepUpdateInputRecipeIDFormKey                  = recipeStepRecipeIDFormKey
)

// parseFormEncodedRecipeStepCreationInput checks a request for an RecipeStepCreationInput.
func (s *service) parseFormEncodedRecipeStepCreationInput(ctx context.Context, req *http.Request) (creationInput *types.RecipeStepCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step creation input")
		return nil
	}

	creationInput = &types.RecipeStepCreationInput{
		Index:                     s.stringToUint(form, recipeStepCreationInputIndexFormKey),
		PreparationID:             s.stringToUint64(form, recipeStepCreationInputPreparationIDFormKey),
		PrerequisiteStep:          s.stringToUint64(form, recipeStepCreationInputPrerequisiteStepFormKey),
		MinEstimatedTimeInSeconds: s.stringToUint32(form, recipeStepCreationInputMinEstimatedTimeInSecondsFormKey),
		MaxEstimatedTimeInSeconds: s.stringToUint32(form, recipeStepCreationInputMaxEstimatedTimeInSecondsFormKey),
		TemperatureInCelsius:      s.stringToPointerToUint16(form, recipeStepCreationInputTemperatureInCelsiusFormKey),
		Notes:                     form.Get(recipeStepCreationInputNotesFormKey),
		Why:                       form.Get(recipeStepCreationInputWhyFormKey),
		RecipeID:                  s.stringToUint64(form, recipeStepCreationInputRecipeIDFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleRecipeStepCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("recipe step creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for recipe step creation route")

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	creationInput := s.parseFormEncodedRecipeStepCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	creationInput.BelongsToRecipe = recipeID

	logger.Debug("recipe step creation input parsed successfully")

	if _, err = s.dataStore.CreateRecipeStep(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing recipe step to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("recipe step created")

	htmxRedirectTo(res, "/recipe_steps")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/recipe_step_editor.gotpl
var recipeStepEditorTemplate string

func (s *service) buildRecipeStepEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStep, err := s.fetchRecipeStep(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe step from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.RecipeStep) string {
				return fmt.Sprintf("RecipeStep #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("RecipeStep #%d", recipeStep.ID),
				ContentData: recipeStep,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeStep, res)
		}
	}
}

func (s *service) fetchRecipeSteps(ctx context.Context, req *http.Request) (recipeSteps *types.RecipeStepList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeSteps = fakes.BuildFakeRecipeStepList()
	} else {
		// determine recipe ID.
		recipeID := s.recipeIDFetcher(req)
		tracing.AttachRecipeIDToSpan(span, recipeID)
		logger = logger.WithValue(keys.RecipeIDKey, recipeID)

		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		recipeSteps, err = s.dataStore.GetRecipeSteps(ctx, recipeID, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step data")
		}
	}

	return recipeSteps, nil
}

//go:embed templates/partials/generated/tables/recipe_steps_table.gotpl
var recipeStepsTableTemplate string

func (s *service) buildRecipeStepsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeSteps, err := s.fetchRecipeSteps(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe steps from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.RecipeStep) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/recipe_steps/%d", x.ID))
			},
			"pushURL": func(x *types.RecipeStep) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/recipe_steps/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(recipeStepsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "RecipeSteps",
				ContentData: recipeSteps,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", recipeStepsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeSteps, res)
		}
	}
}

// parseFormEncodedRecipeStepUpdateInput checks a request for an RecipeStepUpdateInput.
func (s *service) parseFormEncodedRecipeStepUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.RecipeStepUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step creation input")
		return nil
	}

	updateInput = &types.RecipeStepUpdateInput{
		Index:                     s.stringToUint(form, recipeStepUpdateInputIndexFormKey),
		PreparationID:             s.stringToUint64(form, recipeStepUpdateInputPreparationIDFormKey),
		PrerequisiteStep:          s.stringToUint64(form, recipeStepUpdateInputPrerequisiteStepFormKey),
		MinEstimatedTimeInSeconds: s.stringToUint32(form, recipeStepUpdateInputMinEstimatedTimeInSecondsFormKey),
		MaxEstimatedTimeInSeconds: s.stringToUint32(form, recipeStepUpdateInputMaxEstimatedTimeInSecondsFormKey),
		TemperatureInCelsius:      s.stringToPointerToUint16(form, recipeStepUpdateInputTemperatureInCelsiusFormKey),
		Notes:                     form.Get(recipeStepUpdateInputNotesFormKey),
		Why:                       form.Get(recipeStepUpdateInputWhyFormKey),
		RecipeID:                  s.stringToUint64(form, recipeStepUpdateInputRecipeIDFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleRecipeStepUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedRecipeStepUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	recipeStep, err := s.fetchRecipeStep(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := recipeStep.Update(updateInput)

	if err = s.dataStore.UpdateRecipeStep(ctx, recipeStep, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.RecipeStep) string {
			return fmt.Sprintf("RecipeStep #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", recipeStepEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeStep, res)
}

func (s *service) handleRecipeStepArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	if err = s.dataStore.ArchiveRecipeStep(ctx, recipeID, recipeStepID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe steps in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipeSteps, err := s.fetchRecipeSteps(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe steps from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.RecipeStep) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/recipe_steps/%d", x.ID))
		},
		"pushURL": func(x *types.RecipeStep) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/recipe_steps/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", recipeStepsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeSteps, res)
}
