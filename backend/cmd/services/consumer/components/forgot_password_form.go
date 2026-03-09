package components

import (
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/design"
	webappdesign "github.com/dinnerdonebetter/backend/internal/webapp/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// ForgotPasswordFormProps holds the forgot password form state.
type ForgotPasswordFormProps struct {
	EmailError   string
	GeneralError string
	Success      bool
}

// ForgotPasswordForm renders a forgot password form with email input.
func (r *ComponentRenderer) ForgotPasswordForm(props *ForgotPasswordFormProps) g.Node {
	if props == nil {
		props = &ForgotPasswordFormProps{}
	}
	palette := design.StandardPalette

	if props.Success {
		return r.forgotPasswordSuccessCard(&palette)
	}

	return ghtml.Div(
		ghtml.ID("forgot-password-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text("Forgot password"),
			),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-sm mb-4 %s", design.TextColor(palette.Text))),
				g.Text("Enter your email address and we'll send you a link to reset your password."),
			),
			ghtml.Form(
				ghtml.Class("space-y-4"),
				ghtml.Method("post"),
				g.Attr("hx-post", "/forgot_password/submit"),
				g.Attr("hx-ext", "json-enc"),
				g.Attr("hx-target", "#forgot-password-container"),
				g.Attr("hx-swap", "outerHTML"),
				g.Attr("hx-request", `{"credentials":"include"}`),

				wrapInputElement("email", props.EmailError, emailInput("email", "emailAddress", "", &palette), &palette),

				submitButton("Send reset link"),

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

func (r *ComponentRenderer) forgotPasswordSuccessCard(palette *webappdesign.Palette) g.Node {
	return ghtml.Div(
		ghtml.ID("forgot-password-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text("Check your email"),
			),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-sm mb-4 %s", design.TextColor(palette.Text))),
				g.Text("If an account exists with that email address, we've sent you a link to reset your password."),
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

func emailInput(id, fieldName, content string, palette *webappdesign.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("email"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("email"),
	)
}

func wrapInputElement(label, inputError string, input g.Node, palette *webappdesign.Palette) g.Node {
	titleLabel := strings.ToUpper(label[:1]) + strings.ToLower(label[1:])
	if strings.Contains(label, "-") {
		titleLabel = strings.ReplaceAll(titleLabel, "-", " ")
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

func inputClass(palette *webappdesign.Palette) string {
	return fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s",
		palette.Background.Value, palette.Primary.Value, palette.Primary.Value)
}

func submitButton(text string) g.Node {
	return ghtml.Button(
		ghtml.Type("submit"),
		ghtml.Class("w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md shadow hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"),
		g.Text(text),
	)
}
