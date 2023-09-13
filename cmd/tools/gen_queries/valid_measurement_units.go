package main

const validMeasurementUnitsTableName = "valid_measurement_units"

var validMeasurementUnitsColumns = []string{
	idColumn,
	"name",
	"description",
	"icon_path",
	"volumetric",
	"universal",
	"metric",
	"imperial",
	"plural_name",
	"slug",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
