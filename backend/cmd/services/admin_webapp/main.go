package main

import (
	"context"
	"encoding/base64"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	routingcfg "github.com/dinnerdonebetter/backend/internal/routing/config"
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

	cfg := &config{
		APIServerURL:           "https://api.dinnerdonebetter.dev",
		OAuth2APIClientID:      "",
		OAuth2APIClientSecret:  "",
		APIClientCacheCapacity: 64,
		APIClientCacheTTL:      12 * time.Hour,
		Cookies: cookies.Config{
			Base64EncodedHashKey:  tempCookieSecret1,
			Base64EncodedBlockKey: tempCookieSecret2,
		},
	}

	x, err := newServer(cfg, logger, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	if err = setupRoutes(logger, tracerProvider, router, x); err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router.Handler(),
	}

	if err = s.ListenAndServe(); err != nil {
		slog.Info("Error serving", "error", err)
	}
}
