-- name: ArchiveOAuth2Client :exec

UPDATE oauth2_clients SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1;

-- name: CreateOAuth2Client :exec

INSERT INTO oauth2_clients (id,"name",client_id,client_secret) VALUES ($1,$2,$3,$4);

-- name: GetOAuth2ClientByClientID :one

SELECT
	oauth2_clients.id,
	oauth2_clients.name,
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
	oauth2_clients.client_id,
	oauth2_clients.client_secret,
	oauth2_clients.created_at,
	oauth2_clients.archived_at
FROM oauth2_clients
WHERE oauth2_clients.archived_at IS NULL
	AND oauth2_clients.id = $1;
