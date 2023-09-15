package main

import (
	"github.com/cristalhq/builq"
)

const passwordResetTokensTableName = "password_reset_tokens"

var passwordResetTokensColumns = []string{
	idColumn,
	"token",
	"expires_at",
	"redeemed_at",
	"belongs_to_user",
	createdAtColumn,
	lastUpdatedAtColumn,
}

func buildPasswordResetTokensQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
