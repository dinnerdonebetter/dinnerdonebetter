package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/pkg/types"
)

// WritesWorker writes data from the pending writes topic to the database.
type WritesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	dataChangesPublisher                    publishers.Publisher
	dataManager                             database.DataManager
	validInstrumentsIndexManager            search.IndexManager
	validIngredientsIndexManager            search.IndexManager
	validPreparationsIndexManager           search.IndexManager
	validIngredientPreparationsIndexManager search.IndexManager
	recipesIndexManager                     search.IndexManager
	emailSender                             email.Emailer
	customerDataCollector                   customerdata.Collector
}

// ProvideWritesWorker provides a WritesWorker.
func ProvideWritesWorker(
	ctx context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postWritesPublisher publishers.Publisher,
	searchIndexProvider search.IndexManagerProvider,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*WritesWorker, error) {
	const name = "pre_writes"

	validInstrumentsIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_instruments", "name", "variant", "description", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid instruments search index manager: %w", err)
	}

	validIngredientsIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_ingredients", "name", "variant", "description", "warning", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients search index manager: %w", err)
	}

	validPreparationsIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_preparations", "name", "description", "icon")
	if err != nil {
		return nil, fmt.Errorf("setting up valid preparations search index manager: %w", err)
	}

	validIngredientPreparationsIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "valid_ingredient_preparations", "notes", "validPreparationID", "validIngredientID")
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient preparations search index manager: %w", err)
	}

	recipesIndexManager, err := searchIndexProvider.ProvideIndexManager(ctx, logger, "recipes", "name", "source", "description", "inspiredByRecipeID")
	if err != nil {
		return nil, fmt.Errorf("setting up recipes search index manager: %w", err)
	}

	w := &WritesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:                                 encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataChangesPublisher:                    postWritesPublisher,
		dataManager:                             dataManager,
		validInstrumentsIndexManager:            validInstrumentsIndexManager,
		validIngredientsIndexManager:            validIngredientsIndexManager,
		validPreparationsIndexManager:           validPreparationsIndexManager,
		validIngredientPreparationsIndexManager: validIngredientPreparationsIndexManager,
		recipesIndexManager:                     recipesIndexManager,
		emailSender:                             emailSender,
		customerDataCollector:                   customerDataCollector,
	}

	return w, nil
}

func (w *WritesWorker) determineWriteMessageHandler(msg *types.PreWriteMessage) func(context.Context, *types.PreWriteMessage) error {
	funcMap := map[string]func(context.Context, *types.PreWriteMessage) error{
		string(types.ValidInstrumentDataType):            w.createValidInstrument,
		string(types.ValidIngredientDataType):            w.createValidIngredient,
		string(types.ValidPreparationDataType):           w.createValidPreparation,
		string(types.ValidIngredientPreparationDataType): w.createValidIngredientPreparation,
		string(types.MealDataType):                       w.createMeal,
		string(types.RecipeDataType):                     w.createRecipe,
		string(types.RecipeStepDataType):                 w.createRecipeStep,
		string(types.RecipeStepInstrumentDataType):       w.createRecipeStepInstrument,
		string(types.RecipeStepIngredientDataType):       w.createRecipeStepIngredient,
		string(types.RecipeStepProductDataType):          w.createRecipeStepProduct,
		string(types.MealPlanDataType):                   w.createMealPlan,
		string(types.MealPlanOptionDataType):             w.createMealPlanOption,
		string(types.MealPlanOptionVoteDataType):         w.createMealPlanOptionVote,
		string(types.WebhookDataType):                    w.createWebhook,
		string(types.HouseholdInvitationDataType):        w.createHouseholdInvitation,
		string(types.UserMembershipDataType):             func(context.Context, *types.PreWriteMessage) error { return nil },
	}

	f, ok := funcMap[string(msg.DataType)]
	if ok {
		return f
	}

	return nil
}

// HandleMessage handles a pending write.
func (w *WritesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	w.logger.Debug("message received")
	var msg *types.PreWriteMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read successfully")

	f := w.determineWriteMessageHandler(msg)

	if f == nil {
		return fmt.Errorf("no handler assigned to message type %q", msg.DataType)
	}

	return f(ctx, msg)
}
