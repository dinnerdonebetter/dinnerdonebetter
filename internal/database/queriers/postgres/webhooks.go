package postgres

import (
	"context"
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
	_ types.WebhookDataManager = (*SQLQuerier)(nil)

	// householdInvitationsTableColumns are the columns for the webhooks table.
	webhooksTableColumns = []string{
		"webhooks.id",
		"webhooks.name",
		"webhooks.content_type",
		"webhooks.url",
		"webhooks.method",
		"webhooks.events",
		"webhooks.data_types",
		"webhooks.topics",
		"webhooks.created_on",
		"webhooks.last_updated_on",
		"webhooks.archived_on",
		"webhooks.belongs_to_household",
	}
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (q *SQLQuerier) scanWebhook(ctx context.Context, scan database.Scanner, includeCounts bool) (webhook *types.Webhook, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)
	webhook = &types.Webhook{}

	var (
		eventsStr,
		dataTypesStr,
		topicsStr string
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
		&webhook.CreatedOn,
		&webhook.LastUpdatedOn,
		&webhook.ArchivedOn,
		&webhook.BelongsToHousehold,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "scanning webhook")
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

	return webhook, filteredCount, totalCount, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (q *SQLQuerier) scanWebhooks(ctx context.Context, rows database.ResultIterator, includeCounts bool) (webhooks []*types.Webhook, filteredCount, totalCount uint64, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

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
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	if err = rows.Close(); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	return webhooks, filteredCount, totalCount, nil
}

const webhookExistenceQuery = "SELECT EXISTS ( SELECT webhooks.id FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_household = $1 AND webhooks.id = $2 )"

// WebhookExists fetches whether a webhook exists from the database.
func (q *SQLQuerier) WebhookExists(ctx context.Context, webhookID, householdID string) (exists bool, err error) {
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
		return false, observability.PrepareError(err, logger, span, "performing webhook existence check")
	}

	return result, nil
}

const getWebhookQuery = `
	SELECT webhooks.id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_household FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_household = $1 AND webhooks.id = $2
`

// GetWebhook fetches a webhook from the database.
func (q *SQLQuerier) GetWebhook(ctx context.Context, webhookID, householdID string) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if webhookID == "" || householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.WebhookIDKey:   webhookID,
		keys.HouseholdIDKey: householdID,
	})

	args := []interface{}{
		householdID,
		webhookID,
	}

	row := q.getOneRow(ctx, q.db, "webhook", getWebhookQuery, args)

	webhook, _, _, err := q.scanWebhook(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning webhook")
	}

	return webhook, nil
}

const getAllWebhooksCountQuery = `
	SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL
`

// GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter.
func (q *SQLQuerier) GetAllWebhooksCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getAllWebhooksCountQuery, "fetching count of webhooks")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of webhooks")
	}

	return count, nil
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (q *SQLQuerier) GetWebhooks(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.WebhookList, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := q.logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &types.WebhookList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "webhooks", nil, nil, nil, "belongs_to_household", webhooksTableColumns, householdID, false, filter)

	rows, err := q.performReadQuery(ctx, q.db, "webhooks", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching webhook from database")
	}

	if x.Webhooks, x.FilteredCount, x.TotalCount, err = q.scanWebhooks(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning database response")
	}

	return x, nil
}

const createWebhookQuery = `
	INSERT INTO webhooks (id,name,content_type,url,method,events,data_types,topics,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`

// CreateWebhook creates a webhook in a database.
func (q *SQLQuerier) CreateWebhook(ctx context.Context, input *types.WebhookDatabaseCreationInput) (*types.Webhook, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachHouseholdIDToSpan(span, input.BelongsToHousehold)
	logger := q.logger.WithValue(keys.HouseholdIDKey, input.BelongsToHousehold)

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
		return nil, observability.PrepareError(err, logger, span, "performing webhook creation query")
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
		CreatedOn:          q.currentTime(),
	}

	tracing.AttachWebhookIDToSpan(span, x.ID)
	logger = logger.WithValue(keys.WebhookIDKey, x.ID)

	return x, nil
}

const archiveWebhookQuery = `
UPDATE webhooks SET
	last_updated_on = extract(epoch FROM NOW()), 
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL 
AND belongs_to_household = $1
AND id = $2
`

// ArchiveWebhook archives a webhook from the database.
func (q *SQLQuerier) ArchiveWebhook(ctx context.Context, webhookID, householdID string) error {
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
		return observability.PrepareError(err, logger, span, "archiving webhook")
	}

	logger.Info("webhook archived")

	return nil
}
