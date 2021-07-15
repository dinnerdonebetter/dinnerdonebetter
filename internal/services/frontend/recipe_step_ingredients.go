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
	recipeStepIngredientIDURLParamKey = "recipe_step_ingredient"
)

func (s *service) fetchRecipeStepIngredient(ctx context.Context, req *http.Request) (recipeStepIngredient *types.RecipeStepIngredient, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeStepIngredient = fakes.BuildFakeRecipeStepIngredient()
	} else {
		// determine recipe ID.
		recipeID := s.recipeIDFetcher(req)
		tracing.AttachRecipeIDToSpan(span, recipeID)
		logger = logger.WithValue(keys.RecipeIDKey, recipeID)

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

		// determine recipe step ingredient ID.
		recipeStepIngredientID := s.recipeStepIngredientIDFetcher(req)
		tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)
		logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

		recipeStepIngredient, err = s.dataStore.GetRecipeStepIngredient(ctx, recipeID, recipeStepID, recipeStepIngredientID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step ingredient data")
		}
	}

	return recipeStepIngredient, nil
}

//go:embed templates/partials/generated/creators/recipe_step_ingredient_creator.gotpl
var recipeStepIngredientCreatorTemplate string

func (s *service) buildRecipeStepIngredientCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepIngredient := &types.RecipeStepIngredient{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepIngredientCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Recipe Step Ingredient",
				ContentData: recipeStepIngredient,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepIngredientCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepIngredient, res)
		}
	}
}

const (
	recipeStepIngredientIngredientIDFormKey        = "ingredientID"
	recipeStepIngredientNameFormKey                = "name"
	recipeStepIngredientQuantityTypeFormKey        = "quantityType"
	recipeStepIngredientQuantityValueFormKey       = "quantityValue"
	recipeStepIngredientQuantityNotesFormKey       = "quantityNotes"
	recipeStepIngredientProductOfRecipeStepFormKey = "productOfRecipeStep"
	recipeStepIngredientIngredientNotesFormKey     = "ingredientNotes"

	recipeStepIngredientCreationInputIngredientIDFormKey        = recipeStepIngredientIngredientIDFormKey
	recipeStepIngredientCreationInputNameFormKey                = recipeStepIngredientNameFormKey
	recipeStepIngredientCreationInputQuantityTypeFormKey        = recipeStepIngredientQuantityTypeFormKey
	recipeStepIngredientCreationInputQuantityValueFormKey       = recipeStepIngredientQuantityValueFormKey
	recipeStepIngredientCreationInputQuantityNotesFormKey       = recipeStepIngredientQuantityNotesFormKey
	recipeStepIngredientCreationInputProductOfRecipeStepFormKey = recipeStepIngredientProductOfRecipeStepFormKey
	recipeStepIngredientCreationInputIngredientNotesFormKey     = recipeStepIngredientIngredientNotesFormKey

	recipeStepIngredientUpdateInputIngredientIDFormKey        = recipeStepIngredientIngredientIDFormKey
	recipeStepIngredientUpdateInputNameFormKey                = recipeStepIngredientNameFormKey
	recipeStepIngredientUpdateInputQuantityTypeFormKey        = recipeStepIngredientQuantityTypeFormKey
	recipeStepIngredientUpdateInputQuantityValueFormKey       = recipeStepIngredientQuantityValueFormKey
	recipeStepIngredientUpdateInputQuantityNotesFormKey       = recipeStepIngredientQuantityNotesFormKey
	recipeStepIngredientUpdateInputProductOfRecipeStepFormKey = recipeStepIngredientProductOfRecipeStepFormKey
	recipeStepIngredientUpdateInputIngredientNotesFormKey     = recipeStepIngredientIngredientNotesFormKey
)

// parseFormEncodedRecipeStepIngredientCreationInput checks a request for an RecipeStepIngredientCreationInput.
func (s *service) parseFormEncodedRecipeStepIngredientCreationInput(ctx context.Context, req *http.Request) (creationInput *types.RecipeStepIngredientCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step ingredient creation input")
		return nil
	}

	creationInput = &types.RecipeStepIngredientCreationInput{
		IngredientID:        s.stringToPointerToUint64(form, recipeStepIngredientCreationInputIngredientIDFormKey),
		Name:                form.Get(recipeStepIngredientCreationInputNameFormKey),
		QuantityType:        form.Get(recipeStepIngredientCreationInputQuantityTypeFormKey),
		QuantityValue:       s.stringToFloat32(form, recipeStepIngredientCreationInputQuantityValueFormKey),
		QuantityNotes:       form.Get(recipeStepIngredientCreationInputQuantityNotesFormKey),
		ProductOfRecipeStep: s.stringToBool(form, recipeStepIngredientCreationInputProductOfRecipeStepFormKey),
		IngredientNotes:     form.Get(recipeStepIngredientCreationInputIngredientNotesFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step ingredient creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleRecipeStepIngredientCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("recipe step ingredient creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for recipe step ingredient creation route")

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	creationInput := s.parseFormEncodedRecipeStepIngredientCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step ingredient creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	creationInput.BelongsToRecipeStep = recipeStepID

	logger.Debug("recipe step ingredient creation input parsed successfully")

	if _, err = s.dataStore.CreateRecipeStepIngredient(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing recipe step ingredient to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("recipe step ingredient created")

	htmxRedirectTo(res, "/recipe_step_ingredients")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/recipe_step_ingredient_editor.gotpl
var recipeStepIngredientEditorTemplate string

func (s *service) buildRecipeStepIngredientEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepIngredient, err := s.fetchRecipeStepIngredient(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe step ingredient from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.RecipeStepIngredient) string {
				return fmt.Sprintf("RecipeStepIngredient #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepIngredientEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Recipe Step Ingredient #%d", recipeStepIngredient.ID),
				ContentData: recipeStepIngredient,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepIngredientEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepIngredient, res)
		}
	}
}

func (s *service) fetchRecipeStepIngredients(ctx context.Context, req *http.Request) (recipeStepIngredients *types.RecipeStepIngredientList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeStepIngredients = fakes.BuildFakeRecipeStepIngredientList()
	} else {
		// determine recipe ID.
		recipeID := s.recipeIDFetcher(req)
		tracing.AttachRecipeIDToSpan(span, recipeID)
		logger = logger.WithValue(keys.RecipeIDKey, recipeID)

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

		filter := types.ExtractQueryFilter(req)
		tracing.AttachQueryFilterToSpan(span, filter)

		recipeStepIngredients, err = s.dataStore.GetRecipeStepIngredients(ctx, recipeID, recipeStepID, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step ingredient data")
		}
	}

	return recipeStepIngredients, nil
}

//go:embed templates/partials/generated/tables/recipe_step_ingredients_table.gotpl
var recipeStepIngredientsTableTemplate string

func (s *service) buildRecipeStepIngredientsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepIngredients, err := s.fetchRecipeStepIngredients(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe step ingredients from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.RecipeStepIngredient) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/recipe_step_ingredients/%d", x.ID))
			},
			"pushURL": func(x *types.RecipeStepIngredient) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/recipe_step_ingredients/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(recipeStepIngredientsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Recipe Step Ingredients",
				ContentData: recipeStepIngredients,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", recipeStepIngredientsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepIngredients, res)
		}
	}
}

// parseFormEncodedRecipeStepIngredientUpdateInput checks a request for an RecipeStepIngredientUpdateInput.
func (s *service) parseFormEncodedRecipeStepIngredientUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.RecipeStepIngredientUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step ingredient creation input")
		return nil
	}

	updateInput = &types.RecipeStepIngredientUpdateInput{
		IngredientID:        s.stringToPointerToUint64(form, recipeStepIngredientUpdateInputIngredientIDFormKey),
		Name:                form.Get(recipeStepIngredientUpdateInputNameFormKey),
		QuantityType:        form.Get(recipeStepIngredientUpdateInputQuantityTypeFormKey),
		QuantityValue:       s.stringToFloat32(form, recipeStepIngredientUpdateInputQuantityValueFormKey),
		QuantityNotes:       form.Get(recipeStepIngredientUpdateInputQuantityNotesFormKey),
		ProductOfRecipeStep: s.stringToBool(form, recipeStepIngredientUpdateInputProductOfRecipeStepFormKey),
		IngredientNotes:     form.Get(recipeStepIngredientUpdateInputIngredientNotesFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step ingredient creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleRecipeStepIngredientUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedRecipeStepIngredientUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	recipeStepIngredient, err := s.fetchRecipeStepIngredient(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step ingredient from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := recipeStepIngredient.Update(updateInput)

	if err = s.dataStore.UpdateRecipeStepIngredient(ctx, recipeStepIngredient, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step ingredient from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.RecipeStepIngredient) string {
			return fmt.Sprintf("RecipeStepIngredient #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", recipeStepIngredientEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeStepIngredient, res)
}

func (s *service) handleRecipeStepIngredientArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	recipeStepIngredientID := s.recipeStepIngredientIDFetcher(req)
	tracing.AttachRecipeStepIngredientIDToSpan(span, recipeStepIngredientID)
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	if err = s.dataStore.ArchiveRecipeStepIngredient(ctx, recipeStepID, recipeStepIngredientID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe step ingredients in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipeStepIngredients, err := s.fetchRecipeStepIngredients(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step ingredients from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.RecipeStepIngredient) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/recipe_step_ingredients/%d", x.ID))
		},
		"pushURL": func(x *types.RecipeStepIngredient) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/recipe_step_ingredients/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", recipeStepIngredientsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeStepIngredients, res)
}
