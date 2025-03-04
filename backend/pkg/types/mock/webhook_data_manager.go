package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.WebhookDataManager = (*WebhookDataManagerMock)(nil)

// WebhookDataManagerMock is a mocked types.WebhookDataManager for testing.
type WebhookDataManagerMock struct {
	mock.Mock
}

// WebhookExists satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) WebhookExists(ctx context.Context, webhookID, householdID string) (bool, error) {
	returnValues := m.Called(ctx, webhookID, householdID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetWebhook satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhook(ctx context.Context, webhookID, householdID string) (*types.Webhook, error) {
	returnValues := m.Called(ctx, webhookID, householdID)
	return returnValues.Get(0).(*types.Webhook), returnValues.Error(1)
}

// GetWebhooks satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhooks(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Webhook], error) {
	returnValues := m.Called(ctx, householdID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.Webhook]), returnValues.Error(1)
}

// GetWebhooksForHouseholdAndEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) GetWebhooksForHouseholdAndEvent(ctx context.Context, householdID, eventType string) ([]*types.Webhook, error) {
	returnValues := m.Called(ctx, householdID, eventType)
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
func (m *WebhookDataManagerMock) ArchiveWebhook(ctx context.Context, webhookID, householdID string) error {
	return m.Called(ctx, webhookID, householdID).Error(0)
}

// AddWebhookTriggerEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) AddWebhookTriggerEvent(ctx context.Context, householdID string, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	returnValues := m.Called(ctx, householdID, input)
	return returnValues.Get(0).(*types.WebhookTriggerEvent), returnValues.Error(1)
}

// ArchiveWebhookTriggerEvent satisfies our WebhookDataManagerMock interface.
func (m *WebhookDataManagerMock) ArchiveWebhookTriggerEvent(ctx context.Context, webhookID, webhookTriggerEventID string) error {
	return m.Called(ctx, webhookID, webhookTriggerEventID).Error(0)
}
