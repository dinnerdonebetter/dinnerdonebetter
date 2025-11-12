package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	validPreparationIDURLParamKey = "validPreparationID"
)

func (s *AdminFrontendServer) ValidPreparationCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Preparation", s.renderValidPreparationsError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidPreparationCreationRequestInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Preparation", s.renderValidPreparationsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid preparation
	createRes, err := c.CreateValidPreparation(ctx, &mealplanningsvc.CreateValidPreparationRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Preparation", s.renderValidPreparationsError(fmt.Sprintf("Error creating valid preparation: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Preparation", s.renderValidPreparationsError("Error: No valid preparation returned from server")), nil
	}

	// Redirect to the newly created valid preparation's page
	preparationID := createRes.Result.ID
	http.Redirect(res, req, fmt.Sprintf("/valid_preparations/%s", preparationID), http.StatusSeeOther)

	return nil, nil
}

func (s *AdminFrontendServer) ValidPreparationNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Preparation", s.renderValidPreparationsError("Error: No API client available")), nil
	}

	// Create an empty ValidPreparationCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidPreparationCreationRequestInput{}

	// Use the FormPage component for creating a new valid preparation
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidPreparationCreationRequestInput]{
		Title:        "Create New Valid Preparation",
		BaseSubtitle: "Add a new valid preparation method",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidPreparationCreationRequestInput]{
			FormID: "create-valid-preparation-form",
			Action: "/api/valid_preparations",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PastTense",
				"Description",
				"Slug",
				"IconPath",
				"TemperatureRequired",
				"TimeEstimateRequired",
				"ConditionExpressionRequired",
				"ConsumesVessel",
				"OnlyForVessels",
				"RestrictToIngredients",
				"YieldsNothing",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter preparation name...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PastTense": {
					Placeholder: "Enter past tense form (e.g., 'chopped' for 'chop')...",
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

			SubmitButtonText: "Create Valid Preparation",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_preparations",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Preparation", s.renderValidPreparationsError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Preparation", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidPreparationPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError("Error: No API client available")), nil
	}

	validPreparationID := s.validPreparationIDRouteParamFetcher(req)
	if validPreparationID == "" {
		return page("Valid Preparations", s.renderValidPreparationsError("Error: No valid preparation ID provided")), nil
	}

	validPreparationRes, err := c.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationID: validPreparationID})
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError(fmt.Sprintf("Error loading valid preparation: %v", err))), nil
	}

	if validPreparationRes == nil || validPreparationRes.Result == nil {
		return page("Valid Preparations", s.renderValidPreparationsError("Error: Valid preparation not found")), nil
	}

	validPreparation := validPreparationRes.Result

	// Use the FormPage component for viewing valid preparation data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidPreparation]{
		Title:        "Valid Preparation Details",
		BaseSubtitle: "View valid preparation information",
		Palette:      &design.StandardPalette,
		Data:         validPreparation,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidPreparation]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-preparation-form",
			Action:  fmt.Sprintf("/api/valid_preparations/%s", validPreparation.ID),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{"Name", "Description", "PastTense"},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the preparation...",
					InputType:   "textarea",
				},
				"PastTense": {
					Placeholder: "Past tense form (e.g., chopped, diced, sautéed)",
				},
			},

			FormRows: []components.FormRow{
				{
					Fields:  []string{"Name", "PastTense"},
					Columns: 2,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Preparations",
			CancelURL:        "/valid_preparations",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Preparations", URL: "/valid_preparations"},
			{Text: validPreparation.Name, URL: ""},
		},

		// Dynamic subtitle showing preparation info
		SubtitleGenerator: func(vp *mealplanningsvc.ValidPreparation) string {
			return fmt.Sprintf("Viewing preparation: %s", vp.Name)
		},

		// Additional info section showing preparation properties
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("mt-6 space-y-6"),
				// Preparation Properties
				components.Card(&design.StandardPalette,
					ghtml.H3(
						ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
						g.Text("Preparation Properties"),
					),
					ghtml.Div(
						ghtml.Class("grid grid-cols-2 md:grid-cols-3 gap-4"),
						propertyBadge("Restrict to Ingredients", validPreparation.RestrictToIngredients, &design.StandardPalette),
						propertyBadge("Temperature Required", validPreparation.TemperatureRequired, &design.StandardPalette),
						propertyBadge("Time Estimate Required", validPreparation.TimeEstimateRequired, &design.StandardPalette),
						propertyBadge("Condition Expression Required", validPreparation.ConditionExpressionRequired, &design.StandardPalette),
						propertyBadge("Consumes Vessel", validPreparation.ConsumesVessel, &design.StandardPalette),
						propertyBadge("Only for Vessels", validPreparation.OnlyForVessels, &design.StandardPalette),
						propertyBadge("Yields Nothing", validPreparation.YieldsNothing, &design.StandardPalette),
					),
				),
				// Count Requirements
				components.Card(&design.StandardPalette,
					ghtml.H3(
						ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
						g.Text("Count Requirements"),
					),
					ghtml.Div(
						ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
						rangeInfo("Ingredient Count", validPreparation.IngredientCount, &design.StandardPalette),
						rangeInfo("Instrument Count", validPreparation.InstrumentCount, &design.StandardPalette),
						rangeInfo("Vessel Count", validPreparation.VesselCount, &design.StandardPalette),
					),
				),
				// Associated Instruments
				components.Card(&design.StandardPalette,
					ghtml.Div(
						g.Attr("hx-get", fmt.Sprintf("/api/valid_preparations/%s/instruments", validPreparation.ID)),
						g.Attr("hx-trigger", "load"),
						g.Attr("hx-swap", "innerHTML"),
						ghtml.Class("min-h-24"),
						ghtml.P(
							ghtml.Class("text-gray-500 text-center py-8"),
							g.Text("Loading associated instruments..."),
						),
					),
				),
			),
		},
	})
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Preparations", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidPreparationsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError("Error: No API client available")), nil
	}

	validPreparationsRes, err := c.GetValidPreparations(ctx, &mealplanningsvc.GetValidPreparationsRequest{})
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError(fmt.Sprintf("Error loading valid preparations: %v", err))), nil
	}

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidPreparation]{
		Title:             "Valid Preparations",
		BaseSubtitle:      "Manage preparation definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search preparations...",
		HTMXSearchTarget:  "/api/valid_preparations/search",
		Data:              validPreparationsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Preparation", "/valid_preparations/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidPreparation]{
			TableID: "valid-preparations-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Name",
				"PastTense",
				"Description",
				"TemperatureRequired",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"PastTense":           "Past Tense",
				"TemperatureRequired": "Temp Required?",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
				"TemperatureRequired": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidPreparation) string {
			return fmt.Sprintf("/valid_preparations/%s", data.ID)
		},
		EmptyStateTitle:       "No valid preparations found",
		EmptyStateDescription: "No valid preparations have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Preparation", "/valid_preparations/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage preparation definitions"
			}
			return fmt.Sprintf("Manage %d preparation definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Preparations", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidPreparationsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validPreparationsRes, err := c.GetValidPreparations(ctx, &mealplanningsvc.GetValidPreparationsRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid preparations: %v", err)),
			),
		), nil
	}

	// Filter preparations based on search query
	var filteredPreparations []*mealplanningsvc.ValidPreparation
	if searchQuery == "" {
		filteredPreparations = validPreparationsRes.Results
	} else {
		// Filter preparations by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, preparation := range validPreparationsRes.Results {
			if strings.Contains(strings.ToLower(preparation.Name), searchQueryLower) ||
				strings.Contains(strings.ToLower(preparation.PastTense), searchQueryLower) ||
				strings.Contains(strings.ToLower(preparation.Description), searchQueryLower) {
				filteredPreparations = append(filteredPreparations, preparation)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredPreparations) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid preparations found",
				fmt.Sprintf("No valid preparations match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(filteredPreparations, &components.TableOptions[*mealplanningsvc.ValidPreparation]{
		TableID: "valid-preparations-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Name",
			"PastTense",
			"Description",
			"TemperatureRequired",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"PastTense":           "Past Tense",
			"TemperatureRequired": "Temp Required?",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
			"TemperatureRequired": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidPreparation) string {
			return fmt.Sprintf("/valid_preparations/%s", data.ID)
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

// renderValidPreparationsError creates a consistent error display for the valid preparations page
func (s *AdminFrontendServer) renderValidPreparationsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Preparations",
		Subtitle: "Manage preparation definitions",
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

// rangeInfo creates an info display for count range requirements
func rangeInfo(label string, rangeData *types.Uint16RangeWithOptionalMax, palette *design.Palette) g.Node {
	valueText := "-"
	if rangeData != nil {
		if rangeData.Max != nil && *rangeData.Max > 0 {
			valueText = fmt.Sprintf("%d - %d", rangeData.Min, *rangeData.Max)
		} else {
			valueText = fmt.Sprintf("%d+", rangeData.Min)
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
