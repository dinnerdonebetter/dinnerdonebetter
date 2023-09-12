package main

const oauth2ClientTokensTableName = "oauth2_client_tokens"

var oauth2ClientTokensColumns = []string{
	"id",
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
