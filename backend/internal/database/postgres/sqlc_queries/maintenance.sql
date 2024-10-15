-- name: DeleteExpiredOAuth2ClientTokens :execrows

DELETE FROM oauth2_client_tokens WHERE code_expires_at < (NOW() - interval '1 day') AND access_expires_at < (NOW() - interval '1 day') AND refresh_expires_at < (NOW() - interval '1 day');
