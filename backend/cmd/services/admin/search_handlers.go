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

// SearchValidMeasurementUnits handles searching for valid measurement units.
func (s *AdminFrontendServer) SearchValidMeasurementUnits(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get("q")
	searchBoxID := req.URL.Query().Get("search_box_id")

	if query == "" {
		return ghtml.Div(ghtml.Style("display: none;")), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("px-4 py-3 text-sm text-red-600"),
			g.Text("Error: No API client available"),
		), nil
	}

	searchRes, err := c.SearchForValidMeasurementUnits(ctx, &mealplanningsvc.SearchForValidMeasurementUnitsRequest{
		Query: query,
	})
	if err != nil {
		s.logger.Error("error searching valid measurement units", err)
		return ghtml.Div(
			ghtml.Class("px-4 py-3 text-sm text-red-600"),
			g.Text("Error performing search"),
		), nil
	}

	var results []*components.SearchResultItem
	for _, unit := range searchRes.Results {
		extraInfo := ""
		if unit.Metric {
			extraInfo = "Metric"
		} else if unit.Imperial {
			extraInfo = "Imperial"
		}
		if unit.Volumetric {
			if extraInfo != "" {
				extraInfo += ", "
			}
			extraInfo += "Volumetric"
		}

		results = append(results, &components.SearchResultItem{
			ID:          unit.Id,
			Name:        unit.Name,
			Description: unit.Description,
			ExtraInfo:   extraInfo,
		})
	}

	// Generate onclick handler that updates the hidden field and UI
	onSelectJS := fmt.Sprintf(`
		document.getElementById('%s-hidden').value = '%%s';
		document.getElementById('%s-input').style.display = 'none';
		document.getElementById('%s-results').style.display = 'none';
		document.getElementById('%s-query').value = '';
		document.getElementById('%s-results').innerHTML = '';
		const selectedDiv = document.getElementById('%s-selected');
		if (selectedDiv) {
			selectedDiv.style.display = 'flex';
			selectedDiv.querySelector('span').textContent = '%%s';
		} else {
			const display = document.createElement('div');
			display.id = '%s-selected';
			display.className = 'mb-2 px-3 py-2 bg-blue-50 border border-blue-200 rounded-md flex items-center justify-between';
			display.innerHTML = '<span class="text-sm text-gray-700">%%s</span><button type="button" class="text-sm text-blue-600 hover:text-blue-800 font-medium" onclick="document.getElementById(\'%s-hidden\').value = \'\'; this.parentElement.style.display = \'none\'; document.getElementById(\'%s-input\').style.display = \'block\';">Clear</button>';
			document.getElementById('%s-input').parentElement.insertBefore(display, document.getElementById('%s-input'));
		}
	`, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID)

	return components.SearchInputWithResults(searchBoxID, "/admin/search/valid_measurement_units", query, results, onSelectJS, &design.StandardPalette), nil
}

// SearchValidIngredients handles searching for valid ingredients.
func (s *AdminFrontendServer) SearchValidIngredients(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get("q")
	searchBoxID := req.URL.Query().Get("search_box_id")

	if query == "" {
		return ghtml.Div(ghtml.Style("display: none;")), nil
	}

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return ghtml.Div(
			ghtml.Class("px-4 py-3 text-sm text-red-600"),
			g.Text("Error: No API client available"),
		), nil
	}

	searchRes, err := c.SearchForValidIngredients(ctx, &mealplanningsvc.SearchForValidIngredientsRequest{
		Query: query,
	})
	if err != nil {
		s.logger.Error("error searching valid ingredients", err)
		return ghtml.Div(
			ghtml.Class("px-4 py-3 text-sm text-red-600"),
			g.Text("Error performing search"),
		), nil
	}

	var results []*components.SearchResultItem
	for _, ingredient := range searchRes.Results {
		results = append(results, &components.SearchResultItem{
			ID:          ingredient.Id,
			Name:        ingredient.Name,
			Description: ingredient.Description,
		})
	}

	// Generate onclick handler that updates the hidden field and UI
	onSelectJS := fmt.Sprintf(`
		document.getElementById('%s-hidden').value = '%%s';
		document.getElementById('%s-input').style.display = 'none';
		document.getElementById('%s-results').style.display = 'none';
		document.getElementById('%s-query').value = '';
		document.getElementById('%s-results').innerHTML = '';
		const selectedDiv = document.getElementById('%s-selected');
		if (selectedDiv) {
			selectedDiv.style.display = 'flex';
			selectedDiv.querySelector('span').textContent = '%%s';
		} else {
			const display = document.createElement('div');
			display.id = '%s-selected';
			display.className = 'mb-2 px-3 py-2 bg-blue-50 border border-blue-200 rounded-md flex items-center justify-between';
			display.innerHTML = '<span class="text-sm text-gray-700">%%s</span><button type="button" class="text-sm text-blue-600 hover:text-blue-800 font-medium" onclick="document.getElementById(\'%s-hidden\').value = \'\'; this.parentElement.style.display = \'none\'; document.getElementById(\'%s-input\').style.display = \'block\';">Clear</button>';
			document.getElementById('%s-input').parentElement.insertBefore(display, document.getElementById('%s-input'));
		}
	`, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID, searchBoxID)

	return components.SearchInputWithResults(searchBoxID, "/admin/search/valid_ingredients", query, results, onSelectJS, &design.StandardPalette), nil
}
