package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	webhooksTableName             = "webhooks"
	webhookTriggerEventsTableName = "webhook_trigger_events"
	webhookTriggerEventsJoin      = "webhook_trigger_events ON webhooks.id = webhook_trigger_events.belongs_to_webhook"
)

var (
	webhookTriggerEventsColumns = []string{
		"id",
		"trigger_event",
		"belongs_to_webhook",
		createdAtColumn,
		archivedAtColumn,
	}

	webhooksColumns = []string{
		"id",
		"name",
		"content_type",
		"url",
		"method",
		belongsToHouseholdColumn,
		createdAtColumn,
		lastUpdatedAtColumn,
		archivedAtColumn,
	}
)

func buildWebhooksQueries() []Query {
	insertColumns := filterForInsert(webhooksColumns)

	queryFilteredColumns := append(mergeColumns(
		applyToEach(webhooksColumns, func(s string) string {
			return fullColumnName(webhooksTableName, s)
		}),
		applyToEach(webhookTriggerEventsColumns, func(s string) string {
			return fullColumnName(webhookTriggerEventsTableName, s)
		}), 5),
		buildFilteredColumnCountQuery(webhooksTableName, true, "webhooks.belongs_to_household = sqlc.arg(household_id)"),
		buildTotalColumnCountQuery(webhooksTableName, ""),
	)

	return []Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveWebhook",
				Type: ExecRowsType,
			},
			Content: buildQuery(
				queryBuilder.Update(webhooksTableName).
					Set(lastUpdatedAtColumn, now).
					Set(archivedAtColumn, now).
					Where(squirrel.Eq{
						archivedAtColumn:         nil,
						belongsToHouseholdColumn: placeholder,
						idColumn:                 placeholder,
					}),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateWebhook",
				Type: ExecType,
			},
			Content: buildQuery(
				queryBuilder.Insert(webhooksTableName).
					Columns(insertColumns...).
					Values(argsForList(insertColumns)...),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckWebhookExistence",
				Type: OneType,
			},
			Content: buildExistenceCheckQuery(webhooksTableName, " AND webhooks.belongs_to_household = sqlc.arg(household_id) AND webhooks.id = sqlc.arg(webhook_id)"),
			//Content: buildQuery(
			//	queryBuilder.Select(fullColumnName(webhooksTableName, idColumn)).
			//		Prefix("SELECT EXISTS (").
			//		From(webhooksTableName).
			//		Where(squirrel.Eq{
			//			fullColumnName(webhooksTableName, archivedAtColumn):         nil,
			//			fullColumnName(webhooksTableName, belongsToHouseholdColumn): placeholder,
			//			fullColumnName(webhooksTableName, idColumn):                 placeholder,
			//		}).
			//		Suffix(")"),
			//),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhooksForHousehold",
				Type: ManyType,
			},
			Content: buildQuery(
				queryBuilder.Select(queryFilteredColumns...).
					From(webhooksTableName).
					Join(webhookTriggerEventsJoin).
					Where(squirrel.Eq{
						fullColumnName(webhooksTableName, archivedAtColumn):         nil,
						fullColumnName(webhooksTableName, belongsToHouseholdColumn): placeholder,
					}),
			),
		},
	}
}
