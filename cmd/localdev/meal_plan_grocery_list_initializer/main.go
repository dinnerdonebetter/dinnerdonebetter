package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/postgres"
	"github.com/prixfixeco/api_server/internal/features/grocerylistpreparation"
	"github.com/prixfixeco/api_server/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/messagequeue/redis"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/workers"
	"github.com/prixfixeco/api_server/pkg/utils"
)

const (
	dataChangesTopicName         = "data_changes"
	mealPlanGroceryListInitTopic = "meal_plan_grocery_list_initialization"

	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
)

func main() {
	ctx := context.Background()
	logger, err := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting meal plan grocery list worker...")

	// find and validate our configuration filepath.
	configFilepath := os.Getenv(configFilepathEnvVar)
	if configFilepath == "" {
		log.Fatal("no config provided")
	}

	configBytes, err := os.ReadFile(configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
		log.Fatal(err)
	}

	cfg.Observability.Tracing.Jaeger.ServiceName = "meal_plan_task_creation_workers"

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	cfg.Database.RunMigrations = false

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		log.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	urlToUse := utils.DetermineServiceURL().String()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	utils.EnsureServerIsUp(ctx, urlToUse)
	dataManager.IsReady(ctx, 50)

	consumerProvider := redis.ProvideRedisConsumerProvider(logger, tracerProvider, cfg.Events.Consumers.RedisConfig)

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		log.Fatal(err)
	}

	mealPlanGroceryListInitializationWorker := workers.ProvideMealPlanGroceryListInitializer(
		logger,
		dataManager,
		&recipeanalysis.MockRecipeAnalyzer{},
		dataChangesPublisher,
		cdp,
		tracerProvider,
		grocerylistpreparation.NewGroceryListCreator(logger, tracerProvider),
	)

	mealPlanGroceryListInitializationConsumer, err := consumerProvider.ProvideConsumer(ctx, mealPlanGroceryListInitTopic, mealPlanGroceryListInitializationWorker.HandleMessage)
	if err != nil {
		log.Fatal(err)
	}

	go mealPlanGroceryListInitializationConsumer.Consume(nil, nil)

	logger.Info("working...")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
