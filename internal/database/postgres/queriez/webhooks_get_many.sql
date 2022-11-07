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
    SELECT COUNT(webhooks.id)
     FROM webhooks
     WHERE webhooks.archived_at IS NULL
       AND webhooks.created_at > COALESCE(sqlc.narg('created_after'), 0)
       AND webhooks.created_at < COALESCE(sqlc.narg('created_before'), (SELECT NOW() + interval '9999 years'))
       AND (webhooks.last_updated_at IS NULL OR webhooks.last_updated_at > COALESCE(sqlc.narg('updated_after'), 0))
       AND (webhoks.last_updated_at IS NULL OR
            webhooks.last_updated_at < COALESCE(sqlc.narg('updated_before'), (SELECT NOW() + interval '9999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(webhooks.id)
        FROM
            webhooks
        WHERE
            webhooks.archived_at IS NULL
    ) as total_count
FROM webhooks
WHERE webhooks.archived_at IS NULL
  AND webhook_trigger_events.archived_at IS NULL
  AND webhooks.belongs_to_household = $1;