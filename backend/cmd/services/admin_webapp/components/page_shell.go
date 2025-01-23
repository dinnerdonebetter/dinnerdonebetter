package components

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/components"
	ghtml "maragu.dev/gomponents/html"
)

func PageShell(title string, children ...gomponents.Node) gomponents.Node {
	return components.HTML5(components.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []gomponents.Node{
			ghtml.Script(
				ghtml.Src("https://unpkg.com/@tailwindcss/browser@4"),
				ghtml.Integrity("sha384-fsXZ0Oru5QjGkveFx8DdmBAyKdwnJ7TnbRzDN5LROCKt8PAN8h7y7oqCwtk9cN68"),
				ghtml.CrossOrigin("anonymous"),
			),
			ghtml.Script(
				ghtml.Src("https://unpkg.com/htmx.org@2.0.4"),
				ghtml.Integrity("sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"),
				ghtml.CrossOrigin("anonymous"),
			),
			ghtml.Script(
				ghtml.Src("https://unpkg.com/htmx.org@2.0.4/dist/ext/json-enc.js"),
				ghtml.Integrity("sha384-j+tqxLrwDkbeOdjbpaVekgvQL/J7qm/yh/UqSEs6RjEtnBFHqlJViBWG/jBZ6I2p"),
				ghtml.CrossOrigin("anonymous"),
			),
		},
		Body: []gomponents.Node{ghtml.Class("bg-gradient-to-b from-white to-indigo-100 bg-no-repeat"),
			ghtml.Div(ghtml.Class("min-h-screen flex flex-col justify-between"),
				// header
				ghtml.Div(ghtml.ID("main"), ghtml.Class("bg-gray-600 h-8 text-white shadow"),
					container(false,
						ghtml.Div(ghtml.Class("flex items-center space-x-4 h-8"),
							headerLink("/", "Home"),
							headerLink("/login", "Login"),
						),
					),
				),
				// main body
				ghtml.Div(ghtml.Class("grow flex justify-around"),
					container(true,
						gomponents.Group(children),
					),
				),
				// footer
				ghtml.Div(ghtml.Class("bg-gray-600 text-white shadow text-center h-12 flex items-center justify-center"),
					ghtml.A(ghtml.Href("/"), gomponents.Text("dinner-done-better")),
				),
			),
		},
	})
}

func headerLink(href, text string) gomponents.Node {
	return ghtml.A(ghtml.Class("hover:text-indigo-300"), ghtml.Href(href), gomponents.Text(text))
}

func container(padY bool, children ...gomponents.Node) gomponents.Node {
	return ghtml.Div(
		components.Classes{
			"w-full max-w-[90%]":    true,
			"px-4 md:px-8 lg:px-16": true,
			"py-4 md:py-8":          padY,
		},
		gomponents.Group(children),
	)
}
