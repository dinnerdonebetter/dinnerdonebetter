package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()
	client := &http.Client{Timeout: 10 * time.Second}

	cfg, err := config.GetConfigFromParameterStore(true)
	if err != nil {
		logger.Fatal(err)
	}
	cfg.Database.RunMigrations = false

	tracerProvider := trace.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	postChoresPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		logger.Fatal(err)
	}

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		logger.Fatal(err)
	}

	mealPlanFinalizer := workers.ProvideMealPlanFinalizer(
		logger,
		dataManager,
		postChoresPublisher,
		emailer,
		cdp,
		tracerProvider,
	)

	everyMinute := time.Tick(time.Minute)
	for {
		select {
		case <-everyMinute:
			if err = mealPlanFinalizer.HandleMessage(context.Background(), nil); err != nil {
				observability.AcknowledgeError(err, logger, nil, "performing meal plan finalization")
			}
		}
	}
}
