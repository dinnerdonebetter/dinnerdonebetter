package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	validMeasurementUnitsTable = "valid_measurement_units"
)

var (
	// validMeasurementUnitsTableColumns are the columns for the valid_measurement_units table.
	validMeasurementUnitsTableColumns = []string{
		buildColumnName(validMeasurementUnitsTable, id),
		buildColumnName(validMeasurementUnitsTable, "name"),
		buildColumnName(validMeasurementUnitsTable, "description"),
		buildColumnName(validMeasurementUnitsTable, "volumetric"),
		buildColumnName(validMeasurementUnitsTable, "icon_path"),
		buildColumnName(validMeasurementUnitsTable, "universal"),
		buildColumnName(validMeasurementUnitsTable, "metric"),
		buildColumnName(validMeasurementUnitsTable, "imperial"),
		buildColumnName(validMeasurementUnitsTable, "slug"),
		buildColumnName(validMeasurementUnitsTable, "plural_name"),
		buildColumnName(validMeasurementUnitsTable, createdAt),
		buildColumnName(validMeasurementUnitsTable, lastUpdatedAt),
		buildColumnName(validMeasurementUnitsTable, archivedAt),
	}
)

func buildSelectValidMeasurementUnitsNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(validMeasurementUnitsTable)
}

func buildSelectValidMeasurementUnitQuery(_ squirrel.StatementBuilderType) string {
	return buildSelectQuery(validMeasurementUnitsTable, validMeasurementUnitsTableColumns)
}
