package manager

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

var (
	_ MealPlanningManager = (*mealPlanningManager)(nil)
)

type (
	mealPlanningManager struct {
		tracer tracing.Tracer
		logger logging.Logger
		db     types.MealPlanningDataManager
	}
)

func NewMealPlanningManager() (MealPlanningManager, error) {
	m := &mealPlanningManager{
		//
	}

	return m, nil
}

func (m *mealPlanningManager) ListMeals(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Meal, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(m.logger.Clone())

	results, err := m.db.GetMeals(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching meals from database")
	}

	return results.Data, "", nil
}

func (m *mealPlanningManager) CreateMeal(ctx context.Context, input *types.MealCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.Meal{}, nil
}

func (m *mealPlanningManager) SearchMeals(ctx context.Context, query string) ([]*types.Meal, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.Meal{}, nil
}

func (m *mealPlanningManager) ArchiveMeal(ctx context.Context, mealID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlan(ctx context.Context, filter *filtering.QueryFilter) ([]*types.MealPlan, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlan{}, "", nil
}

func (m *mealPlanningManager) CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlan{}, nil
}

func (m *mealPlanningManager) UpdateMealPlan(ctx context.Context, input *types.MealPlanUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlan(ctx context.Context, mealPlanID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) FinalizeMealPlan(ctx context.Context) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlanEvent(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanEvent, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlanEvent{}, "", nil
}

func (m *mealPlanningManager) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlanEvent{}, nil
}

func (m *mealPlanningManager) UpdateMealPlanEvent(ctx context.Context, input *types.MealPlanEventUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlanOption(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanOption, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlanOption{}, "", nil
}

func (m *mealPlanningManager) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlanOption{}, nil
}

func (m *mealPlanningManager) UpdateMealPlanOption(ctx context.Context, input *types.MealPlanOptionUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlanOptionVote(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanOptionVote, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlanOptionVote{}, "", nil
}

func (m *mealPlanningManager) CreateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVoteCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlanOptionVote{}, nil
}

func (m *mealPlanningManager) UpdateMealPlanOptionVote(ctx context.Context, input *types.MealPlanOptionVoteUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionVoteID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanTask, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlanTask{}, "", nil
}

func (m *mealPlanningManager) ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlanTask{}, nil
}

func (m *mealPlanningManager) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) MealPlanTaskStatusChange(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, filter *filtering.QueryFilter) ([]*types.MealPlanGroceryListItem, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.MealPlanGroceryListItem{}, "", nil
}

func (m *mealPlanningManager) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.MealPlanGroceryListItem{}, nil
}

func (m *mealPlanningManager) UpdateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.IngredientPreference, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.IngredientPreference{}, "", nil
}

func (m *mealPlanningManager) CreateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) UpdateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ListInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.InstrumentOwnership, string, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return []*types.InstrumentOwnership{}, "", nil
}

func (m *mealPlanningManager) CreateInstrumentOwnership(ctx context.Context, input *types.InstrumentOwnershipCreationRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ReadInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*types.InstrumentOwnership, error) {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return &types.InstrumentOwnership{}, nil
}

func (m *mealPlanningManager) UpdateInstrumentOwnership(ctx context.Context, input *types.InstrumentOwnershipUpdateRequestInput) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (m *mealPlanningManager) ArchiveInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error {
	_, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}
