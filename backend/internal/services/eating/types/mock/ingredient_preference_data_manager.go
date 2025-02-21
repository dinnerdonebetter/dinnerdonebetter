package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.IngredientPreferenceDataManager = (*IngredientPreferenceDataManagerMock)(nil)

// IngredientPreferenceDataManagerMock is a mocked types.IngredientPreferenceDataManager for testing.
type IngredientPreferenceDataManagerMock struct {
	mock.Mock
}

// IngredientPreferenceExists is a mock function.
func (m *IngredientPreferenceDataManagerMock) IngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (bool, error) {
	returnValues := m.Called(ctx, userIngredientPreferenceID, userID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetIngredientPreference is a mock function.
func (m *IngredientPreferenceDataManagerMock) GetIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*types.IngredientPreference, error) {
	returnValues := m.Called(ctx, userIngredientPreferenceID, userID)
	return returnValues.Get(0).(*types.IngredientPreference), returnValues.Error(1)
}

// GetIngredientPreferences is a mock function.
func (m *IngredientPreferenceDataManagerMock) GetIngredientPreferences(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.IngredientPreference], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.IngredientPreference]), returnValues.Error(1)
}

// CreateIngredientPreference is a mock function.
func (m *IngredientPreferenceDataManagerMock) CreateIngredientPreference(ctx context.Context, input *types.IngredientPreferenceDatabaseCreationInput) ([]*types.IngredientPreference, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).([]*types.IngredientPreference), returnValues.Error(1)
}

// UpdateIngredientPreference is a mock function.
func (m *IngredientPreferenceDataManagerMock) UpdateIngredientPreference(ctx context.Context, updated *types.IngredientPreference) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveIngredientPreference is a mock function.
func (m *IngredientPreferenceDataManagerMock) ArchiveIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error {
	return m.Called(ctx, userIngredientPreferenceID, userID).Error(0)
}
