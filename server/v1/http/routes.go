package httpserver

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	"gitlab.com/prixfixe/prixfixe/services/v1/ingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/instruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	"gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	"gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	"gitlab.com/prixfixe/prixfixe/services/v1/preparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	"gitlab.com/prixfixe/prixfixe/services/v1/reports"
	"gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/users"
	"gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/heptiolabs/healthcheck"
)

const (
	numericIDPattern = "/{%s:[0-9]+}"
	oauth2IDPattern  = "/{%s:[0-9_\\-]+}"
)

func (s *Server) setupRouter(frontendConfig config.FrontendSettings, metricsHandler metrics.Handler) {
	router := chi.NewRouter()

	// Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	ch := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts,
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// Maximum value not ignored by any of major browsers,
		MaxAge: 300,
	})

	router.Use(
		middleware.RequestID,
		middleware.Timeout(maxTimeout),
		s.loggingMiddleware,
		ch.Handler,
	)

	// all middleware must be defined before routes on a mux

	router.Route("/_meta_", func(metaRouter chi.Router) {
		health := healthcheck.NewHandler()
		// Expose a liveness check on /live
		metaRouter.Get("/live", health.LiveEndpoint)
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.Handle("/metrics", metricsHandler)
	}

	// Frontend routes
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}

	for route, handler := range s.frontendService.Routes() {
		router.Get(route, handler)
	}

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler())
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler())
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler())

		userIDPattern := fmt.Sprintf(oauth2IDPattern, users.URIParamKey)

		userRouter.Get("/", s.usersService.ListHandler())
		userRouter.With(s.usersService.UserInputMiddleware).Post("/", s.usersService.CreateHandler())
		userRouter.Get(userIDPattern, s.usersService.ReadHandler())
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler())

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler())

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler())
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler())

		oauth2Router.With(s.oauth2ClientsService.OAuth2ClientInfoMiddleware).
			Post("/authorize", func(res http.ResponseWriter, req *http.Request) {
				s.logger.WithRequest(req).Debug("oauth2 authorize route hit")
				if err := s.oauth2ClientsService.HandleAuthorizeRequest(res, req); err != nil {
					http.Error(res, err.Error(), http.StatusBadRequest)
				}
			})

		oauth2Router.Post("/token", func(res http.ResponseWriter, req *http.Request) {
			if err := s.oauth2ClientsService.HandleTokenRequest(res, req); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
		})
	})

	router.With(s.authService.AuthenticationMiddleware(true)).Route("/api/v1", func(v1Router chi.Router) {
		// Instruments
		v1Router.Route("/instruments", func(instrumentsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, instruments.URIParamKey)
			instrumentsRouter.With(s.instrumentsService.CreationInputMiddleware).Post("/", s.instrumentsService.CreateHandler())
			instrumentsRouter.Get(sr, s.instrumentsService.ReadHandler())
			instrumentsRouter.With(s.instrumentsService.UpdateInputMiddleware).Put(sr, s.instrumentsService.UpdateHandler())
			instrumentsRouter.Delete(sr, s.instrumentsService.ArchiveHandler())
			instrumentsRouter.Get("/", s.instrumentsService.ListHandler())
		})
		// Ingredients
		v1Router.Route("/ingredients", func(ingredientsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, ingredients.URIParamKey)
			ingredientsRouter.With(s.ingredientsService.CreationInputMiddleware).Post("/", s.ingredientsService.CreateHandler())
			ingredientsRouter.Get(sr, s.ingredientsService.ReadHandler())
			ingredientsRouter.With(s.ingredientsService.UpdateInputMiddleware).Put(sr, s.ingredientsService.UpdateHandler())
			ingredientsRouter.Delete(sr, s.ingredientsService.ArchiveHandler())
			ingredientsRouter.Get("/", s.ingredientsService.ListHandler())
		})
		// Preparations
		v1Router.Route("/preparations", func(preparationsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, preparations.URIParamKey)
			preparationsRouter.With(s.preparationsService.CreationInputMiddleware).Post("/", s.preparationsService.CreateHandler())
			preparationsRouter.Get(sr, s.preparationsService.ReadHandler())
			preparationsRouter.With(s.preparationsService.UpdateInputMiddleware).Put(sr, s.preparationsService.UpdateHandler())
			preparationsRouter.Delete(sr, s.preparationsService.ArchiveHandler())
			preparationsRouter.Get("/", s.preparationsService.ListHandler())
		})
		// RequiredPreparationInstruments
		v1Router.Route("/required_preparation_instruments", func(requiredPreparationInstrumentsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, requiredpreparationinstruments.URIParamKey)
			requiredPreparationInstrumentsRouter.With(s.requiredPreparationInstrumentsService.CreationInputMiddleware).Post("/", s.requiredPreparationInstrumentsService.CreateHandler())
			requiredPreparationInstrumentsRouter.Get(sr, s.requiredPreparationInstrumentsService.ReadHandler())
			requiredPreparationInstrumentsRouter.With(s.requiredPreparationInstrumentsService.UpdateInputMiddleware).Put(sr, s.requiredPreparationInstrumentsService.UpdateHandler())
			requiredPreparationInstrumentsRouter.Delete(sr, s.requiredPreparationInstrumentsService.ArchiveHandler())
			requiredPreparationInstrumentsRouter.Get("/", s.requiredPreparationInstrumentsService.ListHandler())
		})
		// Recipes
		v1Router.Route("/recipes", func(recipesRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipes.URIParamKey)
			recipesRouter.With(s.recipesService.CreationInputMiddleware).Post("/", s.recipesService.CreateHandler())
			recipesRouter.Get(sr, s.recipesService.ReadHandler())
			recipesRouter.With(s.recipesService.UpdateInputMiddleware).Put(sr, s.recipesService.UpdateHandler())
			recipesRouter.Delete(sr, s.recipesService.ArchiveHandler())
			recipesRouter.Get("/", s.recipesService.ListHandler())
		})
		// RecipeSteps
		v1Router.Route("/recipe_steps", func(recipeStepsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipesteps.URIParamKey)
			recipeStepsRouter.With(s.recipeStepsService.CreationInputMiddleware).Post("/", s.recipeStepsService.CreateHandler())
			recipeStepsRouter.Get(sr, s.recipeStepsService.ReadHandler())
			recipeStepsRouter.With(s.recipeStepsService.UpdateInputMiddleware).Put(sr, s.recipeStepsService.UpdateHandler())
			recipeStepsRouter.Delete(sr, s.recipeStepsService.ArchiveHandler())
			recipeStepsRouter.Get("/", s.recipeStepsService.ListHandler())
		})
		// RecipeStepInstruments
		v1Router.Route("/recipe_step_instruments", func(recipeStepInstrumentsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipestepinstruments.URIParamKey)
			recipeStepInstrumentsRouter.With(s.recipeStepInstrumentsService.CreationInputMiddleware).Post("/", s.recipeStepInstrumentsService.CreateHandler())
			recipeStepInstrumentsRouter.Get(sr, s.recipeStepInstrumentsService.ReadHandler())
			recipeStepInstrumentsRouter.With(s.recipeStepInstrumentsService.UpdateInputMiddleware).Put(sr, s.recipeStepInstrumentsService.UpdateHandler())
			recipeStepInstrumentsRouter.Delete(sr, s.recipeStepInstrumentsService.ArchiveHandler())
			recipeStepInstrumentsRouter.Get("/", s.recipeStepInstrumentsService.ListHandler())
		})
		// RecipeStepIngredients
		v1Router.Route("/recipe_step_ingredients", func(recipeStepIngredientsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipestepingredients.URIParamKey)
			recipeStepIngredientsRouter.With(s.recipeStepIngredientsService.CreationInputMiddleware).Post("/", s.recipeStepIngredientsService.CreateHandler())
			recipeStepIngredientsRouter.Get(sr, s.recipeStepIngredientsService.ReadHandler())
			recipeStepIngredientsRouter.With(s.recipeStepIngredientsService.UpdateInputMiddleware).Put(sr, s.recipeStepIngredientsService.UpdateHandler())
			recipeStepIngredientsRouter.Delete(sr, s.recipeStepIngredientsService.ArchiveHandler())
			recipeStepIngredientsRouter.Get("/", s.recipeStepIngredientsService.ListHandler())
		})
		// RecipeStepProducts
		v1Router.Route("/recipe_step_products", func(recipeStepProductsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipestepproducts.URIParamKey)
			recipeStepProductsRouter.With(s.recipeStepProductsService.CreationInputMiddleware).Post("/", s.recipeStepProductsService.CreateHandler())
			recipeStepProductsRouter.Get(sr, s.recipeStepProductsService.ReadHandler())
			recipeStepProductsRouter.With(s.recipeStepProductsService.UpdateInputMiddleware).Put(sr, s.recipeStepProductsService.UpdateHandler())
			recipeStepProductsRouter.Delete(sr, s.recipeStepProductsService.ArchiveHandler())
			recipeStepProductsRouter.Get("/", s.recipeStepProductsService.ListHandler())
		})
		// RecipeIterations
		v1Router.Route("/recipe_iterations", func(recipeIterationsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipeiterations.URIParamKey)
			recipeIterationsRouter.With(s.recipeIterationsService.CreationInputMiddleware).Post("/", s.recipeIterationsService.CreateHandler())
			recipeIterationsRouter.Get(sr, s.recipeIterationsService.ReadHandler())
			recipeIterationsRouter.With(s.recipeIterationsService.UpdateInputMiddleware).Put(sr, s.recipeIterationsService.UpdateHandler())
			recipeIterationsRouter.Delete(sr, s.recipeIterationsService.ArchiveHandler())
			recipeIterationsRouter.Get("/", s.recipeIterationsService.ListHandler())
		})
		// RecipeStepEvents
		v1Router.Route("/recipe_step_events", func(recipeStepEventsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, recipestepevents.URIParamKey)
			recipeStepEventsRouter.With(s.recipeStepEventsService.CreationInputMiddleware).Post("/", s.recipeStepEventsService.CreateHandler())
			recipeStepEventsRouter.Get(sr, s.recipeStepEventsService.ReadHandler())
			recipeStepEventsRouter.With(s.recipeStepEventsService.UpdateInputMiddleware).Put(sr, s.recipeStepEventsService.UpdateHandler())
			recipeStepEventsRouter.Delete(sr, s.recipeStepEventsService.ArchiveHandler())
			recipeStepEventsRouter.Get("/", s.recipeStepEventsService.ListHandler())
		})
		// IterationMedias
		v1Router.Route("/iteration_medias", func(iterationMediasRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, iterationmedias.URIParamKey)
			iterationMediasRouter.With(s.iterationMediasService.CreationInputMiddleware).Post("/", s.iterationMediasService.CreateHandler())
			iterationMediasRouter.Get(sr, s.iterationMediasService.ReadHandler())
			iterationMediasRouter.With(s.iterationMediasService.UpdateInputMiddleware).Put(sr, s.iterationMediasService.UpdateHandler())
			iterationMediasRouter.Delete(sr, s.iterationMediasService.ArchiveHandler())
			iterationMediasRouter.Get("/", s.iterationMediasService.ListHandler())
		})
		// Invitations
		v1Router.Route("/invitations", func(invitationsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, invitations.URIParamKey)
			invitationsRouter.With(s.invitationsService.CreationInputMiddleware).Post("/", s.invitationsService.CreateHandler())
			invitationsRouter.Get(sr, s.invitationsService.ReadHandler())
			invitationsRouter.With(s.invitationsService.UpdateInputMiddleware).Put(sr, s.invitationsService.UpdateHandler())
			invitationsRouter.Delete(sr, s.invitationsService.ArchiveHandler())
			invitationsRouter.Get("/", s.invitationsService.ListHandler())
		})
		// Reports
		v1Router.Route("/reports", func(reportsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, reports.URIParamKey)
			reportsRouter.With(s.reportsService.CreationInputMiddleware).Post("/", s.reportsService.CreateHandler())
			reportsRouter.Get(sr, s.reportsService.ReadHandler())
			reportsRouter.With(s.reportsService.UpdateInputMiddleware).Put(sr, s.reportsService.UpdateHandler())
			reportsRouter.Delete(sr, s.reportsService.ArchiveHandler())
			reportsRouter.Get("/", s.reportsService.ListHandler())
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
			webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post("/", s.webhooksService.CreateHandler())
			webhookRouter.Get(sr, s.webhooksService.ReadHandler())
			webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler())
			webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler())
			webhookRouter.Get("/", s.webhooksService.ListHandler())
		})

		// OAuth2 Clients
		v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
			// CreateHandler is not bound to an OAuth2 authentication token
			// UpdateHandler not supported for OAuth2 clients.
			clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler())
			clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler())
			clientRouter.Get("/", s.oauth2ClientsService.ListHandler())
		})
	})

	s.router = router
}
