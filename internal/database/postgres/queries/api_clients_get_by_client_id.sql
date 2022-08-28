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