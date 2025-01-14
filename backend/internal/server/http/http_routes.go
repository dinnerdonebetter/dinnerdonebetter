package http

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/routing"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/core/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/dataprivacy"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/core/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/core/households"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/core/oauth2clients"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/core/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/core/servicesettings"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/core/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/users"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/core/webhooks"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/eating/householdinstrumentownerships"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanevents"
	mealplangrocerylistitemssservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplans"
	mealplantasksservice "github.com/dinnerdonebetter/backend/internal/services/eating/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/eating/meals"
	recipemanagementservice "github.com/dinnerdonebetter/backend/internal/services/eating/recipe_management"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/eating/useringredientpreferences"
	validenumerationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/valid_enumerations"
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
				logger.Error("database not responding to ping", err)
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
		validInstrumentIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidInstrumentIDURIParamKey, "")
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter routing.Router) {
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
				Post(root, s.validEnumerationsService.CreateValidInstrumentHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(root, s.validEnumerationsService.ListValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidInstrumentsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(randomRoot, s.validEnumerationsService.RandomValidInstrumentHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, s.validEnumerationsService.ReadValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidInstrumentHandler)
			})
		})

		// ValidVessels
		validVesselPath := "valid_vessels"
		validVesselsRouteWithPrefix := fmt.Sprintf("/%s", validVesselPath)
		validVesselIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidVesselIDURIParamKey, "")
		v1Router.Route(validVesselsRouteWithPrefix, func(validVesselsRouter routing.Router) {
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidVesselsPermission)).
				Post(root, s.validEnumerationsService.CreateValidVesselHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(root, s.validEnumerationsService.ListValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidVesselsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(randomRoot, s.validEnumerationsService.RandomValidVesselHandler)

			validVesselsRouter.Route(validVesselIDRouteParam, func(singleValidVesselRouter routing.Router) {
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
					Get(root, s.validEnumerationsService.ReadValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidVesselsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidVesselsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidVesselHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter routing.Router) {
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission, authorization.SearchValidIngredientsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(randomRoot, s.validEnumerationsService.RandomValidIngredientHandler)

			validIngredientsByPreparationIDSearchRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientsRouter.Route(validIngredientsByPreparationIDSearchRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientsByPreparationAndIngredientNameHandler)
			})

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientHandler)
			})
		})

		// ValidIngredientGroups
		validIngredientGroupPath := "valid_ingredient_groups"
		validIngredientGroupsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientGroupPath)
		validIngredientGroupIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientGroupIDURIParamKey, "")
		v1Router.Route(validIngredientGroupsRouteWithPrefix, func(validIngredientGroupsRouter routing.Router) {
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientGroupsPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientGroupHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientGroupsHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidIngredientGroupsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidIngredientGroupsHandler)

			validIngredientGroupsRouter.Route(validIngredientGroupIDRouteParam, func(singleValidIngredientGroupRouter routing.Router) {
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientGroupsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientGroupsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientGroupHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey, "")
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter routing.Router) {
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
				Post(root, s.validEnumerationsService.CreateValidPreparationHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(root, s.validEnumerationsService.ListValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(randomRoot, s.validEnumerationsService.RandomValidPreparationHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, s.validEnumerationsService.ReadValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidPreparationHandler)
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
		validMeasurementUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitServiceIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitsRouteWithPrefix, func(validMeasurementUnitsRouter routing.Router) {
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitsPermission)).
				Post(root, s.validEnumerationsService.CreateValidMeasurementUnitHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
				Get(root, s.validEnumerationsService.ListValidMeasurementUnitsHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidMeasurementUnitsHandler)

			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(path.Join("/by_ingredient", validMeasurementUnitServiceIngredientIDRouteParam), s.validEnumerationsService.SearchValidMeasurementUnitsByIngredientIDHandler)

			validMeasurementUnitsRouter.Route(validMeasurementUnitIDRouteParam, func(singleValidMeasurementUnitRouter routing.Router) {
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
					Get(root, s.validEnumerationsService.ReadValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidMeasurementUnitHandler)
			})
		})

		// ValidIngredientStates
		validIngredientStatePath := "valid_ingredient_states"
		validIngredientStatesRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStatePath)
		validIngredientStateIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientStateIDURIParamKey, "")
		v1Router.Route(validIngredientStatesRouteWithPrefix, func(validIngredientStatesRouter routing.Router) {
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStatesPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientStateHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientStatesHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(searchRoot, s.validEnumerationsService.SearchValidIngredientStatesHandler)

			validIngredientStatesRouter.Route(validIngredientStateIDRouteParam, func(singleValidIngredientStateRouter routing.Router) {
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStatesPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStatesPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientStateHandler)
			})
		})

		// ValidMeasurementUnitConversions
		validMeasurementUnitConversionPath := "valid_measurement_conversions"
		validMeasurementUnitConversionsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitConversionPath)
		validMeasurementUnitConversionUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitConversionIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitConversionIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitConversionsRouteWithPrefix, func(validMeasurementUnitConversionsRouter routing.Router) {
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitConversionsPermission)).
				Post(root, s.validEnumerationsService.CreateValidMeasurementUnitConversionHandler)

			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/from_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validEnumerationsService.ValidMeasurementUnitConversionsFromMeasurementUnitHandler)
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/to_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validEnumerationsService.ValidMeasurementUnitConversionsToMeasurementUnitHandler)

			validMeasurementUnitConversionsRouter.Route(validMeasurementUnitConversionIDRouteParam, func(singleValidMeasurementUnitConversionRouter routing.Router) {
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
					Get(root, s.validEnumerationsService.ReadValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitConversionsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitConversionsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidMeasurementUnitConversionHandler)
			})
		})

		// ValidIngredientStateIngredients
		validIngredientStateIngredientPath := "valid_ingredient_state_ingredients"
		validIngredientStateIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStateIngredientPath)
		validIngredientStateIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientStateIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientStateIngredientsRouteWithPrefix, func(validIngredientStateIngredientsRouter routing.Router) {
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStateIngredientsPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientStateIngredientHandler)
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientStateIngredientsHandler)

			validIngredientStateIngredientsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientStateIngredientsByIngredientHandler)
			})

			validIngredientStateIngredientsByStateIngredientIDRouteParam := fmt.Sprintf("/by_ingredient_state%s", buildURLVarChunk(validenumerationsservice.ValidIngredientStateIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByStateIngredientIDRouteParam, func(byValidStateIngredientIDRouter routing.Router) {
				byValidStateIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientStateIngredientsByIngredientStateHandler)
			})

			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientIDRouteParam, func(singleValidIngredientStateIngredientRouter routing.Router) {
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStateIngredientsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStateIngredientsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientStateIngredientHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientPreparationIDURIParamKey, "")
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter routing.Router) {
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientPreparationHandler)
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientPreparationsHandler)

			validIngredientPreparationsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientPreparationsByIngredientHandler)
			})

			validIngredientPreparationsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientPreparationsByPreparationHandler)
			})

			validIngredientPreparationsRouter.Route(validIngredientPreparationIDRouteParam, func(singleValidIngredientPreparationRouter routing.Router) {
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientPreparationHandler)
			})
		})

		// ValidPreparationInstruments
		validPreparationInstrumentPath := "valid_preparation_instruments"
		validPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationInstrumentPath)
		v1Router.Route(validPreparationInstrumentsRouteWithPrefix, func(validPreparationInstrumentsRouter routing.Router) {
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
				Post(root, s.validEnumerationsService.CreateValidPreparationInstrumentHandler)
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
				Get(root, s.validEnumerationsService.ListValidPreparationInstrumentsHandler)

			validPreparationInstrumentsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validEnumerationsService.SearchValidPreparationInstrumentsByPreparationHandler)
			})

			validPreparationInstrumentsByInstrumentIDRouteParam := fmt.Sprintf("/by_instrument%s", buildURLVarChunk(validenumerationsservice.ValidInstrumentIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validEnumerationsService.SearchValidPreparationInstrumentsByInstrumentHandler)
			})

			validPreparationInstrumentIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationInstrumentIDURIParamKey, "")
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentIDRouteParam, func(singleValidPreparationInstrumentRouter routing.Router) {
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validEnumerationsService.ReadValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidPreparationInstrumentHandler)
			})
		})

		// ValidPreparationVessels
		validPreparationVesselPath := "valid_preparation_vessels"
		validPreparationVesselsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationVesselPath)
		v1Router.Route(validPreparationVesselsRouteWithPrefix, func(validPreparationVesselsRouter routing.Router) {
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationVesselsPermission)).
				Post(root, s.validEnumerationsService.CreateValidPreparationVesselHandler)
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
				Get(root, s.validEnumerationsService.ListValidPreparationVesselsHandler)

			validPreparationVesselsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validEnumerationsService.SearchValidPreparationVesselsByPreparationHandler)
			})

			validPreparationVesselsByInstrumentIDRouteParam := fmt.Sprintf("/by_vessel%s", buildURLVarChunk(validenumerationsservice.ValidVesselIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validEnumerationsService.SearchValidPreparationVesselsByVesselHandler)
			})

			validPreparationVesselIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationVesselIDURIParamKey, "")
			validPreparationVesselsRouter.Route(validPreparationVesselIDRouteParam, func(singleValidPreparationVesselRouter routing.Router) {
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validEnumerationsService.ReadValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationVesselsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationVesselsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidPreparationVesselHandler)
			})
		})

		// ValidIngredientMeasurementUnit
		validIngredientMeasurementUnitPath := "valid_ingredient_measurement_units"
		validIngredientMeasurementUnitRouteWithPrefix := fmt.Sprintf("/%s", validIngredientMeasurementUnitPath)
		validIngredientMeasurementUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientMeasurementUnitIDURIParamKey, "")
		v1Router.Route(validIngredientMeasurementUnitRouteWithPrefix, func(validIngredientMeasurementUnitRouter routing.Router) {
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientMeasurementUnitsPermission)).
				Post(root, s.validEnumerationsService.CreateValidIngredientMeasurementUnitHandler)
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
				Get(root, s.validEnumerationsService.ListValidIngredientMeasurementUnitsHandler)

			validIngredientMeasurementUnitsByIngredientRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByIngredientRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientMeasurementUnitsByIngredientHandler)
			})

			validIngredientMeasurementUnitsByMeasurementUnitRouteParam := fmt.Sprintf("/by_measurement_unit%s", buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByMeasurementUnitRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validEnumerationsService.SearchValidIngredientMeasurementUnitsByMeasurementUnitHandler)
			})

			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitIDRouteParam, func(singleValidIngredientMeasurementUnitRouter routing.Router) {
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validEnumerationsService.ReadValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientMeasurementUnitsPermission)).
					Put(root, s.validEnumerationsService.UpdateValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientMeasurementUnitsPermission)).
					Delete(root, s.validEnumerationsService.ArchiveValidIngredientMeasurementUnitHandler)
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
		recipeIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeIDURIParamKey, "")
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, s.recipeManagementService.CreateRecipeHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, s.recipeManagementService.ListRecipesHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(searchRoot, s.recipeManagementService.SearchRecipesHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, s.recipeManagementService.ReadRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/mermaid", s.recipeManagementService.RecipeMermaidHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/prep_steps", s.recipeManagementService.RecipeEstimatedPrepStepsHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Post("/clone", s.recipeManagementService.CloneRecipeHandler)

				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Post("/images", s.recipeManagementService.RecipeImageUploadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeHandler)

				// RecipeRatings
				singleRecipeRouter.Route("/ratings", func(recipeRatingsRouter routing.Router) {
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
						Post(root, s.recipeManagementService.CreateRecipeRatingHandler)
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeRatingsPermission)).
						Get(root, s.recipeManagementService.ListRecipeRatingsHandler)

					recipeRatingIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeRatingIDURIParamKey, "")
					recipeRatingsRouter.Route(recipeRatingIDRouteParam, func(singleRecipeRatingRouter routing.Router) {
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
							Get(root, s.recipeManagementService.ReadRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeRatingsPermission)).
							Put(root, s.recipeManagementService.UpdateRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeRatingsPermission)).
							Delete(root, s.recipeManagementService.ArchiveRecipeRatingHandler)
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
		recipePrepTaskIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipePrepTaskIDURIParamKey, "")
		v1Router.Route(recipePrepTasksRouteWithPrefix, func(recipePrepTasksRouter routing.Router) {
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipePrepTasksPermission)).
				Post(root, s.recipeManagementService.CreateRecipePrepTaskHandler)
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
				Get(root, s.recipeManagementService.ListRecipePrepTaskHandler)

			recipePrepTasksRouter.Route(recipePrepTaskIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
					Get(root, s.recipeManagementService.ReadRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipePrepTasksPermission)).
					Put(root, s.recipeManagementService.UpdateRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipePrepTasksPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipePrepTaskHandler)
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
		recipeStepIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepIDURIParamKey, "")
		v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter routing.Router) {
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepHandler)
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepsHandler)

			recipeStepsRouter.Route(recipeStepIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Post("/images", s.recipeManagementService.RecipeStepImageUploadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepHandler)
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
		recipeStepInstrumentIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepInstrumentIDURIParamKey, "")
		v1Router.Route(recipeStepInstrumentsRouteWithPrefix, func(recipeStepInstrumentsRouter routing.Router) {
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepInstrumentsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepInstrumentHandler)
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepInstrumentsHandler)

			recipeStepInstrumentsRouter.Route(recipeStepInstrumentIDRouteParam, func(singleRecipeStepInstrumentRouter routing.Router) {
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepInstrumentsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepInstrumentsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepInstrumentHandler)
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
		recipeStepVesselIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepVesselIDURIParamKey, "")
		v1Router.Route(recipeStepVesselsRouteWithPrefix, func(recipeStepVesselsRouter routing.Router) {
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepVesselsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepVesselHandler)
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepVesselsHandler)

			recipeStepVesselsRouter.Route(recipeStepVesselIDRouteParam, func(singleRecipeStepVesselRouter routing.Router) {
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepVesselsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepVesselsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepVesselHandler)
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
		recipeStepIngredientIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepIngredientIDURIParamKey, "")
		v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter routing.Router) {
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepIngredientsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepIngredientHandler)
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepIngredientsHandler)

			recipeStepIngredientsRouter.Route(recipeStepIngredientIDRouteParam, func(singleRecipeStepIngredientRouter routing.Router) {
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepIngredientHandler)
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
		recipeStepCompletionConditionIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepCompletionConditionIDURIParamKey, "")
		v1Router.Route(recipeStepCompletionConditionsRouteWithPrefix, func(recipeStepCompletionConditionsRouter routing.Router) {
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepCompletionConditionsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepCompletionConditionHandler)
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepCompletionConditionsHandler)

			recipeStepCompletionConditionsRouter.Route(recipeStepCompletionConditionIDRouteParam, func(singleRecipeStepCompletionConditionRouter routing.Router) {
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepCompletionConditionsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepCompletionConditionsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepCompletionConditionHandler)
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
		recipeStepProductIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepProductIDURIParamKey, "")
		v1Router.Route(recipeStepProductsRouteWithPrefix, func(recipeStepProductsRouter routing.Router) {
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepProductsPermission)).
				Post(root, s.recipeManagementService.CreateRecipeStepProductHandler)
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
				Get(root, s.recipeManagementService.ListRecipeStepProductsHandler)

			recipeStepProductsRouter.Route(recipeStepProductIDRouteParam, func(singleRecipeStepProductRouter routing.Router) {
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Get(root, s.recipeManagementService.ReadRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, s.recipeManagementService.UpdateRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, s.recipeManagementService.ArchiveRecipeStepProductHandler)
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
