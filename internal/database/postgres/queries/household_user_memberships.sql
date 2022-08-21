-- name: GetHouseholdMembershipsForUserQuery :many
SELECT
    household_user_memberships.id,
    household_user_memberships.belongs_to_user,
    household_user_memberships.belongs_to_household,
    household_user_memberships.household_roles,
    household_user_memberships.default_household,
    household_user_memberships.created_on,
    household_user_memberships.last_updated_on,
    household_user_memberships.archived_on
FROM household_user_memberships
JOIN households ON households.id = household_user_memberships.belongs_to_household
WHERE household_user_memberships.archived_on IS NULL
AND household_user_memberships.belongs_to_user = $1;

-- name: GetDefaultHouseholdIDForUserQuery :one
SELECT households.id
    FROM households
    JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
    WHERE household_user_memberships.belongs_to_user = $1
    AND household_user_memberships.default_household = $2;

-- name: MarkHouseholdAsUserDefaultQuery :exec
UPDATE household_user_memberships
	SET default_household = (belongs_to_user = $1 AND belongs_to_household = $2)
	WHERE archived_on IS NULL
	AND belongs_to_user = $3;

-- name: UserIsMemberOfHouseholdQuery :one
SELECT EXISTS (
    SELECT household_user_memberships.id
    FROM household_user_memberships
    WHERE household_user_memberships.archived_on IS NULL
    AND household_user_memberships.belongs_to_household = $1
    AND household_user_memberships.belongs_to_user = $2
);

-- name: ModifyUserPermissionsQuery :exec
UPDATE household_user_memberships SET household_roles = $1 WHERE belongs_to_household = $2 AND belongs_to_user = $3;

-- name: TransferHouseholdOwnershipQuery :exec
UPDATE households SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_user = $2 AND id = $3;

-- name: TransferHouseholdMembershipQuery :exec
UPDATE household_user_memberships SET belongs_to_user = $1 WHERE archived_on IS NULL AND belongs_to_household = $2 AND belongs_to_user = $3;

-- name: AddUserToHouseholdQuery :exec
INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_roles)
	VALUES ($1,$2,$3,$4);

-- name: RemoveUserFromHouseholdQuery :exec
UPDATE household_user_memberships
	SET archived_on = extract(epoch from NOW()),
		default_household = 'false'
	WHERE household_user_memberships.archived_on IS NULL
	AND household_user_memberships.belongs_to_household = $1
	AND household_user_memberships.belongs_to_user = $2;
