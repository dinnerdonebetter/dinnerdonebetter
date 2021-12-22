package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/messagequeue/redis"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

func buildHandler(tracer tracing.Tracer, logger logging.Logger, notificationQueue messagequeue.Publisher) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		ctx, span := tracer.StartSpan(ctx)
		defer span.End()

		logger = logger.WithValue("message_count", len(sqsEvent.Records))
		logger.Info("looping through records")

		for i := 0; i < len(sqsEvent.Records); i++ {
			message := sqsEvent.Records[i]

			var dcm *types.DataChangeMessage
			if err := json.Unmarshal([]byte(message.Body), &dcm); err != nil {
				logger.Error(err, "unmarshalling data change message")
			}

			if err := notificationQueue.Publish(ctx, message); err != nil {
				return observability.PrepareError(err, logger, span, "publishing message to notification queue")
			}

			logger.Info("published message to pubsub")
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

	tracerProvider, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}
	tracer := tracerProvider.Tracer("data_changes_worker")

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

	_, _ = emailer, cdp

	publisherProvider := redis.ProvideRedisPublisherProvider(logger, tracerProvider, cfg.Events.Publishers.RedisConfig)
	publisher, err := publisherProvider.ProviderPublisher("data_changes")

	lambda.Start(buildHandler(tracing.NewTracer(tracer), logger, publisher))
}
