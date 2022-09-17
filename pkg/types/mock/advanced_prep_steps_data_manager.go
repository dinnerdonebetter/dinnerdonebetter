package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AdvancedPrepStepDataManager = (*AdvancedPrepStepDataManager)(nil)

// AdvancedPrepStepDataManager is a mocked types.AdvancedPrepStepDataManager for testing.
type AdvancedPrepStepDataManager struct {
	mock.Mock
}

// AdvancedPrepStepExists is a mock function.
func (m *AdvancedPrepStepDataManager) AdvancedPrepStepExists(ctx context.Context, advancedPrepStepID string) (bool, error) {
	args := m.Called(ctx, advancedPrepStepID)
	return args.Bool(0), args.Error(1)
}

// GetAdvancedPrepStep is a mock function.
func (m *AdvancedPrepStepDataManager) GetAdvancedPrepStep(ctx context.Context, advancedPrepStepID string) (*types.AdvancedPrepStep, error) {
	args := m.Called(ctx, advancedPrepStepID)
	return args.Get(0).(*types.AdvancedPrepStep), args.Error(1)
}

// GetAdvancedPrepSteps is a mock function.
func (m *AdvancedPrepStepDataManager) GetAdvancedPrepSteps(ctx context.Context, filter *types.QueryFilter) (*types.AdvancedPrepStepList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.AdvancedPrepStepList), args.Error(1)
}

// CreateAdvancedPrepStep is a mock function.
func (m *AdvancedPrepStepDataManager) CreateAdvancedPrepStep(ctx context.Context, input *types.AdvancedPrepStepDatabaseCreationInput) (*types.AdvancedPrepStep, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.AdvancedPrepStep), args.Error(1)
}

// MarkAdvancedPrepStepAsComplete is a mock function.
func (m *AdvancedPrepStepDataManager) MarkAdvancedPrepStepAsComplete(ctx context.Context, advancedPrepStepID string) error {
	return m.Called(ctx, advancedPrepStepID).Error(0)
}
