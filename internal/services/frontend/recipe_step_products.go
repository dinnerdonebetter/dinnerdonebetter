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
	recipeStepProductIDURLParamKey = "recipe_step_product"
)

func (s *service) fetchRecipeStepProduct(ctx context.Context, req *http.Request) (recipeStepProduct *types.RecipeStepProduct, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeStepProduct = fakes.BuildFakeRecipeStepProduct()
	} else {
		// determine recipe ID.
		recipeID := s.recipeIDFetcher(req)
		tracing.AttachRecipeIDToSpan(span, recipeID)
		logger = logger.WithValue(keys.RecipeIDKey, recipeID)

		// determine recipe step ID.
		recipeStepID := s.recipeStepIDFetcher(req)
		tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
		logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

		// determine recipe step product ID.
		recipeStepProductID := s.recipeStepProductIDFetcher(req)
		tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)
		logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)

		recipeStepProduct, err = s.dataStore.GetRecipeStepProduct(ctx, recipeID, recipeStepID, recipeStepProductID)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step product data")
		}
	}

	return recipeStepProduct, nil
}

//go:embed templates/partials/generated/creators/recipe_step_product_creator.gotpl
var recipeStepProductCreatorTemplate string

func (s *service) buildRecipeStepProductCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepProduct := &types.RecipeStepProduct{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepProductCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Recipe Step Product",
				ContentData: recipeStepProduct,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepProductCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepProduct, res)
		}
	}
}

const (
	recipeStepProductNameFormKey          = "name"
	recipeStepProductQuantityTypeFormKey  = "quantityType"
	recipeStepProductQuantityValueFormKey = "quantityValue"
	recipeStepProductQuantityNotesFormKey = "quantityNotes"
	recipeStepProductRecipeStepIDFormKey  = "recipeStepID"

	recipeStepProductCreationInputNameFormKey          = recipeStepProductNameFormKey
	recipeStepProductCreationInputQuantityTypeFormKey  = recipeStepProductQuantityTypeFormKey
	recipeStepProductCreationInputQuantityValueFormKey = recipeStepProductQuantityValueFormKey
	recipeStepProductCreationInputQuantityNotesFormKey = recipeStepProductQuantityNotesFormKey
	recipeStepProductCreationInputRecipeStepIDFormKey  = recipeStepProductRecipeStepIDFormKey

	recipeStepProductUpdateInputNameFormKey          = recipeStepProductNameFormKey
	recipeStepProductUpdateInputQuantityTypeFormKey  = recipeStepProductQuantityTypeFormKey
	recipeStepProductUpdateInputQuantityValueFormKey = recipeStepProductQuantityValueFormKey
	recipeStepProductUpdateInputQuantityNotesFormKey = recipeStepProductQuantityNotesFormKey
	recipeStepProductUpdateInputRecipeStepIDFormKey  = recipeStepProductRecipeStepIDFormKey
)

// parseFormEncodedRecipeStepProductCreationInput checks a request for an RecipeStepProductCreationInput.
func (s *service) parseFormEncodedRecipeStepProductCreationInput(ctx context.Context, req *http.Request) (creationInput *types.RecipeStepProductCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step product creation input")
		return nil
	}

	creationInput = &types.RecipeStepProductCreationInput{
		Name:          form.Get(recipeStepProductCreationInputNameFormKey),
		QuantityType:  form.Get(recipeStepProductCreationInputQuantityTypeFormKey),
		QuantityValue: s.stringToFloat32(form, recipeStepProductCreationInputQuantityValueFormKey),
		QuantityNotes: form.Get(recipeStepProductCreationInputQuantityNotesFormKey),
		RecipeStepID:  s.stringToUint64(form, recipeStepProductCreationInputRecipeStepIDFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step product creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleRecipeStepProductCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("recipe step product creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for recipe step product creation route")

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	creationInput := s.parseFormEncodedRecipeStepProductCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step product creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	creationInput.BelongsToRecipeStep = recipeStepID

	logger.Debug("recipe step product creation input parsed successfully")

	if _, err = s.dataStore.CreateRecipeStepProduct(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing recipe step product to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("recipe step product created")

	htmxRedirectTo(res, "/recipe_step_products")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/recipe_step_product_editor.gotpl
var recipeStepProductEditorTemplate string

func (s *service) buildRecipeStepProductEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepProduct, err := s.fetchRecipeStepProduct(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe step product from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.RecipeStepProduct) string {
				return fmt.Sprintf("RecipeStepProduct #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(recipeStepProductEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("RecipeStepProduct #%d", recipeStepProduct.ID),
				ContentData: recipeStepProduct,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", recipeStepProductEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepProduct, res)
		}
	}
}

func (s *service) fetchRecipeStepProducts(ctx context.Context, req *http.Request) (recipeStepProducts *types.RecipeStepProductList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	if s.useFakeData {
		recipeStepProducts = fakes.BuildFakeRecipeStepProductList()
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

		recipeStepProducts, err = s.dataStore.GetRecipeStepProducts(ctx, recipeID, recipeStepID, filter)
		if err != nil {
			return nil, observability.PrepareError(err, logger, span, "fetching recipe step product data")
		}
	}

	return recipeStepProducts, nil
}

//go:embed templates/partials/generated/tables/recipe_step_products_table.gotpl
var recipeStepProductsTableTemplate string

func (s *service) buildRecipeStepProductsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		recipeStepProducts, err := s.fetchRecipeStepProducts(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching recipe step products from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.RecipeStepProduct) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/recipe_step_products/%d", x.ID))
			},
			"pushURL": func(x *types.RecipeStepProduct) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/recipe_step_products/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(recipeStepProductsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Recipe Step Products",
				ContentData: recipeStepProducts,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", recipeStepProductsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, recipeStepProducts, res)
		}
	}
}

// parseFormEncodedRecipeStepProductUpdateInput checks a request for an RecipeStepProductUpdateInput.
func (s *service) parseFormEncodedRecipeStepProductUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.RecipeStepProductUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing recipe step product creation input")
		return nil
	}

	updateInput = &types.RecipeStepProductUpdateInput{
		Name:          form.Get(recipeStepProductUpdateInputNameFormKey),
		QuantityType:  form.Get(recipeStepProductUpdateInputQuantityTypeFormKey),
		QuantityValue: s.stringToFloat32(form, recipeStepProductUpdateInputQuantityValueFormKey),
		QuantityNotes: form.Get(recipeStepProductUpdateInputQuantityNotesFormKey),
		RecipeStepID:  s.stringToUint64(form, recipeStepProductUpdateInputRecipeStepIDFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid recipe step product creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleRecipeStepProductUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedRecipeStepProductUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	recipeStepProduct, err := s.fetchRecipeStepProduct(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step product from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := recipeStepProduct.Update(updateInput)

	if err = s.dataStore.UpdateRecipeStepProduct(ctx, recipeStepProduct, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step product from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.RecipeStepProduct) string {
			return fmt.Sprintf("RecipeStepProduct #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", recipeStepProductEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeStepProduct, res)
}

func (s *service) handleRecipeStepProductArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	recipeStepProductID := s.recipeStepProductIDFetcher(req)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)

	if err = s.dataStore.ArchiveRecipeStepProduct(ctx, recipeStepID, recipeStepProductID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe step products in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipeStepProducts, err := s.fetchRecipeStepProducts(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipe step products from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.RecipeStepProduct) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/recipe_step_products/%d", x.ID))
		},
		"pushURL": func(x *types.RecipeStepProduct) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/recipe_step_products/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", recipeStepProductsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, recipeStepProducts, res)
}
