package frontend

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
)

const (
	numericIDPattern                 = "{%s:[0-9]+}"
	unauthorizedRedirectResponseCode = http.StatusSeeOther
)

// SetupRoutes sets up the routes.
func (s *service) SetupRoutes(router routing.Router) {
	router = router.WithMiddleware(s.authService.UserAttributionMiddleware)

	router.Get("/", s.homepage)
	router.Get("/dashboard", s.homepage)

	// statics
	router.Get("/favicon.svg", s.favicon)

	// auth stuff
	router.Get("/login", s.buildLoginView(true))
	router.Get("/components/login_prompt", s.buildLoginView(false))
	router.Post("/auth/submit_login", s.handleLoginSubmission)
	router.Post("/logout", s.handleLogoutSubmission)

	router.Get("/register", s.registrationView)
	router.Get("/components/registration_prompt", s.registrationComponent)
	router.Post("/auth/submit_registration", s.handleRegistrationSubmission)
	router.Post("/auth/verify_two_factor_secret", s.handleTOTPVerificationSubmission)

	router.Post("/billing/checkout/begin", s.handleCheckoutSessionStart)
	router.Post("/billing/checkout/success", s.handleCheckoutSuccess)
	router.Post("/billing/checkout/cancel", s.handleCheckoutCancel)
	router.Post("/billing/checkout/failures", s.handleCheckoutFailure)

	singleAccountPattern := fmt.Sprintf(numericIDPattern, accountIDURLParamKey)
	router.Get("/accounts", s.buildAccountsTableView(true))
	router.Get(fmt.Sprintf("/accounts/%s", singleAccountPattern), s.buildAccountEditorView(true))
	router.Get("/dashboard_pages/accounts", s.buildAccountsTableView(false))
	router.Get(fmt.Sprintf("/dashboard_pages/accounts/%s", singleAccountPattern), s.buildAccountEditorView(false))

	singleAPIClientPattern := fmt.Sprintf(numericIDPattern, apiClientIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
		Get("/api_clients", s.buildAPIClientsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
		Get("/dashboard_pages/api_clients", s.buildAPIClientsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
		Get(fmt.Sprintf("/api_clients/%s", singleAPIClientPattern), s.buildAPIClientEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAPIClientsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/api_clients/%s", singleAPIClientPattern), s.buildAPIClientEditorView(false))

	singleWebhookPattern := fmt.Sprintf(numericIDPattern, webhookIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
		Get("/account/webhooks", s.buildWebhooksTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
		Get("/dashboard_pages/account/webhooks", s.buildWebhooksTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateWebhooksPermission)).
		Get(fmt.Sprintf("/account/webhooks/%s", singleWebhookPattern), s.buildWebhookEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateWebhooksPermission)).
		Get(fmt.Sprintf("/dashboard_pages/account/webhooks/%s", singleWebhookPattern), s.buildWebhookEditorView(false))

	router.Get("/user/settings", s.buildUserSettingsView(true))
	router.Get("/dashboard_pages/user/settings", s.buildUserSettingsView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateAccountPermission)).
		Get("/account/settings", s.buildAccountSettingsView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateAccountPermission)).
		Get("/dashboard_pages/account/settings", s.buildAccountSettingsView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
		Get("/admin/users/search", s.buildUsersTableView(true, true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
		Get("/dashboard_pages/admin/users/search", s.buildUsersTableView(false, true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
		Get("/admin/users", s.buildUsersTableView(true, false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
		Get("/dashboard_pages/admin/users", s.buildUsersTableView(false, false))
	router.WithMiddleware(s.authService.ServiceAdminMiddleware).
		Get("/admin/settings", s.buildAdminSettingsView(true))
	router.WithMiddleware(s.authService.ServiceAdminMiddleware).
		Get("/dashboard_pages/admin/settings", s.buildAdminSettingsView(false))

	singleValidInstrumentPattern := fmt.Sprintf(numericIDPattern, validInstrumentIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
		Get("/valid_instruments", s.buildValidInstrumentsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
		Get("/dashboard_pages/valid_instruments", s.buildValidInstrumentsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
		Get("/valid_instruments/new", s.buildValidInstrumentCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
		Post("/valid_instruments/new/submit", s.handleValidInstrumentCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/valid_instruments/%s", singleValidInstrumentPattern), s.handleValidInstrumentArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
		Get("/dashboard_pages/valid_instruments/new", s.buildValidInstrumentCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
		Get(fmt.Sprintf("/valid_instruments/%s", singleValidInstrumentPattern), s.buildValidInstrumentEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/valid_instruments/%s", singleValidInstrumentPattern), s.handleValidInstrumentUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/valid_instruments/%s", singleValidInstrumentPattern), s.buildValidInstrumentEditorView(false))

	singleValidPreparationPattern := fmt.Sprintf(numericIDPattern, validPreparationIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
		Get("/valid_preparations", s.buildValidPreparationsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
		Get("/dashboard_pages/valid_preparations", s.buildValidPreparationsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
		Get("/valid_preparations/new", s.buildValidPreparationCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
		Post("/valid_preparations/new/submit", s.handleValidPreparationCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/valid_preparations/%s", singleValidPreparationPattern), s.handleValidPreparationArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
		Get("/dashboard_pages/valid_preparations/new", s.buildValidPreparationCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
		Get(fmt.Sprintf("/valid_preparations/%s", singleValidPreparationPattern), s.buildValidPreparationEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/valid_preparations/%s", singleValidPreparationPattern), s.handleValidPreparationUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/valid_preparations/%s", singleValidPreparationPattern), s.buildValidPreparationEditorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidIngredientsPermission)).
		Get("/elements/valid_preparations/search", s.validPreparationsSearchResults)

	singleValidIngredientPattern := fmt.Sprintf(numericIDPattern, validIngredientIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
		Get("/valid_ingredients", s.buildValidIngredientsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
		Get("/dashboard_pages/valid_ingredients", s.buildValidIngredientsTableView(false))

	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
		Get("/valid_ingredients/new", s.buildValidIngredientCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
		Post("/valid_ingredients/new/submit", s.handleValidIngredientCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/valid_ingredients/%s", singleValidIngredientPattern), s.handleValidIngredientArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
		Get("/dashboard_pages/valid_ingredients/new", s.buildValidIngredientCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
		Get(fmt.Sprintf("/valid_ingredients/%s", singleValidIngredientPattern), s.buildValidIngredientEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/valid_ingredients/%s", singleValidIngredientPattern), s.handleValidIngredientUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/valid_ingredients/%s", singleValidIngredientPattern), s.buildValidIngredientEditorView(false))

	singleValidIngredientPreparationPattern := fmt.Sprintf(numericIDPattern, validIngredientPreparationIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
		Get("/valid_ingredient_preparations", s.buildValidIngredientPreparationsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
		Get("/dashboard_pages/valid_ingredient_preparations", s.buildValidIngredientPreparationsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
		Get("/valid_ingredient_preparations/new", s.buildValidIngredientPreparationCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
		Post("/valid_ingredient_preparations/new/submit", s.handleValidIngredientPreparationCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/valid_ingredient_preparations/%s", singleValidIngredientPreparationPattern), s.handleValidIngredientPreparationArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
		Get("/dashboard_pages/valid_ingredient_preparations/new", s.buildValidIngredientPreparationCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
		Get(fmt.Sprintf("/valid_ingredient_preparations/%s", singleValidIngredientPreparationPattern), s.buildValidIngredientPreparationEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/valid_ingredient_preparations/%s", singleValidIngredientPreparationPattern), s.handleValidIngredientPreparationUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/valid_ingredient_preparations/%s", singleValidIngredientPreparationPattern), s.buildValidIngredientPreparationEditorView(false))

	singleValidPreparationInstrumentPattern := fmt.Sprintf(numericIDPattern, validPreparationInstrumentIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
		Get("/valid_preparation_instruments", s.buildValidPreparationInstrumentsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
		Get("/dashboard_pages/valid_preparation_instruments", s.buildValidPreparationInstrumentsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
		Get("/valid_preparation_instruments/new", s.buildValidPreparationInstrumentCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
		Post("/valid_preparation_instruments/new/submit", s.handleValidPreparationInstrumentCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/valid_preparation_instruments/%s", singleValidPreparationInstrumentPattern), s.handleValidPreparationInstrumentArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
		Get("/dashboard_pages/valid_preparation_instruments/new", s.buildValidPreparationInstrumentCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
		Get(fmt.Sprintf("/valid_preparation_instruments/%s", singleValidPreparationInstrumentPattern), s.buildValidPreparationInstrumentEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/valid_preparation_instruments/%s", singleValidPreparationInstrumentPattern), s.handleValidPreparationInstrumentUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/valid_preparation_instruments/%s", singleValidPreparationInstrumentPattern), s.buildValidPreparationInstrumentEditorView(false))

	singleRecipePattern := fmt.Sprintf(numericIDPattern, recipeIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
		Get("/recipes", s.buildRecipesTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
		Get("/dashboard_pages/recipes", s.buildRecipesTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
		Get("/recipes/new", s.buildRecipeCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/recipes/%s", singleRecipePattern), s.handleRecipeArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
		Get("/dashboard_pages/recipes/new", s.buildRecipeCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
		Get(fmt.Sprintf("/recipes/%s", singleRecipePattern), s.buildRecipeEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
		Put(fmt.Sprintf("/dashboard_pages/recipes/%s", singleRecipePattern), s.handleRecipeUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
		Get(fmt.Sprintf("/dashboard_pages/recipes/%s", singleRecipePattern), s.buildRecipeEditorView(false))

	singleInvitationPattern := fmt.Sprintf(numericIDPattern, invitationIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsPermission)).
		Get("/invitations", s.buildInvitationsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadInvitationsPermission)).
		Get("/dashboard_pages/invitations", s.buildInvitationsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateInvitationsPermission)).
		Get("/invitations/new", s.buildInvitationCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateInvitationsPermission)).
		Post("/invitations/new/submit", s.handleInvitationCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveInvitationsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/invitations/%s", singleInvitationPattern), s.handleInvitationArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveInvitationsPermission)).
		Get("/dashboard_pages/invitations/new", s.buildInvitationCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateInvitationsPermission)).
		Get(fmt.Sprintf("/invitations/%s", singleInvitationPattern), s.buildInvitationEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateInvitationsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/invitations/%s", singleInvitationPattern), s.handleInvitationUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateInvitationsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/invitations/%s", singleInvitationPattern), s.buildInvitationEditorView(false))

	singleReportPattern := fmt.Sprintf(numericIDPattern, reportIDURLParamKey)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsPermission)).
		Get("/reports", s.buildReportsTableView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadReportsPermission)).
		Get("/dashboard_pages/reports", s.buildReportsTableView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateReportsPermission)).
		Get("/reports/new", s.buildReportCreatorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateReportsPermission)).
		Post("/reports/new/submit", s.handleReportCreationRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveReportsPermission)).
		Delete(fmt.Sprintf("/dashboard_pages/reports/%s", singleReportPattern), s.handleReportArchiveRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveReportsPermission)).
		Get("/dashboard_pages/reports/new", s.buildReportCreatorView(false))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateReportsPermission)).
		Get(fmt.Sprintf("/reports/%s", singleReportPattern), s.buildReportEditorView(true))
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateReportsPermission)).
		Put(fmt.Sprintf("/dashboard_pages/reports/%s", singleReportPattern), s.handleReportUpdateRequest)
	router.WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateReportsPermission)).
		Get(fmt.Sprintf("/dashboard_pages/reports/%s", singleReportPattern), s.buildReportEditorView(false))

	router.LogRoutes()
}
