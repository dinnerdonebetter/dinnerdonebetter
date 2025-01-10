package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/securecookie"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
)

const (
	serverName = "admin-frontend-server"
)

func main() {
	ctx := context.Background()

	tracerProvider := tracing.NewNoopTracerProvider()
	logger := loggingcfg.ProvideLogger(&loggingcfg.Config{
		Level:    logging.DebugLevel,
		Provider: loggingcfg.ProviderSlog,
	})
	metricProvider := metrics.NewNoopMetricsProvider()

	router, err := routingcfg.ProvideRouter(&routingcfg.Config{
		Provider: routingcfg.ProviderChi,
		ChiConfig: &chi.Config{
			ServiceName:            serverName,
			ValidDomains:           nil,
			EnableCORSForLocalhost: false,
			SilenceRouteLogging:    false,
		},
	}, logger, tracerProvider, metricProvider)
	if err != nil {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse("https://api.dinnerdonebetter.dev")
	if err != nil {
		log.Fatal(err)
	}

	cookieBuilder := securecookie.New(random.MustGenerateRawBytes(ctx, 32), random.MustGenerateRawBytes(ctx, 32))
	pageBuilder := pages.NewPageBuilder(tracerProvider, logger, router, parsedURL)

	if err = setupRoutes(router, pageBuilder, cookieBuilder); err != nil {
		log.Fatal(err)
	}

	if err = http.ListenAndServe(":8080", router.Handler()); err != nil {
		slog.Info("Error starting", "error", err)
	}
}
