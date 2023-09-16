-- name: ArchiveOAuth2Client :execrows

UPDATE oauth2_clients SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;

-- name: CreateOAuth2Client :exec

INSERT INTO oauth2_clients (
    id,
    name,
    client_id,
    client_secret
) VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: GetOAuth2ClientByClientID :one

SELECT
	oauth2_clients.id,
    oauth2_clients.name,
    oauth2_clients.description,
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.client_id = $1;

-- name: GetOAuth2ClientByDatabaseID :one

SELECT
    oauth2_clients.id,
    oauth2_clients.name,
    oauth2_clients.description,
    oauth2_clients.client_id,
    oauth2_clients.client_secret,
    oauth2_clients.created_at,
    oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
  AND oauth2_clients.id = $1;

-- name: GetOAuth2Clients :many

SELECT
    oauth2_clients.id,
    oauth2_clients.name,
    oauth2_clients.description,
    oauth2_clients.client_id,
    oauth2_clients.client_secret,
    oauth2_clients.created_at,
    oauth2_clients.archived_at,
    (
        SELECT
            COUNT(oauth2_clients.id)
        FROM
            oauth2_clients
        WHERE
            oauth2_clients.archived_at IS NULL
          AND oauth2_clients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND oauth2_clients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    ) as filtered_count,
    (
        SELECT
            COUNT(oauth2_clients.id)
        FROM
            oauth2_clients
        WHERE
            oauth2_clients.archived_at IS NULL
    ) as total_count
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
    AND oauth2_clients.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND oauth2_clients.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);
