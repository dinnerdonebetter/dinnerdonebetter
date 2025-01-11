package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.OAuth2ClientDataManager = (*OAuth2ClientDataManagerMock)(nil)

// OAuth2ClientDataManagerMock is a mocked types.OAuth2ClientDataManager for testing.
type OAuth2ClientDataManagerMock struct {
	mock.Mock
}

// GetOAuth2ClientByClientID is a mock function.
func (m *OAuth2ClientDataManagerMock) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	returnValues := m.Called(ctx, clientID)
	return returnValues.Get(0).(*types.OAuth2Client), returnValues.Error(1)
}

// GetOAuth2ClientByDatabaseID is a mock function.
func (m *OAuth2ClientDataManagerMock) GetOAuth2ClientByDatabaseID(ctx context.Context, clientID string) (*types.OAuth2Client, error) {
	returnValues := m.Called(ctx, clientID)
	return returnValues.Get(0).(*types.OAuth2Client), returnValues.Error(1)
}

// GetOAuth2Clients is a mock function.
func (m *OAuth2ClientDataManagerMock) GetOAuth2Clients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.OAuth2Client], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*types.QueryFilteredResult[types.OAuth2Client]), returnValues.Error(1)
}

// CreateOAuth2Client is a mock function.
func (m *OAuth2ClientDataManagerMock) CreateOAuth2Client(ctx context.Context, input *types.OAuth2ClientDatabaseCreationInput) (*types.OAuth2Client, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.OAuth2Client), returnValues.Error(1)
}

// ArchiveOAuth2Client is a mock function.
func (m *OAuth2ClientDataManagerMock) ArchiveOAuth2Client(ctx context.Context, clientID string) error {
	return m.Called(ctx, clientID).Error(0)
}
