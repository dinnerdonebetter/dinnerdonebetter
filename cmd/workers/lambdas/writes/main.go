package main

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search/elasticsearch"
	"github.com/prixfixeco/api_server/internal/workers"
)

func buildHandler(worker *workers.WritesWorker) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for i := 0; i < len(sqsEvent.Records); i++ {
			message := sqsEvent.Records[i]
			if err := worker.HandleMessage(ctx, []byte(message.Body)); err != nil {
				return observability.PrepareError(err, nil, nil, "handling writes message")
			}
		}

		return nil
	}
}

func main() {
	ctx := context.Background()
	logger := logging.NewZerologLogger()
	client := &http.Client{Timeout: 10 * time.Second}

	cfg, err := config.GetConfigFromParameterStore()
	if err != nil {
		logger.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database)
	if err != nil {
		logger.Fatal(err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	postWritesPublisher, err := publisherProvider.ProviderPublisher("data_changes")
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

	indexManagerProvider, err := elasticsearch.NewIndexManagerProvider(logger, observability.HTTPClient(), &cfg.Search)
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
	)
	if err != nil {
		logger.Fatal(err)
	}

	lambda.Start(buildHandler(preWritesWorker))
}
