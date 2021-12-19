package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/search/elasticsearch"
	"github.com/prixfixeco/api_server/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func buildHandler(logger logging.Logger, worker *workers.WritesWorker) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	logger.Info("starting lambda")
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		logger = logger.WithValue("event_count", len(sqsEvent.Records))
		logger.Info("looping through records")
		for i := 0; i < len(sqsEvent.Records); i++ {
			logger.WithValue("event_count", len(sqsEvent.Records)).Info("record loop")
			message := sqsEvent.Records[i]
			if err := worker.HandleMessage(ctx, []byte(message.Body)); err != nil {
				return observability.PrepareError(err, logger, nil, "handling writes message")
			}
		}

		return nil
	}
}

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()
	client := &http.Client{Timeout: 10 * time.Second}

	cfg, err := config.GetConfigFromParameterStore(true)
	if err != nil {
		logger.Fatal(err)
	}
	cfg.Database.RunMigrations = false

	tracerProvider, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		if shutdownErr := tracerProvider.Shutdown(ctx); shutdownErr != nil {
			fmt.Printf("error shutting down tracer provider: %v", shutdownErr)
		}
	}(ctx)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(xray.Propagator{})

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	postWritesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
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

	indexManagerProvider, err := elasticsearch.NewIndexManagerProvider(ctx, logger, &cfg.Search, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	preWritesWorker, err := workers.ProvideWritesWorker(
		ctx,
		logger,
		dataManager,
		postWritesPublisher,
		indexManagerProvider,
		emailer,
		cdp,
		tracerProvider,
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("starting lambda")

	lambda.Start(buildHandler(logger, preWritesWorker))
}
