package outboundemailerfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("IndexDataForSearch", IndexDataForSearch)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// IndexDataForSearch handles a data change.
func IndexDataForSearch(ctx context.Context, e event.Event) error {
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	if strings.TrimSpace(strings.ToLower(os.Getenv("CEASE_OPERATION"))) == "true" {
		logger.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	envCfg := email.GetConfigForEnvironment(os.Getenv("DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, logger, nil, "getting environment config")
	}

	cfg, err := config.GetOutboundEmailerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	ctx, span := tracing.NewTracer(tracerProvider.Tracer("outbound_emailer_job")).StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "error setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	// manual db timeout until I find out what's wrong
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, &cfg.Database, tracerProvider)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	var searchIndexRequest search.IndexRequest
	if err = json.Unmarshal(msg.Message.Data, &searchIndexRequest); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	switch searchIndexRequest.IndexType {
	case search.IndexTypeRecipes:
		var im search.IndexManager[search.RecipeSearchSubset]
		im, err = searchcfg.ProvideIndexManager[search.RecipeSearchSubset](ctx, logger, tracerProvider, &cfg.Search, search.IndexTypeRecipes)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "initializing index manager")
		}

		var recipes *types.QueryFilteredResult[types.Recipe]
		recipes, err = dataManager.GetRecipes(ctx, types.DefaultQueryFilter())
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting recipes")
		}

		for _, x := range recipes.Data {
			var recipe *types.Recipe
			recipe, err = dataManager.GetRecipe(ctx, x.ID)
			if err != nil {
				return observability.PrepareAndLogError(err, logger, span, "getting recipe")
			}

			toBeIndexed := search.SubsetFromRecipe(recipe)
			if err = im.Index(ctx, recipe.ID, toBeIndexed); err != nil {
				return observability.PrepareAndLogError(err, logger, span, "indexing recipe")
			}
		}
	}

	return nil
}
