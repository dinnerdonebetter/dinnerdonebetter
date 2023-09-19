package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

/* #nosec G101 */
const oauth2ClientTokensTableName = "oauth2_client_tokens"

/* #nosec G101 */
var oauth2ClientTokensColumns = []string{
	idColumn,
	"client_id",
	"belongs_to_user",
	"redirect_uri",
	"scope",
	"code",
	"code_challenge",
	"code_challenge_method",
	"code_created_at",
	"code_expires_at",
	"access",
	"access_created_at",
	"access_expires_at",
	"refresh",
	"refresh_created_at",
	"refresh_expires_at",
}

func buildOAuth2ClientTokensQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveOAuth2ClientTokenByAccess",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`DELETE FROM %s WHERE access = $1;`, oauth2ClientTokensTableName)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveOAuth2ClientTokenByCode",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`DELETE FROM %s WHERE code = $1;`, oauth2ClientTokensTableName)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveOAuth2ClientTokenByRefresh",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`DELETE FROM %s WHERE refresh = $1;`, oauth2ClientTokensTableName)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateOAuth2ClientToken",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				oauth2ClientTokensTableName,
				strings.Join(oauth2ClientTokensColumns, ",\n\t"),
				strings.Join(applyToEach(oauth2ClientTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckOAuth2ClientTokenExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.archived_at IS NULL
		AND %s.id = $1
);`,
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2ClientTokenByAccess",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.access = $1;`,
				strings.Join(applyToEach(oauth2ClientTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientTokensTableName, s)
				}), ",\n\t"),
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2ClientTokenByCode",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.code = $1;`,
				strings.Join(applyToEach(oauth2ClientTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientTokensTableName, s)
				}), ",\n\t"),
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetOAuth2ClientTokenByRefresh",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE %s.refresh = $1;`,
				strings.Join(applyToEach(oauth2ClientTokensColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", oauth2ClientTokensTableName, s)
				}), ",\n\t"),
				oauth2ClientTokensTableName,
				oauth2ClientTokensTableName,
			)),
		},
	}
}
