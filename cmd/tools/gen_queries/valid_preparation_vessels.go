package main

const validPreparationVesselsTableName = "valid_preparation_vessels"

var validPreparationVesselsColumns = []string{
	idColumn,
	"valid_preparation_id",
	"valid_vessel_id",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
