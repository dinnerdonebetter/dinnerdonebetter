package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	webhooksTableName = "webhooks"
)

var (
	webhooksColumns = []string{
		"id",
		nameColumn,
		"content_type",
		"url",
		"method",
		createdAtColumn,
		lastUpdatedAtColumn,
		archivedAtColumn,
		belongsToHouseholdColumn,
	}
)

func buildWebhooksQueries() []*Query {
	insertColumns := filterForInsert(webhooksColumns)
	fullSelectColumns := mergeColumns(
		applyToEach(webhooksColumns, func(_ int, s string) string {
			return fullColumnName(webhooksTableName, s)
		}),
		applyToEach(webhookTriggerEventsColumns, func(_ int, s string) string {
			return fullColumnName(webhookTriggerEventsTableName, s)
		}),
		5,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveWebhook",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s
   SET %s = %s
 WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(household_id);`,
				webhooksTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
				belongsToHouseholdColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateWebhook",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				webhooksTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckWebhookExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.id
	FROM %s
	WHERE %s.%s IS NULL
	AND %s.id = sqlc.arg(id)
	AND %s.belongs_to_household = sqlc.arg(household_id)
);`,
				webhooksTableName,
				webhooksTableName,
				webhooksTableName,
				archivedAtColumn,
				webhooksTableName,
				webhooksTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhooksForHousehold",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM webhooks
	JOIN %s
WHERE webhooks.archived_at IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					webhooksTableName,
					true,
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
				),
				buildTotalCountSelect(
					webhooksTableName,
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
					"webhook_trigger_events.archived_at IS NULL",
				),
				webhookTriggerEventsJoin,
				buildFilterConditions(
					webhooksTableName,
					true,
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
					"webhook_trigger_events.archived_at IS NULL",
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhooksForHouseholdAndEvent",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT 
	%s
FROM webhooks
	JOIN %s
WHERE webhook_trigger_events.archived_at IS NULL
	AND webhook_trigger_events.trigger_event = sqlc.arg(trigger_event)
	AND webhooks.belongs_to_household = sqlc.arg(household_id)
	AND webhooks.archived_at IS NULL;`,
				strings.Join(applyToEach(webhooksColumns, func(_ int, s string) string {
					return fullColumnName(webhooksTableName, s)
				}), ",\n\t"),
				webhookTriggerEventsJoin,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhook",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT 
	%s
FROM webhooks
	JOIN %s
WHERE webhook_trigger_events.archived_at IS NULL
	AND webhooks.belongs_to_household = sqlc.arg(household_id)
	AND webhooks.archived_at IS NULL
	AND webhooks.id = sqlc.arg(id);`,
				strings.Join(applyToEach(fullSelectColumns, func(_ int, s string) string {
					parts := strings.Split(s, ".")
					return fmt.Sprintf("%s as %s_%s",
						s, strings.TrimSuffix(parts[0], "s"), parts[1],
					)
				}), ",\n\t"),
				webhookTriggerEventsJoin,
			)),
		},
	}
}
