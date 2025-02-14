package pages

import (
	"encoding/json"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (b *Builder) AdminLoginSubmit(req *http.Request) (*apiclient.TokenResponse, error) {
	ctx, span := b.tracer.StartSpan(req.Context())
	defer span.End()

	logger := b.logger.WithRequest(req)

	var x apiclient.UserLoginInput
	if err := json.NewDecoder(req.Body).Decode(&x); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "decoding json")
	}

	// TODO: validate

	client, err := apiclient.NewClient(b.apiServerURL, b.tracerProvider)
	if err != nil {
		return nil, err
	}

	jwtResponse, err := client.AdminLoginForToken(ctx, &x)
	if err != nil {
		return nil, err
	}

	return jwtResponse, nil
}

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

func (b *Builder) AdminLoginPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
	ctx, span := b.tracer.StartSpan(req.Context())
	defer span.End()

	output := components.PageShell("Registration",
		ghtml.Form(
			components.FormTextInput(ctx, validatedUsernameInputProps),
			components.FormTextInput(ctx, validatedPasswordInputProps),
			components.FormTextInput(ctx, validatedTOTPCodeInputProps),
			components.Button("Login"),
		),
	)

	return output, nil
}
