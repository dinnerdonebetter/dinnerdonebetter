-- name: GetAPIClientByID :one
SELECT
	api_clients.id,
	api_clients.name,
	api_clients.client_id,
	api_clients.secret_key,
	api_clients.created_at,
	api_clients.last_updated_at,
	api_clients.archived_at,
	api_clients.belongs_to_user
FROM api_clients
WHERE api_clients.archived_at IS NULL
	AND api_clients.belongs_to_user = $1
	AND api_clients.id = $2;
