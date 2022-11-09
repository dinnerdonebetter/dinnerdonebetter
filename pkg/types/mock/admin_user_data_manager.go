package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.AdminUserDataManager = (*AdminUserDataManager)(nil)

// AdminUserDataManager is a mocked types.AdminUserDataManager for testing.
type AdminUserDataManager struct {
	mock.Mock
}

// UpdateUserAccountStatus is a mock function.
func (m *AdminUserDataManager) UpdateUserAccountStatus(ctx context.Context, userID string, input *types.UserAccountStatusUpdateInput) error {
	return m.Called(ctx, userID, input).Error(0)
}
