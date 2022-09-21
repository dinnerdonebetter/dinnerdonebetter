package mealplanfinalizerfunction

import (
	"context"
	"fmt"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/workers"
	testutils "github.com/prixfixeco/api_server/tests/utils"
	"log"
	"net/http"
	"time"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
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
func FinalizeMealPlans(ctx context.Context, _ PubSubMessage) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetMealPlanFinalizerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}

	cfg.Database.RunMigrations = false
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		log.Fatal(err)
	}

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		log.Fatal(err)
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	urlToUse := testutils.DetermineServiceURL().String()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(ctx, urlToUse)
	dataManager.IsReady(ctx, 50)

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		log.Fatal(err)
	}

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		log.Fatal(err)
	}

	mealPlanFinalizationWorker := workers.ProvideMealPlanFinalizationWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		emailer,
		cdp,
		tracerProvider,
	)

	changedCount, mealPlanFinalizationErr := mealPlanFinalizationWorker.FinalizeExpiredMealPlans(ctx, nil)
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
