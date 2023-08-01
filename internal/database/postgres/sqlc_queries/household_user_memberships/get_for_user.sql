SELECT
	household_user_memberships.id,
	household_user_memberships.belongs_to_user,
	household_user_memberships.belongs_to_household,
	household_user_memberships.household_role,
	household_user_memberships.default_household,
	household_user_memberships.created_at,
	household_user_memberships.last_updated_at,
	household_user_memberships.archived_at
FROM household_user_memberships
	JOIN households ON households.id = household_user_memberships.belongs_to_household
WHERE household_user_memberships.archived_at IS NULL
	AND household_user_memberships.belongs_to_user = $1;
