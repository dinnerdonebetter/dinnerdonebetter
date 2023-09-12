package main

const passwordResetTokensTableName = "password_reset_tokens"

var passwordResetTokensColumns = []string{
	"id",
	"token",
	"expires_at",
	"redeemed_at",
	"belongs_to_user",
	createdAtColumn,
	lastUpdatedAtColumn,
}
