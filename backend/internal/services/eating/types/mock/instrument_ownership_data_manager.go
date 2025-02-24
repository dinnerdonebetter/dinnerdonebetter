package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.InstrumentOwnershipDataManager = (*InstrumentOwnershipDataManagerMock)(nil)

// InstrumentOwnershipDataManagerMock is a mocked types.InstrumentOwnershipDataManager for testing.
type InstrumentOwnershipDataManagerMock struct {
	mock.Mock
}

// InstrumentOwnershipExists is a mock function.
func (m *InstrumentOwnershipDataManagerMock) InstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (bool, error) {
	returnValues := m.Called(ctx, householdInstrumentOwnershipID, householdID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetInstrumentOwnership is a mock function.
func (m *InstrumentOwnershipDataManagerMock) GetInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*types.InstrumentOwnership, error) {
	returnValues := m.Called(ctx, householdInstrumentOwnershipID, householdID)
	return returnValues.Get(0).(*types.InstrumentOwnership), returnValues.Error(1)
}

// GetInstrumentOwnerships is a mock function.
func (m *InstrumentOwnershipDataManagerMock) GetInstrumentOwnerships(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.InstrumentOwnership], error) {
	returnValues := m.Called(ctx, householdID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.InstrumentOwnership]), returnValues.Error(1)
}

// CreateInstrumentOwnership is a mock function.
func (m *InstrumentOwnershipDataManagerMock) CreateInstrumentOwnership(ctx context.Context, input *types.InstrumentOwnershipDatabaseCreationInput) (*types.InstrumentOwnership, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.InstrumentOwnership), returnValues.Error(1)
}

// UpdateInstrumentOwnership is a mock function.
func (m *InstrumentOwnershipDataManagerMock) UpdateInstrumentOwnership(ctx context.Context, updated *types.InstrumentOwnership) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveInstrumentOwnership is a mock function.
func (m *InstrumentOwnershipDataManagerMock) ArchiveInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error {
	return m.Called(ctx, householdInstrumentOwnershipID, householdID).Error(0)
}
