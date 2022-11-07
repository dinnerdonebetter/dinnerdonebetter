package main

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	webhooksTableName             = "webhooks"
	webhookTriggerEventsTableName = "webhook_trigger_events"
)

var (
	webhookColumns = []string{
		"id",
		"name",
		"content_type",
		"url",
		"method",
		"created_at",
		"last_updated_at",
		"archived_at",
		"belongs_to_household",
	}
	webhookTriggerEventsColumns = []string{
		"id",
		"trigger_event",
		"belongs_to_webhook",
		"created_at",
		"archived_at",
	}
)

func buildWebhooksGetOne(ctx context.Context, psql squirrel.StatementBuilderType, format bool, filePath string) error {
	columns := mergeStringSlicesAtIndex(
		joinTableName(webhooksTableName, webhookColumns),
		joinTableName(webhookTriggerEventsTableName, webhookTriggerEventsColumns),
		5,
	)

	q := psql.Select(columns...).
		From(webhookTriggerEventsTableName).
		Join("webhooks ON webhook_trigger_events.belongs_to_webhook=webhooks.id").
		Where(squirrel.Eq{
			"webhooks.archived_at":               nil,
			"webhook_trigger_events.archived_at": nil,
			"webhooks.belongs_to_household":      whatever,
			"webhooks.id":                        whatever,
		})

	query, _, err := q.ToSql()
	if err != nil {
		return err
	}

	if format {
		query = formatQuery(query)
	}

	content := fmt.Sprintf("-- name: GetWebhook :one\n%s", query)

	return saveFile(ctx, content, filePath)
}
