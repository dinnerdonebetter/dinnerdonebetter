package mockmanagers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

type MockMealPlanningManager struct {
	mock.Mock
}

func (m *MockMealPlanningManager) ListMeals(ctx context.Context, filter *filtering.QueryFilter) ([]*mealplanning.Meal, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*mealplanning.Meal), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMeal(ctx context.Context, input *mealplanning.MealCreationRequestInput) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMeal(ctx context.Context, mealID string) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, mealID)

	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) SearchMeals(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*mealplanning.Meal), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ArchiveMeal(ctx context.Context, mealID, ownerID string) error {
	returnValues := m.Called(ctx, mealID, ownerID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlans(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlan, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlan), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlan(ctx context.Context, input *mealplanning.MealPlanCreationRequestInput) (*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.MealPlan), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlan(ctx context.Context, mealPlanID, ownerID string) (*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx, mealPlanID, ownerID)

	return returnValues.Get(0).(*mealplanning.MealPlan), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlan(ctx context.Context, mealPlanID, ownerID string, input *mealplanning.MealPlanUpdateRequestInput) error {
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

func (m *MockMealPlanningManager) ListMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlanEvent, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlanEvent), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanEvent(ctx context.Context, input *mealplanning.MealPlanEventCreationRequestInput) (*mealplanning.MealPlanEvent, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.MealPlanEvent), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*mealplanning.MealPlanEvent, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)

	return returnValues.Get(0).(*mealplanning.MealPlanEvent), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string, input *mealplanning.MealPlanEventUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlanOption, string, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlanOption), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanOption(ctx context.Context, input *mealplanning.MealPlanOptionCreationRequestInput) (*mealplanning.MealPlanOption, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.MealPlanOption), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*mealplanning.MealPlanOption, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

	return returnValues.Get(0).(*mealplanning.MealPlanOption), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, input *mealplanning.MealPlanOptionUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlanOptionVote, string, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlanOptionVote), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanOptionVotes(ctx context.Context, input *mealplanning.MealPlanOptionVoteCreationRequestInput) ([]*mealplanning.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).([]*mealplanning.MealPlanOptionVote), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*mealplanning.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

	return returnValues.Get(0).(*mealplanning.MealPlanOptionVote), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string, input *mealplanning.MealPlanOptionVoteUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanTasksByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlanTask, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlanTask), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanTask(ctx context.Context, mealPlanID, mealPlanTaskID string) (*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanTaskID)

	return returnValues.Get(0).(*mealplanning.MealPlanTask), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanTask(ctx context.Context, input *mealplanning.MealPlanTaskCreationRequestInput) (*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.MealPlanTask), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) MealPlanTaskStatusChange(ctx context.Context, input *mealplanning.MealPlanTaskStatusChangeRequestInput) error {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListMealPlanGroceryListItemsByMealPlan(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) ([]*mealplanning.MealPlanGroceryListItem, string, error) {
	returnValues := m.Called(ctx, mealPlanID, filter)

	return returnValues.Get(0).([]*mealplanning.MealPlanGroceryListItem), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateMealPlanGroceryListItem(ctx context.Context, input *mealplanning.MealPlanGroceryListItemCreationRequestInput) (*mealplanning.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.MealPlanGroceryListItem), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*mealplanning.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)

	return returnValues.Get(0).(*mealplanning.MealPlanGroceryListItem), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string, input *mealplanning.MealPlanGroceryListItemUpdateRequestInput) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) error {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListUserIngredientPreferences(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*mealplanning.UserIngredientPreference, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*mealplanning.UserIngredientPreference), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) ReadUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) (*mealplanning.UserIngredientPreference, error) {
	returnValues := m.Called(ctx, ownerID, ingredientPreferenceID)

	return returnValues.Get(0).(*mealplanning.UserIngredientPreference), returnValues.Error(1)
}

func (m *MockMealPlanningManager) CreateUserIngredientPreference(ctx context.Context, ownerID string, input *mealplanning.UserIngredientPreferenceCreationRequestInput) ([]*mealplanning.UserIngredientPreference, error) {
	returnValues := m.Called(ctx, ownerID, input)

	return returnValues.Get(0).([]*mealplanning.UserIngredientPreference), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateUserIngredientPreference(ctx context.Context, ingredientPreferenceID, ownerID string, input *mealplanning.UserIngredientPreferenceUpdateRequestInput) error {
	returnValues := m.Called(ctx, ingredientPreferenceID, ownerID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveUserIngredientPreference(ctx context.Context, ownerID, ingredientPreferenceID string) error {
	returnValues := m.Called(ctx, ownerID, ingredientPreferenceID)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ListAccountInstrumentOwnerships(ctx context.Context, ownerID string, filter *filtering.QueryFilter) ([]*mealplanning.AccountInstrumentOwnership, string, error) {
	returnValues := m.Called(ctx, ownerID, filter)

	return returnValues.Get(0).([]*mealplanning.AccountInstrumentOwnership), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockMealPlanningManager) CreateAccountInstrumentOwnership(ctx context.Context, input *mealplanning.AccountInstrumentOwnershipCreationRequestInput) (*mealplanning.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.AccountInstrumentOwnership), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) ReadAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) (*mealplanning.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, ownerID, instrumentOwnershipID)

	return returnValues.Get(0).(*mealplanning.AccountInstrumentOwnership), returnValues.Get(1).(error)
}

func (m *MockMealPlanningManager) UpdateAccountInstrumentOwnership(ctx context.Context, instrumentOwnershipID, ownerID string, input *mealplanning.AccountInstrumentOwnershipUpdateRequestInput) error {
	returnValues := m.Called(ctx, instrumentOwnershipID, ownerID, input)

	return returnValues.Get(0).(error)
}

func (m *MockMealPlanningManager) ArchiveAccountInstrumentOwnership(ctx context.Context, ownerID, instrumentOwnershipID string) error {
	returnValues := m.Called(ctx, ownerID, instrumentOwnershipID)

	return returnValues.Get(0).(error)
}
