package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	webhooksTableName             = "webhooks"
	webhookTriggerEventsTableName = "webhook_trigger_events"
)

func webhooksColumn(col string) string {
	return buildColumnName(webhooksTableName, col)
}

func webhookTriggerEventsColumns(col string) string {
	return buildColumnName(webhookTriggerEventsTableName, col)
}

var (
	rawWebhooksTableColumns = []string{
		webhooksColumn("id"),
		webhooksColumn("name"),
		webhooksColumn("content_type"),
		webhooksColumn("url"),
		webhooksColumn("method"),
		webhooksColumn("created_at"),
		webhooksColumn("last_updated_at"),
		webhooksColumn("archived_at"),
		webhooksColumn("belongs_to_household"),
	}

	rawWebhookTriggerEventsTableColumns = []string{
		webhookTriggerEventsColumns("id"),
		webhookTriggerEventsColumns("trigger_event"),
		webhookTriggerEventsColumns("belongs_to_webhook"),
		webhookTriggerEventsColumns("created_at"),
		webhookTriggerEventsColumns("archived_at"),
	}

	fullWebhooksTableColumns = mergeSlicesAtIndex(rawWebhooksTableColumns, rawWebhookTriggerEventsTableColumns, 5)
)

func buildGetOneQuery(sqlBuilder squirrel.StatementBuilderType) string {
	query, _, err := sqlBuilder.Select(fullWebhooksTableColumns...).
		From("webhooks").
		Where(squirrel.Eq{
			webhooksColumn(id):                      dummyValue,
			webhooksColumn(belongsToHousehold):      dummyValue,
			webhooksColumn(archivedAt):              nil,
			webhookTriggerEventsColumns(archivedAt): nil,
		}).
		Join(buildJoinStatement(
			webhooksTableName,
			webhookTriggerEventsColumns("belongs_to_webhook"),
			webhooksColumn(id)),
		).
		ToSql()

	if err != nil {
		panic(err)
	}

	return query
}
