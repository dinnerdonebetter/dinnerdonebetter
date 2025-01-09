package pages

import (
	"encoding/json"
	"net/http"

	"maragu.dev/gomponents"
	gcomponents "maragu.dev/gomponents/components"
	ghtml "maragu.dev/gomponents/html"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func (b *PageBuilder) buildAPIClient() (*apiclient.Client, error) {
	return apiclient.NewClient(b.apiServerURL, b.tracerProvider)
}

func (b *PageBuilder) LoginSubmit(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := b.tracer.StartSpan(req.Context())
	defer span.End()

	var x types.UserLoginInput
	if err := json.NewDecoder(req.Body).Decode(&x); err != nil {
		return nil, err
	}

	if err := x.ValidateWithContext(ctx); err != nil {
		res.Header().Set("HX-Redirect", "/login")
		return nil, err
	}

	client, err := b.buildAPIClient()
	if err != nil {
		return nil, err
	}

	jwtResponse, err := client.LoginForJWT(ctx, &x)
	if err != nil {
		return nil, err
	}
	return ghtml.Div(
		ghtml.H1(gomponents.Text(jwtResponse.Token)),
	), nil
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

func (b *PageBuilder) LoginPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := b.tracer.StartSpan(req.Context())
	defer span.End()

	return components.PageShell("Registration",
		ghtml.Form(
			gcomponents.Classes{
				"flex flex-col": true,
				"gap-4":         true,
				"max-w-md":      true,
				"mt-10":         true,
			},
			ghtml.Div(
				components.FormTextInput(ctx, validatedUsernameInputProps),
			),
			ghtml.Div(
				components.FormTextInput(ctx, validatedPasswordInputProps),
			),
			ghtml.Div(
				components.FormTextInput(ctx, validatedTOTPCodeInputProps),
			),
			ghtml.Div(
				components.Button("Login"),
			),
		),
	), nil
}
