package main

const householdInstrumentOwnershipsTableName = "household_instrument_ownerships"

var householdInstrumentOwnershipsColumns = []string{
	"id",
	"notes",
	"quantity",
	"valid_instrument_id",
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
