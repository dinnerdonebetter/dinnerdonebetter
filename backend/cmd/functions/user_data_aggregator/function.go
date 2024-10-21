package webhookexecutor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("AggregateUserData", AggregateUserData)
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// AggregateUserData handles a user data aggregation request.
func AggregateUserData(ctx context.Context, e event.Event) error {
	if strings.TrimSpace(strings.ToLower(os.Getenv("CEASE_OPERATION"))) == "true" {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	logger := (&loggingcfg.Config{Level: logging.DebugLevel, Provider: loggingcfg.ProviderSlog}).ProvideLogger()

	envCfg := email.GetConfigForEnvironment(os.Getenv("DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, logger, nil, "getting environment config")
	}

	logger.Info("getting config")

	cfg, err := config.GetSearchDataIndexerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search_indexer_cloud_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger.Info("connecting to database")

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	var userDataCollectionRequest *types.UserDataAggregationRequest
	if err = json.Unmarshal(msg.Message.Data, &userDataCollectionRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	logger.Info("loaded payload, aggregating data")

	collection, err := dataManager.AggregateUserData(ctx, userDataCollectionRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "collecting user data")
	}
	collection.ReportID = userDataCollectionRequest.ReportID

	// TODO: save this in the cloud somewhere
	logger.Info("compiled payload, would save")

	return nil
}
