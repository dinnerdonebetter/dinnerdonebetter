package managers

import (
	"context"
	"errors"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

const (
	recipeManagerName = "recipe_manager"
)

type (
	RecipeManager interface {
		ListRecipes(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Recipe, string, error)
		CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (*types.Recipe, error)
		ReadRecipe(ctx context.Context, recipeID string) (*types.Recipe, error)
		SearchRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Recipe, string, error)
		UpdateRecipe(ctx context.Context, input *types.RecipeUpdateRequestInput) error
		ArchiveRecipe(ctx context.Context, recipeID string) error
		RecipeEstimatedPrepSteps(ctx context.Context, recipeID string)
		RecipeImageUpload(ctx context.Context)
		RecipeMermaid(ctx context.Context, recipeID string) (string, error)
		CloneRecipe(ctx context.Context, recipeID string) (*types.Recipe, error)

		ListRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeStep, string, error)
		CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error)
		ReadRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepUpdateRequestInput) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error
		RecipeStepImageUpload(ctx context.Context)

		ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepProduct, string, error)
		CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error)
		ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductUpdateRequestInput) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error

		ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepInstrument, string, error)
		CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error)
		ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error)
		UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentUpdateRequestInput) error
		ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error

		ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepIngredient, string, error)
		CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error)
		ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientUpdateRequestInput) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error

		ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error)
		CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error)
		ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskUpdateRequestInput) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error

		ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepCompletionCondition, string, error)
		CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*types.RecipeStepCompletionCondition, error)
		ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error)
		UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error
		ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error

		ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepVessel, string, error)
		CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error)
		ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselUpdateRequestInput) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error

		ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error)
		ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error)
		CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error)
		UpdateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingUpdateRequestInput) error
		ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error
	}

	recipeManager struct {
		db                   database.DataManager
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		recipeSearchIndex    textsearch.IndexSearcher[eatingindexing.RecipeSearchSubset]
	}
)

var errUnimplemented = errors.New("not implemented") // TODO: DELETE ME

func NewRecipeManager(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	db database.DataManager,
	cfg *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
	searchConfig *textsearchcfg.Config,
	metricsProvider metrics.Provider,
) (RecipeManager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide publisher for data changes topic: %w", err)
	}

	validIngredientStatesSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.RecipeSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeRecipes)
	if err != nil {
		return nil, fmt.Errorf("failed to provide search index for %s index", eatingindexing.IndexTypeValidIngredientStates)
	}

	m := &recipeManager{
		db:                   db,
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(recipeManagerName)),
		logger:               logging.EnsureLogger(logger).WithName(recipeManagerName),
		dataChangesPublisher: dataChangesPublisher,
		recipeSearchIndex:    validIngredientStatesSearchIndex,
	}

	return m, nil
}

func (r *recipeManager) ListRecipes(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Recipe, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipe(ctx context.Context, input *types.RecipeCreationRequestInput) (*types.Recipe, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) SearchRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) ([]*types.Recipe, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) UpdateRecipe(ctx context.Context, input *types.RecipeUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipe(ctx context.Context, recipeID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) RecipeEstimatedPrepSteps(ctx context.Context, recipeID string) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()
}

func (r *recipeManager) RecipeImageUpload(ctx context.Context) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()
}

func (r *recipeManager) RecipeMermaid(ctx context.Context, recipeID string) (string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return "", errUnimplemented
}

func (r *recipeManager) CloneRecipe(ctx context.Context, recipeID string) (*types.Recipe, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ListRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeStep, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) RecipeStepImageUpload(ctx context.Context) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()
}

func (r *recipeManager) ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepProduct, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductCreationRequestInput) (*types.RecipeStepProduct, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepProductUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepInstrument, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentCreationRequestInput) (*types.RecipeStepInstrument, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepInstrumentUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepIngredient, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientCreationRequestInput) (*types.RecipeStepIngredient, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepIngredientUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipePrepTask, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskCreationRequestInput) (*types.RecipePrepTask, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipePrepTask(ctx context.Context, recipeID string, input *types.RecipePrepTaskUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepCompletionCondition, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionCreationRequestInput) (*types.RecipeStepCompletionCondition, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepCompletionConditionUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*types.RecipeStepVessel, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselCreationRequestInput) (*types.RecipeStepVessel, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *types.RecipeStepVesselUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*types.RecipeRating, string, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, "", errUnimplemented
}

func (r *recipeManager) ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*types.RecipeRating, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) CreateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (r *recipeManager) UpdateRecipeRating(ctx context.Context, recipeID string, input *types.RecipeRatingUpdateRequestInput) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}

func (r *recipeManager) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	return errUnimplemented
}
