package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
)

type (
	contextKey string
)

const (
	serverName = "admin-frontend-Server"
	cookieName = "dinner-done-better-admin-webapp"

	apiClientContextKey contextKey = "api_client"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfigFromEnvironment[config.AdminWebappConfig]()
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	router, err := cfg.Routing.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		log.Fatal(err)
	}

	webappCfg := &adminWebappCfg{
		APIServerURL:           "https://api.dinnerdonebetter.dev",
		OAuth2APIClientID:      "9819637062b9bbd3c1997cd3b6a264d4",
		OAuth2APIClientSecret:  "0299fececf3f0be3af94adc9a98b2b0b",
		APIClientCacheCapacity: 64,
		APIClientCacheTTL:      12 * time.Hour,
		Cookies: cookies.Config{
			Base64EncodedHashKey:  "OPAu6PFzAAvztqkBiDgF_Qw9RUP2Lnng9aADq0EQeUk",
			Base64EncodedBlockKey: "5KRnutGaUGste3esRtl970KaFmR18EiUnhaeQ-6mYR4",
		},
	}

	x, err := NewServer(webappCfg, logger, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	if err = x.setupRoutes(router); err != nil {
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
