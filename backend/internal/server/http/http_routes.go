package http

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/routing"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplangrocerylistitemssservice "github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	mealplantasksservice "github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
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

func (s *server) setupRouter(ctx context.Context, router routing.Router) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	router.Route("/_meta_", func(metaRouter routing.Router) {
		// Expose a readiness check on /ready
		metaRouter.Get("/live", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			reqCtx, reqSpan := s.tracer.StartSpan(req.Context())
			defer reqSpan.End()

			logger := s.logger.WithRequest(req)

			responseCode := http.StatusOK
			if err := s.dataManager.DB().PingContext(reqCtx); err != nil {
				logger.Error(err, "database not responding to ping")
				responseCode = http.StatusInternalServerError
			}

			res.WriteHeader(responseCode)
		})
	})

	authenticatedRouter := router.WithMiddleware(s.authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", s.authService.StatusHandler)

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.Get("/authorize", s.authService.AuthorizeHandler)
		userRouter.Post("/token", s.authService.TokenHandler)
	})

	// these are routes we don't expect or require users to be authenticated for
	router.Route("/users", func(userRouter routing.Router) {
		userRouter.Post(root, s.usersService.CreateUserHandler)
		userRouter.Post("/login/jwt", s.authService.BuildLoginHandler(false))
		userRouter.Post("/login/jwt/admin", s.authService.BuildLoginHandler(true))
		userRouter.Post("/username/reminder", s.usersService.RequestUsernameReminderHandler)
		userRouter.Post("/password/reset", s.usersService.CreatePasswordResetTokenHandler)
		userRouter.Post("/password/reset/redeem", s.usersService.PasswordResetTokenRedemptionHandler)
		userRouter.Post("/email_address/verify", s.usersService.VerifyUserEmailAddressHandler)
		userRouter.Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
	})

	router.Route("/auth", func(authRouter routing.Router) {
		providerRouteParam := buildURLVarChunk(authservice.AuthProviderParamKey, "")
		authRouter.Get(providerRouteParam, s.authService.SSOLoginHandler)
		authRouter.Get(path.Join(providerRouteParam, "callback"), s.authService.SSOLoginCallbackHandler)
	})

	authenticatedRouter.WithMiddleware(s.authService.AuthorizationMiddleware).Route("/api/v1", func(v1Router routing.Router) {
		adminRouter := v1Router.WithMiddleware(s.authService.ServiceAdminMiddleware)

		// Admin
		adminRouter.Route("/admin", func(adminRouter routing.Router) {
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/users/status", s.adminService.UserAccountStatusChangeHandler)
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/queues/test", s.adminService.UserAccountStatusChangeHandler)
		})

		// Workers
		adminRouter.Route("/workers", func(adminRouter routing.Router) {
			adminRouter.
				Post("/finalize_meal_plans", s.workerService.MealPlanFinalizationHandler)
			adminRouter.
				Post("/meal_plan_grocery_list_init", s.workerService.MealPlanGroceryListInitializationHandler)
			adminRouter.
				Post("/meal_plan_tasks", s.workerService.MealPlanTaskCreationHandler)
		})

		// Users
		v1Router.Route("/users", func(usersRouter routing.Router) {
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
				Get(root, s.usersService.ListUsersHandler)
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
				Get(searchRoot, s.usersService.UsernameSearchHandler)
			usersRouter.Post("/avatar/upload", s.usersService.AvatarUploadHandler)

			usersRouter.Get("/self", s.usersService.SelfHandler)
			usersRouter.Post("/email_address_verification", s.usersService.RequestEmailVerificationEmailHandler)
			usersRouter.Post("/permissions/check", s.usersService.UserPermissionsHandler)
			usersRouter.Put("/password/new", s.usersService.UpdatePasswordHandler)
			usersRouter.Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)
			usersRouter.Put("/username", s.usersService.UpdateUserUsernameHandler)
			usersRouter.Put("/email_address", s.usersService.UpdateUserEmailAddressHandler)
			usersRouter.Put("/details", s.usersService.UpdateUserDetailsHandler)

			singleUserRoute := buildURLVarChunk(usersservice.UserIDURIParamKey, "")
			usersRouter.Route(singleUserRoute, func(singleUserRouter routing.Router) {
				singleUserRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
					Get(root, s.usersService.ReadUserHandler)

				singleUserRouter.Delete(root, s.usersService.ArchiveUserHandler)
			})
		})

		v1Router.Route("/data_privacy", func(dataPrivacyRouter routing.Router) {
			dataPrivacyRouter.Delete("/destroy", s.dataPrivacyService.DataDeletionHandler)
			dataPrivacyRouter.Post("/disclose", s.dataPrivacyService.UserDataAggregationRequestHandler)

			singleReportRoute := buildURLVarChunk(dataprivacyservice.ReportIDURIParamKey, "")
			dataPrivacyRouter.Route("/reports", func(singleReportRouter routing.Router) {
				singleReportRouter.Get(singleReportRoute, s.dataPrivacyService.ReadUserDataAggregationReportHandler)
			})
		})

		// Households
		v1Router.Route("/households", func(householdsRouter routing.Router) {
			householdsRouter.Post(root, s.householdsService.CreateHouseholdHandler)
			householdsRouter.Get(root, s.householdsService.ListHouseholdsHandler)
			householdsRouter.Get("/current", s.householdsService.CurrentInfoHandler)

			singleUserRoute := buildURLVarChunk(householdsservice.UserIDURIParamKey, "")
			singleHouseholdRoute := buildURLVarChunk(householdsservice.HouseholdIDURIParamKey, "")
			householdsRouter.Route(singleHouseholdRoute, func(singleHouseholdRouter routing.Router) {
				singleHouseholdRouter.Get(root, s.householdsService.ReadHouseholdHandler)
				singleHouseholdRouter.Put(root, s.householdsService.UpdateHouseholdHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdPermission)).
					Delete(root, s.householdsService.ArchiveHouseholdHandler)

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
						singleHouseholdInvitationRouter.Get(root, s.householdInvitationsService.ReadHouseholdInviteHandler)
					})
				})
			})

			// HouseholdInstrumentOwnerships
			householdInstrumentOwnershipsRouteWithPrefix := "/instruments"
			householdInstrumentOwnershipIDRouteParam := buildURLVarChunk(householdinstrumentownershipsservice.HouseholdInstrumentOwnershipIDURIParamKey, "")
			householdsRouter.Route(householdInstrumentOwnershipsRouteWithPrefix, func(householdInstrumentOwnershipsRouter routing.Router) {
				householdInstrumentOwnershipsRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateHouseholdInstrumentOwnershipsPermission)).
					Post(root, s.householdInstrumentOwnershipService.CreateHouseholdInstrumentOwnershipHandler)
				householdInstrumentOwnershipsRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
					Get(root, s.householdInstrumentOwnershipService.ListHouseholdInstrumentOwnershipHandler)

				householdInstrumentOwnershipsRouter.Route(householdInstrumentOwnershipIDRouteParam, func(singleHouseholdInstrumentOwnershipRouter routing.Router) {
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
						Get(root, s.householdInstrumentOwnershipService.ReadHouseholdInstrumentOwnershipHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateHouseholdInstrumentOwnershipsPermission)).
						Put(root, s.householdInstrumentOwnershipService.UpdateHouseholdInstrumentOwnershipHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdInstrumentOwnershipsPermission)).
						Delete(root, s.householdInstrumentOwnershipService.ArchiveHouseholdInstrumentOwnershipHandler)
				})
			})
		})

		v1Router.Route("/household_invitations", func(householdInvitationsRouter routing.Router) {
			householdInvitationsRouter.Get("/sent", s.householdInvitationsService.OutboundInvitesHandler)
			householdInvitationsRouter.Get("/received", s.householdInvitationsService.InboundInvitesHandler)

			singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey, "")
			householdInvitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
				singleHouseholdInvitationRouter.Get(root, s.householdInvitationsService.ReadHouseholdInviteHandler)
				singleHouseholdInvitationRouter.Put("/cancel", s.householdInvitationsService.CancelInviteHandler)
				singleHouseholdInvitationRouter.Put("/accept", s.householdInvitationsService.AcceptInviteHandler)
				singleHouseholdInvitationRouter.Put("/reject", s.householdInvitationsService.RejectInviteHandler)
			})
		})

		// OAuth2 Clients
		v1Router.Route("/oauth2_clients", func(clientRouter routing.Router) {
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
				Get(root, s.oauth2ClientsService.ListOAuth2ClientsHandler)
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateOAuth2ClientsPermission)).
				Post(root, s.oauth2ClientsService.CreateOAuth2ClientHandler)

			singleClientRoute := buildURLVarChunk(oauth2clientsservice.OAuth2ClientIDURIParamKey, "")
			clientRouter.Route(singleClientRoute, func(singleClientRouter routing.Router) {
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
					Get(root, s.oauth2ClientsService.ReadOAuth2ClientHandler)
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveOAuth2ClientsPermission)).
					Delete(root, s.oauth2ClientsService.ArchiveOAuth2ClientHandler)
			})
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter routing.Router) {
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
				Get(root, s.webhooksService.ListWebhooksHandler)
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateWebhooksPermission)).
				Post(root, s.webhooksService.CreateWebhookHandler)

			singleWebhookRoute := buildURLVarChunk(webhooksservice.WebhookIDURIParamKey, "")
			webhookRouter.Route(singleWebhookRoute, func(singleWebhookRouter routing.Router) {
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
					Get(root, s.webhooksService.ReadWebhookHandler)
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveWebhooksPermission)).
					Delete(root, s.webhooksService.ArchiveWebhookHandler)

				singleWebhookTriggerEventRoute := buildURLVarChunk(webhooksservice.WebhookTriggerEventIDURIParamKey, "")
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateWebhookTriggerEventsPermission)).
					Post("/trigger_events", s.webhooksService.AddWebhookTriggerEventHandler)

				singleWebhookRouter.Route("/trigger_events"+singleWebhookTriggerEventRoute, func(singleWebhookTriggerEventRouter routing.Router) {
					singleWebhookTriggerEventRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveWebhookTriggerEventsPermission)).
						Delete(root, s.webhooksService.ArchiveWebhookTriggerEventHandler)
				})
			})
		})

		// Audit Log Entries
		v1Router.Route("/audit_log_entries", func(auditLogEntriesRouter routing.Router) {
			singleAuditLogEntryRoute := buildURLVarChunk(auditlogentriesservice.AuditLogEntryIDURIParamKey, "")
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get(singleAuditLogEntryRoute, s.auditLogEntriesService.ReadAuditLogEntryHandler)
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_user", s.auditLogEntriesService.ListUserAuditLogEntriesHandler)
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_household", s.auditLogEntriesService.ListHouseholdAuditLogEntriesHandler)
		})

		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		validInstrumentIDRouteParam := buildURLVarChunk(validinstrumentsservice.ValidInstrumentIDURIParamKey, "")
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter routing.Router) {
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
				Post(root, s.validInstrumentsService.CreateValidInstrumentHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(root, s.validInstrumentsService.ListValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidInstrumentsPermission)).
				Get(searchRoot, s.validInstrumentsService.SearchValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(randomRoot, s.validInstrumentsService.RandomValidInstrumentHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, s.validInstrumentsService.ReadValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, s.validInstrumentsService.UpdateValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, s.validInstrumentsService.ArchiveValidInstrumentHandler)
			})
		})

		// ValidVessels
		validVesselPath := "valid_vessels"
		validVesselsRouteWithPrefix := fmt.Sprintf("/%s", validVesselPath)
		validVesselIDRouteParam := buildURLVarChunk(validvesselsservice.ValidVesselIDURIParamKey, "")
		v1Router.Route(validVesselsRouteWithPrefix, func(validVesselsRouter routing.Router) {
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidVesselsPermission)).
				Post(root, s.validVesselsService.CreateValidVesselHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(root, s.validVesselsService.ListValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidVesselsPermission)).
				Get(searchRoot, s.validVesselsService.SearchValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(randomRoot, s.validVesselsService.RandomValidVesselHandler)

			validVesselsRouter.Route(validVesselIDRouteParam, func(singleValidVesselRouter routing.Router) {
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
					Get(root, s.validVesselsService.ReadValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidVesselsPermission)).
					Put(root, s.validVesselsService.UpdateValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidVesselsPermission)).
					Delete(root, s.validVesselsService.ArchiveValidVesselHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildURLVarChunk(validingredientsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter routing.Router) {
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
				Post(root, s.validIngredientsService.CreateValidIngredientHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(root, s.validIngredientsService.ListValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission, authorization.SearchValidIngredientsPermission)).
				Get(searchRoot, s.validIngredientsService.SearchValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(randomRoot, s.validIngredientsService.RandomValidIngredientHandler)

			validIngredientsByPreparationIDSearchRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validingredientsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientsRouter.Route(validIngredientsByPreparationIDSearchRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientsService.SearchValidIngredientsByPreparationAndIngredientNameHandler)
			})

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, s.validIngredientsService.ReadValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, s.validIngredientsService.UpdateValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, s.validIngredientsService.ArchiveValidIngredientHandler)
			})
		})

		// ValidIngredientGroups
		validIngredientGroupPath := "valid_ingredient_groups"
		validIngredientGroupsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientGroupPath)
		validIngredientGroupIDRouteParam := buildURLVarChunk(validingredientgroupsservice.ValidIngredientGroupIDURIParamKey, "")
		v1Router.Route(validIngredientGroupsRouteWithPrefix, func(validIngredientGroupsRouter routing.Router) {
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientGroupsPermission)).
				Post(root, s.validIngredientGroupsService.CreateValidIngredientGroupHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
				Get(root, s.validIngredientGroupsService.ListValidIngredientGroupsHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidIngredientGroupsPermission)).
				Get(searchRoot, s.validIngredientGroupsService.SearchValidIngredientGroupsHandler)

			validIngredientGroupsRouter.Route(validIngredientGroupIDRouteParam, func(singleValidIngredientGroupRouter routing.Router) {
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
					Get(root, s.validIngredientGroupsService.ReadValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientGroupsPermission)).
					Put(root, s.validIngredientGroupsService.UpdateValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientGroupsPermission)).
					Delete(root, s.validIngredientGroupsService.ArchiveValidIngredientGroupHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildURLVarChunk(validpreparationsservice.ValidPreparationIDURIParamKey, "")
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter routing.Router) {
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
				Post(root, s.validPreparationsService.CreateValidPreparationHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(root, s.validPreparationsService.ListValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(searchRoot, s.validPreparationsService.SearchValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(randomRoot, s.validPreparationsService.RandomValidPreparationHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, s.validPreparationsService.ReadValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, s.validPreparationsService.UpdateValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, s.validPreparationsService.ArchiveValidPreparationHandler)
			})
		})

		// UserIngredientPreferences
		userIngredientPreferencesPath := "user_ingredient_preferences"
		userIngredientPreferencesRouteWithPrefix := fmt.Sprintf("/%s", userIngredientPreferencesPath)
		userIngredientPreferencesIDRouteParam := buildURLVarChunk(useringredientpreferencesservice.UserIngredientPreferenceIDURIParamKey, "")
		v1Router.Route(userIngredientPreferencesRouteWithPrefix, func(userIngredientPreferencesRouter routing.Router) {
			userIngredientPreferencesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateUserIngredientPreferencesPermission)).
				Post(root, s.userIngredientPreferencesService.CreateUserIngredientPreferenceHandler)
			userIngredientPreferencesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserIngredientPreferencesPermission)).
				Get(root, s.userIngredientPreferencesService.ListUserIngredientPreferencesHandler)

			userIngredientPreferencesRouter.Route(userIngredientPreferencesIDRouteParam, func(singleUserIngredientPreferenceRouter routing.Router) {
				singleUserIngredientPreferenceRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserIngredientPreferencesPermission)).
					Put(root, s.userIngredientPreferencesService.UpdateUserIngredientPreferenceHandler)
				singleUserIngredientPreferenceRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveUserIngredientPreferencesPermission)).
					Delete(root, s.userIngredientPreferencesService.ArchiveUserIngredientPreferenceHandler)
			})
		})

		// ValidMeasurementUnits
		validMeasurementUnitPath := "valid_measurement_units"
		validMeasurementUnitsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitPath)
		validMeasurementUnitIDRouteParam := buildURLVarChunk(validmeasurementunitsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitServiceIngredientIDRouteParam := buildURLVarChunk(validmeasurementunitsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitsRouteWithPrefix, func(validMeasurementUnitsRouter routing.Router) {
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitsPermission)).
				Post(root, s.validMeasurementUnitsService.CreateValidMeasurementUnitHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
				Get(root, s.validMeasurementUnitsService.ListValidMeasurementUnitsHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(searchRoot, s.validMeasurementUnitsService.SearchValidMeasurementUnitsHandler)

			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(path.Join("/by_ingredient", validMeasurementUnitServiceIngredientIDRouteParam), s.validMeasurementUnitsService.SearchValidMeasurementUnitsByIngredientIDHandler)

			validMeasurementUnitsRouter.Route(validMeasurementUnitIDRouteParam, func(singleValidMeasurementUnitRouter routing.Router) {
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
					Get(root, s.validMeasurementUnitsService.ReadValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitsPermission)).
					Put(root, s.validMeasurementUnitsService.UpdateValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitsPermission)).
					Delete(root, s.validMeasurementUnitsService.ArchiveValidMeasurementUnitHandler)
			})
		})

		// ValidIngredientStates
		validIngredientStatePath := "valid_ingredient_states"
		validIngredientStatesRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStatePath)
		validIngredientStateIDRouteParam := buildURLVarChunk(validingredientstatesservice.ValidIngredientStateIDURIParamKey, "")
		v1Router.Route(validIngredientStatesRouteWithPrefix, func(validIngredientStatesRouter routing.Router) {
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStatesPermission)).
				Post(root, s.validIngredientStatesService.CreateValidIngredientStateHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(root, s.validIngredientStatesService.ListValidIngredientStatesHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(searchRoot, s.validIngredientStatesService.SearchValidIngredientStatesHandler)

			validIngredientStatesRouter.Route(validIngredientStateIDRouteParam, func(singleValidIngredientStateRouter routing.Router) {
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
					Get(root, s.validIngredientStatesService.ReadValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStatesPermission)).
					Put(root, s.validIngredientStatesService.UpdateValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStatesPermission)).
					Delete(root, s.validIngredientStatesService.ArchiveValidIngredientStateHandler)
			})
		})

		// ValidMeasurementUnitConversions
		validMeasurementUnitConversionPath := "valid_measurement_conversions"
		validMeasurementUnitConversionsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitConversionPath)
		validMeasurementUnitConversionUnitIDRouteParam := buildURLVarChunk(validmeasurementconversionsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitConversionIDRouteParam := buildURLVarChunk(validmeasurementconversionsservice.ValidMeasurementUnitConversionIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitConversionsRouteWithPrefix, func(validMeasurementUnitConversionsRouter routing.Router) {
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitConversionsPermission)).
				Post(root, s.validMeasurementUnitConversionsService.CreateValidMeasurementUnitConversionHandler)

			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/from_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validMeasurementUnitConversionsService.ValidMeasurementUnitConversionsFromMeasurementUnitHandler)
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/to_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validMeasurementUnitConversionsService.ValidMeasurementUnitConversionsToMeasurementUnitHandler)

			validMeasurementUnitConversionsRouter.Route(validMeasurementUnitConversionIDRouteParam, func(singleValidMeasurementUnitConversionRouter routing.Router) {
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
					Get(root, s.validMeasurementUnitConversionsService.ReadValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitConversionsPermission)).
					Put(root, s.validMeasurementUnitConversionsService.UpdateValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitConversionsPermission)).
					Delete(root, s.validMeasurementUnitConversionsService.ArchiveValidMeasurementUnitConversionHandler)
			})
		})

		// ValidIngredientStateIngredients
		validIngredientStateIngredientPath := "valid_ingredient_state_ingredients"
		validIngredientStateIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStateIngredientPath)
		validIngredientStateIngredientIDRouteParam := buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientStateIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientStateIngredientsRouteWithPrefix, func(validIngredientStateIngredientsRouter routing.Router) {
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStateIngredientsPermission)).
				Post(root, s.validIngredientStateIngredientsService.CreateValidIngredientStateIngredientHandler)
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
				Get(root, s.validIngredientStateIngredientsService.ListValidIngredientStateIngredientsHandler)

			validIngredientStateIngredientsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.SearchValidIngredientStateIngredientsByIngredientHandler)
			})

			validIngredientStateIngredientsByStateIngredientIDRouteParam := fmt.Sprintf("/by_ingredient_state%s", buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientStateIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByStateIngredientIDRouteParam, func(byValidStateIngredientIDRouter routing.Router) {
				byValidStateIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.SearchValidIngredientStateIngredientsByIngredientStateHandler)
			})

			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientIDRouteParam, func(singleValidIngredientStateIngredientRouter routing.Router) {
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.ReadValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStateIngredientsPermission)).
					Put(root, s.validIngredientStateIngredientsService.UpdateValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStateIngredientsPermission)).
					Delete(root, s.validIngredientStateIngredientsService.ArchiveValidIngredientStateIngredientHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildURLVarChunk(validingredientpreparationsservice.ValidIngredientPreparationIDURIParamKey, "")
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter routing.Router) {
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
				Post(root, s.validIngredientPreparationsService.CreateValidIngredientPreparationHandler)
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
				Get(root, s.validIngredientPreparationsService.ListValidIngredientPreparationsHandler)

			validIngredientPreparationsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientpreparationsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.SearchValidIngredientPreparationsByIngredientHandler)
			})

			validIngredientPreparationsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validingredientpreparationsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.SearchValidIngredientPreparationsByPreparationHandler)
			})

			validIngredientPreparationsRouter.Route(validIngredientPreparationIDRouteParam, func(singleValidIngredientPreparationRouter routing.Router) {
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.ReadValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, s.validIngredientPreparationsService.UpdateValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, s.validIngredientPreparationsService.ArchiveValidIngredientPreparationHandler)
			})
		})

		// ValidPreparationInstruments
		validPreparationInstrumentPath := "valid_preparation_instruments"
		validPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationInstrumentPath)
		v1Router.Route(validPreparationInstrumentsRouteWithPrefix, func(validPreparationInstrumentsRouter routing.Router) {
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
				Post(root, s.validPreparationInstrumentsService.CreateValidPreparationInstrumentHandler)
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
				Get(root, s.validPreparationInstrumentsService.ListValidPreparationInstrumentsHandler)

			validPreparationInstrumentsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validpreparationinstrumentsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.SearchValidPreparationInstrumentsByPreparationHandler)
			})

			validPreparationInstrumentsByInstrumentIDRouteParam := fmt.Sprintf("/by_instrument%s", buildURLVarChunk(validpreparationinstrumentsservice.ValidInstrumentIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.SearchValidPreparationInstrumentsByInstrumentHandler)
			})

			validPreparationInstrumentIDRouteParam := buildURLVarChunk(validpreparationinstrumentsservice.ValidPreparationVesselIDURIParamKey, "")
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentIDRouteParam, func(singleValidPreparationInstrumentRouter routing.Router) {
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.ReadValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
					Put(root, s.validPreparationInstrumentsService.UpdateValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
					Delete(root, s.validPreparationInstrumentsService.ArchiveValidPreparationInstrumentHandler)
			})
		})

		// ValidPreparationVessels
		validPreparationVesselPath := "valid_preparation_vessels"
		validPreparationVesselsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationVesselPath)
		v1Router.Route(validPreparationVesselsRouteWithPrefix, func(validPreparationVesselsRouter routing.Router) {
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationVesselsPermission)).
				Post(root, s.validPreparationVesselsService.CreateValidPreparationVesselHandler)
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
				Get(root, s.validPreparationVesselsService.ListValidPreparationVesselsHandler)

			validPreparationVesselsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validpreparationvesselsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.SearchValidPreparationVesselsByPreparationHandler)
			})

			validPreparationVesselsByInstrumentIDRouteParam := fmt.Sprintf("/by_vessel%s", buildURLVarChunk(validpreparationvesselsservice.ValidVesselIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.SearchValidPreparationVesselsByVesselHandler)
			})

			validPreparationVesselIDRouteParam := buildURLVarChunk(validpreparationvesselsservice.ValidPreparationVesselIDURIParamKey, "")
			validPreparationVesselsRouter.Route(validPreparationVesselIDRouteParam, func(singleValidPreparationVesselRouter routing.Router) {
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.ReadValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationVesselsPermission)).
					Put(root, s.validPreparationVesselsService.UpdateValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationVesselsPermission)).
					Delete(root, s.validPreparationVesselsService.ArchiveValidPreparationVesselHandler)
			})
		})

		// ValidIngredientMeasurementUnit
		validIngredientMeasurementUnitPath := "valid_ingredient_measurement_units"
		validIngredientMeasurementUnitRouteWithPrefix := fmt.Sprintf("/%s", validIngredientMeasurementUnitPath)
		validIngredientMeasurementUnitIDRouteParam := buildURLVarChunk(validingredientmeasurementunitsservice.ValidIngredientMeasurementUnitIDURIParamKey, "")
		v1Router.Route(validIngredientMeasurementUnitRouteWithPrefix, func(validIngredientMeasurementUnitRouter routing.Router) {
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientMeasurementUnitsPermission)).
				Post(root, s.validIngredientMeasurementUnitsService.CreateValidIngredientMeasurementUnitHandler)
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
				Get(root, s.validIngredientMeasurementUnitsService.ListValidIngredientMeasurementUnitsHandler)

			validIngredientMeasurementUnitsByIngredientRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientmeasurementunitsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByIngredientRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.SearchValidIngredientMeasurementUnitsByIngredientHandler)
			})

			validIngredientMeasurementUnitsByMeasurementUnitRouteParam := fmt.Sprintf("/by_measurement_unit%s", buildURLVarChunk(validingredientmeasurementunitsservice.ValidMeasurementUnitIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByMeasurementUnitRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.SearchValidIngredientMeasurementUnitsByMeasurementUnitHandler)
			})

			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitIDRouteParam, func(singleValidIngredientMeasurementUnitRouter routing.Router) {
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.ReadValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientMeasurementUnitsPermission)).
					Put(root, s.validIngredientMeasurementUnitsService.UpdateValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientMeasurementUnitsPermission)).
					Delete(root, s.validIngredientMeasurementUnitsService.ArchiveValidIngredientMeasurementUnitHandler)
			})
		})

		// Meals
		mealPath := "meals"
		mealsRouteWithPrefix := fmt.Sprintf("/%s", mealPath)
		mealIDRouteParam := buildURLVarChunk(mealsservice.MealIDURIParamKey, "")
		v1Router.Route(mealsRouteWithPrefix, func(mealsRouter routing.Router) {
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealsPermission)).
				Post(root, s.mealsService.CreateMealHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(root, s.mealsService.ListMealsHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(searchRoot, s.mealsService.SearchMealsHandler)

			mealsRouter.Route(mealIDRouteParam, func(singleMealRouter routing.Router) {
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
					Get(root, s.mealsService.ReadMealHandler)
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealsPermission)).
					Delete(root, s.mealsService.ArchiveMealHandler)
			})
		})

		// Components
		recipePath := "recipes"
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		recipeIDRouteParam := buildURLVarChunk(recipesservice.RecipeIDURIParamKey, "")
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, s.recipesService.CreateRecipeHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, s.recipesService.ListRecipesHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(searchRoot, s.recipesService.SearchRecipesHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, s.recipesService.ReadRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/mermaid", s.recipesService.RecipeMermaidHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/prep_steps", s.recipesService.RecipeEstimatedPrepStepsHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, s.recipesService.UpdateRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Post("/clone", s.recipesService.CloneRecipeHandler)

				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Post("/images", s.recipesService.RecipeImageUploadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, s.recipesService.ArchiveRecipeHandler)

				// RecipeRatings
				singleRecipeRouter.Route("/ratings", func(recipeRatingsRouter routing.Router) {
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
						Post(root, s.recipeRatingsService.CreateRecipeRatingHandler)
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeRatingsPermission)).
						Get(root, s.recipeRatingsService.ListRecipeRatingsHandler)

					recipeRatingIDRouteParam := buildURLVarChunk(reciperatingsservice.RecipeRatingIDURIParamKey, "")
					recipeRatingsRouter.Route(recipeRatingIDRouteParam, func(singleRecipeRatingRouter routing.Router) {
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
							Get(root, s.recipeRatingsService.ReadRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeRatingsPermission)).
							Put(root, s.recipeRatingsService.UpdateRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeRatingsPermission)).
							Delete(root, s.recipeRatingsService.ArchiveRecipeRatingHandler)
					})
				})
			})
		})

		// PrepTasks
		recipePrepTaskPath := "prep_tasks"
		recipePrepTasksRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipePrepTaskPath,
		)
		recipePrepTasksRouteWithPrefix := fmt.Sprintf("/%s", recipePrepTasksRoute)
		recipePrepTaskIDRouteParam := buildURLVarChunk(recipepreptasksservice.RecipePrepTaskIDURIParamKey, "")
		v1Router.Route(recipePrepTasksRouteWithPrefix, func(recipePrepTasksRouter routing.Router) {
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipePrepTasksPermission)).
				Post(root, s.recipePrepTasksService.CreateRecipePrepTaskHandler)
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
				Get(root, s.recipePrepTasksService.ListRecipePrepTaskHandler)

			recipePrepTasksRouter.Route(recipePrepTaskIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
					Get(root, s.recipePrepTasksService.ReadRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipePrepTasksPermission)).
					Put(root, s.recipePrepTasksService.UpdateRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipePrepTasksPermission)).
					Delete(root, s.recipePrepTasksService.ArchiveRecipePrepTaskHandler)
			})
		})

		// RecipeSteps
		recipeStepPath := "steps"
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
				Post(root, s.recipeStepsService.CreateRecipeStepHandler)
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
				Get(root, s.recipeStepsService.ListRecipeStepsHandler)

			recipeStepsRouter.Route(recipeStepIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Post("/images", s.recipeStepsService.RecipeStepImageUploadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Get(root, s.recipeStepsService.ReadRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, s.recipeStepsService.UpdateRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, s.recipeStepsService.ArchiveRecipeStepHandler)
			})
		})

		// RecipeStepInstruments
		recipeStepInstrumentPath := "instruments"
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
				Post(root, s.recipeStepInstrumentsService.CreateRecipeStepInstrumentHandler)
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
				Get(root, s.recipeStepInstrumentsService.ListRecipeStepInstrumentsHandler)

			recipeStepInstrumentsRouter.Route(recipeStepInstrumentIDRouteParam, func(singleRecipeStepInstrumentRouter routing.Router) {
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
					Get(root, s.recipeStepInstrumentsService.ReadRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepInstrumentsPermission)).
					Put(root, s.recipeStepInstrumentsService.UpdateRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepInstrumentsPermission)).
					Delete(root, s.recipeStepInstrumentsService.ArchiveRecipeStepInstrumentHandler)
			})
		})

		// RecipeStepVessels
		recipeStepVesselPath := "vessels"
		recipeStepVesselsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepVesselPath,
		)
		recipeStepVesselsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepVesselsRoute)
		recipeStepVesselIDRouteParam := buildURLVarChunk(recipestepvesselsservice.RecipeStepVesselIDURIParamKey, "")
		v1Router.Route(recipeStepVesselsRouteWithPrefix, func(recipeStepVesselsRouter routing.Router) {
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepVesselsPermission)).
				Post(root, s.recipeStepVesselsService.CreateRecipeStepVesselHandler)
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
				Get(root, s.recipeStepVesselsService.ListRecipeStepVesselsHandler)

			recipeStepVesselsRouter.Route(recipeStepVesselIDRouteParam, func(singleRecipeStepVesselRouter routing.Router) {
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
					Get(root, s.recipeStepVesselsService.ReadRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepVesselsPermission)).
					Put(root, s.recipeStepVesselsService.UpdateRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepVesselsPermission)).
					Delete(root, s.recipeStepVesselsService.ArchiveRecipeStepVesselHandler)
			})
		})

		// RecipeStepIngredients
		recipeStepIngredientPath := "ingredients"
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
				Post(root, s.recipeStepIngredientsService.CreateRecipeStepIngredientHandler)
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
				Get(root, s.recipeStepIngredientsService.ListRecipeStepIngredientsHandler)

			recipeStepIngredientsRouter.Route(recipeStepIngredientIDRouteParam, func(singleRecipeStepIngredientRouter routing.Router) {
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Get(root, s.recipeStepIngredientsService.ReadRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, s.recipeStepIngredientsService.UpdateRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, s.recipeStepIngredientsService.ArchiveRecipeStepIngredientHandler)
			})
		})

		// RecipeStepIngredients
		recipeStepCompletionConditionPath := "completion_conditions"
		recipeStepCompletionConditionsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepCompletionConditionPath,
		)
		recipeStepCompletionConditionsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepCompletionConditionsRoute)
		recipeStepCompletionConditionIDRouteParam := buildURLVarChunk(recipestepcompletionconditionsservice.RecipeStepCompletionConditionIDURIParamKey, "")
		v1Router.Route(recipeStepCompletionConditionsRouteWithPrefix, func(recipeStepCompletionConditionsRouter routing.Router) {
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepCompletionConditionsPermission)).
				Post(root, s.recipeStepCompletionConditionsService.CreateRecipeStepCompletionConditionHandler)
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
				Get(root, s.recipeStepCompletionConditionsService.ListRecipeStepCompletionConditionsHandler)

			recipeStepCompletionConditionsRouter.Route(recipeStepCompletionConditionIDRouteParam, func(singleRecipeStepCompletionConditionRouter routing.Router) {
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
					Get(root, s.recipeStepCompletionConditionsService.ReadRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepCompletionConditionsPermission)).
					Put(root, s.recipeStepCompletionConditionsService.UpdateRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepCompletionConditionsPermission)).
					Delete(root, s.recipeStepCompletionConditionsService.ArchiveRecipeStepCompletionConditionHandler)
			})
		})

		// RecipeStepProducts
		recipeStepProductPath := "products"
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
				Post(root, s.recipeStepProductsService.CreateRecipeStepProductHandler)
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
				Get(root, s.recipeStepProductsService.ListRecipeStepProductsHandler)

			recipeStepProductsRouter.Route(recipeStepProductIDRouteParam, func(singleRecipeStepProductRouter routing.Router) {
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Get(root, s.recipeStepProductsService.ReadRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, s.recipeStepProductsService.UpdateRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, s.recipeStepProductsService.ArchiveRecipeStepProductHandler)
			})
		})

		// MealPlans
		mealPlanPath := "meal_plans"
		mealPlansRouteWithPrefix := fmt.Sprintf("/%s", mealPlanPath)
		mealPlanIDRouteParam := buildURLVarChunk(mealplansservice.MealPlanIDURIParamKey, "")
		v1Router.Route(mealPlansRouteWithPrefix, func(mealPlansRouter routing.Router) {
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
				Post(root, s.mealPlansService.CreateMealPlanHandler)
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
				Get(root, s.mealPlansService.ListMealPlanHandler)

			mealPlansRouter.Route(mealPlanIDRouteParam, func(singleMealPlanRouter routing.Router) {
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
					Get(root, s.mealPlansService.ReadMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlansPermission)).
					Put(root, s.mealPlansService.UpdateMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlansPermission)).
					Delete(root, s.mealPlansService.ArchiveMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
					Post("/finalize", s.mealPlansService.FinalizeMealPlanHandler)
			})
		})

		// MealPlanTasks
		mealPlanTaskPath := "tasks"
		mealPlanTasksRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanTaskPath,
		)
		mealPlanTasksRouteWithPrefix := fmt.Sprintf("/%s", mealPlanTasksRoute)
		mealPlanTaskIDRouteParam := buildURLVarChunk(mealplantasksservice.MealPlanTaskIDURIParamKey, "")
		v1Router.Route(mealPlanTasksRouteWithPrefix, func(mealPlanTasksRouter routing.Router) {
			mealPlanTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
				Get(root, s.mealPlanTasksService.ListMealPlanTasksByMealPlanHandler)
			mealPlanTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanTasksPermission)).
				Post(root, s.mealPlanTasksService.CreateMealPlanTaskHandler)

			mealPlanTasksRouter.Route(mealPlanTaskIDRouteParam, func(singleMealPlanTaskRouter routing.Router) {
				singleMealPlanTaskRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
					Get(root, s.mealPlanTasksService.ReadMealPlanTaskHandler)

				singleMealPlanTaskRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanTasksPermission)).
					Patch(root, s.mealPlanTasksService.MealPlanTaskStatusChangeHandler)
			})
		})

		// MealPlanEvents
		mealPlanEventPath := "events"
		mealPlanEventsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
		)
		mealPlanEventsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanEventsRoute)
		mealPlanEventIDRouteParam := buildURLVarChunk(mealplaneventsservice.MealPlanEventIDURIParamKey, "")
		v1Router.Route(mealPlanEventsRouteWithPrefix, func(mealPlanEventsRouter routing.Router) {
			mealPlanEventsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanEventsPermission)).
				Post(root, s.mealPlanEventsService.CreateMealPlanEventHandler)
			mealPlanEventsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
				Get(root, s.mealPlanEventsService.ListMealPlanEventHandler)

			mealPlanEventsRouter.Route(mealPlanEventIDRouteParam, func(singleMealPlanEventRouter routing.Router) {
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
					Get(root, s.mealPlanEventsService.ReadMealPlanEventHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanEventsPermission)).
					Put(root, s.mealPlanEventsService.UpdateMealPlanEventHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionVotesPermission)).
					Post("/vote", s.mealPlanOptionVotesService.CreateMealPlanOptionVoteHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanEventsPermission)).
					Delete(root, s.mealPlanEventsService.ArchiveMealPlanEventHandler)
			})
		})

		// MealPlanGroceryListItems
		mealPlanGroceryListItemPath := "grocery_list_items"
		mealPlanGroceryListItemsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanGroceryListItemPath,
		)
		mealPlanGroceryListItemsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanGroceryListItemsRoute)
		mealPlanGroceryListItemIDRouteParam := buildURLVarChunk(mealplangrocerylistitemssservice.MealPlanGroceryListItemIDURIParamKey, "")
		v1Router.Route(mealPlanGroceryListItemsRouteWithPrefix, func(mealPlanGroceryListItemsRouter routing.Router) {
			mealPlanGroceryListItemsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanGroceryListItemsPermission)).
				Post(root, s.mealPlanGroceryListItemsService.CreateMealPlanGroceryListItemHandler)
			mealPlanGroceryListItemsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
				Get(root, s.mealPlanGroceryListItemsService.ListMealPlanGroceryListItemsByMealPlanHandler)

			mealPlanGroceryListItemsRouter.Route(mealPlanGroceryListItemIDRouteParam, func(singleMealPlanGroceryListItemRouter routing.Router) {
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
					Get(root, s.mealPlanGroceryListItemsService.ReadMealPlanGroceryListItemHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanGroceryListItemsPermission)).
					Put(root, s.mealPlanGroceryListItemsService.UpdateMealPlanGroceryListItemHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanGroceryListItemsPermission)).
					Delete(root, s.mealPlanGroceryListItemsService.ArchiveMealPlanGroceryListItemHandler)
			})
		})

		// MealPlanOptions
		mealPlanOptionPath := "options"
		mealPlanOptionsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
			mealPlanEventIDRouteParam,
			mealPlanOptionPath,
		)
		mealPlanOptionsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionsRoute)
		mealPlanOptionIDRouteParam := buildURLVarChunk(mealplanoptionsservice.MealPlanOptionIDURIParamKey, "")
		v1Router.Route(mealPlanOptionsRouteWithPrefix, func(mealPlanOptionsRouter routing.Router) {
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionsPermission)).
				Post(root, s.mealPlanOptionsService.CreateMealPlanOptionHandler)
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
				Get(root, s.mealPlanOptionsService.ListMealPlanOptionHandler)

			mealPlanOptionsRouter.Route(mealPlanOptionIDRouteParam, func(singleMealPlanOptionRouter routing.Router) {
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
					Get(root, s.mealPlanOptionsService.ReadMealPlanOptionHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionsPermission)).
					Put(root, s.mealPlanOptionsService.UpdateMealPlanOptionHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionsPermission)).
					Delete(root, s.mealPlanOptionsService.ArchiveMealPlanOptionHandler)
			})
		})

		// MealPlanOptionVotes
		mealPlanOptionVotePath := "votes"
		mealPlanOptionVotesRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
			mealPlanEventIDRouteParam,
			mealPlanOptionPath,
			mealPlanOptionIDRouteParam,
			mealPlanOptionVotePath,
		)
		mealPlanOptionVotesRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionVotesRoute)
		mealPlanOptionVoteIDRouteParam := buildURLVarChunk(mealplanoptionvotesservice.MealPlanOptionVoteIDURIParamKey, "")
		v1Router.Route(mealPlanOptionVotesRouteWithPrefix, func(mealPlanOptionVotesRouter routing.Router) {
			mealPlanOptionVotesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
				Get(root, s.mealPlanOptionVotesService.ListMealPlanOptionVoteHandler)

			mealPlanOptionVotesRouter.Route(mealPlanOptionVoteIDRouteParam, func(singleMealPlanOptionVoteRouter routing.Router) {
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
					Get(root, s.mealPlanOptionVotesService.ReadMealPlanOptionVoteHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionVotesPermission)).
					Put(root, s.mealPlanOptionVotesService.UpdateMealPlanOptionVoteHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionVotesPermission)).
					Delete(root, s.mealPlanOptionVotesService.ArchiveMealPlanOptionVoteHandler)
			})
		})

		// ServiceSettings
		serviceSettingPath := "settings"
		serviceSettingsRouteWithPrefix := fmt.Sprintf("/%s", serviceSettingPath)
		serviceSettingIDRouteParam := buildURLVarChunk(servicesettingsservice.ServiceSettingIDURIParamKey, "")
		v1Router.Route(serviceSettingsRouteWithPrefix, func(serviceSettingsRouter routing.Router) {
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateServiceSettingsPermission)).
				Post(root, s.serviceSettingsService.CreateServiceSettingHandler)
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(root, s.serviceSettingsService.ListServiceSettingsHandler)
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(searchRoot, s.serviceSettingsService.SearchServiceSettingsHandler)

			serviceSettingsRouter.Route(serviceSettingIDRouteParam, func(singleServiceSettingRouter routing.Router) {
				singleServiceSettingRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
					Get(root, s.serviceSettingsService.ReadServiceSettingHandler)
				singleServiceSettingRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingsPermission)).
					Delete(root, s.serviceSettingsService.ArchiveServiceSettingHandler)
			})

			serviceSettingConfigurationIDRouteParam := buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationIDURIParamKey, "")
			serviceSettingsRouter.Route("/configurations", func(settingConfigurationRouter routing.Router) {
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateServiceSettingConfigurationsPermission)).
					Post(root, s.serviceSettingConfigurationsService.CreateServiceSettingConfigurationHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user"+buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationNameURIParamKey, ""), s.serviceSettingConfigurationsService.GetServiceSettingConfigurationsForUserByNameHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user", s.serviceSettingConfigurationsService.GetServiceSettingConfigurationsForUserHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/household", s.serviceSettingConfigurationsService.GetServiceSettingConfigurationsForHouseholdHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateServiceSettingConfigurationsPermission)).
					Put(serviceSettingConfigurationIDRouteParam, s.serviceSettingConfigurationsService.UpdateServiceSettingConfigurationHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingConfigurationsPermission)).
					Delete(serviceSettingConfigurationIDRouteParam, s.serviceSettingConfigurationsService.ArchiveServiceSettingConfigurationHandler)
			})
		})

		// User Notifications
		userNotificationPath := "user_notifications"
		userNotificationsRouteWithPrefix := fmt.Sprintf("/%s", userNotificationPath)
		userNotificationIDRouteParam := buildURLVarChunk(usernotificationsservice.UserNotificationIDURIParamKey, "")
		v1Router.Route(userNotificationsRouteWithPrefix, func(userNotificationsRouter routing.Router) {
			userNotificationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateUserNotificationsPermission)).
				Post(root, s.userNotificationsService.CreateUserNotificationHandler)
			userNotificationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
				Get(root, s.userNotificationsService.ListUserNotificationsHandler)

			userNotificationsRouter.Route(userNotificationIDRouteParam, func(singleUserNotificationRouter routing.Router) {
				singleUserNotificationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
					Get(root, s.userNotificationsService.ReadUserNotificationHandler)
				singleUserNotificationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserNotificationsPermission)).
					Patch(root, s.userNotificationsService.UpdateUserNotificationHandler)
			})
		})
	})

	s.router = router
}
