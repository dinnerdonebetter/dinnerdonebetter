package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validIngredientStatesTableName = "valid_ingredient_states"

var validIngredientStatesColumns = []string{
	idColumn,
	"name",
	"past_tense",
	"slug",
	"description",
	"icon_path",
	"attribute_type",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientStatesQueries() []*Query {

	insertColumns := filterForInsert(validIngredientStatesColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredientState",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientStatesTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredientState",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validIngredientStatesTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientStateExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				validIngredientStatesTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStateByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validIngredientStatesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientStatesTableName, s)
				}), ",\n\t"),
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				validIngredientStatesTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStates",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
    %s,
    %s
FROM %s
WHERE
	%s.%s IS NULL
	%s
GROUP BY %s.%s
ORDER BY %s.%s
%s;`,
				strings.Join(applyToEach(validIngredientStatesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientStatesTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validIngredientStatesTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientStatesTableName,
				),
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				buildFilterConditions(
					validIngredientStatesTableName,
					true,
				),
				validIngredientStatesTableName,
				idColumn,
				validIngredientStatesTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStatesNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
    AND (
    %s.%s IS NULL
    OR %s.%s < NOW() - '24 hours'::INTERVAL
);`,
				validIngredientStatesTableName,
				idColumn,
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				validIngredientStatesTableName,
				lastIndexedAtColumn,
				validIngredientStatesTableName,
				lastIndexedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientState",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validIngredientStatesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientStatesTableName, s)
				}), ",\n\t"),
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				validIngredientStatesTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientStatesWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validIngredientStatesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientStatesTableName, s)
				}), ",\n\t"),
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				archivedAtColumn,
				validIngredientStatesTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidIngredientStates",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.name %s
	AND %s.%s IS NULL
LIMIT 50;`,
				strings.Join(applyToEach(validIngredientStatesColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientStatesTableName, s)
				}), ",\n\t"),
				validIngredientStatesTableName,
				validIngredientStatesTableName,
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
				validIngredientStatesTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientState",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validIngredientStatesTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientStatesColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientStateLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validIngredientStatesTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
