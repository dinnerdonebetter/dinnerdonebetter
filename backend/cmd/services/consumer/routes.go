package main

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/internal/platform/version"
	"github.com/dinnerdonebetter/backend/internal/services/auth/handlers/passkey"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"
	"github.com/dinnerdonebetter/backend/pkg/client"

	ghttp "maragu.dev/gomponents/http"
)

const (
	assetsDir = "./cmd/services/consumer/assets"
)

func (s *ConsumerFrontendServer) setupRoutes(router routing.Router) {
	r := router.WithMiddleware(s.authMiddleware)

	// Health checks (public)
	router.Route("/_ops_", func(metaRouter routing.Router) {
		metaRouter.Get("/live", func(res http.ResponseWriter, _ *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
		metaRouter.Get("/ready", func(res http.ResponseWriter, _ *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
		metaRouter.Get("/version", func(res http.ResponseWriter, req *http.Request) {
			s.encoder.EncodeResponseWithStatus(req.Context(), res, version.Get(), http.StatusOK)
		})
	})

	// Apple App Site Association (public)
	router.Get("/.well-known/apple-app-site-association", s.AppleAppSiteAssociationHandler)

	// Auth (public)
	router.Get("/login", ghttp.Adapt(s.LoginPage))
	router.Post("/login/submit", ghttp.Adapt(s.LoginSubmission))
	passkeyHandlers := passkey.NewHandlers(s.tracer, s.logger, s.encoder, s.cookieManager, &s.config.Cookies, s.buildUnauthedGRPCClient, func(r *http.Request) (client.Client, error) {
		return webappauth.ClientFromContext(r.Context())
	})
	router.Post("/auth/passkey/authentication/options", passkeyHandlers.AuthOptionsHandler)
	router.Post("/auth/passkey/authentication/verify", passkeyHandlers.AuthVerifyHandler)
	router.Get("/verify_email_address", ghttp.Adapt(s.VerifyEmailAddressPage))
	router.Get("/forgot_password", ghttp.Adapt(s.ForgotPasswordPage))
	router.Post("/forgot_password/submit", ghttp.Adapt(s.ForgotPasswordSubmission))
	router.Get("/reset_password", ghttp.Adapt(s.ResetPasswordPage))
	router.Post("/reset_password/submit", ghttp.Adapt(s.ResetPasswordSubmission))

	// Accept invitation - Option A: show "Open in app" page (public)
	router.Get("/accept_invitation", ghttp.Adapt(s.AcceptInvitationPage))
	router.Get("/accept_invitation/*", ghttp.Adapt(s.AcceptInvitationPage))

	// Legal (public)
	router.Get("/terms-of-service", ghttp.Adapt(s.TermsPage))
	router.Get("/privacy-policy", ghttp.Adapt(s.PrivacyPage))

	// Logout (public - clearing cookie works regardless of auth state)
	router.Get("/logout", s.LogoutHandler)
	router.Post("/logout", s.LogoutHandler)

	// Protected routes
	r.Get("/", ghttp.Adapt(s.HomePage))
	r.Get("/account/settings", ghttp.Adapt(s.AccountSettingsPage))
	r.Get("/account", ghttp.Adapt(s.AccountSettingsPage))
	r.Get("/account/household-members", ghttp.Adapt(s.HouseholdMembersPage))
	r.Post("/account/household-members/send-invitation", s.SendInvitationHandler)
	r.Post("/account/household-members/cancel-invitation", s.CancelInvitationHandler)
	r.Post("/account/household-members/update-role", s.UpdateMemberRoleHandler)
	r.Get("/account/household-details", ghttp.Adapt(s.HouseholdDetailsPage))
	r.Post("/account/household-details/update", s.UpdateHouseholdDetailsHandler)
	r.Get("/account/preferences", ghttp.Adapt(s.PreferencesPage))
	r.Post("/account/preferences/update", s.UpdatePreferenceHandler)
	r.Get("/account/profile", ghttp.Adapt(s.ProfilePage))
	r.Post("/account/profile/update-username", s.UpdateProfileUsernameHandler)
	r.Post("/account/profile/update-details", s.UpdateProfileDetailsHandler)
	r.Get("/account/passkeys", ghttp.Adapt(s.PasskeysPage))
	r.Post("/account/passkeys/delete", s.DeletePasskeyHandler)
	r.Post("/auth/passkey/registration/options", passkeyHandlers.RegOptionsHandler)
	r.Post("/auth/passkey/registration/verify", passkeyHandlers.RegVerifyHandler)

	// static files - NOTE: this must be registered last
	router.Get("/*", phttp.RootLevelAssetsHandler(assetsDir))
}
