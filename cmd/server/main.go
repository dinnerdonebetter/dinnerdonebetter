package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prixfixeco/api_server/internal/observability"

	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/prixfixeco/api_server/internal/build/server"
	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	useNoOpLoggerEnvVar        = "USE_NOOP_LOGGER"
	configFilepathEnvVar       = "CONFIGURATION_FILEPATH"
	googleCloudIndicatorEnvVar = "RUNNING_IN_GOOGLE_CLOUD_RUN"
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

	if os.Getenv(googleCloudIndicatorEnvVar) != "" {
		cfg, err = config.GetConfigFromGoogleCloudRunEnvironment(ctx)
		if err != nil {
			logger.Fatal(err)
		}
	} else if configFilepath := os.Getenv(configFilepathEnvVar); configFilepath != "" {
		configBytes, configReadErr := os.ReadFile(configFilepath)
		if configReadErr != nil {
			logger.Fatal(configReadErr)
		}

		if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
			logger.Fatal(err)
		}
	} else {
		log.Fatal("no config provided")
	}

	// should make wire do these someday
	tracerProvider := tracing.NewNoopTracerProvider()

	// should make wire do these someday
	metricsProvider, initializeMetricsErr := cfg.Observability.Metrics.ProvideUnitCounterProvider(ctx, logger)
	if initializeMetricsErr != nil {
		logger.Error(initializeMetricsErr, "initializing metrics collector")
	}

	// should make wire do these someday
	metricsHandler, metricsHandlerErr := cfg.Observability.Metrics.ProvideMetricsHandler(logger)
	if metricsHandlerErr != nil {
		logger.Error(metricsHandlerErr, "initializing metrics handler")
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)
	ctx, initSpan := tracing.StartSpan(ctx)

	dbConfig := &cfg.Database
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracerProvider)
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing database client: %w", err))
	}

	// build our server struct.
	// NOTE: we should, at some point, return a function that we can defer that will end any database connections and such
	srv, err := server.Build(ctx, logger, cfg, tracerProvider, metricsProvider, metricsHandler, dataManager)
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing HTTP server: %w", err))
	}

	initSpan.End()
	cancel()

	// I slept and dreamt that life was joy.
	//   I awoke and saw that life was service.
	//   	I acted and behold, service deployed.
	srv.Serve()

	// Run server
	go srv.Serve()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
	}()

	_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if dbCloseErr := dataManager.DB().Close(); dbCloseErr != nil {
		observability.AcknowledgeError(err, logger, nil, "closing database connection")
	}

	cancel()
}
