package main

import (
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/design"
	"github.com/dinnerdonebetter/backend/internal/branding"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

var (
	tailwindImport = ghtml.Script(ghtml.Src("https://cdn.tailwindcss.com?plugins=typography"))
	htmxImport     = g.Group{
		ghtml.Script(
			ghtml.Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.min.js"),
			ghtml.Integrity("sha384-ZBXiYtYQ6hJ2Y0ZNoYuI+Nq5MqWBr+chMrS/RkXpNzQCApHEhOt2aY8EJgqwHLkJ"),
			ghtml.CrossOrigin("anonymous"),
		),
		ghtml.Script(
			ghtml.Src("https://unpkg.com/htmx.org@2.0.7/dist/ext/json-enc.js"),
			ghtml.Integrity("sha384-j+tqxLrwDkbeOdjbpaVekgvQL/J7qm/yh/UqSEs6RjEtnBFHqlJViBWG/jBZ6I2p"),
			ghtml.CrossOrigin("anonymous"),
		),
	}
)

func page(title string, children ...g.Node) g.Node {
	palette := &design.StandardPalette

	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Meta(ghtml.Charset("utf-8")),
			ghtml.Meta(ghtml.Name("viewport"), ghtml.Content("width=device-width, initial-scale=1")),
			ghtml.Title(fmt.Sprintf("%s - %s", title, branding.CompanyName)),
			tailwindImport,
			htmxImport,
		),
		ghtml.Body(
			ghtml.Class(fmt.Sprintf("min-h-screen flex flex-col %s %s",
				design.Background(palette.Background),
				design.TextColor(palette.Text),
			)),
			header(),
			ghtml.Main(
				ghtml.Class("flex-1 flex items-center justify-center p-4"),
				ghtml.Div(
					ghtml.Class("w-full max-w-2xl"),
					g.Group(children),
				),
			),
			footer(),
		),
	)
}

// legalPage renders a page for long-form legal content (Terms, Privacy Policy).
// Uses prose typography and top-aligned layout instead of centered.
func legalPage(title string, children ...g.Node) g.Node {
	palette := &design.StandardPalette

	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Meta(ghtml.Charset("utf-8")),
			ghtml.Meta(ghtml.Name("viewport"), ghtml.Content("width=device-width, initial-scale=1")),
			ghtml.Title(fmt.Sprintf("%s - %s", title, branding.CompanyName)),
			tailwindImport,
			htmxImport,
		),
		ghtml.Body(
			ghtml.Class(fmt.Sprintf("min-h-screen flex flex-col %s %s",
				design.Background(palette.Background),
				design.TextColor(palette.Text),
			)),
			header(),
			ghtml.Main(
				ghtml.Class("flex-1 flex items-start justify-center p-4"),
				ghtml.Div(
					ghtml.Class("w-full max-w-3xl prose prose-lg"),
					g.Group(children),
				),
			),
			footer(),
		),
	)
}

func footer() g.Node {
	palette := &design.StandardPalette

	return ghtml.Footer(
		ghtml.Class(fmt.Sprintf("border-t %s %s",
			design.BorderColor(palette.Text),
			design.Background(palette.Background),
		)),
		ghtml.Div(
			ghtml.Class("max-w-4xl mx-auto px-4 py-4"),
			ghtml.Div(
				ghtml.Class("flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0"),
				ghtml.P(
					ghtml.Class(fmt.Sprintf("text-sm %s", design.TextColor(palette.Text))),
					g.Textf("© %d %s. All rights reserved.", time.Now().Year(), branding.CompanyName),
				),
				ghtml.Div(
					ghtml.Class("flex space-x-4 text-sm"),
					ghtml.A(
						ghtml.Href("/privacy-policy"),
						ghtml.Class(fmt.Sprintf("%s hover:%s transition-colors duration-200",
							design.TextColor(palette.Text),
							design.TextColor(palette.Primary),
						)),
						g.Text("Privacy Policy"),
					),
					ghtml.A(
						ghtml.Href("/terms-of-service"),
						ghtml.Class(fmt.Sprintf("%s hover:%s transition-colors duration-200",
							design.TextColor(palette.Text),
							design.TextColor(palette.Primary),
						)),
						g.Text("Terms of Service"),
					),
				),
			),
		),
	)
}

func header() g.Node {
	palette := &design.StandardPalette

	return ghtml.Header(
		ghtml.Class(fmt.Sprintf("sticky top-0 z-50 %s border-b %s shadow-sm",
			design.Background(palette.Background),
			design.BorderColor(palette.Text),
		)),
		ghtml.Div(
			ghtml.Class("max-w-4xl mx-auto px-4"),
			ghtml.Div(
				ghtml.Class("flex items-center justify-between h-16"),
				ghtml.A(
					ghtml.Href("/"),
					ghtml.H1(
						ghtml.Class(fmt.Sprintf("text-xl font-bold %s hover:opacity-80 transition-opacity duration-200 cursor-pointer", design.TextColor(palette.Primary))),
						g.Text(branding.CompanyName),
					),
				),
				ghtml.Nav(
					ghtml.Class("flex space-x-4"),
					ghtml.A(
						ghtml.Href("/account/settings"),
						ghtml.Class(fmt.Sprintf("px-3 py-2 rounded-md text-sm font-medium %s hover:%s transition-colors duration-200",
							design.TextColor(palette.Text),
							design.Background(palette.Secondary),
						)),
						g.Text("Account"),
					),
					ghtml.A(
						ghtml.Href("/logout"),
						ghtml.Class(fmt.Sprintf("px-3 py-2 rounded-md text-sm font-medium %s hover:%s transition-colors duration-200",
							design.TextColor(palette.Text),
							design.Background(palette.Secondary),
						)),
						g.Text("Sign Out"),
					),
				),
			),
		),
	)
}
