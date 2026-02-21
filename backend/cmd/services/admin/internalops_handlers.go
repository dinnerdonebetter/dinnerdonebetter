package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

var queueOptions = []struct {
	Label string
	Value string
}{
	{Label: "Data Changes", Value: "data-changes"},
	{Label: "Outbound Emails", Value: "outbound-emails"},
	{Label: "Search Index Requests", Value: "search-index-requests"},
	{Label: "User Data Aggregation", Value: "user-data-aggregation"},
	{Label: "Webhook Execution Requests", Value: "webhook-execution-requests"},
}

func (s *AdminFrontendServer) QueueTestPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	if _, err := fetchClientFromContext(ctx); err != nil {
		return page("Queue Test", renderQueueTestError("Error: No API client available")), nil
	}

	palette := &design.StandardPalette

	var options []g.Node
	for _, opt := range queueOptions {
		options = append(options, ghtml.Option(
			ghtml.Value(opt.Value),
			g.Text(opt.Label),
		))
	}

	return page("Queue Test",
		components.ContentContainer(&components.ContentContainerProps{
			Title:    "Queue Test",
			Subtitle: "Send a test message through a queue and verify round-trip delivery",
			Palette:  palette,
		},
			components.Card(palette,
				ghtml.FormEl(
					g.Attr("hx-post", "/api/queue_test"),
					g.Attr("hx-target", "#queue-test-result"),
					g.Attr("hx-swap", "innerHTML"),
					g.Attr("hx-ext", "json-enc"),
					g.Attr("hx-indicator", "#queue-test-spinner"),
					ghtml.Div(
						ghtml.Class("space-y-4"),
						ghtml.Div(
							ghtml.Label(
								ghtml.For("queue_name"),
								ghtml.Class(fmt.Sprintf("block text-sm font-medium %s mb-1", design.TextColor(palette.Text))),
								g.Text("Queue"),
							),
							ghtml.Select(
								ghtml.ID("queue_name"),
								ghtml.Name("queue_name"),
								ghtml.Class(fmt.Sprintf("block w-full px-3 py-2 border %s rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-%s focus:border-%s sm:text-sm",
									design.BorderColor(palette.Background),
									palette.Primary.Value,
									palette.Primary.Value,
								)),
								g.Group(options),
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
								g.Text("Test Queue"),
							),
							ghtml.Div(
								ghtml.ID("queue-test-spinner"),
								ghtml.Class("htmx-indicator"),
								components.LoadingSpinner(palette),
							),
						),
					),
				),
				ghtml.Div(
					ghtml.ID("queue-test-result"),
					ghtml.Class("mt-4"),
				),
			),
		),
	), nil
}

type queueTestInput struct {
	QueueName string `json:"queue_name"`
}

func (s *AdminFrontendServer) QueueTestSubmit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return renderQueueTestError("Error: No API client available"), nil
	}

	var input queueTestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return renderQueueTestError(fmt.Sprintf("Error decoding request: %v", err)), nil
	}

	if input.QueueName == "" {
		return renderQueueTestError("Queue name is required"), nil
	}

	res, err := c.TestQueueMessage(ctx, &internalopssvc.TestQueueMessageRequest{
		QueueName: input.QueueName,
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "testing queue message")
		return renderQueueTestError(fmt.Sprintf("Queue test failed: %v", err)), nil
	}

	return renderQueueTestSuccess(input.QueueName, res.TestId, res.RoundTripMs), nil
}

func renderQueueTestSuccess(queueName, testID string, roundTripMs int64) g.Node {
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
				g.Text("Queue test passed"),
			),
		),
		ghtml.Dl(
			ghtml.Class("grid grid-cols-1 sm:grid-cols-3 gap-3 text-sm"),
			definitionItem("Queue", queueName, palette),
			definitionItem("Test ID", testID, palette),
			definitionItem("Round Trip", fmt.Sprintf("%d ms", roundTripMs), palette),
		),
	)
}

func renderQueueTestError(message string) g.Node {
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

func definitionItem(label, value string, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Dt(
			ghtml.Class(fmt.Sprintf("text-xs font-medium %s", design.TextColor(palette.Text))),
			g.Text(label),
		),
		ghtml.Dd(
			ghtml.Class(fmt.Sprintf("mt-1 font-semibold %s", design.TextColor(palette.Primary))),
			g.Text(value),
		),
	)
}
