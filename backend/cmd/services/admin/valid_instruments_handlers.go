package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	validInstrumentIDURLParamKey = "validInstrumentID"
)

func (s *AdminFrontendServer) ValidInstrumentCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidInstrumentCreationRequestInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid instrument
	createRes, err := c.CreateValidInstrument(ctx, &mealplanningsvc.CreateValidInstrumentRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError(fmt.Sprintf("Error creating valid instrument: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError("Error: No valid instrument returned from server")), nil
	}

	// Redirect to the newly created valid instrument's page
	instrumentID := createRes.Result.ID
	http.Redirect(res, req, fmt.Sprintf("/valid_instruments/%s", instrumentID), http.StatusSeeOther)

	return nil, nil
}

func (s *AdminFrontendServer) ValidInstrumentNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError("Error: No API client available")), nil
	}

	// Create an empty ValidInstrumentCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidInstrumentCreationRequestInput{}

	// Use the FormPage component for creating a new valid instrument
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidInstrumentCreationRequestInput]{
		Title:        "Create New Valid Instrument",
		BaseSubtitle: "Add a new valid instrument",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidInstrumentCreationRequestInput]{
			FormID: "create-valid-instrument-form",
			Action: "/api/valid_instruments",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PluralName",
				"Description",
				"Slug",
				"IconPath",
				"DisplayInSummaryLists",
				"IncludeInGeneratedInstructions",
				"UsableForStorage",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter instrument name (e.g., 'knife', 'whisk', 'spatula')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PluralName": {
					Placeholder: "Enter plural name (e.g., 'knives', 'whisks', 'spatulas')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Enter description...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Slug": {
					Placeholder: "Enter URL-friendly slug...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"IconPath": {
					Placeholder: "Enter icon path (optional)...",
				},
			},

			SubmitButtonText: "Create Valid Instrument",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_instruments",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Instrument", s.renderValidInstrumentsError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Instrument", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidInstrumentPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError("Error: No API client available")), nil
	}

	validInstrumentID := s.validInstrumentIDRouteParamFetcher(req)
	if validInstrumentID == "" {
		return page("Valid Instruments", s.renderValidInstrumentsError("Error: No valid instrument ID provided")), nil
	}

	validInstrumentRes, err := c.GetValidInstrument(ctx, &mealplanningsvc.GetValidInstrumentRequest{ValidInstrumentID: validInstrumentID})
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError(fmt.Sprintf("Error loading valid instrument: %v", err))), nil
	}

	if validInstrumentRes == nil || validInstrumentRes.Result == nil {
		return page("Valid Instruments", s.renderValidInstrumentsError("Error: Valid instrument not found")), nil
	}

	validInstrument := validInstrumentRes.Result

	// Use the FormPage component for viewing valid instrument data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidInstrument]{
		Title:        "Valid Instrument Details",
		BaseSubtitle: "View valid instrument information",
		Palette:      &design.StandardPalette,
		Data:         validInstrument,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidInstrument]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-instrument-form",
			Action:  fmt.Sprintf("/api/valid_instruments/%s", validInstrument.ID),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{"Name", "Description", "PluralName"},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the instrument...",
					InputType:   "textarea",
				},
				"PluralName": {
					Placeholder: "Plural form of the instrument name",
				},
			},

			FormRows: []components.FormRow{
				{
					Fields:  []string{"Name", "PluralName"},
					Columns: 2,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Instruments",
			CancelURL:        "/valid_instruments",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Instruments", URL: "/valid_instruments"},
			{Text: validInstrument.Name, URL: ""},
		},

		// Dynamic subtitle showing instrument info
		SubtitleGenerator: func(vi *mealplanningsvc.ValidInstrument) string {
			return fmt.Sprintf("Viewing instrument: %s", vi.Name)
		},

		// Additional info section showing boolean flags and associations
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("mt-6"),
				components.Card(&design.StandardPalette,
					ghtml.H3(
						ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
						g.Text("Instrument Properties"),
					),
					ghtml.Div(
						ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
						propertyBadge("Display in Summary Lists", validInstrument.DisplayInSummaryLists, &design.StandardPalette),
						propertyBadge("Include in Generated Instructions", validInstrument.IncludeInGeneratedInstructions, &design.StandardPalette),
						propertyBadge("Usable for Storage", validInstrument.UsableForStorage, &design.StandardPalette),
					),
				),
			),
			ghtml.Div(
				ghtml.Class("mt-6"),
				components.Card(&design.StandardPalette,
					ghtml.Div(
						g.Attr("hx-get", fmt.Sprintf("/api/valid_instruments/%s/preparations", validInstrument.ID)),
						g.Attr("hx-trigger", "load"),
						g.Attr("hx-swap", "innerHTML"),
						ghtml.Class("min-h-24"),
						ghtml.P(
							ghtml.Class("text-gray-500 text-center py-8"),
							g.Text("Loading associated preparations..."),
						),
					),
				),
			),
		},
	})
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Instruments", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidInstrumentsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError("Error: No API client available")), nil
	}

	validInstrumentsRes, err := c.GetValidInstruments(ctx, &mealplanningsvc.GetValidInstrumentsRequest{})
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError(fmt.Sprintf("Error loading valid instruments: %v", err))), nil
	}

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidInstrument]{
		Title:             "Valid Instruments",
		BaseSubtitle:      "Manage instrument definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search instruments...",
		HTMXSearchTarget:  "/api/valid_instruments/search",
		Data:              validInstrumentsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Instrument", "/valid_instruments/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidInstrument]{
			TableID: "valid-instruments-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Name",
				"PluralName",
				"Description",
				"UsableForStorage",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"PluralName":       "Plural Name",
				"UsableForStorage": "Storage?",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
				"UsableForStorage": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidInstrument) string {
			return fmt.Sprintf("/valid_instruments/%s", data.ID)
		},
		EmptyStateTitle:       "No valid instruments found",
		EmptyStateDescription: "No valid instruments have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Instrument", "/valid_instruments/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage instrument definitions"
			}
			return fmt.Sprintf("Manage %d instrument definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Instruments", s.renderValidInstrumentsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Instruments", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidInstrumentsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	// Get search query from request
	searchQuery := req.URL.Query().Get("search")

	validInstrumentsRes, err := c.GetValidInstruments(ctx, &mealplanningsvc.GetValidInstrumentsRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid instruments: %v", err)),
			),
		), nil
	}

	// Filter instruments based on search query
	var filteredInstruments []*mealplanningsvc.ValidInstrument
	if searchQuery == "" {
		filteredInstruments = validInstrumentsRes.Results
	} else {
		// Filter instruments by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, instrument := range validInstrumentsRes.Results {
			if strings.Contains(strings.ToLower(instrument.Name), searchQueryLower) ||
				strings.Contains(strings.ToLower(instrument.PluralName), searchQueryLower) ||
				strings.Contains(strings.ToLower(instrument.Description), searchQueryLower) {
				filteredInstruments = append(filteredInstruments, instrument)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredInstruments) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid instruments found",
				fmt.Sprintf("No valid instruments match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(filteredInstruments, &components.TableOptions[*mealplanningsvc.ValidInstrument]{
		TableID: "valid-instruments-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Name",
			"PluralName",
			"Description",
			"UsableForStorage",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"PluralName":       "Plural Name",
			"UsableForStorage": "Storage?",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
			"UsableForStorage": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidInstrument) string {
			return fmt.Sprintf("/valid_instruments/%s", data.ID)
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error creating table: %v", err)),
			),
		), nil
	}

	// Wrap table in the same scrollable container structure for consistency
	return g.El("div",
		g.Attr("class", "overflow-x-auto"),
		table,
	), nil
}

// renderValidInstrumentsError creates a consistent error display for the valid instruments page
func (s *AdminFrontendServer) renderValidInstrumentsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Instruments",
		Subtitle: "Manage instrument definitions",
		Palette:  &design.StandardPalette,
	},
		components.Card(&design.StandardPalette,
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(errorMsg),
			),
		),
	)
}
