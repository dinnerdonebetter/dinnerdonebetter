package mockmanagers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

type MockMealPlanningManager struct {
	mock.Mock
}

func (m *MockMealPlanningManager) ListMeals(ctx context.Context, filter *filtering.QueryFilter) ([]*types.Meal, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.Meal), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMeal(ctx context.Context, input *types.MealCreationRequestInput) (*types.Meal, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	returnValues := m.Called(ctx, mealID)

	return returnValues.Get(0).(*types.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) SearchMeals(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.Meal, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ArchiveMeal(ctx context.Context, mealID, ownerID string) error {
	returnValues := m.Called(ctx, mealID, ownerID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlans(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.MealPlan, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*types.MealPlan), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.MealPlan), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlan(ctx context.Context, mealPlanID, ownerID string) (*types.MealPlan, error) {
	returnValues := m.Called(ctx, mealPlanID, ownerID)

	return returnValues.Get(0).(*types.MealPlan), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlan(ctx context.Context, mealPlanID, ownerID string, input *types.MealPlanUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, ownerID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlan(ctx context.Context, mealPlanID, ownerID string) error {
	returnValues := m.Called(ctx, mealPlanID, ownerID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) FinalizeMealPlan(ctx context.Context, mealPlanID, ownerID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, ownerID)

	return returnValues.Get(0).(bool), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ListMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanEvent, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*types.MealPlanEvent), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.MealPlanEvent), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)

	return returnValues.Get(0).(*types.MealPlanEvent), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string, input *types.MealPlanEventUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) ([]*types.MealPlanOption, string, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, filter)

	return returnValues.Get(0).([]*types.MealPlanOption), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionCreationRequestInput) (*types.MealPlanOption, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.MealPlanOption), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

	return returnValues.Get(0).(*types.MealPlanOption), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, input *types.MealPlanOptionUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) ([]*types.MealPlanOptionVote, string, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)

	return returnValues.Get(0).([]*types.MealPlanOptionVote), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanOptionVotes(ctx context.Context, input *types.MealPlanOptionVoteCreationRequestInput) ([]*types.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).([]*types.MealPlanOptionVote), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*types.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

	return returnValues.Get(0).(*types.MealPlanOptionVote), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string, input *types.MealPlanOptionVoteUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanTask, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*types.MealPlanTask), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanTaskID)

	return returnValues.Get(0).(*types.MealPlanTask), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskCreationRequestInput) (*types.MealPlanTask, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.MealPlanTask), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) MealPlanTaskStatusChange(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*types.MealPlanGroceryListItem, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*types.MealPlanGroceryListItem), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemCreationRequestInput) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)

	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *types.MealPlanGroceryListItemUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.IngredientPreference, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*types.IngredientPreference), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceCreationRequestInput) ([]*types.IngredientPreference, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).([]*types.IngredientPreference), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateIngredientPreference(ctx context.Context, ingredientPreferenceID, ownerID string, input *types.IngredientPreferenceUpdateRequestInput) error {
	returnValues := m.Called(ctx, ingredientPreferenceID, ownerID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error {
	returnValues := m.Called(ctx, ownerID, ingredientPreferenceID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*types.InstrumentOwnership, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*types.InstrumentOwnership), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateInstrumentOwnership(ctx context.Context, input *types.InstrumentOwnershipCreationRequestInput) (*types.InstrumentOwnership, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.InstrumentOwnership), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*types.InstrumentOwnership, error) {
	returnValues := m.Called(ctx, ownerID, instrumentOwnershipID)

	return returnValues.Get(0).(*types.InstrumentOwnership), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateInstrumentOwnership(ctx context.Context, instrumentOwnershipID, ownerID string, input *types.InstrumentOwnershipUpdateRequestInput) error {
	returnValues := m.Called(ctx, instrumentOwnershipID, ownerID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error {
	returnValues := m.Called(ctx, ownerID, instrumentOwnershipID)

	return returnValues.Get(0).(error)
}
