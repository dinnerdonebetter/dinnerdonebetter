package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validPreparationVesselsTableName = "valid_preparation_vessels"

var validPreparationVesselsColumns = []string{
	idColumn,
	"notes",
	"valid_preparation_id",
	"valid_vessel_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationVesselsQueries() []*Query {
	insertColumns := filterForInsert(validPreparationVesselsColumns)

	fullSelectColumns := append(
		mergeColumns(
			applyToEach(filterFromSlice(validPreparationVesselsColumns, "valid_preparation_id", "valid_vessel_id"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_vessel_%s", validPreparationVesselsTableName, s, s)
			}),
			applyToEach(validPreparationsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
			}),
			2,
		),
		mergeColumns(
			applyToEach(filterFromSlice(validVesselsColumns, "capacity_unit"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_vessel_%s", validVesselsTableName, s, s)
			}),
			applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
			}),
			10,
		)...,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidPreparationVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validPreparationVesselsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidPreparationVessel",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				validPreparationVesselsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidPreparationVesselExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
				validPreparationVesselsTableName,
				validPreparationVesselsTableName,
				validPreparationVesselsTableName,
				archivedAtColumn,
				validPreparationVesselsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVesselsForPreparation",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
	JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
WHERE
	valid_preparation_vessels.archived_at IS NULL
	AND valid_vessels.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_preparation_vessels.valid_preparation_id = sqlc.arg(id)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validPreparationVesselsTableName,
					true,
				),
				buildTotalCountSelect(
					validPreparationVesselsTableName,
				),
				validPreparationVesselsTableName,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVesselsForVessel",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
	JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
WHERE
	valid_preparation_vessels.archived_at IS NULL
	AND valid_vessels.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_preparation_vessels.valid_vessel_id = sqlc.arg(id)
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validPreparationVesselsTableName,
					true,
				),
				buildTotalCountSelect(
					validPreparationVesselsTableName,
				),
				validPreparationVesselsTableName,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVessels",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
	JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
WHERE
	valid_preparation_vessels.archived_at IS NULL
	AND valid_vessels.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	%s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					validPreparationVesselsTableName,
					true,
				),
				buildTotalCountSelect(
					validPreparationVesselsTableName,
				),
				validPreparationVesselsTableName,
				buildFilterConditions(validPreparationVesselsTableName, true),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationVessel",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
	JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
	LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
WHERE
	valid_preparation_vessels.archived_at IS NULL
	AND valid_vessels.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_preparation_vessels.id = sqlc.arg(id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
				validPreparationVesselsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidPreparationVesselPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s
	FROM %s
	WHERE valid_vessel_id = sqlc.arg(valid_vessel_id)
	AND valid_preparation_id = sqlc.arg(valid_preparation_id)
	AND %s IS NULL
);`,
				idColumn,
				validPreparationVesselsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidPreparationVessel",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validPreparationVesselsTableName,
				strings.Join(applyToEach(filterForUpdate(validPreparationVesselsColumns), func(i int, s string) string {
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
