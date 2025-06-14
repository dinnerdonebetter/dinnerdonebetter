package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.AccountDataManager = (*AccountDataManagerMock)(nil)

// AccountDataManagerMock is a mocked types.AccountDataManager for testing.
type AccountDataManagerMock struct {
	mock.Mock
}

// AccountExists is a mock function.
func (m *AccountDataManagerMock) AccountExists(ctx context.Context, accountID, userID string) (bool, error) {
	returnValues := m.Called(ctx, accountID, userID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetAccount is a mock function.
func (m *AccountDataManagerMock) GetAccount(ctx context.Context, accountID string) (*types.Account, error) {
	returnValues := m.Called(ctx, accountID)
	return returnValues.Get(0).(*types.Account), returnValues.Error(1)
}

// GetAllAccounts is a mock function.
func (m *AccountDataManagerMock) GetAllAccounts(ctx context.Context, results chan []*types.Account, bucketSize uint16) error {
	return m.Called(ctx, results, bucketSize).Error(0)
}

// GetAccounts is a mock function.
func (m *AccountDataManagerMock) GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Account], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Account]), returnValues.Error(1)
}

// CreateAccount is a mock function.
func (m *AccountDataManagerMock) CreateAccount(ctx context.Context, input *types.AccountDatabaseCreationInput) (*types.Account, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.Account), returnValues.Error(1)
}

// UpdateAccount is a mock function.
func (m *AccountDataManagerMock) UpdateAccount(ctx context.Context, updated *types.Account) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveAccount is a mock function.
func (m *AccountDataManagerMock) ArchiveAccount(ctx context.Context, accountID, userID string) error {
	return m.Called(ctx, accountID, userID).Error(0)
}
