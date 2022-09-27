UPDATE household_user_memberships
SET default_household = (belongs_to_user = $1 AND belongs_to_household = $2)
WHERE archived_at IS NULL
	AND belongs_to_user = $3;
