-- name: GetAPIClientByClientID :one
	SELECT
		api_clients.id,
		api_clients.name,
		api_clients.client_id,
		api_clients.secret_key,
		api_clients.created_on,
		api_clients.last_updated_on,
		api_clients.archived_on,
		api_clients.belongs_to_user
	FROM api_clients
	WHERE api_clients.archived_on IS NULL
	AND api_clients.client_id = $1;

-- name: GetAPIClientByDatabaseID :one
	SELECT
		api_clients.id,
		api_clients.name,
		api_clients.client_id,
		api_clients.secret_key,
		api_clients.created_on,
		api_clients.last_updated_on,
		api_clients.archived_on,
		api_clients.belongs_to_user
	FROM api_clients
	WHERE api_clients.archived_on IS NULL
	AND api_clients.belongs_to_user = $1
	AND api_clients.id = $2;

-- name: GetTotalAPIClientCount :exec
	SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL;

-- name: CreateAPIClient :exec
	INSERT INTO api_clients (id,name,client_id,secret_key,belongs_to_user) VALUES ($1,$2,$3,$4,$5);

-- name: ArchiveAPIClient :exec
	UPDATE api_clients SET
		last_updated_on = extract(epoch FROM NOW()),
		archived_on = extract(epoch FROM NOW())
	WHERE archived_on IS NULL
	AND belongs_to_user = $1 AND id = $2;
