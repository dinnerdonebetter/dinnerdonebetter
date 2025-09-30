package main

import (
	"net/http"
	"time"

	"maragu.dev/gomponents"
	"maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

func main() {
	mux := http.NewServeMux()

	Home(mux)
	About(mux)

	if err := http.ListenAndServe(":8888", mux); err != nil {
		panic(err)
	}
}

type adminFrontendServer struct{}

func HomePage() gomponents.Node {
	return page("Home",
		html.H1(gomponents.Text("Home")),

		html.P(gomponents.Text("This is the gomponents example app!")),
	)
}

func Home(mux *http.ServeMux) {
	mux.Handle("GET /", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		return HomePage(), nil
	}))
}

func AboutPage() gomponents.Node {
	now := time.Now()

	return page("About",
		html.H1(gomponents.Text("About")),

		html.P(gomponents.Textf("Built with gomponents and rendered at %v.", now.Format(time.TimeOnly))),

		html.P(
			gomponents.If(now.Second()%2 == 0, gomponents.Text("It's an even second!")),
			gomponents.If(now.Second()%2 != 0, gomponents.Text("It's an odd second!")),
		),
	)
}

func About(mux *http.ServeMux) {
	mux.Handle("GET /about", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gomponents.Node, error) {
		return AboutPage(), nil
	}))
}

///

func page(title string, children ...gomponents.Node) gomponents.Node {
	return components.HTML5(components.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []gomponents.Node{
			html.Script(html.Src("https://cdn.tailwindcss.com?plugins=typography")),
		},
		Body: []gomponents.Node{html.Class("bg-gradient-to-b from-white to-indigo-100 bg-no-repeat"),
			html.Div(html.Class("min-h-screen flex flex-col justify-between"),
				header(),
				html.Div(html.Class("grow"),
					container(true,
						html.Div(html.Class("prose prose-lg prose-indigo"),
							gomponents.Group(children),
						),
					),
				),
				footer(),
			),
		},
	})
}

func header() gomponents.Node {
	return html.Div(html.Class("bg-indigo-600 text-white shadow"),
		container(false,
			html.Div(html.Class("flex items-center space-x-4 h-8"),
				headerLink("/", "Home"),
				headerLink("/about", "About"),
			),
		),
	)
}

func headerLink(href, text string) gomponents.Node {
	return html.A(html.Class("hover:text-indigo-300"), html.Href(href), gomponents.Text(text))
}

func container(padY bool, children ...gomponents.Node) gomponents.Node {
	return html.Div(
		components.Classes{
			"max-w-7xl mx-auto":     true,
			"px-4 md:px-8 lg:px-16": true,
			"py-4 md:py-8":          padY,
		},
		gomponents.Group(children),
	)
}

func footer() gomponents.Node {
	return html.Div(html.Class("bg-gray-900 text-white shadow text-center h-16 flex items-center justify-center"))
}
