package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validPreparationInstrumentsTableName = "valid_preparation_instruments"
)

var validPreparationInstrumentsColumns = []string{
	idColumn,
	notesColumn,
	"valid_preparation_id",
	"valid_instrument_id",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPreparationInstrumentsQueries() []*Query {
	insertColumns := filterForInsert(validPreparationInstrumentsColumns)

	fullSelectColumns := mergeColumns(
		mergeColumns(
			applyToEach(filterFromSlice(validPreparationInstrumentsColumns, "valid_preparation_id", "valid_instrument_id"), func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_instrument_%s", validPreparationInstrumentsTableName, s, s)
			}),
			applyToEach(validInstrumentsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_instrument_%s", validInstrumentsTableName, s, s)
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
				Name: "ArchiveValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
				validPreparationInstrumentsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateValidPreparationInstrument",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				validPreparationInstrumentsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckValidPreparationInstrumentExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
);`,
				validPreparationInstrumentsTableName,
				validPreparationInstrumentsTableName,
				validPreparationInstrumentsTableName,
				archivedAtColumn,
				validPreparationInstrumentsTableName,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForInstrument",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
			%s
	) as filtered_count,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
	) as total_count
FROM valid_preparation_instruments
	JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
	JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
	AND valid_preparation_instruments.valid_instrument_id = sqlc.arg(id)
	AND valid_instruments.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	%s
GROUP BY
	valid_preparation_instruments.id,
	valid_preparations.id,
	valid_instruments.id
ORDER BY valid_preparation_instruments.id
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrumentsForPreparation",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
			%s
	) as filtered_count,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
	) as total_count
FROM valid_preparation_instruments
	JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
	JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
	AND valid_preparation_instruments.valid_preparation_id = sqlc.arg(id)
	AND valid_instruments.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	%s
GROUP BY
	valid_preparation_instruments.id,
	valid_preparations.id,
	valid_instruments.id
ORDER BY valid_preparation_instruments.id
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstruments",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
			%s
	) as filtered_count,
	(
		SELECT
			COUNT(valid_preparation_instruments.id)
		FROM valid_preparation_instruments
            JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
            JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
		WHERE
			valid_preparation_instruments.archived_at IS NULL
			AND valid_instruments.archived_at IS NULL
			AND valid_preparations.archived_at IS NULL
	) as total_count
FROM valid_preparation_instruments
	JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
	JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
	AND valid_instruments.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	%s
GROUP BY
	valid_preparation_instruments.id,
	valid_preparations.id,
	valid_instruments.id
ORDER BY valid_preparation_instruments.id
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				buildFilterConditions(
					validPreparationInstrumentsTableName,
					true,
				),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetValidPreparationInstrument",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM
    valid_preparation_instruments
    JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
    JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
	AND valid_instruments.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_preparation_instruments.id = sqlc.arg(id);`,
				strings.Join(fullSelectColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ValidPreparationInstrumentPairIsValid",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT id
	FROM %s
	WHERE valid_instrument_id = sqlc.arg(valid_instrument_id)
	AND valid_preparation_id = sqlc.arg(valid_preparation_id)
	AND %s IS NULL
);`,
				validPreparationInstrumentsTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateValidPreparationInstrument",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				validPreparationInstrumentsTableName,
				strings.Join(applyToEach(filterForUpdate(validPreparationInstrumentsColumns), func(i int, s string) string {
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
