package main

import (
	"github.com/Masterminds/squirrel"
)

func buildAdminQueries() []Query {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	updateQuery, _, err := queryBuilder.Update(usersTableName).
		Set("user_account_status", placeholder).
		Set("user_account_status_explanation", placeholder).
		Where(squirrel.Eq{
			archivedAtColumn: nil,
			idColumn:         placeholder,
		}).
		ToSql()

	if err != nil {
		panic(err)
	}

	return []Query{
		{
			Content: updateQuery,
			Annotation: QueryAnnotation{
				Name: "SetUserAccountStatus",
				Type: ExecRows,
			},
		},
	}
}
