package components

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/design"
	webappdesign "github.com/dinnerdonebetter/backend/internal/webapp/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// ResetPasswordFormProps holds the reset password form state.
type ResetPasswordFormProps struct {
	Token         string
	PasswordError string
	ConfirmError  string
	GeneralError  string
	InvalidToken  bool
	MissingToken  bool
}

// ResetPasswordForm renders a reset password form.
func (r *ComponentRenderer) ResetPasswordForm(props *ResetPasswordFormProps) g.Node {
	if props == nil {
		props = &ResetPasswordFormProps{}
	}
	palette := design.StandardPalette

	if props.MissingToken || props.InvalidToken {
		return r.resetPasswordErrorCard(&palette, props)
	}

	return ghtml.Div(
		ghtml.ID("reset-password-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text("Reset your password"),
			),
			ghtml.Form(
				ghtml.Class("space-y-4"),
				ghtml.Method("post"),
				g.Attr("hx-post", "/reset_password/submit"),
				g.Attr("hx-ext", "json-enc"),
				g.Attr("hx-target", "#reset-password-container"),
				g.Attr("hx-swap", "outerHTML"),
				g.Attr("hx-request", `{"credentials":"include"}`),

				ghtml.Input(
					ghtml.Type("hidden"),
					ghtml.Name("token"),
					ghtml.Value(props.Token),
				),

				wrapInputElement("new-password", props.PasswordError, passwordInput("new-password", "newPassword", &palette), &palette),
				wrapInputElement("confirm-password", props.ConfirmError, passwordInput("confirm-password", "confirmPassword", &palette), &palette),

				submitButton("Reset password"),

				g.If(props.GeneralError != "", ghtml.Div(
					ghtml.Class(fmt.Sprintf("mt-2 text-sm %s", design.TextColor(palette.Warning))),
					g.Text(props.GeneralError),
				)),
			),
			ghtml.Div(
				ghtml.Class("mt-4 text-center"),
				ghtml.A(
					ghtml.Href("/login"),
					ghtml.Class(fmt.Sprintf("text-sm %s hover:underline", design.TextColor(palette.Primary))),
					g.Text("Back to sign in"),
				),
			),
		),
	)
}

func (r *ComponentRenderer) resetPasswordErrorCard(palette *webappdesign.Palette, props *ResetPasswordFormProps) g.Node {
	message := "Invalid or expired reset link."
	if props.MissingToken {
		message = "Missing reset token. Please use the link from your email."
	}
	return ghtml.Div(
		ghtml.ID("reset-password-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text("Reset link invalid"),
			),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-sm mb-4 %s", design.TextColor(palette.Text))),
				g.Text(message),
			),
			ghtml.Div(
				ghtml.Class("mt-4 text-center"),
				ghtml.A(
					ghtml.Href("/forgot_password"),
					ghtml.Class(fmt.Sprintf("text-sm %s hover:underline", design.TextColor(palette.Primary))),
					g.Text("Request a new reset link"),
				),
			),
			ghtml.Div(
				ghtml.Class("mt-2 text-center"),
				ghtml.A(
					ghtml.Href("/login"),
					ghtml.Class(fmt.Sprintf("text-sm %s hover:underline", design.TextColor(palette.Primary))),
					g.Text("Back to sign in"),
				),
			),
		),
	)
}

func passwordInput(id, fieldName string, palette *webappdesign.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("password"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("new-password"),
	)
}
