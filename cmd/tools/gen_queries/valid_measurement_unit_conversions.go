package main

const validMeasurementUnitConversionsTableName = "valid_measurement_unit_conversions"

var validMeasurementUnitConversionsColumns = []string{
	"id",
	"from_unit",
	"to_unit",
	"only_for_ingredient",
	"modifier",
	"notes",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
