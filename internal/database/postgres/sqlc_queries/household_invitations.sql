-- name: AttachHouseholdInvitationsToUserID :exec

UPDATE household_invitations SET
    to_user = sqlc.arg(user_id),
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND to_email = LOWER(sqlc.arg(email_address));


-- name: CreateHouseholdInvitation :exec

INSERT INTO household_invitations (id,from_user,to_user,to_name,note,to_email,token,destination_household,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);


-- name: CheckHouseholdInvitationExistence :one

SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_at IS NULL AND household_invitations.id = $1 );

-- name: GetHouseholdInvitationByEmailAndToken :one

SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
    households.longitude,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
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
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_invitations.to_name,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.expires_at,
	household_invitations.created_at,
	household_invitations.last_updated_at,
	household_invitations.archived_at
FROM household_invitations
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.to_email = LOWER($1)
	AND household_invitations.token = $2;


-- name: GetHouseholdInvitationByHouseholdAndID :one

SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
    households.longitude,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
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
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_invitations.to_name,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.expires_at,
	household_invitations.created_at,
	household_invitations.last_updated_at,
	household_invitations.archived_at
FROM household_invitations
	LEFT JOIN households ON household_invitations.destination_household = households.id
	LEFT JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
      AND household_invitations.expires_at > NOW()
      AND household_invitations.expires_at > NOW()
	AND household_invitations.destination_household = $1
	AND household_invitations.id = $2;


-- name: GetHouseholdInvitationByTokenAndID :one

SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
    households.longitude,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
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
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_invitations.to_name,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.expires_at,
	household_invitations.created_at,
	household_invitations.last_updated_at,
	household_invitations.archived_at
FROM household_invitations
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.token = $1
	AND household_invitations.id = $2;


SELECT
    household_invitations.id,
    households.id,
    households.name,
    households.billing_status,
    households.contact_phone,
    households.address_line_1,
    households.address_line_2,
    households.city,
    households.state,
    households.zip_code,
    households.country,
    households.latitude,
    households.longitude,
    households.payment_processor_customer_id,
    households.subscription_plan_id,
    households.created_at,
    households.last_updated_at,
    households.archived_at,
    households.belongs_to_user,
    household_invitations.to_email,
    household_invitations.to_user,
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
    users.created_at,
    users.last_updated_at,
    users.archived_at,
    household_invitations.to_name,
    household_invitations.status,
    household_invitations.note,
    household_invitations.status_note,
    household_invitations.token,
    household_invitations.expires_at,
    household_invitations.created_at,
    household_invitations.last_updated_at,
    household_invitations.archived_at,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.to_user = $1
          AND household_invitations.status = $2
          AND household_invitations.created_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
          AND household_invitations.created_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE($5, (SELECT NOW() - interval '999 years')))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE($6, (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.to_user = $1
          AND household_invitations.status = $2
    ) as total_count
FROM household_invitations
    JOIN households ON household_invitations.destination_household = households.id
    JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
  AND household_invitations.to_user = $1
  AND household_invitations.status = $2
  AND household_invitations.created_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
  AND household_invitations.created_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
  AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE($5, (SELECT NOW() - interval '999 years')))
  AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE($6, (SELECT NOW() + interval '999 years')));


SELECT
    household_invitations.id,
    households.id,
    households.name,
    households.billing_status,
    households.contact_phone,
    households.address_line_1,
    households.address_line_2,
    households.city,
    households.state,
    households.zip_code,
    households.country,
    households.latitude,
    households.longitude,
    households.payment_processor_customer_id,
    households.subscription_plan_id,
    households.created_at,
    households.last_updated_at,
    households.archived_at,
    households.belongs_to_user,
    household_invitations.to_email,
    household_invitations.to_user,
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
    users.created_at,
    users.last_updated_at,
    users.archived_at,
    household_invitations.to_name,
    household_invitations.status,
    household_invitations.note,
    household_invitations.status_note,
    household_invitations.token,
    household_invitations.expires_at,
    household_invitations.created_at,
    household_invitations.last_updated_at,
    household_invitations.archived_at,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.from_user = $1
          AND household_invitations.status = $2
          AND household_invitations.created_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
          AND household_invitations.created_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE($5, (SELECT NOW() - interval '999 years')))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE($6, (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.from_user = $1
          AND household_invitations.status = $2
    ) as total_count
FROM household_invitations
    JOIN households ON household_invitations.destination_household = households.id
    JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
  AND household_invitations.from_user = $1
  AND household_invitations.status = $2
  AND household_invitations.created_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
  AND household_invitations.created_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
  AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE($5, (SELECT NOW() - interval '999 years')))
  AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE($6, (SELECT NOW() + interval '999 years')));


-- name: SetHouseholdInvitationStatus :exec

UPDATE household_invitations SET
	status = $1,
	status_note = $2,
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $3;
