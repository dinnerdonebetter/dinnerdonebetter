package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validMeasurementUnitsTableName = "valid_measurement_units"

var validMeasurementUnitsColumns = []string{
	idColumn,
	"name",
	"description",
	"volumetric",
	"icon_path",
	"universal",
	"metric",
	"imperial",
	"slug",
	"plural_name",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidMeasurementUnitsQueries() []*Query {
	insertColumns := filterForInsert(validMeasurementUnitsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidMeasurementUnit",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidMeasurementUnit",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validMeasurementUnitsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidMeasurementUnitExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				validMeasurementUnitsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnitByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				validMeasurementUnitsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnits",
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
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validMeasurementUnitsTableName,
					true,
				),
				buildTotalCountSelect(
					validMeasurementUnitsTableName,
				),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				buildFilterConditions(
					validMeasurementUnitsTableName,
					true,
				),
				validMeasurementUnitsTableName,
				idColumn,
				validMeasurementUnitsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnitsNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
    AND (
    %s.%s IS NULL
    OR %s.%s < NOW() - '24 hours'::INTERVAL
);`,
				validMeasurementUnitsTableName,
				idColumn,
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				validMeasurementUnitsTableName,
				lastIndexedAtColumn,
				validMeasurementUnitsTableName,
				lastIndexedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnit",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				validMeasurementUnitsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRandomValidMeasurementUnit",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
ORDER BY RANDOM() LIMIT 1;`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnitsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				archivedAtColumn,
				validMeasurementUnitsTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidMeasurementUnits",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.name %s
	AND %s.%s IS NULL
LIMIT 50;`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				validMeasurementUnitsTableName,
				validMeasurementUnitsTableName,
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
				validMeasurementUnitsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchValidMeasurementUnitsByIngredientID",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
    %s,
    %s
FROM valid_measurement_units
	JOIN valid_ingredient_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE
	(
	    valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
	    OR valid_measurement_units.universal = true
	)
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
	AND valid_ingredient_measurement_units.archived_at IS NULL
	%s
%s;`,
				strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					if i == 0 {
						return fmt.Sprintf("DISTINCT(%s.%s)", validMeasurementUnitsTableName, s)
					}
					return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(validMeasurementUnitsTableName, true, ` (
				valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
				OR valid_measurement_units.universal = true
			)`),
				buildTotalCountSelect(validMeasurementUnitsTableName),
				buildFilterConditions(validMeasurementUnitsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidMeasurementUnit",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validMeasurementUnitsTableName,
				strings.Join(applyToEach(filterForUpdate(validMeasurementUnitsColumns), func(i int, s string) string {
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
				Name: "UpdateValidMeasurementUnitLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validMeasurementUnitsTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
