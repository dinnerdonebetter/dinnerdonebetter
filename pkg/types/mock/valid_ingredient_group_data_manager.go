package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientGroupDataManager = (*ValidIngredientGroupDataManager)(nil)

// ValidIngredientGroupDataManager is a mocked types.ValidIngredientGroupDataManager for testing.
type ValidIngredientGroupDataManager struct {
	mock.Mock
}

// ValidIngredientGroupExists is a mock method.
func (m *ValidIngredientGroupDataManager) ValidIngredientGroupExists(ctx context.Context, validIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManager) GetValidIngredientGroup(ctx context.Context, validIngredientID string) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

// GetValidIngredientGroups is a mock method.
func (m *ValidIngredientGroupDataManager) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientGroup], error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).(*types.QueryFilteredResult[types.ValidIngredientGroup]), returnValues.Error(1)
}

// SearchForValidIngredientGroups is a mock method.
func (m *ValidIngredientGroupDataManager) SearchForValidIngredientGroups(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, query, filter)

	return returnValues.Get(0).([]*types.ValidIngredientGroup), returnValues.Error(1)
}

// CreateValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManager) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupDatabaseCreationInput) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

// UpdateValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManager) UpdateValidIngredientGroup(ctx context.Context, updated *types.ValidIngredientGroup) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManager) ArchiveValidIngredientGroup(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}
