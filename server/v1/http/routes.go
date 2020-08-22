package httpserver

import (
	"fmt"
	"net/http"
	"path/filepath"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepeventsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/heptiolabs/healthcheck"
)

const (
	root             = "/"
	searchRoot       = "/search"
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

	// all middleware must be defined before routes on a mux.

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

	// Frontend routes.
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/server_admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler)
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler)

		userIDPattern := fmt.Sprintf(oauth2IDPattern, usersservice.URIParamKey)

		userRouter.Get(root, s.usersService.ListHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Get("/status", s.authService.StatusHandler)
		userRouter.With(s.usersService.UserInputMiddleware).Post(root, s.usersService.CreateHandler)
		userRouter.Get(userIDPattern, s.usersService.ReadHandler)
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler)

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)

		userRouter.With(
			s.usersService.TOTPSecretVerificationInputMiddleware,
		).Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler)
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler)

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

	router.With(s.authService.AuthenticationMiddleware(true)).
		Route("/api/v1", func(v1Router chi.Router) {
			// ValidInstruments
			validInstrumentPath := "valid_instruments"
			validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
			validInstrumentRouteParam := fmt.Sprintf(numericIDPattern, validinstrumentsservice.URIParamKey)
			v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter chi.Router) {
				validInstrumentsRouter.With(s.validInstrumentsService.CreationInputMiddleware).Post(root, s.validInstrumentsService.CreateHandler)
				validInstrumentsRouter.Route(validInstrumentRouteParam, func(singleValidInstrumentRouter chi.Router) {
					singleValidInstrumentRouter.Get(root, s.validInstrumentsService.ReadHandler)
					singleValidInstrumentRouter.With(s.validInstrumentsService.UpdateInputMiddleware).Put(root, s.validInstrumentsService.UpdateHandler)
					singleValidInstrumentRouter.Delete(root, s.validInstrumentsService.ArchiveHandler)
					singleValidInstrumentRouter.Head(root, s.validInstrumentsService.ExistenceHandler)
				})
				validInstrumentsRouter.Get(root, s.validInstrumentsService.ListHandler)
				validInstrumentsRouter.Get(searchRoot, s.validInstrumentsService.SearchHandler)
			})

			// ValidIngredients
			validIngredientPath := "valid_ingredients"
			validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
			validIngredientRouteParam := fmt.Sprintf(numericIDPattern, validingredientsservice.URIParamKey)
			v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter chi.Router) {
				validIngredientsRouter.With(s.validIngredientsService.CreationInputMiddleware).Post(root, s.validIngredientsService.CreateHandler)
				validIngredientsRouter.Route(validIngredientRouteParam, func(singleValidIngredientRouter chi.Router) {
					singleValidIngredientRouter.Get(root, s.validIngredientsService.ReadHandler)
					singleValidIngredientRouter.With(s.validIngredientsService.UpdateInputMiddleware).Put(root, s.validIngredientsService.UpdateHandler)
					singleValidIngredientRouter.Delete(root, s.validIngredientsService.ArchiveHandler)
					singleValidIngredientRouter.Head(root, s.validIngredientsService.ExistenceHandler)
				})
				validIngredientsRouter.Get(root, s.validIngredientsService.ListHandler)
				validIngredientsRouter.Get(searchRoot, s.validIngredientsService.SearchHandler)
			})

			// ValidPreparations
			validPreparationPath := "valid_preparations"
			validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
			validPreparationRouteParam := fmt.Sprintf(numericIDPattern, validpreparationsservice.URIParamKey)
			v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter chi.Router) {
				validPreparationsRouter.With(s.validPreparationsService.CreationInputMiddleware).Post(root, s.validPreparationsService.CreateHandler)
				validPreparationsRouter.Route(validPreparationRouteParam, func(singleValidPreparationRouter chi.Router) {
					singleValidPreparationRouter.Get(root, s.validPreparationsService.ReadHandler)
					singleValidPreparationRouter.With(s.validPreparationsService.UpdateInputMiddleware).Put(root, s.validPreparationsService.UpdateHandler)
					singleValidPreparationRouter.Delete(root, s.validPreparationsService.ArchiveHandler)
					singleValidPreparationRouter.Head(root, s.validPreparationsService.ExistenceHandler)
				})
				validPreparationsRouter.Get(root, s.validPreparationsService.ListHandler)
				validPreparationsRouter.Get(searchRoot, s.validPreparationsService.SearchHandler)
			})

			// ValidIngredientPreparations
			validIngredientPreparationPath := "valid_ingredient_preparations"
			validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
			validIngredientPreparationRouteParam := fmt.Sprintf(numericIDPattern, validingredientpreparationsservice.URIParamKey)
			v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter chi.Router) {
				validIngredientPreparationsRouter.With(s.validIngredientPreparationsService.CreationInputMiddleware).Post(root, s.validIngredientPreparationsService.CreateHandler)
				validIngredientPreparationsRouter.Route(validIngredientPreparationRouteParam, func(singleValidIngredientPreparationRouter chi.Router) {
					singleValidIngredientPreparationRouter.Get(root, s.validIngredientPreparationsService.ReadHandler)
					singleValidIngredientPreparationRouter.With(s.validIngredientPreparationsService.UpdateInputMiddleware).Put(root, s.validIngredientPreparationsService.UpdateHandler)
					singleValidIngredientPreparationRouter.Delete(root, s.validIngredientPreparationsService.ArchiveHandler)
					singleValidIngredientPreparationRouter.Head(root, s.validIngredientPreparationsService.ExistenceHandler)
				})
				validIngredientPreparationsRouter.Get(root, s.validIngredientPreparationsService.ListHandler)
			})

			// RequiredPreparationInstruments
			requiredPreparationInstrumentPath := "required_preparation_instruments"
			requiredPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", requiredPreparationInstrumentPath)
			requiredPreparationInstrumentRouteParam := fmt.Sprintf(numericIDPattern, requiredpreparationinstrumentsservice.URIParamKey)
			v1Router.Route(requiredPreparationInstrumentsRouteWithPrefix, func(requiredPreparationInstrumentsRouter chi.Router) {
				requiredPreparationInstrumentsRouter.With(s.requiredPreparationInstrumentsService.CreationInputMiddleware).Post(root, s.requiredPreparationInstrumentsService.CreateHandler)
				requiredPreparationInstrumentsRouter.Route(requiredPreparationInstrumentRouteParam, func(singleRequiredPreparationInstrumentRouter chi.Router) {
					singleRequiredPreparationInstrumentRouter.Get(root, s.requiredPreparationInstrumentsService.ReadHandler)
					singleRequiredPreparationInstrumentRouter.With(s.requiredPreparationInstrumentsService.UpdateInputMiddleware).Put(root, s.requiredPreparationInstrumentsService.UpdateHandler)
					singleRequiredPreparationInstrumentRouter.Delete(root, s.requiredPreparationInstrumentsService.ArchiveHandler)
					singleRequiredPreparationInstrumentRouter.Head(root, s.requiredPreparationInstrumentsService.ExistenceHandler)
				})
				requiredPreparationInstrumentsRouter.Get(root, s.requiredPreparationInstrumentsService.ListHandler)
			})

			// Recipes
			recipePath := "recipes"
			recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
			recipeRouteParam := fmt.Sprintf(numericIDPattern, recipesservice.URIParamKey)
			v1Router.Route(recipesRouteWithPrefix, func(recipesRouter chi.Router) {
				recipesRouter.With(s.recipesService.CreationInputMiddleware).Post(root, s.recipesService.CreateHandler)
				recipesRouter.Route(recipeRouteParam, func(singleRecipeRouter chi.Router) {
					singleRecipeRouter.Get(root, s.recipesService.ReadHandler)
					singleRecipeRouter.With(s.recipesService.UpdateInputMiddleware).Put(root, s.recipesService.UpdateHandler)
					singleRecipeRouter.Delete(root, s.recipesService.ArchiveHandler)
					singleRecipeRouter.Head(root, s.recipesService.ExistenceHandler)
				})
				recipesRouter.Get(root, s.recipesService.ListHandler)
			})

			// RecipeSteps
			recipeStepPath := "recipe_steps"
			recipeStepsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeStepPath,
			)
			recipeStepsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepsRoute)
			recipeStepRouteParam := fmt.Sprintf(numericIDPattern, recipestepsservice.URIParamKey)
			v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter chi.Router) {
				recipeStepsRouter.With(s.recipeStepsService.CreationInputMiddleware).Post(root, s.recipeStepsService.CreateHandler)
				recipeStepsRouter.Route(recipeStepRouteParam, func(singleRecipeStepRouter chi.Router) {
					singleRecipeStepRouter.Get(root, s.recipeStepsService.ReadHandler)
					singleRecipeStepRouter.With(s.recipeStepsService.UpdateInputMiddleware).Put(root, s.recipeStepsService.UpdateHandler)
					singleRecipeStepRouter.Delete(root, s.recipeStepsService.ArchiveHandler)
					singleRecipeStepRouter.Head(root, s.recipeStepsService.ExistenceHandler)
				})
				recipeStepsRouter.Get(root, s.recipeStepsService.ListHandler)
			})

			// RecipeStepInstruments
			recipeStepInstrumentPath := "recipe_step_instruments"
			recipeStepInstrumentsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeStepPath,
				recipeStepRouteParam,
				recipeStepInstrumentPath,
			)
			recipeStepInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepInstrumentsRoute)
			recipeStepInstrumentRouteParam := fmt.Sprintf(numericIDPattern, recipestepinstrumentsservice.URIParamKey)
			v1Router.Route(recipeStepInstrumentsRouteWithPrefix, func(recipeStepInstrumentsRouter chi.Router) {
				recipeStepInstrumentsRouter.With(s.recipeStepInstrumentsService.CreationInputMiddleware).Post(root, s.recipeStepInstrumentsService.CreateHandler)
				recipeStepInstrumentsRouter.Route(recipeStepInstrumentRouteParam, func(singleRecipeStepInstrumentRouter chi.Router) {
					singleRecipeStepInstrumentRouter.Get(root, s.recipeStepInstrumentsService.ReadHandler)
					singleRecipeStepInstrumentRouter.With(s.recipeStepInstrumentsService.UpdateInputMiddleware).Put(root, s.recipeStepInstrumentsService.UpdateHandler)
					singleRecipeStepInstrumentRouter.Delete(root, s.recipeStepInstrumentsService.ArchiveHandler)
					singleRecipeStepInstrumentRouter.Head(root, s.recipeStepInstrumentsService.ExistenceHandler)
				})
				recipeStepInstrumentsRouter.Get(root, s.recipeStepInstrumentsService.ListHandler)
			})

			// RecipeStepIngredients
			recipeStepIngredientPath := "recipe_step_ingredients"
			recipeStepIngredientsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeStepPath,
				recipeStepRouteParam,
				recipeStepIngredientPath,
			)
			recipeStepIngredientsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepIngredientsRoute)
			recipeStepIngredientRouteParam := fmt.Sprintf(numericIDPattern, recipestepingredientsservice.URIParamKey)
			v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter chi.Router) {
				recipeStepIngredientsRouter.With(s.recipeStepIngredientsService.CreationInputMiddleware).Post(root, s.recipeStepIngredientsService.CreateHandler)
				recipeStepIngredientsRouter.Route(recipeStepIngredientRouteParam, func(singleRecipeStepIngredientRouter chi.Router) {
					singleRecipeStepIngredientRouter.Get(root, s.recipeStepIngredientsService.ReadHandler)
					singleRecipeStepIngredientRouter.With(s.recipeStepIngredientsService.UpdateInputMiddleware).Put(root, s.recipeStepIngredientsService.UpdateHandler)
					singleRecipeStepIngredientRouter.Delete(root, s.recipeStepIngredientsService.ArchiveHandler)
					singleRecipeStepIngredientRouter.Head(root, s.recipeStepIngredientsService.ExistenceHandler)
				})
				recipeStepIngredientsRouter.Get(root, s.recipeStepIngredientsService.ListHandler)
			})

			// RecipeStepProducts
			recipeStepProductPath := "recipe_step_products"
			recipeStepProductsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeStepPath,
				recipeStepRouteParam,
				recipeStepProductPath,
			)
			recipeStepProductsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepProductsRoute)
			recipeStepProductRouteParam := fmt.Sprintf(numericIDPattern, recipestepproductsservice.URIParamKey)
			v1Router.Route(recipeStepProductsRouteWithPrefix, func(recipeStepProductsRouter chi.Router) {
				recipeStepProductsRouter.With(s.recipeStepProductsService.CreationInputMiddleware).Post(root, s.recipeStepProductsService.CreateHandler)
				recipeStepProductsRouter.Route(recipeStepProductRouteParam, func(singleRecipeStepProductRouter chi.Router) {
					singleRecipeStepProductRouter.Get(root, s.recipeStepProductsService.ReadHandler)
					singleRecipeStepProductRouter.With(s.recipeStepProductsService.UpdateInputMiddleware).Put(root, s.recipeStepProductsService.UpdateHandler)
					singleRecipeStepProductRouter.Delete(root, s.recipeStepProductsService.ArchiveHandler)
					singleRecipeStepProductRouter.Head(root, s.recipeStepProductsService.ExistenceHandler)
				})
				recipeStepProductsRouter.Get(root, s.recipeStepProductsService.ListHandler)
			})

			// RecipeIterations
			recipeIterationPath := "recipe_iterations"
			recipeIterationsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeIterationPath,
			)
			recipeIterationsRouteWithPrefix := fmt.Sprintf("/%s", recipeIterationsRoute)
			recipeIterationRouteParam := fmt.Sprintf(numericIDPattern, recipeiterationsservice.URIParamKey)
			v1Router.Route(recipeIterationsRouteWithPrefix, func(recipeIterationsRouter chi.Router) {
				recipeIterationsRouter.With(s.recipeIterationsService.CreationInputMiddleware).Post(root, s.recipeIterationsService.CreateHandler)
				recipeIterationsRouter.Route(recipeIterationRouteParam, func(singleRecipeIterationRouter chi.Router) {
					singleRecipeIterationRouter.Get(root, s.recipeIterationsService.ReadHandler)
					singleRecipeIterationRouter.With(s.recipeIterationsService.UpdateInputMiddleware).Put(root, s.recipeIterationsService.UpdateHandler)
					singleRecipeIterationRouter.Delete(root, s.recipeIterationsService.ArchiveHandler)
					singleRecipeIterationRouter.Head(root, s.recipeIterationsService.ExistenceHandler)
				})
				recipeIterationsRouter.Get(root, s.recipeIterationsService.ListHandler)
			})

			// RecipeStepEvents
			recipeStepEventPath := "recipe_step_events"
			recipeStepEventsRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeStepPath,
				recipeStepRouteParam,
				recipeStepEventPath,
			)
			recipeStepEventsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepEventsRoute)
			recipeStepEventRouteParam := fmt.Sprintf(numericIDPattern, recipestepeventsservice.URIParamKey)
			v1Router.Route(recipeStepEventsRouteWithPrefix, func(recipeStepEventsRouter chi.Router) {
				recipeStepEventsRouter.With(s.recipeStepEventsService.CreationInputMiddleware).Post(root, s.recipeStepEventsService.CreateHandler)
				recipeStepEventsRouter.Route(recipeStepEventRouteParam, func(singleRecipeStepEventRouter chi.Router) {
					singleRecipeStepEventRouter.Get(root, s.recipeStepEventsService.ReadHandler)
					singleRecipeStepEventRouter.With(s.recipeStepEventsService.UpdateInputMiddleware).Put(root, s.recipeStepEventsService.UpdateHandler)
					singleRecipeStepEventRouter.Delete(root, s.recipeStepEventsService.ArchiveHandler)
					singleRecipeStepEventRouter.Head(root, s.recipeStepEventsService.ExistenceHandler)
				})
				recipeStepEventsRouter.Get(root, s.recipeStepEventsService.ListHandler)
			})

			// IterationMedias
			iterationMediaPath := "iteration_medias"
			iterationMediasRoute := filepath.Join(
				recipePath,
				recipeRouteParam,
				recipeIterationPath,
				recipeIterationRouteParam,
				iterationMediaPath,
			)
			iterationMediasRouteWithPrefix := fmt.Sprintf("/%s", iterationMediasRoute)
			iterationMediaRouteParam := fmt.Sprintf(numericIDPattern, iterationmediasservice.URIParamKey)
			v1Router.Route(iterationMediasRouteWithPrefix, func(iterationMediasRouter chi.Router) {
				iterationMediasRouter.With(s.iterationMediasService.CreationInputMiddleware).Post(root, s.iterationMediasService.CreateHandler)
				iterationMediasRouter.Route(iterationMediaRouteParam, func(singleIterationMediaRouter chi.Router) {
					singleIterationMediaRouter.Get(root, s.iterationMediasService.ReadHandler)
					singleIterationMediaRouter.With(s.iterationMediasService.UpdateInputMiddleware).Put(root, s.iterationMediasService.UpdateHandler)
					singleIterationMediaRouter.Delete(root, s.iterationMediasService.ArchiveHandler)
					singleIterationMediaRouter.Head(root, s.iterationMediasService.ExistenceHandler)
				})
				iterationMediasRouter.Get(root, s.iterationMediasService.ListHandler)
			})

			// Invitations
			invitationPath := "invitations"
			invitationsRouteWithPrefix := fmt.Sprintf("/%s", invitationPath)
			invitationRouteParam := fmt.Sprintf(numericIDPattern, invitationsservice.URIParamKey)
			v1Router.Route(invitationsRouteWithPrefix, func(invitationsRouter chi.Router) {
				invitationsRouter.With(s.invitationsService.CreationInputMiddleware).Post(root, s.invitationsService.CreateHandler)
				invitationsRouter.Route(invitationRouteParam, func(singleInvitationRouter chi.Router) {
					singleInvitationRouter.Get(root, s.invitationsService.ReadHandler)
					singleInvitationRouter.With(s.invitationsService.UpdateInputMiddleware).Put(root, s.invitationsService.UpdateHandler)
					singleInvitationRouter.Delete(root, s.invitationsService.ArchiveHandler)
					singleInvitationRouter.Head(root, s.invitationsService.ExistenceHandler)
				})
				invitationsRouter.Get(root, s.invitationsService.ListHandler)
			})

			// Reports
			reportPath := "reports"
			reportsRouteWithPrefix := fmt.Sprintf("/%s", reportPath)
			reportRouteParam := fmt.Sprintf(numericIDPattern, reportsservice.URIParamKey)
			v1Router.Route(reportsRouteWithPrefix, func(reportsRouter chi.Router) {
				reportsRouter.With(s.reportsService.CreationInputMiddleware).Post(root, s.reportsService.CreateHandler)
				reportsRouter.Route(reportRouteParam, func(singleReportRouter chi.Router) {
					singleReportRouter.Get(root, s.reportsService.ReadHandler)
					singleReportRouter.With(s.reportsService.UpdateInputMiddleware).Put(root, s.reportsService.UpdateHandler)
					singleReportRouter.Delete(root, s.reportsService.ArchiveHandler)
					singleReportRouter.Head(root, s.reportsService.ExistenceHandler)
				})
				reportsRouter.Get(root, s.reportsService.ListHandler)
			})

			// Webhooks.
			v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, webhooksservice.URIParamKey)
				webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler)
				webhookRouter.Get(sr, s.webhooksService.ReadHandler)
				webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler)
				webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler)
				webhookRouter.Get(root, s.webhooksService.ListHandler)
			})

			// OAuth2 Clients.
			v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, oauth2clientsservice.URIParamKey)
				// CreateHandler is not bound to an OAuth2 authentication token.
				// UpdateHandler not supported for OAuth2 clients.
				clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler)
				clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler)
				clientRouter.Get(root, s.oauth2ClientsService.ListHandler)
			})
		})

	s.router = router
}
