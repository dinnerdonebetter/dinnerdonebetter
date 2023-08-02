-- name: CreateOAuth2Client :exec

INSERT INTO oauth2_clients (id,"name",client_id,client_secret) VALUES ($1,$2,$3,$4);
