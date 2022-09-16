package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"strings"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// webhooksTableEventsSeparator is what the webhooks table uses to separate event subscriptions.
	webhooksTableEventsSeparator = commaSeparator
	// webhooksTableDataTypesSeparator is what the webhooks table uses to separate data type subscriptions.
	webhooksTableDataTypesSeparator = commaSeparator
	// webhooksTableTopicsSeparator is what the webhooks table uses to separate topic subscriptions.
	webhooksTableTopicsSeparator = commaSeparator
)

var (
	_ types.WebhookDataManager = (*Querier)(nil)

	// webhooksTableColumns are the columns for the webhooks table.
	webhooksTableColumns = []string{
		"webhooks.id",
		"webhooks.name",
		"webhooks.content_type",
		"webhooks.url",
		"webhooks.method",
		"webhooks.events",
		"webhooks.data_types",
		"webhooks.topics",
		"webhooks.created_at",
		"webhooks.last_updated_at",
		"webhooks.archived_at",
		"webhooks.belongs_to_household",
	}
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (q *Querier) scanWebhook(ctx context.Context, scan database.Scanner, includeCounts bool) (webhook *types.Webhook, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	webhook = &types.Webhook{}
	var (
		eventsStr,
		dataTypesStr,
		topicsStr string

		lastUpdatedAt,
		archivedAt sql.NullTime
	)

	targetVars := []interface{}{
		&webhook.ID,
		&webhook.Name,
		&webhook.ContentType,
		&webhook.URL,
		&webhook.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&webhook.CreatedAt,
		&lastUpdatedAt,
		&archivedAt,
		&webhook.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "scanning webhook")
	}

	if events := strings.Split(eventsStr, webhooksTableEventsSeparator); len(events) >= 1 && events[0] != "" {
		webhook.Events = events
	}

	if dataTypes := strings.Split(dataTypesStr, webhooksTableDataTypesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		webhook.DataTypes = dataTypes
	}

	if topics := strings.Split(topicsStr, webhooksTableTopicsSeparator); len(topics) >= 1 && topics[0] != "" {
		webhook.Topics = topics
	}

	if lastUpdatedAt.Valid {
		webhook.LastUpdatedAt = &lastUpdatedAt.Time
	}

	if archivedAt.Valid {
		webhook.ArchivedAt = &archivedAt.Time
	}

	return webhook, filteredCount, totalCount, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (q *Querier) scanWebhooks(ctx context.Context, rows database.ResultIterator, includeCounts bool) (webhooks []*types.Webhook, filteredCount, totalCount uint64, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		webhook, fc, tc, scanErr := q.scanWebhook(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		webhooks = append(webhooks, webhook)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "fetching webhook from database")
	}

	if err = rows.Close(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "fetching webhook from database")
	}

	return webhooks, filteredCount, totalCount, nil
}

//go:embed queries/webhooks/webhooks_get_exists.sql
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

//go:embed queries/webhooks/webhooks_get_one.sql
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

	row := q.getOneRow(ctx, q.db, "webhook", getWebhookQuery, args)

	webhook, _, _, err := q.scanWebhook(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning webhook")
	}

	return webhook, nil
}

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
	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, householdID, false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "webhooks", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching webhook from database")
	}

	if x.Webhooks, x.FilteredCount, x.TotalCount, err = q.scanWebhooks(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning database response")
	}

	return x, nil
}

//go:embed queries/webhooks/webhooks_create.sql
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
		strings.Join(input.Events, webhooksTableEventsSeparator),
		strings.Join(input.DataTypes, webhooksTableDataTypesSeparator),
		strings.Join(input.Topics, webhooksTableTopicsSeparator),
		input.BelongsToHousehold,
	}

	if err := q.performWriteQuery(ctx, q.db, "webhook creation", createWebhookQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook creation query")
	}

	x := &types.Webhook{
		ID:                 input.ID,
		Name:               input.Name,
		ContentType:        input.ContentType,
		URL:                input.URL,
		Method:             input.Method,
		Events:             input.Events,
		DataTypes:          input.DataTypes,
		Topics:             input.Topics,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	tracing.AttachWebhookIDToSpan(span, x.ID)

	return x, nil
}

//go:embed queries/webhooks/webhooks_archive.sql
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
