package main

const validInstrumentsTableName = "valid_instruments"

var validInstrumentsColumns = []string{
	"id",
	"name",
	"description",
	"icon_path",
	"plural_name",
	"usable_for_storage",
	"slug",
	"display_in_summary_lists",
	"include_in_generated_instructions",
	"last_indexed_at",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
