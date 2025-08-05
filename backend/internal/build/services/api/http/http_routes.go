package api

import (
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	routingcfg "github.com/dinnerdonebetter/backend/internal/platform/routing/config"
	authservice "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
)

func buildURLVarChunk(key string) string {
	return fmt.Sprintf("/{%s}", key)
}

func ProvideAPIRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	authService identity.AuthDataService,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	router.Route("/_ops_", func(metaRouter routing.Router) {
		// Expose a readiness check on /ready
		metaRouter.Get("/live", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			// TODO: check readiness here lol
			res.WriteHeader(http.StatusOK)
		})
	})

	authenticatedRouter := router.WithMiddleware(authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", authService.StatusHandler)

	router.Route("/oauth2", func(userRouter routing.Router) {
		//userRouter.Get("/authorize", authService.AuthorizeHandler)
		//userRouter.Post("/token", authService.TokenHandler)
	})

	router.Route("/auth", func(authRouter routing.Router) {
		providerRouteParam := buildURLVarChunk(authservice.AuthProviderParamKey)
		authRouter.Get(providerRouteParam, authService.SSOLoginHandler)
		authRouter.Get(path.Join(providerRouteParam, "callback"), authService.SSOLoginCallbackHandler)
	})

	return router, nil
}
