package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	validInstrumentsTable = "valid_instruments"
)

var (
	validInstrumentsTableColumns = []string{
		buildColumnName(validInstrumentsTable, id),
		buildColumnName(validInstrumentsTable, "name"),
		buildColumnName(validInstrumentsTable, "plural_name"),
		buildColumnName(validInstrumentsTable, "description"),
		buildColumnName(validInstrumentsTable, "icon_path"),
		buildColumnName(validInstrumentsTable, "usable_for_storage"),
		buildColumnName(validInstrumentsTable, "display_in_summary_lists"),
		buildColumnName(validInstrumentsTable, "include_in_generated_instructions"),
		buildColumnName(validInstrumentsTable, "slug"),
		buildColumnName(validInstrumentsTable, createdAt),
		buildColumnName(validInstrumentsTable, lastUpdatedAt),
		buildColumnName(validInstrumentsTable, archivedAt),
	}
)

func buildSelectValidInstrumentsNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(validInstrumentsTable)
}

func buildSelectValidInstrumentQuery(_ squirrel.StatementBuilderType) string {
	return buildSelectQuery(validInstrumentsTable, validInstrumentsTableColumns)
}
