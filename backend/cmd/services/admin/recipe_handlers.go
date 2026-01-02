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
		return page("Recipes", s.renderRecipesError("Error: No recipe MealPlanTaskID provided")), nil
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
					DisplayName: "MealPlanTaskID",
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
							// Generated step description
							ghtml.Div(
								ghtml.Class("text-sm text-gray-700 mb-1"),
								g.Text(generateStepDescription(step)),
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
							// Time and temperature (only show if there's meaningful data)
							g.If(
								(step.EstimatedTimeInSeconds != nil && step.EstimatedTimeInSeconds.Min != nil && *step.EstimatedTimeInSeconds.Min > 0) ||
									(step.TemperatureInCelsius != nil && step.TemperatureInCelsius.Min != nil && *step.TemperatureInCelsius.Min > 0),
								ghtml.Div(
									ghtml.Class("flex gap-4 text-xs text-gray-600 mb-2"),
									g.If(
										step.EstimatedTimeInSeconds != nil && step.EstimatedTimeInSeconds.Min != nil && *step.EstimatedTimeInSeconds.Min > 0,
										ghtml.Div(
											g.Text(func() string {
												if step.EstimatedTimeInSeconds == nil || step.EstimatedTimeInSeconds.Min == nil {
													return ""
												}
												minimum := *step.EstimatedTimeInSeconds.Min
												if step.EstimatedTimeInSeconds.Max != nil && *step.EstimatedTimeInSeconds.Max > 0 {
													maximum := *step.EstimatedTimeInSeconds.Max
													if maximum > minimum {
														return fmt.Sprintf("Time: %s - %s", formatDuration(minimum), formatDuration(maximum))
													} else if maximum == minimum {
														return fmt.Sprintf("Time: %s", formatDuration(minimum))
													}
												}
												return fmt.Sprintf("Time: %s", formatDuration(minimum))
											}()),
										),
									),
									g.If(
										step.TemperatureInCelsius != nil && step.TemperatureInCelsius.Min != nil && *step.TemperatureInCelsius.Min > 0,
										ghtml.Div(
											g.Text(func() string {
												if step.TemperatureInCelsius == nil || step.TemperatureInCelsius.Min == nil {
													return ""
												}
												minC := *step.TemperatureInCelsius.Min
												minF := (minC * 9.0 / 5.0) + 32.0

												if step.TemperatureInCelsius.Max != nil && *step.TemperatureInCelsius.Max > minC {
													maxC := *step.TemperatureInCelsius.Max
													maxF := (maxC * 9.0 / 5.0) + 32.0
													return fmt.Sprintf("Temperature: %.1f-%.1f°F (%.1f-%.1f°C)", minF, maxF, minC, maxC)
												}
												return fmt.Sprintf("Temperature: %.1f°F (%.1f°C)", minF, minC)
											}()),
										),
									),
								),
							),
						),
					),
					// 2x2 Grid: Instruments | Vessels, Ingredients | Products
					ghtml.Div(
						ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
						ghtml.Div(
							ghtml.Class("grid grid-cols-2 gap-4"),
							// Top row: Instruments
							ghtml.Div(
								ghtml.Div(
									ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
									g.Text("Instruments:"),
								),
								ghtml.Div(
									ghtml.Class("space-y-2"),
									func() g.Node {
										if len(step.Instruments) > 0 {
											return g.Group(s.renderStepInstruments(step.Instruments, steps))
										}
										return ghtml.Div(
											ghtml.Class("text-sm pl-4 border-l-2 border-gray-200 text-gray-500"),
											g.Text("none"),
										)
									}(),
								),
							),
							// Top row: Vessels
							ghtml.Div(
								ghtml.Div(
									ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
									g.Text("Vessels:"),
								),
								ghtml.Div(
									ghtml.Class("space-y-2"),
									func() g.Node {
										if len(step.Vessels) > 0 {
											return g.Group(s.renderStepVessels(step.Vessels, steps))
										}
										return ghtml.Div(
											ghtml.Class("text-sm pl-4 border-l-2 border-gray-200 text-gray-500"),
											g.Text("none"),
										)
									}(),
								),
							),
							// Bottom row: Ingredients
							ghtml.Div(
								g.If(len(step.Ingredients) > 0,
									ghtml.Div(
										ghtml.Div(
											ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
											g.Text("Ingredients:"),
										),
										ghtml.Div(
											ghtml.Class("space-y-2"),
											g.Group(s.renderStepIngredients(step.Ingredients, steps)),
										),
									),
								),
							),
							// Bottom row: Products
							ghtml.Div(
								g.If(len(step.Products) > 0,
									ghtml.Div(
										ghtml.Div(
											ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
											g.Text("Products:"),
										),
										ghtml.Div(
											ghtml.Class("space-y-2"),
											g.Group(s.renderStepProducts(step.Products)),
										),
									),
								),
							),
						),
					),
					// Completion conditions section
					g.If(len(step.CompletionConditions) > 0,
						ghtml.Div(
							ghtml.Class("border-t border-gray-200 pt-3 mt-3"),
							ghtml.Div(
								ghtml.Class("text-sm font-semibold text-gray-700 mb-2"),
								g.Text("Completion Conditions:"),
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

		// Get ingredient name - prioritize recipe step ingredient name over base ingredient name
		ingredientName := ""
		if ing.Name != "" {
			ingredientName = ing.Name
		} else if ing.Ingredient != nil {
			ingredientName = ing.Ingredient.Name
		}

		// Check if this ingredient comes from a previous step
		if ing.RecipeStepProductId != nil {
			productStep := s.findStepWithProduct(*ing.RecipeStepProductId, allSteps)
			if productStep != nil {
				// Format as "ingredient name from step X" without quantity
				details = append(details, ghtml.Span(
					ghtml.Class("font-medium"),
					g.Text(fmt.Sprintf("%s from step %d", ingredientName, productStep.Index+1)),
				))
			} else {
				// Fallback if step not found
				details = append(details,
					ghtml.Span(
						ghtml.Class("font-medium"),
						g.Text(ingredientName),
					),
					ghtml.Span(
						ghtml.Class("text-blue-600 ml-2"),
						g.Text(fmt.Sprintf("← Product MealPlanTaskID: %s", *ing.RecipeStepProductId)),
					),
				)
			}
		} else {
			// Regular ingredient - show name and quantity
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(ingredientName),
			))

			// MeasurementQuantity and unit
			if ing.Quantity != nil {
				qtyStr := ""
				if ing.Quantity.Min != 0 {
					qtyStr = formatQuantity(ing.Quantity.Min)
					if ing.Quantity.Max != nil {
						qtyStr += "-" + formatQuantity(*ing.Quantity.Max)
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
				g.Text(fmt.Sprintf("MeasurementQuantity notes: %s", ing.QuantityNotes)),
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

		// MeasurementQuantity (only show if not 1)
		if inst.Quantity != nil {
			qtyStr := ""
			if inst.Quantity.Min != 0 {
				// Only show quantity if it's not 1, or if there's a range
				if inst.Quantity.Min != 1 || (inst.Quantity.Max != nil && *inst.Quantity.Max != 1) {
					qtyStr = fmt.Sprintf("%d", inst.Quantity.Min)
					if inst.Quantity.Max != nil {
						qtyStr += fmt.Sprintf("-%d", *inst.Quantity.Max)
					}
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
					g.Text(fmt.Sprintf("← Product MealPlanTaskID: %s", *inst.RecipeStepProductId)),
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

		// Get vessel name
		vesselName := ""
		if vessel.Vessel != nil {
			vesselName = vessel.Vessel.Name
		} else if vessel.Name != "" {
			vesselName = vessel.Name
		}

		// Check if this vessel comes from a previous step
		if vessel.RecipeStepProductId != nil {
			productStep := s.findStepWithProduct(*vessel.RecipeStepProductId, allSteps)
			if productStep != nil {
				// Format as "vessel name from step X" without quantity
				details = append(details, ghtml.Span(
					ghtml.Class("font-medium"),
					g.Text(fmt.Sprintf("%s from step %d", vesselName, productStep.Index+1)),
				))
			} else {
				// Fallback if step not found
				details = append(details,
					ghtml.Span(
						ghtml.Class("font-medium"),
						g.Text(vesselName),
					),
					ghtml.Span(
						ghtml.Class("text-blue-600 ml-2"),
						g.Text(fmt.Sprintf("← Product MealPlanTaskID: %s", *vessel.RecipeStepProductId)),
					),
				)
			}
		} else {
			// Regular vessel - show name
			details = append(details, ghtml.Span(
				ghtml.Class("font-medium"),
				g.Text(vesselName),
			))
		}

		// MeasurementQuantity (only show if not 1)
		if vessel.Quantity != nil {
			qtyStr := ""
			if vessel.Quantity.Min != 0 {
				// Only show quantity if it's not 1, or if there's a range
				if vessel.Quantity.Min != 1 || (vessel.Quantity.Max != nil && *vessel.Quantity.Max != 1) {
					qtyStr = fmt.Sprintf("%d", vessel.Quantity.Min)
					if vessel.Quantity.Max != nil {
						qtyStr += fmt.Sprintf("-%d", *vessel.Quantity.Max)
					}
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

		// Quantity display: discrete vs continuous products
		isDiscrete := product.ItemQuantity != nil && (product.ItemQuantity.Min != nil || product.ItemQuantity.Max != nil)

		if isDiscrete {
			// Discrete product: Display "4 patties (4 oz each)" format
			itemQtyStr := ""
			if product.ItemQuantity.Min != nil {
				itemQtyStr = formatQuantity(*product.ItemQuantity.Min)
				if product.ItemQuantity.Max != nil {
					itemQtyStr += "-" + formatQuantity(*product.ItemQuantity.Max)
				}
			}

			measurementQtyStr := ""
			if product.MeasurementQuantity != nil && product.MeasurementQuantity.Min != nil {
				measurementQtyStr = formatQuantity(*product.MeasurementQuantity.Min)
				if product.MeasurementQuantity.Max != nil {
					measurementQtyStr += "-" + formatQuantity(*product.MeasurementQuantity.Max)
				}
			}

			if itemQtyStr != "" {
				unitName := formatMeasurementUnitName(product.MeasurementUnit)
				if measurementQtyStr != "" && unitName != "" {
					// Format: "4 patties (4 oz each)"
					details = append(details, ghtml.Span(
						ghtml.Class("text-gray-600 ml-2"),
						g.Text(fmt.Sprintf("%s (%s %s each)", itemQtyStr, measurementQtyStr, unitName)),
					))
				} else {
					// Fallback: just show count if measurement is missing
					details = append(details, ghtml.Span(
						ghtml.Class("text-gray-600 ml-2"),
						g.Text(itemQtyStr),
					))
				}
			}
		} else if product.MeasurementQuantity != nil {
			// Continuous product: Display "16 oz" format (backward compatible)
			qtyStr := ""
			if product.MeasurementQuantity.Min != nil {
				qtyStr = formatQuantity(*product.MeasurementQuantity.Min)
				if product.MeasurementQuantity.Max != nil {
					qtyStr += "-" + formatQuantity(*product.MeasurementQuantity.Max)
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
				g.Text(fmt.Sprintf("MeasurementQuantity notes: %s", product.QuantityNotes)),
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

// findStepWithProduct finds the step that contains a product with the given MealPlanTaskID.
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

// formatDuration formats seconds into a human-readable duration string (e.g., "1h 30m", "45m", "30s").
func formatDuration(seconds uint32) string {
	if seconds == 0 {
		return "0s"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if secs > 0 && hours == 0 {
		// Only show seconds if we don't have hours (to avoid clutter)
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}

	if len(parts) == 0 {
		return "0s"
	}
	return strings.Join(parts, " ")
}

// formatQuantity formats a float32 quantity, showing as integer if decimal part is zero.
func formatQuantity(qty float32) string {
	if qty == float32(int32(qty)) {
		return fmt.Sprintf("%d", int32(qty))
	}
	return fmt.Sprintf("%.2f", qty)
}

// humanizeList formats a list of strings with Oxford comma and "and" (e.g., "red, white, and blue").
func humanizeList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return fmt.Sprintf("%s and %s", items[0], items[1])
	}
	// For 3+ items, use Oxford comma: "a, b, and c"
	allButLast := strings.Join(items[:len(items)-1], ", ")
	return fmt.Sprintf("%s, and %s", allButLast, items[len(items)-1])
}

// generateStepDescription generates a human-readable description of a recipe step.
func generateStepDescription(step *mealplanningsvc.RecipeStep) string {
	// Collect instruments
	var instruments []string
	for _, inst := range step.Instruments {
		name := ""
		if inst.Instrument != nil {
			name = inst.Instrument.Name
		} else if inst.Name != "" {
			name = inst.Name
		}
		if name != "" {
			instruments = append(instruments, name)
		}
	}

	// Collect ingredients
	var ingredients []string
	for _, ing := range step.Ingredients {
		name := ""
		if ing.Name != "" {
			name = ing.Name
		} else if ing.Ingredient != nil {
			name = ing.Ingredient.Name
		}
		if name != "" {
			ingredients = append(ingredients, name)
		}
	}

	// Collect vessels
	var vessels []string
	for _, vessel := range step.Vessels {
		name := ""
		if vessel.Vessel != nil {
			name = vessel.Vessel.Name
		} else if vessel.Name != "" {
			name = vessel.Name
		}
		if name != "" {
			vessels = append(vessels, name)
		}
	}

	// Collect products
	var products []string
	for _, product := range step.Products {
		if product.Name != "" {
			products = append(products, product.Name)
		}
	}

	// Build description following template: "using <instruments>, <preparation> <ingredients> in <vessels> to yield <products>"
	var descriptionParts []string

	// "using <instruments>"
	if len(instruments) > 0 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("using %s", humanizeList(instruments)))
	}

	// "<preparation>"
	prepName := ""
	if step.Preparation != nil {
		prepName = step.Preparation.Name
	}
	if prepName != "" {
		descriptionParts = append(descriptionParts, prepName)
	}

	// "<ingredients>"
	if len(ingredients) > 0 {
		descriptionParts = append(descriptionParts, humanizeList(ingredients))
	}

	// "in <vessels>"
	if len(vessels) > 0 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("in %s", humanizeList(vessels)))
	}

	// "to yield <products>"
	if len(products) > 0 {
		descriptionParts = append(descriptionParts, fmt.Sprintf("to yield %s", humanizeList(products)))
	}

	return strings.Join(descriptionParts, " ")
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
