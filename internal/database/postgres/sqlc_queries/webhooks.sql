-- name: ArchiveWebhook :execrows

UPDATE webhooks
SET
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_household = $1
	AND id = $2;

-- name: CreateWebhook :exec

INSERT INTO webhooks (
	id,
	"name",
	content_type,
	url,
	"method",
	belongs_to_household
)
VALUES
	($1, $2, $3, $4, $5, $6);

-- name: CheckWebhookExistence :one

SELECT EXISTS (
	SELECT webhooks.id
	FROM webhooks
	WHERE webhooks.archived_at IS NULL
	  AND webhooks.belongs_to_household = $1
	  AND webhooks.id = $2
);

-- name: GetWebhooksForUser :many

SELECT webhooks.id,
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
       webhooks.belongs_to_household
  FROM webhooks
       JOIN webhook_trigger_events ON
			webhook_trigger_events.belongs_to_webhook
			= webhooks.id
 WHERE webhook_trigger_events.archived_at IS NULL
       AND webhooks.archived_at IS NULL
       AND webhooks.belongs_to_household = $1
       AND webhooks.id = $2;

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

SELECT
	webhooks.id as webhook_id,
	webhooks.name as webhook_name,
	webhooks.content_type as webhook_content_type,
	webhooks.url as webhook_url,
	webhooks.method as webhook_method,
	webhook_trigger_events.id as webhook_trigger_event_id,
	webhook_trigger_events.trigger_event as webhook_trigger_event_trigger_event,
	webhook_trigger_events.belongs_to_webhook as webhook_trigger_event_belongs_to_webhook,
	webhook_trigger_events.created_at as webhook_trigger_event_created_at,
	webhook_trigger_events.archived_at as webhook_trigger_event_archived_at,
	webhooks.created_at as webhook_created_at,
	webhooks.last_updated_at as webhook_last_updated_at,
	webhooks.archived_at as webhook_archived_at,
	webhooks.belongs_to_household as webhook_belongs_to_household
FROM webhook_trigger_events
 JOIN webhooks ON webhook_trigger_events.belongs_to_webhook=webhooks.id
WHERE webhooks.archived_at IS NULL
  AND webhook_trigger_events.archived_at IS NULL
  AND webhooks.belongs_to_household = sqlc.arg(household_id)
  AND webhooks.id = sqlc.arg(webhook_id);