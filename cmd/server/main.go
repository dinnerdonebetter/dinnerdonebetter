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
	flag "github.com/spf13/pflag"

	"github.com/prixfixeco/api_server/internal/build/server"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	useNoOpLoggerEnvVar  = "USE_NOOP_LOGGER"
	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
)

var (
	configFilepath string
)

func init() {
	flag.StringVarP(&configFilepath, "config", "c", "", "the config filepath")
}

func main() {
	flag.Parse()

	var (
		ctx    = context.Background()
		logger = logging.ProvideLogger(logging.Config{Provider: logging.ProviderZerolog})
	)

	logger.SetLevel(logging.DebugLevel)

	logger.SetRequestIDFunc(func(req *http.Request) string {
		return chimiddleware.GetReqID(req.Context())
	})

	if x, err := strconv.ParseBool(os.Getenv(useNoOpLoggerEnvVar)); x && err == nil {
		logger = logging.NewNoopLogger()
	}

	// find and validate our configuration filepath.

	var (
		cfg *config.InstanceConfig
		err error
	)

	// find and validate our configuration filepath.
	configFilepath := os.Getenv(configFilepathEnvVar)
	if configFilepath != "" {
		configBytes, err := os.ReadFile(configFilepath)
		if err != nil {
			logger.Fatal(err)
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

	flushFunc, initializeTracerErr := cfg.Observability.Tracing.Initialize(logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	// if tracing is disabled, this will be nil
	if flushFunc != nil {
		defer flushFunc()
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)
	ctx, initSpan := tracing.StartSpan(ctx)

	// build our server struct.
	srv, err := server.Build(ctx, logger, cfg)
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
