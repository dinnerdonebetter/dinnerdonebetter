package main

const usersTableName = "users"

var usersColumns = []string{
	"id",
	"username",
	"avatar_src",
	"email_address",
	"hashed_password",
	"password_last_changed_at",
	"requires_password_change",
	"two_factor_secret",
	"two_factor_secret_verified_at",
	"service_role",
	"user_account_status",
	"user_account_status_explanation",
	"birthday",
	"email_address_verification_token",
	"email_address_verified_at",
	"first_name",
	"last_name",
	"last_accepted_terms_of_service",
	"last_accepted_privacy_policy",
	"last_indexed_at",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
