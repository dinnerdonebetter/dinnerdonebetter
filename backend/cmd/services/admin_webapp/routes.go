package main

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/routing"

	"maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func (s *Server) setupRoutes(router routing.Router) error {
	authedRouter := router.WithMiddleware(s.authMiddleware)

	router.Get("/", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		return s.pageBuilder.HomePage(ctx), nil
	}))

	router.Get("/about", ghttp.Adapt(func(res http.ResponseWriter, req *http.Request) (gomponents.Node, error) {
		ctx, span := s.tracer.StartSpan(req.Context())
		defer span.End()

		return s.pageBuilder.AboutPage(ctx), nil
	}))

	router.Get("/login", ghttp.Adapt(s.renderLoginPage))
	router.Post("/login/submit", ghttp.Adapt(s.handleLoginSubmission))

	authedRouter.Get("/users", ghttp.Adapt(s.renderUsersPage))
	authedRouter.Get("/valid_ingredients", ghttp.Adapt(s.renderValidIngredientsPage))
	authedRouter.Get("/valid_ingredients/new", ghttp.Adapt(s.renderValidIngredientCreationForm))
	authedRouter.Post("/valid_ingredients/new/submit", ghttp.Adapt(s.handleValidIngredientSubmission))

	return nil
}
