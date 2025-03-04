package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientPreparationDataManager = (*ValidIngredientPreparationDataManagerMock)(nil)

// ValidIngredientPreparationDataManagerMock is a mocked types.ValidIngredientPreparationDataManager for testing.
type ValidIngredientPreparationDataManagerMock struct {
	mock.Mock
}

// ValidIngredientPreparationExists is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID)
	return returnValues.Get(0).(*types.ValidIngredientPreparation), returnValues.Error(1)
}

// GetValidIngredientPreparations is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) GetValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForIngredient is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) GetValidIngredientPreparationsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, ingredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForPreparation is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) GetValidIngredientPreparationsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientPreparation]), returnValues.Error(1)
}

// GetValidIngredientPreparationsForIngredientNameQuery is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) GetValidIngredientPreparationsForIngredientNameQuery(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientPreparation], error) {
	returnValues := m.Called(ctx, preparationID, query, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientPreparation]), returnValues.Error(1)
}

// CreateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationDatabaseCreationInput) (*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidIngredientPreparation), returnValues.Error(1)
}

// UpdateValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientPreparation is a mock function.
func (m *ValidIngredientPreparationDataManagerMock) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	return m.Called(ctx, validIngredientPreparationID).Error(0)
}
