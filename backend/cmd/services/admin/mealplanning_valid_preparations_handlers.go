package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	uploadedmediagrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"

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
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
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
	preparationID := createRes.Result.Id
	http.Redirect(res, req, fmt.Sprintf("/valid_preparations/%s", preparationID), http.StatusSeeOther)

	return g.El("div"), nil
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

	validPreparationRes, err := c.GetValidPreparation(ctx, &mealplanningsvc.GetValidPreparationRequest{ValidPreparationId: validPreparationID})
	if err != nil {
		return page("Valid Preparations", s.renderValidPreparationsError(fmt.Sprintf("Error loading valid preparation: %v", err))), nil
	}

	if validPreparationRes == nil || validPreparationRes.Result == nil {
		return page("Valid Preparations", s.renderValidPreparationsError("Error: Valid preparation not found")), nil
	}

	validPreparation := validPreparationRes.Result

	// Fetch associations for this preparation
	instrumentsAssociations, err := s.ValidPreparationInstrumentsForPreparation(nil, req)
	if err != nil {
		s.logger.Error("error fetching instrument associations", err)
	}

	ingredientsAssociations, err := s.ValidIngredientPreparationsForPreparation(nil, req)
	if err != nil {
		s.logger.Error("error fetching ingredient associations", err)
	}

	vesselsAssociations, err := s.ValidPreparationVesselsForPreparation(nil, req)
	if err != nil {
		s.logger.Error("error fetching vessel associations", err)
	}

	// Fetch ingredient associations for the "for ingredient" dropdown in media upload
	var ingredientOptions []*mealplanningsvc.ValidIngredient
	if ingRes, err := c.GetValidIngredientPreparationsByPreparation(ctx, &mealplanningsvc.GetValidIngredientPreparationsByPreparationRequest{
		ValidPreparationId: validPreparationID,
	}); err == nil && ingRes != nil {
		for _, assoc := range ingRes.Results {
			if assoc != nil && assoc.Ingredient != nil {
				ingredientOptions = append(ingredientOptions, assoc.Ingredient)
			}
		}
	}

	mediaSection := s.renderPreparationMediaSection(validPreparation, validPreparationID, ingredientOptions)

	// Use the FormPage component for viewing valid preparation data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidPreparation]{
		Title:        "Valid Preparation Details",
		BaseSubtitle: "View valid preparation information",
		Palette:      &design.StandardPalette,
		Data:         validPreparation,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidPreparation]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-preparation-form",
			Action:  fmt.Sprintf("/api/valid_preparations/%s", validPreparation.Id),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{
				"Name", "Description", "PastTense",
				// Preparation properties
				"RestrictToIngredients", "TemperatureRequired", "TimeEstimateRequired",
				"ConditionExpressionRequired", "ConsumesVessel", "OnlyForVessels", "YieldsNothing",
			},

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
				// Preparation property flags
				"RestrictToIngredients":       {InputType: "checkbox"},
				"TemperatureRequired":         {InputType: "checkbox"},
				"TimeEstimateRequired":        {InputType: "checkbox"},
				"ConditionExpressionRequired": {InputType: "checkbox"},
				"ConsumesVessel":              {InputType: "checkbox"},
				"OnlyForVessels":              {InputType: "checkbox"},
				"YieldsNothing":               {InputType: "checkbox"},
			},

			FormRows: []*components.FormRow{
				{
					Fields:  []string{"Name", "PastTense"},
					Columns: 2,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
				{
					Fields:  []string{"RestrictToIngredients", "TemperatureRequired", "TimeEstimateRequired"},
					Columns: 3,
				},
				{
					Fields:  []string{"ConditionExpressionRequired", "ConsumesVessel", "OnlyForVessels"},
					Columns: 3,
				},
				{
					Fields:  []string{"YieldsNothing"},
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
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Preparations", URL: "/valid_preparations"},
			{Text: validPreparation.Name, URL: ""},
		},

		// Dynamic subtitle showing preparation info
		SubtitleGenerator: func(vp *mealplanningsvc.ValidPreparation) string {
			return fmt.Sprintf("Viewing preparation: %s", vp.Name)
		},

		// Additional content - associations and media
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-6 mt-6"),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Instruments",
					Subtitle: "Required or optional instruments",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, instrumentsAssociations)),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Ingredients",
					Subtitle: "Valid ingredients for this preparation",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, ingredientsAssociations)),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Vessels",
					Subtitle: "Valid vessels for this preparation",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, vesselsAssociations)),
			),
			mediaSection,
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
			return fmt.Sprintf("/valid_preparations/%s", data.Id)
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
			return fmt.Sprintf("/valid_preparations/%s", data.Id)
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

// renderPreparationMediaSection renders the Media section with upload form and list of media.
func (s *AdminFrontendServer) renderPreparationMediaSection(vp *mealplanningsvc.ValidPreparation, validPreparationID string, ingredientOptions []*mealplanningsvc.ValidIngredient) g.Node {
	uploadForm := ghtml.Form(
		ghtml.Class("space-y-3"),
		ghtml.Method("POST"),
		ghtml.Action(fmt.Sprintf("/api/valid_preparations/%s/media", validPreparationID)),
		ghtml.EncType("multipart/form-data"),

		ghtml.Div(
			ghtml.Class("flex flex-col gap-2"),
			ghtml.Label(
				ghtml.Class("text-sm font-medium"),
				g.Text("Upload image or video"),
			),
			ghtml.Input(
				ghtml.Type("file"),
				ghtml.Name("file"),
				ghtml.Class("block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"),
				g.Attr("accept", "image/*,video/mp4"),
			),
		),
		g.If(len(ingredientOptions) > 0,
			ghtml.Div(
				ghtml.Class("flex flex-col gap-2"),
				ghtml.Label(
					ghtml.Class("text-sm font-medium"),
					g.Text("For ingredient (optional)"),
				),
				ghtml.Select(
					ghtml.Name("for_ingredient_id"),
					ghtml.Class("block w-full rounded-md border border-gray-300 px-3 py-2 text-sm"),
					ghtml.Option(
						ghtml.Value(""),
						g.Text("— General (all ingredients) —"),
					),
					g.Group(g.Map(ingredientOptions, func(ing *mealplanningsvc.ValidIngredient) g.Node {
						if ing == nil {
							return g.El("")
						}
						return ghtml.Option(
							ghtml.Value(ing.Id),
							g.Text(ing.Name),
						)
					})),
				),
			),
		),
		ghtml.Button(
			ghtml.Type("submit"),
			ghtml.Class("inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"),
			g.Text("Upload"),
		),
	)

	var mediaList g.Node
	if len(vp.Media) == 0 {
		mediaList = ghtml.P(
			ghtml.Class("text-sm text-gray-500 py-4"),
			g.Text("No media uploaded yet."),
		)
	} else {
		mediaList = ghtml.Div(
			ghtml.Class("space-y-2"),
			g.Group(g.Map(vp.Media, func(m *uploadedmediagrpc.UploadedMedia) g.Node {
				if m == nil {
					return g.El("")
				}
				return ghtml.Div(
					ghtml.Class("flex items-center gap-2 py-2 border-b border-gray-100 last:border-0"),
					ghtml.Span(
						ghtml.Class("text-sm font-mono text-gray-600"),
						g.Text(m.Id),
					),
					ghtml.Span(
						ghtml.Class("text-xs text-gray-400"),
						g.Text(m.MimeType.String()),
					),
				)
			})),
		)
	}

	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Media",
		Subtitle: "Images and videos for this preparation",
		Palette:  &design.StandardPalette,
	}, components.Card(&design.StandardPalette, ghtml.Div(
		ghtml.Class("space-y-4"),
		uploadForm,
		ghtml.Div(
			ghtml.Class("mt-4"),
			ghtml.H4(ghtml.Class("text-sm font-medium mb-2"), g.Text("Uploaded media")),
			mediaList,
		),
	)))
}

// renderValidPreparationsError creates a consistent error display for the valid preparations page.
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
