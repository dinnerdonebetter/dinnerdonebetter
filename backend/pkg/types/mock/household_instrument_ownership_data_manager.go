package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.AccountInstrumentOwnershipDataManager = (*AccountInstrumentOwnershipDataManagerMock)(nil)

// AccountInstrumentOwnershipDataManagerMock is a mocked types.AccountInstrumentOwnershipDataManager for testing.
type AccountInstrumentOwnershipDataManagerMock struct {
	mock.Mock
}

// AccountInstrumentOwnershipExists is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) AccountInstrumentOwnershipExists(ctx context.Context, accountInstrumentOwnershipID, accountID string) (bool, error) {
	returnValues := m.Called(ctx, accountInstrumentOwnershipID, accountID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetAccountInstrumentOwnership is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) GetAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) (*types.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, accountInstrumentOwnershipID, accountID)
	return returnValues.Get(0).(*types.AccountInstrumentOwnership), returnValues.Error(1)
}

// GetAccountInstrumentOwnerships is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) GetAccountInstrumentOwnerships(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInstrumentOwnership], error) {
	returnValues := m.Called(ctx, accountID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.AccountInstrumentOwnership]), returnValues.Error(1)
}

// CreateAccountInstrumentOwnership is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) CreateAccountInstrumentOwnership(ctx context.Context, input *types.AccountInstrumentOwnershipDatabaseCreationInput) (*types.AccountInstrumentOwnership, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.AccountInstrumentOwnership), returnValues.Error(1)
}

// UpdateAccountInstrumentOwnership is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) UpdateAccountInstrumentOwnership(ctx context.Context, updated *types.AccountInstrumentOwnership) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveAccountInstrumentOwnership is a mock function.
func (m *AccountInstrumentOwnershipDataManagerMock) ArchiveAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) error {
	return m.Called(ctx, accountInstrumentOwnershipID, accountID).Error(0)
}
