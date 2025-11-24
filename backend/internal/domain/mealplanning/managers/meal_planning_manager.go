package managers

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
)

const (
	mealPlannerName = "meal_planning_manager"
)

type (
	MealPlanningManager interface {
		ListMeals(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Meal], error)
		CreateMeal(ctx context.Context, creatorID string, input *types.MealCreationRequestInput) (*types.Meal, error)
		ReadMeal(ctx context.Context, mealID string) (*types.Meal, error)
		SearchMeals(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.Meal, error)
		ArchiveMeal(ctx context.Context, mealID, ownerID string) error

		ListMealPlans(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlan], error)
		CreateMealPlan(ctx context.Context, ownerID, creatorID string, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error)
		ReadMealPlan(ctx context.Context, mealPlanID, ownerID string) (*types.MealPlan, error)
		UpdateMealPlan(ctx context.Context, mealPlanID, ownerID string, input *types.MealPlanUpdateRequestInput) error
		ArchiveMealPlan(ctx context.Context, mealPlanID, ownerID string) error
		FinalizeMealPlan(ctx context.Context, mealPlanID, ownerID string) (bool, error)

		ListMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanEvent], error)
		CreateMealPlanEvent(ctx context.Context, mealPlanID string, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error)
		ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error)
		UpdateMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanEventUpdateRequestInput) error
		ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error

		ListMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanOption], error)
		CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error)
		CreateMealPlanOptionWithEventID(ctx context.Context, mealPlanEventID string, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error)
		ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error)
		UpdateMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, input *types.MealPlanOptionUpdateRequestInput) error
		ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error

		ListMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanOptionVote], error)
		CreateMealPlanOptionVotes(ctx context.Context, creatorID string, input *types.MealPlanOptionVoteCreationRequestInput) ([]*types.MealPlanOptionVote, error)
		ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error)
		UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string, input *types.MealPlanOptionVoteUpdateRequestInput) error
		ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error

		ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanTask, string, error)
		ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error)
		CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskCreationRequestInput) (*types.MealPlanTask, error)
		MealPlanTaskStatusChange(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error

		ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanGroceryListItem, string, error)
		CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemCreationRequestInput) (*types.MealPlanGroceryListItem, error)
		ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error)
		UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) error
		ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error

		ReadUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) (*types.UserIngredientPreference, error)
		ListUserIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.UserIngredientPreference], error)
		CreateUserIngredientPreference(ctx context.Context, ownerID string, input *types.UserIngredientPreferenceCreationRequestInput) ([]*types.UserIngredientPreference, error)
		UpdateUserIngredientPreference(ctx context.Context, ingredientPreferenceID, ownerID string, input *types.UserIngredientPreferenceUpdateRequestInput) error
		ArchiveUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error

		ListAccountInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInstrumentOwnership], error)
		CreateAccountInstrumentOwnership(ctx context.Context, ownerID string, input *types.AccountInstrumentOwnershipCreationRequestInput) (*types.AccountInstrumentOwnership, error)
		ReadAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*types.AccountInstrumentOwnership, error)
		UpdateAccountInstrumentOwnership(ctx context.Context, instrumentOwnershipID, ownerID string, input *types.AccountInstrumentOwnershipUpdateRequestInput) error
		ArchiveAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error
	}

	mealPlanningManager struct {
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		mealsSearchIndex     textsearch.IndexSearcher[eatingindexing.MealSearchSubset]
		db                   types.Repository
	}
)

var (
	_ MealPlanningManager = (*mealPlanningManager)(nil)
)

func NewMealPlanningManager(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	db types.Repository,
	cfg *msgconfig.QueuesConfig,
	publisherProvider messagequeue.PublisherProvider,
	searchConfig *textsearchcfg.Config,
	metricsProvider metrics.Provider,
) (MealPlanningManager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to provide publisher for data changes topic: %w", err)
	}

	mealsSearchIndex, err := textsearchcfg.ProvideIndex[eatingindexing.MealSearchSubset](ctx, logger, tracerProvider, metricsProvider, searchConfig, eatingindexing.IndexTypeMeals)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing recipe index manager")
	}

	m := &mealPlanningManager{
		db:                   db,
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(mealPlannerName)),
		logger:               logging.EnsureLogger(logger).WithName(mealPlannerName),
		dataChangesPublisher: dataChangesPublisher,
		mealsSearchIndex:     mealsSearchIndex,
	}

	return m, nil
}

func (m *mealPlanningManager) ListMeals(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Meal], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	results, err := m.db.GetMeals(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, m.logger, span, "fetching meals from database")
	}

	return results, nil
}

func (m *mealPlanningManager) CreateMeal(ctx context.Context, creatorID string, input *types.MealCreationRequestInput) (*types.Meal, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealCreationRequestInputToMealDatabaseCreationInput(input)
	convertedInput.CreatedByUser = creatorID
	logger := m.logger.WithSpan(span).WithValue(keys.MealIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealIDKey, convertedInput.ID)

	created, err := m.db.CreateMeal(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealCreatedServiceEventType, map[string]any{
		keys.MealIDKey: created.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	meal, err := m.db.GetMeal(ctx, mealID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal")
	}

	return meal, nil
}

func (m *mealPlanningManager) SearchMeals(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.Meal, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	var results []*types.Meal
	if useDatabase {
		allResults, err := m.db.SearchForMeals(ctx, query, filter)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching for meals")
		}

		results = allResults.Data
	} else {
		mealSubsets, err := m.mealsSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching external provider for meals")
		}

		ids := []string{}
		for _, mealSubset := range mealSubsets {
			ids = append(ids, mealSubset.ID)
		}

		results, err = m.db.GetMealsWithIDs(ctx, ids)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "fetching meals from database")
		}
	}

	return results, nil
}

func (m *mealPlanningManager) ArchiveMeal(ctx context.Context, mealID, ownerID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealIDKey: mealID,
		keys.UserIDKey: ownerID,
	})
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	if err := m.db.ArchiveMeal(ctx, mealID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealArchivedServiceEventType, map[string]any{
		keys.MealIDKey: mealID,
	}))

	return nil
}

func (m *mealPlanningManager) ListMealPlans(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlan], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValue(keys.AccountIDKey, ownerID)
	tracing.AttachToSpan(span, keys.AccountIDKey, ownerID)

	mealPlans, err := m.db.GetMealPlansForAccount(ctx, ownerID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching list of meal plans for account")
	}

	return mealPlans, nil
}

func (m *mealPlanningManager) CreateMealPlan(ctx context.Context, ownerID, creatorID string, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	if creatorID == "" {
		return nil, internalerrors.ErrEmptyInputParameter
	}

	if ownerID == "" {
		return nil, internalerrors.ErrEmptyInputParameter
	}

	convertedInput := converters.ConvertMealPlanCreationRequestInputToMealPlanDatabaseCreationInput(input)
	convertedInput.CreatedByUser = creatorID
	convertedInput.BelongsToAccount = ownerID

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlan(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanCreatedServiceEventType, map[string]any{
		keys.MealPlanIDKey: created.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMealPlan(ctx context.Context, mealPlanID, ownerID string) (*types.MealPlan, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey: mealPlanID,
		keys.UserIDKey:     ownerID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	mealPlan, err := m.db.GetMealPlan(ctx, mealPlanID, ownerID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	return mealPlan, nil
}

func (m *mealPlanningManager) UpdateMealPlan(ctx context.Context, mealPlanID, ownerID string, input *types.MealPlanUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey: mealPlanID,
		keys.UserIDKey:     ownerID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	existingMealPlan, err := m.db.GetMealPlan(ctx, mealPlanID, ownerID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan to update")
	}

	existingMealPlan.Update(input)
	if err = m.db.UpdateMealPlan(ctx, existingMealPlan); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanUpdatedServiceEventType, map[string]any{
		keys.MealPlanIDKey: mealPlanID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlan(ctx context.Context, mealPlanID, ownerID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey: mealPlanID,
		keys.UserIDKey:     ownerID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	if err := m.db.ArchiveMealPlan(ctx, mealPlanID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanArchivedServiceEventType, map[string]any{
		keys.MealPlanIDKey: mealPlanID,
	}))

	return nil
}

func (m *mealPlanningManager) FinalizeMealPlan(ctx context.Context, mealPlanID, ownerID string) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey: mealPlanID,
		keys.UserIDKey:     ownerID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	finalized, err := m.db.AttemptToFinalizeMealPlan(ctx, mealPlanID, ownerID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanFinalizedServiceEventType, map[string]any{
		keys.MealPlanIDKey: mealPlanID,
	}))

	return finalized, nil
}

func (m *mealPlanningManager) ListMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanEvent], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	mealPlanEvents, err := m.db.GetMealPlanEvents(ctx, mealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan events")
	}

	return mealPlanEvents, nil
}

func (m *mealPlanningManager) CreateMealPlanEvent(ctx context.Context, mealPlanID string, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealPlanEventCreationRequestInputToMealPlanEventDatabaseCreationInput(input)
	convertedInput.BelongsToMealPlan = mealPlanID
	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanEventIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlanEvent(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "created meal plan event")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanEventCreatedServiceEventType, map[string]any{
		keys.MealPlanEventIDKey: created.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	mealPlanEvent, err := m.db.GetMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan event")
	}

	return mealPlanEvent, nil
}

func (m *mealPlanningManager) UpdateMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanEventUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	existingMealPlanEvent, err := m.db.GetMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan event to update")
	}

	existingMealPlanEvent.Update(input)
	if err = m.db.UpdateMealPlanEvent(ctx, existingMealPlanEvent); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan event")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanEventUpdatedServiceEventType, map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if err := m.db.ArchiveMealPlanEvent(ctx, mealPlanID, mealPlanEventID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan event")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanEventArchivedServiceEventType, map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	}))

	return nil
}

func (m *mealPlanningManager) ListMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanOption], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:      mealPlanID,
		keys.MealPlanEventIDKey: mealPlanEventID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	results, err := m.db.GetMealPlanOptions(ctx, mealPlanID, mealPlanEventID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan options")
	}

	return results, nil
}

func (m *mealPlanningManager) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealPlanOptionCreationRequestInputToMealPlanOptionDatabaseCreationInput(input)
	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanOptionIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlanOption(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "created meal plan option")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionCreatedServiceEventType, map[string]any{
		keys.MealPlanOptionIDKey: convertedInput.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) CreateMealPlanOptionWithEventID(ctx context.Context, mealPlanEventID string, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	if mealPlanEventID == "" {
		return nil, internalerrors.ErrEmptyInputParameter
	}

	convertedInput := converters.ConvertMealPlanOptionCreationRequestInputToMealPlanOptionDatabaseCreationInput(input)
	// Set the meal plan event ID that was missing before
	convertedInput.BelongsToMealPlanEvent = mealPlanEventID

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanOptionIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlanOption(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "created meal plan option")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionCreatedServiceEventType, map[string]any{
		keys.MealPlanOptionIDKey: convertedInput.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	mealPlanOption, err := m.db.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan option")
	}

	return mealPlanOption, nil
}

func (m *mealPlanningManager) UpdateMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, input *types.MealPlanOptionUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	existingMealPlanOption, err := m.db.GetMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan option to update")
	}

	existingMealPlanOption.Update(input)
	if err = m.db.UpdateMealPlanOption(ctx, existingMealPlanOption); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionUpdatedServiceEventType, map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if err := m.db.ArchiveMealPlanOption(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan option")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionArchivedServiceEventType, map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	}))

	return nil
}

func (m *mealPlanningManager) ListMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlanOptionVote], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanEventIDKey:  mealPlanEventID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	results, err := m.db.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan option votes")
	}

	return results, nil
}

func (m *mealPlanningManager) CreateMealPlanOptionVotes(ctx context.Context, creatorID string, input *types.MealPlanOptionVoteCreationRequestInput) ([]*types.MealPlanOptionVote, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVotesDatabaseCreationInput(input)
	logger := m.logger.WithSpan(span).WithValue("vote_count", len(input.Votes))
	tracing.AttachToSpan(span, "vote_count", len(input.Votes))

	for i := range input.Votes {
		convertedInput.Votes[i].ByUser = creatorID
	}

	created, err := m.db.CreateMealPlanOptionVote(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "created meal plan option votes")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionVoteCreatedServiceEventType, map[string]any{
		"vote_count": len(input.Votes),
		"created":    len(created),
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:           mealPlanID,
		keys.MealPlanEventIDKey:      mealPlanEventID,
		keys.MealPlanOptionIDKey:     mealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: mealPlanOptionVoteID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	mealPlanOptionVote, err := m.db.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan option vote")
	}

	return mealPlanOptionVote, nil
}

func (m *mealPlanningManager) UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string, input *types.MealPlanOptionVoteUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:           mealPlanID,
		keys.MealPlanEventIDKey:      mealPlanEventID,
		keys.MealPlanOptionIDKey:     mealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: mealPlanOptionVoteID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	existingMealPlanOptionVote, err := m.db.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan option vote to update")
	}

	existingMealPlanOptionVote.Update(input)
	if err = m.db.UpdateMealPlanOptionVote(ctx, existingMealPlanOptionVote); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option vote")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionVoteUpdatedServiceEventType, map[string]any{
		keys.MealPlanIDKey:           mealPlanID,
		keys.MealPlanEventIDKey:      mealPlanEventID,
		keys.MealPlanOptionIDKey:     mealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: mealPlanOptionVoteID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:           mealPlanID,
		keys.MealPlanEventIDKey:      mealPlanEventID,
		keys.MealPlanOptionIDKey:     mealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: mealPlanOptionVoteID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	if err := m.db.ArchiveMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan option vote")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanOptionVoteArchivedServiceEventType, map[string]any{
		keys.MealPlanIDKey:           mealPlanID,
		keys.MealPlanEventIDKey:      mealPlanEventID,
		keys.MealPlanOptionIDKey:     mealPlanOptionID,
		keys.MealPlanOptionVoteIDKey: mealPlanOptionVoteID,
	}))

	return nil
}

func (m *mealPlanningManager) ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, _ *filtering.QueryFilter) ([]*types.MealPlanTask, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := m.db.GetMealPlanTasksForMealPlan(ctx, mealPlanID)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "getting meal plan tasks for meal plan")
	}

	return results, "", nil
}

func (m *mealPlanningManager) ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:     mealPlanID,
		keys.MealPlanTaskIDKey: mealPlanTaskID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	result, err := m.db.GetMealPlanTask(ctx, mealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan task")
	}

	return result, nil
}

func (m *mealPlanningManager) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskCreationRequestInput) (*types.MealPlanTask, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput(input)
	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanTaskIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlanTask(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanTaskCreatedServiceEventType, map[string]any{
		keys.MealPlanTaskIDKey: convertedInput.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) MealPlanTaskStatusChange(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanTaskIDKey, input.ID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, input.ID)

	if err := m.db.ChangeMealPlanTaskStatus(ctx, input); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing meal plan task status")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanTaskStatusChangedServiceEventType, map[string]any{
		keys.MealPlanTaskIDKey: input.ID,
	}))

	return nil
}

func (m *mealPlanningManager) ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, mealPlanID string, _ *filtering.QueryFilter) ([]*types.MealPlanGroceryListItem, string, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := m.db.GetMealPlanGroceryListItemsForMealPlan(ctx, mealPlanID)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching meal plan grocery list items for meal plan")
	}

	return results, "", nil
}

func (m *mealPlanningManager) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemCreationRequestInput) (*types.MealPlanGroceryListItem, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	convertedInput := converters.ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(input)
	logger := m.logger.WithSpan(span).WithValue(keys.MealPlanGroceryListItemIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, convertedInput.ID)

	created, err := m.db.CreateMealPlanGroceryListItem(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan grocery list item")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanGroceryListItemCreatedServiceEventType, map[string]any{
		keys.MealPlanGroceryListItemIDKey: convertedInput.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:                mealPlanID,
		keys.MealPlanGroceryListItemIDKey: mealPlanGroceryListItemID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	result, err := m.db.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan grocery list item")
	}

	return result, nil
}

func (m *mealPlanningManager) UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:                mealPlanID,
		keys.MealPlanGroceryListItemIDKey: mealPlanGroceryListItemID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	existingMealPlanGroceryListItem, err := m.db.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching meal plan grocery list item to update")
	}

	existingMealPlanGroceryListItem.Update(input)
	if err = m.db.UpdateMealPlanGroceryListItem(ctx, existingMealPlanGroceryListItem); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list item")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanGroceryListItemUpdatedServiceEventType, map[string]any{
		keys.MealPlanIDKey:                mealPlanID,
		keys.MealPlanGroceryListItemIDKey: mealPlanGroceryListItemID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.MealPlanIDKey:                mealPlanID,
		keys.MealPlanGroceryListItemIDKey: mealPlanGroceryListItemID,
	})
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	if err := m.db.ArchiveMealPlanGroceryListItem(ctx, mealPlanGroceryListItemID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan grocery list item")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.MealPlanGroceryListItemArchivedServiceEventType, map[string]any{
		keys.MealPlanIDKey:                mealPlanID,
		keys.MealPlanGroceryListItemIDKey: mealPlanGroceryListItemID,
	}))

	return nil
}

func (m *mealPlanningManager) ListUserIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.UserIngredientPreference], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValue(keys.UserIDKey, ownerID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	results, err := m.db.GetUserIngredientPreferences(ctx, ownerID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching ingredient preferences")
	}

	return results, nil
}

func (m *mealPlanningManager) ReadUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) (*types.UserIngredientPreference, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValue(keys.UserIDKey, ownerID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	result, err := m.db.GetUserIngredientPreference(ctx, ingredientPreferenceID, ownerID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching ingredient preferences")
	}

	return result, nil
}

func (m *mealPlanningManager) CreateUserIngredientPreference(ctx context.Context, ownerID string, input *types.UserIngredientPreferenceCreationRequestInput) ([]*types.UserIngredientPreference, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}

	if ownerID == "" {
		return nil, internalerrors.ErrEmptyInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.ValidIngredientGroupIDKey: input.ValidIngredientGroupID,
		keys.ValidIngredientIDKey:      input.ValidIngredientID,
	})
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, input.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, input.ValidIngredientID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating ingredient preference creation request input")
	}

	convertedInput := converters.ConvertUserIngredientPreferenceCreationRequestInputToUserIngredientPreferenceDatabaseCreationInput(input)
	convertedInput.BelongsToUser = ownerID

	created, err := m.db.CreateUserIngredientPreference(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating ingredient preference")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.UserIngredientPreferenceCreatedServiceEventType, map[string]any{
		keys.ValidIngredientGroupIDKey: input.ValidIngredientGroupID,
		keys.ValidIngredientIDKey:      input.ValidIngredientID,
		"created":                      len(created),
	}))

	return created, nil
}

func (m *mealPlanningManager) UpdateUserIngredientPreference(ctx context.Context, ingredientPreferenceID, ownerID string, input *types.UserIngredientPreferenceUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: ingredientPreferenceID,
		keys.UserIDKey:                     ownerID,
	})
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, ingredientPreferenceID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	existingUserIngredientPreference, err := m.db.GetUserIngredientPreference(ctx, ingredientPreferenceID, ownerID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching UserIngredientPreference to update")
	}

	existingUserIngredientPreference.Update(input)
	if err = m.db.UpdateUserIngredientPreference(ctx, existingUserIngredientPreference); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating ingredient preference")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.UserIngredientPreferenceUpdatedServiceEventType, map[string]any{
		keys.UserIngredientPreferenceIDKey: ingredientPreferenceID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.UserIngredientPreferenceIDKey: ingredientPreferenceID,
		keys.UserIDKey:                     ownerID,
	})
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, ingredientPreferenceID)
	tracing.AttachToSpan(span, keys.UserIDKey, ownerID)

	if err := m.db.ArchiveUserIngredientPreference(ctx, ingredientPreferenceID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving ingredient preference")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.UserIngredientPreferenceArchivedServiceEventType, map[string]any{
		keys.UserIngredientPreferenceIDKey: ingredientPreferenceID,
	}))

	return nil
}

func (m *mealPlanningManager) ListAccountInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInstrumentOwnership], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	logger := m.logger.WithSpan(span).WithValue(keys.AccountIDKey, ownerID)
	tracing.AttachToSpan(span, keys.AccountIDKey, ownerID)

	results, err := m.db.GetAccountInstrumentOwnerships(ctx, ownerID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching instrument ownerships")
	}

	return results, nil
}

func (m *mealPlanningManager) CreateAccountInstrumentOwnership(ctx context.Context, ownerID string, input *types.AccountInstrumentOwnershipCreationRequestInput) (*types.AccountInstrumentOwnership, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	convertedInput := converters.ConvertAccountInstrumentOwnershipCreationRequestInputToAccountInstrumentOwnershipDatabaseCreationInput(input)
	convertedInput.BelongsToAccount = ownerID

	logger := m.logger.WithSpan(span).WithValue(keys.AccountInstrumentOwnershipIDKey, convertedInput.ID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, convertedInput.ID)

	created, err := m.db.CreateAccountInstrumentOwnership(ctx, convertedInput)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating instrument ownership")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.AccountInstrumentOwnershipCreatedServiceEventType, map[string]any{
		keys.AccountInstrumentOwnershipIDKey: convertedInput.ID,
	}))

	return created, nil
}

func (m *mealPlanningManager) ReadAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*types.AccountInstrumentOwnership, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.AccountIDKey:                    ownerID,
		keys.AccountInstrumentOwnershipIDKey: instrumentOwnershipID,
	})
	tracing.AttachToSpan(span, keys.AccountIDKey, ownerID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, instrumentOwnershipID)

	result, err := m.db.GetAccountInstrumentOwnership(ctx, instrumentOwnershipID, ownerID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching instrument ownership")
	}

	return result, nil
}

func (m *mealPlanningManager) UpdateAccountInstrumentOwnership(ctx context.Context, instrumentOwnershipID, ownerID string, input *types.AccountInstrumentOwnershipUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return internalerrors.ErrNilInputParameter
	}

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.AccountIDKey:                    ownerID,
		keys.AccountInstrumentOwnershipIDKey: instrumentOwnershipID,
	})
	tracing.AttachToSpan(span, keys.AccountIDKey, ownerID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, instrumentOwnershipID)

	existingAccountInstrumentOwnership, err := m.db.GetAccountInstrumentOwnership(ctx, instrumentOwnershipID, ownerID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching instrument ownership to update")
	}

	existingAccountInstrumentOwnership.Update(input)
	if err = m.db.UpdateAccountInstrumentOwnership(ctx, existingAccountInstrumentOwnership); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating instrument ownership")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.AccountInstrumentOwnershipUpdatedServiceEventType, map[string]any{
		keys.AccountInstrumentOwnershipIDKey: instrumentOwnershipID,
	}))

	return nil
}

func (m *mealPlanningManager) ArchiveAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.WithSpan(span).WithValues(map[string]any{
		keys.AccountIDKey:                    ownerID,
		keys.AccountInstrumentOwnershipIDKey: instrumentOwnershipID,
	})
	tracing.AttachToSpan(span, keys.AccountIDKey, ownerID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, instrumentOwnershipID)

	if err := m.db.ArchiveAccountInstrumentOwnership(ctx, instrumentOwnershipID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving instrument ownership")
	}

	m.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, types.AccountInstrumentOwnershipArchivedServiceEventType, map[string]any{
		keys.AccountInstrumentOwnershipIDKey: instrumentOwnershipID,
	}))

	return nil
}
