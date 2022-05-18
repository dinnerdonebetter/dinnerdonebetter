package server

import (
	"context"
	"fmt"
	"path"

	"github.com/prixfixeco/api_server/internal/observability/metrics"

	"github.com/heptiolabs/healthcheck"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/routing"
	apiclientsservice "github.com/prixfixeco/api_server/internal/services/apiclients"
	householdinvitationsservice "github.com/prixfixeco/api_server/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/api_server/internal/services/mealplans"
	mealsservice "github.com/prixfixeco/api_server/internal/services/meals"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/api_server/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/api_server/internal/services/users"
	validingredientpreparationsservice "github.com/prixfixeco/api_server/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/api_server/internal/services/validingredients"
	validinstrumentsservice "github.com/prixfixeco/api_server/internal/services/validinstruments"
	validpreparationsservice "github.com/prixfixeco/api_server/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/api_server/internal/services/webhooks"
)

const (
	root       = "/"
	randomRoot = "/random"
	searchRoot = "/search"
)

func buildURLVarChunk(key, pattern string) string {
	if pattern != "" {
		return fmt.Sprintf("/{%s:%s}", key, pattern)
	}
	return fmt.Sprintf("/{%s}", key)
}

func (s *HTTPServer) setupRouter(ctx context.Context, router routing.Router, metricsHandler metrics.Handler) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	router.Route("/_meta_", func(metaRouter routing.Router) {
		health := healthcheck.NewHandler()
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		logger.Info("setting up metrics handler")
		router.Get("/metrics", metricsHandler.ServeHTTP)
	}

	router.Post("/paseto", s.authService.PASETOHandler)

	authenticatedRouter := router.WithMiddleware(s.authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", s.authService.StatusHandler)

	router.Route("/users", func(userRouter routing.Router) {
		userRouter.Post("/login", s.authService.BuildLoginHandler(false))
		userRouter.Post("/login/admin", s.authService.BuildLoginHandler(true))
		userRouter.WithMiddleware(s.authService.UserAttributionMiddleware, s.authService.CookieRequirementMiddleware).Post("/logout", s.authService.EndSessionHandler)
		userRouter.Post(root, s.usersService.CreateHandler)
		userRouter.Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)

		// need credentials beyond this point
		authedRouter := userRouter.WithMiddleware(s.authService.UserAttributionMiddleware, s.authService.AuthorizationMiddleware)
		authedRouter.Post("/household/select", s.authService.ChangeActiveHouseholdHandler)
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

			singleUserRoute := buildURLVarChunk(usersservice.UserIDURIParamKey, "")
			usersRouter.Route(singleUserRoute, func(singleUserRouter routing.Router) {
				singleUserRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
					Get(root, s.usersService.ReadHandler)

				singleUserRouter.Delete(root, s.usersService.ArchiveHandler)
			})
		})

		// Households
		v1Router.Route("/households", func(householdsRouter routing.Router) {
			householdsRouter.Post(root, s.householdsService.CreateHandler)
			householdsRouter.Get(root, s.householdsService.ListHandler)
			householdsRouter.Get("/current", s.householdsService.InfoHandler)

			singleUserRoute := buildURLVarChunk(householdsservice.UserIDURIParamKey, "")
			singleHouseholdRoute := buildURLVarChunk(householdsservice.HouseholdIDURIParamKey, "")
			householdsRouter.Route(singleHouseholdRoute, func(singleHouseholdRouter routing.Router) {
				singleHouseholdRouter.Get(root, s.householdsService.ReadHandler)
				singleHouseholdRouter.Put(root, s.householdsService.UpdateHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdPermission)).
					Delete(root, s.householdsService.ArchiveHandler)

				singleHouseholdRouter.Post("/default", s.householdsService.MarkAsDefaultHouseholdHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.RemoveMemberHouseholdPermission)).
					Delete("/members"+singleUserRoute, s.householdsService.RemoveMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)).
					Post("/invite", s.householdInvitationsService.InviteMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ModifyMemberPermissionsForHouseholdPermission)).
					Patch("/members"+singleUserRoute+"/permissions", s.householdsService.ModifyMemberPermissionsHandler)
				singleHouseholdRouter.Post("/transfer", s.householdsService.TransferHouseholdOwnershipHandler)

				singleHouseholdRouter.Route("/invitations", func(invitationsRouter routing.Router) {
					invitationsRouter.Post(root, s.householdInvitationsService.InviteMemberHandler)

					singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey, "")
					invitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
						singleHouseholdInvitationRouter.Get(root, s.householdInvitationsService.ReadHandler)
						singleHouseholdInvitationRouter.Put("/cancel", s.householdInvitationsService.CancelInviteHandler)
						singleHouseholdInvitationRouter.Put("/accept", s.householdInvitationsService.AcceptInviteHandler)
						singleHouseholdInvitationRouter.Put("/reject", s.householdInvitationsService.RejectInviteHandler)
					})
				})
			})
		})

		v1Router.Route("/household_invitations", func(householdInvitationsRouter routing.Router) {
			householdInvitationsRouter.Get("/sent", s.householdInvitationsService.OutboundInvitesHandler)
			householdInvitationsRouter.Get("/received", s.householdInvitationsService.InboundInvitesHandler)
		})

		// API Clients
		v1Router.Route("/api_clients", func(clientRouter routing.Router) {
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
				Get(root, s.apiClientsService.ListHandler)
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateAPIClientsPermission)).
				Post(root, s.apiClientsService.CreateHandler)

			singleClientRoute := buildURLVarChunk(apiclientsservice.APIClientIDURIParamKey, "")
			clientRouter.Route(singleClientRoute, func(singleClientRouter routing.Router) {
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
					Get(root, s.apiClientsService.ReadHandler)
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveAPIClientsPermission)).
					Delete(root, s.apiClientsService.ArchiveHandler)
			})
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter routing.Router) {
			singleWebhookRoute := buildURLVarChunk(webhooksservice.WebhookIDURIParamKey, "")
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
			})
		})

		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		validInstrumentIDRouteParam := buildURLVarChunk(validinstrumentsservice.ValidInstrumentIDURIParamKey, "")
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
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(randomRoot, s.validInstrumentsService.RandomHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, s.validInstrumentsService.ReadHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, s.validInstrumentsService.ArchiveHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, s.validInstrumentsService.UpdateHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildURLVarChunk(validingredientsservice.ValidIngredientIDURIParamKey, "")
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
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(randomRoot, s.validIngredientsService.RandomHandler)

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, s.validIngredientsService.ReadHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, s.validIngredientsService.ArchiveHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, s.validIngredientsService.UpdateHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildURLVarChunk(validpreparationsservice.ValidPreparationIDURIParamKey, "")
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
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(randomRoot, s.validPreparationsService.RandomHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, s.validPreparationsService.ReadHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, s.validPreparationsService.ArchiveHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, s.validPreparationsService.UpdateHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildURLVarChunk(validingredientpreparationsservice.ValidIngredientPreparationIDURIParamKey, "")
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
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, s.validIngredientPreparationsService.ArchiveHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, s.validIngredientPreparationsService.UpdateHandler)
			})
		})

		// Meals
		mealPath := "meals"
		mealsRouteWithPrefix := fmt.Sprintf("/%s", mealPath)
		mealIDRouteParam := buildURLVarChunk(mealsservice.MealIDURIParamKey, "")
		v1Router.Route(mealsRouteWithPrefix, func(mealsRouter routing.Router) {
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealsPermission)).
				Post(root, s.mealsService.CreateHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(root, s.mealsService.ListHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get("/search", s.mealsService.SearchHandler)

			mealsRouter.Route(mealIDRouteParam, func(singleMealRouter routing.Router) {
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
					Get(root, s.mealsService.ReadHandler)
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealsPermission)).
					Delete(root, s.mealsService.ArchiveHandler)
			})
		})

		// Recipes
		recipePath := "recipes"
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		recipeIDRouteParam := buildURLVarChunk(recipesservice.RecipeIDURIParamKey, "")
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, s.recipesService.CreateHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, s.recipesService.ListHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get("/search", s.recipesService.SearchHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, s.recipesService.ReadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, s.recipesService.ArchiveHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, s.recipesService.UpdateHandler)
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
		recipeStepIDRouteParam := buildURLVarChunk(recipestepsservice.RecipeStepIDURIParamKey, "")
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
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, s.recipeStepsService.ArchiveHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, s.recipeStepsService.UpdateHandler)
			})
		})

		// RecipeStepInstruments
		recipeStepInstrumentPath := "recipe_step_instruments"
		recipeStepInstrumentsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepInstrumentPath,
		)
		recipeStepInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepInstrumentsRoute)
		recipeStepInstrumentIDRouteParam := buildURLVarChunk(recipestepinstrumentsservice.RecipeStepInstrumentIDURIParamKey, "")
		v1Router.Route(recipeStepInstrumentsRouteWithPrefix, func(recipeStepInstrumentsRouter routing.Router) {
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepInstrumentsPermission)).
				Post(root, s.recipeStepInstrumentsService.CreateHandler)
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
				Get(root, s.recipeStepInstrumentsService.ListHandler)

			recipeStepInstrumentsRouter.Route(recipeStepInstrumentIDRouteParam, func(singleRecipeStepInstrumentRouter routing.Router) {
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
					Get(root, s.recipeStepInstrumentsService.ReadHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepInstrumentsPermission)).
					Delete(root, s.recipeStepInstrumentsService.ArchiveHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepInstrumentsPermission)).
					Put(root, s.recipeStepInstrumentsService.UpdateHandler)
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
		recipeStepIngredientIDRouteParam := buildURLVarChunk(recipestepingredientsservice.RecipeStepIngredientIDURIParamKey, "")
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
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, s.recipeStepIngredientsService.ArchiveHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, s.recipeStepIngredientsService.UpdateHandler)
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
		recipeStepProductIDRouteParam := buildURLVarChunk(recipestepproductsservice.RecipeStepProductIDURIParamKey, "")
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
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, s.recipeStepProductsService.ArchiveHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, s.recipeStepProductsService.UpdateHandler)
			})
		})

		// MealPlans
		mealPlanPath := "meal_plans"
		mealPlansRouteWithPrefix := fmt.Sprintf("/%s", mealPlanPath)
		mealPlanIDRouteParam := buildURLVarChunk(mealplansservice.MealPlanIDURIParamKey, "")
		v1Router.Route(mealPlansRouteWithPrefix, func(mealPlansRouter routing.Router) {
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
				Post(root, s.mealPlansService.CreateHandler)
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
				Get(root, s.mealPlansService.ListHandler)

			mealPlansRouter.Route(mealPlanIDRouteParam, func(singleMealPlanRouter routing.Router) {
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
					Get(root, s.mealPlansService.ReadHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlansPermission)).
					Delete(root, s.mealPlansService.ArchiveHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlansPermission)).
					Put(root, s.mealPlansService.UpdateHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionVotesPermission)).
					Post("/vote", s.mealPlanOptionVotesService.CreateHandler)
			})
		})

		// MealPlanOptions
		mealPlanOptionPath := "meal_plan_options"
		mealPlanOptionsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanOptionPath,
		)
		mealPlanOptionsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionsRoute)
		mealPlanOptionIDRouteParam := buildURLVarChunk(mealplanoptionsservice.MealPlanOptionIDURIParamKey, "")
		v1Router.Route(mealPlanOptionsRouteWithPrefix, func(mealPlanOptionsRouter routing.Router) {
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionsPermission)).
				Post(root, s.mealPlanOptionsService.CreateHandler)
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
				Get(root, s.mealPlanOptionsService.ListHandler)

			mealPlanOptionsRouter.Route(mealPlanOptionIDRouteParam, func(singleMealPlanOptionRouter routing.Router) {
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
					Get(root, s.mealPlanOptionsService.ReadHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionsPermission)).
					Delete(root, s.mealPlanOptionsService.ArchiveHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionsPermission)).
					Put(root, s.mealPlanOptionsService.UpdateHandler)
			})
		})

		// MealPlanOptionVotes
		mealPlanOptionVotePath := "meal_plan_option_votes"
		mealPlanOptionVotesRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanOptionPath,
			mealPlanOptionIDRouteParam,
			mealPlanOptionVotePath,
		)
		mealPlanOptionVotesRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionVotesRoute)
		mealPlanOptionVoteIDRouteParam := buildURLVarChunk(mealplanoptionvotesservice.MealPlanOptionVoteIDURIParamKey, "")
		v1Router.Route(mealPlanOptionVotesRouteWithPrefix, func(mealPlanOptionVotesRouter routing.Router) {
			mealPlanOptionVotesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
				Get(root, s.mealPlanOptionVotesService.ListHandler)

			mealPlanOptionVotesRouter.Route(mealPlanOptionVoteIDRouteParam, func(singleMealPlanOptionVoteRouter routing.Router) {
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
					Get(root, s.mealPlanOptionVotesService.ReadHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionVotesPermission)).
					Delete(root, s.mealPlanOptionVotesService.ArchiveHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionVotesPermission)).
					Put(root, s.mealPlanOptionVotesService.UpdateHandler)
			})
		})
	})

	s.router = router
}
