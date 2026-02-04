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

// ValidIngredientPreparationsForIngredient lists all preparations associated with an ingredient.
func (s *AdminFrontendServer) ValidIngredientPreparationsForIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	ingredientID := s.validIngredientIDRouteParamFetcher(req)
	if ingredientID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No preparations associated with this ingredient.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidIngredientPreparationsByIngredient(ctx, &mealplanningsvc.GetValidIngredientPreparationsByIngredientRequest{
		ValidIngredientId: ingredientID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]*components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.Preparation != nil {
			items = append(items, &components.AssociationItem{
				ID:          assoc.Id,
				Name:        assoc.Preparation.Name,
				Description: assoc.Preparation.Description,
				Notes:       assoc.Notes,
			})
		}
	}

	return components.AssociationList(&components.AssociationListProps{
		Title:                "Associated Preparations",
		Palette:              &design.StandardPalette,
		Items:                items,
		EntityID:             ingredientID,
		AddSearchPlaceholder: "Search for preparations to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_ingredients/%s/preparations/search", ingredientID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_ingredients/%s/preparations", ingredientID),
		DeleteEndpoint:       "/api/valid_ingredient_preparations",
		NoItemsMessage:       "No preparations associated with this ingredient.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// ValidIngredientPreparationsForPreparation lists all ingredients associated with a preparation.
func (s *AdminFrontendServer) ValidIngredientPreparationsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	preparationID := s.validPreparationIDRouteParamFetcher(req)
	if preparationID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Ingredients",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No ingredients associated with this preparation.",
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
	res, err := c.GetValidIngredientPreparationsByPreparation(ctx, &mealplanningsvc.GetValidIngredientPreparationsByPreparationRequest{
		ValidPreparationId: preparationID,
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
		EntityID:             preparationID,
		AddSearchPlaceholder: "Search for ingredients to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_preparations/%s/ingredients/search", preparationID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_preparations/%s/ingredients", preparationID),
		DeleteEndpoint:       "/api/valid_ingredient_preparations",
		NoItemsMessage:       "No ingredients associated with this preparation.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// SearchPreparationsForIngredient searches for preparations to add to an ingredient.
func (s *AdminFrontendServer) SearchPreparationsForIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Search for preparations
	searchRes, err := c.GetValidPreparations(ctx, &mealplanningsvc.GetValidPreparationsRequest{
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
	for _, preparation := range searchRes.Results {
		if contains(preparation.Name, query) || contains(preparation.Description, query) {
			results = append(results, &components.SearchResultItem{
				ID:          preparation.Id,
				Name:        preparation.Name,
				Description: preparation.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_ingredients/%s/preparations", ingredientID),
		HTMXTarget:     "#association-list-container",
		EntityID:       ingredientID,
		NoResultsText:  "No preparations found matching your search.",
	}), nil
}

// SearchIngredientsForPreparation searches for ingredients to add to a preparation.
func (s *AdminFrontendServer) SearchIngredientsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	preparationID := s.validPreparationIDRouteParamFetcher(req)
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
		CreateEndpoint: fmt.Sprintf("/api/valid_preparations/%s/ingredients", preparationID),
		HTMXTarget:     "#association-list-container",
		EntityID:       preparationID,
		NoResultsText:  "No ingredients found matching your search.",
	}), nil
}

// CreateIngredientPreparationFromIngredient creates an association from the ingredient side.
func (s *AdminFrontendServer) CreateIngredientPreparationFromIngredient(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Parse the preparation ID from the request
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
	_, err = c.CreateValidIngredientPreparation(ctx, &mealplanningsvc.CreateValidIngredientPreparationRequest{
		Input: &mealplanningsvc.ValidIngredientPreparationCreationRequestInput{
			ValidIngredientId:  ingredientID,
			ValidPreparationId: input.ID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidIngredientPreparationsForIngredient(nil, req)
}

// CreateIngredientPreparationFromPreparation creates an association from the preparation side.
func (s *AdminFrontendServer) CreateIngredientPreparationFromPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	preparationID := s.validPreparationIDRouteParamFetcher(req)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No API client available"),
		), nil
	}

	// Parse the ingredient ID from the request
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
	_, err = c.CreateValidIngredientPreparation(ctx, &mealplanningsvc.CreateValidIngredientPreparationRequest{
		Input: &mealplanningsvc.ValidIngredientPreparationCreationRequestInput{
			ValidIngredientId:  input.ID,
			ValidPreparationId: preparationID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidIngredientPreparationsForPreparation(nil, req)
}

// DeleteIngredientPreparation deletes an association.
func (s *AdminFrontendServer) DeleteIngredientPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	// Get the association ID from the URL path
	associationID := req.PathValue("associationID")
	if associationID == "" {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text("Error: No association ID provided"),
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
	_, err = c.ArchiveValidIngredientPreparation(ctx, &mealplanningsvc.ArchiveValidIngredientPreparationRequest{
		ValidIngredientPreparationId: associationID,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error deleting association: %v", err)),
		), nil
	}

	// Determine which list to return based on the referer or a parameter
	if ingredientID := s.validIngredientIDRouteParamFetcher(req); ingredientID != "" {
		return s.ValidIngredientPreparationsForIngredient(nil, req)
	} else if preparationID := s.validPreparationIDRouteParamFetcher(req); preparationID != "" {
		return s.ValidIngredientPreparationsForPreparation(nil, req)
	}

	return ghtml.Div(
		ghtml.Class("text-sm text-green-600 py-2"),
		g.Text("Association deleted successfully"),
	), nil
}
