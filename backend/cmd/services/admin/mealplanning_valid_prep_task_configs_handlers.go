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
	validPrepTaskConfigIDURLParamKey = "validPrepTaskConfigID"

	unknown      = "Unknown"
	notSpecified = "Not specified"
)

func (s *AdminFrontendServer) ValidPrepTaskConfigPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Prep Task Config", s.renderValidPrepTaskConfigsError("Error: No API client available")), nil
	}

	validPrepTaskConfigID := s.validPrepTaskConfigIDRouteParamFetcher(req)
	if validPrepTaskConfigID == "" {
		return page("Valid Prep Task Config", s.renderValidPrepTaskConfigsError("Error: No valid prep task config MealPlanTaskID provided")), nil
	}

	validPrepTaskConfigRes, err := c.GetValidPrepTaskConfig(ctx, &mealplanningsvc.GetValidPrepTaskConfigRequest{ValidPrepTaskConfigId: validPrepTaskConfigID})
	if err != nil {
		return page("Valid Prep Task Config", s.renderValidPrepTaskConfigsError(fmt.Sprintf("Error loading valid prep task config: %v", err))), nil
	}

	if validPrepTaskConfigRes == nil || validPrepTaskConfigRes.Result == nil {
		return page("Valid Prep Task Config", s.renderValidPrepTaskConfigsError("Error: Valid prep task config not found")), nil
	}

	validPrepTaskConfig := validPrepTaskConfigRes.Result

	// Use the FormPage component for viewing valid prep task config data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.ValidPrepTaskConfig]{
		Title:        "Valid Prep Task Config Details",
		BaseSubtitle: "View prep task config information",
		Palette:      &design.StandardPalette,
		Data:         validPrepTaskConfig,
		FormOptions: &components.FormOptions[*mealplanningsvc.ValidPrepTaskConfig]{
			Palette: &design.StandardPalette,
			FormID:  "view-valid-prep-task-config-form",
			Action:  fmt.Sprintf("/api/valid_prep_task_configs/%s", validPrepTaskConfig.Id),
			Method:  "PUT",

			// Fields that can be edited
			EnabledFields: []string{
				"StorageType",
				"StorageInstructions",
				"Notes",
				"Source",
			},

			FieldConfigs: map[string]*components.FieldConfig{
				"StorageType": {
					Placeholder: "Storage container type (e.g., airtight, covered)...",
				},
				"StorageInstructions": {
					Placeholder: "Instructions for storing the prepped ingredient...",
					InputType:   "textarea",
				},
				"Notes": {
					Placeholder: "Additional notes...",
					InputType:   "textarea",
				},
				"Source": {
					Placeholder: "Source of this information...",
				},
			},

			FormRows: []*components.FormRow{
				{
					Fields:  []string{"StorageType"},
					Columns: 1,
				},
				{
					Fields:  []string{"StorageInstructions"},
					Columns: 1,
				},
				{
					Fields:  []string{"Notes"},
					Columns: 1,
				},
				{
					Fields:  []string{"Source"},
					Columns: 1,
				},
			},

			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Prep Task Configs",
			CancelURL:        "/valid_prep_task_configs",

			HTMXTarget:  "body",
			HTMXSwap:    "innerHTML",
			HTMXPushURL: true,
		},

		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Enumerations", URL: ""},
			{Text: "Valid Prep Task Configs", URL: "/valid_prep_task_configs"},
			{Text: validPrepTaskConfig.Id, URL: ""},
		},

		// Dynamic subtitle showing config info
		SubtitleGenerator: func(cfg *mealplanningsvc.ValidPrepTaskConfig) string {
			ingredientName := unknown
			preparationName := unknown
			if cfg.Ingredient != nil {
				ingredientName = cfg.Ingredient.Name
			}
			if cfg.Preparation != nil {
				preparationName = cfg.Preparation.Name
			}
			return fmt.Sprintf("%s + %s", preparationName, ingredientName)
		},

		// Additional content - related information
		AdditionalContent: []g.Node{
			ghtml.Div(
				ghtml.Class("grid grid-cols-1 md:grid-cols-2 gap-6 mt-6"),
				// Ingredient info card
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Ingredient",
					Subtitle: "The ingredient this config applies to",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette,
					ghtml.Div(
						ghtml.Class("space-y-2"),
						ghtml.P(
							ghtml.Strong(g.Text("Name: ")),
							g.Text(func() string {
								if validPrepTaskConfig.Ingredient != nil {
									return validPrepTaskConfig.Ingredient.Name
								}
								return unknown
							}()),
						),
						ghtml.P(
							ghtml.Strong(g.Text("Description: ")),
							g.Text(func() string {
								if validPrepTaskConfig.Ingredient != nil {
									return validPrepTaskConfig.Ingredient.Description
								}
								return "No description"
							}()),
						),
					),
				)),
				// Preparation info card
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Preparation",
					Subtitle: "The preparation method this config applies to",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette,
					ghtml.Div(
						ghtml.Class("space-y-2"),
						ghtml.P(
							ghtml.Strong(g.Text("Name: ")),
							g.Text(func() string {
								if validPrepTaskConfig.Preparation != nil {
									return validPrepTaskConfig.Preparation.Name
								}
								return unknown
							}()),
						),
						ghtml.P(
							ghtml.Strong(g.Text("Description: ")),
							g.Text(func() string {
								if validPrepTaskConfig.Preparation != nil {
									return validPrepTaskConfig.Preparation.Description
								}
								return "No description"
							}()),
						),
					),
				)),
			),
			// Storage duration info
			ghtml.Div(
				ghtml.Class("mt-6"),
				components.ContentContainer(&components.ContentContainerProps{
					Title:    "Storage Duration",
					Subtitle: "How long the prepped ingredient can be stored",
					Palette:  &design.StandardPalette,
				}, components.Card(&design.StandardPalette,
					ghtml.Div(
						ghtml.Class("space-y-2"),
						ghtml.P(
							ghtml.Strong(g.Text("Minimum Duration: ")),
							g.Text(func() string {
								if validPrepTaskConfig.StorageDurationInSeconds != nil {
									hours := validPrepTaskConfig.StorageDurationInSeconds.Min / 3600
									return fmt.Sprintf("%d hours", hours)
								}
								return notSpecified
							}()),
						),
						ghtml.P(
							ghtml.Strong(g.Text("Maximum Duration: ")),
							g.Text(func() string {
								if validPrepTaskConfig.StorageDurationInSeconds != nil && validPrepTaskConfig.StorageDurationInSeconds.Max != nil {
									hours := *validPrepTaskConfig.StorageDurationInSeconds.Max / 3600
									return fmt.Sprintf("%d hours", hours)
								}
								return notSpecified
							}()),
						),
						ghtml.P(
							ghtml.Strong(g.Text("Storage Temperature: ")),
							g.Text(func() string {
								if validPrepTaskConfig.StorageTemperatureInCelsius != nil {
									minTemp := ""
									maxTemp := ""
									if validPrepTaskConfig.StorageTemperatureInCelsius.Min != nil {
										roundedMin := roundTemperatureToNearest5(*validPrepTaskConfig.StorageTemperatureInCelsius.Min)
										minTemp = fmt.Sprintf("%.0f°C", roundedMin)
									}
									if validPrepTaskConfig.StorageTemperatureInCelsius.Max != nil {
										roundedMax := roundTemperatureToNearest5(*validPrepTaskConfig.StorageTemperatureInCelsius.Max)
										maxTemp = fmt.Sprintf("%.0f°C", roundedMax)
									}

									switch {
									case minTemp != "" && maxTemp != "":
										return fmt.Sprintf("%s to %s", minTemp, maxTemp)
									case minTemp != "":
										return fmt.Sprintf("Min: %s", minTemp)
									case maxTemp != "":
										return fmt.Sprintf("Max: %s", maxTemp)
									}
								}
								return notSpecified
							}()),
						),
					),
				)),
			),
		},
	})
	if err != nil {
		return page("Valid Prep Task Config", s.renderValidPrepTaskConfigsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Valid Prep Task Config", formPageResult.Node), nil
}

func (s *AdminFrontendServer) ValidPrepTaskConfigsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithSpan(span)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Valid Prep Task Configs", s.renderValidPrepTaskConfigsError("Error: No API client available")), nil
	}

	// Extract QueryFilter and convert to gRPC filter
	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)

	validPrepTaskConfigsRes, err := c.GetValidPrepTaskConfigs(ctx, &mealplanningsvc.GetValidPrepTaskConfigsRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return page("Valid Prep Task Configs", s.renderValidPrepTaskConfigsError(fmt.Sprintf("Error loading valid prep task configs: %v", err))), nil
	}

	logger.WithValue("pagination", validPrepTaskConfigsRes.Pagination).Info("Valid prep task configs retrieved")

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validPrepTaskConfigsRes.Pagination)

	// Use search endpoint for pagination buttons to return just the table content
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/valid_prep_task_configs/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/valid_prep_task_configs", queryFilter)

	// Use the integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.ValidPrepTaskConfig]{
		Title:             "Valid Prep Task Configs",
		BaseSubtitle:      "Manage prep task storage configurations",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search prep task configs...",
		HTMXSearchTarget:  "/api/valid_prep_task_configs/search",
		Data:              validPrepTaskConfigsRes.Results,
		TableOptions: &components.TableOptions[*mealplanningsvc.ValidPrepTaskConfig]{
			TableID: "valid-prep-task-configs-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"MealPlanTaskID",
				"Ingredient",
				"Preparation",
				"StorageType",
				"CreatedAt",
			},
			FieldReplacements: map[string]string{
				"StorageType": "Storage Type",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"CreatedAt": renderTimestamp,
				"Ingredient": func(value any) g.Node {
					if ingredient, ok := value.(*mealplanningsvc.ValidIngredient); ok && ingredient != nil {
						return g.Text(ingredient.Name)
					}
					return g.Text("-")
				},
				"Preparation": func(value any) g.Node {
					if preparation, ok := value.(*mealplanningsvc.ValidPreparation); ok && preparation != nil {
						return g.Text(preparation.Name)
					}
					return g.Text("-")
				},
			},
			Pagination:             pagination,
			PaginationURLGenerator: paginationURLGenerator,
			DeepLinkURLGenerator:   deepLinkURLGenerator,
			PaginationHTMXTarget:   "#search-results",
		},
		RowLinkGenerator: func(data *mealplanningsvc.ValidPrepTaskConfig) string {
			return fmt.Sprintf("/valid_prep_task_configs/%s", data.Id)
		},
		EmptyStateTitle:       "No valid prep task configs found",
		EmptyStateDescription: "No prep task configs have been created yet.",
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage prep task storage configurations"
			}
			return fmt.Sprintf("Manage %d prep task storage configurations", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Valid Prep Task Configs", s.renderValidPrepTaskConfigsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Valid Prep Task Configs", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) ValidPrepTaskConfigsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	validPrepTaskConfigsRes, err := c.GetValidPrepTaskConfigs(ctx, &mealplanningsvc.GetValidPrepTaskConfigsRequest{
		Filter: grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading valid prep task configs: %v", err)),
			),
		), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(validPrepTaskConfigsRes.Pagination)

	// Use search endpoint for pagination buttons to return just the table content
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/valid_prep_task_configs/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/valid_prep_task_configs", queryFilter)

	// Handle empty results
	if len(validPrepTaskConfigsRes.Results) == 0 {
		searchQuery := req.URL.Query().Get("search")

		// Check if we're on a paginated page (not page 1)
		var isOnPage2Plus bool
		if pagination != nil && pagination.AppliedQueryFilter != nil && pagination.AppliedQueryFilter.Cursor != nil {
			cursorValue := strings.TrimSpace(*pagination.AppliedQueryFilter.Cursor)
			isOnPage2Plus = cursorValue != ""
		}

		// If we're on a paginated page, include pagination controls
		if isOnPage2Plus && pagination != nil {
			paginationControls := components.CreatePaginationControls(&components.TableOptions[*mealplanningsvc.ValidPrepTaskConfig]{
				Pagination:             pagination,
				PaginationURLGenerator: paginationURLGenerator,
				DeepLinkURLGenerator:   deepLinkURLGenerator,
				PaginationHTMXTarget:   "#search-results",
			}, &design.StandardPalette)

			return g.El("div",
				g.Attr("class", "overflow-x-auto"),
				components.EmptyState(
					"No valid prep task configs found",
					fmt.Sprintf("No valid prep task configs match the search term '%s'.", searchQuery),
					&design.StandardPalette,
					[]g.Node{},
				),
				paginationControls,
			), nil
		}

		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No valid prep task configs found",
				fmt.Sprintf("No valid prep task configs match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(validPrepTaskConfigsRes.Results, &components.TableOptions[*mealplanningsvc.ValidPrepTaskConfig]{
		TableID: "valid-prep-task-configs-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"MealPlanTaskID",
			"Ingredient",
			"Preparation",
			"StorageType",
			"CreatedAt",
		},
		FieldReplacements: map[string]string{
			"StorageType": "Storage Type",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"CreatedAt": renderTimestamp,
			"Ingredient": func(value any) g.Node {
				if ingredient, ok := value.(*mealplanningsvc.ValidIngredient); ok && ingredient != nil {
					return g.Text(ingredient.Name)
				}
				return g.Text("-")
			},
			"Preparation": func(value any) g.Node {
				if preparation, ok := value.(*mealplanningsvc.ValidPreparation); ok && preparation != nil {
					return g.Text(preparation.Name)
				}
				return g.Text("-")
			},
		},
		Pagination:             pagination,
		PaginationURLGenerator: paginationURLGenerator,
		DeepLinkURLGenerator:   deepLinkURLGenerator,
		PaginationHTMXTarget:   "#search-results",
		RowLinkGenerator: func(data *mealplanningsvc.ValidPrepTaskConfig) string {
			return fmt.Sprintf("/valid_prep_task_configs/%s", data.Id)
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

	return g.El("div",
		g.Attr("class", "overflow-x-auto"),
		table,
	), nil
}

// renderValidPrepTaskConfigsError creates a consistent error display for the valid prep task configs page.
func (s *AdminFrontendServer) renderValidPrepTaskConfigsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Valid Prep Task Configs",
		Subtitle: "Manage prep task storage configurations",
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
