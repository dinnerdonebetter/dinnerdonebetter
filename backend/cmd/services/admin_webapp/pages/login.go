package pages

import (
	"encoding/json"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (b *PageBuilder) AdminLoginSubmit(req *http.Request) (*types.TokenResponse, error) {
	ctx, span := b.tracer.StartSpan(req.Context())
	defer span.End()

	logger := b.logger.WithRequest(req)

	var x types.UserLoginInput
	if err := json.NewDecoder(req.Body).Decode(&x); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding json")
		return nil, err
	}

	if err := x.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	client, err := b.buildAPIClient()
	if err != nil {
		return nil, err
	}

	jwtResponse, err := client.AdminLoginForJWT(ctx, &x)
	if err != nil {
		return nil, err
	}

	apiclient.UsingOAuth2(ctx, "clientID", "clientSecret", []string{"*"}, jwtResponse.Token)

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

func (b *PageBuilder) AdminLoginPage(_ http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
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
