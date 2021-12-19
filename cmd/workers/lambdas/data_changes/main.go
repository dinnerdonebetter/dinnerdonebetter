package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func buildHandler(worker *workers.DataChangesWorker) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for i := 0; i < len(sqsEvent.Records); i++ {
			message := sqsEvent.Records[i]
			if err := worker.HandleMessage(ctx, []byte(message.Body)); err != nil {
				return observability.PrepareError(err, nil, nil, "handling archives message")
			}
		}

		return nil
	}
}

func main() {
	logger := zerolog.NewZerologLogger()
	client := &http.Client{Timeout: 10 * time.Second}
	ctx := context.Background()

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

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		logger.Fatal(err)
	}

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		logger.Fatal(err)
	}

	dataChangesWorker := workers.ProvideDataChangesWorker(logger, emailer, cdp, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	lambda.Start(buildHandler(dataChangesWorker))
}
