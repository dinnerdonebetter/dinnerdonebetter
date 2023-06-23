package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.APIClientDataManager = (*APIClientDataManagerMock)(nil)

// APIClientDataManagerMock is a mocked types.APIClientDataManager for testing.
type APIClientDataManagerMock struct {
	mock.Mock
}

// GetAPIClientByClientID is a mock function.
func (m *APIClientDataManagerMock) GetAPIClientByClientID(ctx context.Context, clientID string) (*types.APIClient, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(*types.APIClient), args.Error(1)
}

// GetAPIClientByDatabaseID is a mock function.
func (m *APIClientDataManagerMock) GetAPIClientByDatabaseID(ctx context.Context, clientID, userID string) (*types.APIClient, error) {
	args := m.Called(ctx, clientID, userID)
	return args.Get(0).(*types.APIClient), args.Error(1)
}

// GetAPIClients is a mock function.
func (m *APIClientDataManagerMock) GetAPIClients(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.APIClient], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.APIClient]), args.Error(1)
}

// CreateAPIClient is a mock function.
func (m *APIClientDataManagerMock) CreateAPIClient(ctx context.Context, input *types.APIClientDatabaseCreationInput) (*types.APIClient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.APIClient), args.Error(1)
}

// ArchiveAPIClient is a mock function.
func (m *APIClientDataManagerMock) ArchiveAPIClient(ctx context.Context, clientID, householdID string) error {
	return m.Called(ctx, clientID, householdID).Error(0)
}
