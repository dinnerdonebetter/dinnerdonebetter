package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// This file demonstrates how the HTMX search functionality works

// DemoUser represents a simple user structure for demonstration
type DemoUser struct {
	ID       string
	Username string
	Email    string
	Status   string
}

// DemoUsersPage shows how to create a table page with HTMX search
func DemoUsersPage() g.Node {
	// Sample data
	users := []DemoUser{
		{"1", "john_doe", "john@example.com", "Active"},
		{"2", "jane_smith", "jane@example.com", "Active"},
		{"3", "bob_wilson", "bob@example.com", "Inactive"},
		{"4", "alice_brown", "alice@example.com", "Active"},
	}

	result, err := components.TablePage(&components.TablePageProps[DemoUser]{
		Title:             "Demo Users",
		BaseSubtitle:      "HTMX search demonstration",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search users by name or email...",

		// HTMX configuration - this makes search work!
		HTMXSearchTarget:  "/api/demo/users/search",
		HTMXSearchTrigger: "keyup changed delay:300ms", // Optional, this is the default

		Data: users,
		Actions: []g.Node{
			components.ActionButton("Add User", "/users/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[DemoUser]{
			TableID: "demo-users-table",
			Palette: &design.StandardPalette,
			Fields:  []string{"ID", "Username", "Email", "Status"},
		},
		EmptyStateTitle:       "No users found",
		EmptyStateDescription: "Try adjusting your search criteria.",
	})
	if err != nil {
		return ghtml.P(g.Text(fmt.Sprintf("Error: %v", err)))
	}

	return result.Node
}

// DemoUsersSearchHandler shows how to implement the search endpoint
func DemoUsersSearchHandler(w http.ResponseWriter, r *http.Request) {
	// This is what your search endpoint should look like

	// 1. Get the search query from URL parameters
	searchQuery := r.URL.Query().Get("search")

	// 2. Fetch your data (from database, API, etc.)
	allUsers := []DemoUser{
		{"1", "john_doe", "john@example.com", "Active"},
		{"2", "jane_smith", "jane@example.com", "Active"},
		{"3", "bob_wilson", "bob@example.com", "Inactive"},
		{"4", "alice_brown", "alice@example.com", "Active"},
	}

	// 3. Filter the data based on search query
	var filteredUsers []DemoUser
	if searchQuery == "" {
		filteredUsers = allUsers
	} else {
		searchLower := strings.ToLower(searchQuery)
		for _, user := range allUsers {
			if strings.Contains(strings.ToLower(user.Username), searchLower) ||
				strings.Contains(strings.ToLower(user.Email), searchLower) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	// 4. Generate the table HTML (NOT a full page)
	var tableHTML g.Node

	if len(filteredUsers) == 0 {
		tableHTML = components.EmptyState(
			"No users found",
			fmt.Sprintf("No users match '%s'", searchQuery),
			&design.StandardPalette,
			[]g.Node{
				components.ActionButton("Add User", "/users/new", &design.StandardPalette, true),
			},
		)
	} else {
		table, err := components.Table(filteredUsers, &components.TableOptions[DemoUser]{
			TableID: "demo-users-table",
			Palette: &design.StandardPalette,
			Fields:  []string{"ID", "Username", "Email", "Status"},
		})
		if err != nil {
			tableHTML = ghtml.P(g.Text(fmt.Sprintf("Error: %v", err)))
		} else {
			tableHTML = table
		}
	}

	// 5. Render just the table HTML and return it
	// (The gomponents/http adapter handles this automatically)
	// In a real handler you would use: ghttp.Adapt(handlerFunc)
	_ = tableHTML // This would be returned by your handler
}

/*
How It Works:

1. **Search Input**: When user types in search box, HTMX automatically sends GET request to `/api/users/search?search=<query>`

2. **HTMX Configuration**:
   - `HTMXSearchTarget: "/api/users/search"` - WHERE to send the request
   - `HTMXSearchTrigger: "keyup changed delay:300ms"` - WHEN to send (after 300ms of no typing)
   - HTMX automatically adds `hx-get`, `hx-target="#search-results"`, `hx-swap="innerHTML"`

3. **Search Handler**:
   - Gets search query from URL params
   - Filters data
   - Returns JUST the table HTML (not full page)

4. **HTMX Response**:
   - HTMX receives the table HTML
   - Replaces content of `#search-results` div with the new table
   - User sees filtered results instantly

5. **Key Points**:
   - Search endpoint returns HTML fragments, not JSON
   - The `#search-results` div wraps the table in TablePage component
   - No JavaScript needed - HTMX handles everything
   - Works with empty states, error states, etc.

Setup Steps:

1. Add HTMX properties to TablePageProps:
   ```go
   HTMXSearchTarget: "/api/users/search",
   ```

2. Create search route:
   ```go
   r.Get("/api/users/search", ghttp.Adapt(s.UsersSearch))
   ```

3. Implement search handler that returns table HTML:
   ```go
   func (s *Server) UsersSearch(w http.ResponseWriter, r *http.Request) (g.Node, error) {
       query := r.URL.Query().Get("search")
       // ... filter data ...
       return components.Table(filteredData, options)
   }
   ```

That's it! The search will work automatically with real-time filtering as users type.
*/
