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

// ValidPreparationInstrumentsForPreparation lists all instruments associated with a preparation
func (s *AdminFrontendServer) ValidPreparationInstrumentsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	preparationID := s.validPreparationIDRouteParamFetcher(req)
	if preparationID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Instruments",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: "No instruments associated with this preparation.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Instruments",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidPreparationInstrumentsByPreparation(ctx, &mealplanningsvc.GetValidPreparationInstrumentsByPreparationRequest{
		ValidPreparationID: preparationID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Instruments",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.Instrument != nil {
			items = append(items, components.AssociationItem{
				ID:          assoc.ID,
				Name:        assoc.Instrument.Name,
				Description: assoc.Instrument.Description,
				Notes:       assoc.Notes,
			})
		}
	}

	return components.AssociationList(&components.AssociationListProps{
		Title:                "Associated Instruments",
		Palette:              &design.StandardPalette,
		Items:                items,
		EntityID:             preparationID,
		AddSearchPlaceholder: "Search for instruments to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_preparations/%s/instruments/search", preparationID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_preparations/%s/instruments", preparationID),
		DeleteEndpoint:       "/api/valid_preparation_instruments",
		NoItemsMessage:       "No instruments associated with this preparation.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// ValidPreparationInstrumentsForInstrument lists all preparations associated with an instrument
func (s *AdminFrontendServer) ValidPreparationInstrumentsForInstrument(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	instrumentID := s.validInstrumentIDRouteParamFetcher(req)
	if instrumentID == "" {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: "No preparations associated with this instrument.",
		}), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: "Error: No API client available",
		}), nil
	}

	// Fetch associations
	res, err := c.GetValidPreparationInstrumentsByInstrument(ctx, &mealplanningsvc.GetValidPreparationInstrumentsByInstrumentRequest{
		ValidInstrumentID: instrumentID,
	})
	if err != nil {
		return components.AssociationList(&components.AssociationListProps{
			Title:          "Associated Preparations",
			Palette:        &design.StandardPalette,
			Items:          []components.AssociationItem{},
			NoItemsMessage: fmt.Sprintf("Error loading associations: %v", err),
		}), nil
	}

	// Convert to AssociationItems
	items := make([]components.AssociationItem, 0, len(res.Results))
	for _, assoc := range res.Results {
		if assoc.Preparation != nil {
			items = append(items, components.AssociationItem{
				ID:          assoc.ID,
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
		EntityID:             instrumentID,
		AddSearchPlaceholder: "Search for preparations to add...",
		AddSearchEndpoint:    fmt.Sprintf("/api/valid_instruments/%s/preparations/search", instrumentID),
		CreateEndpoint:       fmt.Sprintf("/api/valid_instruments/%s/preparations", instrumentID),
		DeleteEndpoint:       "/api/valid_preparation_instruments",
		NoItemsMessage:       "No preparations associated with this instrument.",
		HTMXTarget:           "#association-list-container",
	}), nil
}

// SearchInstrumentsForPreparation searches for instruments to add to a preparation
func (s *AdminFrontendServer) SearchInstrumentsForPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Search for instruments
	searchRes, err := c.GetValidInstruments(ctx, &mealplanningsvc.GetValidInstrumentsRequest{
		Filter: nil, // TODO: Add filtering by query
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error searching: %v", err)),
		), nil
	}

	// Filter results by query (simple substring match)
	var results []components.SearchResultItem
	for _, instrument := range searchRes.Results {
		if contains(instrument.Name, query) || contains(instrument.Description, query) {
			results = append(results, components.SearchResultItem{
				ID:          instrument.ID,
				Name:        instrument.Name,
				Description: instrument.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_preparations/%s/instruments", preparationID),
		HTMXTarget:     "#association-list-container",
		EntityID:       preparationID,
		NoResultsText:  "No instruments found matching your search.",
	}), nil
}

// SearchPreparationsForInstrument searches for preparations to add to an instrument
func (s *AdminFrontendServer) SearchPreparationsForInstrument(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	instrumentID := s.validInstrumentIDRouteParamFetcher(req)
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
	var results []components.SearchResultItem
	for _, preparation := range searchRes.Results {
		if contains(preparation.Name, query) || contains(preparation.Description, query) {
			results = append(results, components.SearchResultItem{
				ID:          preparation.ID,
				Name:        preparation.Name,
				Description: preparation.Description,
			})
		}
	}

	return components.AssociationSearchResults(&components.AssociationSearchResultsProps{
		Results:        results,
		CreateEndpoint: fmt.Sprintf("/api/valid_instruments/%s/preparations", instrumentID),
		HTMXTarget:     "#association-list-container",
		EntityID:       instrumentID,
		NoResultsText:  "No preparations found matching your search.",
	}), nil
}

// CreatePreparationInstrumentFromPreparation creates an association from the preparation side
func (s *AdminFrontendServer) CreatePreparationInstrumentFromPreparation(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Parse the instrument ID from the request
	var input struct {
		ID string `json:"id"`
	}
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error decoding request: %v", err)),
		), nil
	}

	// Create the association
	_, err = c.CreateValidPreparationInstrument(ctx, &mealplanningsvc.CreateValidPreparationInstrumentRequest{
		Input: &mealplanningsvc.ValidPreparationInstrumentCreationRequestInput{
			ValidPreparationID: preparationID,
			ValidInstrumentID:  input.ID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidPreparationInstrumentsForPreparation(nil, req)
}

// CreatePreparationInstrumentFromInstrument creates an association from the instrument side
func (s *AdminFrontendServer) CreatePreparationInstrumentFromInstrument(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	instrumentID := s.validInstrumentIDRouteParamFetcher(req)

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
	if err := s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error decoding request: %v", err)),
		), nil
	}

	// Create the association
	_, err = c.CreateValidPreparationInstrument(ctx, &mealplanningsvc.CreateValidPreparationInstrumentRequest{
		Input: &mealplanningsvc.ValidPreparationInstrumentCreationRequestInput{
			ValidPreparationID: input.ID,
			ValidInstrumentID:  instrumentID,
		},
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error creating association: %v", err)),
		), nil
	}

	// Return the updated list
	return s.ValidPreparationInstrumentsForInstrument(nil, req)
}

// DeletePreparationInstrument deletes an association
func (s *AdminFrontendServer) DeletePreparationInstrument(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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
	_, err = c.ArchiveValidPreparationInstrument(ctx, &mealplanningsvc.ArchiveValidPreparationInstrumentRequest{
		ValidPreparationInstrumentID: associationID,
	})
	if err != nil {
		return ghtml.Div(
			ghtml.Class("text-sm text-red-600 py-2"),
			g.Text(fmt.Sprintf("Error deleting association: %v", err)),
		), nil
	}

	// Determine which list to return based on the referer or a parameter
	// For simplicity, we'll check if we have a preparation or instrument ID in the URL
	if preparationID := s.validPreparationIDRouteParamFetcher(req); preparationID != "" {
		return s.ValidPreparationInstrumentsForPreparation(nil, req)
	} else if instrumentID := s.validInstrumentIDRouteParamFetcher(req); instrumentID != "" {
		return s.ValidPreparationInstrumentsForInstrument(nil, req)
	}

	return ghtml.Div(
		ghtml.Class("text-sm text-green-600 py-2"),
		g.Text("Association deleted successfully"),
	), nil
}

// Helper function for case-insensitive substring match
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
