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

	fullSelectColumns := mergeColumns(
		applyToEach(filterFromSlice(validVesselsColumns, "capacity_unit"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", validVesselsTableName, s)
		}),
		applyToEach(validMeasurementUnitsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as valid_measurement_unit_%s", validMeasurementUnitsTableName, s, s)
		}),
		10,
	)

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
					return fmt.Sprintf("sqlc.arg(%s)", s)
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
					return fmt.Sprintf("%s.%s", validVesselsTableName, s)
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
				Name: "GetValidVesselIDsNeedingIndexing",
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
	JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE %s.%s IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(fullSelectColumns, ",\n\t"),
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
	JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE %s.%s IS NULL
	AND valid_measurement_units.archived_at IS NULL
ORDER BY RANDOM() LIMIT 1;`,
				strings.Join(fullSelectColumns, ",\n\t"),
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
	JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE %s.%s IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(fullSelectColumns, ",\n\t"),
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
