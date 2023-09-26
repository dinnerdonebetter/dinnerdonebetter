package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validMeasurementUnitConversionsTableName = "valid_measurement_unit_conversions"
)

var validMeasurementUnitConversionsColumns = []string{
	idColumn,
	"from_unit",
	"to_unit",
	"only_for_ingredient",
	"modifier",
	notesColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidMeasurementUnitConversionsQueries() []*Query {
	insertColumns := filterForInsert(validMeasurementUnitConversionsColumns)

	fullSelectColumns := mergeColumns(
		applyToEach(validMeasurementUnitConversionsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_measurement_unit_conversion_%s", validMeasurementUnitConversionsTableName, s, s)
		}),
		append(
			append(
				applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("valid_measurement_units_from.%s as from_unit_%s", s, s)
				}),
				applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
					return fmt.Sprintf("valid_measurement_units_to.%s as to_unit_%s", s, s)
				})...,
			),
			applyToEach(validIngredientsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
			})...,
		),
		1,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidMeasurementUnitConversion",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validMeasurementUnitConversionsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidMeasurementUnitConversion",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validMeasurementUnitConversionsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidMeasurementUnitConversionExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitConversionsTableName,
				validMeasurementUnitConversionsTableName,
				archivedAtColumn,
				validMeasurementUnitConversionsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllValidMeasurementUnitConversionsFromMeasurementUnit",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN valid_measurement_units AS valid_measurement_units_from ON valid_measurement_unit_conversions.from_unit = valid_measurement_units_from.id
    JOIN valid_measurement_units AS valid_measurement_units_to ON valid_measurement_unit_conversions.to_unit = valid_measurement_units_to.id
    LEFT JOIN valid_ingredients ON valid_measurement_unit_conversions.only_for_ingredient = valid_ingredients.id
WHERE
    valid_measurement_units_from.id = sqlc.arg(id)
    AND valid_measurement_unit_conversions.archived_at IS NULL
    AND valid_measurement_units_from.archived_at IS NULL
    AND valid_measurement_units_to.archived_at IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAllValidMeasurementUnitConversionsToMeasurementUnit",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN valid_measurement_units AS valid_measurement_units_from ON valid_measurement_unit_conversions.from_unit = valid_measurement_units_from.id
    JOIN valid_measurement_units AS valid_measurement_units_to ON valid_measurement_unit_conversions.to_unit = valid_measurement_units_to.id
    LEFT JOIN valid_ingredients ON valid_measurement_unit_conversions.only_for_ingredient = valid_ingredients.id
WHERE
    valid_measurement_units_to.id = sqlc.arg(id)
    AND valid_measurement_unit_conversions.archived_at IS NULL
    AND valid_measurement_units_from.archived_at IS NULL
    AND valid_measurement_units_to.archived_at IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidMeasurementUnitConversion",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN valid_measurement_units AS valid_measurement_units_from ON valid_measurement_unit_conversions.from_unit = valid_measurement_units_from.id
    JOIN valid_measurement_units AS valid_measurement_units_to ON valid_measurement_unit_conversions.to_unit = valid_measurement_units_to.id
    LEFT JOIN valid_ingredients ON valid_measurement_unit_conversions.only_for_ingredient = valid_ingredients.id
WHERE
    valid_measurement_unit_conversions.id = sqlc.arg(id)
    AND valid_measurement_unit_conversions.archived_at IS NULL
    AND valid_measurement_units_from.archived_at IS NULL
    AND valid_measurement_units_to.archived_at IS NULL;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validMeasurementUnitConversionsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidMeasurementUnitConversion",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validMeasurementUnitConversionsTableName,
				strings.Join(applyToEach(filterForUpdate(validMeasurementUnitConversionsColumns), func(i int, s string) string {
					return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
				}), ",\n\t"),
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
