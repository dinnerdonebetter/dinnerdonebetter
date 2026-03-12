package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	analyticsgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

var analyticsSourceOptions = []*struct {
	Label string
	Value string
}{
	{Label: "iOS", Value: "ios"},
	{Label: "Web", Value: "web"},
}

func (s *AdminFrontendServer) AnalyticsTestPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	if _, err := fetchClientFromContext(ctx); err != nil {
		return page("Analytics Test", renderAnalyticsTestError("Error: No API client available")), nil
	}

	palette := &design.StandardPalette

	var sourceOptions []g.Node
	for _, opt := range analyticsSourceOptions {
		sourceOptions = append(sourceOptions, ghtml.Option(
			ghtml.Value(opt.Value),
			g.Text(opt.Label),
		))
	}

	return page("Analytics Test",
		components.ContentContainer(&components.ContentContainerProps{
			Title:    "Analytics Test",
			Subtitle: "Fire a test event to the analytics proxy for a given source (iOS or Web)",
			Palette:  palette,
		},
			components.Card(palette,
				ghtml.Form(
					g.Attr("hx-post", "/api/analytics_test"),
					g.Attr("hx-target", "#analytics-test-result"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-ext", "json-enc"),
					g.Attr("hx-indicator", "#analytics-test-spinner"),
					ghtml.Div(
						ghtml.Class("space-y-4"),
						ghtml.Div(
							ghtml.Label(
								ghtml.For("source"),
								ghtml.Class(fmt.Sprintf("block text-sm font-medium %s mb-1", design.TextColor(palette.Text))),
								g.Text("Source"),
							),
							ghtml.Select(
								ghtml.ID("source"),
								ghtml.Name("source"),
								ghtml.Class(fmt.Sprintf("block w-full px-3 py-2 border %s rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-%s focus:border-%s sm:text-sm",
									design.BorderColor(palette.Background),
									palette.Primary.Value,
									palette.Primary.Value,
								)),
								g.Group(sourceOptions),
							),
						),
						ghtml.Div(
							ghtml.Class("flex items-center space-x-4"),
							ghtml.Button(
								ghtml.Type("submit"),
								ghtml.Class(fmt.Sprintf("inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white %s hover:%s focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-%s",
									design.Background(palette.Primary),
									design.Background(design.Color{Value: palette.Primary.Value + "-700"}),
									palette.Primary.Value,
								)),
								g.Text("Fire Test Event"),
							),
							ghtml.Div(
								ghtml.ID("analytics-test-spinner"),
								ghtml.Class("htmx-indicator"),
								components.LoadingSpinner(palette),
							),
						),
					),
				),
				ghtml.Div(
					ghtml.ID("analytics-test-result"),
					ghtml.Class("mt-4"),
				),
			),
		),
	), nil
}

type analyticsTestInput struct {
	Source string `json:"source"`
}

func (s *AdminFrontendServer) AnalyticsTestSubmit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return renderAnalyticsTestError("Error: No API client available"), nil
	}

	var input analyticsTestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return renderAnalyticsTestError(fmt.Sprintf("Error decoding request: %v", err)), nil
	}

	if input.Source == "" {
		return renderAnalyticsTestError("Source is required"), nil
	}

	_, err = c.TrackEvent(ctx, &analyticsgrpc.TrackEventRequest{
		Source: input.Source,
		Event:  "admin_test_event",
		Properties: map[string]string{
			"source": input.Source,
			"test":   "true",
		},
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "firing analytics test event")
		return renderAnalyticsTestError(fmt.Sprintf("Analytics test failed: %v", err)), nil
	}

	return renderAnalyticsTestSuccess(input.Source), nil
}

func renderAnalyticsTestSuccess(source string) g.Node {
	palette := &design.StandardPalette

	return ghtml.Div(
		ghtml.Class("rounded-md border border-green-300 bg-green-50 p-4"),
		ghtml.Div(
			ghtml.Class("flex items-center mb-3"),
			ghtml.Div(
				ghtml.Class("h-5 w-5 text-green-500 mr-2"),
				g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>`),
			),
			ghtml.Span(
				ghtml.Class("text-sm font-medium text-green-800"),
				g.Text("Test event sent successfully"),
			),
		),
		ghtml.Dl(
			ghtml.Class("grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm"),
			definitionItem("Source", source, palette),
			definitionItem("Event", "admin_test_event", palette),
		),
	)
}

func renderAnalyticsTestError(message string) g.Node {
	return ghtml.Div(
		ghtml.Class("rounded-md border border-red-300 bg-red-50 p-4"),
		ghtml.Div(
			ghtml.Class("flex items-center"),
			ghtml.Div(
				ghtml.Class("h-5 w-5 text-red-500 mr-2"),
				g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>`),
			),
			ghtml.Span(
				ghtml.Class("text-sm font-medium text-red-800"),
				g.Text(message),
			),
		),
	)
}
