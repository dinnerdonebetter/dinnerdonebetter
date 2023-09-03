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

	"github.com/dinnerdonebetter/backend/internal/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "go.uber.org/automaxprocs"
)

const (
	mealPlanFinalizationTopic = "meal_plan_finalizer"
	mealPlanTaskCreationTopic = "meal_plan_task_creation"

	configFilepathEnvVar = "CONFIGURATION_FILEPATH"
)

func main() {
	ctx := context.Background()

	logger := (&logcfg.Config{Level: logging.DebugLevel, Provider: logcfg.ProviderSlog}).ProvideLogger()
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

	cfg.Observability.Tracing.Otel.ServiceName = "queue_loader"

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	mealPlanFinalizationTicker := time.Tick(time.Second)
	mealPlanFinalizerPublisher, err := publisherProvider.ProvidePublisher(mealPlanFinalizationTopic)
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
	mealPlanTaskCreationPublisher, err := publisherProvider.ProvidePublisher(mealPlanTaskCreationTopic)
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
