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
	validIngredientIDURLParamKey = "valid_ingredient"
)

func (s *service) fetchValidIngredient(ctx context.Context, req *http.Request) (validIngredient *types.ValidIngredient, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	validIngredient, err = s.dataStore.GetValidIngredient(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient data")
	}

	return validIngredient, nil
}

//go:embed templates/partials/generated/creators/valid_ingredient_creator.gotpl
var validIngredientCreatorTemplate string

func (s *service) buildValidIngredientCreatorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredient := &types.ValidIngredient{}
		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validIngredientCreatorTemplate, nil)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "New Valid Ingredient",
				ContentData: validIngredient,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validIngredientCreatorTemplate, nil)

			s.renderTemplateToResponse(ctx, tmpl, validIngredient, res)
		}
	}
}

const (
	validIngredientNameFormKey              = "name"
	validIngredientVariantFormKey           = "variant"
	validIngredientDescriptionFormKey       = "description"
	validIngredientWarningFormKey           = "warning"
	validIngredientContainsEggFormKey       = "containsEgg"
	validIngredientContainsDairyFormKey     = "containsDairy"
	validIngredientContainsPeanutFormKey    = "containsPeanut"
	validIngredientContainsTreeNutFormKey   = "containsTreeNut"
	validIngredientContainsSoyFormKey       = "containsSoy"
	validIngredientContainsWheatFormKey     = "containsWheat"
	validIngredientContainsShellfishFormKey = "containsShellfish"
	validIngredientContainsSesameFormKey    = "containsSesame"
	validIngredientContainsFishFormKey      = "containsFish"
	validIngredientContainsGlutenFormKey    = "containsGluten"
	validIngredientAnimalFleshFormKey       = "animalFlesh"
	validIngredientAnimalDerivedFormKey     = "animalDerived"
	validIngredientVolumetricFormKey        = "volumetric"
	validIngredientIconPathFormKey          = "iconPath"

	validIngredientCreationInputNameFormKey              = validIngredientNameFormKey
	validIngredientCreationInputVariantFormKey           = validIngredientVariantFormKey
	validIngredientCreationInputDescriptionFormKey       = validIngredientDescriptionFormKey
	validIngredientCreationInputWarningFormKey           = validIngredientWarningFormKey
	validIngredientCreationInputContainsEggFormKey       = validIngredientContainsEggFormKey
	validIngredientCreationInputContainsDairyFormKey     = validIngredientContainsDairyFormKey
	validIngredientCreationInputContainsPeanutFormKey    = validIngredientContainsPeanutFormKey
	validIngredientCreationInputContainsTreeNutFormKey   = validIngredientContainsTreeNutFormKey
	validIngredientCreationInputContainsSoyFormKey       = validIngredientContainsSoyFormKey
	validIngredientCreationInputContainsWheatFormKey     = validIngredientContainsWheatFormKey
	validIngredientCreationInputContainsShellfishFormKey = validIngredientContainsShellfishFormKey
	validIngredientCreationInputContainsSesameFormKey    = validIngredientContainsSesameFormKey
	validIngredientCreationInputContainsFishFormKey      = validIngredientContainsFishFormKey
	validIngredientCreationInputContainsGlutenFormKey    = validIngredientContainsGlutenFormKey
	validIngredientCreationInputAnimalFleshFormKey       = validIngredientAnimalFleshFormKey
	validIngredientCreationInputAnimalDerivedFormKey     = validIngredientAnimalDerivedFormKey
	validIngredientCreationInputVolumetricFormKey        = validIngredientVolumetricFormKey
	validIngredientCreationInputIconPathFormKey          = validIngredientIconPathFormKey

	validIngredientUpdateInputNameFormKey              = validIngredientNameFormKey
	validIngredientUpdateInputVariantFormKey           = validIngredientVariantFormKey
	validIngredientUpdateInputDescriptionFormKey       = validIngredientDescriptionFormKey
	validIngredientUpdateInputWarningFormKey           = validIngredientWarningFormKey
	validIngredientUpdateInputContainsEggFormKey       = validIngredientContainsEggFormKey
	validIngredientUpdateInputContainsDairyFormKey     = validIngredientContainsDairyFormKey
	validIngredientUpdateInputContainsPeanutFormKey    = validIngredientContainsPeanutFormKey
	validIngredientUpdateInputContainsTreeNutFormKey   = validIngredientContainsTreeNutFormKey
	validIngredientUpdateInputContainsSoyFormKey       = validIngredientContainsSoyFormKey
	validIngredientUpdateInputContainsWheatFormKey     = validIngredientContainsWheatFormKey
	validIngredientUpdateInputContainsShellfishFormKey = validIngredientContainsShellfishFormKey
	validIngredientUpdateInputContainsSesameFormKey    = validIngredientContainsSesameFormKey
	validIngredientUpdateInputContainsFishFormKey      = validIngredientContainsFishFormKey
	validIngredientUpdateInputContainsGlutenFormKey    = validIngredientContainsGlutenFormKey
	validIngredientUpdateInputAnimalFleshFormKey       = validIngredientAnimalFleshFormKey
	validIngredientUpdateInputAnimalDerivedFormKey     = validIngredientAnimalDerivedFormKey
	validIngredientUpdateInputVolumetricFormKey        = validIngredientVolumetricFormKey
	validIngredientUpdateInputIconPathFormKey          = validIngredientIconPathFormKey
)

// parseFormEncodedValidIngredientCreationInput checks a request for an ValidIngredientCreationInput.
func (s *service) parseFormEncodedValidIngredientCreationInput(ctx context.Context, req *http.Request) (creationInput *types.ValidIngredientCreationInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient creation input")
		return nil
	}

	creationInput = &types.ValidIngredientCreationInput{
		Name:              form.Get(validIngredientCreationInputNameFormKey),
		Variant:           form.Get(validIngredientCreationInputVariantFormKey),
		Description:       form.Get(validIngredientCreationInputDescriptionFormKey),
		Warning:           form.Get(validIngredientCreationInputWarningFormKey),
		ContainsEgg:       s.stringToBool(form, validIngredientCreationInputContainsEggFormKey),
		ContainsDairy:     s.stringToBool(form, validIngredientCreationInputContainsDairyFormKey),
		ContainsPeanut:    s.stringToBool(form, validIngredientCreationInputContainsPeanutFormKey),
		ContainsTreeNut:   s.stringToBool(form, validIngredientCreationInputContainsTreeNutFormKey),
		ContainsSoy:       s.stringToBool(form, validIngredientCreationInputContainsSoyFormKey),
		ContainsWheat:     s.stringToBool(form, validIngredientCreationInputContainsWheatFormKey),
		ContainsShellfish: s.stringToBool(form, validIngredientCreationInputContainsShellfishFormKey),
		ContainsSesame:    s.stringToBool(form, validIngredientCreationInputContainsSesameFormKey),
		ContainsFish:      s.stringToBool(form, validIngredientCreationInputContainsFishFormKey),
		ContainsGluten:    s.stringToBool(form, validIngredientCreationInputContainsGlutenFormKey),
		AnimalFlesh:       s.stringToBool(form, validIngredientCreationInputAnimalFleshFormKey),
		AnimalDerived:     s.stringToBool(form, validIngredientCreationInputAnimalDerivedFormKey),
		Volumetric:        s.stringToBool(form, validIngredientCreationInputVolumetricFormKey),
		IconPath:          form.Get(validIngredientCreationInputIconPathFormKey),
	}

	if err = creationInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", creationInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid ingredient creation input")
		return nil
	}

	return creationInput
}

func (s *service) handleValidIngredientCreationRequest(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("valid ingredient creation route called")

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "no session context data attached to request")
		http.Redirect(res, req, "/login", unauthorizedRedirectResponseCode)
		return
	}

	logger.Debug("session context data retrieved for valid ingredient creation route")

	creationInput := s.parseFormEncodedValidIngredientCreationInput(ctx, req)
	if creationInput == nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient creation input")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("valid ingredient creation input parsed successfully")

	if _, err = s.dataStore.CreateValidIngredient(ctx, creationInput, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "writing valid ingredient to datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug("valid ingredient created")

	htmxRedirectTo(res, "/valid_ingredients")
	res.WriteHeader(http.StatusCreated)
}

//go:embed templates/partials/generated/editors/valid_ingredient_editor.gotpl
var validIngredientEditorTemplate string

func (s *service) buildValidIngredientEditorView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredient, err := s.fetchValidIngredient(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid ingredient from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"componentTitle": func(x *types.ValidIngredient) string {
				return fmt.Sprintf("ValidIngredient #%d", x.ID)
			},
		}

		if includeBaseTemplate {
			view := s.renderTemplateIntoBaseTemplate(validIngredientEditorTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       fmt.Sprintf("Valid Ingredient #%d", validIngredient.ID),
				ContentData: validIngredient,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, view, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "", validIngredientEditorTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validIngredient, res)
		}
	}
}

func (s *service) fetchValidIngredients(ctx context.Context, req *http.Request) (validIngredients *types.ValidIngredientList, err error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger
	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilter(req)
	tracing.AttachQueryFilterToSpan(span, filter)

	validIngredients, err = s.dataStore.GetValidIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid ingredient data")
	}

	return validIngredients, nil
}

//go:embed templates/partials/generated/tables/valid_ingredients_table.gotpl
var validIngredientsTableTemplate string

func (s *service) buildValidIngredientsTableView(includeBaseTemplate bool) func(http.ResponseWriter, *http.Request) {
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

		validIngredients, err := s.fetchValidIngredients(ctx, req)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "fetching valid ingredients from datastore")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmplFuncMap := map[string]interface{}{
			"individualURL": func(x *types.ValidIngredient) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/dashboard_pages/valid_ingredients/%d", x.ID))
			},
			"pushURL": func(x *types.ValidIngredient) template.URL {
				// #nosec G203
				return template.URL(fmt.Sprintf("/valid_ingredients/%d", x.ID))
			},
		}

		if includeBaseTemplate {
			tmpl := s.renderTemplateIntoBaseTemplate(validIngredientsTableTemplate, tmplFuncMap)

			page := &pageData{
				IsLoggedIn:  sessionCtxData != nil,
				Title:       "Valid Ingredients",
				ContentData: validIngredients,
			}
			if sessionCtxData != nil {
				page.IsServiceAdmin = sessionCtxData.Requester.ServicePermissions.IsServiceAdmin()
			}

			s.renderTemplateToResponse(ctx, tmpl, page, res)
		} else {
			tmpl := s.parseTemplate(ctx, "dashboard", validIngredientsTableTemplate, tmplFuncMap)

			s.renderTemplateToResponse(ctx, tmpl, validIngredients, res)
		}
	}
}

// parseFormEncodedValidIngredientUpdateInput checks a request for an ValidIngredientUpdateInput.
func (s *service) parseFormEncodedValidIngredientUpdateInput(ctx context.Context, req *http.Request) (updateInput *types.ValidIngredientUpdateInput) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	form, err := s.extractFormFromRequest(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "parsing valid ingredient creation input")
		return nil
	}

	updateInput = &types.ValidIngredientUpdateInput{
		Name:              form.Get(validIngredientUpdateInputNameFormKey),
		Variant:           form.Get(validIngredientUpdateInputVariantFormKey),
		Description:       form.Get(validIngredientUpdateInputDescriptionFormKey),
		Warning:           form.Get(validIngredientUpdateInputWarningFormKey),
		ContainsEgg:       s.stringToBool(form, validIngredientUpdateInputContainsEggFormKey),
		ContainsDairy:     s.stringToBool(form, validIngredientUpdateInputContainsDairyFormKey),
		ContainsPeanut:    s.stringToBool(form, validIngredientUpdateInputContainsPeanutFormKey),
		ContainsTreeNut:   s.stringToBool(form, validIngredientUpdateInputContainsTreeNutFormKey),
		ContainsSoy:       s.stringToBool(form, validIngredientUpdateInputContainsSoyFormKey),
		ContainsWheat:     s.stringToBool(form, validIngredientUpdateInputContainsWheatFormKey),
		ContainsShellfish: s.stringToBool(form, validIngredientUpdateInputContainsShellfishFormKey),
		ContainsSesame:    s.stringToBool(form, validIngredientUpdateInputContainsSesameFormKey),
		ContainsFish:      s.stringToBool(form, validIngredientUpdateInputContainsFishFormKey),
		ContainsGluten:    s.stringToBool(form, validIngredientUpdateInputContainsGlutenFormKey),
		AnimalFlesh:       s.stringToBool(form, validIngredientUpdateInputAnimalFleshFormKey),
		AnimalDerived:     s.stringToBool(form, validIngredientUpdateInputAnimalDerivedFormKey),
		Volumetric:        s.stringToBool(form, validIngredientUpdateInputVolumetricFormKey),
		IconPath:          form.Get(validIngredientUpdateInputIconPathFormKey),
	}

	if err = updateInput.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue("input", updateInput)
		observability.AcknowledgeError(err, logger, span, "invalid valid ingredient creation input")
		return nil
	}

	return updateInput
}

func (s *service) handleValidIngredientUpdateRequest(res http.ResponseWriter, req *http.Request) {
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

	updateInput := s.parseFormEncodedValidIngredientUpdateInput(ctx, req)
	if updateInput == nil {
		observability.AcknowledgeError(err, logger, span, "no update input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	validIngredient, err := s.fetchValidIngredient(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredient from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	changes := validIngredient.Update(updateInput)

	if err = s.dataStore.UpdateValidIngredient(ctx, validIngredient, sessionCtxData.Requester.UserID, changes); err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredient from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"componentTitle": func(x *types.ValidIngredient) string {
			return fmt.Sprintf("ValidIngredient #%d", x.ID)
		},
	}

	tmpl := s.parseTemplate(ctx, "", validIngredientEditorTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validIngredient, res)
}

func (s *service) handleValidIngredientArchiveRequest(res http.ResponseWriter, req *http.Request) {
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

	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	if err = s.dataStore.ArchiveValidIngredient(ctx, validIngredientID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredients in datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	validIngredients, err := s.fetchValidIngredients(ctx, req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching valid ingredients from datastore")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmplFuncMap := map[string]interface{}{
		"individualURL": func(x *types.ValidIngredient) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/dashboard_pages/valid_ingredients/%d", x.ID))
		},
		"pushURL": func(x *types.ValidIngredient) template.URL {
			// #nosec G203
			return template.URL(fmt.Sprintf("/valid_ingredients/%d", x.ID))
		},
	}

	tmpl := s.parseTemplate(ctx, "dashboard", validIngredientsTableTemplate, tmplFuncMap)

	s.renderTemplateToResponse(ctx, tmpl, validIngredients, res)
}
