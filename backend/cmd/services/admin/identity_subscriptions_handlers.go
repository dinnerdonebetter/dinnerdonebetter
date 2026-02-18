package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/payments"

	"google.golang.org/protobuf/types/known/timestamppb"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// userSubscriptionRow holds subscription display fields for the user page table.
// We use a flat struct to avoid copying protobuf types (which contain mutexes).
type userSubscriptionRow struct {
	CurrentPeriodStart *timestamppb.Timestamp `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *timestamppb.Timestamp `json:"current_period_end,omitempty"`
	AccountName        string                 `json:"account_name"`
	ID                 string                 `json:"id"`
	ProductID          string                 `json:"product_id"`
	Status             string                 `json:"status"`
}

func (s *AdminFrontendServer) AccountSubscriptionsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	_, grpcFilter := buildQueryFilterFromRequest(req)
	if grpcFilter == nil {
		grpcFilter = &filtering.QueryFilter{}
	}
	var pageSize uint32 = 20
	if grpcFilter.MaxResponseSize == nil {
		grpcFilter.MaxResponseSize = &pageSize
	}

	subsRes, err := c.GetSubscriptionsForAccount(ctx, &paymentssvc.GetSubscriptionsForAccountRequest{
		AccountId: accountID,
		Filter:    grpcFilter,
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading subscriptions: %v", err)),
			),
		), nil
	}

	if len(subsRes.Results) == 0 {
		return g.El("div",
			g.Attr("class", "text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No subscriptions for this account."),
			),
			ghtml.Div(
				ghtml.Class("mt-4"),
				components.ActionButton("Create Subscription", "/subscriptions/new", &design.StandardPalette, true),
			),
		), nil
	}

	table, err := components.Table(subsRes.Results, &components.TableOptions[*paymentssvc.Subscription]{
		TableID: "account-subscriptions-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"product_id",
			"status",
			"current_period_start",
			"current_period_end",
			"created_at",
		},
		FieldReplacements: map[string]string{
			"product_id":           "Product",
			"current_period_start": "Period Start",
			"current_period_end":   "Period End",
			"created_at":           "Created",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at":           renderTimestamp,
			"current_period_start": renderTimestamp,
			"current_period_end":   renderTimestamp,
		},
		RowLinkGenerator: func(sub *paymentssvc.Subscription) string {
			return fmt.Sprintf("/subscriptions/%s", sub.Id)
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
	if subsRes.Pagination != nil && subsRes.Pagination.Cursor != "" {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("flex justify-between items-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.Div(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d subscriptions", len(subsRes.Results))),
				),
				ghtml.Button(
					ghtml.Class(fmt.Sprintf("px-4 py-2 text-sm font-medium rounded-md %s %s hover:%s",
						design.TextColor(design.Color{Value: "white"}),
						design.Background(design.StandardPalette.Primary),
						design.Background(design.Color{Value: design.StandardPalette.Primary.Value + "-700"}),
					)),
					g.Attr("hx-get", fmt.Sprintf("/api/accounts/%s/subscriptions?cursor=%s", accountID, subsRes.Pagination.Cursor)),
					g.Attr("hx-target", "#account-subscriptions-container"),
					g.Attr("hx-swap", "innerHTML"),
					g.Text("Load More"),
				),
			),
		)
	} else if len(subsRes.Results) > 0 {
		paginationControls = append(paginationControls,
			ghtml.Div(
				ghtml.Class("text-center mt-4 pt-4 border-t border-gray-200"),
				ghtml.P(
					ghtml.Class("text-sm text-gray-500"),
					g.Text(fmt.Sprintf("Showing %d subscription(s)", len(subsRes.Results))),
				),
			),
		)
	}

	return g.El("div",
		g.Attr("class", "space-y-4"),
		ghtml.Div(
			ghtml.Class("flex justify-end"),
			components.ActionButton("Create Subscription", "/subscriptions/new", &design.StandardPalette, true),
		),
		table,
		g.Group(paginationControls),
	), nil
}

func (s *AdminFrontendServer) UserSubscriptionsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	// Fetch user's accounts (reasonable limit for aggregation)
	var accountsPageSize uint32 = 50
	accountsRes, err := c.GetAccountsForUser(ctx, &identitysvc.GetAccountsForUserRequest{
		UserId: userID,
		Filter: &filtering.QueryFilter{
			MaxResponseSize: &accountsPageSize,
		},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "text-center py-4"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("%s text-sm", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading accounts: %v", err)),
			),
		), nil
	}

	if len(accountsRes.Results) == 0 {
		return g.El("div",
			g.Attr("class", "text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("This user has no accounts, so there are no subscriptions to show."),
			),
		), nil
	}

	// Build account ID to name map
	accountNames := make(map[string]string)
	for _, acc := range accountsRes.Results {
		accountNames[acc.Id] = acc.Name
	}

	// Fetch subscriptions for each account
	var rows []*userSubscriptionRow
	var subsPageSize uint32 = 50
	for _, acc := range accountsRes.Results {
		subsRes, err := c.GetSubscriptionsForAccount(ctx, &paymentssvc.GetSubscriptionsForAccountRequest{
			AccountId: acc.Id,
			Filter: &filtering.QueryFilter{
				MaxResponseSize: &subsPageSize,
			},
		})
		if err != nil {
			continue // Skip accounts that fail
		}
		for _, sub := range subsRes.Results {
			rows = append(rows, &userSubscriptionRow{
				AccountName:        accountNames[acc.Id],
				ID:                 sub.Id,
				ProductID:          sub.ProductId,
				Status:             sub.Status,
				CurrentPeriodStart: sub.CurrentPeriodStart,
				CurrentPeriodEnd:   sub.CurrentPeriodEnd,
			})
		}
	}

	if len(rows) == 0 {
		return g.El("div",
			g.Attr("class", "text-center py-8"),
			ghtml.P(
				ghtml.Class("text-sm text-gray-500"),
				g.Text("No subscriptions across this user's accounts."),
			),
		), nil
	}

	table, err := components.Table(rows, &components.TableOptions[*userSubscriptionRow]{
		TableID: "user-subscriptions-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"account_name",
			"id",
			"product_id",
			"status",
			"current_period_start",
			"current_period_end",
		},
		FieldReplacements: map[string]string{
			"account_name":         "Account",
			"product_id":           "Product",
			"current_period_start": "Period Start",
			"current_period_end":   "Period End",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"current_period_start": renderTimestamp,
			"current_period_end":   renderTimestamp,
		},
		RowLinkGenerator: func(row *userSubscriptionRow) string {
			return fmt.Sprintf("/subscriptions/%s", row.ID)
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

	return g.El("div",
		g.Attr("class", "space-y-4"),
		ghtml.Div(
			ghtml.Class("text-sm text-gray-500"),
			g.Text(fmt.Sprintf("Showing %d subscription(s) across %d account(s)", len(rows), len(accountsRes.Results))),
		),
		table,
	), nil
}
