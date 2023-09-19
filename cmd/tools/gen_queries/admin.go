package main

import (
	"github.com/cristalhq/builq"
)

func buildAdminQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "SetUserAccountStatus",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = NOW(),
	user_account_status = sqlc.arg(user_account_status),
	user_account_status_explanation = sqlc.arg(user_account_status_explanation)
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
	}
}
