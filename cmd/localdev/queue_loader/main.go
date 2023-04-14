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

	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability/keys"
	logcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	serverutils "github.com/prixfixeco/backend/internal/server/utils"
	"github.com/prixfixeco/backend/pkg/types"

	_ "go.uber.org/automaxprocs"
)

const (
	mealPlanFinalizationTopic = "meal_plan_finalizer"
	mealPlanTaskCreationTopic = "meal_plan_task_creation"

	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
)

func main() {
	ctx := context.Background()
	logger, err := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting queue loader...")

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

	cfg.Observability.Tracing.Jaeger.ServiceName = "queue_loader"

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	cfg.Database.RunMigrations = false

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	urlToUse := serverutils.DetermineServiceURL().String()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	serverutils.EnsureServerIsUp(ctx, urlToUse)
	dataManager.IsReady(ctx, 50)

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	mealPlanFinalizationTicker := time.Tick(time.Second)
	mealPlanFinalizerPublisher, err := publisherProvider.ProviderPublisher(mealPlanFinalizationTopic)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for range mealPlanFinalizationTicker {
			if err = mealPlanFinalizerPublisher.Publish(ctx, &types.ChoreMessage{ChoreType: types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType}); err != nil {
				log.Fatal(err)
			}
		}
	}()

	mealPlanTaskCreationTicker := time.Tick(time.Second)
	mealPlanTaskCreationPublisher, err := publisherProvider.ProviderPublisher(mealPlanTaskCreationTopic)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for range mealPlanTaskCreationTicker {
			if err = mealPlanTaskCreationPublisher.Publish(ctx, &types.ChoreMessage{ChoreType: types.CreateMealPlanTasksChoreType}); err != nil {
				log.Fatal(err)
			}
		}
	}()

	logger.Info("working...")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
