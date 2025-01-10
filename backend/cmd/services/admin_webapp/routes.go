package main

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/pkg/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/routing"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

const (
	cookieName = "dinner-done-better-admin-webapp"
)

func mustValidateTextProps(props components.TextInputsProps) components.ValidatedTextInput {
	validatedProps, err := components.BuildValidatedTextInputPrompt(context.Background(), &props)
	if err != nil {
		panic(err)
	}

	return validatedProps
}

var (
	validatedUsernameInputProps = mustValidateTextProps(components.TextInputsProps{
		ID:          "username",
		Name:        "username",
		LabelText:   "Username",
		Type:        "text",
		Placeholder: "username",
	})

	validatedPasswordInputProps = mustValidateTextProps(components.TextInputsProps{
		ID:          "password",
		Name:        "password",
		LabelText:   "Password",
		Type:        "password",
		Placeholder: "hunter2",
	})

	validatedTOTPCodeInputProps = mustValidateTextProps(components.TextInputsProps{
		ID:          "totpToken",
		Name:        "totpToken",
		LabelText:   "TOTP Token",
		Type:        "text",
		Placeholder: "123456",
	})
)

func setupRoutes(router routing.Router, pageBuilder *pages.PageBuilder, cookieBuilder CookieBuilder) error {
	if pageBuilder == nil {
		return internalerrors.NilConfigError("pageBuilder for frontend admin service")
	}

	router.Get("/", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx := req.Context()

		return pageBuilder.HomePage(ctx), nil
	}))

	router.Get("/about", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx := req.Context()

		return pageBuilder.AboutPage(ctx), nil
	}))

	router.Post("/login/submit", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
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

		encoded, err := cookieBuilder.Encode(cookieName, response)
		if err != nil {
			return nil, err
		}

		res.Header().Set("HX-Redirect", "/")
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)

		return nil, nil
	}))

	router.Get("/login", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx := req.Context()

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

	return nil
}
