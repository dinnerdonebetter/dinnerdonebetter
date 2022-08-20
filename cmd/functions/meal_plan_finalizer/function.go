package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	"log"

	"github.com/prixfixeco/api_server/internal/config"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/database/postgres"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go.opentelemetry.io/otel"
)

const (
	dataChangesTopicName = "data_changes"
)

func finalizeMealPlans(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	dataChangesPublisher messagequeue.Publisher,
) (int, error) {
	_, span := tracer.StartSpan(ctx)
	defer span.End()

	mealPlans, fetchMealPlansErr := dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if fetchMealPlansErr != nil {
		log.Fatal(fetchMealPlansErr)
	}

	var changedCount int
	for _, mealPlan := range mealPlans {
		changed, err := dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "finalizing meal plan")
			continue
		}

		if changed {
			changedCount++

			if dataChangePublishErr := dataChangesPublisher.Publish(ctx, &types.DataChangeMessage{
				DataType:                  types.MealPlanDataType,
				EventType:                 types.MealPlanFinalizedCustomerEventType,
				MealPlan:                  mealPlan,
				MealPlanID:                mealPlan.ID,
				Context:                   map[string]string{},
				AttributableToHouseholdID: mealPlan.BelongsToHousehold,
			}); dataChangePublishErr != nil {
				observability.AcknowledgeError(dataChangePublishErr, logger, span, "publishing data change message")
			}
		}
	}

	return changedCount, nil
}

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// FinalizeMealPlans is our cloud function entrypoint.
func FinalizeMealPlans(ctx context.Context, m PubSubMessage) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider := tracing.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	tracer := tracing.NewTracer(tracerProvider.Tracer("meal_plan_finalizer"))

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		return fmt.Errorf("error setting up database client: %w", err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		log.Fatal(err)
	}

	changedCount, mealPlanFinalizationErr := finalizeMealPlans(ctx, logger, tracer, dataManager, dataChangesPublisher)
	if mealPlanFinalizationErr != nil {
		return fmt.Errorf("finalizing meal plans: %w", mealPlanFinalizationErr)
	}

	if closeErr := dataManager.DB().Close(); closeErr != nil {
		observability.AcknowledgeError(closeErr, logger, nil, "closing database connection")
	}

	if changedCount > 0 {
		logger.WithValue("count", changedCount).Info("finalized meal plans")
	}

	return nil
}
