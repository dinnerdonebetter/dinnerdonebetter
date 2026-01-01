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

// ValidPreparationVesselsForPreparation lists all vessels associated with a preparation.
func (s *AdminFrontendServer) ValidPreparationVesselsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	preparationID := s.validPreparationIDRouteParamFetcher(req)
	if preparationID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Vessels",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No vessels associated with this preparation.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Vessels",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidPreparationVesselsByPreparation(ctx, &mealplanningsvc.GetValidPreparationVesselsByPreparationRequest{
		ValidPreparationId: preparationID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Vessels",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]*components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.Vessel != nil {
			items = append(items, &components.AssociationItem{
				ID:          assoc.Id,
				Name:        assoc.Vessel.Name,
				Description: assoc.Vessel.Description,
				Notes:       assoc.Notes,
			})
		}
	}

	return components.AssociationList(&components.AssociationListProps{
		Title:                "Associated Vessels",
		Palette:              &design.StandardPalette,
		Items:                items,
		EntityID:             preparationID,
		AddSearchPlaceholder: "Search for vessels to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_preparations/%s/vessels/search", preparationID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_preparations/%s/vessels", preparationID),
		DeleteEndpoint:       "/api/valid_preparation_vessels",
		NoItemsMessage:       "No vessels associated with this preparation.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// ValidPreparationVesselsForVessel lists all preparations associated with a vessel.
func (s *AdminFrontendServer) ValidPreparationVesselsForVessel(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	vesselID := s.validVesselIDRouteParamFetcher(req)
	if vesselID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []*components.AssociationItem{},
			NoItemsMessage: "No preparations associated with this vessel.",
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
	res, err := c.GetValidPreparationVesselsByVessel(ctx, &mealplanningsvc.GetValidPreparationVesselsByVesselRequest{
		ValidVesselId: vesselID,
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
		EntityID:             vesselID,
		AddSearchPlaceholder: "Search for preparations to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_vessels/%s/preparations/search", vesselID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_vessels/%s/preparations", vesselID),
		DeleteEndpoint:       "/api/valid_preparation_vessels",
		NoItemsMessage:       "No preparations associated with this vessel.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// SearchVesselsForPreparation searches for vessels to add to a preparation.
func (s *AdminFrontendServer) SearchVesselsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Search for vessels
	searchRes, err := c.GetValidVessels(ctx, &mealplanningsvc.GetValidVesselsRequest{
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
	for _, vessel := range searchRes.Results {
		if contains(vessel.Name, query) || contains(vessel.Description, query) {
			results = append(results, &components.SearchResultItem{
				ID:          vessel.Id,
				Name:        vessel.Name,
				Description: vessel.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_preparations/%s/vessels", preparationID),
		HTMXTarget:     "#association-list-container",
		EntityID:       preparationID,
		NoResultsText:  "No vessels found matching your search.",
	}), nil
}

// SearchPreparationsForVessel searches for preparations to add to a vessel.
func (s *AdminFrontendServer) SearchPreparationsForVessel(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	vesselID := s.validVesselIDRouteParamFetcher(req)
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
		CreateEndpoint: fmt.Sprintf("/api/valid_vessels/%s/preparations", vesselID),
		HTMXTarget:     "#association-list-container",
		EntityID:       vesselID,
		NoResultsText:  "No preparations found matching your search.",
	}), nil
}

// CreatePreparationVesselFromPreparation creates an association from the preparation side.
func (s *AdminFrontendServer) CreatePreparationVesselFromPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Parse the vessel ID from the request
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
	_, err = c.CreateValidPreparationVessel(ctx, &mealplanningsvc.CreateValidPreparationVesselRequest{
		Input: &mealplanningsvc.ValidPreparationVesselCreationRequestInput{
			ValidPreparationId: preparationID,
			ValidVesselId:      input.ID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidPreparationVesselsForPreparation(nil, req)
}

// CreatePreparationVesselFromVessel creates an association from the vessel side.
func (s *AdminFrontendServer) CreatePreparationVesselFromVessel(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	vesselID := s.validVesselIDRouteParamFetcher(req)

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
	_, err = c.CreateValidPreparationVessel(ctx, &mealplanningsvc.CreateValidPreparationVesselRequest{
		Input: &mealplanningsvc.ValidPreparationVesselCreationRequestInput{
			ValidPreparationId: input.ID,
			ValidVesselId:      vesselID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidPreparationVesselsForVessel(nil, req)
}

// DeletePreparationVessel deletes an association.
func (s *AdminFrontendServer) DeletePreparationVessel(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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
	_, err = c.ArchiveValidPreparationVessel(ctx, &mealplanningsvc.ArchiveValidPreparationVesselRequest{
		ValidPreparationVesselId: associationID,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error deleting association: %v", err)),
		), nil
	}

	// Determine which list to return based on the referer or a parameter
	if preparationID := s.validPreparationIDRouteParamFetcher(req); preparationID != "" {
		return s.ValidPreparationVesselsForPreparation(nil, req)
	} else if vesselID := s.validVesselIDRouteParamFetcher(req); vesselID != "" {
		return s.ValidPreparationVesselsForVessel(nil, req)
	}

	return ghtml.Div(
		ghtml.Class("text-sm text-green-600 py-2"),
		g.Text("Association deleted successfully"),
	), nil
}
