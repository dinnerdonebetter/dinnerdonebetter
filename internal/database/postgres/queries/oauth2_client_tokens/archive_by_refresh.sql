-- name: ArchiveOAuth2ClientTokenByRefresh :exec

DELETE FROM oauth2_client_tokens WHERE refresh = $1;
