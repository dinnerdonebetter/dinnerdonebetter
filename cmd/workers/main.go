package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prixfixeco/api_server/internal/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/messagequeue/consumers"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search/elasticsearch"
	"github.com/prixfixeco/api_server/internal/secrets"
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
	configStoreEnvVarKey = "PRIXFIXE_WORKERS_LOCAL_CONFIG_STORE_KEY"
)

func initializeLocalSecretManager(ctx context.Context, envVarKey string) secrets.SecretManager {
	logger := logging.NewNoopLogger()

	cfg := &secrets.Config{
		Provider: secrets.ProviderLocal,
		Key:      os.Getenv(envVarKey),
	}

	k, err := secrets.ProvideSecretKeeper(ctx, cfg)
	if err != nil {
		panic(err)
	}

	sm, err := secrets.ProvideSecretManager(logger, k)
	if err != nil {
		panic(err)
	}

	return sm
}

func main() {
	ctx := context.Background()

	logger := logging.ProvideLogger(logging.Config{
		Provider: logging.ProviderZerolog,
	})

	logger.Info("starting workers...")

	// find and validate our configuration filepath.
	configFilepath := os.Getenv(configFilepathEnvVar)
	if configFilepath == "" {
		log.Fatal("no config provided")
	}

	configBytes, err := os.ReadFile(configFilepath)
	if err != nil {
		logger.Fatal(err)
	}

	sm := initializeLocalSecretManager(ctx, configStoreEnvVarKey)

	var cfg *config.InstanceConfig
	if err = sm.Decrypt(ctx, string(configBytes), &cfg); err != nil || cfg == nil {
		logger.Fatal(err)
	}

	cfg.Observability.Tracing.Jaeger.ServiceName = "workers"

	flushFunc, initializeTracerErr := cfg.Observability.Tracing.Initialize(logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	// if tracing is disabled, this will be nil
	if flushFunc != nil {
		defer flushFunc()
	}

	cfg.Database.RunMigrations = false

	dataManager, err := config.ProvideDatabaseClient(ctx, logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	consumerProvider := consumers.ProvideRedisConsumerProvider(logger, string(cfg.Events.RedisConfig.QueueAddress))

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	// post-writes worker

	postWritesWorker := workers.ProvideDataChangesWorker(logger)
	postWritesConsumer, err := consumerProvider.ProviderConsumer(ctx, dataChangesTopicName, postWritesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go postWritesConsumer.Consume(nil, nil)

	// pre-writes worker

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	postWritesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preWritesWorker, err := workers.ProvideWritesWorker(ctx, logger, client, dataManager, postWritesPublisher, cfg.Search.Address, elasticsearch.NewIndexManager)
	if err != nil {
		logger.Fatal(err)
	}

	preWritesConsumer, err := consumerProvider.ProviderConsumer(ctx, preWritesTopicName, preWritesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preWritesConsumer.Consume(nil, nil)

	// pre-updates worker

	postUpdatesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preUpdatesWorker, err := workers.ProvideUpdatesWorker(ctx, logger, client, dataManager, postUpdatesPublisher, cfg.Search.Address, elasticsearch.NewIndexManager)
	if err != nil {
		logger.Fatal(err)
	}

	preUpdatesConsumer, err := consumerProvider.ProviderConsumer(ctx, preUpdatesTopicName, preUpdatesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preUpdatesConsumer.Consume(nil, nil)

	// pre-archives worker

	postArchivesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	preArchivesWorker, err := workers.ProvideArchivesWorker(ctx, logger, client, dataManager, postArchivesPublisher, cfg.Search.Address, elasticsearch.NewIndexManager)
	if err != nil {
		logger.Fatal(err)
	}

	preArchivesConsumer, err := consumerProvider.ProviderConsumer(ctx, preArchivesTopicName, preArchivesWorker.HandleMessage)
	if err != nil {
		logger.Fatal(err)
	}

	go preArchivesConsumer.Consume(nil, nil)

	everySecond := time.Tick(time.Second)
	choresWorker := workers.ProvideChoresWorker(logger, dataManager, postUpdatesPublisher)

	choresConsumer, err := consumerProvider.ProviderConsumer(ctx, choresTopicName, choresWorker.HandleMessage)
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
