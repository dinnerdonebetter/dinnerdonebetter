package main

const validPreparationVesselsTableName = "valid_preparation_vessels"

var validPreparationVesselsColumns = []string{
	"id",
	"valid_preparation_id",
	"valid_vessel_id",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
