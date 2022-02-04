package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/prixfixeco/api_server/internal/build/server"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	useNoOpLoggerEnvVar = "USE_NOOP_LOGGER"
	useNoEnvVars        = "USE_NO_ENV_VARS"
)

func main() {
	var (
		ctx    = context.Background()
		logger = (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()
	)

	logger.SetLevel(logging.DebugLevel)
	logger.SetRequestIDFunc(func(req *http.Request) string {
		return chimiddleware.GetReqID(req.Context())
	})

	if x, parseErr := strconv.ParseBool(os.Getenv(useNoOpLoggerEnvVar)); x && parseErr == nil {
		logger = logging.NewNoopLogger()
	}

	cfg, err := config.GetConfigFromKubernetesEnv(ctx, strings.ToLower(os.Getenv(useNoEnvVars)) == "yes")
	if err != nil {
		logger.Fatal(err)
	}

	// tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	// if initializeTracerErr != nil {
	// 	logger.Error(initializeTracerErr, "initializing tracer")
	// }

	tracerProvider := trace.NewNoopTracerProvider()

	metricsProvider, initializeMetricsErr := cfg.Observability.Metrics.ProvideUnitCounterProvider(ctx, logger)
	if initializeMetricsErr != nil {
		logger.Error(initializeMetricsErr, "initializing metrics collector")
	}

	metricsHandler, metricsHandlerErr := cfg.Observability.Metrics.ProvideMetricsHandler(logger)
	if metricsHandlerErr != nil {
		logger.Error(metricsHandlerErr, "initializing metrics handler")
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)
	ctx, initSpan := tracing.StartSpan(ctx)

	// build our server struct.
	srv, err := server.Build(ctx, logger, cfg, tracerProvider, metricsProvider, metricsHandler)
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing HTTP server: %w", err))
	}

	initSpan.End()
	cancel()

	// I slept and dreamt that life was joy.
	//   I awoke and saw that life was service.
	//   	I acted and behold, service deployed.
	srv.Serve()
}
