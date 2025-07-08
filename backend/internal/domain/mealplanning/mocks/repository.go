package mocks

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

var _ mealplanning.Repository = (*Repository)(nil)

type Repository struct {
	mock.Mock
}

// MealExists is a mock function.
func (m *Repository) MealExists(ctx context.Context, recipeID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMeal is a mock function.
func (m *Repository) GetMeal(ctx context.Context, recipeID string) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Error(1)
}

// GetMealByIDAndUser is a mock function.
func (m *Repository) GetMealByIDAndUser(ctx context.Context, recipeID, userID string) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, recipeID, userID)
	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Error(1)
}

// GetMeals is a mock function.
func (m *Repository) GetMeals(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.Meal], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Meal]), returnValues.Error(1)
}

func (m *Repository) GetMealsCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[mealplanning.Meal], err error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Meal]), returnValues.Error(1)
}

// SearchForMeals is a mock function.
func (m *Repository) SearchForMeals(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.Meal], error) {
	returnValues := m.Called(ctx, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Meal]), returnValues.Error(1)
}

// CreateMeal is a mock function.
func (m *Repository) CreateMeal(ctx context.Context, input *mealplanning.MealDatabaseCreationInput) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Error(1)
}

// CreateMealInDatabase is a mock function.
func (m *Repository) CreateMealInDatabase(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *mealplanning.MealDatabaseCreationInput) (*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, querier, input)
	return returnValues.Get(0).(*mealplanning.Meal), returnValues.Error(1)
}

// UpdateMeal is a mock function.
func (m *Repository) UpdateMeal(ctx context.Context, updated *mealplanning.Meal) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMeal is a mock function.
func (m *Repository) ArchiveMeal(ctx context.Context, recipeID, accountID string) error {
	return m.Called(ctx, recipeID, accountID).Error(0)
}

// MarkMealAsIndexed is a mock function.
func (m *Repository) MarkMealAsIndexed(ctx context.Context, mealID string) error {
	return m.Called(ctx, mealID).Error(0)
}

// GetMealIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetMealIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetMealsWithIDs is a mock function.
func (m *Repository) GetMealsWithIDs(ctx context.Context, ids []string) ([]*mealplanning.Meal, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*mealplanning.Meal), returnValues.Error(1)
}

// RecipeExists is a mock function.
func (m *Repository) RecipeExists(ctx context.Context, recipeID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipe is a mock function.
func (m *Repository) GetRecipe(ctx context.Context, recipeID string) (*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, recipeID)
	return returnValues.Get(0).(*mealplanning.Recipe), returnValues.Error(1)
}

// SearchForRecipes is a mock function.
func (m *Repository) SearchForRecipes(ctx context.Context, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.Recipe], error) {
	returnValues := m.Called(ctx, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Recipe]), returnValues.Error(1)
}

// GetRecipes is a mock function.
func (m *Repository) GetRecipes(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.Recipe], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Recipe]), returnValues.Error(1)
}

// GetRecipesCreatedByUser is a mock function.
func (m *Repository) GetRecipesCreatedByUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.Recipe], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.Recipe]), returnValues.Error(1)
}

// CreateRecipe is a mock function.
func (m *Repository) CreateRecipe(ctx context.Context, input *mealplanning.RecipeDatabaseCreationInput) (*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.Recipe), returnValues.Error(1)
}

// UpdateRecipe is a mock function.
func (m *Repository) UpdateRecipe(ctx context.Context, updated *mealplanning.Recipe) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipe is a mock function.
func (m *Repository) ArchiveRecipe(ctx context.Context, recipeID, accountID string) error {
	return m.Called(ctx, recipeID, accountID).Error(0)
}

// MarkRecipeAsIndexed is a mock function.
func (m *Repository) MarkRecipeAsIndexed(ctx context.Context, recipeID string) error {
	return m.Called(ctx, recipeID).Error(0)
}

// GetRecipeIDsThatNeedSearchIndexing is a mock function.
func (m *Repository) GetRecipeIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetRecipesWithIDs is a mock function.
func (m *Repository) GetRecipesWithIDs(ctx context.Context, ids []string) ([]*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*mealplanning.Recipe), returnValues.Error(1)
}

// RecipeStepExists is a mock function.
func (m *Repository) RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStep is a mock function.
func (m *Repository) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*mealplanning.RecipeStep, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID)
	return returnValues.Get(0).(*mealplanning.RecipeStep), returnValues.Error(1)
}

// GetRecipeSteps is a mock function.
func (m *Repository) GetRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStep], error) {
	returnValues := m.Called(ctx, recipeID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStep]), returnValues.Error(1)
}

// CreateRecipeStep is a mock function.
func (m *Repository) CreateRecipeStep(ctx context.Context, input *mealplanning.RecipeStepDatabaseCreationInput) (*mealplanning.RecipeStep, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStep), returnValues.Error(1)
}

// UpdateRecipeStep is a mock function.
func (m *Repository) UpdateRecipeStep(ctx context.Context, updated *mealplanning.RecipeStep) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStep is a mock function.
func (m *Repository) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	return m.Called(ctx, recipeID, recipeStepID).Error(0)
}

// RecipeStepProductExists is a mock function.
func (m *Repository) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepProduct is a mock function.
func (m *Repository) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*mealplanning.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)
	return returnValues.Get(0).(*mealplanning.RecipeStepProduct), returnValues.Error(1)
}

// GetRecipeStepProducts is a mock function.
func (m *Repository) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStepProduct], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStepProduct]), returnValues.Error(1)
}

// CreateRecipeStepProduct is a mock function.
func (m *Repository) CreateRecipeStepProduct(ctx context.Context, input *mealplanning.RecipeStepProductDatabaseCreationInput) (*mealplanning.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStepProduct), returnValues.Error(1)
}

// UpdateRecipeStepProduct is a mock function.
func (m *Repository) UpdateRecipeStepProduct(ctx context.Context, updated *mealplanning.RecipeStepProduct) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepProduct is a mock function.
func (m *Repository) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	return m.Called(ctx, recipeStepID, recipeStepProductID).Error(0)
}

// RecipeStepInstrumentExists is a mock function.
func (m *Repository) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepInstrument is a mock function.
func (m *Repository) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*mealplanning.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	return returnValues.Get(0).(*mealplanning.RecipeStepInstrument), returnValues.Error(1)
}

// GetRecipeStepInstruments is a mock function.
func (m *Repository) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStepInstrument], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStepInstrument]), returnValues.Error(1)
}

// CreateRecipeStepInstrument is a mock function.
func (m *Repository) CreateRecipeStepInstrument(ctx context.Context, input *mealplanning.RecipeStepInstrumentDatabaseCreationInput) (*mealplanning.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStepInstrument), returnValues.Error(1)
}

// UpdateRecipeStepInstrument is a mock function.
func (m *Repository) UpdateRecipeStepInstrument(ctx context.Context, updated *mealplanning.RecipeStepInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function.
func (m *Repository) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
	return m.Called(ctx, recipeStepID, recipeStepInstrumentID).Error(0)
}

// RecipeStepIngredientExists is a mock function.
func (m *Repository) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepIngredient is a mock function.
func (m *Repository) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*mealplanning.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Get(0).(*mealplanning.RecipeStepIngredient), returnValues.Error(1)
}

// GetRecipeStepIngredients is a mock function.
func (m *Repository) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStepIngredient], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStepIngredient]), returnValues.Error(1)
}

// CreateRecipeStepIngredient is a mock function.
func (m *Repository) CreateRecipeStepIngredient(ctx context.Context, input *mealplanning.RecipeStepIngredientDatabaseCreationInput) (*mealplanning.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStepIngredient), returnValues.Error(1)
}

// UpdateRecipeStepIngredient is a mock function.
func (m *Repository) UpdateRecipeStepIngredient(ctx context.Context, updated *mealplanning.RecipeStepIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function.
func (m *Repository) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}

// MealPlanExists is a mock function.
func (m *Repository) MealPlanExists(ctx context.Context, mealPlanID, accountID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, accountID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlan is a mock function.
func (m *Repository) GetMealPlan(ctx context.Context, mealPlanID, accountID string) (*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx, mealPlanID, accountID)
	return returnValues.Get(0).(*mealplanning.MealPlan), returnValues.Error(1)
}

// GetMealPlansForAccount is a mock function.
func (m *Repository) GetMealPlansForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.MealPlan], error) {
	returnValues := m.Called(ctx, accountID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.MealPlan]), returnValues.Error(1)
}

// CreateMealPlan is a mock function.
func (m *Repository) CreateMealPlan(ctx context.Context, input *mealplanning.MealPlanDatabaseCreationInput) (*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.MealPlan), returnValues.Error(1)
}

// UpdateMealPlan is a mock function.
func (m *Repository) UpdateMealPlan(ctx context.Context, updated *mealplanning.MealPlan) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlan is a mock function.
func (m *Repository) ArchiveMealPlan(ctx context.Context, mealPlanID, accountID string) error {
	return m.Called(ctx, mealPlanID, accountID).Error(0)
}

// AttemptToFinalizeMealPlan is a mock function.
func (m *Repository) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, accountID string) (changed bool, err error) {
	returnValues := m.Called(ctx, mealPlanID, accountID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUnfinalizedMealPlansWithExpiredVotingPeriods is a mock function.
func (m *Repository) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*mealplanning.MealPlan), returnValues.Error(1)
}

// GetFinalizedMealPlanIDsForTheNextWeek is a mock function.
func (m *Repository) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*mealplanning.FinalizedMealPlanDatabaseResult, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*mealplanning.FinalizedMealPlanDatabaseResult), returnValues.Error(1)
}

// GetFinalizedMealPlansWithUninitializedGroceryLists is a mock function.
func (m *Repository) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*mealplanning.MealPlan, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*mealplanning.MealPlan), returnValues.Error(1)
}

// MealPlanOptionExists is a mock function.
func (m *Repository) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanOption is a mock function.
func (m *Repository) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*mealplanning.MealPlanOption, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return returnValues.Get(0).(*mealplanning.MealPlanOption), returnValues.Error(1)
}

// GetMealPlanOptions is a mock function.
func (m *Repository) GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.MealPlanOption], error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.MealPlanOption]), returnValues.Error(1)
}

// CreateMealPlanOption is a mock function.
func (m *Repository) CreateMealPlanOption(ctx context.Context, input *mealplanning.MealPlanOptionDatabaseCreationInput) (*mealplanning.MealPlanOption, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.MealPlanOption), returnValues.Error(1)
}

// UpdateMealPlanOption is a mock function.
func (m *Repository) UpdateMealPlanOption(ctx context.Context, updated *mealplanning.MealPlanOption) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOption is a mock function.
func (m *Repository) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID).Error(0)
}

// FinalizeMealPlanOption is a mock function.
func (m *Repository) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, accountID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, accountID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// MealPlanOptionVoteExists is a mock function.
func (m *Repository) MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanOptionVote is a mock function.
func (m *Repository) GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) (*mealplanning.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	return returnValues.Get(0).(*mealplanning.MealPlanOptionVote), returnValues.Error(1)
}

// GetMealPlanOptionVotes is a mock function.
func (m *Repository) GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.MealPlanOptionVote], error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.MealPlanOptionVote]), returnValues.Error(1)
}

// CreateMealPlanOptionVote is a mock function.
func (m *Repository) CreateMealPlanOptionVote(ctx context.Context, input *mealplanning.MealPlanOptionVotesDatabaseCreationInput) ([]*mealplanning.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).([]*mealplanning.MealPlanOptionVote), returnValues.Error(1)
}

// UpdateMealPlanOptionVote is a mock function.
func (m *Repository) UpdateMealPlanOptionVote(ctx context.Context, updated *mealplanning.MealPlanOptionVote) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanOptionVote is a mock function.
func (m *Repository) ArchiveMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID).Error(0)
}

// GetMealPlanOptionVotesForMealPlanOption is a mock function.
func (m *Repository) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) ([]*mealplanning.MealPlanOptionVote, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID)
	return returnValues.Get(0).([]*mealplanning.MealPlanOptionVote), returnValues.Error(1)
}

// MealPlanEventIsEligibleForVoting is a mock function.
func (m *Repository) MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// MealPlanEventExists is a mock function.
func (m *Repository) MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanEvent is a mock function.
func (m *Repository) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*mealplanning.MealPlanEvent, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanEventID)
	return returnValues.Get(0).(*mealplanning.MealPlanEvent), returnValues.Error(1)
}

// GetMealPlanEvents is a mock function.
func (m *Repository) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.MealPlanEvent], error) {
	returnValues := m.Called(ctx, mealPlanID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.MealPlanEvent]), returnValues.Error(1)
}

// CreateMealPlanEvent is a mock function.
func (m *Repository) CreateMealPlanEvent(ctx context.Context, input *mealplanning.MealPlanEventDatabaseCreationInput) (*mealplanning.MealPlanEvent, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.MealPlanEvent), returnValues.Error(1)
}

// UpdateMealPlanEvent is a mock function.
func (m *Repository) UpdateMealPlanEvent(ctx context.Context, updated *mealplanning.MealPlanEvent) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanEvent is a mock function.
func (m *Repository) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID).Error(0)
}

// AttemptToFinalizeMealPlanEvent is a mock function.
func (m *Repository) AttemptToFinalizeMealPlanEvent(ctx context.Context, mealPlanEventID string) (changed bool, err error) {
	returnValues := m.Called(ctx, mealPlanEventID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods is a mock function.
func (m *Repository) GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods(ctx context.Context) ([]*mealplanning.MealPlanEvent, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]*mealplanning.MealPlanEvent), returnValues.Error(1)
}

// GetFinalizedMealPlanEventIDsForTheNextWeek is a mock function.
func (m *Repository) GetFinalizedMealPlanEventIDsForTheNextWeek(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// MarkMealPlanAsHavingTasksCreated is a mock function.
func (m *Repository) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	return m.Called(ctx, mealPlanID).Error(0)
}

// MealPlanTaskExists is a mock function.
func (m *Repository) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanTaskID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanTask is a mock function.
func (m *Repository) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanTaskID)
	return returnValues.Get(0).(*mealplanning.MealPlanTask), returnValues.Error(1)
}

// CreateMealPlanTask is a mock function.
func (m *Repository) CreateMealPlanTask(ctx context.Context, input *mealplanning.MealPlanTaskDatabaseCreationInput) (*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.MealPlanTask), returnValues.Error(1)
}

// GetMealPlanTasksForMealPlan is a mock function.
func (m *Repository) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, mealPlanID)
	return returnValues.Get(0).([]*mealplanning.MealPlanTask), returnValues.Error(1)
}

// CreateMealPlanTasksForMealPlanOption is a mock function.
func (m *Repository) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*mealplanning.MealPlanTaskDatabaseCreationInput) ([]*mealplanning.MealPlanTask, error) {
	returnValues := m.Called(ctx, inputs)
	return returnValues.Get(0).([]*mealplanning.MealPlanTask), returnValues.Error(1)
}

// ChangeMealPlanTaskStatus is a mock function.
func (m *Repository) ChangeMealPlanTaskStatus(ctx context.Context, input *mealplanning.MealPlanTaskStatusChangeRequestInput) error {
	return m.Called(ctx, input).Error(0)
}

// RecipePrepTaskExists implements the requisite interface.
func (m *Repository) RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipePrepTask implements the requisite interface.
func (m *Repository) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*mealplanning.RecipePrepTask, error) {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID)

	return returnValues.Get(0).(*mealplanning.RecipePrepTask), returnValues.Error(1)
}

// GetRecipePrepTasksForRecipe implements the requisite interface.
func (m *Repository) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*mealplanning.RecipePrepTask, error) {
	returnValues := m.Called(ctx, recipeID)

	return returnValues.Get(0).([]*mealplanning.RecipePrepTask), returnValues.Error(1)
}

// CreateRecipePrepTask implements the requisite interface.
func (m *Repository) CreateRecipePrepTask(ctx context.Context, input *mealplanning.RecipePrepTaskDatabaseCreationInput) (*mealplanning.RecipePrepTask, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.RecipePrepTask), returnValues.Error(1)
}

// UpdateRecipePrepTask implements the requisite interface.
func (m *Repository) UpdateRecipePrepTask(ctx context.Context, updated *mealplanning.RecipePrepTask) error {
	returnValues := m.Called(ctx, updated)

	return returnValues.Error(0)
}

// ArchiveRecipePrepTask implements the requisite interface.
func (m *Repository) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID)

	return returnValues.Error(0)
}

// MealPlanGroceryListItemExists is a mock function.
func (m *Repository) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanGroceryListItem is a mock function.
func (m *Repository) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*mealplanning.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Get(0).(*mealplanning.MealPlanGroceryListItem), returnValues.Error(1)
}

// GetMealPlanGroceryListItemsForMealPlan is a mock function.
func (m *Repository) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*mealplanning.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID)
	return returnValues.Get(0).([]*mealplanning.MealPlanGroceryListItem), returnValues.Error(1)
}

// CreateMealPlanGroceryListItem is a mock function.
func (m *Repository) CreateMealPlanGroceryListItem(ctx context.Context, input *mealplanning.MealPlanGroceryListItemDatabaseCreationInput) (*mealplanning.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.MealPlanGroceryListItem), returnValues.Error(1)
}

// UpdateMealPlanGroceryListItem is a mock function.
func (m *Repository) UpdateMealPlanGroceryListItem(ctx context.Context, updated *mealplanning.MealPlanGroceryListItem) error {
	returnValues := m.Called(ctx, updated)
	return returnValues.Error(0)
}

// ArchiveMealPlanGroceryListItem is a mock function.
func (m *Repository) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	returnValues := m.Called(ctx, mealPlanGroceryListItemID)
	return returnValues.Error(0)
}

// RecipeMediaExists is a mock function.
func (m *Repository) RecipeMediaExists(ctx context.Context, recipeMediaID string) (bool, error) {
	returnValues := m.Called(ctx, recipeMediaID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeMedia is a mock function.
func (m *Repository) GetRecipeMedia(ctx context.Context, recipeMediaID string) (*mealplanning.RecipeMedia, error) {
	returnValues := m.Called(ctx, recipeMediaID)
	return returnValues.Get(0).(*mealplanning.RecipeMedia), returnValues.Error(1)
}

// CreateRecipeMedia is a mock function.
func (m *Repository) CreateRecipeMedia(ctx context.Context, input *mealplanning.RecipeMediaDatabaseCreationInput) (*mealplanning.RecipeMedia, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeMedia), returnValues.Error(1)
}

// UpdateRecipeMedia is a mock function.
func (m *Repository) UpdateRecipeMedia(ctx context.Context, updated *mealplanning.RecipeMedia) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeMedia is a mock function.
func (m *Repository) ArchiveRecipeMedia(ctx context.Context, recipeMediaID string) error {
	return m.Called(ctx, recipeMediaID).Error(0)
}

// RecipeStepCompletionConditionExists is a mock function.
func (m *Repository) RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepCompletionCondition is a mock function.
func (m *Repository) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*mealplanning.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return returnValues.Get(0).(*mealplanning.RecipeStepCompletionCondition), returnValues.Error(1)
}

// GetRecipeStepCompletionConditions is a mock function.
func (m *Repository) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStepCompletionCondition], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStepCompletionCondition]), returnValues.Error(1)
}

// CreateRecipeStepCompletionCondition is a mock function.
func (m *Repository) CreateRecipeStepCompletionCondition(ctx context.Context, input *mealplanning.RecipeStepCompletionConditionDatabaseCreationInput) (*mealplanning.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStepCompletionCondition), returnValues.Error(1)
}

// UpdateRecipeStepCompletionCondition is a mock function.
func (m *Repository) UpdateRecipeStepCompletionCondition(ctx context.Context, updated *mealplanning.RecipeStepCompletionCondition) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepCompletionCondition is a mock function.
func (m *Repository) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}

// RecipeStepVesselExists is a mock function.
func (m *Repository) RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeStepVessel is a mock function.
func (m *Repository) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*mealplanning.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)
	return returnValues.Get(0).(*mealplanning.RecipeStepVessel), returnValues.Error(1)
}

// GetRecipeStepVessels is a mock function.
func (m *Repository) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeStepVessel], error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeStepVessel]), returnValues.Error(1)
}

// CreateRecipeStepVessel is a mock function.
func (m *Repository) CreateRecipeStepVessel(ctx context.Context, input *mealplanning.RecipeStepVesselDatabaseCreationInput) (*mealplanning.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeStepVessel), returnValues.Error(1)
}

// UpdateRecipeStepVessel is a mock function.
func (m *Repository) UpdateRecipeStepVessel(ctx context.Context, updated *mealplanning.RecipeStepVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepVessel is a mock function.
func (m *Repository) ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepVesselID string) error {
	return m.Called(ctx, recipeStepID, recipeStepVesselID).Error(0)
}

// UserIngredientPreferenceExists is a mock function.
func (m *Repository) UserIngredientPreferenceExists(ctx context.Context, userUserIngredientPreferenceID, userID string) (bool, error) {
	returnValues := m.Called(ctx, userUserIngredientPreferenceID, userID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetUserIngredientPreference is a mock function.
func (m *Repository) GetUserIngredientPreference(ctx context.Context, userUserIngredientPreferenceID, userID string) (*mealplanning.UserIngredientPreference, error) {
	returnValues := m.Called(ctx, userUserIngredientPreferenceID, userID)
	return returnValues.Get(0).(*mealplanning.UserIngredientPreference), returnValues.Error(1)
}

// GetUserIngredientPreferences is a mock function.
func (m *Repository) GetUserIngredientPreferences(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.UserIngredientPreference], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.UserIngredientPreference]), returnValues.Error(1)
}

// CreateUserIngredientPreference is a mock function.
func (m *Repository) CreateUserIngredientPreference(ctx context.Context, input *mealplanning.UserIngredientPreferenceDatabaseCreationInput) ([]*mealplanning.UserIngredientPreference, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).([]*mealplanning.UserIngredientPreference), returnValues.Error(1)
}

// UpdateUserIngredientPreference is a mock function.
func (m *Repository) UpdateUserIngredientPreference(ctx context.Context, updated *mealplanning.UserIngredientPreference) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveUserIngredientPreference is a mock function.
func (m *Repository) ArchiveUserIngredientPreference(ctx context.Context, userUserIngredientPreferenceID, userID string) error {
	return m.Called(ctx, userUserIngredientPreferenceID, userID).Error(0)
}

// AccountInstrumentOwnershipExists is a mock function.
func (m *Repository) AccountInstrumentOwnershipExists(ctx context.Context, accountInstrumentOwnershipID, accountID string) (bool, error) {
	returnValues := m.Called(ctx, accountInstrumentOwnershipID, accountID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetAccountInstrumentOwnership is a mock function.
func (m *Repository) GetAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) (*mealplanning.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, accountInstrumentOwnershipID, accountID)
	return returnValues.Get(0).(*mealplanning.AccountInstrumentOwnership), returnValues.Error(1)
}

// GetAccountInstrumentOwnerships is a mock function.
func (m *Repository) GetAccountInstrumentOwnerships(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.AccountInstrumentOwnership], error) {
	returnValues := m.Called(ctx, accountID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.AccountInstrumentOwnership]), returnValues.Error(1)
}

// CreateAccountInstrumentOwnership is a mock function.
func (m *Repository) CreateAccountInstrumentOwnership(ctx context.Context, input *mealplanning.AccountInstrumentOwnershipDatabaseCreationInput) (*mealplanning.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.AccountInstrumentOwnership), returnValues.Error(1)
}

// UpdateAccountInstrumentOwnership is a mock function.
func (m *Repository) UpdateAccountInstrumentOwnership(ctx context.Context, updated *mealplanning.AccountInstrumentOwnership) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveAccountInstrumentOwnership is a mock function.
func (m *Repository) ArchiveAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) error {
	return m.Called(ctx, accountInstrumentOwnershipID, accountID).Error(0)
}

// RecipeRatingExists is a mock function.
func (m *Repository) RecipeRatingExists(ctx context.Context, recipeID, recipeRatingID string) (bool, error) {
	returnValues := m.Called(ctx, recipeID, recipeRatingID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetRecipeRating is a mock function.
func (m *Repository) GetRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*mealplanning.RecipeRating, error) {
	returnValues := m.Called(ctx, recipeID, recipeRatingID)
	return returnValues.Get(0).(*mealplanning.RecipeRating), returnValues.Error(1)
}

// GetRecipeRatingsForRecipe is a mock function.
func (m *Repository) GetRecipeRatingsForRecipe(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeRating], error) {
	returnValues := m.Called(ctx, recipeID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeRating]), returnValues.Error(1)
}

// GetRecipeRatingsForUser is a mock function.
func (m *Repository) GetRecipeRatingsForUser(ctx context.Context, user string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[mealplanning.RecipeRating], error) {
	returnValues := m.Called(ctx, user, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[mealplanning.RecipeRating]), returnValues.Error(1)
}

// CreateRecipeRating is a mock function.
func (m *Repository) CreateRecipeRating(ctx context.Context, input *mealplanning.RecipeRatingDatabaseCreationInput) (*mealplanning.RecipeRating, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*mealplanning.RecipeRating), returnValues.Error(1)
}

// UpdateRecipeRating is a mock function.
func (m *Repository) UpdateRecipeRating(ctx context.Context, updated *mealplanning.RecipeRating) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeRating is a mock function.
func (m *Repository) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	return m.Called(ctx, recipeID, recipeRatingID).Error(0)
}
