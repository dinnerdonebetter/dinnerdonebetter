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
	preWritesTopicName   = "pre_writes"
	dataChangesTopicName = "data_changes"
	preUpdatesTopicName  = "pre_updates"
	preArchivesTopicName = "pre_archives"
	choresTopicName      = "chores"

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

	dataManager.IsReady(ctx, 50)

	consumerProvider := redis.ProvideRedisConsumerProvider(logger, tracerProvider, cfg.Events.RedisConfig)

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	// post-writes worker

	postWritesWorker := workers.ProvideDataChangesWorker(
		logger,
		emailer,
		cdp,
		tracerProvider,
	)
	postWritesConsumer, err := consumerProvider.ProvideConsumer(ctx, dataChangesTopicName, postWritesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go postWritesConsumer.Consume(nil, nil)

	// pre-writes worker

	postWritesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preWritesWorker, err := workers.ProvideWritesWorker(
		ctx,
		logger,
		dataManager,
		postWritesPublisher,
		emailer,
		cdp,
		tracerProvider,
	)
	if err != nil {
		logger.Fatal(err)
	}

	preWritesConsumer, err := consumerProvider.ProvideConsumer(ctx, preWritesTopicName, preWritesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preWritesConsumer.Consume(nil, nil)

	// pre-updates worker

	postUpdatesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preUpdatesWorker, err := workers.ProvideUpdatesWorker(
		ctx,
		logger,
		dataManager,
		postUpdatesPublisher,
		emailer,
		cdp,
		tracerProvider,
	)
	if err != nil {
		logger.Fatal(err)
	}

	preUpdatesConsumer, err := consumerProvider.ProvideConsumer(ctx, preUpdatesTopicName, preUpdatesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preUpdatesConsumer.Consume(nil, nil)

	// pre-archives worker

	postArchivesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preArchivesWorker, err := workers.ProvideArchivesWorker(
		ctx,
		logger,
		dataManager,
		postArchivesPublisher,
		cdp,
		tracerProvider,
	)
	if err != nil {
		logger.Fatal(err)
	}

	preArchivesConsumer, err := consumerProvider.ProvideConsumer(ctx, preArchivesTopicName, preArchivesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preArchivesConsumer.Consume(nil, nil)

	everySecond := time.Tick(time.Second)
	choresWorker := workers.ProvideChoresWorker(
		logger,
		dataManager,
		postUpdatesPublisher,
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
			mealPlans, err := dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
			if err != nil {
				logger.Fatal(err)
			}

			for _, mealPlan := range mealPlans {
				if err = choresPublisher.Publish(ctx, &types.ChoreMessage{
					ChoreType:                 types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType,
					MealPlanID:                mealPlan.ID,
					AttributableToHouseholdID: mealPlan.BelongsToHousehold,
				}); err != nil {
					logger.Fatal(err)
				}
			}
		}
	}()

	logger.Info("working...")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
