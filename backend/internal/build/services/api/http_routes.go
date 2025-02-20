package api

import (
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	dataprivacyservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/dataprivacy"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/households"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/oauth2clients"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/servicesettings"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/users"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/webhooks"
	mealplanningservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/meal_planning"
	recipemanagementservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/recipe_management"
	validenumerationsservice "github.com/dinnerdonebetter/backend/internal/services/eating/handlers/valid_enumerations"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	root       = "/"
	randomRoot = "/random"
	searchRoot = "/search"
)

func buildURLVarChunk(key string) string {
	return fmt.Sprintf("/{%s}", key)
}

func ProvideAPIRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	dataManager database.DataManager,
	authService types.AuthDataService,
	householdsService types.HouseholdDataService,
	householdInvitationsService types.HouseholdInvitationDataService,
	usersService types.UserDataService,
	adminService types.AdminDataService,
	webhooksService types.WebhookDataService,
	serviceSettingsService types.ServiceSettingDataService,
	serviceSettingConfigurationsService types.ServiceSettingConfigurationDataService,
	oauth2ClientsService types.OAuth2ClientDataService,
	userNotificationsService types.UserNotificationDataService,
	workerService types.WorkerService,
	validEnumerationsService types.ValidEnumerationDataService,
	auditLogEntriesService types.AuditLogEntryDataService,
	dataPrivacyService types.DataPrivacyService,
	recipeManagementService types.RecipeManagementDataService,
	mealPlanningService types.MealPlanningDataService,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	tracer := tracing.NewTracer(tracerProvider.Tracer("api_router"))

	router.Route("/_ops_", func(metaRouter routing.Router) {
		// Expose a readiness check on /ready
		metaRouter.Get("/live", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			reqCtx, reqSpan := tracer.StartSpan(req.Context())
			defer reqSpan.End()

			responseCode := http.StatusOK
			if err = dataManager.DB().PingContext(reqCtx); err != nil {
				logger.WithRequest(req).Error("database not responding to ping", err)
				responseCode = http.StatusInternalServerError
			}

			res.WriteHeader(responseCode)
		})
	})

	authenticatedRouter := router.WithMiddleware(authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", authService.StatusHandler)

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.Get("/authorize", authService.AuthorizeHandler)
		userRouter.Post("/token", authService.TokenHandler)
	})

	// these are routes we don't expect or require users to be authenticated for
	router.Route("/users", func(userRouter routing.Router) {
		userRouter.Post(root, usersService.CreateUserHandler)
		userRouter.Post("/login/jwt", authService.BuildLoginHandler(false))
		userRouter.Post("/login/jwt/admin", authService.BuildLoginHandler(true))
		userRouter.Post("/username/reminder", usersService.RequestUsernameReminderHandler)
		userRouter.Post("/password/reset", usersService.CreatePasswordResetTokenHandler)
		userRouter.Post("/password/reset/redeem", usersService.PasswordResetTokenRedemptionHandler)
		userRouter.Post("/email_address/verify", usersService.VerifyUserEmailAddressHandler)
		userRouter.Post("/totp_secret/verify", usersService.TOTPSecretVerificationHandler)
	})

	router.Route("/auth", func(authRouter routing.Router) {
		providerRouteParam := buildURLVarChunk(authservice.AuthProviderParamKey)
		authRouter.Get(providerRouteParam, authService.SSOLoginHandler)
		authRouter.Get(path.Join(providerRouteParam, "callback"), authService.SSOLoginCallbackHandler)
	})

	authenticatedRouter.WithMiddleware(authService.AuthorizationMiddleware).Route("/api/v1", func(v1Router routing.Router) {
		adminRouter := v1Router.WithMiddleware(authService.ServiceAdminMiddleware)

		// Admin
		adminRouter.Route("/admin", func(adminRouter routing.Router) {
			adminRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/users/status", adminService.UserAccountStatusChangeHandler)
			adminRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/queues/test", adminService.UserAccountStatusChangeHandler)
		})

		// Workers
		adminRouter.Route("/workers", func(adminRouter routing.Router) {
			adminRouter.
				Post("/finalize_meal_plans", workerService.MealPlanFinalizationHandler)
			adminRouter.
				Post("/meal_plan_grocery_list_init", workerService.MealPlanGroceryListInitializationHandler)
			adminRouter.
				Post("/meal_plan_tasks", workerService.MealPlanTaskCreationHandler)
		})

		// Users
		v1Router.Route("/users", func(usersRouter routing.Router) {
			usersRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
				Get(root, usersService.ListUsersHandler)
			usersRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
				Get(searchRoot, usersService.UsernameSearchHandler)
			usersRouter.Post("/avatar/upload", usersService.AvatarUploadHandler)

			usersRouter.Get("/self", usersService.SelfHandler)
			usersRouter.Post("/email_address_verification", usersService.RequestEmailVerificationEmailHandler)
			usersRouter.Post("/permissions/check", usersService.UserPermissionsHandler)
			usersRouter.Put("/password/new", usersService.UpdatePasswordHandler)
			usersRouter.Post("/totp_secret/new", usersService.NewTOTPSecretHandler)
			usersRouter.Put("/username", usersService.UpdateUserUsernameHandler)
			usersRouter.Put("/email_address", usersService.UpdateUserEmailAddressHandler)
			usersRouter.Put("/details", usersService.UpdateUserDetailsHandler)

			singleUserRoute := buildURLVarChunk(usersservice.UserIDURIParamKey)
			usersRouter.Route(singleUserRoute, func(singleUserRouter routing.Router) {
				singleUserRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
					Get(root, usersService.ReadUserHandler)

				singleUserRouter.Delete(root, usersService.ArchiveUserHandler)
			})
		})

		v1Router.Route("/data_privacy", func(dataPrivacyRouter routing.Router) {
			dataPrivacyRouter.Delete("/destroy", dataPrivacyService.DataDeletionHandler)
			dataPrivacyRouter.Post("/disclose", dataPrivacyService.UserDataAggregationRequestHandler)

			singleReportRoute := buildURLVarChunk(dataprivacyservice.ReportIDURIParamKey)
			dataPrivacyRouter.Route("/reports", func(singleReportRouter routing.Router) {
				singleReportRouter.Get(singleReportRoute, dataPrivacyService.ReadUserDataAggregationReportHandler)
			})
		})

		// Households
		v1Router.Route("/households", func(householdsRouter routing.Router) {
			householdsRouter.Post(root, householdsService.CreateHouseholdHandler)
			householdsRouter.Get(root, householdsService.ListHouseholdsHandler)
			householdsRouter.Get("/current", householdsService.CurrentInfoHandler)

			singleUserRoute := buildURLVarChunk(householdsservice.UserIDURIParamKey)
			singleHouseholdRoute := buildURLVarChunk(householdsservice.HouseholdIDURIParamKey)
			householdsRouter.Route(singleHouseholdRoute, func(singleHouseholdRouter routing.Router) {
				singleHouseholdRouter.Get(root, householdsService.ReadHouseholdHandler)
				singleHouseholdRouter.Put(root, householdsService.UpdateHouseholdHandler)
				singleHouseholdRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdPermission)).
					Delete(root, householdsService.ArchiveHouseholdHandler)

				singleHouseholdRouter.Post("/default", householdsService.MarkAsDefaultHouseholdHandler)
				singleHouseholdRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.RemoveMemberHouseholdPermission)).
					Delete("/members"+singleUserRoute, householdsService.RemoveMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)).
					Post("/invite", householdInvitationsService.InviteMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ModifyMemberPermissionsForHouseholdPermission)).
					Patch("/members"+singleUserRoute+"/permissions", householdsService.ModifyMemberPermissionsHandler)
				singleHouseholdRouter.Post("/transfer", householdsService.TransferHouseholdOwnershipHandler)

				singleHouseholdRouter.Route("/invitations", func(invitationsRouter routing.Router) {
					invitationsRouter.Post(root, householdInvitationsService.InviteMemberHandler)

					singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey)
					invitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
						singleHouseholdInvitationRouter.Get(root, householdInvitationsService.ReadHouseholdInviteHandler)
					})
				})
			})

			// InstrumentOwnerships
			householdInstrumentOwnershipsRouteWithPrefix := "/instruments"
			householdInstrumentOwnershipIDRouteParam := buildURLVarChunk(mealplanningservice.HouseholdInstrumentOwnershipIDURIParamKey)
			householdsRouter.Route(householdInstrumentOwnershipsRouteWithPrefix, func(householdInstrumentOwnershipsRouter routing.Router) {
				householdInstrumentOwnershipsRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateHouseholdInstrumentOwnershipsPermission)).
					Post(root, mealPlanningService.CreateHouseholdInstrumentOwnershipHandler)
				householdInstrumentOwnershipsRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
					Get(root, mealPlanningService.ListHouseholdInstrumentOwnershipHandler)

				householdInstrumentOwnershipsRouter.Route(householdInstrumentOwnershipIDRouteParam, func(singleHouseholdInstrumentOwnershipRouter routing.Router) {
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
						Get(root, mealPlanningService.ReadHouseholdInstrumentOwnershipHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateHouseholdInstrumentOwnershipsPermission)).
						Put(root, mealPlanningService.UpdateHouseholdInstrumentOwnershipHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdInstrumentOwnershipsPermission)).
						Delete(root, mealPlanningService.ArchiveHouseholdInstrumentOwnershipHandler)
				})
			})
		})

		v1Router.Route("/household_invitations", func(householdInvitationsRouter routing.Router) {
			householdInvitationsRouter.Get("/sent", householdInvitationsService.OutboundInvitesHandler)
			householdInvitationsRouter.Get("/received", householdInvitationsService.InboundInvitesHandler)

			singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey)
			householdInvitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
				singleHouseholdInvitationRouter.Get(root, householdInvitationsService.ReadHouseholdInviteHandler)
				singleHouseholdInvitationRouter.Put("/cancel", householdInvitationsService.CancelInviteHandler)
				singleHouseholdInvitationRouter.Put("/accept", householdInvitationsService.AcceptInviteHandler)
				singleHouseholdInvitationRouter.Put("/reject", householdInvitationsService.RejectInviteHandler)
			})
		})

		// OAuth2 Clients
		v1Router.Route("/oauth2_clients", func(clientRouter routing.Router) {
			clientRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
				Get(root, oauth2ClientsService.ListOAuth2ClientsHandler)
			clientRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateOAuth2ClientsPermission)).
				Post(root, oauth2ClientsService.CreateOAuth2ClientHandler)

			singleClientRoute := buildURLVarChunk(oauth2clientsservice.OAuth2ClientIDURIParamKey)
			clientRouter.Route(singleClientRoute, func(singleClientRouter routing.Router) {
				singleClientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
					Get(root, oauth2ClientsService.ReadOAuth2ClientHandler)
				singleClientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveOAuth2ClientsPermission)).
					Delete(root, oauth2ClientsService.ArchiveOAuth2ClientHandler)
			})
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter routing.Router) {
			webhookRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
				Get(root, webhooksService.ListWebhooksHandler)
			webhookRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateWebhooksPermission)).
				Post(root, webhooksService.CreateWebhookHandler)

			singleWebhookRoute := buildURLVarChunk(webhooksservice.WebhookIDURIParamKey)
			webhookRouter.Route(singleWebhookRoute, func(singleWebhookRouter routing.Router) {
				singleWebhookRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
					Get(root, webhooksService.ReadWebhookHandler)
				singleWebhookRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveWebhooksPermission)).
					Delete(root, webhooksService.ArchiveWebhookHandler)

				singleWebhookTriggerEventRoute := buildURLVarChunk(webhooksservice.WebhookTriggerEventIDURIParamKey)
				singleWebhookRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateWebhookTriggerEventsPermission)).
					Post("/trigger_events", webhooksService.AddWebhookTriggerEventHandler)

				singleWebhookRouter.Route("/trigger_events"+singleWebhookTriggerEventRoute, func(singleWebhookTriggerEventRouter routing.Router) {
					singleWebhookTriggerEventRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveWebhookTriggerEventsPermission)).
						Delete(root, webhooksService.ArchiveWebhookTriggerEventHandler)
				})
			})
		})

		// Audit Log Entries
		v1Router.Route("/audit_log_entries", func(auditLogEntriesRouter routing.Router) {
			singleAuditLogEntryRoute := buildURLVarChunk(auditlogentriesservice.AuditLogEntryIDURIParamKey)
			auditLogEntriesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get(singleAuditLogEntryRoute, auditLogEntriesService.ReadAuditLogEntryHandler)
			auditLogEntriesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_user", auditLogEntriesService.ListUserAuditLogEntriesHandler)
			auditLogEntriesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_household", auditLogEntriesService.ListHouseholdAuditLogEntriesHandler)
		})

		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		validInstrumentIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidInstrumentIDURIParamKey)
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter routing.Router) {
			validInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
				Post(root, validEnumerationsService.CreateValidInstrumentHandler)
			validInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(root, validEnumerationsService.ListValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchValidInstrumentsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidInstrumentsHandler)
			validInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(randomRoot, validEnumerationsService.RandomValidInstrumentHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, validEnumerationsService.ReadValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, validEnumerationsService.UpdateValidInstrumentHandler)
				singleValidInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidInstrumentHandler)
			})
		})

		// ValidVessels
		validVesselPath := "valid_vessels"
		validVesselsRouteWithPrefix := fmt.Sprintf("/%s", validVesselPath)
		validVesselIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidVesselIDURIParamKey)
		v1Router.Route(validVesselsRouteWithPrefix, func(validVesselsRouter routing.Router) {
			validVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidVesselsPermission)).
				Post(root, validEnumerationsService.CreateValidVesselHandler)
			validVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(root, validEnumerationsService.ListValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchValidVesselsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidVesselsHandler)
			validVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(randomRoot, validEnumerationsService.RandomValidVesselHandler)

			validVesselsRouter.Route(validVesselIDRouteParam, func(singleValidVesselRouter routing.Router) {
				singleValidVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
					Get(root, validEnumerationsService.ReadValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidVesselsPermission)).
					Put(root, validEnumerationsService.UpdateValidVesselHandler)
				singleValidVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidVesselsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidVesselHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey)
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter routing.Router) {
			validIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientHandler)
			validIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(root, validEnumerationsService.ListValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission, authorization.SearchValidIngredientsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidIngredientsHandler)
			validIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(randomRoot, validEnumerationsService.RandomValidIngredientHandler)

			validIngredientsByPreparationIDSearchRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey))
			validIngredientsRouter.Route(validIngredientsByPreparationIDSearchRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientsByPreparationAndIngredientNameHandler)
			})

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientHandler)
				singleValidIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientHandler)
			})
		})

		// ValidIngredientGroups
		validIngredientGroupPath := "valid_ingredient_groups"
		validIngredientGroupsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientGroupPath)
		validIngredientGroupIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientGroupIDURIParamKey)
		v1Router.Route(validIngredientGroupsRouteWithPrefix, func(validIngredientGroupsRouter routing.Router) {
			validIngredientGroupsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientGroupsPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientGroupHandler)
			validIngredientGroupsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
				Get(root, validEnumerationsService.ListValidIngredientGroupsHandler)
			validIngredientGroupsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchValidIngredientGroupsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidIngredientGroupsHandler)

			validIngredientGroupsRouter.Route(validIngredientGroupIDRouteParam, func(singleValidIngredientGroupRouter routing.Router) {
				singleValidIngredientGroupRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientGroupsPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientGroupHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientGroupsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientGroupHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey)
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter routing.Router) {
			validPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
				Post(root, validEnumerationsService.CreateValidPreparationHandler)
			validPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(root, validEnumerationsService.ListValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidPreparationsHandler)
			validPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(randomRoot, validEnumerationsService.RandomValidPreparationHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, validEnumerationsService.ReadValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, validEnumerationsService.UpdateValidPreparationHandler)
				singleValidPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidPreparationHandler)
			})
		})

		// IngredientPreferences
		userIngredientPreferencesPath := "user_ingredient_preferences"
		userIngredientPreferencesRouteWithPrefix := fmt.Sprintf("/%s", userIngredientPreferencesPath)
		userIngredientPreferencesIDRouteParam := buildURLVarChunk(mealplanningservice.UserIngredientPreferenceIDURIParamKey)
		v1Router.Route(userIngredientPreferencesRouteWithPrefix, func(userIngredientPreferencesRouter routing.Router) {
			userIngredientPreferencesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateUserIngredientPreferencesPermission)).
				Post(root, mealPlanningService.CreateUserIngredientPreferenceHandler)
			userIngredientPreferencesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadUserIngredientPreferencesPermission)).
				Get(root, mealPlanningService.ListUserIngredientPreferencesHandler)

			userIngredientPreferencesRouter.Route(userIngredientPreferencesIDRouteParam, func(singleUserIngredientPreferenceRouter routing.Router) {
				singleUserIngredientPreferenceRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateUserIngredientPreferencesPermission)).
					Put(root, mealPlanningService.UpdateUserIngredientPreferenceHandler)
				singleUserIngredientPreferenceRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveUserIngredientPreferencesPermission)).
					Delete(root, mealPlanningService.ArchiveUserIngredientPreferenceHandler)
			})
		})

		// ValidMeasurementUnits
		validMeasurementUnitPath := "valid_measurement_units"
		validMeasurementUnitsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitPath)
		validMeasurementUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey)
		validMeasurementUnitServiceIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey)
		v1Router.Route(validMeasurementUnitsRouteWithPrefix, func(validMeasurementUnitsRouter routing.Router) {
			validMeasurementUnitsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitsPermission)).
				Post(root, validEnumerationsService.CreateValidMeasurementUnitHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
				Get(root, validEnumerationsService.ListValidMeasurementUnitsHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidMeasurementUnitsHandler)

			validMeasurementUnitsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(path.Join("/by_ingredient", validMeasurementUnitServiceIngredientIDRouteParam), validEnumerationsService.SearchValidMeasurementUnitsByIngredientIDHandler)

			validMeasurementUnitsRouter.Route(validMeasurementUnitIDRouteParam, func(singleValidMeasurementUnitRouter routing.Router) {
				singleValidMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
					Get(root, validEnumerationsService.ReadValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitsPermission)).
					Put(root, validEnumerationsService.UpdateValidMeasurementUnitHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidMeasurementUnitHandler)
			})
		})

		// ValidIngredientStates
		validIngredientStatePath := "valid_ingredient_states"
		validIngredientStatesRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStatePath)
		validIngredientStateIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientStateIDURIParamKey)
		v1Router.Route(validIngredientStatesRouteWithPrefix, func(validIngredientStatesRouter routing.Router) {
			validIngredientStatesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStatesPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientStateHandler)
			validIngredientStatesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(root, validEnumerationsService.ListValidIngredientStatesHandler)
			validIngredientStatesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(searchRoot, validEnumerationsService.SearchValidIngredientStatesHandler)

			validIngredientStatesRouter.Route(validIngredientStateIDRouteParam, func(singleValidIngredientStateRouter routing.Router) {
				singleValidIngredientStateRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStatesPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientStateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStatesPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientStateHandler)
			})
		})

		// ValidMeasurementUnitConversions
		validMeasurementUnitConversionPath := "valid_measurement_conversions"
		validMeasurementUnitConversionsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitConversionPath)
		validMeasurementUnitConversionUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey)
		validMeasurementUnitConversionIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitConversionIDURIParamKey)
		v1Router.Route(validMeasurementUnitConversionsRouteWithPrefix, func(validMeasurementUnitConversionsRouter routing.Router) {
			validMeasurementUnitConversionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitConversionsPermission)).
				Post(root, validEnumerationsService.CreateValidMeasurementUnitConversionHandler)

			validMeasurementUnitConversionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/from_unit", validMeasurementUnitConversionUnitIDRouteParam), validEnumerationsService.ValidMeasurementUnitConversionsFromMeasurementUnitHandler)
			validMeasurementUnitConversionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/to_unit", validMeasurementUnitConversionUnitIDRouteParam), validEnumerationsService.ValidMeasurementUnitConversionsToMeasurementUnitHandler)

			validMeasurementUnitConversionsRouter.Route(validMeasurementUnitConversionIDRouteParam, func(singleValidMeasurementUnitConversionRouter routing.Router) {
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
					Get(root, validEnumerationsService.ReadValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitConversionsPermission)).
					Put(root, validEnumerationsService.UpdateValidMeasurementUnitConversionHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitConversionsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidMeasurementUnitConversionHandler)
			})
		})

		// ValidIngredientStateIngredients
		validIngredientStateIngredientPath := "valid_ingredient_state_ingredients"
		validIngredientStateIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStateIngredientPath)
		validIngredientStateIngredientIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientStateIngredientIDURIParamKey)
		v1Router.Route(validIngredientStateIngredientsRouteWithPrefix, func(validIngredientStateIngredientsRouter routing.Router) {
			validIngredientStateIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStateIngredientsPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientStateIngredientHandler)
			validIngredientStateIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
				Get(root, validEnumerationsService.ListValidIngredientStateIngredientsHandler)

			validIngredientStateIngredientsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientStateIngredientsByIngredientHandler)
			})

			validIngredientStateIngredientsByStateIngredientIDRouteParam := fmt.Sprintf("/by_ingredient_state%s", buildURLVarChunk(validenumerationsservice.ValidIngredientStateIDURIParamKey))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByStateIngredientIDRouteParam, func(byValidStateIngredientIDRouter routing.Router) {
				byValidStateIngredientIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientStateIngredientsByIngredientStateHandler)
			})

			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientIDRouteParam, func(singleValidIngredientStateIngredientRouter routing.Router) {
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStateIngredientsPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientStateIngredientHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStateIngredientsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientStateIngredientHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientPreparationIDURIParamKey)
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter routing.Router) {
			validIngredientPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientPreparationHandler)
			validIngredientPreparationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
				Get(root, validEnumerationsService.ListValidIngredientPreparationsHandler)

			validIngredientPreparationsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientPreparationsByIngredientHandler)
			})

			validIngredientPreparationsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientPreparationsByPreparationHandler)
			})

			validIngredientPreparationsRouter.Route(validIngredientPreparationIDRouteParam, func(singleValidIngredientPreparationRouter routing.Router) {
				singleValidIngredientPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientPreparationHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientPreparationHandler)
			})
		})

		// ValidPreparationInstruments
		validPreparationInstrumentPath := "valid_preparation_instruments"
		validPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationInstrumentPath)
		v1Router.Route(validPreparationInstrumentsRouteWithPrefix, func(validPreparationInstrumentsRouter routing.Router) {
			validPreparationInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
				Post(root, validEnumerationsService.CreateValidPreparationInstrumentHandler)
			validPreparationInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
				Get(root, validEnumerationsService.ListValidPreparationInstrumentsHandler)

			validPreparationInstrumentsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, validEnumerationsService.SearchValidPreparationInstrumentsByPreparationHandler)
			})

			validPreparationInstrumentsByInstrumentIDRouteParam := fmt.Sprintf("/by_instrument%s", buildURLVarChunk(validenumerationsservice.ValidInstrumentIDURIParamKey))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, validEnumerationsService.SearchValidPreparationInstrumentsByInstrumentHandler)
			})

			validPreparationInstrumentIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationInstrumentIDURIParamKey)
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentIDRouteParam, func(singleValidPreparationInstrumentRouter routing.Router) {
				singleValidPreparationInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, validEnumerationsService.ReadValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
					Put(root, validEnumerationsService.UpdateValidPreparationInstrumentHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidPreparationInstrumentHandler)
			})
		})

		// ValidPreparationVessels
		validPreparationVesselPath := "valid_preparation_vessels"
		validPreparationVesselsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationVesselPath)
		v1Router.Route(validPreparationVesselsRouteWithPrefix, func(validPreparationVesselsRouter routing.Router) {
			validPreparationVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidPreparationVesselsPermission)).
				Post(root, validEnumerationsService.CreateValidPreparationVesselHandler)
			validPreparationVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
				Get(root, validEnumerationsService.ListValidPreparationVesselsHandler)

			validPreparationVesselsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validenumerationsservice.ValidPreparationIDURIParamKey))
			validPreparationVesselsRouter.Route(validPreparationVesselsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, validEnumerationsService.SearchValidPreparationVesselsByPreparationHandler)
			})

			validPreparationVesselsByInstrumentIDRouteParam := fmt.Sprintf("/by_vessel%s", buildURLVarChunk(validenumerationsservice.ValidVesselIDURIParamKey))
			validPreparationVesselsRouter.Route(validPreparationVesselsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, validEnumerationsService.SearchValidPreparationVesselsByVesselHandler)
			})

			validPreparationVesselIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidPreparationVesselIDURIParamKey)
			validPreparationVesselsRouter.Route(validPreparationVesselIDRouteParam, func(singleValidPreparationVesselRouter routing.Router) {
				singleValidPreparationVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, validEnumerationsService.ReadValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationVesselsPermission)).
					Put(root, validEnumerationsService.UpdateValidPreparationVesselHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationVesselsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidPreparationVesselHandler)
			})
		})

		// ValidIngredientMeasurementUnit
		validIngredientMeasurementUnitPath := "valid_ingredient_measurement_units"
		validIngredientMeasurementUnitRouteWithPrefix := fmt.Sprintf("/%s", validIngredientMeasurementUnitPath)
		validIngredientMeasurementUnitIDRouteParam := buildURLVarChunk(validenumerationsservice.ValidIngredientMeasurementUnitIDURIParamKey)
		v1Router.Route(validIngredientMeasurementUnitRouteWithPrefix, func(validIngredientMeasurementUnitRouter routing.Router) {
			validIngredientMeasurementUnitRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateValidIngredientMeasurementUnitsPermission)).
				Post(root, validEnumerationsService.CreateValidIngredientMeasurementUnitHandler)
			validIngredientMeasurementUnitRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
				Get(root, validEnumerationsService.ListValidIngredientMeasurementUnitsHandler)

			validIngredientMeasurementUnitsByIngredientRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validenumerationsservice.ValidIngredientIDURIParamKey))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByIngredientRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientMeasurementUnitsByIngredientHandler)
			})

			validIngredientMeasurementUnitsByMeasurementUnitRouteParam := fmt.Sprintf("/by_measurement_unit%s", buildURLVarChunk(validenumerationsservice.ValidMeasurementUnitIDURIParamKey))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByMeasurementUnitRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, validEnumerationsService.SearchValidIngredientMeasurementUnitsByMeasurementUnitHandler)
			})

			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitIDRouteParam, func(singleValidIngredientMeasurementUnitRouter routing.Router) {
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, validEnumerationsService.ReadValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientMeasurementUnitsPermission)).
					Put(root, validEnumerationsService.UpdateValidIngredientMeasurementUnitHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientMeasurementUnitsPermission)).
					Delete(root, validEnumerationsService.ArchiveValidIngredientMeasurementUnitHandler)
			})
		})

		// MealPlanning
		mealPath := "meals"
		mealsRouteWithPrefix := fmt.Sprintf("/%s", mealPath)
		mealIDRouteParam := buildURLVarChunk(mealplanningservice.MealIDURIParamKey)
		v1Router.Route(mealsRouteWithPrefix, func(mealsRouter routing.Router) {
			mealsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealsPermission)).
				Post(root, mealPlanningService.CreateMealHandler)
			mealsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(root, mealPlanningService.ListMealsHandler)
			mealsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(searchRoot, mealPlanningService.SearchMealsHandler)

			mealsRouter.Route(mealIDRouteParam, func(singleMealRouter routing.Router) {
				singleMealRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
					Get(root, mealPlanningService.ReadMealHandler)
				singleMealRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealsPermission)).
					Delete(root, mealPlanningService.ArchiveMealHandler)
			})
		})

		// Components
		recipePath := "recipes"
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		recipeIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeIDURIParamKey)
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, recipeManagementService.CreateRecipeHandler)
			recipesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, recipeManagementService.ListRecipesHandler)
			recipesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(searchRoot, recipeManagementService.SearchRecipesHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, recipeManagementService.ReadRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/mermaid", recipeManagementService.RecipeMermaidHandler)
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/prep_steps", recipeManagementService.RecipeEstimatedPrepStepsHandler)
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, recipeManagementService.UpdateRecipeHandler)
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Post("/clone", recipeManagementService.CloneRecipeHandler)

				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Post("/images", recipeManagementService.RecipeImageUploadHandler)
				singleRecipeRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeHandler)

				// RecipeRatings
				singleRecipeRouter.Route("/ratings", func(recipeRatingsRouter routing.Router) {
					recipeRatingsRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
						Post(root, recipeManagementService.CreateRecipeRatingHandler)
					recipeRatingsRouter.
						WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeRatingsPermission)).
						Get(root, recipeManagementService.ListRecipeRatingsHandler)

					recipeRatingIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeRatingIDURIParamKey)
					recipeRatingsRouter.Route(recipeRatingIDRouteParam, func(singleRecipeRatingRouter routing.Router) {
						singleRecipeRatingRouter.
							WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
							Get(root, recipeManagementService.ReadRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeRatingsPermission)).
							Put(root, recipeManagementService.UpdateRecipeRatingHandler)
						singleRecipeRatingRouter.
							WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeRatingsPermission)).
							Delete(root, recipeManagementService.ArchiveRecipeRatingHandler)
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
		recipePrepTaskIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipePrepTaskIDURIParamKey)
		v1Router.Route(recipePrepTasksRouteWithPrefix, func(recipePrepTasksRouter routing.Router) {
			recipePrepTasksRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipePrepTasksPermission)).
				Post(root, recipeManagementService.CreateRecipePrepTaskHandler)
			recipePrepTasksRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
				Get(root, recipeManagementService.ListRecipePrepTaskHandler)

			recipePrepTasksRouter.Route(recipePrepTaskIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
					Get(root, recipeManagementService.ReadRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipePrepTasksPermission)).
					Put(root, recipeManagementService.UpdateRecipePrepTaskHandler)
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipePrepTasksPermission)).
					Delete(root, recipeManagementService.ArchiveRecipePrepTaskHandler)
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
		recipeStepIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepIDURIParamKey)
		v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter routing.Router) {
			recipeStepsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepHandler)
			recipeStepsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
				Get(root, recipeManagementService.ListRecipeStepsHandler)

			recipeStepsRouter.Route(recipeStepIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Post("/images", recipeManagementService.RecipeStepImageUploadHandler)
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepHandler)
				singleRecipeStepRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepHandler)
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
		recipeStepInstrumentIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepInstrumentIDURIParamKey)
		v1Router.Route(recipeStepInstrumentsRouteWithPrefix, func(recipeStepInstrumentsRouter routing.Router) {
			recipeStepInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepInstrumentsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepInstrumentHandler)
			recipeStepInstrumentsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
				Get(root, recipeManagementService.ListRecipeStepInstrumentsHandler)

			recipeStepInstrumentsRouter.Route(recipeStepInstrumentIDRouteParam, func(singleRecipeStepInstrumentRouter routing.Router) {
				singleRecipeStepInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepInstrumentsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepInstrumentHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepInstrumentsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepInstrumentHandler)
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
		recipeStepVesselIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepVesselIDURIParamKey)
		v1Router.Route(recipeStepVesselsRouteWithPrefix, func(recipeStepVesselsRouter routing.Router) {
			recipeStepVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepVesselsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepVesselHandler)
			recipeStepVesselsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
				Get(root, recipeManagementService.ListRecipeStepVesselsHandler)

			recipeStepVesselsRouter.Route(recipeStepVesselIDRouteParam, func(singleRecipeStepVesselRouter routing.Router) {
				singleRecipeStepVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepVesselsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepVesselHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepVesselsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepVesselHandler)
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
		recipeStepIngredientIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepIngredientIDURIParamKey)
		v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter routing.Router) {
			recipeStepIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepIngredientsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepIngredientHandler)
			recipeStepIngredientsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
				Get(root, recipeManagementService.ListRecipeStepIngredientsHandler)

			recipeStepIngredientsRouter.Route(recipeStepIngredientIDRouteParam, func(singleRecipeStepIngredientRouter routing.Router) {
				singleRecipeStepIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepIngredientHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepIngredientHandler)
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
		recipeStepCompletionConditionIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepCompletionConditionIDURIParamKey)
		v1Router.Route(recipeStepCompletionConditionsRouteWithPrefix, func(recipeStepCompletionConditionsRouter routing.Router) {
			recipeStepCompletionConditionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepCompletionConditionsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepCompletionConditionHandler)
			recipeStepCompletionConditionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
				Get(root, recipeManagementService.ListRecipeStepCompletionConditionsHandler)

			recipeStepCompletionConditionsRouter.Route(recipeStepCompletionConditionIDRouteParam, func(singleRecipeStepCompletionConditionRouter routing.Router) {
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepCompletionConditionsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepCompletionConditionHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepCompletionConditionsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepCompletionConditionHandler)
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
		recipeStepProductIDRouteParam := buildURLVarChunk(recipemanagementservice.RecipeStepProductIDURIParamKey)
		v1Router.Route(recipeStepProductsRouteWithPrefix, func(recipeStepProductsRouter routing.Router) {
			recipeStepProductsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateRecipeStepProductsPermission)).
				Post(root, recipeManagementService.CreateRecipeStepProductHandler)
			recipeStepProductsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
				Get(root, recipeManagementService.ListRecipeStepProductsHandler)

			recipeStepProductsRouter.Route(recipeStepProductIDRouteParam, func(singleRecipeStepProductRouter routing.Router) {
				singleRecipeStepProductRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Get(root, recipeManagementService.ReadRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, recipeManagementService.UpdateRecipeStepProductHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, recipeManagementService.ArchiveRecipeStepProductHandler)
			})
		})

		// MealPlans
		mealPlanPath := "meal_plans"
		mealPlansRouteWithPrefix := fmt.Sprintf("/%s", mealPlanPath)
		mealPlanIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanIDURIParamKey)
		v1Router.Route(mealPlansRouteWithPrefix, func(mealPlansRouter routing.Router) {
			mealPlansRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
				Post(root, mealPlanningService.CreateMealPlanHandler)
			mealPlansRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
				Get(root, mealPlanningService.ListMealPlanHandler)

			mealPlansRouter.Route(mealPlanIDRouteParam, func(singleMealPlanRouter routing.Router) {
				singleMealPlanRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
					Get(root, mealPlanningService.ReadMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlansPermission)).
					Put(root, mealPlanningService.UpdateMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealPlansPermission)).
					Delete(root, mealPlanningService.ArchiveMealPlanHandler)
				singleMealPlanRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
					Post("/finalize", mealPlanningService.FinalizeMealPlanHandler)
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
		mealPlanTaskIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanTaskIDURIParamKey)
		v1Router.Route(mealPlanTasksRouteWithPrefix, func(mealPlanTasksRouter routing.Router) {
			mealPlanTasksRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
				Get(root, mealPlanningService.ListMealPlanTasksByMealPlanHandler)
			mealPlanTasksRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlanTasksPermission)).
				Post(root, mealPlanningService.CreateMealPlanTaskHandler)

			mealPlanTasksRouter.Route(mealPlanTaskIDRouteParam, func(singleMealPlanTaskRouter routing.Router) {
				singleMealPlanTaskRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
					Get(root, mealPlanningService.ReadMealPlanTaskHandler)

				singleMealPlanTaskRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlanTasksPermission)).
					Patch(root, mealPlanningService.MealPlanTaskStatusChangeHandler)
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
		mealPlanEventIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanEventIDURIParamKey)
		v1Router.Route(mealPlanEventsRouteWithPrefix, func(mealPlanEventsRouter routing.Router) {
			mealPlanEventsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlanEventsPermission)).
				Post(root, mealPlanningService.CreateMealPlanEventHandler)
			mealPlanEventsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
				Get(root, mealPlanningService.ListMealPlanEventHandler)

			mealPlanEventsRouter.Route(mealPlanEventIDRouteParam, func(singleMealPlanEventRouter routing.Router) {
				singleMealPlanEventRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
					Get(root, mealPlanningService.ReadMealPlanEventHandler)
				singleMealPlanEventRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlanEventsPermission)).
					Put(root, mealPlanningService.UpdateMealPlanEventHandler)
				singleMealPlanEventRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionVotesPermission)).
					Post("/vote", mealPlanningService.CreateMealPlanOptionVoteHandler)
				singleMealPlanEventRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanEventsPermission)).
					Delete(root, mealPlanningService.ArchiveMealPlanEventHandler)
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
		mealPlanGroceryListItemIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanGroceryListItemIDURIParamKey)
		v1Router.Route(mealPlanGroceryListItemsRouteWithPrefix, func(mealPlanGroceryListItemsRouter routing.Router) {
			mealPlanGroceryListItemsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlanGroceryListItemsPermission)).
				Post(root, mealPlanningService.CreateMealPlanGroceryListItemHandler)
			mealPlanGroceryListItemsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
				Get(root, mealPlanningService.ListMealPlanGroceryListItemsByMealPlanHandler)

			mealPlanGroceryListItemsRouter.Route(mealPlanGroceryListItemIDRouteParam, func(singleMealPlanGroceryListItemRouter routing.Router) {
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
					Get(root, mealPlanningService.ReadMealPlanGroceryListItemHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlanGroceryListItemsPermission)).
					Put(root, mealPlanningService.UpdateMealPlanGroceryListItemHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanGroceryListItemsPermission)).
					Delete(root, mealPlanningService.ArchiveMealPlanGroceryListItemHandler)
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
		mealPlanOptionIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanOptionIDURIParamKey)
		v1Router.Route(mealPlanOptionsRouteWithPrefix, func(mealPlanOptionsRouter routing.Router) {
			mealPlanOptionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionsPermission)).
				Post(root, mealPlanningService.CreateMealPlanOptionHandler)
			mealPlanOptionsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
				Get(root, mealPlanningService.ListMealPlanOptionHandler)

			mealPlanOptionsRouter.Route(mealPlanOptionIDRouteParam, func(singleMealPlanOptionRouter routing.Router) {
				singleMealPlanOptionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
					Get(root, mealPlanningService.ReadMealPlanOptionHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionsPermission)).
					Put(root, mealPlanningService.UpdateMealPlanOptionHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionsPermission)).
					Delete(root, mealPlanningService.ArchiveMealPlanOptionHandler)
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
		mealPlanOptionVoteIDRouteParam := buildURLVarChunk(mealplanningservice.MealPlanOptionVoteIDURIParamKey)
		v1Router.Route(mealPlanOptionVotesRouteWithPrefix, func(mealPlanOptionVotesRouter routing.Router) {
			mealPlanOptionVotesRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
				Get(root, mealPlanningService.ListMealPlanOptionVoteHandler)

			mealPlanOptionVotesRouter.Route(mealPlanOptionVoteIDRouteParam, func(singleMealPlanOptionVoteRouter routing.Router) {
				singleMealPlanOptionVoteRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
					Get(root, mealPlanningService.ReadMealPlanOptionVoteHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionVotesPermission)).
					Put(root, mealPlanningService.UpdateMealPlanOptionVoteHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionVotesPermission)).
					Delete(root, mealPlanningService.ArchiveMealPlanOptionVoteHandler)
			})
		})

		// ServiceSettings
		serviceSettingPath := "settings"
		serviceSettingsRouteWithPrefix := fmt.Sprintf("/%s", serviceSettingPath)
		serviceSettingIDRouteParam := buildURLVarChunk(servicesettingsservice.ServiceSettingIDURIParamKey)
		v1Router.Route(serviceSettingsRouteWithPrefix, func(serviceSettingsRouter routing.Router) {
			serviceSettingsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateServiceSettingsPermission)).
				Post(root, serviceSettingsService.CreateServiceSettingHandler)
			serviceSettingsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(root, serviceSettingsService.ListServiceSettingsHandler)
			serviceSettingsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(searchRoot, serviceSettingsService.SearchServiceSettingsHandler)

			serviceSettingsRouter.Route(serviceSettingIDRouteParam, func(singleServiceSettingRouter routing.Router) {
				singleServiceSettingRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
					Get(root, serviceSettingsService.ReadServiceSettingHandler)
				singleServiceSettingRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingsPermission)).
					Delete(root, serviceSettingsService.ArchiveServiceSettingHandler)
			})

			serviceSettingConfigurationIDRouteParam := buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationIDURIParamKey)
			serviceSettingsRouter.Route("/configurations", func(settingConfigurationRouter routing.Router) {
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateServiceSettingConfigurationsPermission)).
					Post(root, serviceSettingConfigurationsService.CreateServiceSettingConfigurationHandler)
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user"+buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationNameURIParamKey), serviceSettingConfigurationsService.GetServiceSettingConfigurationsForUserByNameHandler)
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user", serviceSettingConfigurationsService.GetServiceSettingConfigurationsForUserHandler)
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/household", serviceSettingConfigurationsService.GetServiceSettingConfigurationsForHouseholdHandler)
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateServiceSettingConfigurationsPermission)).
					Put(serviceSettingConfigurationIDRouteParam, serviceSettingConfigurationsService.UpdateServiceSettingConfigurationHandler)
				settingConfigurationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingConfigurationsPermission)).
					Delete(serviceSettingConfigurationIDRouteParam, serviceSettingConfigurationsService.ArchiveServiceSettingConfigurationHandler)
			})
		})

		// User Notifications
		userNotificationPath := "user_notifications"
		userNotificationsRouteWithPrefix := fmt.Sprintf("/%s", userNotificationPath)
		userNotificationIDRouteParam := buildURLVarChunk(usernotificationsservice.UserNotificationIDURIParamKey)
		v1Router.Route(userNotificationsRouteWithPrefix, func(userNotificationsRouter routing.Router) {
			userNotificationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.CreateUserNotificationsPermission)).
				Post(root, userNotificationsService.CreateUserNotificationHandler)
			userNotificationsRouter.
				WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
				Get(root, userNotificationsService.ListUserNotificationsHandler)

			userNotificationsRouter.Route(userNotificationIDRouteParam, func(singleUserNotificationRouter routing.Router) {
				singleUserNotificationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
					Get(root, userNotificationsService.ReadUserNotificationHandler)
				singleUserNotificationRouter.
					WithMiddleware(authService.PermissionFilterMiddleware(authorization.UpdateUserNotificationsPermission)).
					Patch(root, userNotificationsService.UpdateUserNotificationHandler)
			})
		})
	})

	return router, nil
}
