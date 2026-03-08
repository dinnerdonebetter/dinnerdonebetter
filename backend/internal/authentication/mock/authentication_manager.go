package mockauthn

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"

	"github.com/stretchr/testify/mock"
)

var _ authentication.Manager = (*Manager)(nil)

// Manager is a mock implementation of the authentication.Manager interface.
type Manager struct {
	mock.Mock
}

// ProcessLogin is a mock method.
func (m *Manager) ProcessLogin(ctx context.Context, adminOnly bool, loginData *auth.UserLoginInput) (*auth.TokenResponse, error) {
	args := m.Called(ctx, adminOnly, loginData)
	return args.Get(0).(*auth.TokenResponse), args.Error(1)
}

// ProcessPasskeyLogin is a mock method.
func (m *Manager) ProcessPasskeyLogin(ctx context.Context, userID, desiredAccountID string) (*auth.TokenResponse, error) {
	args := m.Called(ctx, userID, desiredAccountID)
	return args.Get(0).(*auth.TokenResponse), args.Error(1)
}

// ExchangeTokenForUser is a mock method.
func (m *Manager) ExchangeTokenForUser(ctx context.Context, refreshToken, desiredAccountID string) (*auth.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*auth.TokenResponse), args.Error(1)
}
