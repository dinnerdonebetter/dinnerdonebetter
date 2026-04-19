package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeRatingsTableName = "recipe_ratings"
)

func init() {
	registerTableName(recipeRatingsTableName)
}

var recipeRatingsColumns = []string{
	idColumn,
	belongsToRecipeColumn,
	"taste",
	"difficulty",
	"cleanup",
	"instructions",
	"overall",
	notesColumn,
	createdByUserColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildRecipeRatingsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipeRatingsColumns)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipeRating",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					recipeRatingsTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateRecipeRating",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipeRatingsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckRecipeRatingExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					recipeRatingsTableName, idColumn,
					recipeRatingsTableName,
					recipeRatingsTableName, archivedAtColumn,
					recipeRatingsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeRatingsForRecipe",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE
	%s.%s IS NULL AND
	%s.%s = sqlc.arg(%s)
	%s
GROUP BY %s.%s
%s;`,
					strings.Join(applyToEach(recipeRatingsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", recipeRatingsTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(recipeRatingsTableName, true, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, belongsToRecipeColumn, belongsToRecipeColumn)),
					buildTotalCountSelect(recipeRatingsTableName, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, belongsToRecipeColumn, belongsToRecipeColumn)),
					recipeRatingsTableName,
					recipeRatingsTableName, archivedAtColumn,
					recipeRatingsTableName, belongsToRecipeColumn, belongsToRecipeColumn,
					buildFilterConditions(recipeRatingsTableName, true, true, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, belongsToRecipeColumn, belongsToRecipeColumn)),
					recipeRatingsTableName, idColumn,
					buildCursorLimitClause(recipeRatingsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeRatingsForUser",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE
	%s.%s IS NULL AND
	%s.%s = sqlc.arg(%s)
	%s
GROUP BY %s.%s
%s;`,
					strings.Join(applyToEach(recipeRatingsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", recipeRatingsTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(recipeRatingsTableName, true, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, createdByUserColumn, createdByUserColumn)),
					buildTotalCountSelect(recipeRatingsTableName, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, createdByUserColumn, createdByUserColumn)),
					recipeRatingsTableName,
					recipeRatingsTableName, archivedAtColumn,
					recipeRatingsTableName, createdByUserColumn, createdByUserColumn,
					buildFilterConditions(recipeRatingsTableName, true, true, fmt.Sprintf("%s.%s = sqlc.arg(%s)", recipeRatingsTableName, createdByUserColumn, createdByUserColumn)),
					recipeRatingsTableName, idColumn,
					buildCursorLimitClause(recipeRatingsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeRating",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(applyToEach(recipeRatingsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", recipeRatingsTableName, s)
					}), ",\n\t"),
					recipeRatingsTableName,
					recipeRatingsTableName, archivedAtColumn,
					recipeRatingsTableName, idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipeRating",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					recipeRatingsTableName,
					strings.Join(applyToEach(filterForUpdate(recipeRatingsColumns, createdByUserColumn), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
