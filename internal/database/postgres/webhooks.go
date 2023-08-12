package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
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

	result, err := q.generatedQuerier.CheckWebhookExistence(ctx, q.db, &generated.CheckWebhookExistenceParams{
		BelongsToHousehold: householdID,
		ID:                 webhookID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
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
	tracing.AttachWebhookIDToSpan(span, webhookID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	results, err := q.generatedQuerier.GetWebhook(ctx, q.db, &generated.GetWebhookParams{
		BelongsToHousehold: householdID,
		ID:                 webhookID,
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
		webhook.CreatedAt = result.CreatedAt
		webhook.ArchivedAt = timePointerFromNullTime(result.ArchivedAt)
		webhook.LastUpdatedAt = timePointerFromNullTime(result.LastUpdatedAt)
		webhook.Name = result.Name
		webhook.URL = result.Url
		webhook.Method = result.Method
		webhook.ID = result.ID
		webhook.BelongsToHousehold = result.BelongsToHousehold
		webhook.ContentType = result.ContentType

		webhook.Events = append(webhook.Events, &types.WebhookTriggerEvent{
			CreatedAt:        result.CreatedAt_2,
			ArchivedAt:       timePointerFromNullTime(result.ArchivedAt_2),
			ID:               result.ID_2,
			BelongsToWebhook: result.BelongsToWebhook,
			TriggerEvent:     string(result.TriggerEvent),
		})
	}

	return webhook, nil
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (q *Querier) GetWebhooks(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Webhook], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.Webhook]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetWebhooks(ctx, q.db, &generated.GetWebhooksParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		HouseholdID:   householdID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhooks from database")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Webhook{
			CreatedAt:          result.CreatedAt,
			ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
			Name:               result.Name,
			URL:                result.Url,
			Method:             result.Method,
			ID:                 result.ID,
			BelongsToHousehold: result.BelongsToHousehold,
			ContentType:        result.ContentType,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
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

	logger.Debug("CreateWebhook invoked")

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

// createWebhookTriggerEvent creates a webhook trigger event in a database.
func (q *Querier) createWebhookTriggerEvent(ctx context.Context, querier database.SQLQueryExecutor, input *types.WebhookTriggerEventDatabaseCreationInput) (*types.WebhookTriggerEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachWebhookIDToSpan(span, input.BelongsToWebhook)

	if err := q.generatedQuerier.CreateWebhookTriggerEvent(ctx, querier, &generated.CreateWebhookTriggerEventParams{
		ID:               input.ID,
		TriggerEvent:     generated.WebhookEvent(input.TriggerEvent),
		BelongsToWebhook: input.BelongsToWebhook,
	}); err != nil {
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

// ArchiveWebhook archives a webhook from the database.
func (q *Querier) ArchiveWebhook(ctx context.Context, webhookID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachWebhookIDToSpan(span, webhookID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValues(map[string]any{
		keys.WebhookIDKey:   webhookID,
		keys.HouseholdIDKey: householdID,
	})

	if err := q.generatedQuerier.ArchiveWebhook(ctx, q.db, &generated.ArchiveWebhookParams{
		BelongsToHousehold: householdID,
		ID:                 webhookID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving webhook")
	}

	logger.Info("webhook archived")

	return nil
}
