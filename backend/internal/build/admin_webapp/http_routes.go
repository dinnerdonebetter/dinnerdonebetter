package adminwebapp

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
	adminsvc "github.com/dinnerdonebetter/backend/internal/services/frontend/admin"

	ghttp "maragu.dev/gomponents/http"
)

func buildURLVarChunk(key string) string {
	return fmt.Sprintf("/{%s}", key)
}

func ProvideAdminWebappRouter(
	routingConfig routingcfg.Config,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	webappServer *adminsvc.WebappServer,
) (routing.Router, error) {
	router, err := routingConfig.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	router.Route("/_meta_", func(metaRouter routing.Router) {
		// Expose a readiness check on /ready
		metaRouter.Get("/live", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})

		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	})
	authedRouter := router.WithMiddleware(webappServer.AuthMiddleware)

	router.Get("/", ghttp.Adapt(webappServer.RenderHome))
	router.Get("/about", ghttp.Adapt(webappServer.RenderAbout))

	router.Get("/login", ghttp.Adapt(webappServer.RenderLoginPage))
	router.Post("/login/submit", ghttp.Adapt(webappServer.HandleLoginSubmission))

	authedRouter.Get("/users", ghttp.Adapt(webappServer.RenderUsersPage))
	authedRouter.Get("/valid_ingredients", ghttp.Adapt(webappServer.RenderValidIngredientsPage))
	// This is just here for the moment so that buildURLVarChunk gets used and the linter doesn't yell at me.
	authedRouter.Get("/valid_ingredients/"+buildURLVarChunk(adminsvc.ValidIngredientIDURLParamKey), ghttp.Adapt(webappServer.RenderValidIngredientsPage))
	authedRouter.Get("/valid_ingredients/new", ghttp.Adapt(webappServer.RenderValidIngredientCreationForm))
	authedRouter.Post("/valid_ingredients/new/submit", ghttp.Adapt(webappServer.HandleValidIngredientSubmission))

	return router, nil
}
