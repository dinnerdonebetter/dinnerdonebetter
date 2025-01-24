package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MaintenanceDataManager = (*MaintenanceDataManagerMock)(nil)

type MaintenanceDataManagerMock struct {
	mock.Mock
}

func (m *MaintenanceDataManagerMock) DeleteExpiredOAuth2ClientTokens(ctx context.Context) (int64, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(int64), returnValues.Error(1)
}
