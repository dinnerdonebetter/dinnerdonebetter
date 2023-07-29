package main

import (
	"github.com/Masterminds/squirrel"
)

const (
	validIngredientsTable = "valid_ingredients"
)

var (
	validIngredientsTableColumns = []string{
		buildColumnName(validIngredientsTable, id),
		buildColumnName(validIngredientsTable, "name"),
		buildColumnName(validIngredientsTable, "description"),
		buildColumnName(validIngredientsTable, "warning"),
		buildColumnName(validIngredientsTable, "contains_egg"),
		buildColumnName(validIngredientsTable, "contains_dairy"),
		buildColumnName(validIngredientsTable, "contains_peanut"),
		buildColumnName(validIngredientsTable, "contains_tree_nut"),
		buildColumnName(validIngredientsTable, "contains_soy"),
		buildColumnName(validIngredientsTable, "contains_wheat"),
		buildColumnName(validIngredientsTable, "contains_shellfish"),
		buildColumnName(validIngredientsTable, "contains_sesame"),
		buildColumnName(validIngredientsTable, "contains_fish"),
		buildColumnName(validIngredientsTable, "contains_gluten"),
		buildColumnName(validIngredientsTable, "animal_flesh"),
		buildColumnName(validIngredientsTable, "volumetric"),
		buildColumnName(validIngredientsTable, "is_liquid"),
		buildColumnName(validIngredientsTable, "icon_path"),
		buildColumnName(validIngredientsTable, "animal_derived"),
		buildColumnName(validIngredientsTable, "plural_name"),
		buildColumnName(validIngredientsTable, "restrict_to_preparations"),
		buildColumnName(validIngredientsTable, "minimum_ideal_storage_temperature_in_celsius"),
		buildColumnName(validIngredientsTable, "maximum_ideal_storage_temperature_in_celsius"),
		buildColumnName(validIngredientsTable, "storage_instructions"),
		buildColumnName(validIngredientsTable, "slug"),
		buildColumnName(validIngredientsTable, "contains_alcohol"),
		buildColumnName(validIngredientsTable, "shopping_suggestions"),
		buildColumnName(validIngredientsTable, "is_starch"),
		buildColumnName(validIngredientsTable, "is_protein"),
		buildColumnName(validIngredientsTable, "is_grain"),
		buildColumnName(validIngredientsTable, "is_fruit"),
		buildColumnName(validIngredientsTable, "is_salt"),
		buildColumnName(validIngredientsTable, "is_fat"),
		buildColumnName(validIngredientsTable, "is_acid"),
		buildColumnName(validIngredientsTable, "is_heat"),
		buildColumnName(validIngredientsTable, createdAt),
		buildColumnName(validIngredientsTable, lastUpdatedAt),
		buildColumnName(validIngredientsTable, archivedAt),
	}
)

func buildSelectValidIngredientsNeedingIndexingQuery(_ squirrel.StatementBuilderType) string {
	return buildListOfNeedingIndexingQuery(validIngredientsTable)
}

func buildSelectValidIngredientsQuery(_ squirrel.StatementBuilderType) string {
	return buildListQuery(validIngredientsTable, validIngredientsTableColumns)
}

func buildSelectValidIngredientQuery(_ squirrel.StatementBuilderType) string {
	return buildSelectQuery(validIngredientsTable, validIngredientsTableColumns)
}
