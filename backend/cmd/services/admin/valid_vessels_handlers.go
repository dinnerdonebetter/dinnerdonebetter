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
	validVesselIDURLParamKey = "validVesselID"
)

func (s *AdminFrontendServer) ValidVesselCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Vessel", s.renderValidVesselsError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidVesselCreationRequestInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Vessel", s.renderValidVesselsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid vessel
	createRes, err := c.CreateValidVessel(ctx, &mealplanningsvc.CreateValidVesselRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Vessel", s.renderValidVesselsError(fmt.Sprintf("Error creating valid vessel: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Vessel", s.renderValidVesselsError("Error: No valid vessel returned from server")), nil
	}

	// Redirect to the newly created valid vessel's page
	vesselID := createRes.Result.ID
	http.Redirect(res, req, fmt.Sprintf("/valid_vessels/%s", vesselID), http.StatusSeeOther)

	return nil, nil
}

func (s *AdminFrontendServer) ValidVesselNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Vessel", s.renderValidVesselsError("Error: No API client available")), nil
	}

	// Create an empty ValidVesselCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidVesselCreationRequestInput{}

	// Use the FormPage component for creating a new valid vessel
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidVesselCreationRequestInput]{
		Title:        "Create New Valid Vessel",
		BaseSubtitle: "Add a new valid vessel",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidVesselCreationRequestInput]{
			FormID: "create-valid-vessel-form",
			Action: "/api/valid_vessels",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PluralName",
				"Description",
				"Slug",
				"Shape",
				"IconPath",
				"CapacityUnitID",
				"Capacity",
				"WidthInMillimeters",
				"LengthInMillimeters",
				"HeightInMillimeters",
				"UsableForStorage",
				"IncludeInGeneratedInstructions",
				"DisplayInSummaryLists",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter vessel name (e.g., 'pot', 'pan', 'bowl')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PluralName": {
					Placeholder: "Enter plural name (e.g., 'pots', 'pans', 'bowls')...",
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
				"Shape": {
					Placeholder: "Enter shape (e.g., 'round', 'rectangular', 'square')...",
				},
				"IconPath": {
					Placeholder: "Enter icon path (optional)...",
				},
				"CapacityUnitID": {
					Placeholder: "Enter capacity unit ID (optional)...",
				},
				"Capacity": {
					Placeholder: "Enter capacity (optional)...",
				},
				"WidthInMillimeters": {
					Placeholder: "Enter width in millimeters (optional)...",
				},
				"LengthInMillimeters": {
					Placeholder: "Enter length in millimeters (optional)...",
				},
				"HeightInMillimeters": {
					Placeholder: "Enter height in millimeters (optional)...",
				},
			},

			SubmitButtonText: "Create Valid Vessel",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_vessels",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Vessel", s.renderValidVesselsError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Vessel", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidVesselPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError("Error: No API client available")), nil
	}

	validVesselID := s.validVesselIDRouteParamFetcher(req)
	if validVesselID == "" {
		return page("Valid Vessels", s.renderValidVesselsError("Error: No valid vessel ID provided")), nil
	}

	validVesselRes, err := c.GetValidVessel(ctx, &mealplanningsvc.GetValidVesselRequest{ValidVesselID: validVesselID})
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError(fmt.Sprintf("Error loading valid vessel: %v", err))), nil
	}

	if validVesselRes == nil || validVesselRes.Result == nil {
		return page("Valid Vessels", s.renderValidVesselsError("Error: Valid vessel not found")), nil
	}

	validVessel := validVesselRes.Result

	// Use the FormPage component for viewing valid vessel data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidVessel]{
		Title:        "Valid Vessel Details",
		BaseSubtitle: "View valid vessel information",
		Palette:      &design.StandardPalette,
		Data:         validVessel,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidVessel]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-vessel-form",
			Action:  fmt.Sprintf("/api/valid_vessels/%s", validVessel.ID),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{"Name", "Description", "PluralName", "Shape"},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the vessel...",
					InputType:   "textarea",
				},
				"PluralName": {
					Placeholder: "Plural form of the vessel name",
				},
				"Shape": {
					Placeholder: "Shape of the vessel (e.g., round, rectangular)",
				},
			},

			FormRows: []components.FormRow{
				{
					Fields:  []string{"Name", "PluralName"},
					Columns: 2,
				},
				{
					Fields:  []string{"Shape"},
					Columns: 1,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Vessels",
			CancelURL:        "/valid_vessels",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Vessels", URL: "/valid_vessels"},
			{Text: validVessel.Name, URL: ""},
		},

		// Dynamic subtitle showing vessel info
		SubtitleGenerator: func(vv *mealplanningsvc.ValidVessel) string {
			return fmt.Sprintf("Viewing vessel: %s", vv.Name)
		},

		// Additional info section showing properties and dimensions
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("mt-6 space-y-6"),
				// Vessel Properties
				components.Card(&design.StandardPalette,
					ghtml.H3(
						ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
						g.Text("Vessel Properties"),
					),
					ghtml.Div(
						ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
						propertyBadge("Display in Summary Lists", validVessel.DisplayInSummaryLists, &design.StandardPalette),
						propertyBadge("Include in Generated Instructions", validVessel.IncludeInGeneratedInstructions, &design.StandardPalette),
						propertyBadge("Usable for Storage", validVessel.UsableForStorage, &design.StandardPalette),
					),
				),
				// Vessel Dimensions
				components.Card(&design.StandardPalette,
					ghtml.H3(
						ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
						g.Text("Dimensions & Capacity"),
					),
					ghtml.Div(
						ghtml.Class("grid grid-cols-2 md:grid-cols-4 gap-4"),
						dimensionInfo("Width", validVessel.WidthInMillimeters, "mm", &design.StandardPalette),
						dimensionInfo("Length", validVessel.LengthInMillimeters, "mm", &design.StandardPalette),
						dimensionInfo("Height", validVessel.HeightInMillimeters, "mm", &design.StandardPalette),
						capacityInfo("Capacity", validVessel.Capacity, validVessel.CapacityUnit, &design.StandardPalette),
					),
				),
			),
		},
	})
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Vessels", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidVesselsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError("Error: No API client available")), nil
	}

	validVesselsRes, err := c.GetValidVessels(ctx, &mealplanningsvc.GetValidVesselsRequest{})
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError(fmt.Sprintf("Error loading valid vessels: %v", err))), nil
	}

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidVessel]{
		Title:             "Valid Vessels",
		BaseSubtitle:      "Manage vessel definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search vessels...",
		HTMXSearchTarget:  "/api/valid_vessels/search",
		Data:              validVesselsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Vessel", "/valid_vessels/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidVessel]{
			TableID: "valid-vessels-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Name",
				"PluralName",
				"Shape",
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
		RowLinkGenerator: func(data *mealplanningsvc.ValidVessel) string {
			return fmt.Sprintf("/valid_vessels/%s", data.ID)
		},
		EmptyStateTitle:       "No valid vessels found",
		EmptyStateDescription: "No valid vessels have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Vessel", "/valid_vessels/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage vessel definitions"
			}
			return fmt.Sprintf("Manage %d vessel definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Vessels", s.renderValidVesselsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Vessels", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidVesselsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validVesselsRes, err := c.GetValidVessels(ctx, &mealplanningsvc.GetValidVesselsRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid vessels: %v", err)),
			),
		), nil
	}

	// Filter vessels based on search query
	var filteredVessels []*mealplanningsvc.ValidVessel
	if searchQuery == "" {
		filteredVessels = validVesselsRes.Results
	} else {
		// Filter vessels by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, vessel := range validVesselsRes.Results {
			if strings.Contains(strings.ToLower(vessel.Name), searchQueryLower) ||
				strings.Contains(strings.ToLower(vessel.PluralName), searchQueryLower) ||
				strings.Contains(strings.ToLower(vessel.Description), searchQueryLower) ||
				strings.Contains(strings.ToLower(vessel.Shape), searchQueryLower) {
				filteredVessels = append(filteredVessels, vessel)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredVessels) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid vessels found",
				fmt.Sprintf("No valid vessels match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(filteredVessels, &components.TableOptions[*mealplanningsvc.ValidVessel]{
		TableID: "valid-vessels-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Name",
			"PluralName",
			"Shape",
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
		RowLinkGenerator: func(data *mealplanningsvc.ValidVessel) string {
			return fmt.Sprintf("/valid_vessels/%s", data.ID)
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

// renderValidVesselsError creates a consistent error display for the valid vessels page
func (s *AdminFrontendServer) renderValidVesselsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Vessels",
		Subtitle: "Manage vessel definitions",
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

// dimensionInfo creates an info display for a dimension value
func dimensionInfo(label string, value float32, unit string, palette *design.Palette) g.Node {
	valueText := "-"
	if value > 0 {
		valueText = fmt.Sprintf("%.1f %s", value, unit)
	}

	return ghtml.Div(
		ghtml.Class(fmt.Sprintf("flex flex-col p-3 rounded-lg %s", design.Background(design.Color{Value: "gray-50"}))),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("text-xs font-medium %s opacity-75", design.TextColor(palette.Text))),
			g.Text(label),
		),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("mt-1 text-sm font-semibold %s", design.TextColor(palette.Primary))),
			g.Text(valueText),
		),
	)
}

// capacityInfo creates an info display for capacity with unit
func capacityInfo(label string, value float32, unit *mealplanningsvc.ValidMeasurementUnit, palette *design.Palette) g.Node {
	valueText := "-"
	if value > 0 {
		if unit != nil {
			valueText = fmt.Sprintf("%.2f %s", value, unit.Name)
		} else {
			valueText = fmt.Sprintf("%.2f", value)
		}
	}

	return ghtml.Div(
		ghtml.Class(fmt.Sprintf("flex flex-col p-3 rounded-lg %s", design.Background(design.Color{Value: "gray-50"}))),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("text-xs font-medium %s opacity-75", design.TextColor(palette.Text))),
			g.Text(label),
		),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("mt-1 text-sm font-semibold %s", design.TextColor(palette.Primary))),
			g.Text(valueText),
		),
	)
}

