package main

import (
	"github.com/Masterminds/squirrel"
)

func buildAdminQueries() []Query {
	return []Query{
		{
			Annotation: QueryAnnotation{
				Name: "SetUserAccountStatus",
				Type: ExecRowsType,
			},
			Content: buildQuery(
				queryBuilder.Update(usersTableName).
					Set("user_account_status", placeholder).
					Set("user_account_status_explanation", placeholder).
					Where(squirrel.Eq{
						archivedAtColumn: nil,
						idColumn:         placeholder,
					}),
			),
		},
	}
}
