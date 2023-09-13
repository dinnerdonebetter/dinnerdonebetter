-- name: ArchiveWebhook :execrows

UPDATE webhooks
   SET archived_at = NOW()
 WHERE archived_at IS NULL
   AND id = sqlc.arg(id)
   AND belongs_to_household = sqlc.arg(household_id);

-- name: CreateWebhook :exec

INSERT
  INTO webhooks (id, name, content_type, url, method, belongs_to_household)
VALUES (
		sqlc.arg(id),
		sqlc.arg(name),
		sqlc.arg(content_type),
		sqlc.arg(url),
		sqlc.arg(method),
		sqlc.arg(belongs_to_household)
       );

-- name: CheckWebhookExistence :one

SELECT EXISTS(
		SELECT webhooks.id
		  FROM webhooks
		 WHERE webhooks.archived_at IS NULL
		   AND webhooks.id = sqlc.arg(id)
		   AND webhooks.belongs_to_household = sqlc.arg(household_id)
       );

-- name: GetWebhooksForHousehold :many

SELECT
    webhooks.id,
    webhooks.name,
    webhooks.content_type,
    webhooks.url,
    webhooks.method,
    webhook_trigger_events.id,
    webhook_trigger_events.trigger_event,
    webhook_trigger_events.belongs_to_webhook,
    webhook_trigger_events.created_at,
    webhook_trigger_events.archived_at,
    webhooks.created_at,
    webhooks.last_updated_at,
    webhooks.archived_at,
    webhooks.belongs_to_household,
    (
        SELECT
            COUNT(webhooks.id)
        FROM
            webhooks
        WHERE
            webhooks.archived_at IS NULL
            AND webhooks.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND webhooks.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (
                webhooks.last_updated_at IS NULL
                OR webhooks.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
            AND (
                webhooks.last_updated_at IS NULL
                OR webhooks.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
            AND webhooks.belongs_to_household = sqlc.arg(household_id)
    ) as filtered_count,
    (
        SELECT
            COUNT(webhooks.id)
        FROM
            webhooks
        WHERE
            webhooks.archived_at IS NULL
            AND webhooks.belongs_to_household = sqlc.arg(household_id)
    ) as total_count
FROM webhooks
    JOIN webhook_trigger_events ON webhook_trigger_events.belongs_to_webhook = webhooks.id
WHERE webhook_trigger_events.archived_at IS NULL
    AND webhooks.created_at > (COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years')))
    AND webhooks.created_at < (COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years')))
    AND (
        webhooks.last_updated_at IS NULL
        OR webhooks.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
    AND (
        webhooks.last_updated_at IS NULL
        OR webhooks.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
    AND webhooks.archived_at IS NULL
    AND webhooks.belongs_to_household = sqlc.arg(household_id)
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: GetWebhooksForHouseholdAndEvent :many

SELECT
    webhooks.id,
    webhooks.name,
    webhooks.content_type,
    webhooks.url,
    webhooks.method,
    webhooks.created_at,
    webhooks.last_updated_at,
    webhooks.archived_at,
    webhooks.belongs_to_household
FROM webhooks
    JOIN webhook_trigger_events ON webhook_trigger_events.belongs_to_webhook = webhooks.id
WHERE webhooks.archived_at IS NULL
    AND webhook_trigger_events.archived_at IS NULL
    AND webhook_trigger_events.trigger_event = sqlc.arg(trigger_event)
    AND webhooks.belongs_to_household = sqlc.arg(household_id);

-- name: GetWebhooks :many

SELECT
	webhooks.id,
	webhooks.name,
	webhooks.content_type,
	webhooks.url,
	webhooks.method,
	webhooks.created_at,
	webhooks.last_updated_at,
	webhooks.archived_at,
	webhooks.belongs_to_household,
	(
	    SELECT
	        COUNT(webhooks.id)
	    FROM
	        webhooks
	    WHERE
            webhooks.archived_at IS NULL
            AND webhooks.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND webhooks.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
            AND (
                webhooks.last_updated_at IS NULL
                OR webhooks.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
            AND (
                webhooks.last_updated_at IS NULL
                OR webhooks.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )
            AND webhooks.belongs_to_household = sqlc.arg(household_id)
	) as filtered_count,
	(
	    SELECT
	        COUNT(webhooks.id)
	    FROM
	        webhooks
	    WHERE
	        webhooks.archived_at IS NULL
            AND webhooks.belongs_to_household = sqlc.arg(household_id)
	) as total_count
FROM
	webhooks
WHERE
	webhooks.archived_at IS NULL
	AND webhooks.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
	AND webhooks.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
	AND (
	    webhooks.last_updated_at IS NULL
	    OR webhooks.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
	)
	AND (
	    webhooks.last_updated_at IS NULL
	    OR webhooks.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
	)
	AND webhooks.belongs_to_household = sqlc.arg(household_id)
	OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: GetWebhook :many

SELECT webhooks.id AS webhook_id,
       webhooks.name AS webhook_name,
       webhooks.content_type AS webhook_content_type,
       webhooks.url AS webhook_url,
       webhooks.method AS webhook_method,
       webhook_trigger_events.id AS webhook_trigger_event_id,
       webhook_trigger_events.trigger_event
		AS webhook_trigger_event_trigger_event,
       webhook_trigger_events.belongs_to_webhook
		AS webhook_trigger_event_belongs_to_webhook,
       webhook_trigger_events.created_at AS webhook_trigger_event_created_at,
       webhook_trigger_events.archived_at AS webhook_trigger_event_archived_at,
       webhooks.created_at AS webhook_created_at,
       webhooks.last_updated_at AS webhook_last_updated_at,
       webhooks.archived_at AS webhook_archived_at,
       webhooks.belongs_to_household AS webhook_belongs_to_household
  FROM webhooks
  JOIN webhook_trigger_events ON webhooks.id
                                 = webhook_trigger_events.belongs_to_webhook
 WHERE webhook_trigger_events.archived_at IS NULL
   AND webhooks.belongs_to_household = sqlc.arg(household_id)
   AND webhooks.archived_at IS NULL
   AND webhooks.id = sqlc.arg(id);