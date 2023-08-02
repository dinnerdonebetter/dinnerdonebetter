// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_by_email_verification_token.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetUserByEmailAddressVerificationToken = `-- name: GetUserByEmailAddressVerificationToken :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.email_address_verification_token = $1
`

type GetUserByEmailAddressVerificationTokenRow struct {
	CreatedAt                    time.Time      `db:"created_at"`
	Birthday                     sql.NullTime   `db:"birthday"`
	ArchivedAt                   sql.NullTime   `db:"archived_at"`
	PasswordLastChangedAt        sql.NullTime   `db:"password_last_changed_at"`
	LastUpdatedAt                sql.NullTime   `db:"last_updated_at"`
	EmailAddressVerifiedAt       sql.NullTime   `db:"email_address_verified_at"`
	LastAcceptedPrivacyPolicy    sql.NullTime   `db:"last_accepted_privacy_policy"`
	LastAcceptedTermsOfService   sql.NullTime   `db:"last_accepted_terms_of_service"`
	TwoFactorSecretVerifiedAt    sql.NullTime   `db:"two_factor_secret_verified_at"`
	UserAccountStatusExplanation string         `db:"user_account_status_explanation"`
	FirstName                    string         `db:"first_name"`
	ServiceRole                  string         `db:"service_role"`
	UserAccountStatus            string         `db:"user_account_status"`
	LastName                     string         `db:"last_name"`
	ID                           string         `db:"id"`
	HashedPassword               string         `db:"hashed_password"`
	TwoFactorSecret              string         `db:"two_factor_secret"`
	EmailAddress                 string         `db:"email_address"`
	Username                     string         `db:"username"`
	AvatarSrc                    sql.NullString `db:"avatar_src"`
	RequiresPasswordChange       bool           `db:"requires_password_change"`
}

func (q *Queries) GetUserByEmailAddressVerificationToken(ctx context.Context, db DBTX, emailAddressVerificationToken sql.NullString) (*GetUserByEmailAddressVerificationTokenRow, error) {
	row := db.QueryRowContext(ctx, GetUserByEmailAddressVerificationToken, emailAddressVerificationToken)
	var i GetUserByEmailAddressVerificationTokenRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.EmailAddress,
		&i.EmailAddressVerifiedAt,
		&i.AvatarSrc,
		&i.HashedPassword,
		&i.RequiresPasswordChange,
		&i.PasswordLastChangedAt,
		&i.TwoFactorSecret,
		&i.TwoFactorSecretVerifiedAt,
		&i.ServiceRole,
		&i.UserAccountStatus,
		&i.UserAccountStatusExplanation,
		&i.Birthday,
		&i.LastAcceptedTermsOfService,
		&i.LastAcceptedPrivacyPolicy,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}