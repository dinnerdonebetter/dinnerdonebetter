-- name: HouseholdInvitationExists :one
SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_on IS NULL AND household_invitations.id = $1 );

-- name: GetHouseholdInvitationByHouseholdAndID :many
SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_email,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_on,
	households.last_updated_on,
	households.archived_on,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_on,
	users.two_factor_secret,
	users.two_factor_secret_verified_on,
	users.service_roles,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_on,
	users.last_updated_on,
	users.archived_on,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations
JOIN households ON household_invitations.destination_household = households.id
JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_on IS NULL
AND household_invitations.destination_household = $1
AND household_invitations.id = $2;

-- name: GetHouseholdInvitationByTokenAndID :one
SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_email,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_on,
	households.last_updated_on,
	households.archived_on,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_on,
	users.two_factor_secret,
	users.two_factor_secret_verified_on,
	users.service_roles,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_on,
	users.last_updated_on,
	users.archived_on,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations
JOIN households ON household_invitations.destination_household = households.id
JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_on IS NULL
AND household_invitations.token = $1
AND household_invitations.id = $2;

-- name: GetHouseholdInvitationByEmailAndToken :one
SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_email,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_on,
	households.last_updated_on,
	households.archived_on,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_on,
	users.two_factor_secret,
	users.two_factor_secret_verified_on,
	users.service_roles,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_on,
	users.last_updated_on,
	users.archived_on,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.created_on,
	household_invitations.last_updated_on,
	household_invitations.archived_on
FROM household_invitations
JOIN households ON household_invitations.destination_household = households.id
JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_on IS NULL
AND household_invitations.to_email = LOWER($1)
AND household_invitations.token = $2;

-- name: GetAllHouseholdInvitationsCount :one
SELECT COUNT(household_invitations.id) FROM household_invitations WHERE household_invitations.archived_on IS NULL;

-- name: CreateHouseholdInvitation :exec
INSERT INTO household_invitations (id,from_user,to_user,note,to_email,token,destination_household) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: SetInvitationStatus :exec
UPDATE household_invitations SET
	status = $1,
	status_note = $2,
	last_updated_on = extract(epoch FROM NOW()),
	archived_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND id = $3;

-- name: AttachInvitationsToUserID :exec
UPDATE household_invitations SET
	to_user = $1,
	last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
AND to_email = LOWER($2);
