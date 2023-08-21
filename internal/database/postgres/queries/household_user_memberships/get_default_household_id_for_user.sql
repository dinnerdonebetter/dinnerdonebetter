-- name: GetDefaultHouseholdIDForUser :one

SELECT households.id
FROM households
 JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
WHERE household_user_memberships.belongs_to_user = $1
	AND household_user_memberships.default_household = TRUE;
