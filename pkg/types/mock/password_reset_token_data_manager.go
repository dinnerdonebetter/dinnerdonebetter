package mocktypes

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
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
