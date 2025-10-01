package components

import (
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func loginInput(label, inputError, inputType string, p design.Palette) g.Node {
	input := ghtml.Input(
		ghtml.Type("text"),
		ghtml.ID(label),
		ghtml.Name(label),
		ghtml.Class(fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s", p.Background.Value, p.Primary.Value, p.Primary.Value)),
		ghtml.AutoComplete(label),
	)

	switch strings.ToLower(strings.TrimSpace(inputType)) {
	case "password":
		input = ghtml.Input(
			ghtml.Type("password"),
			ghtml.ID("password"),
			ghtml.Name("password"),
			ghtml.Class(fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s", p.Background.Value, p.Primary.Value, p.Primary.Value)),
			ghtml.AutoComplete("current-password"),
		)
	case "totp":
		input = ghtml.Input(
			ghtml.Type("text"),
			ghtml.ID("totp"),
			ghtml.Name("totp"),
			ghtml.Class(fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s", p.Background.Value, p.Primary.Value, p.Primary.Value)),
			ghtml.MaxLength("6"),
			g.Attr("inputmode", "numeric"),
			ghtml.Pattern("[0-9]{6}"),
		)
	}

	return ghtml.Div(
		ghtml.Class("space-y-1"),
		ghtml.Label(
			ghtml.For(label),
			ghtml.Class(fmt.Sprintf("block text-sm font-medium %s", design.TextColor(p.Primary))),
			g.Text(strings.Title(label)),
		),
		input,
		g.If(inputError != "", ghtml.Span(
			ghtml.Class(fmt.Sprintf("text-sm %s mt-1 block", design.TextColor(p.Warning))),
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

// LoginForm renders a login form. You can optionally provide error messages for each field and a general error.
// Pass empty strings to ignore them.
func LoginForm(
	usernameErr string,
	passwordErr string,
	totpErr string,
	generalErr string,
) g.Node {
	p := design.StandardPalette

	return ghtml.Div(
		ghtml.ID("login-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(p.Primary))),
				g.Text("Login"),
			),

			ghtml.Form(
				ghtml.Class("space-y-4"),
				ghtml.Method("post"),
				g.Attr("hx-post", "/login/submit"),
				g.Attr("hx-target", "#login-container"),
				g.Attr("hx-swap", "outerHTML"),

				loginInput("username", usernameErr, "default", p),
				loginInput("password", passwordErr, "password", p),
				loginInput("TOTP code", totpErr, "totp", p),

				submitButton("Log In"),

				g.If(generalErr != "", ghtml.Div(
					ghtml.Class(fmt.Sprintf("mt-2 text-sm %s", design.TextColor(p.Warning))),
					g.Text(generalErr),
				)),
			),
		),
	)
}
