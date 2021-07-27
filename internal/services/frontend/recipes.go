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
)

const (
	recipeIDURLParamKey = "recipe"
)

func (s *service) fetchRecipe(ctx context.Context, req *http.Request) (recipe *types.Recipe, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	recipe, err = s.dataStore.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe data")
	}

	return recipe, nil
}

//go:embed templates/partials/creators/recipe/recipe.gotpl
var recipeCreatorTemplate string

func (s *service) buildRecipeCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipe := &types.Recipe{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Recipe",
				ContentData: recipe,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, recipe, res)
		}
	}
}

const (
	recipeNameFormKey               = "name"
	recipeSourceFormKey             = "source"
	recipeDescriptionFormKey        = "description"
	recipeInspiredByRecipeIDFormKey = "inspiredByRecipeID"

	recipeUpdateInputNameFormKey               = recipeNameFormKey
	recipeUpdateInputSourceFormKey             = recipeSourceFormKey
	recipeUpdateInputDescriptionFormKey        = recipeDescriptionFormKey
	recipeUpdateInputInspiredByRecipeIDFormKey = recipeInspiredByRecipeIDFormKey
)

//go:embed templates/partials/editors/recipe/recipe.gotpl
var recipeEditorTemplate string

func (s *service) buildRecipeEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipe, err := s.fetchRecipe(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.Recipe) string {
				return fmt.Sprintf("Recipe #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Recipe #%d", recipe.ID),
				ContentData: recipe,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipe, res)
		}
	}
}

func (s *service) fetchRecipes(ctx context.Context, req *http.Request) (recipes *types.RecipeList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachQueryFilterToSpan(span, filter)

	recipes, err = s.dataStore.GetRecipes(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe data")
	}

	return recipes, nil
}

//go:embed templates/partials/generated/tables/recipes_table.gotpl
var recipesTableTemplate string

func (s *service) buildRecipesTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipes, err := s.fetchRecipes(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipes from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.Recipe) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/recipes/%d", x.ID))
			},
			"pushURL": func(x *types.Recipe) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/recipes/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(recipesTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Recipes",
				ContentData: recipes,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", recipesTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipes, res)
		}
	}
}

// parseFormEncodedRecipeUpdateInput checks a request for an RecipeUpdateInput.
func (s *service) parseFormEncodedRecipeUpdateInput(ctx context.Context, req *http.Request, sessionCtxData *types.SessionContextData) (updateInput *types.RecipeUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe creation input")
		return nil
	}

	updateInput = &types.RecipeUpdateInput{
		Name:               form.Get(recipeUpdateInputNameFormKey),
		Source:             form.Get(recipeUpdateInputSourceFormKey),
		Description:        form.Get(recipeUpdateInputDescriptionFormKey),
		InspiredByRecipeID: s.stringToPointerToUint64(form, recipeUpdateInputInspiredByRecipeIDFormKey),
		BelongsToAccount:   sessionCtxData.ActiveAccountID,
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleRecipeUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedRecipeUpdateInput(ctx, req, sessionCtxData)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe, err := s.fetchRecipe(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := recipe.Update(updateInput)

	if err = s.dataStore.UpdateRecipe(ctx, recipe, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.Recipe) string {
			return fmt.Sprintf("Recipe #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", recipeEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipe, res)
}

func (s *service) handleRecipeArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	if err = s.dataStore.ArchiveRecipe(ctx, recipeID, sessionCtxData.ActiveAccountID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipes in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipes, err := s.fetchRecipes(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipes from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.Recipe) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/recipes/%d", x.ID))
		},
		"pushURL": func(x *types.Recipe) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/recipes/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", recipesTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipes, res)
}
