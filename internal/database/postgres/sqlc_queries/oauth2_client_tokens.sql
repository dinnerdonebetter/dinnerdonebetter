-- name: ArchiveOAuth2ClientTokenByAccess :exec

DELETE FROM oauth2_client_tokens WHERE access = $1;

-- name: ArchiveOAuth2ClientTokenByCode :exec

DELETE FROM oauth2_client_tokens WHERE code = $1;

-- name: ArchiveOAuth2ClientTokenByRefresh :exec

DELETE FROM oauth2_client_tokens WHERE refresh = $1;

-- name: CreateOAuth2ClientToken :exec

INSERT INTO oauth2_client_tokens (id,client_id,belongs_to_user,redirect_uri,scope,code,code_challenge,code_challenge_method,code_created_at,code_expires_at,access,access_created_at,access_expires_at,refresh,refresh_created_at,refresh_expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);

-- name: CheckOAuth2ClientTokenExistence :one

SELECT EXISTS ( SELECT oauth2_client_tokens.id FROM oauth2_client_tokens WHERE oauth2_client_tokens.archived_at IS NULL AND oauth2_client_tokens.id = $1 );

-- name: GetOAuth2ClientTokenByAccess :one

SELECT
	oauth2_client_tokens.id,
	oauth2_client_tokens.client_id,
	oauth2_client_tokens.belongs_to_user,
	oauth2_client_tokens.redirect_uri,
	oauth2_client_tokens.scope,
	oauth2_client_tokens.code,
	oauth2_client_tokens.code_challenge,
	oauth2_client_tokens.code_challenge_method,
	oauth2_client_tokens.code_created_at,
	oauth2_client_tokens.code_expires_at,
	oauth2_client_tokens.access,
	oauth2_client_tokens.access_created_at,
	oauth2_client_tokens.access_expires_at,
	oauth2_client_tokens.refresh,
	oauth2_client_tokens.refresh_created_at,
	oauth2_client_tokens.refresh_expires_at
FROM oauth2_client_tokens
WHERE oauth2_client_tokens.access = $1;

-- name: GetOAuth2ClientTokenByCode :one

SELECT
	oauth2_client_tokens.id,
	oauth2_client_tokens.client_id,
	oauth2_client_tokens.belongs_to_user,
	oauth2_client_tokens.redirect_uri,
	oauth2_client_tokens.scope,
	oauth2_client_tokens.code,
	oauth2_client_tokens.code_challenge,
	oauth2_client_tokens.code_challenge_method,
	oauth2_client_tokens.code_created_at,
	oauth2_client_tokens.code_expires_at,
	oauth2_client_tokens.access,
	oauth2_client_tokens.access_created_at,
	oauth2_client_tokens.access_expires_at,
	oauth2_client_tokens.refresh,
	oauth2_client_tokens.refresh_created_at,
	oauth2_client_tokens.refresh_expires_at
FROM oauth2_client_tokens
WHERE oauth2_client_tokens.code = $1;

-- name: GetOAuth2ClientTokenByRefresh :one

SELECT
	oauth2_client_tokens.id,
	oauth2_client_tokens.client_id,
	oauth2_client_tokens.belongs_to_user,
	oauth2_client_tokens.redirect_uri,
	oauth2_client_tokens.scope,
	oauth2_client_tokens.code,
	oauth2_client_tokens.code_challenge,
	oauth2_client_tokens.code_challenge_method,
	oauth2_client_tokens.code_created_at,
	oauth2_client_tokens.code_expires_at,
	oauth2_client_tokens.access,
	oauth2_client_tokens.access_created_at,
	oauth2_client_tokens.access_expires_at,
	oauth2_client_tokens.refresh,
	oauth2_client_tokens.refresh_created_at,
	oauth2_client_tokens.refresh_expires_at
FROM oauth2_client_tokens
WHERE oauth2_client_tokens.refresh = $1;
