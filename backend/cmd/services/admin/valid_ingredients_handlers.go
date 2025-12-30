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
	validIngredientIDURLParamKey = "validIngredientID"
)

func (s *AdminFrontendServer) ValidIngredientCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError("Error: No API client available")), nil
	}

	// Decode JSON request body
	var input *mealplanningsvc.ValidIngredientCreationRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	// Call gRPC service to create the valid ingredient
	createRes, err := c.CreateValidIngredient(ctx, &mealplanningsvc.CreateValidIngredientRequest{
		Input: input,
	})
	if err != nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError(fmt.Sprintf("Error creating valid ingredient: %v", err))), nil
	}

	if createRes == nil || createRes.Result == nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError("Error: No valid ingredient returned from server")), nil
	}

	// Redirect to the newly created valid ingredient's page
	ingredientID := createRes.Result.Id
	http.Redirect(res, req, fmt.Sprintf("/valid_ingredients/%s", ingredientID), http.StatusSeeOther)

	return g.El("div"), nil
}

func (s *AdminFrontendServer) ValidIngredientNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	_, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError("Error: No API client available")), nil
	}

	// Create an empty ValidIngredientCreationRequestInput for the form
	emptyInput := &mealplanningsvc.ValidIngredientCreationRequestInput{}

	// Use the FormPage component for creating a new valid ingredient
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidIngredientCreationRequestInput]{
		Title:        "Create New Valid Ingredient",
		BaseSubtitle: "Add a new valid ingredient",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidIngredientCreationRequestInput]{
			FormID: "create-valid-ingredient-form",
			Action: "/api/valid_ingredients",
			Method: "POST",

			// All fields should be enabled for creation
			EnabledFields: []string{
				"Name",
				"PluralName",
				"Description",
				"Slug",
				"IconPath",
				"Warning",
				"StorageInstructions",
				"ShoppingSuggestions",
				"IsLiquid",
				"ContainsPeanut",
				"ContainsSoy",
				"ContainsEgg",
				"ContainsWheat",
				"ContainsShellfish",
				"ContainsFish",
				"ContainsDairy",
				"ContainsSesame",
				"ContainsTreeNut",
				"ContainsAlcohol",
				"ContainsGluten",
				"AnimalFlesh",
				"AnimalDerived",
				"RestrictToPreparations",
				"IsStarch",
				"IsProtein",
				"IsGrain",
				"IsFruit",
				"IsSalt",
				"IsFat",
				"IsAcid",
				"IsHeat",
			},

			// Configure field validation
			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Placeholder: "Enter ingredient name...",
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"PluralName": {
					Placeholder: "Enter plural name...",
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
				"Warning": {
					Placeholder: "Enter any warnings (optional)...",
				},
				"StorageInstructions": {
					Placeholder: "Enter storage instructions (optional)...",
				},
				"ShoppingSuggestions": {
					Placeholder: "Enter shopping suggestions (optional)...",
				},
			},

			SubmitButtonText: "Create Valid Ingredient",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/valid_ingredients",

			// HTMX configuration
			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},
	})
	if err != nil {
		return page("New Valid Ingredient", s.renderValidIngredientsError(fmt.Sprintf("Error rendering form: %v", err))), nil
	}

	return page("Create Valid Ingredient", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError("Error: No API client available")), nil
	}

	validIngredientID := s.validIngredientIDRouteParamFetcher(req)
	if validIngredientID == "" {
		return page("Valid Ingredients", s.renderValidIngredientsError("Error: No valid ingredient MealPlanTaskID provided")), nil
	}

	validIngredientRes, err := c.GetValidIngredient(ctx, &mealplanningsvc.GetValidIngredientRequest{ValidIngredientId: validIngredientID})
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError(fmt.Sprintf("Error loading valid ingredient: %v", err))), nil
	}

	if validIngredientRes == nil || validIngredientRes.Result == nil {
		return page("Valid Ingredients", s.renderValidIngredientsError("Error: Valid ingredient not found")), nil
	}

	validIngredient := validIngredientRes.Result

	// Fetch associations for this ingredient
	measurementUnitsAssociations, err := s.ValidIngredientMeasurementUnitsForIngredient(nil, req)
	if err != nil {
		s.logger.Error("error fetching measurement unit associations", err)
	}

	preparationsAssociations, err := s.ValidIngredientPreparationsForIngredient(nil, req)
	if err != nil {
		s.logger.Error("error fetching preparation associations", err)
	}

	// Use the FormPage component for viewing valid ingredient data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidIngredient]{
		Title:        "Valid Ingredient Details",
		BaseSubtitle: "View valid ingredient information",
		Palette:      &design.StandardPalette,
		Data:         validIngredient,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidIngredient]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-ingredient-form",
			Action:  fmt.Sprintf("/api/valid_ingredients/%s", validIngredient.Id),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{
				"Name", "Description", "PluralName", "Warning", "StorageInstructions", "ShoppingSuggestions",
				// Allergen flags
				"ContainsEgg", "ContainsDairy", "ContainsFish", "ContainsShellfish",
				"ContainsPeanut", "ContainsTreeNut", "ContainsWheat", "ContainsSoy",
				"ContainsSesame", "ContainsGluten", "ContainsAlcohol",
				// Property flags
				"AnimalDerived", "AnimalFlesh", "IsLiquid", "IsStarch", "IsProtein",
				"IsGrain", "IsFruit", "IsSalt", "IsFat", "IsAcid", "IsHeat",
			},

			FieldConfigs: map[string]*components.FieldConfig{
				"Name": {
					Validation: &components.FieldValidation{
						Required: true,
					},
				},
				"Description": {
					Placeholder: "Description of the ingredient...",
					InputType:   "textarea",
				},
				"PluralName": {
					Placeholder: "Plural form of the ingredient name",
				},
				"Warning": {
					Placeholder: "Any warnings or cautions...",
					InputType:   "textarea",
				},
				"StorageInstructions": {
					Placeholder: "How to store this ingredient...",
					InputType:   "textarea",
				},
				"ShoppingSuggestions": {
					Placeholder: "Tips for purchasing...",
					InputType:   "textarea",
				},
				// Allergen flags
				"ContainsEgg":       {InputType: "checkbox"},
				"ContainsDairy":     {InputType: "checkbox"},
				"ContainsFish":      {InputType: "checkbox"},
				"ContainsShellfish": {InputType: "checkbox"},
				"ContainsPeanut":    {InputType: "checkbox"},
				"ContainsTreeNut":   {InputType: "checkbox"},
				"ContainsWheat":     {InputType: "checkbox"},
				"ContainsSoy":       {InputType: "checkbox"},
				"ContainsSesame":    {InputType: "checkbox"},
				"ContainsGluten":    {InputType: "checkbox"},
				"ContainsAlcohol":   {InputType: "checkbox"},
				// Property flags
				"AnimalDerived": {InputType: "checkbox"},
				"AnimalFlesh":   {InputType: "checkbox"},
				"IsLiquid":      {InputType: "checkbox"},
				"IsStarch":      {InputType: "checkbox"},
				"IsProtein":     {InputType: "checkbox"},
				"IsGrain":       {InputType: "checkbox"},
				"IsFruit":       {InputType: "checkbox"},
				"IsSalt":        {InputType: "checkbox"},
				"IsFat":         {InputType: "checkbox"},
				"IsAcid":        {InputType: "checkbox"},
				"IsHeat":        {InputType: "checkbox"},
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
					Fields:  []string{"Warning"},
					Columns: 1,
				},
				{
					Fields:  []string{"StorageInstructions"},
					Columns: 1,
				},
				{
					Fields:  []string{"ShoppingSuggestions"},
					Columns: 1,
				},
				{
					Fields:  []string{"ContainsEgg", "ContainsDairy", "ContainsFish"},
					Columns: 3,
				},
				{
					Fields:  []string{"ContainsShellfish", "ContainsPeanut", "ContainsTreeNut"},
					Columns: 3,
				},
				{
					Fields:  []string{"ContainsWheat", "ContainsSoy", "ContainsSesame"},
					Columns: 3,
				},
				{
					Fields:  []string{"ContainsGluten", "ContainsAlcohol", "AnimalDerived"},
					Columns: 3,
				},
				{
					Fields:  []string{"AnimalFlesh", "IsLiquid", "IsStarch"},
					Columns: 3,
				},
				{
					Fields:  []string{"IsProtein", "IsGrain", "IsFruit"},
					Columns: 3,
				},
				{
					Fields:  []string{"IsSalt", "IsFat", "IsAcid", "IsHeat"},
					Columns: 4,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Valid Ingredients",
			CancelURL:        "/valid_ingredients",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Ingredients", URL: "/valid_ingredients"},
			{Text: validIngredient.Name, URL: ""},
		},

		// Dynamic subtitle showing ingredient info
		SubtitleGenerator: func(vi *mealplanningsvc.ValidIngredient) string {
			return fmt.Sprintf("Viewing ingredient: %s", vi.Name)
		},

		// Additional content - associations
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("grid grid-cols-1 md:grid-cols-2 gap-6 mt-6"),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Measurement Units",
					Subtitle: "Valid measurement units for this ingredient",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, measurementUnitsAssociations)),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Preparations",
					Subtitle: "Valid preparations for this ingredient",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette, preparationsAssociations)),
			),
		},
	})
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Ingredients", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithSpan(span)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError("Error: No API client available")), nil
	}

	// Extract QueryFilter and convert to gRPC filter
	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)

	validIngredientsRes, err := c.GetValidIngredients(ctx, &mealplanningsvc.GetValidIngredientsRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError(fmt.Sprintf("Error loading valid ingredients: %v", err))), nil
	}

	logger.WithValue("pagination", validIngredientsRes.Pagination).Info("Valid ingredients retrieved")

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validIngredientsRes.Pagination)

	// Use search endpoint for pagination buttons to return just the table content
	// The main page URL is used for deep linking via hx-push-url
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/valid_ingredients/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/valid_ingredients", queryFilter)

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidIngredient]{
		Title:             "Valid Ingredients",
		BaseSubtitle:      "Manage ingredient definitions",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search ingredients...",
		HTMXSearchTarget:  "/api/valid_ingredients/search",
		Data:              validIngredientsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Ingredient", "/valid_ingredients/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidIngredient]{
			TableID: "valid-ingredients-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"MealPlanTaskID",
				"Name",
				"PluralName",
				"Description",
				"IsLiquid",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"PluralName": "Plural Name",
				"IsLiquid":   "Liquid?",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
				"IsLiquid": func(value any) g.Node {
					if b, ok := value.(bool); ok && b {
						return g.Text("Yes")
					}
					return g.Text("No")
				},
			},
			Pagination:             pagination,
			PaginationURLGenerator: paginationURLGenerator,
			DeepLinkURLGenerator:   deepLinkURLGenerator,
			PaginationHTMXTarget:   "#search-results",
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredient) string {
			return fmt.Sprintf("/valid_ingredients/%s", data.Id)
		},
		EmptyStateTitle:       "No valid ingredients found",
		EmptyStateDescription: "No valid ingredients have been created yet.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Ingredient", "/valid_ingredients/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage ingredient definitions"
			}
			return fmt.Sprintf("Manage %d ingredient definitions", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Ingredients", s.renderValidIngredientsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Ingredients", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidIngredientsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validIngredientsRes, err := c.GetValidIngredients(ctx, &mealplanningsvc.GetValidIngredientsRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid ingredients: %v", err)),
			),
		), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validIngredientsRes.Pagination)

	// Diagnostic logging for pagination passed to table
	var appliedCursorStr string
	if pagination != nil && pagination.AppliedQueryFilter != nil && pagination.AppliedQueryFilter.Cursor != nil {
		appliedCursorStr = *pagination.AppliedQueryFilter.Cursor
	}
	logger := s.logger.WithSpan(span)
	logger.WithValue("table_pagination", map[string]interface{}{
		"cursor":                   pagination.Cursor,
		"appliedCursor":            appliedCursorStr,
		"hasAppliedFilter":         pagination.AppliedQueryFilter != nil,
		"appliedFilterCursorIsNil": pagination.AppliedQueryFilter != nil && pagination.AppliedQueryFilter.Cursor == nil,
		"appliedFilterCursorValue": appliedCursorStr,
	}).Info("Pagination object being passed to table component (search)")

	// Use search endpoint for pagination buttons to return just the table content
	// The main page URL is used for deep linking via hx-push-url
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/valid_ingredients/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/valid_ingredients", queryFilter)

	// Generate just the table (not the full page)
	if len(validIngredientsRes.Results) == 0 {
		searchQuery := req.URL.Query().Get("search")

		// Check if we're on a paginated page (not page 1) - if so, we need to show pagination controls
		var isOnPage2Plus bool
		if pagination != nil && pagination.AppliedQueryFilter != nil && pagination.AppliedQueryFilter.Cursor != nil {
			cursorValue := strings.TrimSpace(*pagination.AppliedQueryFilter.Cursor)
			isOnPage2Plus = cursorValue != ""
		}

		// If we're on a paginated page, include pagination controls
		if isOnPage2Plus && pagination != nil {
			paginationControls := components.CreatePaginationControls(&components.TableOptions[*mealplanningsvc.ValidIngredient]{
				Pagination:             pagination,
				PaginationURLGenerator: paginationURLGenerator,
				DeepLinkURLGenerator:   deepLinkURLGenerator,
				PaginationHTMXTarget:   "#search-results",
			}, &design.StandardPalette)

			return g.El("div",
				g.Attr("class", "overflow-x-auto"),
				components.EmptyState(
					"No valid ingredients found",
					fmt.Sprintf("No valid ingredients match the search term '%s'.", searchQuery),
					&design.StandardPalette,
					[]g.Node{},
				),
				paginationControls,
			), nil
		}

		// Otherwise, just show empty state (we're on page 1)
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid ingredients found",
				fmt.Sprintf("No valid ingredients match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(validIngredientsRes.Results, &components.TableOptions[*mealplanningsvc.ValidIngredient]{
		TableID: "valid-ingredients-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"MealPlanTaskID",
			"Name",
			"PluralName",
			"Description",
			"IsLiquid",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"PluralName": "Plural Name",
			"IsLiquid":   "Liquid?",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
			"IsLiquid": func(value any) g.Node {
				if b, ok := value.(bool); ok && b {
					return g.Text("Yes")
				}
				return g.Text("No")
			},
		},
		Pagination:             pagination,
		PaginationURLGenerator: paginationURLGenerator,
		DeepLinkURLGenerator:   deepLinkURLGenerator,
		PaginationHTMXTarget:   "#search-results",
		RowLinkGenerator: func(data *mealplanningsvc.ValidIngredient) string {
			return fmt.Sprintf("/valid_ingredients/%s", data.Id)
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

// renderValidIngredientsError creates a consistent error display for the valid ingredients page.
func (s *AdminFrontendServer) renderValidIngredientsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Ingredients",
		Subtitle: "Manage ingredient definitions",
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
