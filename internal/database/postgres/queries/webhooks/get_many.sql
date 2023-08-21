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
