-- name: CheckOAuth2ClientTokenExistence :one

SELECT EXISTS ( SELECT oauth2_client_tokens.id FROM oauth2_client_tokens WHERE oauth2_client_tokens.archived_at IS NULL AND oauth2_client_tokens.id = $1 );