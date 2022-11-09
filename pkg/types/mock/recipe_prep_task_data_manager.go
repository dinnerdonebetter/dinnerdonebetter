package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

type RecipePrepTaskDataManager struct {
	mock.Mock
}

// RecipePrepTaskExists implements the requisite interface.
func (m *RecipePrepTaskDataManager) RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error) {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Bool(0), retVals.Error(1)
}

// GetRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManager) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

// GetRecipePrepTasksForRecipe implements the requisite interface.
func (m *RecipePrepTaskDataManager) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipeID)

	return retVals.Get(0).([]*types.RecipePrepTask), retVals.Error(1)
}

// CreateRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManager) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, input)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

// UpdateRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManager) UpdateRecipePrepTask(ctx context.Context, updated *types.RecipePrepTask) error {
	retVals := m.Called(ctx, updated)

	return retVals.Error(0)
}

// ArchiveRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManager) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Error(0)
}
