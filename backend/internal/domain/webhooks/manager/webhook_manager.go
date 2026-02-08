package manager

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "webhook_data_manager"
)

var _ WebhookDataManager = (*webhookManager)(nil)

type webhookManager struct {
	tracer tracing.Tracer
	logger logging.Logger
	repo   webhooks.Repository
}

// NewWebhookDataManager returns a new WebhookDataManager that delegates to the webhooks repository.
func NewWebhookDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	repo webhooks.Repository,
) WebhookDataManager {
	return &webhookManager{
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger: logging.EnsureLogger(logger).WithName(o11yName),
		repo:   repo,
	}
}

func (m *webhookManager) WebhookExists(ctx context.Context, webhookID, accountID string) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.WebhookExists(ctx, webhookID, accountID)
}

func (m *webhookManager) CreateWebhook(ctx context.Context, userID, accountID string, input *webhooks.WebhookCreationRequestInput) (*webhooks.Webhook, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, internalerrors.ErrNilInputParameter
	}
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating webhook creation input")
	}

	webhookID := identifiers.New()
	dbInput := &webhooks.WebhookDatabaseCreationInput{
		ID:               webhookID,
		Name:             input.Name,
		ContentType:      input.ContentType,
		URL:              input.URL,
		Method:           input.Method,
		CreatedByUser:    userID,
		BelongsToAccount: accountID,
		TriggerConfigs:   nil,
	}
	for _, ev := range input.Events {
		triggerEventID := ev.ID
		if triggerEventID == "" {
			catalogInput := &webhooks.WebhookTriggerEventDatabaseCreationInput{
				ID:          identifiers.New(),
				Name:        ev.Name,
				Description: ev.Description,
			}
			created, err := m.repo.CreateWebhookTriggerEvent(ctx, catalogInput)
			if err != nil {
				return nil, observability.PrepareAndLogError(err, m.logger, span, "creating catalog trigger event")
			}
			triggerEventID = created.ID
		}
		dbInput.TriggerConfigs = append(dbInput.TriggerConfigs, &webhooks.WebhookTriggerConfigDatabaseCreationInput{
			ID:               identifiers.New(),
			BelongsToWebhook: webhookID,
			TriggerEventID:   triggerEventID,
		})
	}

	return m.repo.CreateWebhook(ctx, dbInput)
}

func (m *webhookManager) GetWebhook(ctx context.Context, webhookID, accountID string) (*webhooks.Webhook, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetWebhook(ctx, webhookID, accountID)
}

func (m *webhookManager) GetWebhooks(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[webhooks.Webhook], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetWebhooks(ctx, accountID, filter)
}

func (m *webhookManager) ArchiveWebhook(ctx context.Context, webhookID, accountID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.ArchiveWebhook(ctx, webhookID, accountID)
}

func (m *webhookManager) AddWebhookTriggerConfig(ctx context.Context, accountID string, input *webhooks.WebhookTriggerConfigCreationRequestInput) (*webhooks.WebhookTriggerConfig, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, observability.PrepareError(errors.New("nil trigger config creation input"), span, "nil trigger config creation input")
	}
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating trigger config creation input")
	}

	dbInput := &webhooks.WebhookTriggerConfigDatabaseCreationInput{
		ID:               identifiers.New(),
		BelongsToWebhook: input.BelongsToWebhook,
		TriggerEventID:   input.TriggerEventID,
	}
	return m.repo.AddWebhookTriggerConfig(ctx, accountID, dbInput)
}

func (m *webhookManager) ArchiveWebhookTriggerConfig(ctx context.Context, webhookID, configID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.ArchiveWebhookTriggerConfig(ctx, webhookID, configID)
}

func (m *webhookManager) CreateWebhookTriggerEvent(ctx context.Context, input *webhooks.WebhookTriggerEventCreationRequestInput) (*webhooks.WebhookTriggerEvent, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, observability.PrepareError(errors.New("nil trigger event creation input"), span, "nil trigger event creation input")
	}
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating trigger event creation input")
	}

	dbInput := &webhooks.WebhookTriggerEventDatabaseCreationInput{
		ID:          identifiers.New(),
		Name:        input.Name,
		Description: input.Description,
	}
	return m.repo.CreateWebhookTriggerEvent(ctx, dbInput)
}

func (m *webhookManager) GetWebhookTriggerEvent(ctx context.Context, id string) (*webhooks.WebhookTriggerEvent, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetWebhookTriggerEvent(ctx, id)
}

func (m *webhookManager) GetWebhookTriggerEvents(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[webhooks.WebhookTriggerEvent], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetWebhookTriggerEvents(ctx, filter)
}

func (m *webhookManager) UpdateWebhookTriggerEvent(ctx context.Context, id string, input *webhooks.WebhookTriggerEventUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.UpdateWebhookTriggerEvent(ctx, id, input)
}

func (m *webhookManager) ArchiveWebhookTriggerEvent(ctx context.Context, id string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.ArchiveWebhookTriggerEvent(ctx, id)
}
