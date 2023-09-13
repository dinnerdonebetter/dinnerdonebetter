package main

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
		createdAtColumn,
		lastUpdatedAtColumn,
		archivedAtColumn,
		belongsToHouseholdColumn,
	}
)

func buildWebhooksQueries() []*Query {
	insertColumns := filterForInsert(webhooksColumns)
	fullSelectColumns := mergeColumns(
		applyToEach(webhooksColumns, func(s string) string {
			return fullColumnName(webhooksTableName, s)
		}),
		applyToEach(webhookTriggerEventsColumns, func(s string) string {
			return fullColumnName(webhookTriggerEventsTableName, s)
		}), 5)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveWebhook",
				Type: ExecRowsType,
			},
			Content: formatQuery(
				buildArchiveQuery(webhooksTableName, belongsToHouseholdColumn),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateWebhook",
				Type: ExecType,
			},
			Content: formatQuery(
				buildCreateQuery(webhooksTableName, insertColumns),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckWebhookExistence",
				Type: OneType,
			},
			Content: formatQuery(
				buildExistenceCheckQuery(
					webhooksTableName,
					" AND webhooks.belongs_to_household = sqlc.arg(household_id)",
				),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhooksForHousehold",
				Type: ManyType,
			},
			Content: formatQuery(
				buildListQuery(
					webhooksTableName,
					fullSelectColumns,
					[]string{webhookTriggerEventsJoin},
					false,
					belongsToHouseholdColumn,
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
					"webhook_trigger_events.archived_at IS NULL",
				),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhooksForHouseholdAndEvent",
				Type: ManyType,
			},
			Content: formatQuery(
				buildSelectQuery(webhooksTableName, applyToEach(webhooksColumns, func(s string) string {
					return fullColumnName(webhooksTableName, s)
				}), []string{webhookTriggerEventsJoin},
					false,
					false,
					"webhook_trigger_events.archived_at IS NULL",
					"webhook_trigger_events.trigger_event = sqlc.arg(trigger_event)",
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
				),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetWebhook",
				Type: ManyType,
			},
			Content: formatQuery(
				buildSelectQuery(
					webhooksTableName,
					fullSelectColumns,
					[]string{webhookTriggerEventsJoin},
					true,
					true,
					"webhook_trigger_events.archived_at IS NULL",
					"webhooks.belongs_to_household = sqlc.arg(household_id)",
				),
			),
		},
	}
}
