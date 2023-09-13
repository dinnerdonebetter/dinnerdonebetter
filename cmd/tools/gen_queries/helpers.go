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
		LineWidth: 80,
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

func buildFilteredColumnConditionsForListQuery(tableName string, hasUpdateColumn bool) string {
	var filteredCountQueryBuilder builq.Builder

	builder := filteredCountQueryBuilder.Addf(`
            %s.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND %s.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
           `,
		tableName,
		tableName,
	)

	if hasUpdateColumn {
		builder.Addf(`AND (
                %s.last_updated_at IS NULL
                OR %s.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
            AND (
                %s.last_updated_at IS NULL
                OR %s.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )`,
			tableName,
			tableName,
			tableName,
			tableName,
		)
	}

	return strings.TrimSpace(buildRawQuery(builder))
}

func buildFilteredColumnCountQuery(tableName string, hasUpdateColumn bool, ownershipColumn string) string {
	var filteredCountQueryBuilder builq.Builder

	builder := filteredCountQueryBuilder.Addf(`(
	    SELECT
	        COUNT(%s.%s)
	    FROM
	        %s
	    WHERE
            %s.archived_at IS NULL
            AND %s.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND %s.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
           `,
		tableName,
		idColumn,
		tableName,
		tableName,
		tableName,
		tableName,
	)

	if hasUpdateColumn {
		builder.Addf(`AND (
                %s.last_updated_at IS NULL
                OR %s.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
            AND (
                %s.last_updated_at IS NULL
                OR %s.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )`,
			tableName,
			tableName,
			tableName,
			tableName,
		)
	}

	if ownershipColumn != "" {
		parts := strings.Split(ownershipColumn, "_")
		builder.Addf(" AND %s = sqlc.arg(%s_id)", ownershipColumn, parts[len(parts)-1])
	}

	builder.Addf(`
	) as filtered_count`)

	return strings.TrimSpace(buildRawQuery(builder))
}

func buildTotalColumnCountQuery(tableName string, addendumConditions ...string) string {
	var totalCountQueryBuilder builq.Builder

	builder := totalCountQueryBuilder.Addf(`(
	    SELECT
	        COUNT(%s.%s)
	    FROM
	        %s
	    WHERE
            %s.archived_at IS NULL
           `,
		tableName,
		idColumn,
		tableName,
		tableName,
	)

	for _, addendum := range addendumConditions {
		builder.Addf(`
            AND %s`, addendum)
	}

	builder.Addf(`
	) as total_count`)

	return strings.TrimSpace(buildRawQuery(builder))
}

func buildCreateQuery(tableName string, columns []string) string {
	var createQueryBuilder builq.Builder

	values := applyToEach(columns, func(s string) string {
		return fmt.Sprintf("sqlc.arg(%s)", s)
	})

	builder := createQueryBuilder.Addf(
		`INSERT INTO %s (%s) VALUES (%s)`,
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

func buildSelectQuery(tableName string, columnNames, joins []string, byID, addAliases bool, conditions ...string) string {
	var selectQueryBuilder builq.Builder

	joinStatements := applyToEach(joins, func(s string) string {
		return fmt.Sprintf("JOIN %s", s)
	})

	and := ""
	if len(conditions) > 0 {
		and = " AND "
	}

	allConditions := strings.Join(conditions, " AND ")

	columns := applyToEach(columnNames, func(s string) string {
		if addAliases {
			parts := strings.Split(s, ".")
			parts[0] = strings.TrimSuffix(parts[0], "s")

			return fmt.Sprintf("%s AS %s", s, strings.Join(parts, "_"))
		}
		return s
	})

	idAddendum := ""
	if byID {
		idAddendum = fmt.Sprintf(" AND %s.%s = sqlc.arg(%s)", tableName, idColumn, idColumn)
	}

	builder := selectQueryBuilder.Addf(
		`SELECT %s FROM %s %s WHERE %s %s %s.archived_at IS NULL%s`,
		strings.Join(columns, ",\n\t"),
		tableName,
		strings.Join(joinStatements, "\n\t"),
		allConditions,
		and,
		tableName,
		idAddendum,
	)

	return buildRawQuery(builder)
}

func buildListQuery(tableName string, columnNames, joins []string, addAliases bool, ownershipColumn string, conditions ...string) string {
	var selectQueryBuilder builq.Builder

	joinStatements := applyToEach(joins, func(s string) string {
		return fmt.Sprintf("JOIN %s", s)
	})

	and := ""
	if len(conditions) > 0 {
		and = " AND "
	}

	conditionsX := append([]string{buildFilteredColumnConditionsForListQuery(tableName, true)}, conditions...)
	allConditions := strings.Join(conditionsX, " AND ")

	columns := append(applyToEach(columnNames, func(s string) string {
		if addAliases {
			parts := strings.Split(s, ".")
			parts[0] = strings.TrimSuffix(parts[0], "s")

			return fmt.Sprintf("%s AS %s", s, strings.Join(parts, "_"))
		}
		return s
	}),
		buildFilteredColumnCountQuery(webhooksTableName, true, ownershipColumn),
		buildTotalColumnCountQuery(webhooksTableName, conditions...),
	)

	columnsToUse := strings.Join(columns, ",\n\t")

	builder := selectQueryBuilder.Addf(
		`SELECT %s FROM %s %s WHERE %s %s %s.archived_at IS NULL OFFSET sqlc.narg(query_offset) LIMIT sqlc.narg(query_limit)`,
		columnsToUse,
		tableName,
		strings.Join(joinStatements, "\n\t"),
		allConditions,
		and,
		tableName,
	)

	return buildRawQuery(builder)
}

func buildExistenceCheckQuery(tableName, addendum string) string {
	var existenceCheckQueryBuilder builq.Builder

	builder := existenceCheckQueryBuilder.Addf(
		`SELECT EXISTS ( SELECT %s.%s FROM %s WHERE %s.archived_at IS NULL AND %s.%s = sqlc.arg(%s) %s )`,
		tableName,
		idColumn,
		tableName,
		tableName,
		tableName,
		idColumn,
		idColumn,
		addendum,
	)

	return buildRawQuery(builder)
}

func buildUpdateQuery(tableName string, columns []string, ownershipColumn string) string {
	var updateQueryBuilder builq.Builder

	columnUpdates := applyToEach(columns, func(s string) string {
		return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
	})

	addendum := ""
	if ownershipColumn != "" {
		addendum = fmt.Sprintf(" AND %s.%s = sqlc.arg(%s_%s)", tableName, ownershipColumn, ownershipColumn, idColumn)
	}

	builder := updateQueryBuilder.Addf(
		`UPDATE %s SET %s = NOW(), %s WHERE archived_at IS NULL %s AND %s = sqlc.arg(%s)`,
		tableName,
		lastUpdatedAtColumn,
		strings.Join(columnUpdates, ",\n\t"),
		addendum,
		idColumn,
		idColumn,
	)

	return buildRawQuery(builder)
}
