package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationInstrumentDataManager = (*ValidPreparationInstrumentDataManagerMock)(nil)

// ValidPreparationInstrumentDataManagerMock is a mocked types.ValidPreparationInstrumentDataManager for testing.
type ValidPreparationInstrumentDataManagerMock struct {
	mock.Mock
}

// ValidPreparationInstrumentExists is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)
	return returnValues.Get(0).(*types.ValidPreparationInstrument), returnValues.Error(1)
}

// GetValidPreparationInstruments is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) GetValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationInstrument]), returnValues.Error(1)
}

// GetValidPreparationInstrumentsForPreparation is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationInstrument]), returnValues.Error(1)
}

// GetValidPreparationInstrumentsForInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationInstrument], error) {
	returnValues := m.Called(ctx, instrumentID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationInstrument]), returnValues.Error(1)
}

// CreateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidPreparationInstrument), returnValues.Error(1)
}

// UpdateValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationInstrument is a mock function.
func (m *ValidPreparationInstrumentDataManagerMock) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	return m.Called(ctx, validPreparationInstrumentID).Error(0)
}
