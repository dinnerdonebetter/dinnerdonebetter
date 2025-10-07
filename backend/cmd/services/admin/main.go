package main

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/encoding"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	o11yName = "admin_frontend"
)

func main() {
	mux := http.NewServeMux()

	ctx := context.Background()

	fs, err := NewAdminFrontendServer(
		ctx,
		nil,
		nil,
		encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		mux,
	)
	if err != nil {
		panic(err)
	}

	if err = http.ListenAndServe(":8888", fs); err != nil {
		panic(err)
	}
}

///

func header() g.Node {
	return ghtml.Header(
		ghtml.Class("text-center py-6"),
		ghtml.H1(
			ghtml.Class("text-3xl font-bold text-indigo-700"),
			g.Text("My App"),
		),
	)
}

var (
	tailwindImport = ghtml.Script(ghtml.Src("https://cdn.tailwindcss.com?plugins=typography"))

	htmxImport = g.Group{
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
	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Title(title),
			tailwindImport,
			htmxImport,
		),
		ghtml.Body(
			ghtml.Class("bg-gradient-to-b from-white to-indigo-100 min-h-screen flex flex-col"),
			header(),
			ghtml.Main(
				ghtml.Class("flex-grow flex justify-center items-center w-full"),
				ghtml.Div(
					ghtml.Class("w-full max-w-md"),
					g.Group(children),
				),
			),
			footer(),
		),
	)
}

func footer() g.Node {
	return ghtml.Footer(
		ghtml.Class("text-center py-4 text-sm text-gray-600"),
		g.Text("© 2025 My App. All rights reserved."),
	)
}
