package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *AdminFrontendServer) renderMeasurementUnitConversionMismatchesError(message string) g.Node {
	return ghtml.Div(
		ghtml.Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8"),
		ghtml.Div(
			ghtml.Class("rounded-md p-4 bg-red-50 border border-red-200"),
			ghtml.P(
				ghtml.Class("text-sm font-medium text-red-800"),
				g.Text(message),
			),
		),
	)
}

func (s *AdminFrontendServer) MeasurementUnitConversionMismatchesPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Measurement Unit Conversion Mismatches", s.renderMeasurementUnitConversionMismatchesError("Error: No API client available")), nil
	}

	res, err := c.GetMeasurementUnitConversionMismatches(ctx, &mealplanningsvc.GetMeasurementUnitConversionMismatchesRequest{})
	if err != nil {
		return page("Measurement Unit Conversion Mismatches", s.renderMeasurementUnitConversionMismatchesError(fmt.Sprintf("Error loading mismatches: %v", err))), nil
	}

	content := s.renderMeasurementUnitConversionMismatchesList(res.Mismatches)
	return page("Measurement Unit Conversion Mismatches", content), nil
}

func (s *AdminFrontendServer) MeasurementUnitConversionMismatchesList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	res, err := c.GetMeasurementUnitConversionMismatches(ctx, &mealplanningsvc.GetMeasurementUnitConversionMismatchesRequest{})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading mismatches: %v", err)),
			),
		), nil
	}

	return s.renderMeasurementUnitConversionMismatchesList(res.Mismatches), nil
}

func (s *AdminFrontendServer) renderMeasurementUnitConversionMismatchesList(mismatches []*mealplanningsvc.MeasurementUnitConversionMismatch) g.Node {
	emptyState := ghtml.Div(
		ghtml.Class("text-center py-12 px-4 bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg"),
		ghtml.P(
			ghtml.Class("text-gray-600 font-medium mb-1"),
			g.Text("No conversion mismatches found"),
		),
		ghtml.P(
			ghtml.Class("text-sm text-gray-500"),
			g.Text("All ingredients with multiple measurement units have conversions defined."),
		),
	)

	if len(mismatches) == 0 {
		return components.ContentContainer(&components.ContentContainerProps{
			Title:    "Measurement Unit Conversion Mismatches",
			Subtitle: "Ingredients with multiple units that lack conversion data",
			Palette:  &design.StandardPalette,
		}, emptyState)
	}

	var rows []g.Node
	for _, m := range mismatches {
		ingredientName := "-"
		ingredientID := ""
		if m.Ingredient != nil {
			ingredientName = m.Ingredient.Name
			ingredientID = m.Ingredient.Id
		}
		fromUnitName := "-"
		fromUnitID := ""
		if m.FromUnit != nil {
			fromUnitName = m.FromUnit.Name
			fromUnitID = m.FromUnit.Id
		}
		toUnitName := "-"
		toUnitID := ""
		if m.ToUnit != nil {
			toUnitName = m.ToUnit.Name
			toUnitID = m.ToUnit.Id
		}

		addConversionURL := fmt.Sprintf("/measurement_unit_conversion_mismatches/add_conversion?from=%s&to=%s&ingredient=%s", fromUnitID, toUnitID, ingredientID)

		rows = append(rows, ghtml.Tr(
			ghtml.Class("bg-white hover:bg-gray-50 border-b border-gray-200"),
			ghtml.Td(
				ghtml.Class("px-4 py-3 text-sm text-gray-900"),
				g.If(
					ingredientID != "",
					ghtml.A(
						ghtml.Href(fmt.Sprintf("/valid_ingredients/%s", ingredientID)),
						ghtml.Class("text-blue-600 hover:text-blue-800 hover:underline"),
						g.Text(ingredientName),
					),
				),
				g.If(
					ingredientID == "",
					g.Text(ingredientName),
				),
			),
			ghtml.Td(
				ghtml.Class("px-4 py-3 text-sm text-gray-900"),
				g.Text(fromUnitName),
			),
			ghtml.Td(
				ghtml.Class("px-4 py-3 text-sm text-gray-900"),
				g.Text(toUnitName),
			),
			ghtml.Td(
				ghtml.Class("px-4 py-3 text-sm"),
				ghtml.A(
					ghtml.Href(addConversionURL),
					ghtml.Class(fmt.Sprintf("px-3 py-1.5 text-sm font-medium rounded-md %s %s hover:opacity-90",
						design.Background(design.StandardPalette.Primary),
						design.TextColor(design.Color{Value: "white"}),
					)),
					g.Text("Add Conversion"),
				),
			),
		))
	}

	table := ghtml.Div(
		ghtml.Class("overflow-x-auto"),
		ghtml.Table(
			ghtml.Class("min-w-full divide-y divide-gray-200"),
			ghtml.THead(
				ghtml.Class("bg-gray-50"),
				ghtml.Tr(
					ghtml.Th(
						ghtml.Class("px-4 py-2 text-left text-sm font-semibold text-gray-800 uppercase tracking-wide"),
						g.Text("Ingredient"),
					),
					ghtml.Th(
						ghtml.Class("px-4 py-2 text-left text-sm font-semibold text-gray-800 uppercase tracking-wide"),
						g.Text("From Unit"),
					),
					ghtml.Th(
						ghtml.Class("px-4 py-2 text-left text-sm font-semibold text-gray-800 uppercase tracking-wide"),
						g.Text("To Unit"),
					),
					ghtml.Th(
						ghtml.Class("px-4 py-2 text-left text-sm font-semibold text-gray-800 uppercase tracking-wide"),
						g.Text("Action"),
					),
				),
			),
			ghtml.TBody(
				ghtml.Class("bg-white divide-y divide-gray-200"),
				g.Group(rows),
			),
		),
	)

	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Measurement Unit Conversion Mismatches",
		Subtitle: fmt.Sprintf("Viewing %d conversion mismatches", len(mismatches)),
		Palette:  &design.StandardPalette,
	}, table)
}

func (s *AdminFrontendServer) AddConversionPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	fromID := req.URL.Query().Get("from")
	toID := req.URL.Query().Get("to")
	ingredientID := req.URL.Query().Get("ingredient")

	if fromID == "" || toID == "" {
		return page("Add Conversion", s.renderMeasurementUnitConversionMismatchesError("Error: Missing from or to unit in URL")), nil
	}

	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Add Conversion", s.renderMeasurementUnitConversionMismatchesError("Error: No API client available")), nil
	}

	// Fetch unit and ingredient names for display
	fromName := fromID
	toName := toID
	ingredientName := ""

	if resp, getErr := c.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitId: fromID}); getErr == nil && resp != nil && resp.Result != nil {
		fromName = resp.Result.Name
	}
	if resp, getErr := c.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitId: toID}); getErr == nil && resp != nil && resp.Result != nil {
		toName = resp.Result.Name
	}
	if ingredientID != "" {
		if resp, getErr := c.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientId: ingredientID}); getErr == nil && resp != nil && resp.Result != nil {
			ingredientName = resp.Result.Name
		}
	}

	form := ghtml.Form(
		ghtml.Action("/measurement_unit_conversion_mismatches/add_conversion"),
		ghtml.Method("POST"),
		ghtml.Class("space-y-4 max-w-md"),
		ghtml.Div(
			ghtml.Class("space-y-2"),
			ghtml.Label(ghtml.Class("block text-sm font-medium text-gray-700"), g.Text("From Unit")),
			ghtml.Input(
				ghtml.Type("hidden"),
				ghtml.Name("from"),
				ghtml.Value(fromID),
			),
			ghtml.P(ghtml.Class("text-gray-900"), g.Text(fromName)),
		),
		ghtml.Div(
			ghtml.Class("space-y-2"),
			ghtml.Label(ghtml.Class("block text-sm font-medium text-gray-700"), g.Text("To Unit")),
			ghtml.Input(
				ghtml.Type("hidden"),
				ghtml.Name("to"),
				ghtml.Value(toID),
			),
			ghtml.P(ghtml.Class("text-gray-900"), g.Text(toName)),
		),
		g.If(
			ingredientID != "",
			ghtml.Div(
				ghtml.Class("space-y-2"),
				ghtml.Label(ghtml.Class("block text-sm font-medium text-gray-700"), g.Text("Only For Ingredient")),
				ghtml.Input(
					ghtml.Type("hidden"),
					ghtml.Name("onlyForIngredient"),
					ghtml.Value(ingredientID),
				),
				ghtml.P(ghtml.Class("text-gray-900"), g.Text(ingredientName)),
			),
		),
		ghtml.Div(
			ghtml.Class("space-y-2"),
			ghtml.Label(
				ghtml.For("modifier"),
				ghtml.Class("block text-sm font-medium text-gray-700"),
				g.Text("Modifier"),
				ghtml.Span(ghtml.Class("text-red-500 ml-1"), g.Text("*")),
			),
			ghtml.Input(
				ghtml.Type("number"),
				ghtml.ID("modifier"),
				ghtml.Name("modifier"),
				ghtml.Step("any"),
				ghtml.Required(),
				ghtml.Placeholder("e.g., 14.2 (1 from unit = 14.2 to units)"),
				ghtml.Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
			),
			ghtml.P(
				ghtml.Class("text-xs text-gray-500"),
				g.Text("How many 'to' units equal 1 'from' unit"),
			),
		),
		ghtml.Div(
			ghtml.Class("space-y-2"),
			ghtml.Label(
				ghtml.For("notes"),
				ghtml.Class("block text-sm font-medium text-gray-700"),
				g.Text("Notes (optional)"),
			),
			ghtml.Textarea(
				ghtml.ID("notes"),
				ghtml.Name("notes"),
				ghtml.Rows("2"),
				ghtml.Placeholder("Additional notes..."),
				ghtml.Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
			),
		),
		ghtml.Div(
			ghtml.Class("flex gap-2"),
			ghtml.Button(
				ghtml.Type("submit"),
				ghtml.Class(fmt.Sprintf("px-4 py-2 %s text-white rounded-md hover:opacity-90 font-medium", design.Background(design.StandardPalette.Primary))),
				g.Text("Add Conversion"),
			),
			ghtml.A(
				ghtml.Href("/measurement_unit_conversion_mismatches"),
				ghtml.Class("px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium"),
				g.Text("Cancel"),
			),
		),
	)

	return page("Add Conversion", components.ContentContainer(&components.ContentContainerProps{
		Title:    "Add Conversion",
		Subtitle: fmt.Sprintf("Add conversion: %s ⇄ %s", fromName, toName),
		Palette:  &design.StandardPalette,
	}, form)), nil
}

func (s *AdminFrontendServer) AddConversionSubmit(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	if err := req.ParseForm(); err != nil {
		http.Error(res, "Invalid form", http.StatusBadRequest)
		return g.El("div"), nil
	}

	fromID := req.FormValue("from")
	toID := req.FormValue("to")
	modifierStr := req.FormValue("modifier")
	notes := req.FormValue("notes")
	onlyForIngredient := req.FormValue("onlyForIngredient")

	if fromID == "" || toID == "" || modifierStr == "" {
		http.Error(res, "Missing required fields", http.StatusBadRequest)
		return g.El("div"), nil
	}

	var modifier float32
	if _, err := fmt.Sscanf(modifierStr, "%f", &modifier); err != nil {
		http.Error(res, "Invalid modifier", http.StatusBadRequest)
		return g.El("div"), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		http.Error(res, "No API client", http.StatusInternalServerError)
		return g.El("div"), nil
	}

	input := &mealplanningsvc.ValidMeasurementUnitConversionCreationRequestInput{
		From:     fromID,
		To:       toID,
		Modifier: modifier,
		Notes:    notes,
	}
	if onlyForIngredient != "" {
		input.OnlyForIngredient = &onlyForIngredient
	}

	_, err = c.CreateValidMeasurementUnitConversion(ctx, &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{
		Input: input,
	})
	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to create conversion: %v", err), http.StatusInternalServerError)
		return g.El("div"), nil
	}

	http.Redirect(res, req, "/measurement_unit_conversion_mismatches", http.StatusSeeOther)
	return g.El("div"), nil
}
