package httpserver

import (
	"fmt"
	"net/http"
	"path/filepath"

	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	ingredienttagmappingsservice "gitlab.com/prixfixe/prixfixe/services/v1/ingredienttagmappings"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipeiterationstepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterationsteps"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipesteppreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteppreparations"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	recipetagsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipetags"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validingredienttagsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredienttags"
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

		userIDPattern := fmt.Sprintf(oauth2IDPattern, usersservice.URIParamKey)

		userRouter.Get(root, s.usersService.ListHandler())
		userRouter.With(s.usersService.UserInputMiddleware).Post(root, s.usersService.CreateHandler())
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
		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentRouteParam := fmt.Sprintf(numericIDPattern, validinstrumentsservice.URIParamKey)
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter chi.Router) {
			validInstrumentsRouter.With(s.validInstrumentsService.CreationInputMiddleware).Post(root, s.validInstrumentsService.CreateHandler())
			validInstrumentsRouter.Get(validInstrumentRouteParam, s.validInstrumentsService.ReadHandler())
			validInstrumentsRouter.Head(validInstrumentRouteParam, s.validInstrumentsService.ExistenceHandler())
			validInstrumentsRouter.With(s.validInstrumentsService.UpdateInputMiddleware).Put(validInstrumentRouteParam, s.validInstrumentsService.UpdateHandler())
			validInstrumentsRouter.Delete(validInstrumentRouteParam, s.validInstrumentsService.ArchiveHandler())
			validInstrumentsRouter.Get(root, s.validInstrumentsService.ListHandler())
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientRouteParam := fmt.Sprintf(numericIDPattern, validingredientsservice.URIParamKey)
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter chi.Router) {
			validIngredientsRouter.With(s.validIngredientsService.CreationInputMiddleware).Post(root, s.validIngredientsService.CreateHandler())
			validIngredientsRouter.Get(validIngredientRouteParam, s.validIngredientsService.ReadHandler())
			validIngredientsRouter.Head(validIngredientRouteParam, s.validIngredientsService.ExistenceHandler())
			validIngredientsRouter.With(s.validIngredientsService.UpdateInputMiddleware).Put(validIngredientRouteParam, s.validIngredientsService.UpdateHandler())
			validIngredientsRouter.Delete(validIngredientRouteParam, s.validIngredientsService.ArchiveHandler())
			validIngredientsRouter.Get(root, s.validIngredientsService.ListHandler())
		})

		// ValidIngredientTags
		validIngredientTagPath := "valid_ingredient_tags"
		validIngredientTagRouteParam := fmt.Sprintf(numericIDPattern, validingredienttagsservice.URIParamKey)
		validIngredientTagsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientTagPath)
		v1Router.Route(validIngredientTagsRouteWithPrefix, func(validIngredientTagsRouter chi.Router) {
			validIngredientTagsRouter.With(s.validIngredientTagsService.CreationInputMiddleware).Post(root, s.validIngredientTagsService.CreateHandler())
			validIngredientTagsRouter.Get(validIngredientTagRouteParam, s.validIngredientTagsService.ReadHandler())
			validIngredientTagsRouter.Head(validIngredientTagRouteParam, s.validIngredientTagsService.ExistenceHandler())
			validIngredientTagsRouter.With(s.validIngredientTagsService.UpdateInputMiddleware).Put(validIngredientTagRouteParam, s.validIngredientTagsService.UpdateHandler())
			validIngredientTagsRouter.Delete(validIngredientTagRouteParam, s.validIngredientTagsService.ArchiveHandler())
			validIngredientTagsRouter.Get(root, s.validIngredientTagsService.ListHandler())
		})

		// IngredientTagMappings
		ingredientTagMappingPath := "ingredient_tag_mappings"
		ingredientTagMappingRouteParam := fmt.Sprintf(numericIDPattern, ingredienttagmappingsservice.URIParamKey)
		ingredientTagMappingsRoute := filepath.Join(
			validIngredientPath,
			validIngredientRouteParam,
			ingredientTagMappingPath,
		)
		ingredientTagMappingsRouteWithPrefix := fmt.Sprintf("/%s", ingredientTagMappingsRoute)
		v1Router.Route(ingredientTagMappingsRouteWithPrefix, func(ingredientTagMappingsRouter chi.Router) {
			ingredientTagMappingsRouter.With(s.ingredientTagMappingsService.CreationInputMiddleware).Post(root, s.ingredientTagMappingsService.CreateHandler())
			ingredientTagMappingsRouter.Get(ingredientTagMappingRouteParam, s.ingredientTagMappingsService.ReadHandler())
			ingredientTagMappingsRouter.Head(ingredientTagMappingRouteParam, s.ingredientTagMappingsService.ExistenceHandler())
			ingredientTagMappingsRouter.With(s.ingredientTagMappingsService.UpdateInputMiddleware).Put(ingredientTagMappingRouteParam, s.ingredientTagMappingsService.UpdateHandler())
			ingredientTagMappingsRouter.Delete(ingredientTagMappingRouteParam, s.ingredientTagMappingsService.ArchiveHandler())
			ingredientTagMappingsRouter.Get(root, s.ingredientTagMappingsService.ListHandler())
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationRouteParam := fmt.Sprintf(numericIDPattern, validpreparationsservice.URIParamKey)
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter chi.Router) {
			validPreparationsRouter.With(s.validPreparationsService.CreationInputMiddleware).Post(root, s.validPreparationsService.CreateHandler())
			validPreparationsRouter.Get(validPreparationRouteParam, s.validPreparationsService.ReadHandler())
			validPreparationsRouter.Head(validPreparationRouteParam, s.validPreparationsService.ExistenceHandler())
			validPreparationsRouter.With(s.validPreparationsService.UpdateInputMiddleware).Put(validPreparationRouteParam, s.validPreparationsService.UpdateHandler())
			validPreparationsRouter.Delete(validPreparationRouteParam, s.validPreparationsService.ArchiveHandler())
			validPreparationsRouter.Get(root, s.validPreparationsService.ListHandler())
		})

		// RequiredPreparationInstruments
		requiredPreparationInstrumentPath := "required_preparation_instruments"
		requiredPreparationInstrumentRouteParam := fmt.Sprintf(numericIDPattern, requiredpreparationinstrumentsservice.URIParamKey)
		requiredPreparationInstrumentsRoute := filepath.Join(
			validPreparationPath,
			validPreparationRouteParam,
			requiredPreparationInstrumentPath,
		)
		requiredPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", requiredPreparationInstrumentsRoute)
		v1Router.Route(requiredPreparationInstrumentsRouteWithPrefix, func(requiredPreparationInstrumentsRouter chi.Router) {
			requiredPreparationInstrumentsRouter.With(s.requiredPreparationInstrumentsService.CreationInputMiddleware).Post(root, s.requiredPreparationInstrumentsService.CreateHandler())
			requiredPreparationInstrumentsRouter.Get(requiredPreparationInstrumentRouteParam, s.requiredPreparationInstrumentsService.ReadHandler())
			requiredPreparationInstrumentsRouter.Head(requiredPreparationInstrumentRouteParam, s.requiredPreparationInstrumentsService.ExistenceHandler())
			requiredPreparationInstrumentsRouter.With(s.requiredPreparationInstrumentsService.UpdateInputMiddleware).Put(requiredPreparationInstrumentRouteParam, s.requiredPreparationInstrumentsService.UpdateHandler())
			requiredPreparationInstrumentsRouter.Delete(requiredPreparationInstrumentRouteParam, s.requiredPreparationInstrumentsService.ArchiveHandler())
			requiredPreparationInstrumentsRouter.Get(root, s.requiredPreparationInstrumentsService.ListHandler())
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationRouteParam := fmt.Sprintf(numericIDPattern, validingredientpreparationsservice.URIParamKey)
		validIngredientPreparationsRoute := filepath.Join(
			validIngredientPath,
			validIngredientRouteParam,
			validIngredientPreparationPath,
		)
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationsRoute)
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter chi.Router) {
			validIngredientPreparationsRouter.With(s.validIngredientPreparationsService.CreationInputMiddleware).Post(root, s.validIngredientPreparationsService.CreateHandler())
			validIngredientPreparationsRouter.Get(validIngredientPreparationRouteParam, s.validIngredientPreparationsService.ReadHandler())
			validIngredientPreparationsRouter.Head(validIngredientPreparationRouteParam, s.validIngredientPreparationsService.ExistenceHandler())
			validIngredientPreparationsRouter.With(s.validIngredientPreparationsService.UpdateInputMiddleware).Put(validIngredientPreparationRouteParam, s.validIngredientPreparationsService.UpdateHandler())
			validIngredientPreparationsRouter.Delete(validIngredientPreparationRouteParam, s.validIngredientPreparationsService.ArchiveHandler())
			validIngredientPreparationsRouter.Get(root, s.validIngredientPreparationsService.ListHandler())
		})

		// Recipes
		recipePath := "recipes"
		recipeRouteParam := fmt.Sprintf(numericIDPattern, recipesservice.URIParamKey)
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter chi.Router) {
			recipesRouter.With(s.recipesService.CreationInputMiddleware).Post(root, s.recipesService.CreateHandler())
			recipesRouter.Get(recipeRouteParam, s.recipesService.ReadHandler())
			recipesRouter.Head(recipeRouteParam, s.recipesService.ExistenceHandler())
			recipesRouter.With(s.recipesService.UpdateInputMiddleware).Put(recipeRouteParam, s.recipesService.UpdateHandler())
			recipesRouter.Delete(recipeRouteParam, s.recipesService.ArchiveHandler())
			recipesRouter.Get(root, s.recipesService.ListHandler())
		})

		// RecipeTags
		recipeTagPath := "recipe_tags"
		recipeTagRouteParam := fmt.Sprintf(numericIDPattern, recipetagsservice.URIParamKey)
		recipeTagsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeTagPath,
		)
		recipeTagsRouteWithPrefix := fmt.Sprintf("/%s", recipeTagsRoute)
		v1Router.Route(recipeTagsRouteWithPrefix, func(recipeTagsRouter chi.Router) {
			recipeTagsRouter.With(s.recipeTagsService.CreationInputMiddleware).Post(root, s.recipeTagsService.CreateHandler())
			recipeTagsRouter.Get(recipeTagRouteParam, s.recipeTagsService.ReadHandler())
			recipeTagsRouter.Head(recipeTagRouteParam, s.recipeTagsService.ExistenceHandler())
			recipeTagsRouter.With(s.recipeTagsService.UpdateInputMiddleware).Put(recipeTagRouteParam, s.recipeTagsService.UpdateHandler())
			recipeTagsRouter.Delete(recipeTagRouteParam, s.recipeTagsService.ArchiveHandler())
			recipeTagsRouter.Get(root, s.recipeTagsService.ListHandler())
		})

		// RecipeSteps
		recipeStepPath := "recipe_steps"
		recipeStepRouteParam := fmt.Sprintf(numericIDPattern, recipestepsservice.URIParamKey)
		recipeStepsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeStepPath,
		)
		recipeStepsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepsRoute)
		v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter chi.Router) {
			recipeStepsRouter.With(s.recipeStepsService.CreationInputMiddleware).Post(root, s.recipeStepsService.CreateHandler())
			recipeStepsRouter.Get(recipeStepRouteParam, s.recipeStepsService.ReadHandler())
			recipeStepsRouter.Head(recipeStepRouteParam, s.recipeStepsService.ExistenceHandler())
			recipeStepsRouter.With(s.recipeStepsService.UpdateInputMiddleware).Put(recipeStepRouteParam, s.recipeStepsService.UpdateHandler())
			recipeStepsRouter.Delete(recipeStepRouteParam, s.recipeStepsService.ArchiveHandler())
			recipeStepsRouter.Get(root, s.recipeStepsService.ListHandler())
		})

		// RecipeStepPreparations
		recipeStepPreparationPath := "recipe_step_preparations"
		recipeStepPreparationRouteParam := fmt.Sprintf(numericIDPattern, recipesteppreparationsservice.URIParamKey)
		recipeStepPreparationsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeStepPath,
			recipeStepRouteParam,
			recipeStepPreparationPath,
		)
		recipeStepPreparationsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepPreparationsRoute)
		v1Router.Route(recipeStepPreparationsRouteWithPrefix, func(recipeStepPreparationsRouter chi.Router) {
			recipeStepPreparationsRouter.With(s.recipeStepPreparationsService.CreationInputMiddleware).Post(root, s.recipeStepPreparationsService.CreateHandler())
			recipeStepPreparationsRouter.Get(recipeStepPreparationRouteParam, s.recipeStepPreparationsService.ReadHandler())
			recipeStepPreparationsRouter.Head(recipeStepPreparationRouteParam, s.recipeStepPreparationsService.ExistenceHandler())
			recipeStepPreparationsRouter.With(s.recipeStepPreparationsService.UpdateInputMiddleware).Put(recipeStepPreparationRouteParam, s.recipeStepPreparationsService.UpdateHandler())
			recipeStepPreparationsRouter.Delete(recipeStepPreparationRouteParam, s.recipeStepPreparationsService.ArchiveHandler())
			recipeStepPreparationsRouter.Get(root, s.recipeStepPreparationsService.ListHandler())
		})

		// RecipeStepIngredients
		recipeStepIngredientPath := "recipe_step_ingredients"
		recipeStepIngredientRouteParam := fmt.Sprintf(numericIDPattern, recipestepingredientsservice.URIParamKey)
		recipeStepIngredientsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeStepPath,
			recipeStepRouteParam,
			recipeStepIngredientPath,
		)
		recipeStepIngredientsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepIngredientsRoute)
		v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter chi.Router) {
			recipeStepIngredientsRouter.With(s.recipeStepIngredientsService.CreationInputMiddleware).Post(root, s.recipeStepIngredientsService.CreateHandler())
			recipeStepIngredientsRouter.Get(recipeStepIngredientRouteParam, s.recipeStepIngredientsService.ReadHandler())
			recipeStepIngredientsRouter.Head(recipeStepIngredientRouteParam, s.recipeStepIngredientsService.ExistenceHandler())
			recipeStepIngredientsRouter.With(s.recipeStepIngredientsService.UpdateInputMiddleware).Put(recipeStepIngredientRouteParam, s.recipeStepIngredientsService.UpdateHandler())
			recipeStepIngredientsRouter.Delete(recipeStepIngredientRouteParam, s.recipeStepIngredientsService.ArchiveHandler())
			recipeStepIngredientsRouter.Get(root, s.recipeStepIngredientsService.ListHandler())
		})

		// RecipeIterations
		recipeIterationPath := "recipe_iterations"
		recipeIterationRouteParam := fmt.Sprintf(numericIDPattern, recipeiterationsservice.URIParamKey)
		recipeIterationsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeIterationPath,
		)
		recipeIterationsRouteWithPrefix := fmt.Sprintf("/%s", recipeIterationsRoute)
		v1Router.Route(recipeIterationsRouteWithPrefix, func(recipeIterationsRouter chi.Router) {
			recipeIterationsRouter.With(s.recipeIterationsService.CreationInputMiddleware).Post(root, s.recipeIterationsService.CreateHandler())
			recipeIterationsRouter.Get(recipeIterationRouteParam, s.recipeIterationsService.ReadHandler())
			recipeIterationsRouter.Head(recipeIterationRouteParam, s.recipeIterationsService.ExistenceHandler())
			recipeIterationsRouter.With(s.recipeIterationsService.UpdateInputMiddleware).Put(recipeIterationRouteParam, s.recipeIterationsService.UpdateHandler())
			recipeIterationsRouter.Delete(recipeIterationRouteParam, s.recipeIterationsService.ArchiveHandler())
			recipeIterationsRouter.Get(root, s.recipeIterationsService.ListHandler())
		})

		// RecipeIterationSteps
		recipeIterationStepPath := "recipe_iteration_steps"
		recipeIterationStepRouteParam := fmt.Sprintf(numericIDPattern, recipeiterationstepsservice.URIParamKey)
		recipeIterationStepsRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeIterationStepPath,
		)
		recipeIterationStepsRouteWithPrefix := fmt.Sprintf("/%s", recipeIterationStepsRoute)
		v1Router.Route(recipeIterationStepsRouteWithPrefix, func(recipeIterationStepsRouter chi.Router) {
			recipeIterationStepsRouter.With(s.recipeIterationStepsService.CreationInputMiddleware).Post(root, s.recipeIterationStepsService.CreateHandler())
			recipeIterationStepsRouter.Get(recipeIterationStepRouteParam, s.recipeIterationStepsService.ReadHandler())
			recipeIterationStepsRouter.Head(recipeIterationStepRouteParam, s.recipeIterationStepsService.ExistenceHandler())
			recipeIterationStepsRouter.With(s.recipeIterationStepsService.UpdateInputMiddleware).Put(recipeIterationStepRouteParam, s.recipeIterationStepsService.UpdateHandler())
			recipeIterationStepsRouter.Delete(recipeIterationStepRouteParam, s.recipeIterationStepsService.ArchiveHandler())
			recipeIterationStepsRouter.Get(root, s.recipeIterationStepsService.ListHandler())
		})

		// IterationMedias
		iterationMediaPath := "iteration_medias"
		iterationMediaRouteParam := fmt.Sprintf(numericIDPattern, iterationmediasservice.URIParamKey)
		iterationMediasRoute := filepath.Join(
			recipePath,
			recipeRouteParam,
			recipeIterationPath,
			recipeIterationRouteParam,
			iterationMediaPath,
		)
		iterationMediasRouteWithPrefix := fmt.Sprintf("/%s", iterationMediasRoute)
		v1Router.Route(iterationMediasRouteWithPrefix, func(iterationMediasRouter chi.Router) {
			iterationMediasRouter.With(s.iterationMediasService.CreationInputMiddleware).Post(root, s.iterationMediasService.CreateHandler())
			iterationMediasRouter.Get(iterationMediaRouteParam, s.iterationMediasService.ReadHandler())
			iterationMediasRouter.Head(iterationMediaRouteParam, s.iterationMediasService.ExistenceHandler())
			iterationMediasRouter.With(s.iterationMediasService.UpdateInputMiddleware).Put(iterationMediaRouteParam, s.iterationMediasService.UpdateHandler())
			iterationMediasRouter.Delete(iterationMediaRouteParam, s.iterationMediasService.ArchiveHandler())
			iterationMediasRouter.Get(root, s.iterationMediasService.ListHandler())
		})

		// Invitations
		invitationPath := "invitations"
		invitationRouteParam := fmt.Sprintf(numericIDPattern, invitationsservice.URIParamKey)
		invitationsRouteWithPrefix := fmt.Sprintf("/%s", invitationPath)
		v1Router.Route(invitationsRouteWithPrefix, func(invitationsRouter chi.Router) {
			invitationsRouter.With(s.invitationsService.CreationInputMiddleware).Post(root, s.invitationsService.CreateHandler())
			invitationsRouter.Get(invitationRouteParam, s.invitationsService.ReadHandler())
			invitationsRouter.Head(invitationRouteParam, s.invitationsService.ExistenceHandler())
			invitationsRouter.With(s.invitationsService.UpdateInputMiddleware).Put(invitationRouteParam, s.invitationsService.UpdateHandler())
			invitationsRouter.Delete(invitationRouteParam, s.invitationsService.ArchiveHandler())
			invitationsRouter.Get(root, s.invitationsService.ListHandler())
		})

		// Reports
		reportPath := "reports"
		reportRouteParam := fmt.Sprintf(numericIDPattern, reportsservice.URIParamKey)
		reportsRouteWithPrefix := fmt.Sprintf("/%s", reportPath)
		v1Router.Route(reportsRouteWithPrefix, func(reportsRouter chi.Router) {
			reportsRouter.With(s.reportsService.CreationInputMiddleware).Post(root, s.reportsService.CreateHandler())
			reportsRouter.Get(reportRouteParam, s.reportsService.ReadHandler())
			reportsRouter.Head(reportRouteParam, s.reportsService.ExistenceHandler())
			reportsRouter.With(s.reportsService.UpdateInputMiddleware).Put(reportRouteParam, s.reportsService.UpdateHandler())
			reportsRouter.Delete(reportRouteParam, s.reportsService.ArchiveHandler())
			reportsRouter.Get(root, s.reportsService.ListHandler())
		})

		// Webhooks.
		v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, webhooksservice.URIParamKey)
			webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler())
			webhookRouter.Get(sr, s.webhooksService.ReadHandler())
			webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler())
			webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler())
			webhookRouter.Get(root, s.webhooksService.ListHandler())
		})

		// OAuth2 Clients.
		v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, oauth2clientsservice.URIParamKey)
			// CreateHandler is not bound to an OAuth2 authentication token.
			// UpdateHandler not supported for OAuth2 clients.
			clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler())
			clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler())
			clientRouter.Get(root, s.oauth2ClientsService.ListHandler())
		})
	})

	s.router = router
}
