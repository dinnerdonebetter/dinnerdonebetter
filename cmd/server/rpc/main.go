package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/proto"
	htttp "github.com/dinnerdonebetter/backend/internal/server/http"
	"github.com/dinnerdonebetter/backend/internal/server/rpc/build"

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
	service, err := build.Build(buildCtx, cfg)
	if err != nil {
		panic(err)
	}

	cancel()

	s := htttp.ProvideStdLibHTTPServer(8080)
	s.Handler = proto.NewDinnerDoneBetterServer(service)
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}
