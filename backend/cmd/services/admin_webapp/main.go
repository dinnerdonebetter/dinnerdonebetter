package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/pages"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"

	"github.com/gorilla/securecookie"
)

const (
	serverName = "admin-frontend-server"
)

func main() {
	ctx := context.Background()

	logger, err := loggingcfg.ProvideLogger(ctx, &loggingcfg.Config{
		Provider: loggingcfg.ProviderSlog,
		Level:    logging.InfoLevel,
	})
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := tracing.NewNoopTracerProvider()
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
	pageBuilder := pages.NewPageBuilder(tracerProvider, logger, parsedURL)

	if err = setupRoutes(router, pageBuilder, cookieBuilder); err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router.Handler(),
	}

	if err = server.ListenAndServe(); err != nil {
		slog.Info("Error serving", "error", err)
	}
}
