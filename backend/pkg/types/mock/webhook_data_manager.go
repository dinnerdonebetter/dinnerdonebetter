package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.WebhookDataManager = (*WebhookDataManagerMock)(nil)

// WebhookDataManagerMock is a mocked types.WebhookDataManager for testing.
type WebhookDataManagerMock struct {
	mock.Mock
}

// WebhookExists satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) WebhookExists(ctx context.Context, webhookID, accountID string) (bool, error) {
	returnValues := m.Called(ctx, webhookID, accountID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetWebhook satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhook(ctx context.Context, webhookID, accountID string) (*types.Webhook, error) {
	returnValues := m.Called(ctx, webhookID, accountID)
	return returnValues.Get(0).(*types.Webhook), returnValues.Error(1)
}

// GetWebhooks satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhooks(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Webhook], error) {
	returnValues := m.Called(ctx, accountID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Webhook]), returnValues.Error(1)
}

// GetWebhooksForAccountAndEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhooksForAccountAndEvent(ctx context.Context, accountID, eventType string) ([]*types.Webhook, error) {
	returnValues := m.Called(ctx, accountID, eventType)
	return returnValues.Get(0).([]*types.Webhook), returnValues.Error(1)
}

// GetAllWebhooks satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetAllWebhooks(ctx context.Context, results chan []*types.Webhook, bucketSize uint16) error {
	return m.Called(ctx, results, bucketSize).Error(0)
}

// CreateWebhook satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.Webhook), returnValues.Error(1)
}

// ArchiveWebhook satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) ArchiveWebhook(ctx context.Context, webhookID, accountID string) error {
	return m.Called(ctx, webhookID, accountID).Error(0)
}

// AddWebhookTriggerEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) AddWebhookTriggerEvent(ctx context.Context, accountID string, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	returnValues := m.Called(ctx, accountID, input)
	return returnValues.Get(0).(*types.WebhookTriggerEvent), returnValues.Error(1)
}

// ArchiveWebhookTriggerEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) ArchiveWebhookTriggerEvent(ctx context.Context, webhookID, webhookTriggerEventID string) error {
	return m.Called(ctx, webhookID, webhookTriggerEventID).Error(0)
}
