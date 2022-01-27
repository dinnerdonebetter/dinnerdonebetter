package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/email"

	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"

	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

func buildHandler(tracer tracing.Tracer, logger logging.Logger, _ customerdata.Collector, _ email.Emailer) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		_, span := tracer.StartSpan(ctx)
		defer span.End()

		logger = logger.WithValue("message_count", len(sqsEvent.Records))
		logger.Info("looping through records")

		for i := 0; i < len(sqsEvent.Records); i++ {
			message := sqsEvent.Records[i]

			var dcm *types.DataChangeMessage
			if err := json.Unmarshal([]byte(message.Body), &dcm); err != nil {
				logger.Error(err, "unmarshalling data change message")
			}
		}

		logger.Info("all messages handled")

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

	cfg.Events.Publishers.Provider = msgconfig.ProviderRedis

	tracerProvider, err := cfg.Observability.Tracing.Initialize(ctx, logger)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}
	tracer := tracerProvider.Tracer("data_changes_worker")

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(xray.Propagator{})

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		logger.Fatal(err)
	}

	cdc, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		logger.Fatal(err)
	}

	lambda.Start(buildHandler(tracing.NewTracer(tracer), logger, cdc, emailer))
}
