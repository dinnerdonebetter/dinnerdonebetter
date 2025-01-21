package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type contextKey string

const (
	cookieName         = "dinner-done-better-admin-webapp"
	oauth2ClientID     = "9819637062b9bbd3c1997cd3b6a264d4"
	oauth2ClientSecret = "0299fececf3f0be3af94adc9a98b2b0b"

	apiClientContextKey contextKey = "api_client"
)

var (
	apiServerURL *url.URL
)

func init() {
	var err error
	apiServerURL, err = url.Parse("https://api.dinnerdonebetter.dev/")
	if err != nil {
		panic(err)
	}
}

func mustValidateTextProps(props *components.TextInputsProps) components.ValidatedTextInput {
	if props == nil {
		panic("props cannot be nil")
	}

	validatedProps, err := components.BuildValidatedTextInputPrompt(context.Background(), props)
	if err != nil {
		panic(err)
	}

	return validatedProps
}

// begin things that need to be moved or simplified elsewhere

var (
	validatedUsernameInputProps = mustValidateTextProps(&components.TextInputsProps{
		ID:          "username",
		Name:        "username",
		LabelText:   "Username",
		Type:        "text",
		Placeholder: "username",
	})

	validatedPasswordInputProps = mustValidateTextProps(&components.TextInputsProps{
		ID:          "password",
		Name:        "password",
		LabelText:   "Password",
		Type:        "password",
		Placeholder: "hunter2",
	})

	validatedTOTPCodeInputProps = mustValidateTextProps(&components.TextInputsProps{
		ID:          "totpToken",
		Name:        "totpToken",
		LabelText:   "TOTP Token",
		Type:        "text",
		Placeholder: "123456",
	})
)

type userSessionDetails struct {
	Token       string `json:"token"`
	UserID      string `json:"userID"`
	HouseholdID string `json:"householdID"`
}

// end things that need to be moved or simplified elsewhere

func setupRoutes(logger logging.Logger, tracer tracing.Tracer, router routing.Router, pageBuilder *pages.PageBuilder, cookieBuilder CookieManager) error {
	if pageBuilder == nil {
		return internalerrors.NilConfigError("pageBuilder for frontend admin service")
	}

	router = router.WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx, span := tracer.StartSpan(req.Context())
			defer span.End()

			cookie, err := req.Cookie(cookieName)
			if err != nil {
				logger.Error("fetching request cookie", err)
				next.ServeHTTP(res, req)
				return
			} else if cookie == nil {
				logger.Info("cookie was nil!")
				next.ServeHTTP(res, req)
				return
			}

			var usd *userSessionDetails
			if err = cookieBuilder.Decode(cookieName, cookie.Value, &usd); err != nil {
				logger.Error("decoding cookie", err)
				next.ServeHTTP(res, req)
				return
			}

			logger.WithValue("user.id", usd.UserID).Info("user session retrieved from middleware")

			client, err := apiclient.NewClient(
				apiServerURL,
				tracing.NewNoopTracerProvider(),
				apiclient.UsingOAuth2(
					ctx,
					oauth2ClientID,
					oauth2ClientSecret,
					[]string{authorization.HouseholdAdminRoleName},
					usd.Token,
				),
			)
			if err != nil {
				logger.Error("establishing API client", err)
				next.ServeHTTP(res, req)
				return
			}

			req = req.WithContext(context.WithValue(ctx, apiClientContextKey, client))

			next.ServeHTTP(res, req)
		})
	})

	router.Get("/", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := tracer.StartSpan(req.Context())
		defer span.End()

		return pageBuilder.HomePage(ctx), nil
	}))

	router.Get("/about", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := tracer.StartSpan(req.Context())
		defer span.End()

		return pageBuilder.AboutPage(ctx), nil
	}))

	router.Post("/login/submit", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		_, span := tracer.StartSpan(req.Context())
		defer span.End()

		response, err := pageBuilder.AdminLoginSubmit(req)
		if err != nil {
			return nil, err
		}

		if response == nil {
			res.Header().Set("HX-Redirect", "/login")
			return ghtml.Div(
				ghtml.H1(gomponents.Text("bad")),
			), nil
		}

		usd := &userSessionDetails{
			Token:       response.Token,
			UserID:      response.UserID,
			HouseholdID: response.HouseholdID,
		}

		encoded, err := cookieBuilder.Encode(cookieName, usd)
		if err != nil {
			return nil, err
		}

		res.Header().Set("HX-Redirect", "/")
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Domain:   req.URL.Host,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)

		// obligatory div return
		return ghtml.Div(), nil
	}))

	router.Get("/login", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := tracer.StartSpan(req.Context())
		defer span.End()

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
						ghtml.Div(components.FormTextInput(ctx, validatedUsernameInputProps)),
						ghtml.Div(components.FormTextInput(ctx, validatedPasswordInputProps)),
						ghtml.Div(components.FormTextInput(ctx, validatedTOTPCodeInputProps)),
					),
					ghtml.Div(
						ghtml.ID("result"),
						ghtml.Class("mt-4 text-sm text-gray-700"),
					),
				),
			),
		), nil
	}))

	router.Get("/valid_ingredients", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := tracer.StartSpan(req.Context())
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
	}))

	router.Get("/valid_ingredients/new", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := tracer.StartSpan(req.Context())
		defer span.End()

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
			TargetID:    "TODO_formBody",
		}, fieldNameMap)
		if err != nil {
			logger.Error("generating inputs", err)
			return nil, err
		}

		return components.PageShell(
			"New Valid Ingredient",
			elements...,
		), nil
	}))

	return nil
}
