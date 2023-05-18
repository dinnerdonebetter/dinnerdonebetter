package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ServiceSettingConfigurationDataManager = (*ServiceSettingConfigurationDataManager)(nil)

// ServiceSettingConfigurationDataManager is a mocked types.ServiceSettingConfigurationDataManager for testing.
type ServiceSettingConfigurationDataManager struct {
	mock.Mock
}

// GetServiceSettingConfigurationForUserByName is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetServiceSettingConfigurationForUserByName(ctx context.Context, userID, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, userID, serviceSettingConfigurationID)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// GetServiceSettingConfigurationForHouseholdByName is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetServiceSettingConfigurationForHouseholdByName(ctx context.Context, householdID, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	returnValues := m.Called(ctx, householdID, serviceSettingConfigurationID)
	return returnValues.Get(0).(*types.ServiceSettingConfiguration), returnValues.Error(1)
}

// GetServiceSettingConfigurationsForUser is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetServiceSettingConfigurationsForUser(ctx context.Context, userID string) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	returnValues := m.Called(ctx, userID)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.ServiceSettingConfiguration]), returnValues.Error(1)
}

// GetServiceSettingConfigurationsForHousehold is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetServiceSettingConfigurationsForHousehold(ctx context.Context, householdID string) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	returnValues := m.Called(ctx, householdID)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.ServiceSettingConfiguration]), returnValues.Error(1)
}

// ServiceSettingConfigurationExists is a mock function.
func (m *ServiceSettingConfigurationDataManager) ServiceSettingConfigurationExists(ctx context.Context, serviceSettingConfigurationID string) (bool, error) {
	args := m.Called(ctx, serviceSettingConfigurationID)
	return args.Bool(0), args.Error(1)
}

// GetServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	args := m.Called(ctx, serviceSettingConfigurationID)
	return args.Get(0).(*types.ServiceSettingConfiguration), args.Error(1)
}

// GetRandomServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManager) GetRandomServiceSettingConfiguration(ctx context.Context) (*types.ServiceSettingConfiguration, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ServiceSettingConfiguration), args.Error(1)
}

// CreateServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManager) CreateServiceSettingConfiguration(ctx context.Context, input *types.ServiceSettingConfigurationDatabaseCreationInput) (*types.ServiceSettingConfiguration, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ServiceSettingConfiguration), args.Error(1)
}

// UpdateServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManager) UpdateServiceSettingConfiguration(ctx context.Context, updated *types.ServiceSettingConfiguration) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveServiceSettingConfiguration is a mock function.
func (m *ServiceSettingConfigurationDataManager) ArchiveServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) error {
	return m.Called(ctx, serviceSettingConfigurationID).Error(0)
}
