package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prixfixeco/backend/internal/build/server"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/observability"
	_ "go.uber.org/automaxprocs"
)

const (
	configFilepathEnvVar       = "CONFIGURATION_FILEPATH"
	googleCloudIndicatorEnvVar = "RUNNING_IN_GOOGLE_CLOUD_RUN"
)

func getConfig(ctx context.Context) *config.InstanceConfig {
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

	return cfg
}

func main() {
	ctx := context.Background()
	cfg := getConfig(ctx)

	// find and validate our configuration filepath.
	logger, err := cfg.Observability.Logging.ProvideLogger(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetRequestIDFunc(func(req *http.Request) string {
		return chimiddleware.GetReqID(req.Context())
	})

	// only allow initialization to take so long.
	ctx, cancel := context.WithTimeout(ctx, cfg.Server.StartupDeadline)

	// build our server struct.
	srv, err := server.Build(ctx, logger, cfg)
	if err != nil {
		observability.AcknowledgeError(err, logger, nil, "initializing HTTP server")
		return
	}

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
		observability.AcknowledgeError(err, logger, nil, "server shutdown failed")
	}

	cancel()
}
