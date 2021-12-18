package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/pkg/types"
)

// ArchivesWorker archives data from the pending archives topic to the database.
type ArchivesWorker struct {
	logger                                  logging.Logger
	tracer                                  tracing.Tracer
	encoder                                 encoding.ClientEncoder
	postArchivesPublisher                   publishers.Publisher
	dataManager                             database.DataManager
	validInstrumentsIndexManager            search.IndexManager
	validIngredientsIndexManager            search.IndexManager
	validPreparationsIndexManager           search.IndexManager
	validIngredientPreparationsIndexManager search.IndexManager
	recipesIndexManager                     search.IndexManager
	customerDataCollector                   customerdata.Collector
}

// ProvideArchivesWorker provides an ArchivesWorker.
func ProvideArchivesWorker(
	ctx context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postArchivesPublisher publishers.Publisher,
	searchIndexProvider search.IndexManagerProvider,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*ArchivesWorker, error) {
	const name = "pre_archives"

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

	w := &ArchivesWorker{
		logger:                                  logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                                  tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:                                 encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		postArchivesPublisher:                   postArchivesPublisher,
		dataManager:                             dataManager,
		validInstrumentsIndexManager:            validInstrumentsIndexManager,
		validIngredientsIndexManager:            validIngredientsIndexManager,
		validPreparationsIndexManager:           validPreparationsIndexManager,
		validIngredientPreparationsIndexManager: validIngredientPreparationsIndexManager,
		recipesIndexManager:                     recipesIndexManager,
		customerDataCollector:                   customerDataCollector,
	}

	return w, nil
}

func (w *ArchivesWorker) determineArchiveMessageHandler(msg *types.PreArchiveMessage) func(context.Context, *types.PreArchiveMessage) error {
	funcMap := map[string]func(context.Context, *types.PreArchiveMessage) error{
		string(types.ValidIngredientDataType):            w.archiveValidIngredient,
		string(types.ValidIngredientPreparationDataType): w.archiveValidIngredientPreparation,
		string(types.MealDataType):                       w.archiveMeal,
		string(types.RecipeDataType):                     w.archiveRecipe,
		string(types.RecipeStepDataType):                 w.archiveRecipeStep,
		string(types.RecipeStepInstrumentDataType):       w.archiveRecipeStepInstrument,
		string(types.RecipeStepIngredientDataType):       w.archiveRecipeStepIngredient,
		string(types.RecipeStepProductDataType):          w.archiveRecipeStepProduct,
		string(types.MealPlanDataType):                   w.archiveMealPlan,
		string(types.MealPlanOptionDataType):             w.archiveMealPlanOption,
		string(types.MealPlanOptionVoteDataType):         w.archiveMealPlanOptionVote,
		string(types.UserMembershipDataType):             func(context.Context, *types.PreArchiveMessage) error { return nil },
		string(types.HouseholdInvitationDataType):        func(context.Context, *types.PreArchiveMessage) error { return nil },
	}

	f, ok := funcMap[string(msg.DataType)]
	if ok {
		return f
	}

	return nil
}

// HandleMessage handles a pending archive.
func (w *ArchivesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.PreArchiveMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read")

	f := w.determineArchiveMessageHandler(msg)

	if f == nil {
		return fmt.Errorf("no handler assigned to message type %q", msg.DataType)
	}

	return f(ctx, msg)
}
