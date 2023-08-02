-- name: ArchiveOAuth2ClientTokenByCode :exec

DELETE FROM oauth2_client_tokens WHERE code = $1;
