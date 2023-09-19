package main

import (
	"fmt"
	"github.com/cristalhq/builq"
	"strings"
)

const validIngredientPreparationsTableName = "valid_ingredient_preparations"

var validIngredientPreparationsColumns = []string{
	idColumn,
	"notes",
	"valid_preparation_id",
	"valid_ingredient_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientPreparationsQueries() []*Query {
	insertColumns := filterForInsert(validIngredientPreparationsColumns)

	fullSelectColumns := mergeColumns(
		mergeColumns(
			applyToEach(filterFromSlice(validIngredientPreparationsColumns, "valid_preparation_id", "valid_ingredient_id"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_preparation_%s", validIngredientPreparationsTableName, s, s)
			}),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			}),
			2,
		),
		applyToEach(validPreparationsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
		}),
		2,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientPreparation",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientPreparationsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientPreparation",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validIngredientPreparationsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientPreparationExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientPreparationsTableName,
				validIngredientPreparationsTableName,
				validIngredientPreparationsTableName,
				archivedAtColumn,
				validIngredientPreparationsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientPreparationsForIngredient",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
	valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredient_preparations.valid_ingredient_id = sqlc.arg(ids)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validInstrumentsTableName,
					true,
				),
				buildTotalCountSelect(
					validInstrumentsTableName,
				),
				buildFilterConditions(validInstrumentsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientPreparationsForPreparation",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
	valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredient_preparations.valid_preparation_id = sqlc.arg(ids)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validInstrumentsTableName,
					true,
				),
				buildTotalCountSelect(
					validInstrumentsTableName,
				),
				buildFilterConditions(validInstrumentsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientPreparations",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
	valid_ingredient_preparations.archived_at IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validInstrumentsTableName,
					true,
				),
				buildTotalCountSelect(
					validInstrumentsTableName,
				),
				buildFilterConditions(validInstrumentsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientPreparation",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
	valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_ingredient_preparations.id = sqlc.arg(id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidIngredientPreparationPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT id
	FROM %s
	WHERE valid_ingredient_id = sqlc.arg(valid_ingredient_id)
	AND valid_preparation_id = sqlc.arg(valid_preparation_id)
	AND %s IS NULL
);`,
				validIngredientPreparationsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchValidIngredientPreparationsByPreparationAndIngredientName",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
	valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_preparations.id = sqlc.arg(id)
	AND valid_ingredients.name %s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientPreparation",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validIngredientPreparationsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientPreparationsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
