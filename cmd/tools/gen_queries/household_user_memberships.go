package main

const householdUserMembershipsTableName = "household_user_memberships"

var householdUserMembershipsColumns = []string{
	idColumn,
	belongsToHouseholdColumn,
	"belongs_to_user",
	"default_household",
	"household_role",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
