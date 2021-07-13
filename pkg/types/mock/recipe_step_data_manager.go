package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepDataManager = (*RecipeStepDataManager)(nil)

// RecipeStepDataManager is a mocked types.RecipeStepDataManager for testing.
type RecipeStepDataManager struct {
	mock.Mock
}

// RecipeStepExists is a mock function.
func (m *RecipeStepDataManager) RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStep is a mock function.
func (m *RecipeStepDataManager) GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*types.RecipeStep, error) {
	args := m.Called(ctx, recipeID, recipeStepID)
	return args.Get(0).(*types.RecipeStep), args.Error(1)
}

// GetAllRecipeStepsCount is a mock function.
func (m *RecipeStepDataManager) GetAllRecipeStepsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipeSteps is a mock function.
func (m *RecipeStepDataManager) GetAllRecipeSteps(ctx context.Context, results chan []*types.RecipeStep, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetRecipeSteps is a mock function.
func (m *RecipeStepDataManager) GetRecipeSteps(ctx context.Context, recipeID uint64, filter *types.QueryFilter) (*types.RecipeStepList, error) {
	args := m.Called(ctx, recipeID, filter)
	return args.Get(0).(*types.RecipeStepList), args.Error(1)
}

// GetRecipeStepsWithIDs is a mock function.
func (m *RecipeStepDataManager) GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]*types.RecipeStep, error) {
	args := m.Called(ctx, recipeID, limit, ids)
	return args.Get(0).([]*types.RecipeStep), args.Error(1)
}

// CreateRecipeStep is a mock function.
func (m *RecipeStepDataManager) CreateRecipeStep(ctx context.Context, input *types.RecipeStepCreationInput, createdByUser uint64) (*types.RecipeStep, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.RecipeStep), args.Error(1)
}

// UpdateRecipeStep is a mock function.
func (m *RecipeStepDataManager) UpdateRecipeStep(ctx context.Context, updated *types.RecipeStep, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveRecipeStep is a mock function.
func (m *RecipeStepDataManager) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID, archivedBy uint64) error {
	return m.Called(ctx, recipeID, recipeStepID, archivedBy).Error(0)
}

// GetAuditLogEntriesForRecipeStep is a mock function.
func (m *RecipeStepDataManager) GetAuditLogEntriesForRecipeStep(ctx context.Context, recipeStepID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, recipeStepID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
