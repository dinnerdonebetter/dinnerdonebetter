package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validIngredientsTableName = "valid_ingredients"

var validIngredientsColumns = []string{
	idColumn,
	"name",
	"description",
	"warning",
	"contains_egg",
	"contains_dairy",
	"contains_peanut",
	"contains_tree_nut",
	"contains_soy",
	"contains_wheat",
	"contains_shellfish",
	"contains_sesame",
	"contains_fish",
	"contains_gluten",
	"animal_flesh",
	"volumetric",
	"is_liquid",
	"icon_path",
	"animal_derived",
	"plural_name",
	"restrict_to_preparations",
	"minimum_ideal_storage_temperature_in_celsius",
	"maximum_ideal_storage_temperature_in_celsius",
	"storage_instructions",
	"slug",
	"contains_alcohol",
	"shopping_suggestions",
	"is_starch",
	"is_protein",
	"is_grain",
	"is_fruit",
	"is_salt",
	"is_fat",
	"is_acid",
	"is_heat",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientsQueries() []*Query {
	insertColumns := filterForInsert(validIngredientsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validIngredientsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidIngredient",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validIngredientsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					switch s {
					case "minimum_ideal_storage_temperature_in_celsius",
						"maximum_ideal_storage_temperature_in_celsius":
						return fmt.Sprintf("sqlc.narg(%s)", s)
					default:
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidIngredientExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validIngredientsTableName,
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
				validIngredientsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredients",
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
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validIngredientsTableName,
					true,
				),
				buildTotalCountSelect(
					validIngredientsTableName,
				),
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
				buildFilterConditions(
					validIngredientsTableName,
					true,
				),
				validIngredientsTableName,
				idColumn,
				validIngredientsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientsNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
    AND (
    %s.%s IS NULL
    OR %s.%s < NOW() - '24 hours'::INTERVAL
);`,
				validIngredientsTableName,
				idColumn,
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
				validIngredientsTableName,
				lastIndexedAtColumn,
				validIngredientsTableName,
				lastIndexedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredient",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
				validIngredientsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRandomValidIngredient",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
ORDER BY RANDOM() LIMIT 1;`,
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidIngredientsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				validIngredientsTableName,
				validIngredientsTableName,
				archivedAtColumn,
				validIngredientsTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidIngredients",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.name %s
	AND %s.%s IS NULL
LIMIT 50;`,
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				validIngredientsTableName,
				validIngredientsTableName,
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
				validIngredientsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchValidIngredientsByPreparationAndIngredientName",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND (
		valid_ingredient_preparations.valid_preparation_id = sqlc.arg(valid_preparation_id)
		OR valid_preparations.restrict_to_ingredients IS FALSE
	)
	AND valid_ingredients.name %s;`,
				strings.Join(applyToEach(validIngredientsColumns, func(i int, s string) string {
					if i == 0 {
						return fmt.Sprintf("DISTINCT(%s.%s)", validIngredientsTableName, s)
					}
					return fmt.Sprintf("%s.%s", validIngredientsTableName, s)
				}), ",\n\t"),
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredient",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validIngredientsTableName,
				strings.Join(applyToEach(filterForUpdate(validIngredientsColumns), func(i int, s string) string {
					switch s {
					case "minimum_ideal_storage_temperature_in_celsius",
						"maximum_ideal_storage_temperature_in_celsius":
						return fmt.Sprintf("%s = sqlc.narg(%s)", s, s)
					default:
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}
				}), ",\n\t"),
				lastUpdatedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidIngredientLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validIngredientsTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
