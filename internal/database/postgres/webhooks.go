package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	_ types.WebhookDataManager = (*Querier)(nil)
)

// scanWebhookWithEvents is a consistent way to turn a *sql.Row into a webhook struct.
func (q *Querier) scanWebhook(ctx context.Context, rows database.ResultIterator) (webhook *types.Webhook, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	webhook = &types.Webhook{}
	var (
		lastUpdatedAt,
		archivedAt sql.NullTime
	)

	for rows.Next() {
		webhookTriggerEvent := &types.WebhookTriggerEvent{}

		targetVars := []interface{}{
			&webhook.ID,
			&webhook.Name,
			&webhook.ContentType,
			&webhook.URL,
			&webhook.Method,
			&webhookTriggerEvent.ID,
			&webhookTriggerEvent.TriggerEvent,
			&webhookTriggerEvent.BelongsToWebhook,
			&webhookTriggerEvent.CreatedAt,
			&webhookTriggerEvent.ArchivedAt,
			&webhook.CreatedAt,
			&lastUpdatedAt,
			&archivedAt,
			&webhook.BelongsToHousehold,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning webhook")
		}

		webhook.Events = append(webhook.Events, webhookTriggerEvent)
	}

	if err = rows.Err(); err != nil {
		return nil, observability.PrepareError(err, span, "fetching webhook from database")
	}

	if webhook.ID == "" {
		return nil, sql.ErrNoRows
	}

	if lastUpdatedAt.Valid {
		webhook.LastUpdatedAt = &lastUpdatedAt.Time
	}

	if archivedAt.Valid {
		webhook.ArchivedAt = &archivedAt.Time
	}

	return webhook, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (q *Querier) scanWebhooks(ctx context.Context, rows database.ResultIterator) (webhooks []*types.Webhook, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x := &types.Webhook{}
	for rows.Next() {
		webhook := &types.Webhook{}
		webhookTriggerEvent := &types.WebhookTriggerEvent{}

		var (
			lastUpdatedAt,
			archivedAt sql.NullTime
		)

		targetVars := []interface{}{
			&webhook.ID,
			&webhook.Name,
			&webhook.ContentType,
			&webhook.URL,
			&webhook.Method,
			&webhookTriggerEvent.ID,
			&webhookTriggerEvent.TriggerEvent,
			&webhookTriggerEvent.BelongsToWebhook,
			&webhookTriggerEvent.CreatedAt,
			&webhookTriggerEvent.ArchivedAt,
			&webhook.CreatedAt,
			&lastUpdatedAt,
			&archivedAt,
			&webhook.BelongsToHousehold,
		}

		if err = rows.Scan(targetVars...); err != nil {
			return nil, observability.PrepareError(err, span, "scanning webhook")
		}

		if lastUpdatedAt.Valid {
			webhook.LastUpdatedAt = &lastUpdatedAt.Time
		}
		if archivedAt.Valid {
			webhook.ArchivedAt = &archivedAt.Time
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

	if err = rows.Err(); err != nil {
		return nil, observability.PrepareError(err, span, "fetching webhook from database")
	}

	if len(webhooks) == 0 {
		return nil, sql.ErrNoRows
	}

	return webhooks, nil
}

//go:embed queries/webhooks/exists.sql
var webhookExistenceQuery string

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

	args := []interface{}{
		householdID,
		webhookID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, webhookExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
}

//go:embed queries/webhooks/get_one.sql
var getWebhookQuery string

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

	args := []interface{}{
		householdID,
		webhookID,
	}

	rows, err := q.getRows(ctx, q.db, "webhook", getWebhookQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook")
	}

	webhook, err := q.scanWebhook(ctx, rows)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning webhook")
	}

	return webhook, nil
}

//go:embed queries/webhooks/get_many.sql
var getWebhooksForHouseholdQuery string

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (q *Querier) GetWebhooks(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.WebhookList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &types.WebhookList{}
	args := []interface{}{
		householdID,
	}

	rows, err := q.getRows(ctx, q.db, "webhooks", getWebhooksForHouseholdQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
	}

	if x.Webhooks, err = q.scanWebhooks(ctx, rows); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning database response")
	}

	return x, nil
}

//go:embed queries/webhooks/create.sql
var createWebhookQuery string

// CreateWebhook creates a webhook in a database.
func (q *Querier) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)
	logger := q.logger.WithValue(keys.HouseholdIDKey, input.BelongsToHousehold)

	logger.Debug("CreateWebhook invoked")

	args := []interface{}{
		input.ID,
		input.Name,
		input.ContentType,
		input.URL,
		input.Method,
		input.BelongsToHousehold,
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = q.performWriteQuery(ctx, tx, "webhook creation", createWebhookQuery, args); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook creation query")
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

		e, webhookTriggerEventCreationErr := q.createWebhookTriggerEvent(ctx, tx, evt)
		if webhookTriggerEventCreationErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(webhookTriggerEventCreationErr, logger, span, "performing webhook creation query")
		}

		x.Events = append(x.Events, e)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	tracing.AttachWebhookIDToSpan(span, x.ID)

	return x, nil
}

//go:embed queries/webhook_trigger_events/create.sql
var createWebhookTriggerEventQuery string

// createWebhookTriggerEvent creates a webhook trigger event in a database.
func (q *Querier) createWebhookTriggerEvent(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachWebhookIDToSpan(span, input.BelongsToWebhook)

	createWebhookTriggerEventArgs := []interface{}{
		input.ID,
		input.TriggerEvent,
		input.BelongsToWebhook,
	}

	if err := q.performWriteQuery(ctx, querier, "webhook trigger event creation", createWebhookTriggerEventQuery, createWebhookTriggerEventArgs); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook trigger event creation query")
	}

	x := &types.WebhookTriggerEvent{
		ID:               input.ID,
		TriggerEvent:     input.TriggerEvent,
		BelongsToWebhook: input.BelongsToWebhook,
		CreatedAt:        q.currentTime(),
	}

	return x, nil
}

//go:embed queries/webhooks/archive.sql
var archiveWebhookQuery string

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

	args := []interface{}{householdID, webhookID}

	if err := q.performWriteQuery(ctx, q.db, "webhook archive", archiveWebhookQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}

	logger.Info("webhook archived")

	return nil
}
