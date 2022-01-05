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

// WritesWorker writes data from the pending writes topic to the database.
type WritesWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataChangesPublisher  messagequeue.Publisher
	dataManager           database.DataManager
	emailSender           email.Emailer
	customerDataCollector customerdata.Collector
}

// ProvideWritesWorker provides a WritesWorker.
func ProvideWritesWorker(
	_ context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postWritesPublisher messagequeue.Publisher,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*WritesWorker, error) {
	const name = "pre_writes"

	w := &WritesWorker{
		logger:                logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataChangesPublisher:  postWritesPublisher,
		dataManager:           dataManager,
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}

	return w, nil
}

func (w *WritesWorker) determineWriteMessageHandler(msg *types.PreWriteMessage) func(context.Context, *types.PreWriteMessage) error {
	logger := w.logger.WithValue("data_type", msg.DataType)
	logger.Debug("determining message handler for msg with type")

	funcMap := map[string]func(context.Context, *types.PreWriteMessage) error{
		string(types.MealDataType):                 w.createMeal,
		string(types.RecipeDataType):               w.createRecipe,
		string(types.RecipeStepDataType):           w.createRecipeStep,
		string(types.RecipeStepInstrumentDataType): w.createRecipeStepInstrument,
		string(types.RecipeStepIngredientDataType): w.createRecipeStepIngredient,
		string(types.RecipeStepProductDataType):    w.createRecipeStepProduct,
		string(types.MealPlanDataType):             w.createMealPlan,
		string(types.MealPlanOptionDataType):       w.createMealPlanOption,
		string(types.MealPlanOptionVoteDataType):   w.createMealPlanOptionVote,
		string(types.UserMembershipDataType):       func(context.Context, *types.PreWriteMessage) error { return nil },
	}

	f, ok := funcMap[string(msg.DataType)]
	if ok {
		return f
	}

	logger.Debug("no message handler found for msg with type")

	return nil
}

// HandleMessage handles a pending write.
func (w *WritesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.WithValue("msg", string(message))

	logger.Debug("message received")
	var msg *types.PreWriteMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger = logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read successfully")

	f := w.determineWriteMessageHandler(msg)

	logger.Debug("determined message handler for msg with type")

	if f == nil {
		return fmt.Errorf("no handler assigned to message type %q", msg.DataType)
	}

	return f(ctx, msg)
}
