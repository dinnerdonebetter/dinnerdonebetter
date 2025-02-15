package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ServiceSettingConfigurationDataManager = (*ServiceSettingConfigurationDataManagerMock)(nil)

// ServiceSettingConfigurationDataManagerMock is a mocked types.ServiceSettingConfigurationDataManager for testing.
type ServiceSettingConfigurationDataManagerMock struct {
	mock.Mock
}

// GetServiceSettingConfigurationForUserByName is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetServiceSettingConfigurationForUserByName(ctx context.Context, userID, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, userID, serviceSettingConfigurationID)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// GetServiceSettingConfigurationForHouseholdByName is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetServiceSettingConfigurationForHouseholdByName(ctx context.Context, householdID, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, householdID, serviceSettingConfigurationID)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// GetServiceSettingConfigurationsForUser is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetServiceSettingConfigurationsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ServiceSettingConfiguration]), returnValues.Error(1)
}

// GetServiceSettingConfigurationsForHousehold is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetServiceSettingConfigurationsForHousehold(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	returnValues := m.Called(ctx, householdID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ServiceSettingConfiguration]), returnValues.Error(1)
}

// ServiceSettingConfigurationExists is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) ServiceSettingConfigurationExists(ctx context.Context, serviceSettingConfigurationID string) (bool, error) {
	returnValues := m.Called(ctx, serviceSettingConfigurationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, serviceSettingConfigurationID)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// GetRandomServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) GetRandomServiceSettingConfiguration(ctx context.Context) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// CreateServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) CreateServiceSettingConfiguration(ctx context.Context, input *types.ServiceSettingConfigurationDatabaseCreationInput) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// UpdateServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) UpdateServiceSettingConfiguration(ctx context.Context, updated *types.ServiceSettingConfiguration) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManagerMock) ArchiveServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) error {
	return m.Called(ctx, serviceSettingConfigurationID).Error(0)
}
