package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.WebhookDataManager = (*WebhookDataManager)(nil)

// WebhookDataManager is a mocked types.WebhookDataManager for testing.
type WebhookDataManager struct {
	mock.Mock
}

// WebhookExists satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) WebhookExists(ctx context.Context, webhookID, householdID string) (bool, error) {
	args := m.Called(ctx, webhookID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhook(ctx context.Context, webhookID, householdID string) (*types.Webhook, error) {
	args := m.Called(ctx, webhookID, householdID)
	return args.Get(0).(*types.Webhook), args.Error(1)
}

// GetWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhooks(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.WebhookList, error) {
	args := m.Called(ctx, householdID, filter)
	return args.Get(0).(*types.WebhookList), args.Error(1)
}

// GetAllWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetAllWebhooks(ctx context.Context, results chan []*types.Webhook, bucketSize uint16) error {
	return m.Called(ctx, results, bucketSize).Error(0)
}

// CreateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Webhook), args.Error(1)
}

// UpdateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) UpdateWebhook(ctx context.Context, updated *types.Webhook) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) ArchiveWebhook(ctx context.Context, webhookID, householdID string) error {
	return m.Called(ctx, webhookID, householdID).Error(0)
}
