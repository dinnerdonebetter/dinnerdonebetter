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

// ValidIngredientMeasurementUnitsForIngredient lists all measurement units associated with an ingredient.
func (s *AdminFrontendServer) ValidIngredientMeasurementUnitsForIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	ingredientID := s.validIngredientIDRouteParamFetcher(req)
	if ingredientID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Measurement Units",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No measurement units associated with this ingredient.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Measurement Units",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidIngredientMeasurementUnitsByIngredient(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitsByIngredientRequest{
		ValidIngredientId: ingredientID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Measurement Units",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]*components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.MeasurementUnit != nil {
			items = append(items, &components.AssociationItem{
				ID:          assoc.Id,
				Name:        assoc.MeasurementUnit.Name,
				Description: assoc.MeasurementUnit.Description,
				Notes:       assoc.Notes,
			})
		}
	}

	return components.AssociationList(&components.AssociationListProps{
		Title:                "Associated Measurement Units",
		Palette:              &design.StandardPalette,
		Items:                items,
		EntityID:             ingredientID,
		AddSearchPlaceholder: "Search for measurement units to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_ingredients/%s/measurement_units/search", ingredientID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_ingredients/%s/measurement_units", ingredientID),
		DeleteEndpoint:       "/api/valid_ingredient_measurement_units",
		NoItemsMessage:       "No measurement units associated with this ingredient.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// ValidIngredientMeasurementUnitsForMeasurementUnit lists all ingredients associated with a measurement unit.
func (s *AdminFrontendServer) ValidIngredientMeasurementUnitsForMeasurementUnit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)
	if measurementUnitID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Ingredients",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No ingredients associated with this measurement unit.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Ingredients",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx, &mealplanningsvc.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest{
		ValidMeasurementUnitId: measurementUnitID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Ingredients",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]*components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.Ingredient != nil {
			items = append(items, &components.AssociationItem{
				ID:          assoc.Id,
				Name:        assoc.Ingredient.Name,
				Description: assoc.Ingredient.Description,
				Notes:       assoc.Notes,
			})
		}
	}

	return components.AssociationList(&components.AssociationListProps{
		Title:                "Associated Ingredients",
		Palette:              &design.StandardPalette,
		Items:                items,
		EntityID:             measurementUnitID,
		AddSearchPlaceholder: "Search for ingredients to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_measurement_units/%s/ingredients/search", measurementUnitID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_measurement_units/%s/ingredients", measurementUnitID),
		DeleteEndpoint:       "/api/valid_ingredient_measurement_units",
		NoItemsMessage:       "No ingredients associated with this measurement unit.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// SearchMeasurementUnitsForIngredient searches for measurement units to add to an ingredient.
func (s *AdminFrontendServer) SearchMeasurementUnitsForIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	ingredientID := s.validIngredientIDRouteParamFetcher(req)
	query := req.URL.Query().Get("q")

	if query == "" {
		return ghtml.Div(), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Search for measurement units
	searchRes, err := c.GetValidMeasurementUnits(ctx, &mealplanningsvc.GetValidMeasurementUnitsRequest{
		Filter: nil, // TODO: Add filtering by query
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error searching: %v", err)),
		), nil
	}

	// Filter results by query (simple substring match)
	var results []*components.SearchResultItem
	for _, unit := range searchRes.Results {
		if contains(unit.Name, query) || contains(unit.Description, query) {
			results = append(results, &components.SearchResultItem{
				ID:          unit.Id,
				Name:        unit.Name,
				Description: unit.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_ingredients/%s/measurement_units", ingredientID),
		HTMXTarget:     "#association-list-container",
		EntityID:       ingredientID,
		NoResultsText:  "No measurement units found matching your search.",
	}), nil
}

// SearchIngredientsForMeasurementUnit searches for ingredients to add to a measurement unit.
func (s *AdminFrontendServer) SearchIngredientsForMeasurementUnit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)
	query := req.URL.Query().Get("q")

	if query == "" {
		return ghtml.Div(), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Search for ingredients
	searchRes, err := c.GetValidIngredients(ctx, &mealplanningsvc.GetValidIngredientsRequest{
		Filter: nil, // TODO: Add filtering by query
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error searching: %v", err)),
		), nil
	}

	// Filter results by query (simple substring match)
	var results []*components.SearchResultItem
	for _, ingredient := range searchRes.Results {
		if contains(ingredient.Name, query) || contains(ingredient.Description, query) {
			results = append(results, &components.SearchResultItem{
				ID:          ingredient.Id,
				Name:        ingredient.Name,
				Description: ingredient.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_measurement_units/%s/ingredients", measurementUnitID),
		HTMXTarget:     "#association-list-container",
		EntityID:       measurementUnitID,
		NoResultsText:  "No ingredients found matching your search.",
	}), nil
}

// CreateIngredientMeasurementUnitFromIngredient creates an association from the ingredient side.
func (s *AdminFrontendServer) CreateIngredientMeasurementUnitFromIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	ingredientID := s.validIngredientIDRouteParamFetcher(req)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Parse the measurement unit MealPlanTaskID from the request
	var input struct {
		ID string `json:"id"`
	}
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error decoding request: %v", err)),
		), nil
	}

	// Create the association
	_, err = c.CreateValidIngredientMeasurementUnit(ctx, &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{
		Input: &mealplanningsvc.ValidIngredientMeasurementUnitCreationRequestInput{
			ValidIngredientId:      ingredientID,
			ValidMeasurementUnitId: input.ID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidIngredientMeasurementUnitsForIngredient(nil, req)
}

// CreateIngredientMeasurementUnitFromMeasurementUnit creates an association from the measurement unit side.
func (s *AdminFrontendServer) CreateIngredientMeasurementUnitFromMeasurementUnit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Parse the ingredient MealPlanTaskID from the request
	var input struct {
		ID string `json:"id"`
	}
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error decoding request: %v", err)),
		), nil
	}

	// Create the association
	_, err = c.CreateValidIngredientMeasurementUnit(ctx, &mealplanningsvc.CreateValidIngredientMeasurementUnitRequest{
		Input: &mealplanningsvc.ValidIngredientMeasurementUnitCreationRequestInput{
			ValidIngredientId:      input.ID,
			ValidMeasurementUnitId: measurementUnitID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidIngredientMeasurementUnitsForMeasurementUnit(nil, req)
}

// DeleteIngredientMeasurementUnit deletes an association.
func (s *AdminFrontendServer) DeleteIngredientMeasurementUnit(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	// Get the association MealPlanTaskID from the URL path
	associationID := req.PathValue("associationID")
	if associationID == "" {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No association MealPlanTaskID provided"),
		), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Archive (delete) the association
	_, err = c.ArchiveValidIngredientMeasurementUnit(ctx, &mealplanningsvc.ArchiveValidIngredientMeasurementUnitRequest{
		ValidIngredientMeasurementUnitId: associationID,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error deleting association: %v", err)),
		), nil
	}

	// Determine which list to return based on the referer or a parameter
	if ingredientID := s.validIngredientIDRouteParamFetcher(req); ingredientID != "" {
		return s.ValidIngredientMeasurementUnitsForIngredient(nil, req)
	} else if measurementUnitID := s.validMeasurementUnitIDRouteParamFetcher(req); measurementUnitID != "" {
		return s.ValidIngredientMeasurementUnitsForMeasurementUnit(nil, req)
	}

	return ghtml.Div(
		ghtml.Class("text-sm text-green-600 py-2"),
		g.Text("Association deleted successfully"),
	), nil
}
