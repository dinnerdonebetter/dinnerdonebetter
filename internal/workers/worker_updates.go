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
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/pkg/types"
)

// UpdatesWorker updates data from the pending updates topic to the database.
type UpdatesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	postUpdatesPublisher                    messagequeue.Publisher
	dataManager                             database.DataManager
	validInstrumentsIndexManager            search.IndexManager
	validIngredientsIndexManager            search.IndexManager
	validPreparationsIndexManager           search.IndexManager
	validIngredientPreparationsIndexManager search.IndexManager
	recipesIndexManager                     search.IndexManager
	emailSender                             email.Emailer
	customerDataCollector                   customerdata.Collector
}

// ProvideUpdatesWorker provides a UpdatesWorker.
func ProvideUpdatesWorker(
	ctx context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	searchIndexProvider search.IndexManagerProvider,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*UpdatesWorker, error) {
	const name = "pre_updates"

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

	w := &UpdatesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:                                 encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		postUpdatesPublisher:                    postUpdatesPublisher,
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

func (w *UpdatesWorker) determineUpdateMessageHandler(msg *types.PreUpdateMessage) func(context.Context, *types.PreUpdateMessage) error {
	funcMap := map[string]func(context.Context, *types.PreUpdateMessage) error{
		string(types.ValidInstrumentDataType):            w.updateValidInstrument,
		string(types.ValidIngredientDataType):            w.updateValidIngredient,
		string(types.ValidPreparationDataType):           w.updateValidPreparation,
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
