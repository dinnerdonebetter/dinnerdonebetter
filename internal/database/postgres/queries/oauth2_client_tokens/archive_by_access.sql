-- name: ArchiveOAuth2ClientTokenByAccess :exec

DELETE FROM oauth2_client_tokens WHERE access = $1;
