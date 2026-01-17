package main

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

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

	// Create steps section (includes sub-recipe steps inline)
	stepsSection := components.CardWithHeader(
		"Recipe Steps",
		&design.StandardPalette,
		nil,
		s.renderRecipeSteps(recipe.Steps, recipe.AssociatedRecipes),
	)

	// Create aggregated ingredients section
	aggregatedIngredients := s.getAggregatedIngredients(recipe.Steps, recipe.AssociatedRecipes)
	ingredientsSection := components.CardWithHeader(
		"Aggregated Ingredients",
		&design.StandardPalette,
		nil,
		s.renderAggregatedIngredients(aggregatedIngredients),
	)

	// Create aggregated instruments and vessels sections
	aggregatedInstrumentsVessels := s.getAggregatedInstrumentsAndVessels(recipe.Steps, recipe.AssociatedRecipes)

	// Separate instruments and vessels
	var instruments []*AggregatedInstrumentVessel
	var vessels []*AggregatedInstrumentVessel
	for _, item := range aggregatedInstrumentsVessels {
		if item.Type == "instrument" {
			instruments = append(instruments, item)
		} else {
			vessels = append(vessels, item)
		}
	}

	var sectionNodes []g.Node
	sectionNodes = append(sectionNodes, formPageResult.Node, ingredientsSection)

	// Only add instruments section if there are instruments
	if len(instruments) > 0 {
		instrumentsSection := components.CardWithHeader(
			"Aggregated Instruments",
			&design.StandardPalette,
			nil,
			s.renderAggregatedInstrumentsAndVessels(instruments),
		)
		sectionNodes = append(sectionNodes, instrumentsSection)
	}

	// Only add vessels section if there are vessels (referenced more than twice)
	if len(vessels) > 0 {
		vesselsSection := components.CardWithHeader(
			"Aggregated Vessels",
			&design.StandardPalette,
			nil,
			s.renderAggregatedInstrumentsAndVessels(vessels),
		)
		sectionNodes = append(sectionNodes, vesselsSection)
	}

	sectionNodes = append(sectionNodes, stepsSection)

	// Combine form, steps, and aggregated lists
	return page("Recipes",
		ghtml.Div(
			ghtml.Class("space-y-6"),
			g.Group(sectionNodes),
		),
	), nil
}

// buildProductUsageMap builds a map of product ID -> list of step indices (1-based) that use the product.
func (s *AdminFrontendServer) buildProductUsageMap(steps []*mealplanningsvc.RecipeStep) map[string][]uint32 {
	usageMap := make(map[string][]uint32)

	for _, step := range steps {
		// Check ingredients for product references
		for _, ing := range step.Ingredients {
			if ing.RecipeStepProductId != nil && *ing.RecipeStepProductId != "" {
				productID := *ing.RecipeStepProductId
				usageMap[productID] = append(usageMap[productID], step.Index+1)
			}
		}
		// Check instruments for product references
		for _, inst := range step.Instruments {
			if inst.RecipeStepProductId != nil && *inst.RecipeStepProductId != "" {
				productID := *inst.RecipeStepProductId
				usageMap[productID] = append(usageMap[productID], step.Index+1)
			}
		}
		// Check vessels for product references
		for _, vessel := range step.Vessels {
			if vessel.RecipeStepProductId != nil && *vessel.RecipeStepProductId != "" {
				productID := *vessel.RecipeStepProductId
				usageMap[productID] = append(usageMap[productID], step.Index+1)
			}
		}
	}

	return usageMap
}

// AggregatedIngredient represents an aggregated ingredient with total quantities.
type AggregatedIngredient struct {
	IngredientID     string
	Name             string
	Quantity         *types.Float32RangeWithOptionalMax
	QuantityNotes    string
	MeasurementUnit  *mealplanningsvc.ValidMeasurementUnit
	SourceRecipeID   string
	SourceRecipeName string
}

// AggregatedInstrumentVessel represents an aggregated instrument or vessel with total quantities.
type AggregatedInstrumentVessel struct {
	QuantityUint32   *types.Uint32RangeWithOptionalMax
	ItemID           string
	Name             string
	Type             string
	SourceRecipeID   string
	SourceRecipeName string
	ReferenceCount   uint32
}

// getAggregatedIngredients aggregates all ingredients from recipe steps and associated recipes.
func (s *AdminFrontendServer) getAggregatedIngredients(steps []*mealplanningsvc.RecipeStep, associatedRecipes []*mealplanningsvc.Recipe) []*AggregatedIngredient {
	aggregated := make(map[string]*AggregatedIngredient)

	// Process main recipe steps
	for _, step := range steps {
		for _, ing := range step.Ingredients {
			if ing.Ingredient == nil {
				continue
			}
			ingredientID := ing.Ingredient.Id
			if ingredientID == "" {
				continue
			}

			if aggregated[ingredientID] == nil {
				aggregated[ingredientID] = &AggregatedIngredient{
					IngredientID:    ingredientID,
					Name:            ing.Name,
					QuantityNotes:   ing.QuantityNotes,
					MeasurementUnit: ing.MeasurementUnit,
				}
			}

			// Add quantity if present
			if ing.Quantity != nil {
				if aggregated[ingredientID].Quantity == nil {
					aggregated[ingredientID].Quantity = &types.Float32RangeWithOptionalMax{
						Min: ing.Quantity.Min,
						Max: ing.Quantity.Max,
					}
				} else {
					aggregated[ingredientID].Quantity.Min += ing.Quantity.Min
					if ing.Quantity.Max != nil {
						if aggregated[ingredientID].Quantity.Max != nil {
							*aggregated[ingredientID].Quantity.Max += *ing.Quantity.Max
						} else {
							maximum := *ing.Quantity.Max
							aggregated[ingredientID].Quantity.Max = &maximum
						}
					} else {
						aggregated[ingredientID].Quantity.Max = nil
					}
				}
			}
		}
	}

	// Process associated recipes
	for _, associatedRecipe := range associatedRecipes {
		for _, step := range associatedRecipe.Steps {
			for _, ing := range step.Ingredients {
				if ing.Ingredient == nil {
					continue
				}
				ingredientID := ing.Ingredient.Id
				if ingredientID == "" {
					continue
				}

				if aggregated[ingredientID] == nil {
					aggregated[ingredientID] = &AggregatedIngredient{
						IngredientID:     ingredientID,
						Name:             ing.Name,
						QuantityNotes:    ing.QuantityNotes,
						MeasurementUnit:  ing.MeasurementUnit,
						SourceRecipeID:   associatedRecipe.Id,
						SourceRecipeName: associatedRecipe.Name,
					}
				}

				// Add quantity if present
				if ing.Quantity != nil {
					if aggregated[ingredientID].Quantity == nil {
						aggregated[ingredientID].Quantity = &types.Float32RangeWithOptionalMax{
							Min: ing.Quantity.Min,
							Max: ing.Quantity.Max,
						}
					} else {
						aggregated[ingredientID].Quantity.Min += ing.Quantity.Min
						if ing.Quantity.Max != nil {
							if aggregated[ingredientID].Quantity.Max != nil {
								*aggregated[ingredientID].Quantity.Max += *ing.Quantity.Max
							} else {
								maximum := *ing.Quantity.Max
								aggregated[ingredientID].Quantity.Max = &maximum
							}
						} else {
							aggregated[ingredientID].Quantity.Max = nil
						}
					}
				}
			}
		}
	}

	// Convert map to slice and sort by name
	result := slices.Collect(maps.Values(aggregated))
	slices.SortFunc(result, func(a, b *AggregatedIngredient) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

// getAggregatedInstrumentsAndVessels aggregates all instruments and vessels from recipe steps and associated recipes.
func (s *AdminFrontendServer) getAggregatedInstrumentsAndVessels(steps []*mealplanningsvc.RecipeStep, associatedRecipes []*mealplanningsvc.Recipe) []*AggregatedInstrumentVessel {
	aggregated := make(map[string]*AggregatedInstrumentVessel)

	// Process main recipe steps
	for _, step := range steps {
		// Process instruments
		for _, inst := range step.Instruments {
			if inst.Instrument == nil {
				continue
			}
			// Only include instruments that should be displayed in summary lists
			if !inst.Instrument.DisplayInSummaryLists {
				continue
			}
			itemID := inst.Instrument.Id
			if itemID == "" {
				continue
			}

			if aggregated[itemID] == nil {
				aggregated[itemID] = &AggregatedInstrumentVessel{
					ItemID:         itemID,
					Name:           inst.Name,
					Type:           "instrument",
					ReferenceCount: 0,
				}
			}

			// Increment reference count
			aggregated[itemID].ReferenceCount++

			// Add quantity if present
			if inst.Quantity != nil {
				if aggregated[itemID].QuantityUint32 == nil {
					aggregated[itemID].QuantityUint32 = &types.Uint32RangeWithOptionalMax{
						Min: inst.Quantity.Min,
						Max: inst.Quantity.Max,
					}
				} else {
					qty := aggregated[itemID].QuantityUint32
					qty.Min += inst.Quantity.Min
					if inst.Quantity.Max != nil {
						if qty.Max != nil {
							*qty.Max += *inst.Quantity.Max
						} else {
							maximum := *inst.Quantity.Max
							qty.Max = &maximum
						}
					} else {
						qty.Max = nil
					}
				}
			}
		}

		// Process vessels
		for _, vessel := range step.Vessels {
			if vessel.Vessel == nil {
				continue
			}
			// Only include vessels that should be displayed in summary lists
			if !vessel.Vessel.DisplayInSummaryLists {
				continue
			}
			itemID := vessel.Vessel.Id
			if itemID == "" {
				continue
			}

			if aggregated[itemID] == nil {
				aggregated[itemID] = &AggregatedInstrumentVessel{
					ItemID:         itemID,
					Name:           vessel.Name,
					Type:           "vessel",
					ReferenceCount: 0,
				}
			}

			// Increment reference count
			aggregated[itemID].ReferenceCount++

			// Add quantity if present
			// Note: protobuf Uint16RangeWithOptionalMax actually uses uint32 internally
			if vessel.Quantity != nil {
				if aggregated[itemID].QuantityUint32 == nil {
					aggregated[itemID].QuantityUint32 = &types.Uint32RangeWithOptionalMax{
						Min: vessel.Quantity.Min,
						Max: vessel.Quantity.Max,
					}
				} else {
					qty := aggregated[itemID].QuantityUint32
					qty.Min += vessel.Quantity.Min
					if vessel.Quantity.Max != nil {
						if qty.Max != nil {
							*qty.Max += *vessel.Quantity.Max
						} else {
							maximum := *vessel.Quantity.Max
							qty.Max = &maximum
						}
					} else {
						qty.Max = nil
					}
				}
			}
		}
	}

	// Process associated recipes
	for _, associatedRecipe := range associatedRecipes {
		for _, step := range associatedRecipe.Steps {
			// Process instruments
			for _, inst := range step.Instruments {
				if inst.Instrument == nil {
					continue
				}
				if !inst.Instrument.DisplayInSummaryLists {
					continue
				}
				itemID := inst.Instrument.Id
				if itemID == "" {
					continue
				}

				if aggregated[itemID] == nil {
					aggregated[itemID] = &AggregatedInstrumentVessel{
						ItemID:           itemID,
						Name:             inst.Name,
						Type:             "instrument",
						ReferenceCount:   0,
						SourceRecipeID:   associatedRecipe.Id,
						SourceRecipeName: associatedRecipe.Name,
					}
				}

				// Increment reference count
				aggregated[itemID].ReferenceCount++

				// Add quantity if present
				if inst.Quantity != nil {
					if aggregated[itemID].QuantityUint32 == nil {
						aggregated[itemID].QuantityUint32 = &types.Uint32RangeWithOptionalMax{
							Min: inst.Quantity.Min,
							Max: inst.Quantity.Max,
						}
					} else {
						qty := aggregated[itemID].QuantityUint32
						qty.Min += inst.Quantity.Min
						if inst.Quantity.Max != nil {
							if qty.Max != nil {
								*qty.Max += *inst.Quantity.Max
							} else {
								maximum := *inst.Quantity.Max
								qty.Max = &maximum
							}
						} else {
							qty.Max = nil
						}
					}
				}
			}

			// Process vessels
			for _, vessel := range step.Vessels {
				if vessel.Vessel == nil {
					continue
				}
				if !vessel.Vessel.DisplayInSummaryLists {
					continue
				}
				itemID := vessel.Vessel.Id
				if itemID == "" {
					continue
				}

				if aggregated[itemID] == nil {
					aggregated[itemID] = &AggregatedInstrumentVessel{
						ItemID:           itemID,
						Name:             vessel.Name,
						Type:             "vessel",
						ReferenceCount:   0,
						SourceRecipeID:   associatedRecipe.Id,
						SourceRecipeName: associatedRecipe.Name,
					}
				}

				// Increment reference count
				aggregated[itemID].ReferenceCount++

				// Add quantity if present
				// Note: protobuf Uint16RangeWithOptionalMax actually uses uint32 internally
				if vessel.Quantity != nil {
					if aggregated[itemID].QuantityUint32 == nil {
						aggregated[itemID].QuantityUint32 = &types.Uint32RangeWithOptionalMax{
							Min: vessel.Quantity.Min,
							Max: vessel.Quantity.Max,
						}
					} else {
						qty := aggregated[itemID].QuantityUint32
						qty.Min += vessel.Quantity.Min
						if vessel.Quantity.Max != nil {
							if qty.Max != nil {
								*qty.Max += *vessel.Quantity.Max
							} else {
								maximum := *vessel.Quantity.Max
								qty.Max = &maximum
							}
						} else {
							qty.Max = nil
						}
					}
				}
			}
		}
	}

	// Convert map to slice, filter vessels to only those referenced more than twice, and sort by name
	result := make([]*AggregatedInstrumentVessel, 0, len(aggregated))
	for _, item := range aggregated {
		// For vessels, only include if referenced more than twice
		// For instruments, include all
		if item.Type == "vessel" {
			if item.ReferenceCount > 2 {
				result = append(result, item)
			}
		} else {
			// Include all instruments
			result = append(result, item)
		}
	}
	slices.SortFunc(result, func(a, b *AggregatedInstrumentVessel) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

// renderAggregatedIngredients renders the aggregated ingredients list.
func (s *AdminFrontendServer) renderAggregatedIngredients(ingredients []*AggregatedIngredient) g.Node {
	if len(ingredients) == 0 {
		return ghtml.Div(
			ghtml.Class("text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No ingredients found."),
			),
		)
	}

	var items []g.Node
	for _, ing := range ingredients {
		var quantityText string
		if ing.Quantity != nil {
			if ing.Quantity.Max != nil {
				quantityText = fmt.Sprintf("%.2f - %.2f", ing.Quantity.Min, *ing.Quantity.Max)
			} else {
				quantityText = fmt.Sprintf("%.2f+", ing.Quantity.Min)
			}
			if ing.MeasurementUnit != nil {
				quantityText += " " + ing.MeasurementUnit.Name
			}
		}

		var sourceText string
		if ing.SourceRecipeName != "" {
			sourceText = fmt.Sprintf(" (from: %s)", ing.SourceRecipeName)
		}

		itemContent := []g.Node{
			ghtml.Div(
				ghtml.Class("font-medium text-gray-900"),
				g.Text(ing.Name),
			),
		}

		if quantityText != "" {
			itemContent = append(itemContent, ghtml.Div(
				ghtml.Class("text-sm text-gray-600"),
				g.Text(quantityText),
			))
		}

		if ing.QuantityNotes != "" {
			itemContent = append(itemContent, ghtml.Div(
				ghtml.Class("text-xs text-gray-500 italic"),
				g.Text(ing.QuantityNotes),
			))
		}

		if sourceText != "" {
			itemContent = append(itemContent, ghtml.Div(
				ghtml.Class("text-xs text-purple-600"),
				g.Text(sourceText),
			))
		}

		items = append(items, ghtml.Div(
			ghtml.Class("py-2 border-b border-gray-200 last:border-b-0"),
			g.Group(itemContent),
		))
	}

	return ghtml.Div(
		ghtml.Class("space-y-0"),
		g.Group(items),
	)
}

// renderAggregatedInstrumentsAndVessels renders the aggregated instruments and vessels list.
func (s *AdminFrontendServer) renderAggregatedInstrumentsAndVessels(items []*AggregatedInstrumentVessel) g.Node {
	if len(items) == 0 {
		return ghtml.Div(
			ghtml.Class("text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No instruments or vessels found."),
			),
		)
	}

	var itemNodes []g.Node
	for _, item := range items {
		// Show reference count (number of times referenced in steps)
		referenceText := fmt.Sprintf("referenced %d time", item.ReferenceCount)
		if item.ReferenceCount != 1 {
			referenceText += "s"
		}

		var sourceText string
		if item.SourceRecipeName != "" {
			sourceText = fmt.Sprintf(" (from: %s)", item.SourceRecipeName)
		}

		typeLabel := item.Type
		if typeLabel == "" {
			typeLabel = strings.ToUpper(typeLabel[:1]) + typeLabel[1:]
		}
		typeLabel += " • " + referenceText

		itemContent := []g.Node{
			ghtml.Div(
				ghtml.Class("font-medium text-gray-900"),
				g.Text(item.Name),
			),
			ghtml.Div(
				ghtml.Class("text-sm text-gray-600"),
				g.Text(typeLabel),
			),
		}

		if sourceText != "" {
			itemContent = append(itemContent, ghtml.Div(
				ghtml.Class("text-xs text-purple-600"),
				g.Text(sourceText),
			))
		}

		itemNodes = append(itemNodes, ghtml.Div(
			ghtml.Class("py-2 border-b border-gray-200 last:border-b-0"),
			g.Group(itemContent),
		))
	}

	return ghtml.Div(
		ghtml.Class("space-y-0"),
		g.Group(itemNodes),
	)
}

func (s *AdminFrontendServer) renderRecipeSteps(steps []*mealplanningsvc.RecipeStep, associatedRecipes []*mealplanningsvc.Recipe) g.Node {
	if len(steps) == 0 && len(associatedRecipes) == 0 {
		return ghtml.Div(
			ghtml.Class("text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("This recipe has no steps."),
			),
		)
	}

	// Build a map of product ID -> list of steps that use it (for showing "used in step X")
	// Include products from both main recipe and sub-recipes
	productUsageMap := s.buildProductUsageMap(steps)

	var allStepNodes []g.Node

	// First, render all sub-recipe steps (grouped by sub-recipe)
	for _, subRecipe := range associatedRecipes {
		if len(subRecipe.Steps) == 0 {
			continue
		}

		// Build product usage map for this sub-recipe's steps
		subRecipeProductUsageMap := s.buildProductUsageMap(subRecipe.Steps)
		// Also include usage from main recipe steps that reference sub-recipe products
		s.addMainRecipeProductUsage(subRecipeProductUsageMap, subRecipe.Id, steps)

		// Add a header for this sub-recipe
		allStepNodes = append(allStepNodes, ghtml.Div(
			ghtml.Class("border-l-4 border-purple-400 pl-4 py-2 mb-4 bg-purple-50/50 rounded-r"),
			ghtml.Div(
				ghtml.Class("flex items-center gap-2"),
				ghtml.Span(
					ghtml.Class("text-xs font-semibold text-purple-600 uppercase tracking-wide"),
					g.Text("Prerequisite:"),
				),
				ghtml.A(
					ghtml.Href(fmt.Sprintf("/recipes/%s", subRecipe.Id)),
					ghtml.Class("font-semibold text-purple-700 hover:text-purple-900 hover:underline"),
					g.Text(subRecipe.Name),
				),
			),
			g.If(subRecipe.Description != "",
				ghtml.P(
					ghtml.Class("text-xs text-gray-600 mt-1"),
					g.Text(subRecipe.Description),
				),
			),
		))

		// Render each step from this sub-recipe
		for i, step := range subRecipe.Steps {
			allStepNodes = append(allStepNodes, s.renderSubRecipeStep(step, i, subRecipe, subRecipeProductUsageMap))
		}

		// Add a visual separator before the main recipe steps
		allStepNodes = append(allStepNodes, ghtml.Div(
			ghtml.Class("my-6"),
		))
	}

	// If there are sub-recipes and main steps, add a header for the main recipe
	if len(associatedRecipes) > 0 && len(steps) > 0 {
		allStepNodes = append(allStepNodes, ghtml.Div(
			ghtml.Class("border-l-4 border-blue-400 pl-4 py-2 mb-4 bg-blue-50/50 rounded-r"),
			ghtml.Div(
				ghtml.Class("flex items-center gap-2"),
				ghtml.Span(
					ghtml.Class("text-xs font-semibold text-blue-600 uppercase tracking-wide"),
					g.Text("Main Recipe Steps"),
				),
			),
		))
	}

	// Render the main recipe steps
	for i, step := range steps {
		allStepNodes = append(allStepNodes, s.renderSingleRecipeStep(step, i, steps, productUsageMap))
	}

	return ghtml.Div(
		ghtml.Class("space-y-4"),
		g.Group(allStepNodes),
	)
}

// addMainRecipeProductUsage adds usage information for sub-recipe products that are used in main recipe steps.
func (s *AdminFrontendServer) addMainRecipeProductUsage(usageMap map[string][]uint32, subRecipeID string, mainSteps []*mealplanningsvc.RecipeStep) {
	for _, step := range mainSteps {
		for _, ing := range step.Ingredients {
			if ing.RecipeStepProductRecipeId != nil && *ing.RecipeStepProductRecipeId == subRecipeID {
				if ing.RecipeStepProductId != nil && *ing.RecipeStepProductId != "" {
					productID := *ing.RecipeStepProductId
					// Add a special marker for main recipe usage (use step index + 1000 to distinguish)
					usageMap[productID] = append(usageMap[productID], step.Index+1)
				}
			}
		}
	}
}

// renderSubRecipeStep renders a single step from a sub-recipe with distinct styling.
func (s *AdminFrontendServer) renderSubRecipeStep(step *mealplanningsvc.RecipeStep, stepIndex int, subRecipe *mealplanningsvc.Recipe, productUsageMap map[string][]uint32) g.Node {
	return ghtml.Div(
		ghtml.Class("border border-purple-200 rounded-lg p-4 mb-4 last:mb-0 bg-white"),
		ghtml.Div(
			ghtml.Class("space-y-4"),
			// Step header
			ghtml.Div(
				ghtml.Class("flex items-start gap-4"),
				// Step number (with sub-recipe indicator)
				ghtml.Div(
					ghtml.Class("flex-shrink-0 w-10 h-10 rounded-full bg-purple-100 text-purple-800 flex items-center justify-center font-semibold text-sm"),
					g.Text(fmt.Sprintf("%d", stepIndex+1)),
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
					// Show explicit instructions if present, otherwise show generated description
					g.If(step.ExplicitInstructions != "",
						ghtml.Div(
							ghtml.Class("text-sm text-gray-700 mb-1"),
							g.Text(step.ExplicitInstructions),
						),
					),
					g.If(step.ExplicitInstructions == "",
						ghtml.Div(
							ghtml.Class("text-sm text-gray-700 mb-1"),
							g.Text(generateStepDescription(step)),
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
					s.renderStepTimeAndTemperature(step),
				),
			),
			// 2x2 Grid: Instruments | Vessels, Ingredients | Products
			ghtml.Div(
				ghtml.Class("border-t border-purple-200 pt-3 mt-3"),
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
									return g.Group(s.renderStepInstruments(step.Instruments, subRecipe.Steps))
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
									return g.Group(s.renderStepVessels(step.Vessels, subRecipe.Steps))
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
									g.Group(s.renderStepIngredients(step.Ingredients, subRecipe.Steps)),
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
									g.Group(s.renderStepProducts(step.Products, productUsageMap)),
								),
							),
						),
					),
				),
			),
			// Completion conditions section
			g.If(len(step.CompletionConditions) > 0,
				ghtml.Div(
					ghtml.Class("border-t border-purple-200 pt-3 mt-3"),
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
	)
}

// renderSingleRecipeStep renders a single recipe step with all its components.
func (s *AdminFrontendServer) renderSingleRecipeStep(step *mealplanningsvc.RecipeStep, stepIndex int, allSteps []*mealplanningsvc.RecipeStep, productUsageMap map[string][]uint32) g.Node {
	return ghtml.Div(
		ghtml.Class("border border-gray-200 rounded-lg p-4 mb-4 last:mb-0 bg-white"),
		ghtml.Div(
			ghtml.Class("space-y-4"),
			// Step header
			ghtml.Div(
				ghtml.Class("flex items-start gap-4"),
				// Step number
				ghtml.Div(
					ghtml.Class("flex-shrink-0 w-10 h-10 rounded-full bg-blue-100 text-blue-800 flex items-center justify-center font-semibold text-sm"),
					g.Text(fmt.Sprintf("%d", stepIndex+1)),
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
					// Show explicit instructions if present, otherwise show generated description
					g.If(step.ExplicitInstructions != "",
						ghtml.Div(
							ghtml.Class("text-sm text-gray-700 mb-1"),
							g.Text(step.ExplicitInstructions),
						),
					),
					g.If(step.ExplicitInstructions == "",
						ghtml.Div(
							ghtml.Class("text-sm text-gray-700 mb-1"),
							g.Text(generateStepDescription(step)),
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
					s.renderStepTimeAndTemperature(step),
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
									return g.Group(s.renderStepInstruments(step.Instruments, allSteps))
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
									return g.Group(s.renderStepVessels(step.Vessels, allSteps))
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
									g.Group(s.renderStepIngredients(step.Ingredients, allSteps)),
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
									g.Group(s.renderStepProducts(step.Products, productUsageMap)),
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
	)
}

// renderStepTimeAndTemperature renders the time and temperature for a step.
func (s *AdminFrontendServer) renderStepTimeAndTemperature(step *mealplanningsvc.RecipeStep) g.Node {
	hasTime := step.EstimatedTimeInSeconds != nil && step.EstimatedTimeInSeconds.Min != nil && *step.EstimatedTimeInSeconds.Min > 0
	hasTemp := step.TemperatureInCelsius != nil && step.TemperatureInCelsius.Min != nil && *step.TemperatureInCelsius.Min > 0

	if !hasTime && !hasTemp {
		return nil
	}

	return ghtml.Div(
		ghtml.Class("flex gap-4 text-xs text-gray-600 mb-2"),
		g.If(hasTime,
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
		g.If(hasTemp,
			ghtml.Div(
				g.Text(func() string {
					if step.TemperatureInCelsius == nil || step.TemperatureInCelsius.Min == nil {
						return ""
					}
					minC := roundTemperatureToNearest5(*step.TemperatureInCelsius.Min)
					minF := roundTemperatureToNearest5((minC * 9.0 / 5.0) + 32.0)

					if step.TemperatureInCelsius.Max != nil && *step.TemperatureInCelsius.Max > *step.TemperatureInCelsius.Min {
						maxC := roundTemperatureToNearest5(*step.TemperatureInCelsius.Max)
						maxF := roundTemperatureToNearest5((maxC * 9.0 / 5.0) + 32.0)
						return fmt.Sprintf("Temperature: %.0f-%.0f°F (%.0f-%.0f°C)", minF, maxF, minC, maxC)
					}
					return fmt.Sprintf("Temperature: %.0f°F (%.0f°C)", minF, minC)
				}()),
			),
		),
	)
}

func (s *AdminFrontendServer) renderStepIngredients(ingredients []*mealplanningsvc.RecipeStepIngredient, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	// Group ingredients by Index
	indexGroups := make(map[uint32][]*mealplanningsvc.RecipeStepIngredient)
	var indexOrder []uint32
	for _, ing := range ingredients {
		if _, exists := indexGroups[ing.Index]; !exists {
			indexOrder = append(indexOrder, ing.Index)
		}
		indexGroups[ing.Index] = append(indexGroups[ing.Index], ing)
	}

	var nodes []g.Node
	for _, idx := range indexOrder {
		group := indexGroups[idx]

		if len(group) == 1 {
			// Single ingredient - render normally
			nodes = append(nodes, s.renderSingleIngredient(group[0], allSteps, false))
		} else {
			// Multiple options - render as option group with "OR"
			nodes = append(nodes, s.renderIngredientOptionGroup(group, allSteps))
		}
	}
	return nodes
}

// renderIngredientOptionGroup renders a group of ingredient options with "OR" between them.
func (s *AdminFrontendServer) renderIngredientOptionGroup(options []*mealplanningsvc.RecipeStepIngredient, allSteps []*mealplanningsvc.RecipeStep) g.Node {
	var optionNodes []g.Node

	for i, ing := range options {
		optionNodes = append(optionNodes, s.renderSingleIngredient(ing, allSteps, true))

		// Add "OR" separator between options (but not after the last one)
		if i < len(options)-1 {
			optionNodes = append(optionNodes, ghtml.Div(
				ghtml.Class("text-xs font-semibold text-gray-500 uppercase tracking-wide py-1 pl-6"),
				g.Text("— or —"),
			))
		}
	}

	// Wrap in an indented container to show grouping
	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-amber-300 bg-amber-50/30 rounded-r py-1"),
		g.Group(optionNodes),
	)
}

// renderSingleIngredient renders a single ingredient item.
func (s *AdminFrontendServer) renderSingleIngredient(ing *mealplanningsvc.RecipeStepIngredient, allSteps []*mealplanningsvc.RecipeStep, isInOptionGroup bool) g.Node {
	var details []g.Node

	// Get ingredient name - prioritize recipe step ingredient name over base ingredient name
	ingredientName := ""
	if ing.Name != "" {
		ingredientName = ing.Name
	} else if ing.Ingredient != nil {
		ingredientName = ing.Ingredient.Name
	}

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

	// Check if this ingredient comes from a previous step
	if ing.RecipeStepProductId != nil {
		productStep := s.findStepWithProduct(*ing.RecipeStepProductId, allSteps)
		if productStep != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("inline-block px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded ml-2"),
				g.Text(fmt.Sprintf("← from step %d", productStep.Index+1)),
			))
		} else {
			// Fallback if step not found
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

	// Notes
	if ing.IngredientNotes != "" {
		details = append(details, ghtml.Div(
			ghtml.Class("text-xs text-gray-500 mt-1"),
			g.Text(fmt.Sprintf("Ingredient notes: %s", ing.IngredientNotes)),
		))
	}
	// Only show quantity notes if this is not a product from a prior step
	if ing.RecipeStepProductId == nil && ing.QuantityNotes != "" {
		details = append(details, ghtml.Div(
			ghtml.Class("text-xs text-gray-500 mt-1"),
			g.Text(fmt.Sprintf("MeasurementQuantity notes: %s", ing.QuantityNotes)),
		))
	}

	// Use different styling if inside an option group vs standalone
	if isInOptionGroup {
		return ghtml.Div(
			ghtml.Class("py-1 pl-2"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
		ghtml.Div(
			ghtml.Class("flex flex-wrap items-center gap-1"),
			g.Group(details),
			g.Group(flags),
		),
	)
}

func (s *AdminFrontendServer) renderStepInstruments(instruments []*mealplanningsvc.RecipeStepInstrument, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	// Group instruments by Index
	indexGroups := make(map[uint32][]*mealplanningsvc.RecipeStepInstrument)
	var indexOrder []uint32
	for _, inst := range instruments {
		if _, exists := indexGroups[inst.Index]; !exists {
			indexOrder = append(indexOrder, inst.Index)
		}
		indexGroups[inst.Index] = append(indexGroups[inst.Index], inst)
	}

	var nodes []g.Node
	for _, idx := range indexOrder {
		group := indexGroups[idx]

		if len(group) == 1 {
			// Single instrument - render normally
			nodes = append(nodes, s.renderSingleInstrument(group[0], allSteps, false))
		} else {
			// Multiple options - render as option group with "OR"
			nodes = append(nodes, s.renderInstrumentOptionGroup(group, allSteps))
		}
	}
	return nodes
}

// renderInstrumentOptionGroup renders a group of instrument options with "OR" between them.
func (s *AdminFrontendServer) renderInstrumentOptionGroup(options []*mealplanningsvc.RecipeStepInstrument, allSteps []*mealplanningsvc.RecipeStep) g.Node {
	var optionNodes []g.Node

	for i, inst := range options {
		optionNodes = append(optionNodes, s.renderSingleInstrument(inst, allSteps, true))

		// Add "OR" separator between options (but not after the last one)
		if i < len(options)-1 {
			optionNodes = append(optionNodes, ghtml.Div(
				ghtml.Class("text-xs font-semibold text-gray-500 uppercase tracking-wide py-1 pl-6"),
				g.Text("— or —"),
			))
		}
	}

	// Wrap in an indented container to show grouping
	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-amber-300 bg-amber-50/30 rounded-r py-1"),
		g.Group(optionNodes),
	)
}

// renderSingleInstrument renders a single instrument item.
func (s *AdminFrontendServer) renderSingleInstrument(inst *mealplanningsvc.RecipeStepInstrument, allSteps []*mealplanningsvc.RecipeStep, isInOptionGroup bool) g.Node {
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
				ghtml.Class("inline-block px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded ml-2"),
				g.Text(fmt.Sprintf("← from step %d", productStep.Index+1)),
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

	// Use different styling if inside an option group vs standalone
	if isInOptionGroup {
		return ghtml.Div(
			ghtml.Class("py-1 pl-2"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
		ghtml.Div(
			ghtml.Class("flex flex-wrap items-center gap-1"),
			g.Group(details),
			g.Group(flags),
		),
	)
}

func (s *AdminFrontendServer) renderStepVessels(vessels []*mealplanningsvc.RecipeStepVessel, allSteps []*mealplanningsvc.RecipeStep) []g.Node {
	// Group vessels by Index
	indexGroups := make(map[uint32][]*mealplanningsvc.RecipeStepVessel)
	var indexOrder []uint32
	for _, vessel := range vessels {
		if _, exists := indexGroups[vessel.Index]; !exists {
			indexOrder = append(indexOrder, vessel.Index)
		}
		indexGroups[vessel.Index] = append(indexGroups[vessel.Index], vessel)
	}

	var nodes []g.Node
	for _, idx := range indexOrder {
		group := indexGroups[idx]

		if len(group) == 1 {
			// Single vessel - render normally
			nodes = append(nodes, s.renderSingleVessel(group[0], allSteps, false))
		} else {
			// Multiple options - render as option group with "OR"
			nodes = append(nodes, s.renderVesselOptionGroup(group, allSteps))
		}
	}
	return nodes
}

// renderVesselOptionGroup renders a group of vessel options with "OR" between them.
func (s *AdminFrontendServer) renderVesselOptionGroup(options []*mealplanningsvc.RecipeStepVessel, allSteps []*mealplanningsvc.RecipeStep) g.Node {
	var optionNodes []g.Node

	for i, vessel := range options {
		optionNodes = append(optionNodes, s.renderSingleVessel(vessel, allSteps, true))

		// Add "OR" separator between options (but not after the last one)
		if i < len(options)-1 {
			optionNodes = append(optionNodes, ghtml.Div(
				ghtml.Class("text-xs font-semibold text-gray-500 uppercase tracking-wide py-1 pl-6"),
				g.Text("— or —"),
			))
		}
	}

	// Wrap in an indented container to show grouping
	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-amber-300 bg-amber-50/30 rounded-r py-1"),
		g.Group(optionNodes),
	)
}

// renderSingleVessel renders a single vessel item.
func (s *AdminFrontendServer) renderSingleVessel(vessel *mealplanningsvc.RecipeStepVessel, allSteps []*mealplanningsvc.RecipeStep, isInOptionGroup bool) g.Node {
	var details []g.Node

	// Get vessel name
	vesselName := ""
	if vessel.Vessel != nil {
		vesselName = vessel.Vessel.Name
	} else if vessel.Name != "" {
		vesselName = vessel.Name
	}

	// Regular vessel - show name
	details = append(details, ghtml.Span(
		ghtml.Class("font-medium"),
		g.Text(vesselName),
	))

	// Check if this vessel comes from a previous step
	if vessel.RecipeStepProductId != nil {
		productStep := s.findStepWithProduct(*vessel.RecipeStepProductId, allSteps)
		if productStep != nil {
			details = append(details, ghtml.Span(
				ghtml.Class("inline-block px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded ml-2"),
				g.Text(fmt.Sprintf("← from step %d", productStep.Index+1)),
			))
		} else {
			// Fallback if step not found
			details = append(details, ghtml.Span(
				ghtml.Class("text-blue-600 ml-2"),
				g.Text(fmt.Sprintf("← Product ID: %s", *vessel.RecipeStepProductId)),
			))
		}
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

	// Use different styling if inside an option group vs standalone
	if isInOptionGroup {
		return ghtml.Div(
			ghtml.Class("py-1 pl-2"),
			ghtml.Div(
				ghtml.Class("flex flex-wrap items-center gap-1"),
				g.Group(details),
				g.Group(flags),
			),
		)
	}

	return ghtml.Div(
		ghtml.Class("text-sm pl-4 border-l-2 border-gray-200"),
		ghtml.Div(
			ghtml.Class("flex flex-wrap items-center gap-1"),
			g.Group(details),
			g.Group(flags),
		),
	)
}

func (s *AdminFrontendServer) renderStepProducts(products []*mealplanningsvc.RecipeStepProduct, usageMap map[string][]uint32) []g.Node {
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

		// Show which later steps use this product
		if usedInSteps, ok := usageMap[product.Id]; ok && len(usedInSteps) > 0 {
			stepList := make([]string, len(usedInSteps))
			for i, stepNum := range usedInSteps {
				stepList[i] = fmt.Sprintf("%d", stepNum)
			}
			details = append(details, ghtml.Span(
				ghtml.Class("inline-block px-2 py-0.5 text-xs bg-green-100 text-green-800 rounded ml-2"),
				g.Text(fmt.Sprintf("→ used in step %s", strings.Join(stepList, ", "))),
			))
		}

		// Quantity display: discrete vs continuous products
		isDiscrete := product.ItemQuantity != nil && (product.ItemQuantity.Min != nil || product.ItemQuantity.Max != nil)

		if isDiscrete {
			// Discrete product: Display "4 fillets (6 oz each)" format
			itemQtyStr := ""
			var itemQtyMin, itemQtyMax float32
			if product.ItemQuantity.Min != nil {
				itemQtyMin = *product.ItemQuantity.Min
				itemQtyStr = formatQuantity(itemQtyMin)
				if product.ItemQuantity.Max != nil {
					itemQtyMax = *product.ItemQuantity.Max
					itemQtyStr += "-" + formatQuantity(itemQtyMax)
				} else {
					itemQtyMax = itemQtyMin
				}
			}

			measurementQtyStr := ""
			if product.MeasurementQuantity != nil && product.MeasurementQuantity.Min != nil && itemQtyMin > 0 {
				// Calculate per-item quantity by dividing total by item count
				totalMin := *product.MeasurementQuantity.Min
				perItemMin := totalMin / itemQtyMin
				measurementQtyStr = formatQuantity(perItemMin)

				if product.MeasurementQuantity.Max != nil {
					totalMax := *product.MeasurementQuantity.Max
					if product.ItemQuantity.Max != nil && itemQtyMax > 0 {
						perItemMax := totalMax / itemQtyMax
						measurementQtyStr += "-" + formatQuantity(perItemMax)
					} else if itemQtyMin > 0 {
						perItemMax := totalMax / itemQtyMin
						measurementQtyStr += "-" + formatQuantity(perItemMax)
					}
				}
			}

			if itemQtyStr != "" {
				unitName := formatMeasurementUnitName(product.MeasurementUnit)
				if measurementQtyStr != "" && unitName != "" {
					// Format: "4 fillets (6 oz each)"
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
				minTemp := roundTemperatureToNearest5(*product.StorageTemperatureInCelsius.Min)
				tempStr = fmt.Sprintf("%.0f", minTemp)
				if product.StorageTemperatureInCelsius.Max != nil {
					maxTemp := roundTemperatureToNearest5(*product.StorageTemperatureInCelsius.Max)
					tempStr += fmt.Sprintf("-%.0f", maxTemp)
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

// roundTemperatureToNearest5 rounds a temperature to the nearest 5.
func roundTemperatureToNearest5(temp float32) float32 {
	return float32(int((temp+2.5)/5.0) * 5)
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
