-- name: GetHousehold :many
SELECT
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
		users.id,
		users.username,
		users.email_address,
		users.avatar_src,
		users.requires_password_change,
		users.password_last_changed_on,
		users.two_factor_secret_verified_on,
		users.service_roles,
		users.user_account_status,
		users.user_account_status_explanation,
		users.birth_day,
		users.birth_month,
		users.created_on,
		users.last_updated_on,
		users.archived_on,
		household_user_memberships.id,
		household_user_memberships.belongs_to_user,
		household_user_memberships.belongs_to_household,
		household_user_memberships.household_roles,
		household_user_memberships.default_household,
		household_user_memberships.created_on,
		household_user_memberships.last_updated_on,
		household_user_memberships.archived_on
	FROM household_user_memberships
	JOIN households ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
	WHERE households.archived_on IS NULL
	AND household_user_memberships.archived_on IS NULL
	AND households.id = $1;

-- name: GetHouseholdByID :many
SELECT
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
        users.id,
        users.username,
        users.email_address,
        users.avatar_src,
        users.requires_password_change,
        users.password_last_changed_on,
        users.two_factor_secret_verified_on,
        users.service_roles,
        users.user_account_status,
        users.user_account_status_explanation,
        users.birth_day,
        users.birth_month,
        users.created_on,
        users.last_updated_on,
        users.archived_on,
		household_user_memberships.id,
		household_user_memberships.belongs_to_user,
		household_user_memberships.belongs_to_household,
		household_user_memberships.household_roles,
		household_user_memberships.default_household,
		household_user_memberships.created_on,
		household_user_memberships.last_updated_on,
		household_user_memberships.archived_on
	FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
	WHERE households.archived_on IS NULL
	AND households.id = $1;

-- name: GetAllHouseholdsCount :one
SELECT COUNT(households.id) FROM households WHERE households.archived_on IS NULL;

-- name: CreateHousehold :exec
INSERT INTO households (id,name,billing_status,contact_email,contact_phone,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6);


-- name: AddUserToHouseholdDuringCreation :exec
INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_roles)
	VALUES ($1,$2,$3,$4);


-- name: UpdateHousehold :exec
UPDATE households SET name = $1, contact_email = $2, contact_phone = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $4 AND id = $5;

-- name: ArchiveHousehold :exec
UPDATE households SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2;
