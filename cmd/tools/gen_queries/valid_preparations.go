package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	validPreparationsTable = "valid_preparations"
)

var (
	// validPreparationsTableColumns are the columns for the valid_preparations table.
	validPreparationsTableColumns = []string{
		buildColumnName(validPreparationsTable, id),
		buildColumnName(validPreparationsTable, "name"),
		buildColumnName(validPreparationsTable, "description"),
		buildColumnName(validPreparationsTable, "icon_path"),
		buildColumnName(validPreparationsTable, "yields_nothing"),
		buildColumnName(validPreparationsTable, "restrict_to_ingredients"),
		buildColumnName(validPreparationsTable, "minimum_ingredient_count"),
		buildColumnName(validPreparationsTable, "maximum_ingredient_count"),
		buildColumnName(validPreparationsTable, "minimum_instrument_count"),
		buildColumnName(validPreparationsTable, "maximum_instrument_count"),
		buildColumnName(validPreparationsTable, "temperature_required"),
		buildColumnName(validPreparationsTable, "time_estimate_required"),
		buildColumnName(validPreparationsTable, "condition_expression_required"),
		buildColumnName(validPreparationsTable, "consumes_vessel"),
		buildColumnName(validPreparationsTable, "only_for_vessels"),
		buildColumnName(validPreparationsTable, "minimum_vessel_count"),
		buildColumnName(validPreparationsTable, "maximum_vessel_count"),
		buildColumnName(validPreparationsTable, "slug"),
		buildColumnName(validPreparationsTable, "past_tense"),
		buildColumnName(validPreparationsTable, createdAt),
		buildColumnName(validPreparationsTable, lastUpdatedAt),
		buildColumnName(validPreparationsTable, archivedAt),
	}
)

func buildSelectValidPreparationsNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(validPreparationsTable)
}

func buildSelectValidPreparationQuery(_ squirrel.StatementBuilderType) string {
	return buildSelectQuery(validPreparationsTable, validPreparationsTableColumns)
}
