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

const (
	recipeIDURLParamKey = "recipeID"
)

func (s *AdminFrontendServer) RecipePage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Recipes", s.renderRecipesError("Error: No API client available")), nil
	}

	recipeID := s.recipeIDRouteParamFetcher(req)
	if recipeID == "" {
		return page("Recipes", s.renderRecipesError("Error: No recipe ID provided")), nil
	}

	recipeRes, err := c.GetRecipe(ctx, &mealplanningsvc.GetRecipeRequest{RecipeId: recipeID})
	if err != nil {
		return page("Recipes", s.renderRecipesError(fmt.Sprintf("Error loading recipe: %v", err))), nil
	}
	recipe := recipeRes.Result

	// Use FormPage component for viewing recipe data
	formPageResult, err := components.FormPage(&components.FormPageProps[*mealplanningsvc.Recipe]{
		Title:        "Recipe Details",
		BaseSubtitle: "View recipe information",
		Palette:      &design.StandardPalette,
		Data:         recipe,
		FormOptions: &components.FormOptions[*mealplanningsvc.Recipe]{
			FormID: "view-recipe-form",
			// Read-only form, no action needed
			EnabledFields: []string{
				"id",
				"name",
				"description",
				"source",
				"status",
				"slug",
				"portion_name",
				"plural_portion_name",
				"created_by_user",
				"created_at",
				"last_updated_at",
				"archived_at",
				"eligible_for_meals",
			},
			FieldConfigs: map[string]*components.FieldConfig{
				"id": {
					DisplayName: "ID",
				},
				"name": {
					DisplayName: "Name",
				},
				"description": {
					DisplayName: "Description",
				},
				"source": {
					DisplayName: "Source",
				},
				"status": {
					DisplayName: "Status",
				},
				"slug": {
					DisplayName: "Slug",
				},
				"portion_name": {
					DisplayName: "Portion Name",
				},
				"plural_portion_name": {
					DisplayName: "Plural Portion Name",
				},
				"created_by_user": {
					DisplayName: "Created By User",
				},
				"created_at": {
					DisplayName: "Created At",
				},
				"last_updated_at": {
					DisplayName: "Last Updated At",
				},
				"archived_at": {
					DisplayName: "Archived At",
				},
				"eligible_for_meals": {
					DisplayName: "Eligible For Meals",
				},
			},
			FormRows: []*components.FormRow{
				{
					Fields:  []string{"id"},
					Columns: 1,
				},
				{
					Fields:  []string{"name"},
					Columns: 1,
				},
				{
					Fields:  []string{"description"},
					Columns: 1,
				},
				{
					Fields:  []string{"source", "status"},
					Columns: 2,
				},
				{
					Fields:  []string{"slug"},
					Columns: 1,
				},
				{
					Fields:  []string{"portion_name", "plural_portion_name"},
					Columns: 2,
				},
				{
					Fields:  []string{"created_by_user"},
					Columns: 1,
				},
				{
					Fields:  []string{"created_at", "last_updated_at"},
					Columns: 2,
				},
				{
					Fields:  []string{"archived_at"},
					Columns: 1,
				},
				{
					Fields:  []string{"eligible_for_meals"},
					Columns: 1,
				},
			},
		},
		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Recipes", URL: "/recipes"},
			{Text: recipe.Name, URL: ""},
		},

		SubtitleGenerator: func(r *mealplanningsvc.Recipe) string {
			return fmt.Sprintf("Viewing recipe: %s", r.Name)
		},
	})
	if err != nil {
		return page("Recipes", s.renderRecipesError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	// Create steps section
	stepsSection := components.CardWithHeader(
		"Recipe Steps",
		&design.StandardPalette,
		nil,
		s.renderRecipeSteps(recipe.Steps),
	)

	// Combine form and steps section
	return page("Recipes",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			formPageResult.Node,
			stepsSection,
		),
	), nil
}

func (s *AdminFrontendServer) renderRecipeSteps(steps []*mealplanningsvc.RecipeStep) g.Node {
	if len(steps) == 0 {
		return ghtml.Div(
			ghtml.Class("text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("This recipe has no steps."),
			),
		)
	}

	var stepNodes []g.Node
	for i, step := range steps {
		stepNodes = append(stepNodes,
			ghtml.Div(
				ghtml.Class("border border-gray-200 rounded-lg p-4 mb-4 last:mb-0 bg-white"),
				ghtml.Div(
					ghtml.Class("space-y-4"),
					// Step header
					ghtml.Div(
						ghtml.Class("flex items-start gap-4"),
						// Step number
						ghtml.Div(
							ghtml.Class("flex-shrink-0 w-10 h-10 rounded-full bg-blue-100 text-blue-800 flex items-center justify-center font-semibold text-sm"),
							g.Text(fmt.Sprintf("%d", i+1)),
						),
						// Step header content
						ghtml.Div(
							ghtml.Class("flex-1"),
							ghtml.Div(
								ghtml.Class("flex items-center gap-2 mb-2"),
								ghtml.Div(
									ghtml.Class("text-xs text-gray-500"),
									g.Text(fmt.Sprintf("Index: %d | ID: %s", step.Index, step.Id)),
								),
								g.If(step.Optional,
									ghtml.Span(
										ghtml.Class("inline-block px-2 py-1 text-xs font-medium bg-yellow-100 text-yellow-800 rounded"),
										g.Text("Optional"),
									),
								),
								g.If(step.StartTimerAutomatically,
									ghtml.Span(
										ghtml.Class("inline-block px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded"),
										g.Text("Auto Timer"),
									),
								),
							),
							// Preparation
							g.If(step.Preparation != nil,
								ghtml.Div(
									ghtml.Class("font-semibold text-lg text-gray-900 mb-1"),
									g.Text(step.Preparation.Name),
								),
							),
							// Notes
							g.If(step.Notes != "",
								ghtml.Div(
									ghtml.Class("text-sm text-gray-700 mb-1"),
									g.Text(step.Notes),
								),
							),
							// Explicit instructions
							g.If(step.ExplicitInstructions != "",
								ghtml.Div(
									ghtml.Class("text-sm text-gray-600 italic mb-1"),
									g.Text(step.ExplicitInstructions),
								),
							),
							// Condition expression
							g.If(step.ConditionExpression != "",
								ghtml.Div(
									ghtml.Class("text-xs text-purple-600 font-mono bg-purple-50 p-2 rounded mb-1"),
									ghtml.Span(
										ghtml.Class("font-semibold"),
										g.Text("Condition: "),
									),
									g.Text(step.ConditionExpression),
								),
							),
							// Time and temperature
							ghtml.Div(
								ghtml.Class("flex gap-4 text-xs text-gray-600 mb-2"),
								g.If(step.EstimatedTimeInSeconds != nil,
									ghtml.Div(
										g.Text(fmt.Sprintf("Time: %d-%d seconds",
											func() uint32 {
												if step.EstimatedTimeInSeconds.Min != nil {
													return *step.EstimatedTimeInSeconds.Min
												}
												return 0
											}(),
											func() uint32 {
												if step.EstimatedTimeInSeconds.Max != nil {
													return *step.EstimatedTimeInSeconds.Max
												}
												return 0
											}(),
										)),
									),
								),
								g.If(step.TemperatureInCelsius != nil,
									ghtml.Div(
										g.Text(fmt.Sprintf("Temperature: %.1f-%.1f°C",
											func() float32 {
												if step.TemperatureInCelsius.Min != nil {
													return *step.TemperatureInCelsius.Min
												}
												return 0
											}(),
											func() float32 {
												if step.TemperatureInCelsius.Max != nil {
													return *step.TemperatureInCelsius.Max
												}
												return 0
											}(),
										)),
									),
								),
							),
						),
					),
					// Ingredients section
					g.If(len(step.Ingredients) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text(fmt.Sprintf("Ingredients (%d):", len(step.Ingredients))),
							),
							ghtml.Div(
								ghtml.Class("space-y-2"),
								g.Group(s.renderStepIngredients(step.Ingredients, steps)),
							),
						),
					),
					// Instruments section
					g.If(len(step.Instruments) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text(fmt.Sprintf("Instruments (%d):", len(step.Instruments))),
							),
							ghtml.Div(
								ghtml.Class("space-y-2"),
								g.Group(s.renderStepInstruments(step.Instruments, steps)),
							),
						),
					),
					// Vessels section
					g.If(len(step.Vessels) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text(fmt.Sprintf("Vessels (%d):", len(step.Vessels))),
							),
							ghtml.Div(
								ghtml.Class("space-y-2"),
								g.Group(s.renderStepVessels(step.Vessels, steps)),
							),
						),
					),
					// Products section
					g.If(len(step.Products) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text(fmt.Sprintf("Products (%d):", len(step.Products))),
							),
							ghtml.Div(
								ghtml.Class("space-y-2"),
								g.Group(s.renderStepProducts(step.Products)),
							),
						),
					),
					// Completion conditions section
					g.If(len(step.CompletionConditions) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text(fmt.Sprintf("Completion Conditions (%d):", len(step.CompletionConditions))),
							),
							ghtml.Div(
								ghtml.Class("space-y-2"),
								g.Group(s.renderStepCompletionConditions(step.CompletionConditions)),
							),
						),
					),
				),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("space-y-4"),
		g.Group(stepNodes),
	)
}

func (s *AdminFrontendServer) renderStepIngredients(ingredients []*mealplanningsvc.RecipeStepIngredient, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	var nodes []g.Node
	for _, ing := range ingredients {
		var details []g.Node

		// Ingredient name
		if ing.Ingredient != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(ing.Ingredient.Name),
			))
		} else if ing.Name != "" {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(ing.Name),
			))
		}

		// Quantity and unit
		if ing.Quantity != nil {
			qtyStr := ""
			if ing.Quantity.Min != 0 {
				qtyStr = fmt.Sprintf("%.2f", ing.Quantity.Min)
				if ing.Quantity.Max != nil {
					qtyStr += fmt.Sprintf("-%.2f", *ing.Quantity.Max)
				}
			}
			if qtyStr != "" {
				unitName := formatMeasurementUnitName(ing.MeasurementUnit)
				details = append(details, ghtml.Span(
					ghtml.Class("text-gray-600"),
					g.Text(fmt.Sprintf(" (%s %s)", qtyStr, unitName)),
				))
			}
		}

		// Product reference (from previous step)
		if ing.RecipeStepProductId != nil {
			productStep := s.findStepWithProduct(*ing.RecipeStepProductId, allSteps)
			if productStep != nil {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 font-semibold ml-2"),
					g.Text(fmt.Sprintf("← Product from Step %d", productStep.Index+1)),
				))
			} else {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 ml-2"),
					g.Text(fmt.Sprintf("← Product ID: %s", *ing.RecipeStepProductId)),
				))
			}
		}

		// Recipe product reference
		if ing.RecipeStepProductRecipeId != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("text-purple-600 ml-2"),
				g.Text(fmt.Sprintf("← Recipe Product: %s", *ing.RecipeStepProductRecipeId)),
			))
		}

		// Vessel index
		if ing.VesselIndex != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("text-gray-500 ml-2"),
				g.Text(fmt.Sprintf("(Vessel %d)", *ing.VesselIndex)),
			))
		}

		// Product percentage
		if ing.ProductPercentageToUse != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("text-gray-500 ml-2"),
				g.Text(fmt.Sprintf("(%d%% of product)", int(*ing.ProductPercentageToUse*100))),
			))
		}

		// Flags
		var flags []g.Node
		if ing.Optional {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-yellow-100 text-yellow-800 rounded ml-2"),
				g.Text("Optional"),
			))
		}
		if ing.ToTaste {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-orange-100 text-orange-800 rounded ml-2"),
				g.Text("To Taste"),
			))
		}
		if ing.OptionIndex > 0 {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-gray-100 text-gray-800 rounded ml-2"),
				g.Text(fmt.Sprintf("Option %d", ing.OptionIndex)),
			))
		}

		// Notes
		if ing.IngredientNotes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Ingredient notes: %s", ing.IngredientNotes)),
			))
		}
		if ing.QuantityNotes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Quantity notes: %s", ing.QuantityNotes)),
			))
		}

		nodes = append(nodes, ghtml.Div(
			ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		))
	}
	return nodes
}

func (s *AdminFrontendServer) renderStepInstruments(instruments []*mealplanningsvc.RecipeStepInstrument, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	var nodes []g.Node
	for _, inst := range instruments {
		var details []g.Node

		// Instrument name
		if inst.Instrument != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(inst.Instrument.Name),
			))
		} else if inst.Name != "" {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(inst.Name),
			))
		}

		// Quantity
		if inst.Quantity != nil {
			qtyStr := ""
			if inst.Quantity.Min != 0 {
				qtyStr = fmt.Sprintf("%d", inst.Quantity.Min)
				if inst.Quantity.Max != nil {
					qtyStr += fmt.Sprintf("-%d", *inst.Quantity.Max)
				}
			}
			if qtyStr != "" {
				details = append(details, ghtml.Span(
					ghtml.Class("text-gray-600"),
					g.Text(fmt.Sprintf(" (%s)", qtyStr)),
				))
			}
		}

		// Product reference
		if inst.RecipeStepProductId != nil {
			productStep := s.findStepWithProduct(*inst.RecipeStepProductId, allSteps)
			if productStep != nil {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 font-semibold ml-2"),
					g.Text(fmt.Sprintf("← Product from Step %d", productStep.Index+1)),
				))
			} else {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 ml-2"),
					g.Text(fmt.Sprintf("← Product ID: %s", *inst.RecipeStepProductId)),
				))
			}
		}

		// Flags
		var flags []g.Node
		if inst.Optional {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-yellow-100 text-yellow-800 rounded ml-2"),
				g.Text("Optional"),
			))
		}
		if inst.OptionIndex > 0 {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-gray-100 text-gray-800 rounded ml-2"),
				g.Text(fmt.Sprintf("Option %d", inst.OptionIndex)),
			))
		}
		if inst.PreferenceRank > 0 {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-indigo-100 text-indigo-800 rounded ml-2"),
				g.Text(fmt.Sprintf("Rank %d", inst.PreferenceRank)),
			))
		}

		// Notes
		if inst.Notes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Notes: %s", inst.Notes)),
			))
		}

		nodes = append(nodes, ghtml.Div(
			ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		))
	}
	return nodes
}

func (s *AdminFrontendServer) renderStepVessels(vessels []*mealplanningsvc.RecipeStepVessel, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	var nodes []g.Node
	for _, vessel := range vessels {
		var details []g.Node

		// Vessel name
		if vessel.Vessel != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(vessel.Vessel.Name),
			))
		} else if vessel.Name != "" {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(vessel.Name),
			))
		}

		// Quantity
		if vessel.Quantity != nil {
			qtyStr := ""
			if vessel.Quantity.Min != 0 {
				qtyStr = fmt.Sprintf("%d", vessel.Quantity.Min)
				if vessel.Quantity.Max != nil {
					qtyStr += fmt.Sprintf("-%d", *vessel.Quantity.Max)
				}
			}
			if qtyStr != "" {
				details = append(details, ghtml.Span(
					ghtml.Class("text-gray-600"),
					g.Text(fmt.Sprintf(" (%s)", qtyStr)),
				))
			}
		}

		// Preposition
		if vessel.VesselPreposition != "" {
			details = append(details, ghtml.Span(
				ghtml.Class("text-gray-500"),
				g.Text(fmt.Sprintf(" (%s)", vessel.VesselPreposition)),
			))
		}

		// Product reference
		if vessel.RecipeStepProductId != nil {
			productStep := s.findStepWithProduct(*vessel.RecipeStepProductId, allSteps)
			if productStep != nil {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 font-semibold ml-2"),
					g.Text(fmt.Sprintf("← Product from Step %d", productStep.Index+1)),
				))
			} else {
				details = append(details, ghtml.Span(
					ghtml.Class("text-blue-600 ml-2"),
					g.Text(fmt.Sprintf("← Product ID: %s", *vessel.RecipeStepProductId)),
				))
			}
		}

		// Flags
		var flags []g.Node
		if vessel.UnavailableAfterStep {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-red-100 text-red-800 rounded ml-2"),
				g.Text("Unavailable After Step"),
			))
		}

		// Notes
		if vessel.Notes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Notes: %s", vessel.Notes)),
			))
		}

		nodes = append(nodes, ghtml.Div(
			ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		))
	}
	return nodes
}

func (s *AdminFrontendServer) renderStepProducts(products []*mealplanningsvc.RecipeStepProduct) []g.Node {
	var nodes []g.Node
	for _, product := range products {
		var details []g.Node

		// Product name and type
		details = append(details, ghtml.Span(
			ghtml.Class("font-medium"),
			g.Text(product.Name),
		))
		productTypeStr := formatProductType(product.Type)
		details = append(details, ghtml.Span(
			ghtml.Class("text-gray-500 text-xs ml-2"),
			g.Text(fmt.Sprintf("(%s)", productTypeStr)),
		))

		// Quantity and unit
		if product.Quantity != nil {
			qtyStr := ""
			if product.Quantity.Min != nil {
				qtyStr = fmt.Sprintf("%.2f", *product.Quantity.Min)
				if product.Quantity.Max != nil {
					qtyStr += fmt.Sprintf("-%.2f", *product.Quantity.Max)
				}
			}
			if qtyStr != "" {
				unitName := formatMeasurementUnitName(product.MeasurementUnit)
				details = append(details, ghtml.Span(
					ghtml.Class("text-gray-600 ml-2"),
					g.Text(fmt.Sprintf("%s %s", qtyStr, unitName)),
				))
			}
		}

		// Vessel index
		if product.ContainedInVesselIndex != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("text-gray-500 ml-2"),
				g.Text(fmt.Sprintf("(Vessel %d)", *product.ContainedInVesselIndex)),
			))
		}

		// Storage info
		if product.StorageTemperatureInCelsius != nil {
			tempStr := ""
			if product.StorageTemperatureInCelsius.Min != nil {
				tempStr = fmt.Sprintf("%.1f", *product.StorageTemperatureInCelsius.Min)
				if product.StorageTemperatureInCelsius.Max != nil {
					tempStr += fmt.Sprintf("-%.1f", *product.StorageTemperatureInCelsius.Max)
				}
				tempStr += "°C"
			}
			if tempStr != "" {
				details = append(details, ghtml.Div(
					ghtml.Class("text-xs text-gray-500 mt-1"),
					g.Text(fmt.Sprintf("Storage temp: %s", tempStr)),
				))
			}
		}
		if product.StorageDurationInSeconds != nil {
			durStr := ""
			if product.StorageDurationInSeconds.Min != nil {
				durStr = fmt.Sprintf("%d", *product.StorageDurationInSeconds.Min)
				if product.StorageDurationInSeconds.Max != nil {
					durStr += fmt.Sprintf("-%d", *product.StorageDurationInSeconds.Max)
				}
				durStr += "s"
			}
			if durStr != "" {
				details = append(details, ghtml.Div(
					ghtml.Class("text-xs text-gray-500 mt-1"),
					g.Text(fmt.Sprintf("Storage duration: %s", durStr)),
				))
			}
		}
		if product.StorageInstructions != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Storage: %s", product.StorageInstructions)),
			))
		}

		// Flags
		var flags []g.Node
		if product.IsWaste {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-red-100 text-red-800 rounded ml-2"),
				g.Text("Waste"),
			))
		}
		if product.IsLiquid {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-blue-100 text-blue-800 rounded ml-2"),
				g.Text("Liquid"),
			))
		}
		if product.Compostable {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-green-100 text-green-800 rounded ml-2"),
				g.Text("Compostable"),
			))
		}

		// Notes
		if product.QuantityNotes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Quantity notes: %s", product.QuantityNotes)),
			))
		}

		nodes = append(nodes, ghtml.Div(
			ghtml.Class("text-sm pl-4 border-l-2 border-blue-200"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		))
	}
	return nodes
}

func (s *AdminFrontendServer) renderStepCompletionConditions(conditions []*mealplanningsvc.RecipeStepCompletionCondition) []g.Node {
	var nodes []g.Node
	for _, cond := range conditions {
		var details []g.Node

		// Ingredient state
		if cond.IngredientState != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(cond.IngredientState.Name),
			))
		}

		// Ingredients list
		if len(cond.Ingredients) > 0 {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-600 mt-1"),
				g.Text(fmt.Sprintf("Ingredients: %d", len(cond.Ingredients))),
			))
		}

		// Flags
		var flags []g.Node
		if cond.Optional {
			flags = append(flags, ghtml.Span(
				ghtml.Class("inline-block px-1.5 py-0.5 text-xs bg-yellow-100 text-yellow-800 rounded ml-2"),
				g.Text("Optional"),
			))
		}

		// Notes
		if cond.Notes != "" {
			details = append(details, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 mt-1"),
				g.Text(fmt.Sprintf("Notes: %s", cond.Notes)),
			))
		}

		nodes = append(nodes, ghtml.Div(
			ghtml.Class("text-sm pl-4 border-l-2 border-purple-200"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		))
	}
	return nodes
}

// findStepWithProduct finds the step that contains a product with the given ID.
func (s *AdminFrontendServer) findStepWithProduct(productID string, allSteps []*mealplanningsvc.RecipeStep) *mealplanningsvc.RecipeStep {
	for _, step := range allSteps {
		for _, product := range step.Products {
			if product.Id == productID {
				return step
			}
		}
	}
	return nil
}

// formatProductType formats a RecipeStepProductType enum value into a readable string.
func formatProductType(productType mealplanningsvc.RecipeStepProductType) string {
	switch productType {
	case mealplanningsvc.RecipeStepProductType_RECIPE_STEP_PRODUCT_TYPE_INGREDIENT:
		return "Ingredient"
	case mealplanningsvc.RecipeStepProductType_RECIPE_STEP_PRODUCT_TYPE_INSTRUMENT:
		return "Instrument"
	case mealplanningsvc.RecipeStepProductType_RECIPE_STEP_PRODUCT_TYPE_VESSEL:
		return "Vessel"
	default:
		return productType.String()
	}
}

// formatMeasurementUnitName formats a measurement unit name for display, removing test data numbering.
func formatMeasurementUnitName(unit *mealplanningsvc.ValidMeasurementUnit) string {
	if unit == nil {
		return ""
	}
	// Use PluralName if available (it's cleaner and doesn't have test numbering)
	if unit.PluralName != "" {
		return unit.PluralName
	}
	// Fall back to Name, but strip trailing numbers and spaces (for test data like "gram 1")
	name := unit.Name
	// Remove trailing space and number pattern (e.g., "gram 1" -> "gram")
	// This is a simple approach - if the name ends with " <number>", remove it
	if name != "" {
		// Check if it ends with a pattern like " 1", " 2", etc.
		lastSpaceIdx := -1
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == ' ' {
				lastSpaceIdx = i
				break
			}
		}
		if lastSpaceIdx >= 0 && lastSpaceIdx < len(name)-1 {
			// Check if everything after the space is a number
			trailing := name[lastSpaceIdx+1:]
			isNumber := true
			for _, r := range trailing {
				if r < '0' || r > '9' {
					isNumber = false
					break
				}
			}
			if isNumber {
				return name[:lastSpaceIdx]
			}
		}
	}
	return name
}

func (s *AdminFrontendServer) RecipesList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Recipes", s.renderRecipesError("Error: No API client available")), nil
	}

	// Extract QueryFilter and convert to gRPC filter
	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)

	// Extract status from query params, default to "submitted"
	status := req.URL.Query().Get("status")
	if status == "" {
		status = "submitted"
	}

	recipesRes, err := c.GetRecipes(ctx, &mealplanningsvc.GetRecipesRequest{
		Filter: grpcFilter,
		Status: status,
	})
	if err != nil {
		return page("Recipes", s.renderRecipesError(fmt.Sprintf("Error loading recipes: %v", err))), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(recipesRes.Pagination)
	// Use search endpoint for pagination buttons to return just the table content
	// The main page URL is still used for deep linking via hx-push-url
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/recipes/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/recipes", queryFilter)

	// Create status filter toggle as a search modifier
	statusFilter := components.CreateStatusFilter(
		"/api/recipes/search",
		status,
		"recipes-table",
	)

	// Use the new integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*mealplanningsvc.Recipe]{
		Title:             "Recipes",
		BaseSubtitle:      "Manage recipes",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search recipes...",
		HTMXSearchTarget:  "/api/recipes/search",
		Data:              recipesRes.Results,
		Actions:           []g.Node{},
		SearchModifiers:   []g.Node{statusFilter},
		TableOptions: &components.TableOptions[*mealplanningsvc.Recipe]{
			TableID: "recipes-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"id",
				"name",
				"source",
				"description",
				"status",
				"created_by_user",
				"created_at",
				"last_updated_at",
				"archived_at",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"created_at":      renderTimestamp,
				"last_updated_at": renderTimestamp,
				"archived_at":     renderTimestamp,
			},
			Pagination:             pagination,
			PaginationURLGenerator: paginationURLGenerator,
			DeepLinkURLGenerator:   deepLinkURLGenerator,
			PaginationHTMXTarget:   "#search-results",
		},
		RowLinkGenerator: func(data *mealplanningsvc.Recipe) string {
			return fmt.Sprintf("/recipes/%s", data.Id)
		},
		EmptyStateTitle:       "No recipes found",
		EmptyStateDescription: "Get started by creating your first recipe.",
		EmptyStateActions:     []g.Node{},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage recipes"
			}
			return fmt.Sprintf("Manage %d recipes", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Recipes", s.renderRecipesError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Recipes", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) RecipesSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Extract status from query params, default to "submitted"
	status := req.URL.Query().Get("status")
	if status == "" {
		status = "submitted"
	}

	recipesRes, err := c.GetRecipes(ctx, &mealplanningsvc.GetRecipesRequest{
		Filter: grpcFilter,
		Status: status,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading recipes: %v", err)),
			),
		), nil
	}

	// Build pagination from response
	pagination := buildPaginationFromGRPCResponse(recipesRes.Pagination)
	// Use search endpoint for pagination buttons to return just the table content
	// The main page URL is still used for deep linking via hx-push-url
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/recipes/search", queryFilter)
	deepLinkURLGenerator := buildPaginationURLGenerator(req, "/recipes", queryFilter)

	// Generate just the table (not the full page)
	if len(recipesRes.Results) == 0 {
		searchQuery := req.URL.Query().Get("search")
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No recipes found",
				fmt.Sprintf("No recipes match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{},
			),
		), nil
	}

	table, err := components.Table(recipesRes.Results, &components.TableOptions[*mealplanningsvc.Recipe]{
		TableID: "recipes-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"name",
			"source",
			"description",
			"status",
			"created_by_user",
			"created_at",
			"last_updated_at",
			"archived_at",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at":      renderTimestamp,
			"last_updated_at": renderTimestamp,
			"archived_at":     renderTimestamp,
		},
		Pagination:             pagination,
		PaginationURLGenerator: paginationURLGenerator,
		DeepLinkURLGenerator:   deepLinkURLGenerator,
		PaginationHTMXTarget:   "#search-results",
		RowLinkGenerator: func(data *mealplanningsvc.Recipe) string {
			return fmt.Sprintf("/recipes/%s", data.Id)
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

// renderRecipesError creates a consistent error display for the recipes page.
func (s *AdminFrontendServer) renderRecipesError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Recipes",
		Subtitle: "Manage recipes",
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
