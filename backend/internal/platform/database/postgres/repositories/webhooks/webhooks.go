package webhooks

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	types "github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/webhooks/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	resourceTypeWebhooks             = "webhooks"
	resourceTypeWebhookTriggerEvents = "webhook_trigger_events"
)

var (
	_ types.WebhookDataManager = (*Querier)(nil)
)

// WebhookExists fetches whether a webhook exists from the database.
func (q *Querier) WebhookExists(ctx context.Context, webhookID, accountID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if webhookID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	if accountID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	result, err := q.generatedQuerier.CheckWebhookExistence(ctx, q.db, &generated.CheckWebhookExistenceParams{
		BelongsToAccount: accountID,
		ID:               webhookID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
}

// GetWebhook fetches a webhook from the database.
func (q *Querier) GetWebhook(ctx context.Context, webhookID, accountID string) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if webhookID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.WebhookIDKey, webhookID)
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	results, err := q.generatedQuerier.GetWebhook(ctx, q.db, &generated.GetWebhookParams{
		BelongsToAccount: accountID,
		ID:               webhookID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook")
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	webhook := &types.Webhook{
		Events: []*types.WebhookTriggerEvent{},
	}
	for _, result := range results {
		webhook.CreatedAt = result.WebhookCreatedAt
		webhook.ArchivedAt = database.TimePointerFromNullTime(result.WebhookArchivedAt)
		webhook.LastUpdatedAt = database.TimePointerFromNullTime(result.WebhookLastUpdatedAt)
		webhook.Name = result.WebhookName
		webhook.URL = result.WebhookUrl
		webhook.Method = result.WebhookMethod
		webhook.ID = result.WebhookID
		webhook.BelongsToAccount = result.WebhookBelongsToAccount
		webhook.ContentType = result.WebhookContentType

		webhook.Events = append(webhook.Events, &types.WebhookTriggerEvent{
			CreatedAt:        result.WebhookTriggerEventCreatedAt,
			ArchivedAt:       database.TimePointerFromNullTime(result.WebhookTriggerEventArchivedAt),
			ID:               result.WebhookTriggerEventID,
			BelongsToWebhook: result.WebhookTriggerEventBelongsToWebhook,
			TriggerEvent:     string(result.WebhookTriggerEventTriggerEvent),
		})
	}

	return webhook, nil
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (q *Querier) GetWebhooks(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Webhook], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &filtering.QueryFilteredResult[types.Webhook]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetWebhooksForAccount(ctx, q.db, &generated.GetWebhooksForAccountParams{
		BelongsToAccount: accountID,
		CreatedBefore:    database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:     database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:    database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:     database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:      database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:       database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived:  database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhooks from database")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Webhook{
			CreatedAt:        result.CreatedAt,
			ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
			Name:             result.Name,
			URL:              result.URL,
			Method:           result.Method,
			ID:               result.ID,
			BelongsToAccount: result.BelongsToAccount,
			ContentType:      result.ContentType,
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetWebhooksForAccountAndEvent fetches a list of webhooks from the database that meet a particular filter.
func (q *Querier) GetWebhooksForAccountAndEvent(ctx context.Context, accountID, eventType string) ([]*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	databaseResults, err := q.generatedQuerier.GetWebhooksForAccountAndEvent(ctx, q.db, &generated.GetWebhooksForAccountAndEventParams{
		BelongsToAccount: accountID,
		TriggerEvent:     generated.WebhookEvent(eventType),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhooks from database")
	}

	results := []*types.Webhook{}
	for _, result := range databaseResults {
		results = append(results, &types.Webhook{
			CreatedAt:        result.CreatedAt,
			ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
			Name:             result.Name,
			URL:              result.URL,
			Method:           result.Method,
			ID:               result.ID,
			BelongsToAccount: result.BelongsToAccount,
			ContentType:      result.ContentType,
		})
	}

	return results, nil
}

// CreateWebhook creates a webhook in a database.
func (q *Querier) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, input.BelongsToAccount)
	logger = logger.WithValue(keys.AccountIDKey, input.BelongsToAccount)

	logger.Debug("CreateWebhook invoked")

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = q.generatedQuerier.CreateWebhook(ctx, tx, &generated.CreateWebhookParams{
		ID:               input.ID,
		Name:             input.Name,
		ContentType:      input.ContentType,
		URL:              input.URL,
		Method:           input.Method,
		BelongsToAccount: input.BelongsToAccount,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook creation query")
	}

	x := &types.Webhook{
		ID:               input.ID,
		Name:             input.Name,
		ContentType:      input.ContentType,
		URL:              input.URL,
		Method:           input.Method,
		BelongsToAccount: input.BelongsToAccount,
		CreatedAt:        q.CurrentTime(),
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &x.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeWebhooks,
		RelevantID:       x.ID,
		EventType:        audit.AuditLogEventTypeCreated,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	for i := range input.Events {
		evt := input.Events[i]
		evt.BelongsToWebhook = input.ID

		e, webhookTriggerEventCreationErr := q.createWebhookTriggerEvent(ctx, tx, x.BelongsToAccount, evt)
		if webhookTriggerEventCreationErr != nil {
			q.RollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(webhookTriggerEventCreationErr, logger, span, "performing webhook creation query")
		}

		x.Events = append(x.Events, e)
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	tracing.AttachToSpan(span, keys.WebhookIDKey, x.ID)

	return x, nil
}

// createWebhookTriggerEvent creates a webhook trigger event in a database.
func (q *Querier) createWebhookTriggerEvent(ctx context.Context, querier database.SQLQueryExecutor, accountID string, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.WebhookIDKey, input.BelongsToWebhook)

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)
	logger = logger.WithValue(keys.AccountIDKey, accountID)

	if err := q.generatedQuerier.CreateWebhookTriggerEvent(ctx, querier, &generated.CreateWebhookTriggerEventParams{
		ID:               input.ID,
		TriggerEvent:     generated.WebhookEvent(input.TriggerEvent),
		BelongsToWebhook: input.BelongsToWebhook,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook trigger event creation query")
	}

	if _, err := q.auditLogEntryRepo.CreateAuditLogEntry(ctx, querier, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &accountID,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeWebhookTriggerEvents,
		RelevantID:       input.ID,
		EventType:        audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	x := &types.WebhookTriggerEvent{
		ID:               input.ID,
		TriggerEvent:     input.TriggerEvent,
		BelongsToWebhook: input.BelongsToWebhook,
		CreatedAt:        q.CurrentTime(),
		ArchivedAt:       nil,
	}

	return x, nil
}

// ArchiveWebhook archives a webhook from the database.
func (q *Querier) ArchiveWebhook(ctx context.Context, webhookID, accountID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	if accountID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	logger := q.logger.WithValues(map[string]any{
		keys.WebhookIDKey: webhookID,
		keys.AccountIDKey: accountID,
	})

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.ArchiveWebhook(ctx, tx, &generated.ArchiveWebhookParams{
		BelongsToAccount: accountID,
		ID:               webhookID,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &accountID,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeWebhooks,
		RelevantID:       webhookID,
		EventType:        audit.AuditLogEventTypeArchived,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	logger.Info("webhook archived")

	return nil
}

// AddWebhookTriggerEvent adds a webhook trigger event from the database.
func (q *Querier) AddWebhookTriggerEvent(ctx context.Context, accountID string, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.WebhookIDKey, input.BelongsToWebhook)

	logger := q.logger.WithValues(map[string]any{
		keys.WebhookIDKey:             input.BelongsToWebhook,
		keys.WebhookTriggerEventIDKey: input.ID,
		keys.AccountIDKey:             accountID,
	})

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	e, err := q.createWebhookTriggerEvent(ctx, tx, accountID, input)
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook creation query")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	logger.Info("webhook trigger event archived")

	return e, nil
}

// ArchiveWebhookTriggerEvent archives a webhook trigger event from the database.
func (q *Querier) ArchiveWebhookTriggerEvent(ctx context.Context, webhookID, webhookTriggerEventID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.WebhookIDKey, webhookID)

	if webhookTriggerEventID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.WebhookTriggerEventIDKey, webhookTriggerEventID)

	logger := q.logger.WithValues(map[string]any{
		keys.WebhookIDKey:             webhookID,
		keys.WebhookTriggerEventIDKey: webhookTriggerEventID,
	})

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.ArchiveWebhookTriggerEvent(ctx, tx, &generated.ArchiveWebhookTriggerEventParams{
		BelongsToWebhook: webhookID,
		ID:               webhookTriggerEventID,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook trigger event")
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeWebhookTriggerEvents,
		RelevantID:   webhookID,
		EventType:    audit.AuditLogEventTypeArchived,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing database transaction")
	}

	logger.Info("webhook trigger event archived")

	return nil
}
