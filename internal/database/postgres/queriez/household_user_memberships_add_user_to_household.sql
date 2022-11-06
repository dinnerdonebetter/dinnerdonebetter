-- name: CreateHouseholdUserMembership :exec
INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_roles)
VALUES ($1,$2,$3,$4);
