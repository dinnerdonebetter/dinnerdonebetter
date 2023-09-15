package main

/* #nosec G101 */
const oauth2ClientTokensTableName = "oauth2_client_tokens"

/* #nosec G101 */
var oauth2ClientTokensColumns = []string{
	idColumn,
	"client_id",
	"belongs_to_user",
	"redirect_uri",
	"scope",
	"code",
	"code_challenge",
	"code_challenge_method",
	"code_created_at",
	"code_expires_at",
	"access",
	"access_created_at",
	"access_expires_at",
	"refresh",
	"refresh_created_at",
	"refresh_expires_at",
}

func buildOAuth2ClientTokensQueries() []*Query {
	return []*Query{
		//
	}
}
