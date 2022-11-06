-- name: ArchiveWebhook :exec
UPDATE webhooks SET
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_household = $1
	AND id = $2;
