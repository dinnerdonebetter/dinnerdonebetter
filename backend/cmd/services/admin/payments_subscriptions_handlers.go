package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/payments"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const (
	subscriptionIDURLParamKey = "subscriptionID"
)

func (s *AdminFrontendServer) SubscriptionCreate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("New Subscription", s.renderSubscriptionsError("Error: No API client available")), nil
	}

	var input *paymentssvc.SubscriptionCreationRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("New Subscription", s.renderSubscriptionsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	createRes, err := c.CreateSubscription(ctx, &paymentssvc.CreateSubscriptionRequest{
		Input: input,
	})
	if err != nil {
		return page("New Subscription", s.renderSubscriptionsError(fmt.Sprintf("Error creating subscription: %v", err))), nil
	}

	if createRes == nil || createRes.Created == nil {
		return page("New Subscription", s.renderSubscriptionsError("Error: No subscription returned from server")), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/subscriptions/%s", createRes.Created.Id), http.StatusSeeOther)
	return g.El("div"), nil
}

func (s *AdminFrontendServer) SubscriptionNewPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	if _, err := fetchClientFromContext(ctx); err != nil {
		return page("New Subscription", s.renderSubscriptionsError("Error: No API client available")), nil
	}

	emptyInput := &paymentssvc.SubscriptionCreationRequestInput{
		Status: "active",
	}

	formPageResult, err := components.FormPage(&components.FormPageProps[*paymentssvc.SubscriptionCreationRequestInput]{
		Title:        "Create New Subscription",
		BaseSubtitle: "Add a new subscription (manual grant)",
		Palette:      &design.StandardPalette,
		Data:         emptyInput,
		FormOptions: &components.FormOptions[*paymentssvc.SubscriptionCreationRequestInput]{
			FormID: "create-subscription-form",
			Action: "/api/subscriptions",
			Method: "POST",
			EnabledFields: []string{
				"belongs_to_account",
				"product_id",
				"external_subscription_id",
				"status",
				"current_period_start",
				"current_period_end",
			},
			FieldConfigs: map[string]*components.FieldConfig{
				"belongs_to_account": {
					Placeholder: "Account ID (from /accounts)...",
					Validation:  &components.FieldValidation{Required: true, CustomMessage: "Account ID is required"},
				},
				"product_id": {
					Placeholder: "Product ID (from /products)...",
					Validation:  &components.FieldValidation{Required: true, CustomMessage: "Product ID is required"},
				},
				"external_subscription_id": {
					Placeholder: "External subscription ID from payment provider (optional)...",
				},
				"status": {
					Options: []*components.SelectOption{
						{Value: "active", Label: "Active", IsDefault: true},
						{Value: "cancelled", Label: "Cancelled"},
						{Value: "past_due", Label: "Past Due"},
						{Value: "trialing", Label: "Trialing"},
						{Value: "incomplete", Label: "Incomplete"},
					},
					Placeholder: "Select status...",
					Validation:  &components.FieldValidation{Required: true},
				},
				"current_period_start": {
					InputType: "datetime-local",
				},
				"current_period_end": {
					InputType: "datetime-local",
				},
			},
			FormRows: []*components.FormRow{
				{Fields: []string{"belongs_to_account", "product_id"}, Columns: 2},
				{Fields: []string{"external_subscription_id"}, Columns: 1},
				{Fields: []string{"status"}, Columns: 1},
				{Fields: []string{"current_period_start", "current_period_end"}, Columns: 2},
			},
			SubmitButtonText: "Create Subscription",
			ShowCancelButton: true,
			CancelButtonText: "Cancel",
			CancelURL:        "/subscriptions",
			HTMXTarget:       "body",
			HTMXSwap:         "innerHTML",
			HTMXPushURL:      true,
			HTMXExtension:    "json-enc",
		},
		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Subscriptions", URL: "/subscriptions"},
			{Text: "New Subscription", URL: ""},
		},
	})
	if err != nil {
		return page("New Subscription", s.renderSubscriptionsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("New Subscription", formPageResult.Node), nil
}

func (s *AdminFrontendServer) SubscriptionPage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError("Error: No API client available")), nil
	}

	subscriptionID := s.subscriptionIDRouteParamFetcher(req)
	if subscriptionID == "" {
		return page("Subscriptions", s.renderSubscriptionsError("Error: No subscription ID provided")), nil
	}

	subRes, err := c.GetSubscription(ctx, &paymentssvc.GetSubscriptionRequest{SubscriptionId: subscriptionID})
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error loading subscription: %v", err))), nil
	}
	sub := subRes.Result

	formPageResult, err := components.FormPage(&components.FormPageProps[*paymentssvc.Subscription]{
		Title:        "Subscription Details",
		BaseSubtitle: "View and edit subscription information",
		Palette:      &design.StandardPalette,
		Data:         sub,
		FormOptions: &components.FormOptions[*paymentssvc.Subscription]{
			FormID: "edit-subscription-form",
			Action: fmt.Sprintf("/api/subscriptions/%s", sub.Id),
			Method: "POST",
			EnabledFields: []string{
				"status",
				"current_period_start",
				"current_period_end",
			},
			FieldConfigs: map[string]*components.FieldConfig{
				"status": {
					Options: []*components.SelectOption{
						{Value: "active", Label: "Active"},
						{Value: "cancelled", Label: "Cancelled"},
						{Value: "past_due", Label: "Past Due"},
						{Value: "trialing", Label: "Trialing"},
						{Value: "incomplete", Label: "Incomplete"},
					},
					Validation: &components.FieldValidation{Required: true},
				},
				"current_period_start": {InputType: "datetime-local"},
				"current_period_end":   {InputType: "datetime-local"},
			},
			FormRows: []*components.FormRow{
				{Fields: []string{"status"}, Columns: 1},
				{Fields: []string{"current_period_start", "current_period_end"}, Columns: 2},
			},
			SubmitButtonText: "Save Changes",
			ShowCancelButton: true,
			CancelButtonText: "Back to Subscriptions",
			CancelURL:        fmt.Sprintf("/subscriptions?account_id=%s", sub.BelongsToAccount),
			HTMXTarget:       "body",
			HTMXSwap:         "innerHTML",
			HTMXPushURL:      true,
			HTMXExtension:    "json-enc",
		},
		ShowBreadcrumbs: true,
		Breadcrumbs: []*components.Breadcrumb{
			{Text: "Dashboard", URL: "/"},
			{Text: "Subscriptions", URL: "/subscriptions"},
			{Text: sub.Id, URL: ""},
		},
		SubtitleGenerator: func(sb *paymentssvc.Subscription) string {
			return fmt.Sprintf("Subscription %s — %s", sb.Id, sb.Status)
		},
		AdditionalContent: []g.Node{
			components.ContentContainer(&components.ContentContainerProps{
				Title:    "Subscription Info",
				Subtitle: "Read-only fields",
				Palette:  &design.StandardPalette,
			}, components.Card(&design.StandardPalette,
				ghtml.Dl(
					ghtml.Class("grid grid-cols-2 gap-4"),
					ghtml.Dt(ghtml.Class("font-medium"), g.Text("Account")),
					ghtml.Dd(ghtml.Class("text-gray-600"), g.Text(sub.BelongsToAccount)),
					ghtml.Dt(ghtml.Class("font-medium"), g.Text("Product")),
					ghtml.Dd(ghtml.Class("text-gray-600"), g.Text(sub.ProductId)),
					ghtml.Dt(ghtml.Class("font-medium"), g.Text("External ID")),
					ghtml.Dd(ghtml.Class("text-gray-600"), g.Text(sub.ExternalSubscriptionId)),
				),
			)),
		},
	})
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error creating form: %v", err))), nil
	}

	return page("Subscriptions", formPageResult.Node), nil
}

func (s *AdminFrontendServer) SubscriptionUpdate(res http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError("Error: No API client available")), nil
	}

	subscriptionID := s.subscriptionIDRouteParamFetcher(req)
	if subscriptionID == "" {
		return page("Subscriptions", s.renderSubscriptionsError("Error: No subscription ID provided")), nil
	}

	var input *paymentssvc.SubscriptionUpdateRequestInput
	if err = s.encoder.DecodeRequest(ctx, req, &input); err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error decoding request: %v", err))), nil
	}

	_, err = c.UpdateSubscription(ctx, &paymentssvc.UpdateSubscriptionRequest{
		SubscriptionId: subscriptionID,
		Input:          input,
	})
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error updating subscription: %v", err))), nil
	}

	http.Redirect(res, req, fmt.Sprintf("/subscriptions/%s", subscriptionID), http.StatusSeeOther)
	return g.El("div"), nil
}

func (s *AdminFrontendServer) SubscriptionsList(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError("Error: No API client available")), nil
	}

	accountID := req.URL.Query().Get("account_id")
	if accountID == "" {
		return page("Subscriptions", components.ContentContainer(&components.ContentContainerProps{
			Title:    "Subscriptions",
			Subtitle: "View subscriptions for an account",
			Palette:  &design.StandardPalette,
		},
			components.Card(&design.StandardPalette,
				ghtml.P(
					ghtml.Class(fmt.Sprintf("py-8 %s", design.TextColor(design.StandardPalette.Text))),
					g.Text("Enter an account ID in the URL to view subscriptions, e.g. /subscriptions?account_id=xxx"),
				),
				components.ActionButton("Create New Subscription", "/subscriptions/new", &design.StandardPalette, true),
			),
		)), nil
	}

	queryFilter, grpcFilter := buildQueryFilterFromRequest(req)
	subsRes, err := c.GetSubscriptionsForAccount(ctx, &paymentssvc.GetSubscriptionsForAccountRequest{
		AccountId: accountID,
		Filter:    grpcFilter,
	})
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error loading subscriptions: %v", err))), nil
	}

	pagination := buildPaginationFromGRPCResponse(subsRes.Pagination)
	paginationURLGenerator := buildPaginationURLGeneratorForSearch(req, "/api/subscriptions/search?account_id="+accountID, queryFilter)

	tablePageResult, err := components.TablePage(&components.TablePageProps[*paymentssvc.Subscription]{
		Title:             "Subscriptions",
		BaseSubtitle:      fmt.Sprintf("Subscriptions for account %s", accountID),
		Palette:           &design.StandardPalette,
		ShowSearch:        true,
		SearchPlaceholder: "Search subscriptions...",
		HTMXSearchTarget:  "/api/subscriptions/search?account_id=" + accountID,
		Data:              subsRes.Results,
		Actions: []g.Node{
			components.ActionButton("Create New Subscription", "/subscriptions/new", &design.StandardPalette, true),
		},
		TableOptions: &components.TableOptions[*paymentssvc.Subscription]{
			TableID:                "subscriptions-table",
			Palette:                &design.StandardPalette,
			Pagination:             pagination,
			PaginationURLGenerator: paginationURLGenerator,
			Fields: []string{
				"id",
				"belongs_to_account",
				"product_id",
				"external_subscription_id",
				"status",
				"current_period_start",
				"current_period_end",
				"created_at",
				"last_updated_at",
				"archived_at",
			},
			FieldReplacements: map[string]string{
				"belongs_to_account":       "Account",
				"product_id":               "Product",
				"external_subscription_id": "External ID",
				"current_period_start":     "Period Start",
				"current_period_end":       "Period End",
			},
			FieldRenderers: map[string]components.FieldRenderer{
				"created_at":           renderTimestamp,
				"last_updated_at":      renderTimestamp,
				"archived_at":          renderTimestamp,
				"current_period_start": renderTimestamp,
				"current_period_end":   renderTimestamp,
			},
		},
		RowLinkGenerator: func(data *paymentssvc.Subscription) string {
			return fmt.Sprintf("/subscriptions/%s", data.Id)
		},
		EmptyStateTitle:       "No subscriptions found",
		EmptyStateDescription: "Create a subscription for this account.",
		EmptyStateActions: []g.Node{
			components.ActionButton("Create New Subscription", "/subscriptions/new", &design.StandardPalette, true),
		},
		SubtitleGenerator: func(metadata components.TablePageMetadata) string {
			if metadata.EmptyState {
				return fmt.Sprintf("Subscriptions for account %s", accountID)
			}
			return fmt.Sprintf("%d subscriptions for account %s", metadata.TotalCount, accountID)
		},
	})
	if err != nil {
		return page("Subscriptions", s.renderSubscriptionsError(fmt.Sprintf("Error creating table: %v", err))), nil
	}

	return page("Subscriptions", tablePageResult.Node), nil
}

func (s *AdminFrontendServer) SubscriptionsSearch(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
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

	accountID := req.URL.Query().Get("account_id")
	if accountID == "" {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState("No account selected", "Add account_id to the URL to search subscriptions.", &design.StandardPalette, nil),
		), nil
	}

	searchQuery := req.URL.Query().Get("search")
	subsRes, err := c.GetSubscriptionsForAccount(ctx, &paymentssvc.GetSubscriptionsForAccountRequest{
		AccountId: accountID,
		Filter:    &grpcfiltering.QueryFilter{},
	})
	if err != nil {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(fmt.Sprintf("Error loading subscriptions: %v", err)),
			),
		), nil
	}

	var filtered []*paymentssvc.Subscription
	if searchQuery == "" {
		filtered = subsRes.Results
	} else {
		q := strings.ToLower(searchQuery)
		for _, sub := range subsRes.Results {
			if strings.Contains(strings.ToLower(sub.Id), q) ||
				strings.Contains(strings.ToLower(sub.Status), q) ||
				strings.Contains(strings.ToLower(sub.ProductId), q) ||
				strings.Contains(strings.ToLower(sub.ExternalSubscriptionId), q) {
				filtered = append(filtered, sub)
			}
		}
	}

	if len(filtered) == 0 {
		return g.El("div",
			g.Attr("class", "overflow-x-auto"),
			components.EmptyState(
				"No subscriptions found",
				fmt.Sprintf("No subscriptions match the search term '%s'.", searchQuery),
				&design.StandardPalette,
				[]g.Node{
					components.ActionButton("Create New Subscription", "/subscriptions/new", &design.StandardPalette, true),
				},
			),
		), nil
	}

	table, err := components.Table(filtered, &components.TableOptions[*paymentssvc.Subscription]{
		TableID: "subscriptions-table",
		Palette: &design.StandardPalette,
		Fields: []string{
			"id",
			"belongs_to_account",
			"product_id",
			"external_subscription_id",
			"status",
			"current_period_start",
			"current_period_end",
			"created_at",
			"last_updated_at",
			"archived_at",
		},
		FieldReplacements: map[string]string{
			"belongs_to_account":       "Account",
			"product_id":               "Product",
			"external_subscription_id": "External ID",
			"current_period_start":     "Period Start",
			"current_period_end":       "Period End",
		},
		FieldRenderers: map[string]components.FieldRenderer{
			"created_at":           renderTimestamp,
			"last_updated_at":      renderTimestamp,
			"archived_at":          renderTimestamp,
			"current_period_start": renderTimestamp,
			"current_period_end":   renderTimestamp,
		},
		RowLinkGenerator: func(data *paymentssvc.Subscription) string {
			return fmt.Sprintf("/subscriptions/%s", data.Id)
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

	return g.El("div",
		g.Attr("class", "overflow-x-auto"),
		table,
	), nil
}

func (s *AdminFrontendServer) renderSubscriptionsError(errorMsg string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Subscriptions",
		Subtitle: "Manage subscriptions",
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
