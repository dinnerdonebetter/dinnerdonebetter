package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *AdminFrontendServer) UserAuditLogList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	userID := s.userIDRouteParamFetcher(req)
	if userID == "" {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No user ID provided"),
			),
		), nil
	}

	pageParam := req.URL.Query().Get("page")
	var pageSize uint32 = 20
	var nextCursor *string
	if pageParam != "" {
		nextCursor = &pageParam
	}

	_, grpcFilter := buildQueryFilterFromRequest(req)
	if grpcFilter == nil {
		grpcFilter = &filtering.QueryFilter{}
	}
	if grpcFilter.MaxResponseSize == nil {
		grpcFilter.MaxResponseSize = &pageSize
	}
	if nextCursor != nil {
		grpcFilter.Cursor = nextCursor
	}

	res, err := c.GetAuditLogEntriesForUser(ctx, &auditsvc.GetAuditLogEntriesForUserRequest{
		UserId: userID,
		Filter: grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading audit log: %v", err)),
			),
		), nil
	}

	if len(res.Results) == 0 {
		return g.El("div",
			g.Attr("class", "text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No audit log entries for this user."),
			),
		), nil
	}

	table, err := components.Table(res.Results, &components.TableOptions[*auditsvc.AuditLogEntry]{
		TableID: "user-audit-log-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"resource_type",
			"relevant_id",
			"event_type",
			"created_at",
			"changes",
		},
		FieldReplacements: map[string]string{
			"resource_type": "Resource",
			"relevant_id":   "Relevant ID",
			"event_type":    "Event",
			"created_at":    "Created",
			"changes":       "Changes",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at": renderTimestamp,
			"changes": func(value any) g.Node {
				if value == nil {
					return g.Text("-")
				}
				if m, ok := value.(map[string]*auditsvc.ChangeLog); ok {
					n := len(m)
					if n == 0 {
						return g.Text("-")
					}
					return g.Text(fmt.Sprintf("%d field(s)", n))
				}
				return g.Text("-")
			},
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error creating table: %v", err)),
			),
		), nil
	}

	var paginationControls []g.Node
	if res.Pagination != nil && res.Pagination.Cursor != "" {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("flex justify-between items-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.Div(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d entries", len(res.Results))),
				),
				ghtml.Button(
					ghtml.Class(fmt.Sprintf("px-4 py-2 text-sm font-medium rounded-md %s %s hover:%s",
						design.TextColor(design.Color{Value: "white"}),
						design.Background(design.StandardPalette.Primary),
						design.Background(design.Color{Value: design.StandardPalette.Primary.Value + "-700"}),
					)),
					g.Attr("hx-get", fmt.Sprintf("/api/users/%s/audit-log?page=%s", userID, res.Pagination.Cursor)),
					g.Attr("hx-target", "#user-audit-log-container"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text("Load More"),
				),
			),
		)
	} else if len(res.Results) > 0 {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("text-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.P(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d audit log entry(ies)", len(res.Results))),
				),
			),
		)
	}

	return g.El("div",
		g.Attr("class", "space-y-4"),
		table,
		g.Group(paginationControls),
	), nil
}

func (s *AdminFrontendServer) AccountAuditLogList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No API client available"),
			),
		), nil
	}

	accountID := s.accountIDRouteParamFetcher(req)
	if accountID == "" {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text("Error: No account ID provided"),
			),
		), nil
	}

	pageParam := req.URL.Query().Get("page")
	var pageSize uint32 = 20
	var nextCursor *string
	if pageParam != "" {
		nextCursor = &pageParam
	}

	_, grpcFilter := buildQueryFilterFromRequest(req)
	if grpcFilter == nil {
		grpcFilter = &filtering.QueryFilter{}
	}
	if grpcFilter.MaxResponseSize == nil {
		grpcFilter.MaxResponseSize = &pageSize
	}
	if nextCursor != nil {
		grpcFilter.Cursor = nextCursor
	}

	res, err := c.GetAuditLogEntriesForAccount(ctx, &auditsvc.GetAuditLogEntriesForAccountRequest{
		AccountId: accountID,
		Filter:    grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading audit log: %v", err)),
			),
		), nil
	}

	if len(res.Results) == 0 {
		return g.El("div",
			g.Attr("class", "text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No audit log entries for this account."),
			),
		), nil
	}

	table, err := components.Table(res.Results, &components.TableOptions[*auditsvc.AuditLogEntry]{
		TableID: "account-audit-log-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"resource_type",
			"relevant_id",
			"event_type",
			"created_at",
			"changes",
		},
		FieldReplacements: map[string]string{
			"resource_type": "Resource",
			"relevant_id":   "Relevant ID",
			"event_type":    "Event",
			"created_at":    "Created",
			"changes":       "Changes",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at": renderTimestamp,
			"changes": func(value any) g.Node {
				if value == nil {
					return g.Text("-")
				}
				if m, ok := value.(map[string]*auditsvc.ChangeLog); ok {
					n := len(m)
					if n == 0 {
						return g.Text("-")
					}
					return g.Text(fmt.Sprintf("%d field(s)", n))
				}
				return g.Text("-")
			},
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error creating table: %v", err)),
			),
		), nil
	}

	var paginationControls []g.Node
	if res.Pagination != nil && res.Pagination.Cursor != "" {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("flex justify-between items-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.Div(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d entries", len(res.Results))),
				),
				ghtml.Button(
					ghtml.Class(fmt.Sprintf("px-4 py-2 text-sm font-medium rounded-md %s %s hover:%s",
						design.TextColor(design.Color{Value: "white"}),
						design.Background(design.StandardPalette.Primary),
						design.Background(design.Color{Value: design.StandardPalette.Primary.Value + "-700"}),
					)),
					g.Attr("hx-get", fmt.Sprintf("/api/accounts/%s/audit-log?page=%s", accountID, res.Pagination.Cursor)),
					g.Attr("hx-target", "#account-audit-log-container"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text("Load More"),
				),
			),
		)
	} else if len(res.Results) > 0 {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("text-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.P(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d audit log entry(ies)", len(res.Results))),
				),
			),
		)
	}

	return g.El("div",
		g.Attr("class", "space-y-4"),
		table,
		g.Group(paginationControls),
	), nil
}
