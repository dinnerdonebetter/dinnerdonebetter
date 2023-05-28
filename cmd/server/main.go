package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/server/http/build"
	"github.com/dinnerdonebetter/backend/internal/server/http/config"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
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
	rootCtx := context.Background()
	cfg := getConfig(rootCtx)

	// only allow initialization to take so long.
	buildCtx, cancel := context.WithTimeout(rootCtx, cfg.Server.StartupDeadline)

	// build our server struct.
	srv, err := build.Build(buildCtx, cfg)
	if err != nil {
		panic(err)
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

	cancelCtx, cancelShutdown := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancelShutdown()

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err = srv.Shutdown(cancelCtx); err != nil {
		panic(err)
	}
}
