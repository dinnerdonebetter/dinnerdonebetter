package main

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

func main() {
	mux := http.NewServeMux()

	Home(mux)
	Login(mux)

	if err := http.ListenAndServe(":8888", mux); err != nil {
		panic(err)
	}
}

func HomePage() g.Node {
	return page("Home",
		ghtml.H1(g.Text("Home")),

		ghtml.P(g.Text("This is the g example app!")),
	)
}

func Home(mux *http.ServeMux) {
	mux.Handle("GET /", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (g.Node, error) {
		return HomePage(), nil
	}))
}

func LoginPage() g.Node {
	return page("Login",
		components.LoginForm("", "", "", ""),
	)
}

func Login(mux *http.ServeMux) {
	mux.Handle("GET /login", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		return LoginPage(), nil
	}))
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

func page(title string, children ...g.Node) g.Node {
	return ghtml.HTML(
		ghtml.Lang("en"),
		ghtml.Head(
			ghtml.Title(title),
			ghtml.Script(ghtml.Src("https://cdn.tailwindcss.com?plugins=typography")),
			ghtml.Script(
				ghtml.Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.min.js"),
				ghtml.Integrity("sha384-ZBXiYtYQ6hJ2Y0ZNoYuI+Nq5MqWBr+chMrS/RkXpNzQCApHEhOt2aY8EJgqwHLkJ"),
				ghtml.CrossOrigin("anonymous"),
			),
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
