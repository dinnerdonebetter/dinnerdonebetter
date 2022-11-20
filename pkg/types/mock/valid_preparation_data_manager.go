package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.ValidPreparationDataManager = (*ValidPreparationDataManager)(nil)

// ValidPreparationDataManager is a mocked types.ValidPreparationDataManager for testing.
type ValidPreparationDataManager struct {
	mock.Mock
}

// ValidPreparationExists is a mock function.
func (m *ValidPreparationDataManager) ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparation is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// GetRandomValidPreparation is a mock function.
func (m *ValidPreparationDataManager) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// SearchForValidPreparations is a mock function.
func (m *ValidPreparationDataManager) SearchForValidPreparations(ctx context.Context, query string) ([]*types.ValidPreparation, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidPreparation), args.Error(1)
}

// GetValidPreparations is a mock function.
func (m *ValidPreparationDataManager) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparation], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidPreparation]), args.Error(1)
}

// CreateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationDatabaseCreationInput) (*types.ValidPreparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// UpdateValidPreparation is a mock function.
func (m *ValidPreparationDataManager) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparation is a mock function.
func (m *ValidPreparationDataManager) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
