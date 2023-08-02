-- name: ArchiveOAuth2Client :exec

UPDATE oauth2_clients SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;