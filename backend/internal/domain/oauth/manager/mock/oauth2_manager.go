package managermock

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth/manager"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

var _ manager.OAuth2Manager = (*OAuth2Manager)(nil)

// OAuth2Manager is a mock implementation of the OAuth2Manager interface.
type OAuth2Manager struct {
	mock.Mock
}

// CreateOAuth2Client is a mock function.
func (m *OAuth2Manager) CreateOAuth2Client(ctx context.Context, input *oauth.OAuth2ClientCreationRequestInput) (*oauth.OAuth2Client, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*oauth.OAuth2Client), returnValues.Error(1)
}

// GetOAuth2Client is a mock function.
func (m *OAuth2Manager) GetOAuth2Client(ctx context.Context, oauth2ClientID string) (*oauth.OAuth2Client, error) {
	returnValues := m.Called(ctx, oauth2ClientID)
	return returnValues.Get(0).(*oauth.OAuth2Client), returnValues.Error(1)
}

// GetOAuth2Clients is a mock function.
func (m *OAuth2Manager) GetOAuth2Clients(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[oauth.OAuth2Client], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[oauth.OAuth2Client]), returnValues.Error(1)
}

// ArchiveOAuth2Client is a mock function.
func (m *OAuth2Manager) ArchiveOAuth2Client(ctx context.Context, oauth2ClientID string) error {
	return m.Called(ctx, oauth2ClientID).Error(0)
}

// NewOAuth2Manager creates a new mock OAuth2Manager.
func NewOAuth2Manager() *OAuth2Manager {
	return &OAuth2Manager{}
}
