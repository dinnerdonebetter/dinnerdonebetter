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
	validMeasurementUnitIDURLParamKey = "validMeasurementUnitID"
)

func (s *AdminFrontendServer) ValidMeasurementUnitCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidMeasurementUnitCreationRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid measurement unit
	createRes, err := c.CreateValidMeasurementUnit(ctx, &mealplanningsvc.CreateValidMeasurementUnitRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error creating valid measurement unit: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError("Error: No valid measurement unit returned from server")), nil
	}

	// Redirect to the newly created valid measurement unit's page
	measurementUnitID := createRes.Result.ID
	http.Redirect(res, req, fmt.Sprintf("/valid_measurement_units/%s", measurementUnitID), http.StatusSeeOther)

	return g.El("div"), nil
}

func (s *AdminFrontendServer) ValidMeasurementUnitNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError("Error: No API client available")), nil
	}

	// Create an empty ValidMeasurementUnitCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidMeasurementUnitCreationRequestInput{}

	// Use the FormPage component for creating a new valid measurement unit
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidMeasurementUnitCreationRequestInput]{
		Title:        "Create New Valid Measurement Unit",
		BaseSubtitle: "Add a new valid measurement unit",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidMeasurementUnitCreationRequestInput]{
			FormID: "create-valid-measurement-unit-form",
			Action: "/api/valid_measurement_units",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PluralName",
				"Description",
				"Slug",
				"IconPath",
				"Volumetric",
				"Universal",
				"Metric",
				"Imperial",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter measurement unit name (e.g., 'cup', 'tablespoon', 'gram')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PluralName": {
					Placeholder: "Enter plural name (e.g., 'cups', 'tablespoons', 'grams')...",
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

			SubmitButtonText: "Create Valid Measurement Unit",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_measurement_units",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Measurement Unit", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Measurement Unit", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidMeasurementUnitPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError("Error: No API client available")), nil
	}

	validMeasurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)
	if validMeasurementUnitID == "" {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError("Error: No valid measurement unit ID provided")), nil
	}

	validMeasurementUnitRes, err := c.GetValidMeasurementUnit(ctx, &mealplanningsvc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: validMeasurementUnitID})
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error loading valid measurement unit: %v", err))), nil
	}

	if validMeasurementUnitRes == nil || validMeasurementUnitRes.Result == nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError("Error: Valid measurement unit not found")), nil
	}

	validMeasurementUnit := validMeasurementUnitRes.Result

	// Determine the measurement system value for the select element
	measurementSystem := ""
	if validMeasurementUnit.Metric {
		measurementSystem = "metric"
	} else if validMeasurementUnit.Imperial {
		measurementSystem = "imperial"
	}

	// Fetch associations for this measurement unit
	associationListNode, err := s.ValidIngredientMeasurementUnitsForMeasurementUnit(nil, req)
	if err != nil {
		s.logger.Error("error fetching ingredient associations", err)
	}

	// Fetch conversions for this measurement unit
	conversionsNode, err := s.ValidMeasurementUnitConversionsForUnit(nil, req)
	if err != nil {
		s.logger.Error("error fetching conversions", err)
	}

	// Use the FormPage component for viewing valid measurement unit data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidMeasurementUnit]{
		Title:        "Valid Measurement Unit Details",
		BaseSubtitle: "View valid measurement unit information",
		Palette:      &design.StandardPalette,
		Data:         validMeasurementUnit,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidMeasurementUnit]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-measurement-unit-form",
			Action:  fmt.Sprintf("/api/valid_measurement_units/%s", validMeasurementUnit.ID),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{"Name", "Description", "PluralName", "Volumetric", "Universal", "Metric"},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the measurement unit...",
					InputType:   "textarea",
				},
				"PluralName": {
					Placeholder: "Plural form of the measurement unit name",
				},
				"Volumetric": {
					InputType: "checkbox",
				},
				"Universal": {
					InputType: "checkbox",
				},
				"Metric": {
					DisplayName: "Measurement System",
					CustomRenderer: func(fieldName string, value any, config *components.FieldConfig, palette *design.Palette) g.Node {
						// Build the select element for measurement system
						// This select will update hidden Metric and Imperial fields via JavaScript
						selectAttrs := []g.Node{
							ghtml.ID("MeasurementSystem"),
							ghtml.Name("MeasurementSystem"),
							ghtml.Class(fmt.Sprintf("block w-full px-3 py-2 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-%s focus:border-%s sm:text-sm",
								palette.Primary.Value,
								palette.Primary.Value,
							)),
							// Update hidden Metric/Imperial fields when selection changes
							g.Attr("onchange", `
								const metric = document.getElementById('Metric');
								const imperial = document.getElementById('Imperial');
								if (this.value === 'metric') {
									metric.checked = true;
									imperial.checked = false;
								} else if (this.value === 'imperial') {
									metric.checked = false;
									imperial.checked = true;
								}
							`),
						}

						// Build options - only Metric and Imperial (no "Neither")
						options := []g.Node{ghtml.Option(
							ghtml.Value("metric"),
							g.Text("Metric"),
							g.If(measurementSystem == "metric" || measurementSystem == "", ghtml.Selected()), // Default to Metric if neither is set
						), ghtml.Option(
							ghtml.Value("imperial"),
							g.Text("Imperial"),
							g.If(measurementSystem == "imperial", ghtml.Selected()),
						),
						}
						return ghtml.Div(
							ghtml.Class("flex flex-col"),
							ghtml.Label(
								ghtml.For("MeasurementSystem"),
								ghtml.Class(fmt.Sprintf("block text-sm font-medium %s mb-2", design.TextColor(palette.Text))),
								g.Text("Measurement System"),
							),
							ghtml.Select(append(selectAttrs, g.Group(options))...),
							// Hidden fields for Metric and Imperial that will be submitted with the form
							// We use hidden inputs with value="false" followed by checkboxes with value="true"
							// This ensures that false values are sent when checkboxes are unchecked
							ghtml.Input(
								ghtml.Type("hidden"),
								ghtml.Name("Metric"),
								ghtml.Value("false"),
							),
							ghtml.Input(
								ghtml.Type("checkbox"),
								ghtml.ID("Metric"),
								ghtml.Name("Metric"),
								ghtml.Value("true"),
								ghtml.Class("hidden"),
								g.If(validMeasurementUnit.Metric, ghtml.Checked()),
							),
							ghtml.Input(
								ghtml.Type("hidden"),
								ghtml.Name("Imperial"),
								ghtml.Value("false"),
							),
							ghtml.Input(
								ghtml.Type("checkbox"),
								ghtml.ID("Imperial"),
								ghtml.Name("Imperial"),
								ghtml.Value("true"),
								ghtml.Class("hidden"),
								g.If(validMeasurementUnit.Imperial, ghtml.Checked()),
							),
						)
					},
					Enabled: true,
				},
			},

			FormRows: []*components.FormRow{
				{
					Fields:  []string{"Name", "PluralName"},
					Columns: 2,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
				{
					Fields:  []string{"Volumetric", "Universal", "Metric"},
					Columns: 3,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Measurement Units",
			CancelURL:        "/valid_measurement_units",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Measurement Units", URL: "/valid_measurement_units"},
			{Text: validMeasurementUnit.Name, URL: ""},
		},

		// Dynamic subtitle showing measurement unit info
		SubtitleGenerator: func(vmu *mealplanningsvc.ValidMeasurementUnit) string {
			return fmt.Sprintf("Viewing measurement unit: %s", vmu.Name)
		},

		// Additional content - associations and conversions
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("grid grid-cols-1 md:grid-cols-2 gap-6 mt-6"),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Ingredient Associations",
					Subtitle: "Ingredients that can be measured with this unit",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, associationListNode)),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Unit Conversions",
					Subtitle: "Conversion factors to/from other units",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette,
					ghtml.Div(
						ghtml.ID("conversions-container"),
						conversionsNode,
					),
				)),
			),
		},
	})
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Measurement Units", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidMeasurementUnitsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError("Error: No API client available")), nil
	}

	validMeasurementUnitsRes, err := c.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{})
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error loading valid measurement units: %v", err))), nil
	}

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidMeasurementUnit]{
		Title:             "Valid Measurement Units",
		BaseSubtitle:      "Manage measurement unit definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search measurement units...",
		HTMXSearchTarget:  "/api/valid_measurement_units/search",
		Data:              validMeasurementUnitsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Measurement Unit", "/valid_measurement_units/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidMeasurementUnit]{
			TableID: "valid-measurement-units-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Name",
				"PluralName",
				"Description",
				"Volumetric",
				"Universal",
				"Metric",
				"Imperial",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"PluralName": "Plural Name",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
				"Volumetric": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
				"Universal": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
				"Metric": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
				"Imperial": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidMeasurementUnit) string {
			return fmt.Sprintf("/valid_measurement_units/%s", data.ID)
		},
		EmptyStateTitle:       "No valid measurement units found",
		EmptyStateDescription: "No valid measurement units have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Measurement Unit", "/valid_measurement_units/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage measurement unit definitions"
			}
			return fmt.Sprintf("Manage %d measurement unit definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Measurement Units", s.renderValidMeasurementUnitsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Measurement Units", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidMeasurementUnitsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validMeasurementUnitsRes, err := c.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid measurement units: %v", err)),
			),
		), nil
	}

	// Filter measurement units based on search query
	var filteredUnits []*mealplanningsvc.ValidMeasurementUnit
	if searchQuery == "" {
		filteredUnits = validMeasurementUnitsRes.Results
	} else {
		// Filter measurement units by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, unit := range validMeasurementUnitsRes.Results {
			if strings.Contains(strings.ToLower(unit.Name), searchQueryLower) ||
				strings.Contains(strings.ToLower(unit.PluralName), searchQueryLower) ||
				strings.Contains(strings.ToLower(unit.Description), searchQueryLower) {
				filteredUnits = append(filteredUnits, unit)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredUnits) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid measurement units found",
				fmt.Sprintf("No valid measurement units match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(filteredUnits, &components.TableOptions[*mealplanningsvc.ValidMeasurementUnit]{
		TableID: "valid-measurement-units-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Name",
			"PluralName",
			"Description",
			"Volumetric",
			"Metric",
			"Imperial",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"PluralName": "Plural Name",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
			"Volumetric": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
			"Metric": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
			"Imperial": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidMeasurementUnit) string {
			return fmt.Sprintf("/valid_measurement_units/%s", data.ID)
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

// renderValidMeasurementUnitsError creates a consistent error display for the valid measurement units page.
func (s *AdminFrontendServer) renderValidMeasurementUnitsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Measurement Units",
		Subtitle: "Manage measurement unit definitions",
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
