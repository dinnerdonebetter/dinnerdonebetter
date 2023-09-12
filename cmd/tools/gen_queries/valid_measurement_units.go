package main

const validMeasurementUnitsTableName = "valid_measurement_units"

var validMeasurementUnitsColumns = []string{
	"id",
	"name",
	"description",
	"icon_path",
	"volumetric",
	"universal",
	"metric",
	"imperial",
	"plural_name",
	"slug",
	"last_indexed_at",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
