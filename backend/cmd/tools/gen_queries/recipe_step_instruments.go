package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeStepInstrumentsTableName = "recipe_step_instruments"

	recipeStepInstrumentIDColumn = "recipe_step_instrument_id"
	instrumentIDColumn           = "instrument_id"
)

var recipeStepInstrumentsColumns = []string{
	idColumn,
	instrumentIDColumn,
	"recipe_step_product_id",
	nameColumn,
	notesColumn,
	"preference_rank",
	"optional",
	"minimum_quantity",
	"maximum_quantity",
	"option_index",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToRecipeStepColumn,
}

func buildRecipeStepInstrumentsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipeStepInstrumentsColumns)

		fullSelectColumn := mergeColumns(
			applyToEach(filterFromSlice(recipeStepInstrumentsColumns, instrumentIDColumn, measurementUnitColumn), func(i int, s string) string {
				return fmt.Sprintf("%s.%s", recipeStepInstrumentsTableName, s)
			}),
			applyToEach(validInstrumentsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_instrument_%s", validInstrumentsTableName, s, s)
			}),
			1,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipeStepInstrument",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
					recipeStepInstrumentsTableName,
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
					Name: "CreateRecipeStepInstrument",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipeStepInstrumentsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckRecipeStepInstrumentExistence",
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
					recipeStepInstrumentsTableName, idColumn,
					recipeStepInstrumentsTableName,
					recipeStepsTableName, recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipeStepInstrumentsTableName, archivedAtColumn,
					recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepInstrumentsTableName, idColumn, recipeStepInstrumentIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepInstrumentsForRecipe",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	LEFT JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumn, ",\n\t"),
					recipeStepInstrumentsTableName,
					validInstrumentsTableName, recipeStepInstrumentsTableName, instrumentIDColumn, validInstrumentsTableName, idColumn,
					recipeStepsTableName, recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipeStepInstrumentsTableName, archivedAtColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepInstrument",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	LEFT JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumn, ",\n\t"),
					recipeStepInstrumentsTableName,
					validInstrumentsTableName, recipeStepInstrumentsTableName, instrumentIDColumn, validInstrumentsTableName, idColumn,
					recipeStepsTableName, recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipeStepInstrumentsTableName, archivedAtColumn,
					recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepInstrumentsTableName, idColumn, recipeStepInstrumentIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeStepInstruments",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT COUNT(%s.%s)
		FROM %s
		WHERE
			%s.%s IS NULL
			AND %s.%s = sqlc.arg(%s)
			%s
	) AS filtered_count,
	%s
FROM %s
	LEFT JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
	JOIN %s ON %s.%s=%s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumn, ",\n\t"),
					recipeStepInstrumentsTableName, idColumn,
					recipeStepInstrumentsTableName,
					recipeStepInstrumentsTableName, archivedAtColumn,
					recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					strings.Join(strings.Split(buildFilterConditions(recipeStepInstrumentsTableName, true, false, nil), "\n"), "\n\t\t"),
					buildTotalCountSelect(recipeStepInstrumentsTableName, true, []string{}),
					recipeStepInstrumentsTableName,
					validInstrumentsTableName, recipeStepInstrumentsTableName, instrumentIDColumn, validInstrumentsTableName, idColumn,
					recipeStepsTableName, recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepsTableName, idColumn,
					recipesTableName, recipeStepsTableName, belongsToRecipeColumn, recipesTableName, idColumn,
					recipeStepInstrumentsTableName, archivedAtColumn,
					recipeStepInstrumentsTableName, belongsToRecipeStepColumn, recipeStepIDColumn,
					recipeStepsTableName, archivedAtColumn,
					recipeStepsTableName, belongsToRecipeColumn, recipeIDColumn,
					recipeStepsTableName, idColumn, recipeStepIDColumn,
					recipesTableName, archivedAtColumn,
					recipesTableName, idColumn, recipeIDColumn,
					buildFilterConditions(recipeStepInstrumentsTableName, true, false, nil),
					offsetLimitAddendum,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipeStepInstrument",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					recipeStepInstrumentsTableName,
					strings.Join(applyToEach(filterForUpdate(recipeStepInstrumentsColumns, belongsToRecipeStepColumn), func(i int, s string) string {
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
