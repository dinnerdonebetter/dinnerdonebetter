package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/prixfixeco/api_server/internal/build/server"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	useNoOpLoggerEnvVar  = "USE_NOOP_LOGGER"
	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
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

	// find and validate our configuration filepath.

	var (
		cfg *config.InstanceConfig
		err error
	)

	// find and validate our configuration filepath.
	if configFilepath := os.Getenv(configFilepathEnvVar); configFilepath != "" {
		configBytes, configReadErr := os.ReadFile(configFilepath)
		if configReadErr != nil {
			logger.Fatal(configReadErr)
		}

		if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
			logger.Fatal(err)
		}
	} else {
		cfg, err = config.GetConfigFromParameterStore()
		if err != nil {
			logger.Fatal(err)
		}
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	metricsProvider, initializeMetricsErr := cfg.Observability.Metrics.ProvideUnitCounterProvider(ctx, logger)
	if initializeMetricsErr != nil {
		logger.Error(initializeTracerErr, "initializing metrics collector")
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)
	ctx, initSpan := tracing.StartSpan(ctx)

	// build our server struct.
	srv, err := server.Build(ctx, logger, cfg, tracerProvider, metricsProvider)
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
