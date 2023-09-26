package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	passwordResetTokensTableName = "password_reset_tokens"
)

var passwordResetTokensColumns = []string{
	idColumn,
	"token",
	"expires_at",
	"redeemed_at",
	belongsToUserColumn,
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
	sqlc.arg(%s),
	sqlc.arg(token),
	%s + (30 * '1 minutes'::INTERVAL),
    sqlc.arg(%s)
);`,
				passwordResetTokensTableName,
				strings.Join(insertColumns, ",\n\t"),
				idColumn,
				currentTimeExpression,
				belongsToUserColumn,
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
	AND %s < %s.expires_at
	AND %s.token = sqlc.arg(token);`,
				strings.Join(applyToEach(passwordResetTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("password_reset_tokens.%s", s)
				}), ",\n\t"),
				passwordResetTokensTableName,
				passwordResetTokensTableName,
				currentTimeExpression,
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
    redeemed_at = %s
WHERE redeemed_at IS NULL
    AND id = sqlc.arg(id);`,
				passwordResetTokensTableName,
				currentTimeExpression,
			)),
		},
	}
}
