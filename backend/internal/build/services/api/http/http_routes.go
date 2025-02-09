package api

import (
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	routingcfg "github.com/dinnerdonebetter/backend/internal/lib/routing/config"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func ProvideAPIRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	dataManager database.DataManager,
	authService types.AuthDataService,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	tracer := tracing.NewTracer(tracerProvider.Tracer("http_server"))

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

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.Get("/authorize", authService.AuthorizeHandler)
		userRouter.Post("/token", authService.TokenHandler)
	})

	router.Route("/auth", func(authRouter routing.Router) {
		providerRouteParam := fmt.Sprintf("/{%s}", authservice.AuthProviderParamKey)
		authRouter.Get(providerRouteParam, authService.SSOLoginHandler)
		authRouter.Get(path.Join(providerRouteParam, "callback"), authService.SSOLoginCallbackHandler)
	})

	return router, nil
}
