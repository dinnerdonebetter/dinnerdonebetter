package api

import (
	"net/http"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth"
	paymentswebhook "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/payments/http"

	"github.com/verygoodsoftwarenotvirus/platform/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/routing"
	routingcfg "github.com/verygoodsoftwarenotvirus/platform/routing/config"
	"github.com/verygoodsoftwarenotvirus/platform/version"
)

func ProvideAPIRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	authService auth.AuthDataService,
	paymentsWebhookHandler *paymentswebhook.WebhookHandler,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	encoder := encoding.ProvideServerEncoderDecoder(logger, tracerProvider, encoding.ContentTypeJSON)

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

		metaRouter.Get("/version", func(res http.ResponseWriter, req *http.Request) {
			encoder.EncodeResponseWithStatus(req.Context(), res, version.Get(), http.StatusOK)
		})
	})

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.Get("/authorize", authService.AuthorizeHandler)
		userRouter.Post("/token", authService.TokenHandler)
		userRouter.Post("/revoke", authService.RevokeHandler)
	})

	router.Route("/api/payments/webhooks", func(paymentsRouter routing.Router) {
		paymentsRouter.Post("/{provider}", paymentsWebhookHandler.Handle)
	})

	return router, nil
}
