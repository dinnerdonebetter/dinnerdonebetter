package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "go.uber.org/automaxprocs"

	"github.com/prixfixeco/backend/internal/build/server"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

const (
	useNoOpLoggerEnvVar        = "USE_NOOP_LOGGER"
	configFilepathEnvVar       = "CONFIGURATION_FILEPATH"
	googleCloudIndicatorEnvVar = "RUNNING_IN_GOOGLE_CLOUD_RUN"
)

func main() {
	ctx := context.Background()

	var cfg *config.InstanceConfig
	if os.Getenv(googleCloudIndicatorEnvVar) != "" {
		client, secretManagerCreationErr := secretmanager.NewClient(ctx)
		if secretManagerCreationErr != nil {
			log.Fatal(secretManagerCreationErr)
		}

		c, cfgHydrateErr := config.GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx, client)
		if cfgHydrateErr != nil {
			log.Fatal(cfgHydrateErr)
		}
		cfg = c
	} else if configFilepath := os.Getenv(configFilepathEnvVar); configFilepath != "" {
		configBytes, configReadErr := os.ReadFile(configFilepath)
		if configReadErr != nil {
			log.Fatal(configReadErr)
		}

		if err := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("no config provided")
	}

	// find and validate our configuration filepath.
	logger, err := cfg.Observability.Logging.ProvideLogger(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetRequestIDFunc(func(req *http.Request) string {
		return chimiddleware.GetReqID(req.Context())
	})

	if x, parseErr := strconv.ParseBool(os.Getenv(useNoOpLoggerEnvVar)); x && parseErr == nil {
		logger = logging.NewNoopLogger()
	}

	logger.SetLevel(cfg.Observability.Logging.Level)

	// should make wire do these someday
	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	defer func() {
		if flushErr := tracerProvider.ForceFlush(ctx); flushErr != nil {
			logger.Error(flushErr, "flushing traces")
		}
	}()

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

	client := tracing.BuildTracedHTTPClient()
	// should make wire do these someday
	emailer, emailerSetupErr := cfg.Email.ProvideEmailer(logger, client)
	if emailerSetupErr != nil {
		logger.Error(emailerSetupErr, "initializing metrics handler")
	}

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)
	ctx, span := tracing.StartSpan(ctx)

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "initializing database client")
		return
	}

	// build our server struct.
	// NOTE: we should, at some point, return a function that we can defer that will end any database connections and such
	srv, err := server.Build(ctx, logger, cfg, tracerProvider, metricsProvider, metricsHandler, dataManager, emailer)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "initializing HTTP server")
		return
	}

	span.End()
	cancel()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// Run server
	go srv.Serve()

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
	}()

	_, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err = srv.Shutdown(ctx); err != nil {
		observability.AcknowledgeError(err, logger, span, "server shutdown failed")
	}

	if err = dataManager.DB().Close(); err != nil {
		observability.AcknowledgeError(err, logger, span, "closing database connection")
	}

	cancel()
}
