package main

const householdInstrumentOwnershipsTableName = "household_instrument_ownerships"

var householdInstrumentOwnershipsColumns = []string{
	idColumn,
	"notes",
	"quantity",
	"valid_instrument_id",
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
