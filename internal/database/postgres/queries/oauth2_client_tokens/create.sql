-- name: CreateOAuth2ClientToken :exec

INSERT INTO oauth2_client_tokens (id,client_id,belongs_to_user,redirect_uri,scope,code,code_challenge,code_challenge_method,code_created_at,code_expires_at,access,access_created_at,access_expires_at,refresh,refresh_created_at,refresh_expires_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);
