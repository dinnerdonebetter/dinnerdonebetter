package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ServiceSettingDataManager = (*ServiceSettingDataManagerMock)(nil)

// ServiceSettingDataManagerMock is a mocked types.ServiceSettingDataManager for testing.
type ServiceSettingDataManagerMock struct {
	mock.Mock
}

// ServiceSettingExists is a mock function.
func (m *ServiceSettingDataManagerMock) ServiceSettingExists(ctx context.Context, serviceSettingID string) (bool, error) {
	returnValues := m.Called(ctx, serviceSettingID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetServiceSetting is a mock function.
func (m *ServiceSettingDataManagerMock) GetServiceSetting(ctx context.Context, serviceSettingID string) (*types.ServiceSetting, error) {
	returnValues := m.Called(ctx, serviceSettingID)
	return returnValues.Get(0).(*types.ServiceSetting), returnValues.Error(1)
}

// GetRandomServiceSetting is a mock function.
func (m *ServiceSettingDataManagerMock) GetRandomServiceSetting(ctx context.Context) (*types.ServiceSetting, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*types.ServiceSetting), returnValues.Error(1)
}

// SearchForServiceSettings is a mock function.
func (m *ServiceSettingDataManagerMock) SearchForServiceSettings(ctx context.Context, query string) ([]*types.ServiceSetting, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*types.ServiceSetting), returnValues.Error(1)
}

// GetServiceSettings is a mock function.
func (m *ServiceSettingDataManagerMock) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSetting], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.ServiceSetting]), returnValues.Error(1)
}

// CreateServiceSetting is a mock function.
func (m *ServiceSettingDataManagerMock) CreateServiceSetting(ctx context.Context, input *types.ServiceSettingDatabaseCreationInput) (*types.ServiceSetting, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ServiceSetting), returnValues.Error(1)
}

// UpdateServiceSetting is a mock function.
func (m *ServiceSettingDataManagerMock) UpdateServiceSetting(ctx context.Context, updated *types.ServiceSetting) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveServiceSetting is a mock function.
func (m *ServiceSettingDataManagerMock) ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error {
	return m.Called(ctx, serviceSettingID).Error(0)
}
