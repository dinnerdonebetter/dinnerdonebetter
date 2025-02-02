// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package adminwebapp

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/lib/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http"
	"github.com/dinnerdonebetter/backend/internal/services/frontend/admin"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, cfg *config.AdminWebappConfig) (http.Server, error) {
	httpConfig := cfg.HTTPServer
	observabilityConfig := &cfg.Observability
	loggingcfgConfig := &observabilityConfig.Logging
	logger, err := loggingcfg.ProvideLogger(ctx, loggingcfgConfig)
	if err != nil {
		return nil, err
	}
	routingcfgConfig := cfg.Routing
	tracingcfgConfig := &observabilityConfig.Tracing
	tracerProvider, err := tracingcfg.ProvideTracerProvider(ctx, tracingcfgConfig, logger)
	if err != nil {
		return nil, err
	}
	metricscfgConfig := &observabilityConfig.Metrics
	provider, err := metricscfg.ProvideMetricsProvider(ctx, logger, metricscfgConfig)
	if err != nil {
		return nil, err
	}
	webappServer, err := admin.NewServer(cfg, logger, tracerProvider)
	if err != nil {
		return nil, err
	}
	router, err := ProvideAdminWebappRouter(routingcfgConfig, logger, tracerProvider, provider, webappServer)
	if err != nil {
		return nil, err
	}
	server, err := http.ProvideHTTPServer(httpConfig, logger, router, tracerProvider)
	if err != nil {
		return nil, err
	}
	return server, nil
}
