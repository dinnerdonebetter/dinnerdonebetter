package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MaintenanceDataManagerMock struct {
	mock.Mock
}

func (m *MaintenanceDataManagerMock) DeleteExpiredOAuth2ClientTokens(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
