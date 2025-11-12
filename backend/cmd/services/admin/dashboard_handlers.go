package main

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

func (s *AdminFrontendServer) HomePage(_ http.ResponseWriter, req *http.Request) (g.Node, error) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)

	c, err := fetchClientFromContext(ctx)
	if err != nil {
		return page("Dashboard", s.renderDashboardError("Error: No API client available")), nil
	}

	// Fetch user count
	usersRes, err := c.GetUsers(ctx, &identitysvc.GetUsersRequest{})
	userCountStr := "-"
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching users for dashboard")
	} else if usersRes != nil {
		userCountStr = fmt.Sprintf("%d", len(usersRes.Result))
	}

	// Fetch account count
	accountsRes, err := c.GetAccounts(ctx, &identitysvc.GetAccountsRequest{})
	accountCountStr := "-"
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching accounts for dashboard")
	} else if accountsRes != nil {
		accountCountStr = fmt.Sprintf("%d", len(accountsRes.Result))
	}

	// Fetch recipe count
	recipesRes, err := c.GetRecipes(ctx, &mealplanningsvc.GetRecipesRequest{})
	recipeCountStr := "-"
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching recipes for dashboard")
	} else if recipesRes != nil {
		recipeCountStr = fmt.Sprintf("%d", len(recipesRes.Results))
	}

	return page("Dashboard",
		components.ContentContainer(&components.ContentContainerProps{
			Title:    "Dashboard",
			Subtitle: "Welcome to the admin dashboard",
			Palette:  &design.StandardPalette,
		},
			components.Card(&design.StandardPalette,
				ghtml.H2(
					ghtml.Class(fmt.Sprintf("text-lg font-medium %s mb-4", design.TextColor(design.StandardPalette.Primary))),
					g.Text("Quick Stats"),
				),
				ghtml.Div(
					ghtml.Class("grid grid-cols-1 md:grid-cols-3 gap-4"),
					statCard("Users", userCountStr, &design.StandardPalette),
					statCard("Accounts", accountCountStr, &design.StandardPalette),
					statCard("Recipes", recipeCountStr, &design.StandardPalette),
				),
			),
		),
	), nil
}

func (s *AdminFrontendServer) renderDashboardError(message string) g.Node {
	return components.ContentContainer(&components.ContentContainerProps{
		Title:    "Dashboard",
		Subtitle: "Welcome to the admin dashboard",
		Palette:  &design.StandardPalette,
	},
		components.Card(&design.StandardPalette,
			ghtml.P(
				ghtml.Class(fmt.Sprintf("text-center py-8 %s", design.TextColor(design.StandardPalette.Warning))),
				g.Text(message),
			),
		),
	)
}

func statCard(title, value string, palette *design.Palette) g.Node {
	return ghtml.Div(
		ghtml.Class(fmt.Sprintf("p-4 %s rounded-lg border %s",
			design.Background(design.Color{Value: "gray-50"}),
			design.BorderColor(palette.Background),
		)),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("text-sm font-medium %s", design.TextColor(palette.Text))),
			g.Text(title),
		),
		ghtml.Div(
			ghtml.Class(fmt.Sprintf("mt-1 text-2xl font-bold %s", design.TextColor(palette.Primary))),
			g.Text(value),
		),
	)
}
