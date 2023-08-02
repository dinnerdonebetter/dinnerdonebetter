// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_one_by_refresh.sql

package generated

import (
	"context"
)

const GetOAuth2ClientTokenByRefresh = `-- name: GetOAuth2ClientTokenByRefresh :one

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
WHERE oauth2_client_tokens.refresh = $1
`

func (q *Queries) GetOAuth2ClientTokenByRefresh(ctx context.Context, db DBTX, refresh string) (*Oauth2ClientTokens, error) {
	row := db.QueryRowContext(ctx, GetOAuth2ClientTokenByRefresh, refresh)
	var i Oauth2ClientTokens
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.BelongsToUser,
		&i.RedirectUri,
		&i.Scope,
		&i.Code,
		&i.CodeChallenge,
		&i.CodeChallengeMethod,
		&i.CodeCreatedAt,
		&i.CodeExpiresAt,
		&i.Access,
		&i.AccessCreatedAt,
		&i.AccessExpiresAt,
		&i.Refresh,
		&i.RefreshCreatedAt,
		&i.RefreshExpiresAt,
	)
	return &i, err
}
