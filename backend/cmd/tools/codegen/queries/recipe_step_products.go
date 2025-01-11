package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepProductsTableName = "recipe_step_products"

	recipeStepProductIDColumn = "recipe_step_product_id"
)

func init() {
	registerTableName(recipeStepProductsTableName)
}

var recipeStepProductsColumns = []string{
	idColumn,
	nameColumn,
	"type",
	measurementUnitColumn,
	"minimum_quantity_value",
	"maximum_quantity_value",
	"quantity_notes",
	"compostable",
	"maximum_storage_duration_in_seconds",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	"storage_instructions",
	"is_liquid",
	"is_waste",
	"index",
	"contained_in_vessel_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToRecipeStepColumn,
}

func buildRecipeStepProductsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipeStepProductsColumns)

		fullSelectColumns := mergeColumns(
			applyToEach(filterFromSlice(recipeStepProductsColumns, measurementUnitColumn), func(i int, s string) string {
				return fmt.Sprintf("%s.%s", recipeStepProductsTableName, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			3,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipeStepProduct",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
					recipeStepProductsTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					belongsToRecipeStepColumn,
					belongsToRecipeStepColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateRecipeStepProduct",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipeStepProductsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckRecipeStepProductExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
		JOIN %s ON %s.%s=%s.%s
		JOIN %s ON %s.%s=%s.%s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					recipeStepProductsTableName, idColumn,
					recipeStepProductsTableName,
					recipeStepsTableName, recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipeStepProductsTableName, archivedAtColumn,
					recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepProductsTableName, idColumn, recipeStepProductIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepProductsForRecipe",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	LEFT JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					recipeStepProductsTableName,
					recipeStepsTableName, recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					validMeasurementUnitsTableName, recipeStepProductsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,
					recipeStepProductsTableName, archivedAtColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepProducts",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	(
		SELECT COUNT(%s.%s)
		FROM %s
		WHERE %s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
	) AS total_count
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	LEFT JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(recipeStepProductsTableName, true, true, nil),
					recipeStepProductsTableName, idColumn,
					recipeStepProductsTableName,
					recipeStepProductsTableName, archivedAtColumn,
					recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepProductsTableName,
					recipeStepsTableName, recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					validMeasurementUnitsTableName, recipeStepProductsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,
					recipeStepProductsTableName, archivedAtColumn,
					recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
					buildFilterConditions(recipeStepProductsTableName, true, true),
					offsetLimitAddendum,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepProduct",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	LEFT JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					recipeStepProductsTableName,
					recipeStepsTableName, recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					validMeasurementUnitsTableName, recipeStepProductsTableName, measurementUnitColumn, validMeasurementUnitsTableName, idColumn,
					recipeStepProductsTableName, archivedAtColumn,
					recipeStepProductsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepProductsTableName, idColumn, recipeStepProductIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipeStepProduct",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					recipeStepProductsTableName,
					strings.Join(applyToEach(filterForUpdate(recipeStepProductsColumns, belongsToRecipeStepColumn), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					belongsToRecipeStepColumn, belongsToRecipeStepColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
