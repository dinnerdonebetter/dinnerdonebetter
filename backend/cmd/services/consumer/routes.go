package main

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
	"github.com/dinnerdonebetter/backend/internal/platform/version"

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
		metaRouter.Get("/commit", func(res http.ResponseWriter, req *http.Request) {
			s.encoder.EncodeResponseWithStatus(req.Context(), res, version.Get(), http.StatusOK)
		})
	})

	// Apple App Site Association (public)
	router.Get("/.well-known/apple-app-site-association", s.AppleAppSiteAssociationHandler)

	// Auth (public)
	router.Get("/login", ghttp.Adapt(s.LoginPage))
	router.Post("/login/submit", ghttp.Adapt(s.LoginSubmission))

	// Accept invitation - Option A: show "Open in app" page (public)
	router.Get("/accept_invitation", ghttp.Adapt(s.AcceptInvitationPage))
	router.Get("/accept_invitation/*", ghttp.Adapt(s.AcceptInvitationPage))

	// Legal (public)
	router.Get("/terms-of-service", ghttp.Adapt(s.TermsPage))
	router.Get("/privacy-policy", ghttp.Adapt(s.PrivacyPage))

	// Protected routes
	r.Get("/", ghttp.Adapt(s.HomePage))
	r.Get("/account/settings", ghttp.Adapt(s.AccountSettingsPage))
	r.Get("/account", ghttp.Adapt(s.AccountSettingsPage))

	// static files - NOTE: this must be registered last
	router.Get("/*", phttp.RootLevelAssetsHandler(assetsDir))
}
