SELECT
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
	users.id,
	users.username,
	users.email_address,
	users.avatar_src,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birth_day,
	users.birth_month,
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_user_memberships.id,
	household_user_memberships.belongs_to_user,
	household_user_memberships.belongs_to_household,
	household_user_memberships.household_role,
	household_user_memberships.default_household,
	household_user_memberships.created_at,
	household_user_memberships.last_updated_at,
	household_user_memberships.archived_at
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
	AND household_user_memberships.archived_at IS NULL
	AND households.id = $1;
