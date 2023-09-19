package main

import (
	"fmt"
	"strings"

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
	insertColumns := filterForInsert(passwordResetTokensColumns, "redeemed_at")

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "CreatePasswordResetToken",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	sqlc.arg(id),
	sqlc.arg(token),
	NOW() + (30 * interval '1 minutes'),
    sqlc.arg(belongs_to_user)
);`,
				passwordResetTokensTableName,
				strings.Join(insertColumns, ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetPasswordResetToken",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.redeemed_at IS NULL
	AND NOW() < %s.expires_at
	AND %s.token = sqlc.arg(token);`,
				strings.Join(applyToEach(passwordResetTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("password_reset_tokens.%s", s)
				}), ",\n\t"),
				passwordResetTokensTableName,
				passwordResetTokensTableName,
				passwordResetTokensTableName,
				passwordResetTokensTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "RedeemPasswordResetToken",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
    redeemed_at = NOW()
WHERE redeemed_at IS NULL
    AND id = sqlc.arg(id);`,
				passwordResetTokensTableName,
			)),
		},
	}
}
