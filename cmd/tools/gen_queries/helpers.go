package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cristalhq/builq"
	"github.com/mjibson/sqlfmt"
)

const (
	idColumn                 = "id"
	createdAtColumn          = "created_at"
	lastUpdatedAtColumn      = "last_updated_at"
	archivedAtColumn         = "archived_at"
	lastIndexedAtColumn      = "last_indexed_at"
	belongsToHouseholdColumn = "belongs_to_household"
	belongsToUserColumn      = "belongs_to_user"
)

func applyToEach(x []string, f func(string) string) []string {
	output := []string{}

	for _, v := range x {
		output = append(output, f(v))
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
		if column == archivedAtColumn || column == createdAtColumn || column == lastUpdatedAtColumn {
			continue
		}

		if slices.Contains(exceptions, column) {
			continue
		}

		output = append(output, column)
	}

	return output
}

func fullColumnName(tableName, columnName string) string {
	return fmt.Sprintf("%s.%s", tableName, columnName)
}

func formatQuery(query string) string {
	cfg := tree.PrettyCfg{
		LineWidth: 128,
		Align:     tree.PrettyAlignAndDeindent,
		Simplify:  true,
		TabWidth:  4,
		UseTabs:   true,
		Case:      strings.ToUpper,
		JSONFmt:   true,
	}

	formatted, err := sqlfmt.FmtSQL(cfg, []string{query})
	if err != nil {
		panic(err)
	}

	output := strings.NewReplacer("now()", "NOW()", "count(", "COUNT(").Replace(formatted)

	return regexp.MustCompile(`sqlc\.arg\(\s+(\w+)\s+\)`).ReplaceAllStringFunc(output, func(s string) string {
		replacement := regexp.MustCompile(`\s+`).ReplaceAllString(s, "")

		return replacement
	})
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
		allConditions += fmt.Sprintf("\n\tAND %s", condition)
	}

	return strings.TrimSpace(buildRawQuery((&builq.Builder{}).Addf(`AND %s.%s > COALESCE(sqlc.narg(created_before), (SELECT NOW() - '999 years'::INTERVAL))
    AND %s.%s < COALESCE(sqlc.narg(created_after), (SELECT NOW() + '999 years'::INTERVAL))%s%s`,
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
			AND %s.%s > COALESCE(sqlc.narg(created_before), (SELECT NOW() - '999 years'::INTERVAL))
			AND %s.%s < COALESCE(sqlc.narg(created_after), (SELECT NOW() + '999 years'::INTERVAL))%s%s
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

func buildCreateQuery(tableName string, columns []string) string {
	var createQueryBuilder builq.Builder

	values := applyToEach(columns, func(s string) string {
		return fmt.Sprintf("sqlc.arg(%s)", s)
	})

	builder := createQueryBuilder.Addf(
		`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
		tableName,
		strings.Join(columns, ",\n\t"),
		strings.Join(values, ",\n\t"),
	)

	return buildRawQuery(builder)
}

func buildArchiveQuery(tableName, ownershipColumn string) string {
	var archiveQueryBuilder builq.Builder

	addendum := ""
	if ownershipColumn != "" {
		parts := strings.Split(ownershipColumn, "_")
		addendum = fmt.Sprintf(" AND %s = sqlc.arg(%s_id)", ownershipColumn, parts[len(parts)-1])
	}

	builder := archiveQueryBuilder.Addf(
		`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND %s = sqlc.arg(%s)%s`,
		tableName,
		archivedAtColumn,
		archivedAtColumn,
		idColumn,
		idColumn,
		addendum,
	)

	return buildRawQuery(builder)
}
