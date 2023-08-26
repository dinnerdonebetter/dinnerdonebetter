-- name: AddToHouseholdDuringCreation :exec

INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_role)
VALUES ($1,$2,$3,$4);

-- name: ArchiveHousehold :execrows

UPDATE households SET last_updated_at = NOW(), archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_user = $1 AND id = $2;

-- name: CreateHousehold :exec

INSERT INTO households (id,"name",billing_status,contact_phone,address_line_1,address_line_2,city,state,zip_code,country,latitude,longitude,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);

-- name: GetHouseholdByIDWithMemberships :many

SELECT
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
	users.id as user_id,
	users.first_name as user_first_name,
	users.last_name as user_last_name,
	users.username as user_username,
	users.email_address as user_email_address,
	users.email_address_verified_at as user_email_address_verified_at,
	users.avatar_src as user_avatar_src,
	users.requires_password_change as user_requires_password_change,
	users.password_last_changed_at as user_password_last_changed_at,
	users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
	users.service_role as user_service_role,
	users.user_account_status as user_user_account_status,
	users.user_account_status_explanation as user_user_account_status_explanation,
	users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
	users.created_at as user_created_at,
	users.last_updated_at as user_last_updated_at,
	users.archived_at as user_archived_at,
	household_user_memberships.id as membership_id,
	household_user_memberships.belongs_to_user as membership_belongs_to_user,
	household_user_memberships.belongs_to_household as membership_belongs_to_household,
	household_user_memberships.household_role as membership_household_role,
	household_user_memberships.default_household as membership_default_household,
	household_user_memberships.created_at as membership_created_at,
	household_user_memberships.last_updated_at as membership_last_updated_at,
	household_user_memberships.archived_at as membership_archived_at
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
	AND household_user_memberships.archived_at IS NULL
	AND households.id = $1;

-- name: GetHouseholdsForUser :many

SELECT
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
    (
        SELECT
            COUNT(households.id)
        FROM
            households
            JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
        WHERE households.archived_at IS NULL
            AND household_user_memberships.belongs_to_user = sqlc.arg(user_id)
            AND household_user_memberships.belongs_to_user = sqlc.arg(user_id)
            AND households.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND households.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (households.last_updated_at IS NULL OR households.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (households.last_updated_at IS NULL OR households.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(households.id)
        FROM
            households
        WHERE
            households.archived_at IS NULL
    ) as total_count
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
    JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
    AND household_user_memberships.archived_at IS NULL
    AND household_user_memberships.belongs_to_user = sqlc.arg(user_id)
    AND households.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND households.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (households.last_updated_at IS NULL OR households.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (households.last_updated_at IS NULL OR households.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);


-- name: UpdateHousehold :execrows

UPDATE households
SET
	name = $1,
	contact_phone = $2,
	address_line_1 = $3,
	address_line_2 = $4,
	city = $5,
	state = $6,
	zip_code = $7,
	country = $8,
	latitude = $9,
    longitude = $10,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $11
	AND id = $12;
