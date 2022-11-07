package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.WebhookDataManager = (*Querier)(nil)
)

// WebhookExists fetches whether a webhook exists from the database.
func (q *Querier) WebhookExists(ctx context.Context, webhookID, householdID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if webhookID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	result, err := q.generatedQuerier.WebhookExists(ctx, q.db, &generated.WebhookExistsParams{
		BelongsToHousehold: householdID,
		ID:                 webhookID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, observability.PrepareAndLogError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
}

// convertWebhookRow is a consistent way to turn a *sql.Row into a webhook struct.
func (q *Querier) convertWebhookRow(ctx context.Context, rows []*generated.GetWebhookRow) (webhook *types.Webhook, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for _, row := range rows {
		if webhook == nil {
			webhook = &types.Webhook{
				CreatedAt:          row.CreatedAt,
				Name:               row.Name,
				URL:                row.Url,
				Method:             row.Method,
				ID:                 row.ID,
				BelongsToHousehold: row.BelongsToHousehold,
				ContentType:        row.ContentType,
				Events:             nil,
			}

			if row.LastUpdatedAt.Valid {
				webhook.LastUpdatedAt = &row.LastUpdatedAt.Time
			}

			if row.ArchivedAt.Valid {
				webhook.ArchivedAt = &row.ArchivedAt.Time
			}
		}

		webhookTriggerEvent := &types.WebhookTriggerEvent{
			CreatedAt:        row.CreatedAt_2,
			ArchivedAt:       nil,
			ID:               row.ID_2,
			BelongsToWebhook: row.BelongsToWebhook,
			TriggerEvent:     string(row.TriggerEvent),
		}

		if row.ArchivedAt_2.Valid {
			webhookTriggerEvent.ArchivedAt = &row.ArchivedAt_2.Time
		}

		webhook.Events = append(webhook.Events, webhookTriggerEvent)
	}

	if webhook.ID == "" {
		return nil, sql.ErrNoRows
	}

	return webhook, nil
}

// GetWebhook fetches a webhook from the database.
func (q *Querier) GetWebhook(ctx context.Context, webhookID, householdID string) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if webhookID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	row, err := q.generatedQuerier.GetWebhook(ctx, q.db, &generated.GetWebhookParams{
		BelongsToHousehold: householdID,
		ID:                 webhookID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook")
	}

	webhook, err := q.convertWebhookRow(ctx, row)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning webhook")
	}

	return webhook, nil
}

// convertWebhookRows provides a consistent way to turn sql rows into a slice of webhooks.
func (q *Querier) convertWebhookRows(ctx context.Context, rows []*generated.GetWebhooksRow) (webhooks []*types.Webhook, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x := &types.Webhook{}
	for _, row := range rows {
		webhook := &types.Webhook{
			CreatedAt:          row.CreatedAt,
			Name:               row.Name,
			URL:                row.Url,
			Method:             row.Method,
			ID:                 row.ID,
			BelongsToHousehold: row.BelongsToHousehold,
			ContentType:        row.ContentType,
		}

		if row.LastUpdatedAt.Valid {
			webhook.LastUpdatedAt = &row.LastUpdatedAt.Time
		}

		if row.ArchivedAt.Valid {
			webhook.ArchivedAt = &row.ArchivedAt.Time
		}

		webhookTriggerEvent := &types.WebhookTriggerEvent{
			CreatedAt:        row.CreatedAt_2,
			ID:               row.ID_2,
			BelongsToWebhook: row.BelongsToWebhook,
			TriggerEvent:     string(row.TriggerEvent),
		}

		if row.ArchivedAt_2.Valid {
			webhookTriggerEvent.ArchivedAt = &row.ArchivedAt_2.Time
		}

		if x.ID == "" {
			events := x.Events
			x = webhook
			x.Events = events
		}

		if x.ID != webhook.ID {
			webhooks = append(webhooks, x)
			x = webhook
		}

		x.Events = append(x.Events, webhookTriggerEvent)
	}

	if x.ID != "" {
		webhooks = append(webhooks, x)
	}

	if len(webhooks) == 0 {
		return nil, sql.ErrNoRows
	}

	return webhooks, nil
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (q *Querier) GetWebhooks(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.WebhookList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachQueryFilterToSpan(span, filter)

	filterArgs := filter.ToDatabaseArgs()
	x := &types.WebhookList{}

	rows, err := q.generatedQuerier.GetWebhooks(ctx, q.db, &generated.GetWebhooksParams{
		BelongsToHousehold: householdID,
		CreatedAfter:       filterArgs.CreatedAfter,
		CreatedBefore:      filterArgs.CreatedBefore,
		UpdatedAfter:       filterArgs.UpdatedAfter,
		UpdatedBefore:      filterArgs.UpdatedBefore,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
	}

	x.Webhooks, err = q.convertWebhookRows(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "converting response from database")
	}

	return x, nil
}

// CreateWebhook creates a webhook in a database.
func (q *Querier) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)
	logger := q.logger.WithValue(keys.HouseholdIDKey, input.BelongsToHousehold)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = q.generatedQuerier.CreateWebhook(ctx, tx, &generated.CreateWebhookParams{
		ID:                 input.ID,
		Name:               input.Name,
		ContentType:        input.ContentType,
		Url:                input.URL,
		Method:             input.Method,
		BelongsToHousehold: input.BelongsToHousehold,
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating webhook trigger event")
	}

	x := &types.Webhook{
		ID:                 input.ID,
		Name:               input.Name,
		ContentType:        input.ContentType,
		URL:                input.URL,
		Method:             input.Method,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	for i := range input.Events {
		evt := input.Events[i]
		evt.BelongsToWebhook = input.ID

		err = q.generatedQuerier.CreateWebhookTriggerEvent(ctx, tx, &generated.CreateWebhookTriggerEventParams{
			ID:               evt.ID,
			TriggerEvent:     generated.WebhookEvent(evt.TriggerEvent),
			BelongsToWebhook: evt.BelongsToWebhook,
		})
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "creating webhook trigger event")
		}

		x.Events = append(x.Events, &types.WebhookTriggerEvent{
			ID:               evt.ID,
			TriggerEvent:     evt.TriggerEvent,
			BelongsToWebhook: evt.BelongsToWebhook,
			CreatedAt:        q.currentTime(),
		})
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	tracing.AttachWebhookIDToSpan(span, x.ID)

	logger.Debug("webhook created")

	return x, nil
}

// ArchiveWebhook archives a webhook from the database.
func (q *Querier) ArchiveWebhook(ctx context.Context, webhookID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" || householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachWebhookIDToSpan(span, webhookID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.WebhookIDKey:   webhookID,
		keys.HouseholdIDKey: householdID,
	})

	if err := q.generatedQuerier.ArchiveWebhook(ctx, q.db, &generated.ArchiveWebhookParams{
		ID:                 webhookID,
		BelongsToHousehold: householdID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}

	logger.Info("webhook archived")

	return nil
}
