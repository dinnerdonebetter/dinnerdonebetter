package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cristalhq/builq"
	"github.com/mjibson/sqlfmt"
)

const (
	placeholder              = "placeholder"
	idColumn                 = "id"
	createdAtColumn          = "created_at"
	lastUpdatedAtColumn      = "last_updated_at"
	archivedAtColumn         = "archived_at"
	belongsToHouseholdColumn = "belongs_to_household"
)

var (
	queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	now          = squirrel.Expr("NOW()")
)

func applyToEach(x []string, f func(string) string) []string {
	output := []string{}

	for _, v := range x {
		output = append(output, f(v))
	}

	return output
}

func buildQuery(builderType any) string {
	switch v := builderType.(type) {
	case squirrel.SelectBuilder:
		query, _, err := v.ToSql()
		if err != nil {
			panic(err)
		}

		return query
	case squirrel.UpdateBuilder:
		query, _, err := v.ToSql()
		if err != nil {
			panic(err)
		}

		return query
	case squirrel.InsertBuilder:
		query, _, err := v.ToSql()
		if err != nil {
			panic(err)
		}

		return query
	case squirrel.DeleteBuilder:
		query, _, err := v.ToSql()
		if err != nil {
			panic(err)
		}

		return query
	default:
		panic("invalid type")
	}
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

func argsForList(columns []string) []any {
	output := []any{}

	for range columns {
		output = append(output, placeholder)
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

	return strings.NewReplacer("now()", "NOW()", "count(", "COUNT(").Replace(formatted)
}

func mergeColumns(columns1, columns2 []string, indexToInsertSecondSet uint) []string {
	output := []string{}

	for i, col1 := range columns1 {
		if i == int(indexToInsertSecondSet+1) {
			for _, col2 := range columns2 {
				output = append(output, col2)
			}
		} else {
			output = append(output, col1)
		}
	}

	return output
}

func buildFilteredColumnCountQuery(tableName string, hasUpdateColumn bool, addendum string) string {
	var filteredCountQueryBuilder builq.Builder

	builder := filteredCountQueryBuilder.Addf(`(
	    SELECT
	        COUNT(%s.id)
	    FROM
	        %s
	    WHERE
            %s.archived_at IS NULL
            AND %s.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND %s.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
           `,
		tableName,
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

	if addendum != "" {
		builder.Addf(`
            AND %s`, addendum)
	}

	builder.Addf(`
	) as filtered_count`)

	return buildRawQuery(builder)
}

func buildTotalColumnCountQuery(tableName, addendum string) string {
	var totalCountQueryBuilder builq.Builder

	builder := totalCountQueryBuilder.Addf(`(
	    SELECT
	        COUNT(%s.id)
	    FROM
	        %s
	    WHERE
            %s.archived_at IS NULL
           `,
		tableName,
		tableName,
		tableName,
	)

	if addendum != "" {
		builder.Addf(`
            AND %s`, addendum)
	}

	builder.Addf(`
	) as total_count`)

	return buildRawQuery(builder)
}

func buildExistenceCheckQuery(tableName, addendum string) string {
	var existenceCheckQueryBuilder builq.Builder

	builder := existenceCheckQueryBuilder.Addf(`SELECT EXISTS ( SELECT %s.id FROM %s WHERE %s.archived_at IS NULL%s )`,
		tableName,
		tableName,
		tableName,
		addendum,
	)

	return buildRawQuery(builder)
}
