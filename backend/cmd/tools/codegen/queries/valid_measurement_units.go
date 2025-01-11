package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validMeasurementUnitsTableName = "valid_measurement_units"

	validMeasurementUnitsUniversalColumn = "universal"
)

func init() {
	registerTableName(validMeasurementUnitsTableName)
}

var validMeasurementUnitsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	"volumetric",
	iconPathColumn,
	validMeasurementUnitsUniversalColumn,
	"metric",
	"imperial",
	slugColumn,
	pluralNameColumn,
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidMeasurementUnitsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(validMeasurementUnitsColumns)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveValidMeasurementUnit",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					validMeasurementUnitsTableName,
					archivedAtColumn,
					currentTimeExpression,
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
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					validMeasurementUnitsTableName, idColumn,
					validMeasurementUnitsTableName,
					validMeasurementUnitsTableName, archivedAtColumn,
					validMeasurementUnitsTableName, idColumn, idColumn,
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
					buildFilterCountSelect(validMeasurementUnitsTableName, true, true, []string{}),
					buildTotalCountSelect(validMeasurementUnitsTableName, true, []string{}),
					validMeasurementUnitsTableName,
					validMeasurementUnitsTableName,
					archivedAtColumn,
					buildFilterConditions(
						validMeasurementUnitsTableName,
						true,
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
	OR %s.%s < %s - '24 hours'::INTERVAL
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
					currentTimeExpression,
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
WHERE %s.%s %s
	AND %s.%s IS NULL
LIMIT 50;`,
					strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
					}), ",\n\t"),
					validMeasurementUnitsTableName,
					validMeasurementUnitsTableName, nameColumn, buildILIKEForArgument("name_query"),
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
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	(
		%s.%s = sqlc.arg(%s)
		OR %s.%s = TRUE
	)
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	%s
%s;`,
					strings.Join(applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
						if i == 0 {
							return fmt.Sprintf("DISTINCT(%s.%s)", validMeasurementUnitsTableName, s)
						}
						return fmt.Sprintf("%s.%s", validMeasurementUnitsTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(validMeasurementUnitsTableName, true, true, []string{}, ` (
				valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
				OR valid_measurement_units.universal = true
			)`),
					buildTotalCountSelect(validMeasurementUnitsTableName, true, []string{}),
					validMeasurementUnitsTableName,
					validIngredientMeasurementUnitsTableName, validIngredientMeasurementUnitsTableName, validMeasurementUnitIDColumn, validMeasurementUnitsTableName, idColumn,
					validIngredientsTableName, validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientsTableName, idColumn,
					validIngredientMeasurementUnitsTableName, validIngredientIDColumn, validIngredientIDColumn,
					validMeasurementUnitsTableName, validMeasurementUnitsUniversalColumn,
					validMeasurementUnitsTableName, archivedAtColumn,
					validIngredientsTableName, archivedAtColumn,
					validIngredientMeasurementUnitsTableName, archivedAtColumn,
					buildFilterConditions(validMeasurementUnitsTableName, true, false),
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
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					validMeasurementUnitsTableName,
					strings.Join(applyToEach(filterForUpdate(validMeasurementUnitsColumns), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn,
					currentTimeExpression,
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
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
					validMeasurementUnitsTableName,
					lastIndexedAtColumn,
					currentTimeExpression,
					idColumn,
					idColumn,
					archivedAtColumn,
				)),
			},
		}
	default:
		return nil
	}
}
