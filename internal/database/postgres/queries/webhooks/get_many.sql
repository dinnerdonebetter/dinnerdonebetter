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
	      AND webhooks.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
	      AND webhooks.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
	      AND (
	            webhooks.last_updated_at IS NULL
	            OR webhooks.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
	        )
	      AND (
	            webhooks.last_updated_at IS NULL
	            OR webhooks.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
	        )
	    OFFSET COALESCE($6, 0)
	) as filtered_count,
	(
	    SELECT
	        COUNT(webhooks.id)
	    FROM
	        webhooks
	    WHERE
	        webhooks.archived_at IS NULL
	) as total_count
FROM
	webhooks
	    JOIN webhook_trigger_events ON webhook_trigger_events.belongs_to_webhook = webhooks.id
WHERE
	webhooks.archived_at IS NULL
	AND webhook_trigger_events.archived_at IS NULL
	AND webhooks.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
	AND webhooks.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
	AND (
	    webhooks.last_updated_at IS NULL
	    OR webhooks.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
	)
	AND (
	    webhooks.last_updated_at IS NULL
	    OR webhooks.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
	)
	AND webhooks.belongs_to_household = $1
	OFFSET COALESCE($6, 0);
