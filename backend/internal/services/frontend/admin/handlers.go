package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/services/frontend/admin/components"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/donseba/go-htmx"
	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	ValidIngredientIDURLParamKey = "valid_ingredient_id"
)

// static pages

func (s *WebappServer) RenderHome(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	return s.pageBuilder.HomePage(ctx), nil
}

func (s *WebappServer) RenderAbout(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	return s.pageBuilder.AboutPage(ctx), nil
}

// auth handlers

func (s *WebappServer) HandleLoginSubmission(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	var x types.UserLoginInput
	if err := json.NewDecoder(req.Body).Decode(&x); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "decoding json")
	}

	if err := x.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	client, err := apiclient.NewClient(s.apiServerURL, s.tracerProvider)
	if err != nil {
		return nil, err
	}

	response, err := client.AdminLoginForToken(ctx, &x)
	if err != nil {
		return nil, err
	}

	h := htmx.New().NewHandler(res, req)

	if response == nil {
		h.Redirect("/login")
		return ghtml.Div(
			ghtml.H1(gomponents.Text("bad")),
		), nil
	}

	usd := &userSessionDetails{
		Token:       response.Token,
		UserID:      response.UserID,
		HouseholdID: response.HouseholdID,
	}

	encoded, err := s.cookieManager.Encode(ctx, s.cookiesConfig.CookieName, usd)
	if err != nil {
		return nil, err
	}

	h.Redirect("/")
	cookie := &http.Cookie{
		Name:     s.cookiesConfig.CookieName,
		Value:    encoded,
		Domain:   req.URL.Host,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	// obligatory div return
	return ghtml.Div(), nil
}

func (s *WebappServer) RenderLoginPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	inputNodes := []gomponents.Node{}
	for _, field := range components.GetFieldNames[types.UserLoginInput]() {
		fieldType := "text"
		if field == "Password" {
			fieldType = "password"
		}

		props := &components.TextInputsProps{
			ID:          field,
			Name:        field,
			LabelText:   field,
			Type:        fieldType,
			Placeholder: field,
		}

		validatedProps, err := components.BuildValidatedTextInputPrompt(ctx, props)
		if err != nil {
			panic(err)
		}

		inputNodes = append(inputNodes, ghtml.Div(components.FormTextInput(ctx, validatedProps)))
	}

	return components.PageShell(
		"Fart",
		ghtml.Div(
			ghtml.Div(
				ghtml.Class("w-full max-w-sm"),
				components.BuildHTMXPoweredSubmissionForm(
					components.SubmissionFormProps{
						PostAddress: "/login/submit",
						TargetID:    "result",
					},
					inputNodes...,
				),
				ghtml.Div(
					ghtml.ID("result"),
					ghtml.Class("mt-4 text-sm text-gray-700"),
				),
			),
		),
	), nil
}

// users handlers

func (s *WebappServer) RenderUsersPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	client, ok := ctx.Value(apiClientContextKey).(*apiclient.Client)
	if !ok {
		return nil, errors.New("missing api client")
	}

	users, err := client.GetUsers(ctx, types.ExtractQueryFilterFromRequest(req))
	if err != nil {
		return nil, err
	}

	return components.PageShell(
		"Users",
		components.TableView("/users/new", users),
	), nil
}

// valid ingredients handlers

func (s *WebappServer) RenderValidIngredientsPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	client, ok := ctx.Value(apiClientContextKey).(*apiclient.Client)
	if !ok {
		return nil, errors.New("missing api client")
	}

	validIngredients, err := client.GetValidIngredients(ctx, types.ExtractQueryFilterFromRequest(req))
	if err != nil {
		return nil, err
	}

	return components.PageShell(
		"Valid Ingredients",
		components.TableView("/valid_ingredients/new", validIngredients),
	), nil
}

func (s *WebappServer) RenderValidIngredientCreationForm(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)

	fieldNameMap := map[string]string{
		"Warning":                "Warning",
		"IconPath":               "Icon Path",
		"PluralName":             "Plural Name",
		"StorageInstructions":    "Storage Instructions",
		"Name":                   "Name",
		"Description":            "Description",
		"Slug":                   "Slug",
		"ShoppingSuggestions":    "Shopping Suggestions",
		"AnimalFlesh":            "Animal Flesh",
		"ContainsEgg":            "Contains Egg",
		"IsLiquid":               "Is Liquid",
		"AnimalDerived":          "Animal Derived",
		"RestrictToPreparations": "Restrict To Preparations",
		"ContainsFish":           "Contains Fish",
		"ContainsShellfish":      "Contains Shellfish",
		"ContainsSoy":            "Contains Soy",
		"ContainsPeanut":         "Contains Peanut",
		"ContainsDairy":          "Contains Dairy",
		"ContainsSesame":         "Contains Sesame",
		"ContainsTreeNut":        "Contains TreeNut",
		"ContainsWheat":          "Contains Wheat",
		"ContainsAlcohol":        "Contains Alcohol",
		"ContainsGluten":         "Contains Gluten",
		"IsStarch":               "Starch?",
		"IsProtein":              "Protein?",
		"IsGrain":                "Grain?",
		"IsFruit":                "Fruit?",
		"IsSalt":                 "Salt?",
		"IsFat":                  "Fat?",
		"IsAcid":                 "Acid?",
		"IsHeat":                 "Heat?",
	}

	elements, err := components.GenerateInputs[types.ValidIngredientCreationRequestInput](ctx, components.SubmissionFormProps{
		PostAddress: "/valid_ingredients/new/submit",
		TargetID:    "main",
	}, fieldNameMap)
	if err != nil {
		logger.Error("generating inputs", err)
		return nil, err
	}

	return components.PageShell(
		"New Valid Ingredient",
		elements...,
	), nil
}

func (s *WebappServer) HandleValidIngredientSubmission(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	client, ok := ctx.Value(apiClientContextKey).(*apiclient.Client)
	if !ok {
		return nil, errors.New("missing api client")
	}

	var x *types.ValidIngredientCreationRequestInput
	if err := json.NewDecoder(req.Body).Decode(&x); err != nil {
		return nil, err
	}

	createdValidIngredient, err := client.CreateValidIngredient(ctx, x)
	if err != nil {
		return nil, err
	}

	h := htmx.New().NewHandler(res, req)
	h.Redirect(fmt.Sprintf("/valid_ingredients/%s", createdValidIngredient.ID))

	// I think what we put in here doesn't matter
	return ghtml.Div(gomponents.Text(createdValidIngredient.ID)), nil
}
