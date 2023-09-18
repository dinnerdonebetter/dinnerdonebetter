package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validVesselsTableName = "valid_vessels"

var validVesselsColumns = []string{
	idColumn,
	"name",
	"plural_name",
	"description",
	"icon_path",
	"usable_for_storage",
	"slug",
	"display_in_summary_lists",
	"include_in_generated_instructions",
	"capacity",
	"capacity_unit",
	"width_in_millimeters",
	"length_in_millimeters",
	"height_in_millimeters",
	"shape",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidVesselsQueries() []*Query {
	insertColumns := filterForInsert(validVesselsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validVesselsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validVesselsTableName,
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
				Name: "CheckValidVesselExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validVesselsTableName,
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				validVesselsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVesselByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL;`,
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("valid_ingredients.%s", s)
				}), ",\n\t"),
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVessels",
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
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("valid_ingredients.%s", s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validVesselsTableName,
					true,
				),
				buildTotalCountSelect(
					validVesselsTableName,
				),
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				buildFilterConditions(
					validVesselsTableName,
					true,
				),
				validVesselsTableName,
				idColumn,
				validVesselsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVesselsNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
    AND (
    %s.%s IS NULL
    OR %s.%s < NOW() - '24 hours'::INTERVAL
);`,
				validVesselsTableName,
				idColumn,
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				validVesselsTableName,
				lastIndexedAtColumn,
				validVesselsTableName,
				lastIndexedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
				}), ",\n\t"),
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				validVesselsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRandomValidVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
ORDER BY RANDOM() LIMIT 1;`,
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
				}), ",\n\t"),
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidVesselsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
				}), ",\n\t"),
				validVesselsTableName,
				validVesselsTableName,
				archivedAtColumn,
				validVesselsTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.name %s
	AND %s.%s IS NULL
LIMIT 50;`,
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
				}), ",\n\t"),
				validVesselsTableName,
				validVesselsTableName,
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
				validVesselsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchValidVesselsByPreparationAndIngredientName",
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
				strings.Join(applyToEach(validVesselsColumns, func(i int, s string) string {
					if i == 0 {
						return fmt.Sprintf("DISTINCT(%s.%s)", validVesselsTableName, s)
					}
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
				}), ",\n\t"),
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validVesselsTableName,
				strings.Join(applyToEach(filterForUpdate(validVesselsColumns), func(i int, s string) string {
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
				Name: "UpdateValidVesselLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validVesselsTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
