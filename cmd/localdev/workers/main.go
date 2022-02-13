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

	"github.com/prixfixeco/api_server/internal/observability/keys"
	testutils "github.com/prixfixeco/api_server/tests/utils"

	"github.com/prixfixeco/api_server/internal/messagequeue/redis"

	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/workers"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	dataChangesTopicName = "data_changes"
	choresTopicName      = "meal_plan_finalizer"

	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
)

func main() {
	ctx := context.Background()

	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()

	logger.Info("starting workers...")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// find and validate our configuration filepath.
	configFilepath := os.Getenv(configFilepathEnvVar)
	if configFilepath == "" {
		log.Fatal("no config provided")
	}

	configBytes, err := os.ReadFile(configFilepath)
	if err != nil {
		logger.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
		logger.Fatal(err)
	}

	cfg.Observability.Tracing.Jaeger.ServiceName = "workers"

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	cfg.Database.RunMigrations = false

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		logger.Fatal(err)
	}

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		logger.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	urlToUse := testutils.DetermineServiceURL().String()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(ctx, urlToUse)
	dataManager.IsReady(ctx, 50)

	consumerProvider := redis.ProvideRedisConsumerProvider(logger, tracerProvider, cfg.Events.Consumers.RedisConfig)

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	everySecond := time.Tick(time.Second)
	choresWorker := workers.ProvideChoresWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		emailer,
		cdp,
		tracerProvider,
	)

	choresConsumer, err := consumerProvider.ProvideConsumer(ctx, choresTopicName, choresWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go choresConsumer.Consume(nil, nil)

	choresPublisher, err := publisherProvider.ProviderPublisher(choresTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		for range everySecond {
			if err = choresPublisher.Publish(ctx, &types.ChoreMessage{ChoreType: types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType}); err != nil {
				logger.Fatal(err)
			}
		}
	}()

	logger.Info("working...")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
