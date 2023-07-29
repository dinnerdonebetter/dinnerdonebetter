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
	rawWebhooksTableCreationColumns = []string{
		"id",
		"name",
		"content_type",
		"url",
		"method",
		belongsToHousehold,
	}

	rawWebhooksTableColumns = []string{
		webhooksColumn(id),
		webhooksColumn("name"),
		webhooksColumn("content_type"),
		webhooksColumn("url"),
		webhooksColumn("method"),
		webhooksColumn(createdAt),
		webhooksColumn(lastUpdatedAt),
		webhooksColumn(archivedAt),
		webhooksColumn(belongsToHousehold),
	}

	rawWebhookTriggerEventsTableColumns = []string{
		webhookTriggerEventsColumns(id),
		webhookTriggerEventsColumns("trigger_event"),
		webhookTriggerEventsColumns("belongs_to_webhook"),
		webhookTriggerEventsColumns(createdAt),
		webhookTriggerEventsColumns(archivedAt),
	}

	fullWebhooksTableColumns = mergeSlicesAtIndex(rawWebhooksTableColumns, rawWebhookTriggerEventsTableColumns, 5)
)

func fleshOutVariablesForColumns(columns []string) []any {
	out := []any{}

	for range columns {
		out = append(out, dummyValue)
	}

	return out
}

func buildCreateWebhookQuery(sqlBuilder squirrel.StatementBuilderType) string {
	query, _, err := sqlBuilder.Insert(webhooksTableName).
		Columns(rawWebhooksTableCreationColumns...).
		Values(fleshOutVariablesForColumns(rawWebhooksTableCreationColumns)...).
		ToSql()

	if err != nil {
		panic(err)
	}

	return query
}

func buildGetOneWebhookQuery(sqlBuilder squirrel.StatementBuilderType) string {
	query, _, err := sqlBuilder.Select(fullWebhooksTableColumns...).
		From(webhooksTableName).
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

func buildArchiveWebhookQuery(sqlBuilder squirrel.StatementBuilderType) string {
	query, _, err := sqlBuilder.Update(webhooksTableName).
		Set(lastUpdatedAt, postgresNow).
		Set(archivedAt, postgresNow).
		Where(squirrel.Eq{
			id:                 dummyValue,
			belongsToHousehold: dummyValue,
			archivedAt:         nil,
		}).
		ToSql()

	if err != nil {
		panic(err)
	}

	return query
}
