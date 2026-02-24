package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	routingcfg "github.com/dinnerdonebetter/backend/internal/platform/routing/config"
	paymentswebhook "github.com/dinnerdonebetter/backend/internal/services/payments/http"
)

// appleAppSiteAssociation represents the structure of the apple-app-site-association file
// used by iOS for Universal Links and Password AutoFill (webcredentials).
type appleAppSiteAssociation struct {
	AppLinks       appLinks       `json:"applinks"`
	WebCredentials webCredentials `json:"webcredentials"`
}

type webCredentials struct {
	Apps []string `json:"apps"`
}

type appLinks struct {
	Apps    []string    `json:"apps"`
	Details []appDetail `json:"details"`
}

type appDetail struct {
	AppID string   `json:"appID"`
	Paths []string `json:"paths"`
}

func ProvideAPIRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	authService auth.AuthDataService,
	aasaConfig config.AppleAppSiteAssociationConfig,
	paymentsWebhookHandler *paymentswebhook.WebhookHandler,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	router.Route("/_ops_", func(metaRouter routing.Router) {
		// Expose a liveness check on /live
		metaRouter.Get("/live", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			// TODO: check readiness here lol
			res.WriteHeader(http.StatusOK)
		})
	})

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.Get("/authorize", authService.AuthorizeHandler)
		userRouter.Post("/token", authService.TokenHandler)
	})

	router.Route("/api/payments/webhooks", func(paymentsRouter routing.Router) {
		paymentsRouter.Post("/{provider}", paymentsWebhookHandler.Handle)
	})

	// Apple App Site Association for iOS Universal Links
	// This endpoint is required for iOS to recognize the app as a handler for URLs on this domain.
	// See: https://developer.apple.com/documentation/xcode/supporting-associated-domains
	if aasaConfig.TeamID != "" && aasaConfig.BundleID != "" {
		appID := fmt.Sprintf("%s.%s", aasaConfig.TeamID, aasaConfig.BundleID)
		router.Route("/.well-known", func(wellKnownRouter routing.Router) {
			wellKnownRouter.Get("/apple-app-site-association", func(res http.ResponseWriter, req *http.Request) {
				aasa := appleAppSiteAssociation{
					AppLinks: appLinks{
						Apps: []string{},
						Details: []appDetail{
							{
								// App ID format: <TeamID>.<BundleID>
								AppID: appID,
								Paths: []string{
									"/accept_invitation",
									"/accept_invitation/*",
								},
							},
						},
					},
					WebCredentials: webCredentials{
						Apps: []string{appID},
					},
				}

				res.Header().Set("Content-Type", "application/json")
				if err = json.NewEncoder(res).Encode(aasa); err != nil {
					logger.Error("encoding apple-app-site-association", err)
					res.WriteHeader(http.StatusInternalServerError)
				}
			})
		})
	}

	return router, nil
}
