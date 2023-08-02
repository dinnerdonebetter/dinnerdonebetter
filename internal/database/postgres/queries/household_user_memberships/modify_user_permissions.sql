-- name: ModifyHouseholdUserPermissions :exec

UPDATE household_user_memberships SET household_role = $1 WHERE belongs_to_household = $2 AND belongs_to_user = $3;