package server

import (
	"context"
	"fmt"
	"path"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	accountsservice "gitlab.com/prixfixe/prixfixe/internal/services/accounts"
	apiclientsservice "gitlab.com/prixfixe/prixfixe/internal/services/apiclients"
	auditservice "gitlab.com/prixfixe/prixfixe/internal/services/audit"
	invitationsservice "gitlab.com/prixfixe/prixfixe/internal/services/invitations"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepingredients"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/internal/services/reports"
	usersservice "gitlab.com/prixfixe/prixfixe/internal/services/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/internal/services/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validinstruments"
	validpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparationinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/internal/services/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/internal/services/webhooks"

	"github.com/heptiolabs/healthcheck"
)

const (
	root             = "/"
	auditRoute       = "/audit"
	searchRoot       = "/search"
	numericIDPattern = "{%s:[0-9]+}"
)

func buildNumericIDURLChunk(key string) string {
	return fmt.Sprintf(root+numericIDPattern, key)
}

func (s *HTTPServer) setupRouter(ctx context.Context, router routing.Router, metricsHandler metrics.Handler) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	router.Route("/_meta_", func(metaRouter routing.Router) {
		health := healthcheck.NewHandler()
		// Expose a liveness check on /live
		metaRouter.Get("/live", health.LiveEndpoint)
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.HandleFunc("/metrics", metricsHandler.ServeHTTP)
	}

	// Frontend routes.
	s.frontendService.SetupRoutes(router)

	router.Post("/paseto", s.authService.PASETOHandler)

	authenticatedRouter := router.WithMiddleware(s.authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", s.authService.StatusHandler)

	router.Route("/users", func(userRouter routing.Router) {
		userRouter.Post("/login", s.authService.BeginSessionHandler)
		userRouter.WithMiddleware(s.authService.UserAttributionMiddleware, s.authService.CookieRequirementMiddleware).Post("/logout", s.authService.EndSessionHandler)
		userRouter.Post(root, s.usersService.CreateHandler)
		userRouter.Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)

		// need credentials beyond this point
		authedRouter := userRouter.WithMiddleware(s.authService.UserAttributionMiddleware, s.authService.AuthorizationMiddleware)
		authedRouter.Post("/account/select", s.authService.ChangeActiveAccountHandler)
		authedRouter.Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)
		authedRouter.Put("/password/new", s.usersService.UpdatePasswordHandler)
	})

	authenticatedRouter.WithMiddleware(s.authService.AuthorizationMiddleware).Route("/api/v1", func(v1Router routing.Router) {
		adminRouter := v1Router.WithMiddleware(s.authService.ServiceAdminMiddleware)

		// Admin
		adminRouter.Route("/admin", func(adminRouter routing.Router) {
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CycleCookieSecretPermission)).
				Post("/cycle_cookie_secret", s.authService.CycleCookieSecretHandler)
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/users/status", s.adminService.UserReputationChangeHandler)

			adminRouter.Route("/audit_log", func(auditRouter routing.Router) {
				entryIDRouteParam := buildNumericIDURLChunk(auditservice.LogEntryURIParamKey)
				auditRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAllAuditLogEntriesPermission)).
					Get(root, s.auditService.ListHandler)
				auditRouter.Route(entryIDRouteParam, func(singleEntryRouter routing.Router) {
					singleEntryRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAllAuditLogEntriesPermission)).
						Get(root, s.auditService.ReadHandler)
				})
			})
		})

		// Users
		v1Router.Route("/users", func(usersRouter routing.Router) {
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
				Get(root, s.usersService.ListHandler)
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
				Get("/search", s.usersService.UsernameSearchHandler)
			usersRouter.Post("/avatar/upload", s.usersService.AvatarUploadHandler)
			usersRouter.Get("/self", s.usersService.SelfHandler)

			singleUserRoute := buildNumericIDURLChunk(usersservice.UserIDURIParamKey)
			usersRouter.Route(singleUserRoute, func(singleUserRouter routing.Router) {
				singleUserRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
					Get(root, s.usersService.ReadHandler)
				singleUserRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserAuditLogEntriesPermission)).
					Get(auditRoute, s.usersService.AuditEntryHandler)

				singleUserRouter.Delete(root, s.usersService.ArchiveHandler)
			})
		})

		// Accounts
		v1Router.Route("/accounts", func(accountsRouter routing.Router) {
			accountsRouter.Post(root, s.accountsService.CreateHandler)
			accountsRouter.Get(root, s.accountsService.ListHandler)

			singleUserRoute := buildNumericIDURLChunk(accountsservice.UserIDURIParamKey)
			singleAccountRoute := buildNumericIDURLChunk(accountsservice.AccountIDURIParamKey)
			accountsRouter.Route(singleAccountRoute, func(singleAccountRouter routing.Router) {
				singleAccountRouter.Get(root, s.accountsService.ReadHandler)
				singleAccountRouter.Put(root, s.accountsService.UpdateHandler)
				singleAccountRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveAccountPermission)).
					Delete(root, s.accountsService.ArchiveHandler)
				singleAccountRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAccountAuditLogEntriesPermission)).
					Get(auditRoute, s.accountsService.AuditEntryHandler)

				singleAccountRouter.Post("/default", s.accountsService.MarkAsDefaultAccountHandler)
				singleAccountRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.RemoveMemberAccountPermission)).
					Delete("/members"+singleUserRoute, s.accountsService.RemoveMemberHandler)
				singleAccountRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.AddMemberAccountPermission)).
					Post("/member", s.accountsService.AddMemberHandler)
				singleAccountRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ModifyMemberPermissionsForAccountPermission)).
					Patch("/members"+singleUserRoute+"/permissions", s.accountsService.ModifyMemberPermissionsHandler)
				singleAccountRouter.Post("/transfer", s.accountsService.TransferAccountOwnershipHandler)
			})
		})

		// API Clients
		v1Router.Route("/api_clients", func(clientRouter routing.Router) {
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
				Get(root, s.apiClientsService.ListHandler)
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateAPIClientsPermission)).
				Post(root, s.apiClientsService.CreateHandler)

			singleClientRoute := buildNumericIDURLChunk(apiclientsservice.APIClientIDURIParamKey)
			clientRouter.Route(singleClientRoute, func(singleClientRouter routing.Router) {
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
					Get(root, s.apiClientsService.ReadHandler)
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveAPIClientsPermission)).
					Delete(root, s.apiClientsService.ArchiveHandler)
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientAuditLogEntriesPermission)).
					Get(auditRoute, s.apiClientsService.AuditEntryHandler)
			})
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter routing.Router) {
			singleWebhookRoute := buildNumericIDURLChunk(webhooksservice.WebhookIDURIParamKey)
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
				Get(root, s.webhooksService.ListHandler)
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateWebhooksPermission)).
				Post(root, s.webhooksService.CreateHandler)
			webhookRouter.Route(singleWebhookRoute, func(singleWebhookRouter routing.Router) {
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
					Get(root, s.webhooksService.ReadHandler)
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveWebhooksPermission)).
					Delete(root, s.webhooksService.ArchiveHandler)
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateWebhooksPermission)).
					Put(root, s.webhooksService.UpdateHandler)
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksAuditLogEntriesPermission)).
					Get(auditRoute, s.webhooksService.AuditEntryHandler)
			})
		})

		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		validInstrumentIDRouteParam := buildNumericIDURLChunk(validinstrumentsservice.ValidInstrumentIDURIParamKey)
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter routing.Router) {
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
				Post(root, s.validInstrumentsService.CreateHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(root, s.validInstrumentsService.ListHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(searchRoot, s.validInstrumentsService.SearchHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, s.validInstrumentsService.ReadHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Head(root, s.validInstrumentsService.ExistenceHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, s.validInstrumentsService.ArchiveHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, s.validInstrumentsService.UpdateHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsAuditLogEntriesPermission)).
					Get(auditRoute, s.validInstrumentsService.AuditEntryHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildNumericIDURLChunk(validpreparationsservice.ValidPreparationIDURIParamKey)
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter routing.Router) {
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
				Post(root, s.validPreparationsService.CreateHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(root, s.validPreparationsService.ListHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(searchRoot, s.validPreparationsService.SearchHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, s.validPreparationsService.ReadHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Head(root, s.validPreparationsService.ExistenceHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, s.validPreparationsService.ArchiveHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, s.validPreparationsService.UpdateHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsAuditLogEntriesPermission)).
					Get(auditRoute, s.validPreparationsService.AuditEntryHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildNumericIDURLChunk(validingredientsservice.ValidIngredientIDURIParamKey)
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter routing.Router) {
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
				Post(root, s.validIngredientsService.CreateHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(root, s.validIngredientsService.ListHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(searchRoot, s.validIngredientsService.SearchHandler)

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, s.validIngredientsService.ReadHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Head(root, s.validIngredientsService.ExistenceHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, s.validIngredientsService.ArchiveHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, s.validIngredientsService.UpdateHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsAuditLogEntriesPermission)).
					Get(auditRoute, s.validIngredientsService.AuditEntryHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildNumericIDURLChunk(validingredientpreparationsservice.ValidIngredientPreparationIDURIParamKey)
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter routing.Router) {
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
				Post(root, s.validIngredientPreparationsService.CreateHandler)
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
				Get(root, s.validIngredientPreparationsService.ListHandler)

			validIngredientPreparationsRouter.Route(validIngredientPreparationIDRouteParam, func(singleValidIngredientPreparationRouter routing.Router) {
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.ReadHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Head(root, s.validIngredientPreparationsService.ExistenceHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, s.validIngredientPreparationsService.ArchiveHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, s.validIngredientPreparationsService.UpdateHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsAuditLogEntriesPermission)).
					Get(auditRoute, s.validIngredientPreparationsService.AuditEntryHandler)
			})
		})

		// ValidPreparationInstruments
		validPreparationInstrumentPath := "valid_preparation_instruments"
		validPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationInstrumentPath)
		validPreparationInstrumentIDRouteParam := buildNumericIDURLChunk(validpreparationinstrumentsservice.ValidPreparationInstrumentIDURIParamKey)
		v1Router.Route(validPreparationInstrumentsRouteWithPrefix, func(validPreparationInstrumentsRouter routing.Router) {
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
				Post(root, s.validPreparationInstrumentsService.CreateHandler)
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
				Get(root, s.validPreparationInstrumentsService.ListHandler)

			validPreparationInstrumentsRouter.Route(validPreparationInstrumentIDRouteParam, func(singleValidPreparationInstrumentRouter routing.Router) {
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.ReadHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Head(root, s.validPreparationInstrumentsService.ExistenceHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
					Delete(root, s.validPreparationInstrumentsService.ArchiveHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
					Put(root, s.validPreparationInstrumentsService.UpdateHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsAuditLogEntriesPermission)).
					Get(auditRoute, s.validPreparationInstrumentsService.AuditEntryHandler)
			})
		})

		// Recipes
		recipePath := "recipes"
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		recipeIDRouteParam := buildNumericIDURLChunk(recipesservice.RecipeIDURIParamKey)
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, s.recipesService.CreateHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, s.recipesService.ListHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, s.recipesService.ReadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Head(root, s.recipesService.ExistenceHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, s.recipesService.ArchiveHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, s.recipesService.UpdateHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesAuditLogEntriesPermission)).
					Get(auditRoute, s.recipesService.AuditEntryHandler)
			})
		})

		// RecipeSteps
		recipeStepPath := "recipe_steps"
		recipeStepsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
		)
		recipeStepsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepsRoute)
		recipeStepIDRouteParam := buildNumericIDURLChunk(recipestepsservice.RecipeStepIDURIParamKey)
		v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter routing.Router) {
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepsPermission)).
				Post(root, s.recipeStepsService.CreateHandler)
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
				Get(root, s.recipeStepsService.ListHandler)

			recipeStepsRouter.Route(recipeStepIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Get(root, s.recipeStepsService.ReadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Head(root, s.recipeStepsService.ExistenceHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, s.recipeStepsService.ArchiveHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, s.recipeStepsService.UpdateHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsAuditLogEntriesPermission)).
					Get(auditRoute, s.recipeStepsService.AuditEntryHandler)
			})
		})

		// RecipeStepIngredients
		recipeStepIngredientPath := "recipe_step_ingredients"
		recipeStepIngredientsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepIngredientPath,
		)
		recipeStepIngredientsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepIngredientsRoute)
		recipeStepIngredientIDRouteParam := buildNumericIDURLChunk(recipestepingredientsservice.RecipeStepIngredientIDURIParamKey)
		v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter routing.Router) {
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepIngredientsPermission)).
				Post(root, s.recipeStepIngredientsService.CreateHandler)
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
				Get(root, s.recipeStepIngredientsService.ListHandler)

			recipeStepIngredientsRouter.Route(recipeStepIngredientIDRouteParam, func(singleRecipeStepIngredientRouter routing.Router) {
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Get(root, s.recipeStepIngredientsService.ReadHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Head(root, s.recipeStepIngredientsService.ExistenceHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, s.recipeStepIngredientsService.ArchiveHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, s.recipeStepIngredientsService.UpdateHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsAuditLogEntriesPermission)).
					Get(auditRoute, s.recipeStepIngredientsService.AuditEntryHandler)
			})
		})

		// RecipeStepProducts
		recipeStepProductPath := "recipe_step_products"
		recipeStepProductsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepProductPath,
		)
		recipeStepProductsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepProductsRoute)
		recipeStepProductIDRouteParam := buildNumericIDURLChunk(recipestepproductsservice.RecipeStepProductIDURIParamKey)
		v1Router.Route(recipeStepProductsRouteWithPrefix, func(recipeStepProductsRouter routing.Router) {
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepProductsPermission)).
				Post(root, s.recipeStepProductsService.CreateHandler)
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
				Get(root, s.recipeStepProductsService.ListHandler)

			recipeStepProductsRouter.Route(recipeStepProductIDRouteParam, func(singleRecipeStepProductRouter routing.Router) {
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Get(root, s.recipeStepProductsService.ReadHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Head(root, s.recipeStepProductsService.ExistenceHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, s.recipeStepProductsService.ArchiveHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, s.recipeStepProductsService.UpdateHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsAuditLogEntriesPermission)).
					Get(auditRoute, s.recipeStepProductsService.AuditEntryHandler)
			})
		})

		// Invitations
		invitationPath := "invitations"
		invitationsRouteWithPrefix := fmt.Sprintf("/%s", invitationPath)
		invitationIDRouteParam := buildNumericIDURLChunk(invitationsservice.InvitationIDURIParamKey)
		v1Router.Route(invitationsRouteWithPrefix, func(invitationsRouter routing.Router) {
			invitationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateInvitationsPermission)).
				Post(root, s.invitationsService.CreateHandler)
			invitationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsPermission)).
				Get(root, s.invitationsService.ListHandler)

			invitationsRouter.Route(invitationIDRouteParam, func(singleInvitationRouter routing.Router) {
				singleInvitationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsPermission)).
					Get(root, s.invitationsService.ReadHandler)
				singleInvitationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsPermission)).
					Head(root, s.invitationsService.ExistenceHandler)
				singleInvitationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveInvitationsPermission)).
					Delete(root, s.invitationsService.ArchiveHandler)
				singleInvitationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateInvitationsPermission)).
					Put(root, s.invitationsService.UpdateHandler)
				singleInvitationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsAuditLogEntriesPermission)).
					Get(auditRoute, s.invitationsService.AuditEntryHandler)
			})
		})

		// Reports
		reportPath := "reports"
		reportsRouteWithPrefix := fmt.Sprintf("/%s", reportPath)
		reportIDRouteParam := buildNumericIDURLChunk(reportsservice.ReportIDURIParamKey)
		v1Router.Route(reportsRouteWithPrefix, func(reportsRouter routing.Router) {
			reportsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateReportsPermission)).
				Post(root, s.reportsService.CreateHandler)
			reportsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsPermission)).
				Get(root, s.reportsService.ListHandler)

			reportsRouter.Route(reportIDRouteParam, func(singleReportRouter routing.Router) {
				singleReportRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsPermission)).
					Get(root, s.reportsService.ReadHandler)
				singleReportRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsPermission)).
					Head(root, s.reportsService.ExistenceHandler)
				singleReportRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveReportsPermission)).
					Delete(root, s.reportsService.ArchiveHandler)
				singleReportRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateReportsPermission)).
					Put(root, s.reportsService.UpdateHandler)
				singleReportRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsAuditLogEntriesPermission)).
					Get(auditRoute, s.reportsService.AuditEntryHandler)
			})
		})
	})

	s.router = router
}
