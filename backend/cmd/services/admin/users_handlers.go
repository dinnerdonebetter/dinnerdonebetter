package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	userIDURLParamKey = "userID"
)

func (s *AdminFrontendServer) UserPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Users", s.renderUsersError("Error: No API client available")), nil
	}

	userID := s.userIDRouteParamFetcher(req)
	if userID == "" {
		return page("Users", s.renderUsersError("Error: No user ID provided")), nil
	}

	usersRes, err := c.GetUser(ctx, &identitysvc.GetUserRequest{UserID: userID})
	if err != nil {
		return page("Users", s.renderUsersError(fmt.Sprintf("Error loading users: %v", err))), nil
	}

	// Use the new integrated TablePage component

	return page("Users", g.Text(usersRes.Result.Username)), nil
}

func (s *AdminFrontendServer) UsersList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Users", s.renderUsersError("Error: No API client available")), nil
	}

	usersRes, err := c.GetUsers(ctx, &identitysvc.GetUsersRequest{})
	if err != nil {
		return page("Users", s.renderUsersError(fmt.Sprintf("Error loading users: %v", err))), nil
	}

	// Use the new integrated TablePage component
	tablePageResult, err := components.TablePage(&components.TablePageProps[*identitysvc.User]{
		Title:             "Users",
		BaseSubtitle:      "Manage user accounts",
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search users...",
		HTMXSearchTarget:  "/api/users/search",
		Data:              usersRes.Result,
		Actions:           []g.Node{},
		TableOptions: &components.TableOptions[*identitysvc.User]{
			TableID: "users-table",
			Palette: &design.StandardPalette,
			Fields: []string{
				"ID",
				"Username",
				"FirstName",
				"LastName",
				"EmailAddress",
				"ServiceRole",
				"AccountStatus",
				"AccountStatusExplanation",
				"Birthday",
				"PasswordLastChangedAt",
				"LastAcceptedTermsOfService",
				"LastAcceptedPrivacyPolicy",
				"TwoFactorSecretVerifiedAt",
				"EmailAddressVerifiedAt",
				"CreatedAt",
				"LastUpdatedAt",
				"ArchivedAt",
			},
			FieldReplacements: map[string]string{
				"EmailAddressVerifiedAt": "Email Verified At",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"Birthday":                  renderTimestamp,
				"CreatedAt":                 renderTimestamp,
				"TwoFactorSecretVerifiedAt": renderTimestamp,
				"LastUpdatedAt":             renderTimestamp,
				"ArchivedAt":                renderTimestamp,
			},
		},
		EmptyStateTitle:       "No users found",
		EmptyStateDescription: "Get started by creating your first user account.",
		EmptyStateActions:     []g.Node{},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return "Manage user accounts"
			}
			return fmt.Sprintf("Manage %d user accounts", metadata.TotalCount)
		},
	})
	if err != nil {
		return page("Users", s.renderUsersError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Users", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) UsersSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Get search query from request
	searchQuery := req.URL.Query().Get("search")

	usersRes, err := c.GetUsers(ctx, &identitysvc.GetUsersRequest{})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading users: %v", err)),
			),
		), nil
	}

	// Filter users based on search query
	var filteredUsers []*identitysvc.User
	if searchQuery == "" {
		// No search query, return all users
		filteredUsers = usersRes.Result
	} else {
		// Filter users by search query (case insensitive)
		searchQueryLower := strings.ToLower(searchQuery)
		for _, user := range usersRes.Result {
			if strings.Contains(strings.ToLower(user.Username), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.FirstName), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.LastName), searchQueryLower) ||
				strings.Contains(strings.ToLower(user.EmailAddress), searchQueryLower) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	// Generate just the table (not the full page)
	if len(filteredUsers) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No users found",
				fmt.Sprintf("No users match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{
					components.ActionButton("Add User", "/users/new", &design.StandardPalette, true),
				},
			),
		), nil
	}

	table, err := components.Table(filteredUsers, &components.TableOptions[*identitysvc.User]{
		TableID: "users-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"ID",
			"Username",
			"FirstName",
			"LastName",
			"EmailAddress",
			"ServiceRole",
			"AccountStatus",
			"AccountStatusExplanation",
			"Birthday",
			"PasswordLastChangedAt",
			"LastAcceptedTermsOfService",
			"LastAcceptedPrivacyPolicy",
			"TwoFactorSecretVerifiedAt",
			"EmailAddressVerifiedAt",
			"CreatedAt",
			"LastUpdatedAt",
			"ArchivedAt",
		},
		FieldReplacements: map[string]string{
			"EmailAddressVerifiedAt": "Email Verified At",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"Birthday":                  renderTimestamp,
			"CreatedAt":                 renderTimestamp,
			"TwoFactorSecretVerifiedAt": renderTimestamp,
			"LastUpdatedAt":             renderTimestamp,
			"ArchivedAt":                renderTimestamp,
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

// renderUsersError creates a consistent error display for the users page
func (s *AdminFrontendServer) renderUsersError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Users",
		Subtitle: "Manage user accounts",
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
