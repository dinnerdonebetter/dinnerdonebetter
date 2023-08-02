-- name: TransferHouseholdMembership :exec

UPDATE household_user_memberships SET belongs_to_user = $1 WHERE archived_at IS NULL AND belongs_to_household = $2 AND belongs_to_user = $3;
