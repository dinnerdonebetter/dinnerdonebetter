package types

import (
	"context"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// WebhookCreatedCustomerEventType indicates a webhook was created.
	WebhookCreatedCustomerEventType CustomerEventType = "webhook_created"
	// WebhookArchivedCustomerEventType indicates a webhook was archived.
	WebhookArchivedCustomerEventType CustomerEventType = "webhook_archived"
)

type (
	// Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event.
	Webhook struct {
		_ struct{}

		CreatedAt          time.Time              `json:"createdAt"`
		ArchivedAt         *time.Time             `json:"archivedAt"`
		LastUpdatedAt      *time.Time             `json:"lastUpdatedAt"`
		Name               string                 `json:"name"`
		URL                string                 `json:"url"`
		Method             string                 `json:"method"`
		ID                 string                 `json:"id"`
		BelongsToHousehold string                 `json:"belongsToHousehold"`
		ContentType        string                 `json:"contentType"`
		Events             []*WebhookTriggerEvent `json:"events"`
	}

	// WebhookTriggerEvent represents a webhook trigger event.
	WebhookTriggerEvent struct {
		_ struct{}

		CreatedAt        time.Time  `json:"createdAt"`
		ArchivedAt       *time.Time `json:"archivedAt"`
		ID               string     `json:"id"`
		BelongsToWebhook string     `json:"belongsToWebhook"`
		TriggerEvent     string     `json:"triggerEvent"`
	}

	// WebhookCreationRequestInput represents what a User could set as input for creating a webhook.
	WebhookCreationRequestInput struct {
		_ struct{}

		Name        string   `json:"name"`
		ContentType string   `json:"contentType"`
		URL         string   `json:"url"`
		Method      string   `json:"method"`
		Events      []string `json:"events"`
	}

	// WebhookDatabaseCreationInput is used for creating a webhook.
	WebhookDatabaseCreationInput struct {
		_ struct{}

		ID                 string
		Name               string
		ContentType        string
		URL                string
		Method             string
		BelongsToHousehold string
		Events             []*WebhookTriggerEventDatabaseCreationInput
	}

	// WebhookTriggerEventDatabaseCreationInput is used for creating a webhook trigger event.
	WebhookTriggerEventDatabaseCreationInput struct {
		_ struct{}

		ID               string
		BelongsToWebhook string
		TriggerEvent     string
	}

	// WebhookExecutionRequest represents a webhook listener, an endpoint to send an HTTP request to upon an event.
	WebhookExecutionRequest struct {
		_            struct{}
		Payload      any    `json:"payload"`
		WebhookID    string `json:"webhookID"`
		HouseholdID  string `json:"householdID"`
		TriggerEvent string `json:"triggerEvent"`
	}

	// WebhookDataManager describes a structure capable of storing webhooks.
	WebhookDataManager interface {
		WebhookExists(ctx context.Context, webhookID, householdID string) (bool, error)
		GetWebhook(ctx context.Context, webhookID, householdID string) (*Webhook, error)
		GetWebhooks(ctx context.Context, householdID string, filter *QueryFilter) (*QueryFilteredResult[Webhook], error)
		CreateWebhook(ctx context.Context, input *WebhookDatabaseCreationInput) (*Webhook, error)
		ArchiveWebhook(ctx context.Context, webhookID, householdID string) error
	}

	// WebhookDataService describes a structure capable of serving traffic related to webhooks.
	WebhookDataService interface {
		ListWebhooksHandler(http.ResponseWriter, *http.Request)
		CreateWebhookHandler(http.ResponseWriter, *http.Request)
		ReadWebhookHandler(http.ResponseWriter, *http.Request)
		ArchiveWebhookHandler(http.ResponseWriter, *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*WebhookCreationRequestInput)(nil)

// ValidateWithContext validates a WebhookCreationRequestInput.
func (w *WebhookCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, w,
		validation.Field(&w.Name, validation.Required),
		validation.Field(&w.URL, validation.Required, is.URL),
		validation.Field(&w.Method, validation.Required, validation.In(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete)),
		validation.Field(&w.ContentType, validation.Required, validation.In("application/json", "application/xml")),
		validation.Field(&w.Events, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*WebhookDatabaseCreationInput)(nil)

// ValidateWithContext validates a WebhookDatabaseCreationInput.
func (w *WebhookDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, w,
		validation.Field(&w.ID, validation.Required),
		validation.Field(&w.Name, validation.Required),
		validation.Field(&w.URL, validation.Required, is.URL),
		validation.Field(&w.Method, validation.Required, validation.In(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete)),
		validation.Field(&w.ContentType, validation.Required, validation.In("application/json", "application/xml")),
		validation.Field(&w.Events, validation.Required),
		validation.Field(&w.BelongsToHousehold, validation.Required),
	)
}
