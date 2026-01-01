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

// ValidMeasurementUnitConversionsForUnit displays conversions for a measurement unit.
func (s *AdminFrontendServer) ValidMeasurementUnitConversionsForUnit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)
	if measurementUnitID == "" {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No measurement unit ID provided"),
		), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	fromRes, err := c.GetValidMeasurementUnitConversionsForUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitConversionsForUnitRequest{
		ValidMeasurementUnitId: measurementUnitID,
	})
	if err != nil {
		s.logger.Error("error fetching conversions from unit", err)
	}

	return renderMeasurementUnitConversions(measurementUnitID, fromRes.Results, &design.StandardPalette), nil
}

// renderMeasurementUnitConversions creates a custom display for unit conversions.
func renderMeasurementUnitConversions(measurementUnitID string, conversions []*mealplanningsvc.ValidMeasurementUnitConversion, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Class("space-y-4"),

		// Form for creating new conversions
		ghtml.Div(
			ghtml.ID("conversion-create-form"),
			ghtml.Class("bg-gray-50 p-6 rounded-lg border border-gray-200"),
			ghtml.H3(
				ghtml.Class("text-lg font-medium text-gray-900 mb-4"),
				g.Text("Add New Conversion"),
			),
			ghtml.Form(
				g.Attr("hx-post", "/api/valid_measurement_unit_conversions"),
				g.Attr("hx-target", "#conversions-container"),
				g.Attr("hx-swap", "innerHTML"),
				g.Attr("hx-on::after-request", "if(event.detail.successful) this.reset(); document.getElementById('to-unit-search-hidden').value = ''; document.getElementById('ingredient-search-hidden').value = '';"),
				ghtml.Class("space-y-4"),

				// Hidden field for "from" unit (current unit)
				ghtml.Input(
					ghtml.Type("hidden"),
					ghtml.Name("from"),
					ghtml.Value(measurementUnitID),
				),

				// To Unit search
				ghtml.Div(
					ghtml.Class("space-y-2"),
					ghtml.Label(
						ghtml.Class("block text-sm font-medium text-gray-700"),
						g.Text("To Unit"),
						ghtml.Span(ghtml.Class("text-red-500 ml-1"), g.Text("*")),
					),
					components.SearchBox(&components.SearchBoxProps{
						ID:              "to-unit-search",
						Placeholder:     "Search for a measurement unit...",
						SearchEndpoint:  "/admin/search/valid_measurement_units",
						HiddenFieldName: "to",
						Palette:         palette,
					}),
				),

				// Modifier field
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
						ghtml.Placeholder("e.g., 1000 (if 1 from = 1000 to)"),
						ghtml.Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
					),
					ghtml.P(
						ghtml.Class("text-xs text-gray-500"),
						g.Text("How many 'to' units equal 1 'from' unit"),
					),
				),

				// Only For Ingredient search (optional)
				ghtml.Div(
					ghtml.Class("space-y-2"),
					ghtml.Label(
						ghtml.Class("block text-sm font-medium text-gray-700"),
						g.Text("Only For Ingredient (optional)"),
					),
					components.SearchBox(&components.SearchBoxProps{
						ID:              "ingredient-search",
						Placeholder:     "Search for an ingredient...",
						SearchEndpoint:  "/admin/search/valid_ingredients",
						HiddenFieldName: "onlyForIngredient",
						Palette:         palette,
					}),
					ghtml.P(
						ghtml.Class("text-xs text-gray-500"),
						g.Text("If specified, this conversion only applies to this ingredient"),
					),
				),

				// Notes field
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
						ghtml.Placeholder("Additional notes about this conversion..."),
						ghtml.Class("w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"),
					),
				),

				// Submit button
				ghtml.Div(
					ghtml.Class("flex justify-end"),
					ghtml.Button(
						ghtml.Type("submit"),
						ghtml.Class(fmt.Sprintf("px-4 py-2 bg-%s text-white rounded-md hover:opacity-90 font-medium", palette.Primary.Value)),
						g.Text("Add Conversion"),
					),
				),
			),
		),

		// All conversions in one list
		g.If(
			len(conversions) == 0,
			ghtml.Div(
				ghtml.Class("text-center py-12 px-4 bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg"),
				ghtml.Div(
					ghtml.Class("mx-auto w-16 h-16 mb-4 flex items-center justify-center rounded-full bg-gray-100"),
					g.Raw(`<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
					</svg>`),
				),
				ghtml.P(
					ghtml.Class("text-gray-600 font-medium mb-1"),
					g.Text("No conversions defined"),
				),
				ghtml.P(
					ghtml.Class("text-sm text-gray-500"),
					g.Text("Use the search box above to add conversions."),
				),
			),
		),
		g.If(
			len(conversions) > 0,
			ghtml.Div(
				ghtml.Class("space-y-2"),
				g.Map(conversions, func(conv *mealplanningsvc.ValidMeasurementUnitConversion) g.Node {
					return renderConversionItem(conv, measurementUnitID)
				}),
			),
		),
	)
}

// pluralizeUnit returns the appropriate form of the unit name based on quantity.
func pluralizeUnit(singular, plural string, quantity float32) string {
	if quantity == 1.0 {
		return singular
	}
	if plural != "" {
		return plural
	}
	return singular
}

// renderConversionItem renders a single conversion item.
func renderConversionItem(conv *mealplanningsvc.ValidMeasurementUnitConversion, currentUnitID string) g.Node {
	// Safely extract ingredient name if available
	var ingredientName string
	if conv.OnlyForIngredient != nil {
		ingredientName = conv.OnlyForIngredient.Name
	}

	// Generate both example calculations
	forwardCalculation := fmt.Sprintf("1 %s = %.4g %s",
		conv.From.Name,
		conv.Modifier,
		pluralizeUnit(conv.To.Name, conv.To.PluralName, conv.Modifier))

	// Calculate reverse (1 / modifier)
	reverseModifier := 1.0 / float64(conv.Modifier)
	reverseCalculation := fmt.Sprintf("1 %s = %.4g %s",
		conv.To.Name,
		reverseModifier,
		pluralizeUnit(conv.From.Name, conv.From.PluralName, float32(reverseModifier)))

	// Determine which unit is current for highlighting
	fromIsCurrent := conv.From.Id == currentUnitID
	toIsCurrent := conv.To.Id == currentUnitID

	return ghtml.Div(
		ghtml.Class("flex items-start justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-300 transition-colors"),

		// Conversion info
		ghtml.Div(
			ghtml.Class("flex-1"),
			// Unit names with arrow between them
			ghtml.Div(
				ghtml.Class("flex items-center gap-3 mb-3"),
				// From unit
				ghtml.Span(
					ghtml.Class(func() string {
						if fromIsCurrent {
							return "font-semibold text-gray-900 text-lg"
						}
						return "font-medium text-gray-700 text-lg"
					}()),
					g.Text(conv.From.Name),
				),
				// Bidirectional arrow
				ghtml.Span(
					ghtml.Class("text-2xl text-blue-500"),
					g.Text("⇄"),
				),
				// To unit
				ghtml.Span(
					ghtml.Class(func() string {
						if toIsCurrent {
							return "font-semibold text-gray-900 text-lg"
						}
						return "font-medium text-gray-700 text-lg"
					}()),
					g.Text(conv.To.Name),
				),
			),
			// Both example calculations
			ghtml.Div(
				ghtml.Class("space-y-1.5"),
				ghtml.Div(
					ghtml.Class("text-sm font-mono text-blue-600 bg-blue-50 px-2 py-1 rounded inline-block"),
					g.Text(forwardCalculation),
				),
				ghtml.Div(
					ghtml.Class("text-sm font-mono text-green-600 bg-green-50 px-2 py-1 rounded inline-block"),
					g.Text(reverseCalculation),
				),
			),
			ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-2"),
				g.Text(fmt.Sprintf("Modifier: %.6g", conv.Modifier)),
			),
			g.If(
				ingredientName != "",
				ghtml.Div(
					ghtml.Class("mt-2 text-xs text-amber-600 bg-amber-50 px-2 py-1 rounded inline-block"),
					g.Text(fmt.Sprintf("Only for: %s", ingredientName)),
				),
			),
			g.If(
				conv.Notes != "",
				ghtml.P(
					ghtml.Class("text-xs text-gray-500 mt-2 italic"),
					g.Text("Notes: "+conv.Notes),
				),
			),
		),

		// Delete button
		ghtml.Button(
			ghtml.Type("button"),
			ghtml.Class("ml-3 p-2 text-red-600 hover:bg-red-50 rounded-md transition-colors"),
			g.Attr("hx-delete", "/api/valid_measurement_unit_conversions/"+conv.Id),
			g.Attr("hx-target", "#conversions-container"),
			g.Attr("hx-swap", "outerHTML"),
			g.Attr("hx-confirm", "Are you sure you want to remove this conversion?"),
			ghtml.Title("Remove conversion"),
			g.Raw(`<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
			</svg>`),
		),
	)
}

// SearchMeasurementUnitsForConversion searches for measurement units to create conversions.
func (s *AdminFrontendServer) SearchMeasurementUnitsForConversion(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)
	query := req.URL.Query().Get("q")

	if query == "" {
		return ghtml.Div(), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Search for measurement units
	searchRes, err := c.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{
		Filter: nil,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error searching: %v", err)),
		), nil
	}

	// Filter results by query and exclude the current unit
	type searchResult struct {
		ID          string
		Name        string
		Description string
	}
	var results []searchResult
	for _, unit := range searchRes.Results {
		if unit.Id != measurementUnitID && (contains(unit.Name, query) || contains(unit.Description, query)) {
			results = append(results, searchResult{
				ID:          unit.Id,
				Name:        unit.Name,
				Description: unit.Description,
			})
		}
	}

	if len(results) == 0 {
		return ghtml.Div(
			ghtml.Class("text-sm text-gray-500 py-2"),
			g.Text("No measurement units found matching your search."),
		), nil
	}

	return ghtml.Div(
		ghtml.Class("max-h-60 overflow-y-auto space-y-1 border border-gray-200 rounded-md bg-white"),
		g.Map(results, func(item searchResult) g.Node {
			return ghtml.Div(
				ghtml.Class("px-3 py-2 border-b border-gray-100 last:border-b-0"),
				ghtml.Div(
					ghtml.Class("font-medium text-gray-900"),
					g.Text(item.Name),
				),
				g.If(
					item.Description != "",
					ghtml.Div(
						ghtml.Class("text-sm text-gray-600 mb-2"),
						g.Text(item.Description),
					),
				),
				// Two buttons for choosing direction
				ghtml.Div(
					ghtml.Class("flex gap-2 mt-2"),
					ghtml.Button(
						ghtml.Type("button"),
						ghtml.Class("flex-1 px-3 py-1.5 text-sm bg-blue-50 text-blue-700 rounded hover:bg-blue-100 transition-colors"),
						g.Attr("hx-post", fmt.Sprintf("/api/valid_measurement_units/%s/conversions?direction=from", measurementUnitID)),
						g.Attr("hx-target", "#conversions-container"),
						g.Attr("hx-swap", "outerHTML"),
						g.Attr("hx-vals", fmt.Sprintf(`{"id": %q}`, item.ID)),
						g.Text("Current → "+item.Name),
					),
					ghtml.Button(
						ghtml.Type("button"),
						ghtml.Class("flex-1 px-3 py-1.5 text-sm bg-green-50 text-green-700 rounded hover:bg-green-100 transition-colors"),
						g.Attr("hx-post", fmt.Sprintf("/api/valid_measurement_units/%s/conversions?direction=to", measurementUnitID)),
						g.Attr("hx-target", "#conversions-container"),
						g.Attr("hx-swap", "outerHTML"),
						g.Attr("hx-vals", fmt.Sprintf(`{"id": %q}`, item.ID)),
						g.Text(item.Name+" → Current"),
					),
				),
			)
		}),
	), nil
}

// CreateMeasurementUnitConversion creates a new conversion.
func (s *AdminFrontendServer) CreateMeasurementUnitConversion(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Parse form data
	if err = req.ParseForm(); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error parsing form: %v", err)),
		), nil
	}

	// Extract form values
	fromID := req.FormValue("from")
	toID := req.FormValue("to")
	modifierStr := req.FormValue("modifier")
	notes := req.FormValue("notes")
	onlyForIngredient := req.FormValue("onlyForIngredient")

	// Validate required fields
	if fromID == "" || toID == "" || modifierStr == "" {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: Missing required fields (from, to, modifier)"),
		), nil
	}

	// Parse modifier as float32
	var modifier float32
	if _, err = fmt.Sscanf(modifierStr, "%f", &modifier); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error: Invalid modifier value: %v", err)),
		), nil
	}

	// Build the creation input
	input := &mealplanningsvc.ValidMeasurementUnitConversionCreationRequestInput{
		From:     fromID,
		To:       toID,
		Modifier: modifier,
		Notes:    notes,
	}

	// Add optional ingredient restriction if provided
	if onlyForIngredient != "" {
		input.OnlyForIngredient = &onlyForIngredient
	}

	// Create the conversion
	_, err = c.CreateValidMeasurementUnitConversion(ctx, &mealplanningsvc.CreateValidMeasurementUnitConversionRequest{
		Input: input,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating conversion: %v", err)),
		), nil
	}

	// Return the updated conversions list - we need to construct a proper request with the measurement unit ID
	// Extract the measurement unit ID from the 'from' field since that's the current unit
	updatedReq := req.Clone(ctx)
	// Set the path value for the measurement unit ID fetcher
	updatedReq.SetPathValue(validMeasurementUnitIDURLParamKey, fromID)

	return s.ValidMeasurementUnitConversionsForUnit(nil, updatedReq)
}

// DeleteMeasurementUnitConversion deletes a conversion.
func (s *AdminFrontendServer) DeleteMeasurementUnitConversion(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	// Get the conversion ID from the URL path
	conversionID := req.PathValue("conversionID")
	if conversionID == "" {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No conversion ID provided"),
		), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Archive (delete) the conversion
	_, err = c.ArchiveValidMeasurementUnitConversion(ctx, &mealplanningsvc.ArchiveValidMeasurementUnitConversionRequest{
		ValidMeasurementUnitConversionId: conversionID,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error deleting conversion: %v", err)),
		), nil
	}

	// Return the updated conversions list
	return s.ValidMeasurementUnitConversionsForUnit(nil, req)
}
