package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	validIngredientStatesTable = "valid_ingredient_states"
)

var (
	// validIngredientStatesTableColumns are the columns for the valid_ingredient_states table.
	validIngredientStatesTableColumns = []string{
		buildColumnName(validIngredientStatesTable, id),
		buildColumnName(validIngredientsTable, "name"),
		buildColumnName(validIngredientsTable, "description"),
		buildColumnName(validIngredientsTable, "icon_path"),
		buildColumnName(validIngredientsTable, "slug"),
		buildColumnName(validIngredientsTable, "past_tense"),
		buildColumnName(validIngredientsTable, "attribute_type"),
		buildColumnName(validIngredientStatesTable, createdAt),
		buildColumnName(validIngredientStatesTable, lastUpdatedAt),
		buildColumnName(validIngredientStatesTable, archivedAt),
	}
)

func buildSelectValidIngredientStatesNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(validIngredientStatesTable)
}

func buildSelectValidIngredientStateQuery(_ squirrel.StatementBuilderType) string {
	return buildSelectQuery(validIngredientStatesTable, validIngredientStatesTableColumns)
}
