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
