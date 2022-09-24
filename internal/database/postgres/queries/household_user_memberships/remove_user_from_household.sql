UPDATE household_user_memberships
SET archived_at = NOW(),
	default_household = 'false'
WHERE household_user_memberships.archived_at IS NULL
  AND household_user_memberships.belongs_to_household = $1
  AND household_user_memberships.belongs_to_user = $2;
