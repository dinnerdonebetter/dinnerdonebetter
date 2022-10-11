package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

type RecipePrepTaskDataManager struct {
	mock.Mock
}

func (m *RecipePrepTaskDataManager) RecipePrepTaskExists(ctx context.Context, recipePrepTaskID string) (bool, error) {
	retVals := m.Called(ctx, recipePrepTaskID)

	return retVals.Bool(0), retVals.Error(1)
}

func (m *RecipePrepTaskDataManager) GetRecipePrepTask(ctx context.Context, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipePrepTaskID)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

func (m *RecipePrepTaskDataManager) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipeID)

	return retVals.Get(0).([]*types.RecipePrepTask), retVals.Error(1)
}

func (m *RecipePrepTaskDataManager) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, input)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

func (m *RecipePrepTaskDataManager) UpdateRecipePrepTask(ctx context.Context, updated *types.RecipePrepTask) error {
	retVals := m.Called(ctx, updated)

	return retVals.Error(0)
}

func (m *RecipePrepTaskDataManager) ArchiveRecipePrepTask(ctx context.Context, recipePrepTaskID string) error {
	retVals := m.Called(ctx, recipePrepTaskID)

	return retVals.Error(0)
}
