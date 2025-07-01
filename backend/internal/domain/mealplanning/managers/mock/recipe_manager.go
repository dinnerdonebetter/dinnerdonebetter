package mockmanagers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

type MockRecipeManager struct {
	mock.Mock
}

func (m *MockRecipeManager) ListRecipes(ctx context.Context, filter *filtering.QueryFilter) ([]*mealplanning.Recipe, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*mealplanning.Recipe), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipe(ctx context.Context, input *mealplanning.RecipeCreationRequestInput) (*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*mealplanning.Recipe), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipe(ctx context.Context, recipeID string) (*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, recipeID)

	return returnValues.Get(0).(*mealplanning.Recipe), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) SearchRecipes(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*mealplanning.Recipe, string, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*mealplanning.Recipe), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) UpdateRecipe(ctx context.Context, recipeID string, input *mealplanning.RecipeUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipe(ctx context.Context, recipeID, ownerID string) error {
	returnValues := m.Called(ctx, recipeID, ownerID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) RecipeEstimatedPrepSteps(ctx context.Context, recipeID string) ([]*mealplanning.MealPlanTaskDatabaseCreationEstimate, error) {
	returnValues := m.Called(ctx, recipeID)

	return returnValues.Get(0).([]*mealplanning.MealPlanTaskDatabaseCreationEstimate), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) RecipeMermaid(ctx context.Context, recipeID string) (string, error) {
	returnValues := m.Called(ctx, recipeID)

	return returnValues.Get(0).(string), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) CloneRecipe(ctx context.Context, recipeID, newOwnerID string) (*mealplanning.Recipe, error) {
	returnValues := m.Called(ctx, recipeID, newOwnerID)

	return returnValues.Get(0).(*mealplanning.Recipe), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) RecipeImageUpload(ctx context.Context) error {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeSteps(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStep, string, error) {
	returnValues := m.Called(ctx, recipeID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStep), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStep(ctx context.Context, recipeID string, input *mealplanning.RecipeStepCreationRequestInput) (*mealplanning.RecipeStep, error) {
	returnValues := m.Called(ctx, recipeID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStep), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*mealplanning.RecipeStep, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID)

	return returnValues.Get(0).(*mealplanning.RecipeStep), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStep(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) RecipeStepImageUpload(ctx context.Context) error {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStepProduct, string, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStepProduct), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepProductCreationRequestInput) (*mealplanning.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStepProduct), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*mealplanning.RecipeStepProduct, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)

	return returnValues.Get(0).(*mealplanning.RecipeStepProduct), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string, input *mealplanning.RecipeStepProductUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepProductID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStepInstrument, string, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStepInstrument), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepInstrumentCreationRequestInput) (*mealplanning.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStepInstrument), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*mealplanning.RecipeStepInstrument, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)

	return returnValues.Get(0).(*mealplanning.RecipeStepInstrument), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string, input *mealplanning.RecipeStepInstrumentUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepInstrumentID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStepIngredient, string, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStepIngredient), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepIngredientCreationRequestInput) (*mealplanning.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStepIngredient), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*mealplanning.RecipeStepIngredient, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)

	return returnValues.Get(0).(*mealplanning.RecipeStepIngredient), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string, input *mealplanning.RecipeStepIngredientUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipePrepTask(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipePrepTask, string, error) {
	returnValues := m.Called(ctx, recipeID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipePrepTask), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipePrepTask(ctx context.Context, recipeID string, input *mealplanning.RecipePrepTaskCreationRequestInput) (*mealplanning.RecipePrepTask, error) {
	returnValues := m.Called(ctx, recipeID, input)

	return returnValues.Get(0).(*mealplanning.RecipePrepTask), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*mealplanning.RecipePrepTask, error) {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID)

	return returnValues.Get(0).(*mealplanning.RecipePrepTask), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string, input *mealplanning.RecipePrepTaskUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	returnValues := m.Called(ctx, recipeID, recipePrepTaskID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStepCompletionCondition, string, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStepCompletionCondition), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) (*mealplanning.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStepCompletionCondition), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*mealplanning.RecipeStepCompletionCondition, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)

	return returnValues.Get(0).(*mealplanning.RecipeStepCompletionCondition), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string, input *mealplanning.RecipeStepCompletionConditionUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeStepVessel, string, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeStepVessel), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) CreateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID string, input *mealplanning.RecipeStepVesselCreationRequestInput) (*mealplanning.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, input)

	return returnValues.Get(0).(*mealplanning.RecipeStepVessel), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) ReadRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*mealplanning.RecipeStepVessel, error) {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)

	return returnValues.Get(0).(*mealplanning.RecipeStepVessel), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string, input *mealplanning.RecipeStepVesselUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) error {
	returnValues := m.Called(ctx, recipeID, recipeStepID, recipeStepVesselID)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ListRecipeRatings(ctx context.Context, recipeID string, filter *filtering.QueryFilter) ([]*mealplanning.RecipeRating, string, error) {
	returnValues := m.Called(ctx, recipeID, filter)

	return returnValues.Get(0).([]*mealplanning.RecipeRating), returnValues.Get(1).(string), returnValues.Get(2).(error)
}

func (m *MockRecipeManager) ReadRecipeRating(ctx context.Context, recipeID, recipeRatingID string) (*mealplanning.RecipeRating, error) {
	returnValues := m.Called(ctx, recipeID, recipeRatingID)

	return returnValues.Get(0).(*mealplanning.RecipeRating), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) CreateRecipeRating(ctx context.Context, recipeID string, input *mealplanning.RecipeRatingCreationRequestInput) (*mealplanning.RecipeRating, error) {
	returnValues := m.Called(ctx, recipeID, input)

	return returnValues.Get(0).(*mealplanning.RecipeRating), returnValues.Get(1).(error)
}

func (m *MockRecipeManager) UpdateRecipeRating(ctx context.Context, recipeID, recipeRatingID string, input *mealplanning.RecipeRatingUpdateRequestInput) error {
	returnValues := m.Called(ctx, recipeID, recipeRatingID, input)

	return returnValues.Get(0).(error)
}

func (m *MockRecipeManager) ArchiveRecipeRating(ctx context.Context, recipeID, recipeRatingID string) error {
	returnValues := m.Called(ctx, recipeID, recipeRatingID)

	return returnValues.Get(0).(error)
}
