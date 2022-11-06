-- name: CreateAPIClient :exec
INSERT INTO api_clients (id,name,client_id,secret_key,belongs_to_user) VALUES ($1,$2,$3,$4,$5);
