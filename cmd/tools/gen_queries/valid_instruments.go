package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const validInstrumentsTableName = "valid_instruments"

var validInstrumentsColumns = []string{
	idColumn,
	"name",
	"description",
	"icon_path",
	"plural_name",
	"usable_for_storage",
	"slug",
	"display_in_summary_lists",
	"include_in_generated_instructions",
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidInstrumentsQueries() []*Query {
	insertColumns := filterForInsert(validInstrumentsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveValidInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validInstrumentsTableName,
				archivedAtColumn,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidInstrument",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validInstrumentsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidInstrumentExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validInstrumentsTableName,
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				validInstrumentsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidInstrumentByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				validInstrumentsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidInstruments",
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
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					validInstrumentsTableName,
					true,
				),
				buildTotalCountSelect(
					validInstrumentsTableName,
				),
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				buildFilterConditions(
					validInstrumentsTableName,
					true,
				),
				validInstrumentsTableName,
				idColumn,
				validInstrumentsTableName,
				idColumn,
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidInstrumentsNeedingIndexing",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT %s.%s
FROM %s
WHERE %s.%s IS NULL
    AND (
    %s.%s IS NULL
    OR %s.%s < NOW() - '24 hours'::INTERVAL
);`,
				validInstrumentsTableName,
				idColumn,
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				validInstrumentsTableName,
				lastIndexedAtColumn,
				validInstrumentsTableName,
				lastIndexedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidInstrument",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
AND %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				validInstrumentsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetRandomValidInstrument",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
ORDER BY RANDOM() LIMIT 1;`,
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidInstrumentsWithIDs",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				validInstrumentsTableName,
				validInstrumentsTableName,
				archivedAtColumn,
				validInstrumentsTableName,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForValidInstruments",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.name %s
	AND %s.%s IS NULL
LIMIT 50;`,
				strings.Join(applyToEach(validInstrumentsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", validInstrumentsTableName, s)
				}), ",\n\t"),
				validInstrumentsTableName,
				validInstrumentsTableName,
				"ILIKE '%' || sqlc.arg(name_query)::text || '%'",
				validInstrumentsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = NOW()
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				validInstrumentsTableName,
				strings.Join(applyToEach(filterForUpdate(validInstrumentsColumns), func(i int, s string) string {
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
				Name: "UpdateValidInstrumentLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s = sqlc.arg(%s) AND %s IS NULL;`,
				validInstrumentsTableName,
				lastIndexedAtColumn,
				idColumn,
				idColumn,
				archivedAtColumn,
			)),
		},
	}
}
