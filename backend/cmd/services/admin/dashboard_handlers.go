package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *AdminFrontendServer) homeRoute(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	return s.HomePage(""), nil
}

func (s *AdminFrontendServer) HomePage(title string) g.Node {
	if title == "" {
		title = "Dashboard"
	}

	return page(title,
		components.ContentContainer(&components.ContentContainerProps{
			Title:    title,
			Subtitle: "Welcome to the admin dashboard",
			Palette:  &design.StandardPalette,
		},
			components.Card(&design.StandardPalette,
				ghtml.H2(
					ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
					g.Text("Quick Stats"),
				),
				ghtml.Div(
					ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
					statCard("Total Users", "1,234", &design.StandardPalette),
					statCard("Active Sessions", "89", &design.StandardPalette),
					statCard("System Status", "Healthy", &design.StandardPalette),
				),
			),
		),
	)
}

func statCard(title, value string, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Class(fmt.Sprintf("p-4 %s rounded-lg border %s",
			design.Background(design.Color{Value: "gray-50"}),
			design.BorderColor(palette.Background),
		)),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("text-sm font-medium %s", design.TextColor(palette.Text))),
			g.Text(title),
		),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("mt-1 text-2xl font-bold %s", design.TextColor(palette.Primary))),
			g.Text(value),
		),
	)
}

