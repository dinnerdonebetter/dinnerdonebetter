package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	idColumn                 = "id"
	createdAtColumn          = "created_at"
	lastUpdatedAtColumn      = "last_updated_at"
	archivedAtColumn         = "archived_at"
	lastIndexedAtColumn      = "last_indexed_at"
	belongsToHouseholdColumn = "belongs_to_household"
	belongsToUserColumn      = "belongs_to_user"

	offsetLimitAddendum = `LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset)`
)

func applyToEach(x []string, f func(int, string) string) []string {
	output := []string{}

	for i, v := range x {
		output = append(output, f(i, v))
	}

	return output
}

func buildRawQuery(builder *builq.Builder) string {
	query, _, err := builder.Build()
	if err != nil {
		panic(err)
	}

	return query
}

func filterForInsert(columns []string, exceptions ...string) []string {
	output := []string{}

	for _, column := range columns {
		if column == archivedAtColumn ||
			column == createdAtColumn ||
			column == lastUpdatedAtColumn ||
			column == lastIndexedAtColumn {
			continue
		}

		if slices.Contains(exceptions, column) {
			continue
		}

		output = append(output, column)
	}

	return output
}

func filterForUpdate(columns []string, exceptions ...string) []string {
	return filterForInsert(columns, append(exceptions, idColumn)...)
}

func fullColumnName(tableName, columnName string) string {
	return fmt.Sprintf("%s.%s", tableName, columnName)
}

func filterFromSlice(slice []string, filtered ...string) []string {
	output := []string{}

	for _, s := range slice {
		if !slices.Contains(filtered, s) {
			output = append(output, s)
		}
	}

	return output
}

func mergeColumns(columns1, columns2 []string, indexToInsertSecondSet uint) []string {
	output := []string{}

	for i, col1 := range columns1 {
		if i == int(indexToInsertSecondSet) {
			output = append(output, columns2...)
		}
		output = append(output, col1)
	}

	return output
}

func buildFilterConditions(tableName string, withUpdateColumn bool, conditions ...string) string {
	updateAddendum := ""
	if withUpdateColumn {
		updateAddendum = fmt.Sprintf("\n\t%s", strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`
	AND (
		%s.%s IS NULL
		OR %s.%s > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		%s.%s IS NULL
		OR %s.%s < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
		`,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
		))))
	}

	allConditions := ""
	for _, condition := range conditions {
		allConditions += fmt.Sprintf("\n\tAND %s", condition)
	}

	return strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`AND %s.%s > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
    AND %s.%s < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))%s%s`,
		tableName,
		createdAtColumn,
		tableName,
		createdAtColumn,
		updateAddendum,
		allConditions,
	)))
}

func buildFilterCountSelect(tableName string, withUpdateColumn bool, conditions ...string) string {
	updateAddendum := ""
	if withUpdateColumn {
		updateAddendum = fmt.Sprintf("\n\t\t\t%s", strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`
			AND (
				%s.%s IS NULL
				OR %s.%s > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				%s.%s IS NULL
				OR %s.%s < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
		`,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
			tableName,
			lastUpdatedAtColumn,
		))))
	}

	allConditions := ""
	for _, condition := range conditions {
		allConditions += fmt.Sprintf("\n\t\t\tAND %s", condition)
	}

	return strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`(
		SELECT COUNT(%s.id)
		FROM %s
		WHERE %s.%s IS NULL
			AND %s.%s > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND %s.%s < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))%s%s
	) AS filtered_count`,
		tableName,
		tableName,
		tableName,
		archivedAtColumn,
		tableName,
		createdAtColumn,
		tableName,
		createdAtColumn,
		updateAddendum,
		allConditions,
	)))
}

func buildTotalCountSelect(tableName string, conditions ...string) string {
	allConditons := ""
	for _, condition := range conditions {
		allConditons += fmt.Sprintf("\n\t\t\tAND %s", condition)
	}

	return strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`(
        SELECT COUNT(%s.id)
        FROM %s
        WHERE %s.%s IS NULL%s
    ) AS total_count`,
		tableName,
		tableName,
		tableName,
		archivedAtColumn,
		allConditons,
	)))
}
