package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.PasswordResetTokenDataManager = (*PasswordResetTokenDataManager)(nil)

// PasswordResetTokenDataManager is a mocked types.PasswordResetTokenDataManager for testing.
type PasswordResetTokenDataManager struct {
	mock.Mock
}

// GetPasswordResetTokenByToken implements our interface requirements.
func (m *PasswordResetTokenDataManager) GetPasswordResetTokenByToken(ctx context.Context, passwordResetTokenID string) (*types.PasswordResetToken, error) {
	args := m.Called(ctx, passwordResetTokenID)

	return args.Get(0).(*types.PasswordResetToken), args.Error(1)
}

// GetTotalPasswordResetTokenCount implements our interface requirements.
func (m *PasswordResetTokenDataManager) GetTotalPasswordResetTokenCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)

	return args.Get(0).(uint64), args.Error(1)
}

// CreatePasswordResetToken implements our interface requirements.
func (m *PasswordResetTokenDataManager) CreatePasswordResetToken(ctx context.Context, input *types.PasswordResetTokenDatabaseCreationInput) (*types.PasswordResetToken, error) {
	args := m.Called(ctx, input)

	return args.Get(0).(*types.PasswordResetToken), args.Error(1)
}

// RedeemPasswordResetToken implements our interface requirements.
func (m *PasswordResetTokenDataManager) RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error {
	args := m.Called(ctx, passwordResetTokenID)

	return args.Error(0)
}
