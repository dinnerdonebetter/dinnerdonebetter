package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// This file demonstrates how to use the new layout system with custom palettes

// CustomPalette demonstrates creating a custom color palette
var CustomPalette = design.Palette{
	Primary:    design.Color{Name: "Custom Primary", Value: "indigo-600"},
	Secondary:  design.Color{Name: "Custom Secondary", Value: "purple-100"},
	Accent:     design.Color{Name: "Custom Accent", Value: "amber-400"},
	Warning:    design.Color{Name: "Custom Warning", Value: "red-500"},
	Background: design.Color{Name: "Custom Background", Value: "slate-50"},
	Text:       design.Color{Name: "Custom Text", Value: "slate-800"},
}

// ExampleListPage demonstrates how to create a typical admin list page
func ExampleListPage(title string, data []string) g.Node {
	config := &LayoutConfig{
		Palette:     &CustomPalette,
		AppName:     "My Custom Admin",
		MaxWidth:    "6xl",
		Margin:      "6",
		ShowSidebar: true,
	}

	return pageWithConfig(title, config,
		components.ContentContainer(&components.ContentContainerProps{
			Title:             title,
			Subtitle:          fmt.Sprintf("Manage %d items", len(data)),
			Palette:           &CustomPalette,
			ShowSearch:        true,
			SearchPlaceholder: "Search items...",
			Actions: []g.Node{
				components.ActionButton("Add New", "/items/new", &CustomPalette, true),
				components.ActionButton("Import", "/items/import", &CustomPalette, false),
				components.ActionButton("Export", "/items/export", &CustomPalette, false),
			},
		},
			// Main content card
			components.CardWithHeader("Items List", &CustomPalette,
				[]g.Node{
					ghtml.Button(
						ghtml.Class(fmt.Sprintf("text-sm px-3 py-1 rounded %s %s hover:%s",
							design.TextColor(CustomPalette.Primary),
							design.Background(design.Color{Value: "white"}),
							design.Background(CustomPalette.Secondary),
						)),
						g.Text("Refresh"),
					),
				},
				// Table or list content would go here
				ghtml.Div(
					ghtml.Class("space-y-2"),
					g.Group(createExampleItems(data, &CustomPalette)),
				),
			),

			// Statistics card
			components.Card(&CustomPalette,
				ghtml.H3(
					ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(CustomPalette.Primary))),
					g.Text("Statistics"),
				),
				ghtml.Div(
					ghtml.Class("grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"),
					statCard("Total Items", fmt.Sprintf("%d", len(data)), &CustomPalette),
					statCard("Active", "89%", &CustomPalette),
					statCard("Pending", "11%", &CustomPalette),
					statCard("Errors", "0", &CustomPalette),
				),
			),
		),
	)
}

// ExampleDetailPage demonstrates how to create a typical admin detail page
func ExampleDetailPage(title, itemID string) g.Node {
	config := &LayoutConfig{
		Palette:     &CustomPalette,
		AppName:     "My Custom Admin",
		MaxWidth:    "4xl", // Narrower for detail pages
		Margin:      "4",
		ShowSidebar: true,
	}

	return pageWithConfig(title, config,
		components.ContentContainer(&components.ContentContainerProps{
			Title:    fmt.Sprintf("%s Details", title),
			Subtitle: fmt.Sprintf("ID: %s", itemID),
			Palette:  &CustomPalette,
			Actions: []g.Node{
				components.ActionButton("Edit", fmt.Sprintf("/items/%s/edit", itemID), &CustomPalette, true),
				components.ActionButton("Delete", fmt.Sprintf("/items/%s/delete", itemID), &CustomPalette, false),
				components.ActionButton("Back to List", "/items", &CustomPalette, false),
			},
		},
			// Main details card
			components.CardWithHeader("Item Information", &CustomPalette, nil,
				ghtml.Div(
					ghtml.Class("space-y-4"),
					detailRow("Name", "Example Item", &CustomPalette),
					detailRow("Status", "Active", &CustomPalette),
					detailRow("Created", "2025-01-01 12:00:00", &CustomPalette),
					detailRow("Last Updated", "2025-01-02 15:30:00", &CustomPalette),
				),
			),

			// Related actions card
			components.Card(&CustomPalette,
				ghtml.H3(
					ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(CustomPalette.Primary))),
					g.Text("Related Actions"),
				),
				ghtml.Div(
					ghtml.Class("flex flex-wrap gap-2"),
					components.ActionButton("View History", fmt.Sprintf("/items/%s/history", itemID), &CustomPalette, false),
					components.ActionButton("Download", fmt.Sprintf("/items/%s/download", itemID), &CustomPalette, false),
					components.ActionButton("Share", fmt.Sprintf("/items/%s/share", itemID), &CustomPalette, false),
				),
			),
		),
	)
}

// ExampleEmptyPage demonstrates how to handle empty states
func ExampleEmptyPage(title string) g.Node {
	return pageWithConfig(title, nil, // Using default config
		components.ContentContainer(&components.ContentContainerProps{
			Title:      title,
			Subtitle:   "Get started by adding your first item",
			Palette:    &design.StandardPalette,
			ShowSearch: false, // No search for empty state
		},
			components.EmptyState(
				"No items found",
				"There are no items to display. Create your first item to get started.",
				&design.StandardPalette,
				[]g.Node{
					components.ActionButton("Add First Item", "/items/new", &design.StandardPalette, true),
					components.ActionButton("Import Items", "/items/import", &design.StandardPalette, false),
				},
			),
		),
	)
}

// Helper functions for the examples

func createExampleItems(data []string, palette *design.Palette) []g.Node {
	if len(data) == 0 {
		return []g.Node{
			components.EmptyState(
				"No items",
				"No items to display",
				palette,
				nil,
			),
		}
	}

	var items []g.Node
	for i, item := range data {
		items = append(items, ghtml.Div(
			ghtml.Class(fmt.Sprintf("flex items-center justify-between p-3 %s border %s rounded-md",
				design.Background(design.Color{Value: "white"}),
				design.BorderColor(palette.Background),
			)),
			ghtml.Div(
				ghtml.Class("flex items-center space-x-3"),
				ghtml.Div(
					ghtml.Class(fmt.Sprintf("w-8 h-8 rounded-full %s %s flex items-center justify-center text-sm font-medium",
						design.Background(palette.Primary),
						design.TextColor(design.Color{Value: "white"}),
					)),
					g.Text(fmt.Sprintf("%d", i+1)),
				),
				ghtml.Div(
					ghtml.H4(
						ghtml.Class(fmt.Sprintf("text-sm font-medium %s", design.TextColor(palette.Text))),
						g.Text(item),
					),
					ghtml.P(
						ghtml.Class(fmt.Sprintf("text-xs %s", design.TextColor(palette.Text))),
						g.Text("Sample description"),
					),
				),
			),
			ghtml.Div(
				ghtml.Class("flex items-center space-x-2"),
				ghtml.Button(
					ghtml.Class(fmt.Sprintf("text-xs px-2 py-1 rounded %s hover:%s",
						design.TextColor(palette.Primary),
						design.Background(palette.Secondary),
					)),
					g.Text("Edit"),
				),
				ghtml.Button(
					ghtml.Class(fmt.Sprintf("text-xs px-2 py-1 rounded %s hover:%s",
						design.TextColor(palette.Warning),
						design.Background(design.Color{Value: "red-50"}),
					)),
					g.Text("Delete"),
				),
			),
		))
	}

	return items
}

func detailRow(label, value string, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Class("flex justify-between py-2 border-b border-gray-100"),
		ghtml.Dt(
			ghtml.Class(fmt.Sprintf("text-sm font-medium %s", design.TextColor(palette.Text))),
			g.Text(label),
		),
		ghtml.Dd(
			ghtml.Class(fmt.Sprintf("text-sm %s", design.TextColor(palette.Primary))),
			g.Text(value),
		),
	)
}

// ExampleWithHTMXSearch demonstrates HTMX integration for live search
func ExampleWithHTMXSearch(title string) g.Node {
	return page(title,
		components.ContentContainer(&components.ContentContainerProps{
			Title:    title,
			Subtitle: "Live search example with HTMX",
			Palette:  &design.StandardPalette,
		},
			// Custom search with HTMX
			ghtml.Div(
				ghtml.Class("mb-6"),
				components.SearchInput(&components.SearchInputProps{
					Placeholder: "Search with live results...",
					Palette:     &design.StandardPalette,
					HTMXTarget:  "/search",
					HTMXTrigger: "keyup changed delay:300ms",
				}),
			),

			// Search results container
			ghtml.Div(
				ghtml.ID("search-results"),
				ghtml.Class("space-y-4"),
				components.Card(&design.StandardPalette,
					ghtml.P(
						ghtml.Class("text-center py-8 text-gray-500"),
						g.Text("Start typing to search..."),
					),
				),
			),
		),
	)
}
