package main

import (
	"context"
	"encoding/base64"
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
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"

	"github.com/gorilla/securecookie"
)

const (
	serverName = "admin-frontend-server"

	tempCookieSecret1 = "OPAu6PFzAAvztqkBiDgF_Qw9RUP2Lnng9aADq0EQeUk"
	tempCookieSecret2 = "5KRnutGaUGste3esRtl970KaFmR18EiUnhaeQ-6mYR4"
)

func mustBase64Decode(s string) []byte {
	val, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return val
}

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
	tracer := tracing.NewTracer(tracerProvider.Tracer("admin_webapp"))

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

	cookieBuilder := securecookie.New(mustBase64Decode(tempCookieSecret1), mustBase64Decode(tempCookieSecret2))
	pageBuilder := pages.NewPageBuilder(tracerProvider, logger, parsedURL)

	if err = setupRoutes(logger, tracer, router, pageBuilder, cookieBuilder); err != nil {
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
