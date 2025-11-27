package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	validIngredientStateIDURLParamKey = "validIngredientStateID"
)

func (s *AdminFrontendServer) ValidIngredientStateCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidIngredientStateCreationRequestInput
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid ingredient state
	createRes, err := c.CreateValidIngredientState(ctx, &mealplanningsvc.CreateValidIngredientStateRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError(fmt.Sprintf("Error creating valid ingredient state: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError("Error: No valid ingredient state returned from server")), nil
	}

	// Redirect to the newly created valid ingredient state's page
	ingredientStateID := createRes.Result.ID
	http.Redirect(res, req, fmt.Sprintf("/valid_ingredient_states/%s", ingredientStateID), http.StatusSeeOther)

	return nil, nil
}

func (s *AdminFrontendServer) ValidIngredientStateNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError("Error: No API client available")), nil
	}

	// Create an empty ValidIngredientStateCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidIngredientStateCreationRequestInput{}

	// Use the FormPage component for creating a new valid ingredient state
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidIngredientStateCreationRequestInput]{
		Title:        "Create New Valid Ingredient State",
		BaseSubtitle: "Add a new valid ingredient state",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidIngredientStateCreationRequestInput]{
			FormID: "create-valid-ingredient-state-form",
			Action: "/api/valid_ingredient_states",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PastTense",
				"Slug",
				"Description",
				"AttributeType",
				"IconPath",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter ingredient state name (e.g., 'diced', 'minced', 'whole')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PastTense": {
					Placeholder: "Enter past tense form (e.g., 'diced', 'minced', 'whole')...",
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
				"Description": {
					Placeholder: "Enter description...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"AttributeType": {
					Placeholder: "Enter attribute type (e.g., 'texture', 'form', 'temperature')...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"IconPath": {
					Placeholder: "Enter icon path (optional)...",
				},
			},

			SubmitButtonText: "Create Valid Ingredient State",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_ingredient_states",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Ingredient State", s.renderValidIngredientStatesError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Ingredient State", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientStatePage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError("Error: No API client available")), nil
	}

	validIngredientStateID := s.validIngredientStateIDRouteParamFetcher(req)
	if validIngredientStateID == "" {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError("Error: No valid ingredient state ID provided")), nil
	}

	validIngredientStateRes, err := c.GetValidIngredientState(ctx, &mealplanningsvc.GetValidIngredientStateRequest{ValidIngredientStateID: validIngredientStateID})
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError(fmt.Sprintf("Error loading valid ingredient state: %v", err))), nil
	}

	if validIngredientStateRes == nil || validIngredientStateRes.Result == nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError("Error: Valid ingredient state not found")), nil
	}

	validIngredientState := validIngredientStateRes.Result

	// Parse filter from request
	queryFilter := filtering.ExtractQueryFilterFromRequest(req)
	grpcFilter := grpcconverters.ConvertQueryFilterToGRPCQueryFilter(queryFilter, filtering.Pagination{})

	// Fetch ingredients for this ingredient state with filter
	ingredientsRes, err := c.GetValidIngredientStateIngredientsByIngredientState(ctx, &mealplanningsvc.GetValidIngredientStateIngredientsByIngredientStateRequest{
		ValidIngredientStateID: validIngredientStateID,
		Filter:                 grpcFilter,
	})
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError(fmt.Sprintf("Error loading ingredients: %v", err))), nil
	}

	// Use the ValidIngredientStateIngredient results directly for the table
	stateIngredients := []*mealplanningsvc.ValidIngredientStateIngredient{}
	if ingredientsRes != nil && ingredientsRes.Results != nil {
		stateIngredients = ingredientsRes.Results
	}

	// Use the FormPage component for viewing valid ingredient state data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidIngredientState]{
		Title:        "Valid Ingredient State Details",
		BaseSubtitle: "View valid ingredient state information",
		Palette:      &design.StandardPalette,
		Data:         validIngredientState,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidIngredientState]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-ingredient-state-form",
			Action:  fmt.Sprintf("/api/valid_ingredient_states/%s", validIngredientState.ID),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{"Name", "Description", "PastTense", "AttributeType"},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the ingredient state...",
					InputType:   "textarea",
				},
				"PastTense": {
					Placeholder: "Past tense form (e.g., chopped, diced)",
				},
				"AttributeType": {
					Placeholder: "Type: texture, consistency, temperature, color, etc.",
				},
			},

			FormRows: []components.FormRow{
				{
					Fields:  []string{"Name", "PastTense"},
					Columns: 2,
				},
				{
					Fields:  []string{"AttributeType"},
					Columns: 1,
				},
				{
					Fields:  []string{"Description"},
					Columns: 1,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Ingredient States",
			CancelURL:        "/valid_ingredient_states",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Ingredient States", URL: "/valid_ingredient_states"},
			{Text: validIngredientState.Name, URL: ""},
		},

		// Dynamic subtitle showing ingredient state info
		SubtitleGenerator: func(vis *mealplanningsvc.ValidIngredientState) string {
			return fmt.Sprintf("Viewing ingredient state: %s", vis.Name)
		},

		// Additional info section showing ingredients in a paginated table
		AdditionalContent: []g.Node{
			renderIngredientsTableForState(validIngredientStateID, stateIngredients, ingredientsRes, &design.StandardPalette),
		},
	})
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Ingredient States", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientStatesList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError("Error: No API client available")), nil
	}

	// Extract QueryFilter and convert to gRPC filter
	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)

	validIngredientStatesRes, err := c.GetValidIngredientStates(ctx, &mealplanningsvc.GetValidIngredientStatesRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError(fmt.Sprintf("Error loading valid ingredient states: %v", err))), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validIngredientStatesRes.Pagination)
	// Use search endpoint for pagination to return just the table, not the full page
	paginationURLGenerator := buildPaginationURLGenerator(req, "/api/valid_ingredient_states/search", queryFilter)

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidIngredientState]{
		Title:             "Valid Ingredient States",
		BaseSubtitle:      "Manage ingredient state definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search ingredient states...",
		HTMXSearchTarget:  "/api/valid_ingredient_states/search",
		Data:              validIngredientStatesRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Ingredient State", "/valid_ingredient_states/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidIngredientState]{
			TableID: "valid-ingredient-states-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Name",
				"PastTense",
				"AttributeType",
				"Description",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"PastTense":     "Past Tense",
				"AttributeType": "Attribute Type",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
			},
			Pagination:             pagination,
			PaginationURLGenerator: paginationURLGenerator,
			PaginationHTMXTarget:   "#search-results",
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredientState) string {
			return fmt.Sprintf("/valid_ingredient_states/%s", data.ID)
		},
		EmptyStateTitle:       "No valid ingredient states found",
		EmptyStateDescription: "No valid ingredient states have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Ingredient State", "/valid_ingredient_states/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage ingredient state definitions"
			}
			return fmt.Sprintf("Manage %d ingredient state definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Ingredient States", s.renderValidIngredientStatesError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Ingredient States", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientStatesSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Extract QueryFilter and convert to gRPC filter
	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)

	validIngredientStatesRes, err := c.GetValidIngredientStates(ctx, &mealplanningsvc.GetValidIngredientStatesRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid ingredient states: %v", err)),
			),
		), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validIngredientStatesRes.Pagination)
	paginationURLGenerator := buildPaginationURLGenerator(req, "/api/valid_ingredient_states/search", queryFilter)

	// Generate just the table (not the full page)
	if len(validIngredientStatesRes.Results) == 0 {
		searchQuery := req.URL.Query().Get("search")
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid ingredient states found",
				fmt.Sprintf("No valid ingredient states match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(validIngredientStatesRes.Results, &components.TableOptions[*mealplanningsvc.ValidIngredientState]{
		TableID: "valid-ingredient-states-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Name",
			"PastTense",
			"AttributeType",
			"Description",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"PastTense":     "Past Tense",
			"AttributeType": "Attribute Type",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
		},
		Pagination:             pagination,
		PaginationURLGenerator: paginationURLGenerator,
		PaginationHTMXTarget:   "#search-results",
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredientState) string {
			return fmt.Sprintf("/valid_ingredient_states/%s", data.ID)
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

// renderValidIngredientStatesError creates a consistent error display for the valid ingredient states page
func (s *AdminFrontendServer) renderValidIngredientStatesError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Ingredient States",
		Subtitle: "Manage ingredient state definitions",
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

// renderIngredientsTableForState creates a paginated table displaying the ingredients that this state is valid for
func renderIngredientsTableForState(
	ingredientStateID string,
	stateIngredients []*mealplanningsvc.ValidIngredientStateIngredient,
	response *mealplanningsvc.GetValidIngredientStateIngredientsByIngredientStateResponse,
	palette *design.Palette,
) g.Node {
	// Build search URL for this ingredient state
	searchURL := fmt.Sprintf("/api/valid_ingredient_states/%s/ingredients/search", ingredientStateID)

	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidIngredientStateIngredient]{
		Title:             "Valid For Ingredients",
		BaseSubtitle:      "Ingredients that can use this state",
		Palette:           palette,
		ShowSearch:        true,
		SearchPlaceholder: "Search ingredients...",
		HTMXSearchTarget:  searchURL,
		Data:              stateIngredients,
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidIngredientStateIngredient]{
			TableID: "ingredient-state-ingredients-table",
			Palette: palette,
			Fields: []string{
				"IngredientName",
				"IngredientDescription",
				"Notes",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"IngredientName":        "Ingredient",
				"IngredientDescription": "Description",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"IngredientName": func(value any) g.Node {
					// Value will be the whole item for custom fields, extract ingredient name
					if stateIngredient, ok := value.(*mealplanningsvc.ValidIngredientStateIngredient); ok && stateIngredient != nil && stateIngredient.Ingredient != nil {
						ingredientURL := fmt.Sprintf("/valid_ingredients/%s", stateIngredient.Ingredient.ID)
						return ghtml.A(
							ghtml.Href(ingredientURL),
							g.Attr("hx-get", ingredientURL),
							g.Attr("hx-target", "body"),
							g.Attr("hx-swap", "innerHTML"),
							g.Attr("hx-push-url", "true"),
							ghtml.Class(fmt.Sprintf("font-medium %s hover:underline", design.TextColor(palette.Primary))),
							g.Text(stateIngredient.Ingredient.Name),
						)
					}
					return g.Text("-")
				},
				"IngredientDescription": func(value any) g.Node {
					if stateIngredient, ok := value.(*mealplanningsvc.ValidIngredientStateIngredient); ok && stateIngredient != nil && stateIngredient.Ingredient != nil {
						desc := stateIngredient.Ingredient.Description
						if desc == "" {
							return g.Text("-")
						}
						return ghtml.Span(
							ghtml.Class(fmt.Sprintf("text-sm %s opacity-75", design.TextColor(palette.Text))),
							g.Text(desc),
						)
					}
					return g.Text("-")
				},
				"CreatedAt": renderTimestamp,
			},
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredientStateIngredient) string {
			if data != nil && data.Ingredient != nil {
				return fmt.Sprintf("/valid_ingredients/%s", data.Ingredient.ID)
			}
			return ""
		},
		EmptyStateTitle:       "No ingredients found",
		EmptyStateDescription: "No ingredients are associated with this ingredient state.",
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Ingredients that can use this state"
			}
			return fmt.Sprintf("Showing %d ingredient(s)", metadata.TotalCount)
		},
	})

	if err != nil {
		return ghtml.Div(
			ghtml.Class("mt-6"),
			components.Card(palette,
				ghtml.P(
					ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(palette.Warning))),
					g.Text(fmt.Sprintf("Error creating table: %v", err)),
				),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("mt-6"),
		tablePageResult.Node,
	)
}

func (s *AdminFrontendServer) ValidIngredientStateIngredientsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validIngredientStateID := s.validIngredientStateIDRouteParamFetcher(req)
	if validIngredientStateID == "" {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No valid ingredient state ID provided"),
			),
		), nil
	}

	// Parse filter from request
	queryFilter := filtering.ExtractQueryFilterFromRequest(req)
	grpcFilter := grpcconverters.ConvertQueryFilterToGRPCQueryFilter(queryFilter, filtering.Pagination{})

	// Get search query from request and add to filter
	searchQuery := req.URL.Query().Get("search")
	if searchQuery != "" && queryFilter != nil {
		grpcFilter = grpcconverters.ConvertQueryFilterToGRPCQueryFilter(queryFilter, filtering.Pagination{})
	}

	ingredientsRes, err := c.GetValidIngredientStateIngredientsByIngredientState(ctx, &mealplanningsvc.GetValidIngredientStateIngredientsByIngredientStateRequest{
		ValidIngredientStateID: validIngredientStateID,
		Filter:                 grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading ingredients: %v", err)),
			),
		), nil
	}

	stateIngredients := []*mealplanningsvc.ValidIngredientStateIngredient{}
	if ingredientsRes != nil && ingredientsRes.Results != nil {
		stateIngredients = ingredientsRes.Results
	}

	// Filter by search query if provided (client-side filtering for now)
	var filteredIngredients []*mealplanningsvc.ValidIngredientStateIngredient
	if searchQuery == "" {
		filteredIngredients = stateIngredients
	} else {
		searchQueryLower := strings.ToLower(searchQuery)
		for _, stateIngredient := range stateIngredients {
			if stateIngredient != nil && stateIngredient.Ingredient != nil {
				if strings.Contains(strings.ToLower(stateIngredient.Ingredient.Name), searchQueryLower) ||
					strings.Contains(strings.ToLower(stateIngredient.Ingredient.Description), searchQueryLower) ||
					strings.Contains(strings.ToLower(stateIngredient.Notes), searchQueryLower) {
					filteredIngredients = append(filteredIngredients, stateIngredient)
				}
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredIngredients) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No ingredients found",
				fmt.Sprintf("No ingredients match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(filteredIngredients, &components.TableOptions[*mealplanningsvc.ValidIngredientStateIngredient]{
		TableID: "ingredient-state-ingredients-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"IngredientName",
			"IngredientDescription",
			"Notes",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"IngredientName":        "Ingredient",
			"IngredientDescription": "Description",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"IngredientName": func(value any) g.Node {
				if stateIngredient, ok := value.(*mealplanningsvc.ValidIngredientStateIngredient); ok && stateIngredient != nil && stateIngredient.Ingredient != nil {
					ingredientURL := fmt.Sprintf("/valid_ingredients/%s", stateIngredient.Ingredient.ID)
					return ghtml.A(
						ghtml.Href(ingredientURL),
						g.Attr("hx-get", ingredientURL),
						g.Attr("hx-target", "body"),
						g.Attr("hx-swap", "innerHTML"),
						g.Attr("hx-push-url", "true"),
						ghtml.Class(fmt.Sprintf("font-medium %s hover:underline", design.TextColor(design.StandardPalette.Primary))),
						g.Text(stateIngredient.Ingredient.Name),
					)
				}
				return g.Text("-")
			},
			"IngredientDescription": func(value any) g.Node {
				if stateIngredient, ok := value.(*mealplanningsvc.ValidIngredientStateIngredient); ok && stateIngredient != nil && stateIngredient.Ingredient != nil {
					desc := stateIngredient.Ingredient.Description
					if desc == "" {
						return g.Text("-")
					}
					return ghtml.Span(
						ghtml.Class(fmt.Sprintf("text-sm %s opacity-75", design.TextColor(design.StandardPalette.Text))),
						g.Text(desc),
					)
				}
				return g.Text("-")
			},
			"CreatedAt": renderTimestamp,
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredientStateIngredient) string {
			if data != nil && data.Ingredient != nil {
				return fmt.Sprintf("/valid_ingredients/%s", data.Ingredient.ID)
			}
			return ""
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
