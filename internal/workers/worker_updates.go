package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// UpdatesWorker updates data from the pending updates topic to the database.
type UpdatesWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataChangesPublisher  messagequeue.Publisher
	dataManager           database.DataManager
	emailSender           email.Emailer
	customerDataCollector customerdata.Collector
}

// ProvideUpdatesWorker provides a UpdatesWorker.
func ProvideUpdatesWorker(
	_ context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*UpdatesWorker, error) {
	const name = "pre_updates"

	w := &UpdatesWorker{
		logger:                logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataChangesPublisher:  postUpdatesPublisher,
		dataManager:           dataManager,
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}

	return w, nil
}

func (w *UpdatesWorker) determineUpdateMessageHandler(msg *types.PreUpdateMessage) func(context.Context, *types.PreUpdateMessage) error {
	funcMap := map[string]func(context.Context, *types.PreUpdateMessage) error{
		string(types.ValidIngredientPreparationDataType): w.updateValidIngredientPreparation,
		string(types.RecipeDataType):                     w.updateRecipe,
		string(types.RecipeStepDataType):                 w.updateRecipeStep,
		string(types.RecipeStepInstrumentDataType):       w.updateRecipeStepInstrument,
		string(types.RecipeStepIngredientDataType):       w.updateRecipeStepIngredient,
		string(types.RecipeStepProductDataType):          w.updateRecipeStepProduct,
		string(types.MealPlanDataType):                   w.updateMealPlan,
		string(types.MealPlanOptionDataType):             w.updateMealPlanOption,
		string(types.MealPlanOptionVoteDataType):         w.updateMealPlanOptionVote,
		string(types.UserMembershipDataType):             func(context.Context, *types.PreUpdateMessage) error { return nil },
		string(types.WebhookDataType):                    func(context.Context, *types.PreUpdateMessage) error { return nil },
		string(types.HouseholdInvitationDataType):        func(context.Context, *types.PreUpdateMessage) error { return nil },
	}

	f, ok := funcMap[string(msg.DataType)]
	if ok {
		return f
	}

	return nil
}

// HandleMessage handles a pending update.
func (w *UpdatesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.PreUpdateMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read")

	f := w.determineUpdateMessageHandler(msg)

	if f == nil {
		return fmt.Errorf("no handler assigned to message type %q", msg.DataType)
	}

	return f(ctx, msg)
}
