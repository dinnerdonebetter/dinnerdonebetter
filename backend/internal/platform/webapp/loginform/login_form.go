package loginform

import (
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/webapp/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// Props holds the login form state and error messages.
type Props struct {
	UsernameError,
	PasswordError,
	TOTPError,
	GeneralError string
	// ResetSuccessMessage is shown when the user has successfully reset their password (e.g. "Your password has been reset. Sign in with your new password.").
	ResetSuccessMessage string
}

// Config allows customizing the form copy per app.
type Config struct {
	// Title is the form heading (e.g. "Login" or "Sign In").
	Title string
	// SubmitButtonText is the submit button label (e.g. "Log In" or "Sign In").
	SubmitButtonText string
	// ForgotPasswordLink is the href for the "Forgot password?" link. If empty, the link is not shown.
	ForgotPasswordLink string
}

// DefaultConfig returns config with "Login" / "Log In" copy.
func DefaultConfig() Config {
	return Config{
		Title:            "Login",
		SubmitButtonText: "Log In",
	}
}

// SignInConfig returns config with "Sign In" copy.
func SignInConfig() Config {
	return Config{
		Title:              "Sign In",
		SubmitButtonText:   "Sign In",
		ForgotPasswordLink: "/forgot_password",
	}
}

// Form renders a login form with username, password, and TOTP fields.
// Uses HTMX for submission. Pass empty strings in props for no errors.
func Form(props *Props, cfg Config, palette *design.Palette) g.Node {
	if props == nil {
		props = &Props{}
	}
	if palette == nil {
		palette = &design.StandardPalette
	}
	if cfg.Title == "" {
		cfg.Title = "Login"
	}
	if cfg.SubmitButtonText == "" {
		cfg.SubmitButtonText = "Log In"
	}

	return ghtml.Div(
		ghtml.ID("login-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text(cfg.Title),
			),

			g.If(props.ResetSuccessMessage != "", ghtml.Div(
				ghtml.Class(fmt.Sprintf("mb-4 p-3 rounded-md text-sm bg-green-50 border border-green-200 %s", design.TextColor(palette.Primary))),
				g.Text(props.ResetSuccessMessage),
			)),

			ghtml.Form(
				ghtml.Class("space-y-4"),
				ghtml.Method("post"),
				g.Attr("hx-post", "/login/submit"),
				g.Attr("hx-ext", "json-enc"),
				g.Attr("hx-target", "#login-container"),
				g.Attr("hx-swap", "outerHTML"),
				g.Attr("hx-request", `{"credentials":"include"}`),

				wrapInputElement("username", props.UsernameError, usernameInput("username", "username", "", palette), palette),
				wrapInputElement("password", props.PasswordError, passwordInput("password", "password", "", palette), palette),
				g.If(cfg.ForgotPasswordLink != "", ghtml.Div(
					ghtml.Class("text-right -mt-2"),
					ghtml.A(
						ghtml.Href(cfg.ForgotPasswordLink),
						ghtml.Class(fmt.Sprintf("text-sm %s hover:underline", design.TextColor(palette.Primary))),
						g.Text("Forgot password?"),
					),
				)),
				wrapInputElement("TOTP code", props.TOTPError, totpTokenInput("totp", "totpToken", "", palette), palette),

				submitButton(cfg.SubmitButtonText),

				g.If(props.GeneralError != "", ghtml.Div(
					ghtml.Class(fmt.Sprintf("mt-2 text-sm %s", design.TextColor(palette.Warning))),
					g.Text(props.GeneralError),
				)),
			),
		),
	)
}

func usernameInput(label, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("text"),
		ghtml.ID(label),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("username"),
	)
}

func passwordInput(id, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("password"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("current-password"),
	)
}

func totpTokenInput(id, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("text"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.MaxLength("6"),
		g.Attr("inputmode", "numeric"),
		ghtml.Pattern("[0-9]{6}"),
		ghtml.AutoComplete("one-time-code"),
	)
}

func inputClass(palette *design.Palette) string {
	return fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s",
		palette.Background.Value, palette.Primary.Value, palette.Primary.Value)
}

func wrapInputElement(label, inputError string, input g.Node, palette *design.Palette) g.Node {
	titleLabel := strings.ToUpper(label[:1]) + strings.ToLower(label[1:])
	if label == "TOTP code" {
		titleLabel = "TOTP code"
	}

	return ghtml.Div(
		ghtml.Class("space-y-1"),
		ghtml.Label(
			ghtml.For(label),
			ghtml.Class(fmt.Sprintf("block text-sm font-medium %s", design.TextColor(palette.Primary))),
			g.Text(titleLabel),
		),
		input,
		g.If(inputError != "", ghtml.Span(
			ghtml.Class(fmt.Sprintf("text-sm %s mt-1 block", design.TextColor(palette.Warning))),
			g.Text(inputError),
		)),
	)
}

func submitButton(text string) g.Node {
	return ghtml.Button(
		ghtml.Type("submit"),
		ghtml.Class("w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md shadow hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"),
		g.Text(text),
	)
}
