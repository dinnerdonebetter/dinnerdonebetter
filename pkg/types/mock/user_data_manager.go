package mocktypes

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.UserDataManager = (*UserDataManager)(nil)

// UserDataManager is a mocked types.UserDataManager for testing.
type UserDataManager struct {
	mock.Mock
}

// GetUserByEmailAddressVerificationToken is a mock function.
func (m *UserDataManager) GetUserByEmailAddressVerificationToken(ctx context.Context, token string) (*types.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*types.User), args.Error(1)
}

// MarkUserEmailAddressAsVerified is a mock function.
func (m *UserDataManager) MarkUserEmailAddressAsVerified(ctx context.Context, userID, token string) error {
	return m.Called(ctx, userID, token).Error(0)
}

// UserHasStatus is a mock function.
func (m *UserDataManager) UserHasStatus(ctx context.Context, userID string, statuses ...string) (bool, error) {
	args := m.Called(ctx, userID, statuses)

	return args.Bool(0), args.Error(1)
}

// GetUser is a mock function.
func (m *UserDataManager) GetUser(ctx context.Context, userID string) (*types.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*types.User), args.Error(1)
}

// GetUserWithUnverifiedTwoFactorSecret is a mock function.
func (m *UserDataManager) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID string) (*types.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*types.User), args.Error(1)
}

// GetUserByEmail is a mock function.
func (m *UserDataManager) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*types.User), args.Error(1)
}

// MarkUserTwoFactorSecretAsVerified is a mock function.
func (m *UserDataManager) MarkUserTwoFactorSecretAsVerified(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

// MarkUserTwoFactorSecretAsUnverified is a mock function.
func (m *UserDataManager) MarkUserTwoFactorSecretAsUnverified(ctx context.Context, userID, newSecret string) error {
	return m.Called(ctx, userID, newSecret).Error(0)
}

// GetUserByUsername is a mock function.
func (m *UserDataManager) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*types.User), args.Error(1)
}

// GetAdminUserByUsername is a mock function.
func (m *UserDataManager) GetAdminUserByUsername(ctx context.Context, username string) (*types.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*types.User), args.Error(1)
}

// SearchForUsersByUsername is a mock function.
func (m *UserDataManager) SearchForUsersByUsername(ctx context.Context, usernameQuery string) ([]*types.User, error) {
	args := m.Called(ctx, usernameQuery)
	return args.Get(0).([]*types.User), args.Error(1)
}

// GetUsers is a mock function.
func (m *UserDataManager) GetUsers(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.User], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.User]), args.Error(1)
}

// CreateUser is a mock function.
func (m *UserDataManager) CreateUser(ctx context.Context, input *types.UserDatabaseCreationInput) (*types.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.User), args.Error(1)
}

// UpdateUser is a mock function.
func (m *UserDataManager) UpdateUser(ctx context.Context, updated *types.User) error {
	return m.Called(ctx, updated).Error(0)
}

// UpdateUserPassword is a mock function.
func (m *UserDataManager) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	return m.Called(ctx, userID, newHash).Error(0)
}

// ArchiveUser is a mock function.
func (m *UserDataManager) ArchiveUser(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

// GetEmailAddressVerificationTokenForUser is a mock function.
func (m *UserDataManager) GetEmailAddressVerificationTokenForUser(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
