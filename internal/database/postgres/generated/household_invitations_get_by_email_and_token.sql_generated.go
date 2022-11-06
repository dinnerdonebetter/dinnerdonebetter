// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: household_invitations_get_by_email_and_token.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetHouseholdInvitationByEmailAndToken = `-- name: GetHouseholdInvitationByEmailAndToken :one
SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_email,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.time_zone,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_roles,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_at,
	household_invitations.last_updated_at,
	household_invitations.archived_at
FROM household_invitations
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.to_email = LOWER($1)
	AND household_invitations.token = $2
`

type GetHouseholdInvitationByEmailAndTokenParams struct {
	Lower string
	Token string
}

type GetHouseholdInvitationByEmailAndTokenRow struct {
	ID                           string
	ID_2                         string
	Name                         string
	BillingStatus                string
	ContactEmail                 string
	ContactPhone                 string
	PaymentProcessorCustomerID   string
	SubscriptionPlanID           sql.NullString
	TimeZone                     TimeZone
	CreatedAt                    time.Time
	LastUpdatedAt                sql.NullTime
	ArchivedAt                   sql.NullTime
	BelongsToUser                string
	ToEmail                      string
	ToUser                       sql.NullString
	ID_3                         string
	Username                     string
	EmailAddress                 string
	AvatarSrc                    sql.NullString
	HashedPassword               string
	RequiresPasswordChange       bool
	PasswordLastChangedAt        sql.NullTime
	TwoFactorSecret              string
	TwoFactorSecretVerifiedAt    sql.NullTime
	ServiceRoles                 string
	UserAccountStatus            string
	UserAccountStatusExplanation string
	BirthDay                     sql.NullInt16
	BirthMonth                   sql.NullInt16
	CreatedAt_2                  time.Time
	LastUpdatedAt_2              sql.NullTime
	ArchivedAt_2                 sql.NullTime
	Status                       InvitationState
	Note                         string
	StatusNote                   string
	Token                        string
	CreatedAt_3                  time.Time
	LastUpdatedAt_3              sql.NullTime
	ArchivedAt_3                 sql.NullTime
}

func (q *Queries) GetHouseholdInvitationByEmailAndToken(ctx context.Context, arg *GetHouseholdInvitationByEmailAndTokenParams) (*GetHouseholdInvitationByEmailAndTokenRow, error) {
	row := q.db.QueryRowContext(ctx, GetHouseholdInvitationByEmailAndToken, arg.Lower, arg.Token)
	var i GetHouseholdInvitationByEmailAndTokenRow
	err := row.Scan(
		&i.ID,
		&i.ID_2,
		&i.Name,
		&i.BillingStatus,
		&i.ContactEmail,
		&i.ContactPhone,
		&i.PaymentProcessorCustomerID,
		&i.SubscriptionPlanID,
		&i.TimeZone,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToUser,
		&i.ToEmail,
		&i.ToUser,
		&i.ID_3,
		&i.Username,
		&i.EmailAddress,
		&i.AvatarSrc,
		&i.HashedPassword,
		&i.RequiresPasswordChange,
		&i.PasswordLastChangedAt,
		&i.TwoFactorSecret,
		&i.TwoFactorSecretVerifiedAt,
		&i.ServiceRoles,
		&i.UserAccountStatus,
		&i.UserAccountStatusExplanation,
		&i.BirthDay,
		&i.BirthMonth,
		&i.CreatedAt_2,
		&i.LastUpdatedAt_2,
		&i.ArchivedAt_2,
		&i.Status,
		&i.Note,
		&i.StatusNote,
		&i.Token,
		&i.CreatedAt_3,
		&i.LastUpdatedAt_3,
		&i.ArchivedAt_3,
	)
	return &i, err
}
